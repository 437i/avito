package courier

import (
	mc "avito/internal/model/courier"
	"errors"
	"log"
	"net/http"
)

func (c CreateRequest) toModel() mc.Courier {
	return mc.Courier{
		Name: c.Name,
		Phone: c.Phone,
		Status: mc.CourierStatus(c.Status),
	}
}

func (u UpdateRequest) toModel() mc.Courier {
	return mc.Courier{
		ID: u.ID,
		Name: u.Name,
		Phone: u.Phone,
		Status: mc.CourierStatus(u.Status),
	}
}

func modelToResponse(courier mc.Courier) Courier {
	return Courier{
		ID: courier.ID,
		Name: courier.Name,
		Phone: courier.Phone,
		Status: string(courier.Status),
	}
}

func mapErrorToHTTP(err error) (int, string) {
	switch {
	// service
	case errors.Is(err, mc.ErrEmptyRequest):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, mc.ErrNameEmpty):
		return http.StatusBadRequest, "name can't be empty"
	case errors.Is(err, mc.ErrPhoneEmpty):
		return http.StatusBadRequest, "phone can't be empty"
	case errors.Is(err, mc.ErrStatusEmpty):
		return http.StatusBadRequest, "status can't be empty"
	case errors.Is(err, mc.ErrInvalidPhone):
		return http.StatusBadRequest, "invalid phone number"
	case errors.Is(err, mc.ErrInvalidId):
		return http.StatusBadRequest, "invalid id"
	case errors.Is(err, mc.ErrInvalidStatus):
		return http.StatusBadRequest, "invalid status"
	case errors.Is(err, mc.ErrInvalidName):
		return http.StatusBadRequest, "invalid name"
	// repo
	case errors.Is(err, mc.ErrPhoneExists):
		return http.StatusConflict, "invalid phone"
	case errors.Is(err, mc.ErrIdNotFound):
		return http.StatusNotFound, "id not found"
	}
	log.Printf("internal error: %s\n", err)
	return http.StatusInternalServerError, "unknown error"
}