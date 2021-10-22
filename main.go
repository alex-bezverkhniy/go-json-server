package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

type Record map[string]interface{}
type Records []Record
type Collections map[string]Records

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewError(m string) ErrorResponse {
	return ErrorResponse{Status: "error", Message: m}
}

func main() {
	// Get path to .json file
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Printf("Please provide path to .json file as a first argument\n\n")
		fmt.Printf("e.g: go-json-server mydb.json\n\n")
		os.Exit(-1)
	}

	// Read file
	jsonBytes, err := os.ReadFile(argsWithoutProg[0])
	if err != nil {
		fmt.Println("Cannot read file: ", argsWithoutProg[0])
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// Map it to Collections
	var mapData Collections
	if err = json.Unmarshal(jsonBytes, &mapData); err != nil {
		fmt.Println("Cannot unmarshal json: ", argsWithoutProg[0])
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Println("json: ", mapData)

	app := fiber.New()

	// Loop throug all data and and create handlers
	for path, records := range mapData {
		g := app.Group(fmt.Sprintf("/%s", path))

		// Get all
		g.Get("", func(c *fiber.Ctx) error {
			return c.JSON(records)
		})

		// Get by ID
		g.Get("/:id", func(c *fiber.Ctx) error {
			id := c.Params("id")
			if id == "" {
				return c.Status(fiber.StatusBadRequest).JSON(NewError("please provide ID"))
			}

			val := getValueByID(records, id)
			if val == nil {
				return c.Status(fiber.StatusNotFound).JSON(NewError("no data found"))
			}
			return c.JSON(val)
		})
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
		})
	})

	app.Listen(":3000")
}

func getValueByID(m Records, id string) Record {
	for _, record := range m {
		storedId := fmt.Sprintf("%v", record["id"])
		if storedId == id {
			return record
		}
	}
	return nil
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}
