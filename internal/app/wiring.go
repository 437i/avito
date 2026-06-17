package app

import (
	hc "avito/internal/handler/courier"
	rc "avito/internal/repository/courier"
	sc "avito/internal/service/courier"
)

func NewService(repo *rc.Repository) *sc.Service {
    return sc.NewService(repo)
}

func NewHandler(svc *sc.Service) *hc.Handler {
    return hc.NewHandler(svc)
}