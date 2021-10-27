package db

import (
	"encoding/binary"
	"os"
	"time"

	"github.com/mitchellh/go-homedir"
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
	ensureFileutilsDir()
}

func OpenDB() *bolt.DB {
	if db == nil {
		db, err = bolt.Open(dbPath, 0755, &bolt.Options{Timeout: 2 * time.Second})
		if err != nil {
			db.Close()
			panic(err)
		}
	}
	return db
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func ensureFileutilsDir() {
	os.MkdirAll(home+"/.fileutils", 0755)
}
