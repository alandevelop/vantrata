package dto

import (
	"github.com/google/uuid"
)

type Product struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty" validate:"required" form:"name"`
	Capacity int       `json:"capacity,omitempty" validate:"required" form:"capacity"`
	Price    int       `json:"price,omitempty" validate:"required" form:"price"`
	Currency *Currency `json:"currency,omitempty"`
}

type ProductWeb struct {
	Id       string `json:"id,omitempty" form:"id"`
	Name     string `json:"name,omitempty" form:"name" validate:"required"`
	Capacity int    `json:"capacity,omitempty" form:"capacity" validate:"required"`
	Price    int    `json:"price,omitempty" form:"price" validate:"required"`
	Currency string `json:"currency,omitempty" form:"price"`
}
