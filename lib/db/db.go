package db

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"strings"

	"github.com/jhotmann/go-fileutils-cli/lib/util"
	bolt "go.etcd.io/bbolt"
)

var (
	db     *bolt.DB
	dbPath string = util.GetUserDir() + "/.fileutils/fu.db"
)

type Batch struct {
	Id            int
	CommandType   string
	Command       []string
	CommandString string
	WorkingDir    string
	Undoable      bool
	Undone        bool
}

type Operation struct {
	Id      int
	BatchId int
	Input   string
	Output  string
	Undone  bool
}

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
		buff, err := json.Marshal(batch)
		return b.Put(itob(batch.Id), buff)
	})
	if err != nil {
		panic(err)
	}
	return batch
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

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func ensureFileutilsDir() {
	os.MkdirAll(util.GetUserDir()+"/.fileutils", 0755)
}
