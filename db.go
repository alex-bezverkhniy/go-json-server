package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

type Record map[string]interface{}
type Records []Record
type DB map[string]Records

func NewDB(fn string) (DB, error) {
	db := DB{}
	err := db.loadJson(fn)
	return db, err
}

func (c DB) Create(path string, r Record) error {
	if r["id"] == nil {
		return errors.New("record should has id")
	}

	if c[path] == nil {
		return errors.New("cannot get records by path")
	}

	c[path] = append(c[path], r)

	return nil
}

func (c DB) Persist(fn string) error {
	bytes, _ := json.MarshalIndent(c, "", "   ")

	return os.WriteFile(fn, bytes, 0644)
}

func (db DB) loadJson(fn string) error {
	// Read file
	jsonBytes, err := os.ReadFile(fn)
	if err != nil || len(jsonBytes) == 0 {
		fmt.Println("Cannot read file: ", fn)
		fmt.Println("Error:", err)
		return err
	}

	// Unmarshal it to DB
	if err = json.Unmarshal(jsonBytes, &db); err != nil {
		fmt.Println("Cannot unmarshal json: ", fn)
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func (r Records) GetByID(id string) Record {
	for _, record := range r {
		storedId := fmt.Sprintf("%v", record["id"])
		if storedId == id {
			return record
		}
	}
	return nil
}

func (r *Record) bind(c *fiber.Ctx) error {

	if err := c.BodyParser(r); err != nil {
		return err
	}

	//fmt.Printf("%v", *r)
	return nil
}
