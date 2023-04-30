package product

import (
	"context"
	"time"

	"github.com/flukis/go-skulatir/pkg/entities"
	"github.com/google/uuid"
)

type ProductUsecase interface {
	Store(context.Context, StoreProductParams) (entities.Product, error)
	Get(context.Context, uuid.UUID) (entities.Product, error)
	Delete(context.Context, uuid.UUID) error
	GetBySKU(context.Context, string) (entities.Product, error)
	Fetch(context.Context, string, int) ([]entities.Product, string, error)
	Update(context.Context, *entities.Product) (entities.Product, error)
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

func (a *productUsecase) Get(c context.Context, id uuid.UUID) (res entities.Product, err error) {
	ctx, cancel := context.WithTimeout(c, a.ctxTimeout)
	defer cancel()
	res, err = a.productRepo.GetById(ctx, id)
	return
}

func (a *productUsecase) GetBySKU(c context.Context, sku string) (res entities.Product, err error) {
	ctx, cancel := context.WithTimeout(c, a.ctxTimeout)
	defer cancel()
	res, err = a.productRepo.GetBySKU(ctx, sku)
	return
}

func (a *productUsecase) Fetch(c context.Context, cursor string, limit int) (res []entities.Product, nextCursor string, err error) {
	ctx, cancel := context.WithTimeout(c, a.ctxTimeout)
	defer cancel()
	res, nextCursor, err = a.productRepo.Fetch(ctx, cursor, limit)
	return
}

func (a *productUsecase) Update(c context.Context, arg *entities.Product) (res entities.Product, err error) {
	ctx, cancel := context.WithTimeout(c, a.ctxTimeout)
	defer cancel()

	arg.UpdatedAt = time.Now()
	res, err = a.productRepo.Update(ctx, arg)
	return
}

func (a *productUsecase) Delete(c context.Context, id uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(c, a.ctxTimeout)
	defer cancel()

	err = a.productRepo.Delete(ctx, id)
	return
}
