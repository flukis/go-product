package http_handler

import (
	"fmt"
	"net/http"

	"github.com/flukis/go-skulatir/api/presenter"
	"github.com/flukis/go-skulatir/pkg/entities"
	"github.com/flukis/go-skulatir/pkg/product"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHttpHandler struct {
	ProductUseCase product.ProductUsecase
}

func NewProductHttpHandler(f *fiber.App, us product.ProductUsecase) {
	handler := &ProductHttpHandler{ProductUseCase: us}

	f.Post("/product", handler.CreateProduct)
	f.Get("/product/:id", handler.GetProduct)
	f.Get("/products", handler.FetchProduct)
	f.Put("/product", handler.UpdateProduct)
	f.Delete("/product/:id", handler.DeleteProduct)
}

func (p *ProductHttpHandler) CreateProduct(c *fiber.Ctx) error {
	var reqBody product.StoreProductParams
	err := c.BodyParser(&reqBody)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(presenter.ProductErrorResponse(err))
	}

	if reqBody.Name == "" || reqBody.Price == 0 || reqBody.SKU == "" {
		newErr := fmt.Errorf("product name, SKU, and price is required")
		return c.Status(http.StatusBadRequest).JSON(presenter.ProductErrorResponse(newErr))
	}

	res, err := p.ProductUseCase.Store(c.Context(), reqBody)
	if err != nil {
		if err == entities.ErrNotFound {
			return c.Status(http.StatusNotFound).JSON(presenter.ProductErrorResponse(err))
		}
		if err == entities.ErrConflict {
			return c.Status(http.StatusConflict).JSON(presenter.ProductErrorResponse(err))
		}
		return c.Status(http.StatusInternalServerError).JSON(presenter.ProductErrorResponse(err))
	}

	return c.Status(http.StatusOK).JSON(presenter.ProductSuccessResponse(&res))
}

func (p *ProductHttpHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(presenter.ProductErrorResponse(err))
	}

	res, err := p.ProductUseCase.Get(c.Context(), uid)
	if err != nil {
		if err == entities.ErrNotFound {
			return c.Status(http.StatusNotFound).JSON(presenter.ProductErrorResponse(err))
		}
		return c.Status(http.StatusInternalServerError).JSON(presenter.ProductErrorResponse(err))
	}
	return c.Status(http.StatusOK).JSON(presenter.ProductSuccessResponse(&res))
}

func (p *ProductHttpHandler) FetchProduct(c *fiber.Ctx) error {
	var reqBody presenter.Pagination
	err := c.BodyParser(&reqBody)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(presenter.ProductErrorResponse(err))
	}

	if reqBody.Limit < 10 {
		reqBody.Limit = 10
	}

	res, nextCursor, err := p.ProductUseCase.Fetch(c.Context(), reqBody.Cursor, reqBody.Limit)
	if err != nil {
		if err == entities.ErrNotFound {
			return c.Status(http.StatusNotFound).JSON(presenter.ProductErrorResponse(err))
		}
		return c.Status(http.StatusInternalServerError).JSON(presenter.ProductErrorResponse(err))
	}

	pagination := presenter.Pagination{
		Limit:  reqBody.Limit,
		Cursor: nextCursor,
	}

	return c.Status(http.StatusOK).JSON(presenter.ProductsSuccessResponse(&res, pagination))
}

type updateProductParams struct {
	ID          string `json:"id"`
	SKU         string `json:"sku"`
	Price       int    `json:"price"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Images      string `json:"images"`
	Stock       int    `json:"stock"`
}

func (p *ProductHttpHandler) UpdateProduct(c *fiber.Ctx) error {
	var reqBody updateProductParams
	err := c.BodyParser(&reqBody)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(presenter.ProductErrorResponse(err))
	}

	uid, err := uuid.Parse(reqBody.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(presenter.ProductErrorResponse(err))
	}

	currentProduct, err := p.ProductUseCase.Get(c.Context(), uid)
	if err != nil {
		if err == entities.ErrNotFound {
			return c.Status(http.StatusNotFound).JSON(presenter.ProductErrorResponse(err))
		}
		return c.Status(http.StatusInternalServerError).JSON(presenter.ProductErrorResponse(err))
	}

	if currentProduct.SKU != reqBody.SKU {
		_, err := p.ProductUseCase.Get(c.Context(), uid)
		if err != entities.ErrNotFound {
			err = fmt.Errorf("SKU is already taken")
			return c.Status(http.StatusConflict).JSON(presenter.ProductErrorResponse(err))
		}
	}

	arg := entities.Product{
		ID:          uid,
		Name:        reqBody.Name,
		SKU:         reqBody.SKU,
		Description: reqBody.Description,
		Stock:       reqBody.Stock,
		Price:       reqBody.Price,
		Images:      reqBody.Images,
	}

	res, err := p.ProductUseCase.Update(c.Context(), &arg)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(presenter.ProductErrorResponse(err))
	}
	return c.Status(http.StatusOK).JSON(presenter.ProductSuccessResponse(&res))
}

func (p *ProductHttpHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(presenter.ProductErrorResponse(err))
	}

	err = p.ProductUseCase.Delete(c.Context(), uid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(presenter.ProductErrorResponse(err))
	}
	return c.Status(http.StatusOK).JSON(presenter.ProductSuccessResponse(nil))
}
