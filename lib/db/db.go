package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jhotmann/go-fileutils-cli/lib/util"
	"github.com/mitchellh/go-homedir"
	"github.com/pterm/pterm"
	bolt "go.etcd.io/bbolt"
)

var (
	db     *bolt.DB
	home   string
	err    error
	dbPath string
)

func init() {
	home, err = homedir.Dir()
	if err != nil {
		home = "/"
	}
	dbPath = home + "/.fileutils/fu.db"
}

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

type Operation struct {
	Id      int    `json:"Id"`
	BatchId int    `json:"BatchId"`
	Input   string `json:"Input"`
	Output  string `json:"Output"`
	Undone  bool   `json:"Undone"`
}

type OperationList []Operation

func NewBatch(commandType string, command []string, workingDir string) Batch {
	ensureFileutilsDir()
	var err error
	db, err = bolt.Open(dbPath, 0755, nil)
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
		batch.CommandString = strings.Join(command, " ")
		batch.WorkingDir = workingDir
		batch.Undoable = util.IndexOf(commandType, []string{"move", "copy", "link-soft", "link-hard"}) > -1
		batch.Undone = false
		batch.Date = time.Now()
		buff, err := json.Marshal(batch)
		return b.Put(itob(batch.Id), buff)
	})
	if err != nil {
		panic(err)
	}
	return batch
}

func GetBatches() (BatchList, error) {
	var err error
	var batches BatchList
	if db == nil {
		db, err = bolt.Open(dbPath, 0755, nil)
		if err != nil {
			db.Close()
			panic(err)
		}
	}
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

func GetOperationsForBatch(batchId int) (OperationList, error) {
	var err error
	var operations []Operation
	if db == nil {
		db, err = bolt.Open(dbPath, 0755, nil)
		if err != nil {
			db.Close()
			panic(err)
		}
	}
	db.View(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("operations"))
		b := tx.Bucket([]byte("operations"))
		b.ForEach(func(k, v []byte) error {
			var op Operation
			err = json.Unmarshal(v, &op)
			if err != nil {
				return err
			}
			if op.BatchId == batchId {
				operations = append(operations, op)
			}
			return nil
		})
		return nil
	})
	return operations, err
}

func WriteOperation(batchId int, input string, output string) error {
	return db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("operations"))
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		op := Operation{
			Id:      int(id),
			BatchId: batchId,
			Input:   input,
			Output:  output,
			Undone:  false,
		}
		buff, err := json.Marshal(op)
		if err != nil {
			return err
		}
		return b.Put(itob(op.Id), buff)
	})
}

func (b Batch) Close() {
	if db == nil {
		return
	}
	db.Close()
}

func (b BatchList) ToTableData() pterm.TableData {
	ret := [][]string{}
	ret = append(ret, []string{"ID", "Date", "Type", "Undone"})
	for _, batch := range b {
		data := []string{fmt.Sprintf("%d", batch.Id), batch.Date.Format("Jan 2, 06 3:04:05"), batch.CommandType, fmt.Sprintf("%t", batch.Undone)}
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
	if max >= last {
		max = last
		pageAfter = false
	}
	return b[min:max], pageBefore, pageAfter
}

func (ops OperationList) ToTableData(cwd string) pterm.TableData {
	ret := [][]string{}
	ret = append(ret, []string{"ID", "Input", "Output", "Undone"})
	for _, op := range ops {
		ret = append(ret, []string{fmt.Sprintf("%d", op.Id), strings.Replace(op.Input, cwd, "", 1), strings.Replace(op.Output, cwd, "", 1), fmt.Sprintf("%t", op.Undone)})
	}
	return ret
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func ensureFileutilsDir() {
	os.MkdirAll(home+"/.fileutils", 0755)
}
