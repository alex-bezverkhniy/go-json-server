package handlers

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewError(m string) ErrorResponse {
	return ErrorResponse{Status: "error", Message: m}
}

func Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "OK",
	})
}
