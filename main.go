package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "sample.db"
const rootBucket = "users"
const tomBucket = "Tom"
const kenBucket = "Ken"

func main() {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = insert(db)
	if err != nil {
		log.Fatal(err)
	}

	err = reference(db)
	if err != nil {
		log.Fatal(err)
	}
}

func insert(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// create root bucket.
		users, err := tx.CreateBucketIfNotExists([]byte(rootBucket))
		if err != nil {
			return err
		}

		// create nested bucket.
		tom, err := users.CreateBucketIfNotExists([]byte(tomBucket))
		if err != nil {
			return err
		}

		// insert.
		err = tom.Put([]byte("key1"), []byte("tom's todo"))
		if err != nil {
			return err
		}

		// create nested bucket.
		ken, err := users.CreateBucketIfNotExists([]byte(kenBucket))
		if err != nil {
			return err
		}

		// insert.
		err = ken.Put([]byte("key1"), []byte("ken's todo1"))
		if err != nil {
			return err
		}
		err = ken.Put([]byte("key2"), []byte("ken's todo2"))
		if err != nil {
			return err
		}

		return nil
	})
}

func reference(db *bolt.DB) error {
	return db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(rootBucket))

		// select from nested bucket.
		tom := users.Bucket([]byte(tomBucket))
		c := tom.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key:%v, value:%s\n", k, string(v))
		}

		// select from nested bucket.
		ken := users.Bucket([]byte(kenBucket))
		c = ken.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key:%v, value:%s\n", k, string(v))
		}

		return nil
	})
}
