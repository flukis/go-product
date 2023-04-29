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
	f.Get("/product", handler.FetchProduct)
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

func (p ProductHttpHandler) FetchProduct(c *fiber.Ctx) error {
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
