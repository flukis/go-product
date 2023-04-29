package presenter

import (
	"github.com/flukis/go-skulatir/pkg/entities"
	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	Cursor string `json:"next"`
	Limit  int    `json:"limit"`
}

// ProductuccessResponse is the singular SuccessResponse that will be passed in the response by
// Handler
func ProductSuccessResponse(data *entities.Product) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data":   data,
		"error":  nil,
	}
}

// ProductSuccessResponse is the list SuccessResponse that will be passed in the response by Handler
func ProductsSuccessResponse(data *[]entities.Product, meta Pagination) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data":   data,
		"error":  nil,
		"meta":   meta,
	}
}

// ProductErrorResponse is the ErrorResponse that will be passed in the response by Handler
func ProductErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"data":   "",
		"error":  err.Error(),
	}
}
