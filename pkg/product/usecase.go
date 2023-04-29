package product

import (
	"context"
	"time"

	"github.com/flukis/go-skulatir/pkg/entities"
)

type ProductUsecase interface {
	Store(context.Context, StoreProductParams) (entities.Product, error)
}
type productUsecase struct {
	productRepo ProductRepository
	ctxTimeout  time.Duration
}

func NewProductUsecase(p ProductRepository, to time.Duration) ProductUsecase {
	return &productUsecase{
		productRepo: p,
		ctxTimeout:  to,
	}
}

func (a *productUsecase) Store(c context.Context, arg StoreProductParams) (res entities.Product, err error) {
	ctx, cancel := context.WithTimeout(c, a.ctxTimeout)
	defer cancel()

	existedProductWithSKU, _ := a.productRepo.GetBySKU(ctx, arg.SKU)
	if existedProductWithSKU.SKU == (arg.SKU) {
		return res, entities.ErrConflict
	}

	id, err := a.productRepo.Store(ctx, arg)
	if err != nil {
		return
	}

	res, err = a.productRepo.GetById(ctx, id)
	return
}
