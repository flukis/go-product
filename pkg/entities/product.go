package entities

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id" db:"id"`
	SKU         string    `json:"sku" db:"sku"`
	Price       int       `json:"price" db:"price"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Images      string    `json:"images" db:"images"`
	Stock       int       `json:"stock" db:"stock"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
