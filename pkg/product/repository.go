package product

import (
	"context"

	"github.com/flukis/go-skulatir/pkg/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Store(context.Context, StoreProductParams) (uuid.UUID, error)
	GetBySKU(context.Context, string) (entities.Product, error)
	GetById(context.Context, uuid.UUID) (entities.Product, error)
}

type psqlProductRepository struct {
	dbconn *sqlx.DB
}

func NewPsqlProductRepository(dbconn *sqlx.DB) ProductRepository {
	return &psqlProductRepository{dbconn: dbconn}
}

// Store Product Repository
type StoreProductParams struct {
	Name   string `json:"name"`
	SKU    string `json:"sku"`
	Desc   string `json:"description"`
	Price  int    `json:"price"`
	Stock  int    `json:"stock"`
	Images string `json:"images"`
}

func (p *psqlProductRepository) Store(ctx context.Context, arg StoreProductParams) (id uuid.UUID, err error) {
	queryString := `
		INSERT INTO products (
			name,
			sku,
			description,
			price,
			stock,
			images
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		) RETURNING id
	`
	err = p.dbconn.QueryRowContext(ctx, queryString, arg.Name, arg.SKU, arg.Desc, arg.Price, arg.Stock, arg.Images).Scan(&id)
	return
}

func (p *psqlProductRepository) GetBySKU(ctx context.Context, sku string) (res entities.Product, err error) {
	queryString := `
		SELECT id, name, sku, description, price, stock, images, created_at, updated_at FROM products WHERE sku = $1 LIMIT 1
	`
	err = p.dbconn.GetContext(ctx, &res, queryString, sku)
	return
}

func (p *psqlProductRepository) GetById(ctx context.Context, id uuid.UUID) (res entities.Product, err error) {
	queryString := `
	SELECT id, name, sku, description, price, stock, images, created_at, updated_at FROM products WHERE id = $1 LIMIT 1
	`
	err = p.dbconn.GetContext(ctx, &res, queryString, id)
	return
}
