package database

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

type Task struct {
	ID       int
	Task     []byte
	Complete bool
}

const fp = "/.task/tasks.db"

func GetTasks() []Task {
	db := createDBInstance()
	defer db.Close()

	tasks := []Task{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			t := Task{}
			json.Unmarshal(v, &t)
			tasks = append(tasks, t)
		}

		return nil
	})

	return tasks
}

func GetUncompletedTasks() []Task {
	db := createDBInstance()
	defer db.Close()

	tasks := []Task{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			t := Task{}
			json.Unmarshal(v, &t)

			if !t.Complete {
				tasks = append(tasks, t)
			}
		}

		return nil
	})

	return tasks
}

func MarkTaskAsCompleted(id int) {
	db := createDBInstance()
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
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
		b := tx.Bucket([]byte("Tasks"))

		id, _ := b.NextSequence()

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
	createDbIfDoesNotExist()
	db, err := bolt.Open(fp, 0600, nil)
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

func createDbIfDoesNotExist() {
	if _, err := os.Stat(fp); err == nil {
		return
	} else if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(fp), 0770); err != nil {
			fmt.Println("There was an error creating the local tasks record file.")
			os.Exit(1)
		}
		os.Create(fp)

		db, err := bolt.Open(fp, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
		return
	} else {
		fmt.Println(err)
		fmt.Println("There was an error creating the local tasks record.")
		os.Exit(1)
	}
}
