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
		insertSubBucket := func(sub *bolt.Bucket, data string) error {
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
		fmt.Printf("insert data to %s...\n", subBucket1)
		err = insertSubBucket(sub1, "value1")
		if err != nil {
			return err
		}

		// create sub bucket 2 and register data.
		sub2, err := root.CreateBucketIfNotExists([]byte(subBucket2))
		if err != nil {
			return err
		}
		fmt.Printf("insert data to %s...\n", subBucket2)
		err = insertSubBucket(sub2, "value1")
		if err != nil {
			return err
		}

		return nil
	})

	err = db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(rootBucket))

		readSubBucket := func(subBucket string) {
			sub := root.Bucket([]byte(subBucket))
			c := sub.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("key : %v, value : %s\n", k, string(v))
			}
		}

		fmt.Printf("print data from %s...\n", subBucket1)
		readSubBucket(subBucket1)
		fmt.Printf("print data from %s...\n", subBucket2)
		readSubBucket(subBucket2)

		return nil
	})

	fmt.Println("main end...")
}

func keybytes(u uint64) []byte {
	return []byte(string(u))
}
