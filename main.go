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

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Config struct {
	FileName string
}

func NewError(m string) ErrorResponse {
	return ErrorResponse{Status: "error", Message: m}
}

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

var config = Config{}

func main() {
	// Get path to .json file
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Printf("Please provide path to .json file as a first argument\n\n")
		fmt.Printf("e.g: go-json-server mydb.json\n\n")
		os.Exit(-1)
	}

	config.FileName = argsWithoutProg[0]

	db, err := NewDB(config.FileName)
	if err != nil {
		os.Exit(-1)
	}

	fmt.Println("json: ", db)

	app := fiber.New()

	// Loop throug all data and and create handlers
	for path, _ := range db {
		g := app.Group(fmt.Sprintf("/%s", path))

		// Get all
		g.Get("", func(c *fiber.Ctx) error {
			return c.JSON(db[path])
		})

		// Get by ID
		g.Get("/:id", func(c *fiber.Ctx) error {
			id := c.Params("id")
			if id == "" {
				return c.Status(fiber.StatusBadRequest).JSON(NewError("please provide ID"))
			}

			val := db[path].GetByID(id) //getValueByID(records, id)
			if val == nil {
				return c.Status(fiber.StatusNotFound).JSON(NewError("no data found"))
			}
			return c.JSON(val)
		})

		// Create new record
		g.Post("", func(c *fiber.Ctx) error {
			r := Record{}
			if err := r.bind(c); err != nil {
				fmt.Println("ERROR", err.Error())
				return c.Status(fiber.StatusUnprocessableEntity).JSON(NewError(err.Error()))
			}

			db.Create(path, r)
			err := db.Persist(config.FileName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(NewError(err.Error()))
			}

			return c.Status(fiber.StatusCreated).SendString("Created")
		})
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
		})
	})

	app.Listen(":3000")
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}
