package db

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	bolt "go.etcd.io/bbolt"
)

type Operation struct {
	Id      int    `json:"Id"`
	BatchId int    `json:"BatchId"`
	Input   string `json:"Input"`
	Output  string `json:"Output"`
	Undone  bool   `json:"Undone"`
}

type OperationList []Operation

func GetOperationsForBatch(batchId int) (OperationList, error) {
	var err error
	var operations []Operation
	OpenDB()
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

func (ops OperationList) ToTableData(cwd string) pterm.TableData {
	ret := [][]string{}
	ret = append(ret, []string{"ID", "Input", "Output", "Undone"})
	for _, op := range ops {
		ret = append(ret, []string{fmt.Sprintf("%d", op.Id), strings.Replace(op.Input, cwd, "", 1), strings.Replace(op.Output, cwd, "", 1), fmt.Sprintf("%t", op.Undone)})
	}
	return ret
}

func (ops OperationList) Undo(commandType string, cwd string) error {
	var err error
	return db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("operations"))
		for _, op := range ops {
			input := strings.Replace(op.Input, cwd, "", 1)
			output := strings.Replace(op.Output, cwd, "", 1)
			if op.Undone {
				pterm.Info.Printfln("%s already undone", output)
				continue
			}
			if commandType == "move" {
				err = os.Rename(op.Output, op.Input)
				if err != nil {
					pterm.Warning.Printfln("Could not move %s back to %s", output, input)
					return err
				}
				pterm.Success.Printfln("%s → %s", output, input)
			} else {
				err = os.RemoveAll(op.Output)
				if err != nil {
					pterm.Warning.Printfln("Could not delete %s", op.Output)
					return err
				}
				pterm.Success.Printfln("Deleted %s", output)
			}
			op.Undone = true
			buff, _ := json.Marshal(op)
			err = b.Put(itob(op.Id), buff)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
