package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

const dbFile = "sample.db"
const rootBucket = "root"
const subBucket1 = "sub1"
const subBucket2 = "sub2"

func main() {
	fmt.Println("main start...")

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		// create root bucket.
		root, err := tx.CreateBucketIfNotExists([]byte(rootBucket))
		if err != nil {
			return err
		}

		// create sub bucket 1 and register data.
		sub1, err := root.CreateBucketIfNotExists([]byte(subBucket1))
		if err != nil {
			return err
		}
		err = putDataToSubBucket("value1", sub1)
		if err != nil {
			return err
		}

		// create sub bucket 2 and register data.
		sub2, err := root.CreateBucketIfNotExists([]byte(subBucket2))
		if err != nil {
			return err
		}
		err = putDataToSubBucket("value1", sub2)
		if err != nil {
			return err
		}

		return nil
	})

	err = db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(rootBucket))
		sub1 := root.Bucket([]byte(subBucket1))
		printDataFromSubBucket(sub1)
		sub2 := root.Bucket([]byte(subBucket2))
		printDataFromSubBucket(sub2)
		return nil
	})

	fmt.Println("main end...")
}

func printDataFromSubBucket(sub *bolt.Bucket) {
	c := sub.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		fmt.Printf("key : %v, value : %s\n", k, string(v))
	}
}

func putDataToSubBucket(data string, sub *bolt.Bucket) error {
	id, err := sub.NextSequence()
	if err != nil {
		return err
	}
	err = sub.Put(keybytes(id), []byte(data))
	if err != nil {
		return err
	}
	return nil
}

func keybytes(u uint64) []byte {
	return []byte(string(u))
}
