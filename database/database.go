package database

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type Task struct {
	ID       int
	Task     []byte
	Complete bool
}

func GetTasks() {
	db := createDBInstance()
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%d, value%s\n", btoi(k), v)
		}

		return nil
	})
}

func GetUncompletedTasks() {
	db := createDBInstance()
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			t := Task{}
			json.Unmarshal(v, &t)

			fmt.Println(t.Complete)

			if t.Complete {
				fmt.Printf("key=%d, value%s\n", btoi(k), v)
			}
		}

		return nil
	})
}

func MarkTaskAsCompleted(id int) {
	db := createDBInstance()
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte("Tasks"))
		c := b.Cursor()
		_, v := c.Seek(itob(id))

		t := Task{}
		json.Unmarshal(v, &t)

		t.Complete = true

		buf, _ := json.Marshal(t)

		return b.Put(itob(t.ID), buf)
	})
}

func AddTask(task string) {
	db := createDBInstance()
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		if err != nil {
			log.Fatal(err)
		}

		id, _ := b.NextSequence()

		fmt.Println(id)

		t := Task{
			ID:       int(id),
			Task:     []byte(task),
			Complete: false,
		}

		buf, _ := json.Marshal(t)

		return b.Put(itob(t.ID), buf)
	})
}

func createDBInstance() *bolt.DB {
	db, err := bolt.Open("./tasks.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// itob returns an 8-byte big endian representation of v.
func btoi(b []byte) int {
	i := binary.BigEndian.Uint64(b)
	return int(i)
}