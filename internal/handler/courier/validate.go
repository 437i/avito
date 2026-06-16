package courier

import (
	modelCourier "avito/internal/model/courier"
	"strings"
	"unicode"
)

func (r *CreateRequest) validate() error {
	if err := validateName(r.Name); err != nil {
		return err
	}
	if err := validatePhone(r.Phone); err != nil {
		return err
	}
	if err := validateStatus(r.Status); err != nil {
		return err
	}
	return nil
}

func (r *UpdateRequest) validate() error {
	if r.ID <= 0 {
		return modelCourier.ErrInvalidId
	}
	if r.Name == "" && r.Phone == "" && r.Status == "" {
		return modelCourier.ErrEmptyRequest
	}
	if err := validateName(r.Name); r.Name != "" && err != nil {
		return err
	}
	if err := validatePhone(r.Phone); r.Phone != "" && err != nil {
		return err
	}
	if err := validateStatus(r.Status); r.Status != "" && err != nil {
		return err
	}
	return nil
}

func validateName(name string) error {
	if name == "" {
		return modelCourier.ErrNameEmpty
	}
	if len(name) > 100 || strings.ContainsFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r)
	}) {
		return modelCourier.ErrInvalidName
	}
	return nil
}

func validatePhone(phone string) error {
	if phone == "" {
		return modelCourier.ErrPhoneEmpty
	}
	if phone[0] != '+' || len(phone) < 12 {
		return modelCourier.ErrInvalidPhone
	}
	if strings.ContainsFunc(phone[1:], func(r rune) bool {
		return !unicode.IsNumber(r)
	}) {
		return modelCourier.ErrInvalidPhone
	}
	return nil
}

func validateStatus(status string) error {
	switch status {
	case modelCourier.StatusBusy, modelCourier.StatusFree, modelCourier.StatusPaused:
		return nil
	case "":
		return modelCourier.ErrStatusEmpty
	default:
		return modelCourier.ErrInvalidStatus
	}
}