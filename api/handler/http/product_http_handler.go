package http_handler

import (
	"fmt"
	"net/http"

	"github.com/flukis/go-skulatir/api/presenter"
	"github.com/flukis/go-skulatir/pkg/entities"
	"github.com/flukis/go-skulatir/pkg/product"
	"github.com/gofiber/fiber/v2"
)

type ProductHttpHandler struct {
	ProductUseCase product.ProductUsecase
}

func NewProductHttpHandler(f *fiber.App, us product.ProductUsecase) {
	handler := &ProductHttpHandler{ProductUseCase: us}

	f.Post("/product", handler.CreateProduct)
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
