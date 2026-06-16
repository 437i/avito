package courier

import (
	modelCourier "avito/internal/model/courier"
	"context"
)

type repo interface {
	GetById(ctx context.Context, id int) (modelCourier.Courier, error)
	Create(ctx context.Context, courier modelCourier.Courier) (int, error)
	GetAll(ctx context.Context) ([]modelCourier.Courier, error)
	Update(ctx context.Context, courier modelCourier.Courier) error
	Delete(ctx context.Context, id int) error
}
