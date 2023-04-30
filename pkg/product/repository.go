package product

import (
	"context"
	"database/sql"

	"github.com/flukis/go-skulatir/pkg/entities"
	"github.com/flukis/go-skulatir/utils/helpers"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Store(context.Context, StoreProductParams) (uuid.UUID, error)
	GetBySKU(context.Context, string) (entities.Product, error)
	Fetch(context.Context, string, int) ([]entities.Product, string, error)
	GetById(context.Context, uuid.UUID) (entities.Product, error)
	Update(context.Context, *entities.Product) (entities.Product, error)
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
		SELECT
			id, name, sku, description, price, stock, images, created_at, updated_at
		FROM products WHERE sku = $1 LIMIT 1
	`
	err = p.dbconn.GetContext(ctx, &res, queryString, sku)
	if err == sql.ErrNoRows {
		err = entities.ErrNotFound
	}
	return
}

func (p *psqlProductRepository) GetById(ctx context.Context, id uuid.UUID) (res entities.Product, err error) {
	queryString := `
		SELECT
			id, name, sku, description, price, stock, images, created_at, updated_at
		FROM products WHERE id = $1 LIMIT 1
	`
	err = p.dbconn.GetContext(ctx, &res, queryString, id)
	if err == sql.ErrNoRows {
		err = entities.ErrNotFound
	}
	return
}

func (p *psqlProductRepository) Fetch(ctx context.Context, cursor string, num int) (res []entities.Product, nextCursor string, err error) {
	queryString := `
		SELECT
			id, name, sku, description, price, stock, images, created_at, updated_at
		FROM
			products
		WHERE created_at > $1
		ORDER BY created_at
		LIMIT $2
	`
	decodedCursor, err := helpers.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return res, "", entities.ErrBadParamInput
	}
	err = p.dbconn.SelectContext(ctx, &res, queryString, decodedCursor, num)
	if err == sql.ErrNoRows {
		err = entities.ErrNotFound
	}

	if len(res) == int(num) {
		nextCursor = helpers.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}

func (p *psqlProductRepository) Update(ctx context.Context, arg *entities.Product) (res entities.Product, err error) {
	queryString := `
		UPDATE products
		SET
			name = $2,
			sku = $3,
			description = $4,
			price = $5,
			stock = $6,
			images = $7,
			updated_at = $8
		WHERE
			id = $1
		RETURNING
			id, name, sku, description, price, stock, images, created_at, updated_at
	`
	err = p.dbconn.QueryRowContext(ctx, queryString, arg.ID, arg.Name, arg.SKU, arg.Description, arg.Price, arg.Stock, arg.Images, arg.UpdatedAt).Scan(
		&res.ID,
		&res.Name,
		&res.SKU,
		&res.Description,
		&res.Price,
		&res.Stock,
		&res.Images,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	return
}
