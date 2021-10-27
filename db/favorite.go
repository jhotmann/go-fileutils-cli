package db

import (
	"encoding/json"
	"errors"

	bolt "go.etcd.io/bbolt"
)

type Favorite struct {
	Id            int      `json:"Id"`
	Name          string   `json:"Name"`
	Command       []string `json:"Command"`
	CommandType   string   `json:"CommandType"`
	CommandString string   `json:"CommandString"`
}

type FavoriteList []Favorite

func NewFavorite(name string, commandType string, command string) (Favorite, error) {
	OpenDB()
	if err != nil {
		return Favorite{}, err
	}
	favorite := GetFavoriteByName(name)
	if favorite.Id != 0 {
		return favorite, errors.New("Favorite with name " + name + " already exists")
	}
	err = db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("favorites"))
		b := tx.Bucket([]byte("favorites"))
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		favorite.Id = int(id)
		favorite.Name = name
		favorite.CommandType = commandType
		favorite.CommandString = command
		buff, err := json.Marshal(favorite)
		return b.Put(itob(favorite.Id), buff)
	})
	return favorite, err
}

func GetFavoriteByName(name string) Favorite {
	OpenDB()
	if err != nil {
		panic(err)
	}
	ret := Favorite{}
	err = db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("favorites"))
		b := tx.Bucket([]byte("favorites"))
		b.ForEach(func(k, v []byte) error {
			if ret.Id != 0 {
				fav := Favorite{}
				err = json.Unmarshal(v, &fav)
				if err != nil {
					return err
				}
				if fav.Name == name {
					ret = fav
					return nil
				}
			}
			return nil
		})
		return nil
	})
	return ret
}
