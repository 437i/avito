package courier

import (
	modelCourier "avito/internal/model/courier"
	"context"
)

type service interface {
	GetByID(context.Context, int) (modelCourier.Courier, error)
	Create(context.Context, modelCourier.Courier) (int, error)
	GetAll(context.Context) ([]modelCourier.Courier, error)
	Update(context.Context, modelCourier.Courier) error
	Delete(context.Context, int) error
}