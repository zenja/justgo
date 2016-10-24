package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/zenja/justgo/model"
	"log"
	"sort"
)

const (
	bucketName string = "justgo"
)

var (
	DB *bolt.DB
)

func encTutorial(t *model.Tutorial) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(t)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decTutorial(b []byte) (*model.Tutorial, error) {
	var t model.Tutorial
	err := json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}
	return &t, err
}

func OpenDB(file string) error {
	log.Printf("Open DB file %s...\n", file)
	if DB != nil {
		panic(fmt.Errorf("OpenDB: there is already an opened DB"))
	}
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		return err
	}
	DB = db
	return initDB()
}

func CloseDB() error {
	return DB.Close()
}

func initDB() error {
	log.Println("Init DB...")
	return DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("init DB: %s", err)
		}
		return nil
	})
}

func AddTutorial(t *model.Tutorial) error {
	return DB.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucketName))
		data, err := encTutorial(t)
		if err != nil {
			return fmt.Errorf("add tutorial: %s", err)
		}
		bkt.Put([]byte(t.Key), data)
		return nil
	})
}

func RemoveTutorial(key string) error {
	return DB.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucketName))
		if err := bkt.Delete([]byte(key)); err != nil {
			return fmt.Errorf("Remove tutorial failed: %s", err)
		}
		return nil
	})
}

func FetchAllKeys() ([]string, error) {
	var keys []string
	DB.View(func(tx *bolt.Tx) error {
		// FIXME Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucketName))
		b.ForEach(func(k, _ []byte) error {
			keys = append(keys, string(k))
			return nil
		})
		return nil
	})
	return keys, nil
}

func FetchTutorial(key string) (*model.Tutorial, error) {
	var tuBytes []byte
	DB.View(func(tx *bolt.Tx) error {
		// FIXME Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucketName))
		tuBytes = b.Get([]byte(key))
		return nil
	})
	if tuBytes == nil {
		return nil, nil
	}
	return decTutorial(tuBytes)
}

func FetchPreNextKey(key string) (string, string, error) {
	indexOf := func(sortedStrs []string, str string) int {
		poIndex := sort.SearchStrings(sortedStrs, str)
		if poIndex < len(sortedStrs) && sortedStrs[poIndex] == str {
			return poIndex
		} else {
			return -1
		}
	}
	keys, err := FetchAllKeys()
	if err != nil {
		return "", "", err
	}
	if len(keys) == 0 {
		return "", "", fmt.Errorf("there is no keys at all")
	}
	sort.Strings(keys)
	index := indexOf(keys, key)
	if index == -1 {
		return "", "", fmt.Errorf("key %s not found", key)
	}
	if len(keys) == 1 {
		return "", "", nil
	}
	if index == 0 {
		return "", keys[index+1], nil
	}
	if index == len(keys)-1 {
		return keys[index-1], "", nil
	}
	return keys[index-1], keys[index+1], nil
}
