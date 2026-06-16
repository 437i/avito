package courier

import (
	modelCourier "avito/internal/model/courier"
	"context"
)

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{repo}
}

func (s *Service) GetByID(ctx context.Context, id int) (modelCourier.Courier, error) {
	return s.repo.GetById(ctx, id)
}

func (s *Service) Create(ctx context.Context, courier modelCourier.Courier) (int, error) {
	return s.repo.Create(ctx, courier)
}

func (s *Service) GetAll(ctx context.Context) ([]modelCourier.Courier, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) Update(ctx context.Context, courier modelCourier.Courier) error {
	return s.repo.Update(ctx, courier)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}