package main

import (
	"log"
	"time"

	httpHandler "github.com/flukis/go-skulatir/api/handler/http"
	"github.com/flukis/go-skulatir/pkg/product"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS "products" (
		"id" uuid PRIMARY KEY DEFAULT uuid_generate_v4 (), 
		"name" varchar NOT NULL,
		"description" varchar NULL,
		"images" varchar NOT NULL,
		"sku" varchar unique NOT NULL,
		"price" bigint NOT NULL,
		"stock" int NOT NULL,
		"created_at" timestamptz NOT NULL DEFAULT (now()),
		"updated_at" timestamptz NOT NULL DEFAULT (now())
	);

	CREATE INDEX IF NOT EXISTS idx_product_sku ON products(sku);
`

func main() {
	db, err := sqlx.Connect("postgres", "user=root dbname=ecommerce password=wap12345 sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	db.MustExec(schema)

	f := fiber.New()

	productRepo := product.NewPsqlProductRepository(db)
	to := time.Duration(10) * time.Second

	productUsecase := product.NewProductUsecase(productRepo, to)
	httpHandler.NewProductHttpHandler(f, productUsecase)

	if err := f.Listen("0.0.0.0:5000"); err != nil {
		log.Fatalf("Server is not running, %v", err)
	}
}
