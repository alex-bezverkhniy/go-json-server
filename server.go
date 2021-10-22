package main

import (
	"fmt"
	"go-json-server/handlers"

	"github.com/gofiber/fiber/v2"
)

type JsonServer struct {
	config Config
	app    *fiber.App
}

const DefaultPort = ":3000"

// Creates new JsonServer
func NewJsonServer(config Config) (*JsonServer, error) {

	db, err := NewDB(config.FileName)
	if err != nil {
		return nil, err
	}

	app := fiber.New()

	// Loop throug all data and and create handlers
	for path := range db {
		g := app.Group(fmt.Sprintf("/%s", path))

		// Get all
		g.Get("", func(c *fiber.Ctx) error {
			return c.JSON(db[path])
		})

		// Get by ID
		g.Get("/:id", func(c *fiber.Ctx) error {
			id := c.Params("id")
			if id == "" {
				return c.Status(fiber.StatusBadRequest).JSON(handlers.NewError("please provide ID"))
			}

			val := db[path].GetByID(id) //getValueByID(records, id)
			if val == nil {
				return c.Status(fiber.StatusNotFound).JSON(handlers.NewError("no data found"))
			}
			return c.JSON(val)
		})

		// Create new record
		g.Post("", func(c *fiber.Ctx) error {
			r := Record{}
			if err := r.bind(c); err != nil {
				fmt.Println("ERROR", err.Error())
				return c.Status(fiber.StatusUnprocessableEntity).JSON(handlers.NewError(err.Error()))
			}

			db.Create(path, r)
			err := db.Persist(config.FileName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(handlers.NewError(err.Error()))
			}

			return c.Status(fiber.StatusCreated).SendString("Created")
		})
	}

	app.Get("/health", handlers.Health)

	server := JsonServer{
		config: config,
		app:    app,
	}

	return &server, nil
}

func (s *JsonServer) Start() {
	if s.config.Port == "" {
		s.config.Port = DefaultPort
	}
	s.app.Listen(s.config.Port)
}
