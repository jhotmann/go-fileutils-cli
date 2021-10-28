package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jhotmann/go-fileutils-cli/util"
	"github.com/pterm/pterm"
	bolt "go.etcd.io/bbolt"
)

type Batch struct {
	Id            int       `json:"Id"`
	CommandType   string    `json:"CommandType"`
	Command       []string  `json:"Command"`
	CommandString string    `json:"CommandString"`
	WorkingDir    string    `json:"WorkingDir"`
	Undoable      bool      `json:"Undoable"`
	Undone        bool      `json:"Undone"`
	Date          time.Time `json:"Date"`
}

type BatchList []Batch

func (b BatchList) Reverse() BatchList {
	ret := BatchList{}
	for _, batch := range b {
		ret = append([]Batch{batch}, ret...)
	}
	return ret
}

func NewBatch(commandType string, command []string, workingDir string) Batch {
	OpenDB()
	if err != nil {
		panic(err)
	}
	batch := Batch{}
	err = db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("batches"))
		tx.CreateBucketIfNotExists([]byte("operations"))
		b := tx.Bucket([]byte("batches"))
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		batch.Id = int(id)
		batch.CommandType = commandType
		batch.Command = command
		batch.CommandString = formatCommandString(command)
		batch.WorkingDir = workingDir
		batch.Undoable = util.IndexOf(commandType, []string{"move", "copy", "link-soft", "link-hard"}) > -1
		batch.Undone = false
		batch.Date = time.Now()
		buff, err := json.Marshal(batch)
		return b.Put(itob(batch.Id), buff)
	})
	if err != nil {
		db.Close()
		panic(err)
	}
	return batch
}

func formatCommandString(command []string) string {
	ret := ""
	for i, token := range command {
		space := ""
		if i > 0 {
			space = " "
		}
		quotes := ""
		if strings.ContainsAny(token, " |'?%*+[];&<>!$") {
			quotes = `"`
		}
		ret = fmt.Sprintf("%s%s%s%s%s", ret, space, quotes, token, quotes)
	}
	return ret
}

func GetBatches() (BatchList, error) {
	var err error
	var batches BatchList
	OpenDB()
	err = db.View(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("batches"))
		b := tx.Bucket([]byte("batches"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var batch Batch
			err = json.Unmarshal(v, &batch)
			if err == nil {
				batches = append(batches, batch)
			}
		}
		return err
	})
	return batches, err
}

func GetLastNonUndone() (Batch, error) {
	var err error
	var batch Batch
	OpenDB()
	err = db.View(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("batches"))
		b := tx.Bucket([]byte("batches"))
		c := b.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			err = json.Unmarshal(v, &batch)
			if err == nil && !batch.Undone && batch.Undoable {
				return nil
			}
		}
		return errors.New("No undoable batches found")
	})
	return batch, err
}

func (batch Batch) Undo() error {
	if batch.Undone {
		pterm.Warning.Println("Batch already undone")
		return errors.New("Batch already undone")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("batches"))
		batch.Undone = true
		buff, _ := json.Marshal(batch)
		return b.Put(itob(batch.Id), buff)
	})
	if err != nil {
		return err
	}
	operations, err := GetOperationsForBatch(batch.Id)
	if err != nil {
		return err
	}
	return operations.Undo(batch.CommandType, batch.WorkingDir)
}

func (b BatchList) ToTableData() pterm.TableData {
	ret := [][]string{}
	ret = append(ret, []string{"ID", "Date", "Type", "Undone"})
	for _, batch := range b {
		data := []string{fmt.Sprintf("%d", batch.Id), batch.Date.Format("Jan 2, 2006 15:04:05"), batch.CommandType, fmt.Sprintf("%t", batch.Undone)}
		ret = append(ret, data)
	}
	return ret
}

func (b BatchList) GetPage(page int, batchesPerPage int) (batchList BatchList, pageBefore bool, pageAfter bool) {
	batchList = BatchList{}
	pageBefore = page > 1
	pageAfter = true
	max := batchesPerPage * page
	min := max - batchesPerPage
	last := len(b) - 1
	if max > last {
		max = last + 1
		pageAfter = false
	}
	return b[min:max], pageBefore, pageAfter
}

func (b Batch) Close() {
	if db == nil {
		return
	}
	db.Close()
}

func (b BatchList) Close() {
	if db == nil {
		return
	}
	db.Close()
}
