package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"ventrata_task/dto"
	"ventrata_task/lib/types"
	"ventrata_task/service"
)

type ProductController struct {
	services *service.Manager
}

func NewProductController(services *service.Manager) *ProductController {
	return &ProductController{services}
}

func (ctr *ProductController) GetProducts(ctx echo.Context) error {
	p, err := ctr.services.Product.List(ctx.Request().Context())
	if err != nil {
		switch {
		case errors.Is(err, types.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err)
		case errors.Is(err, types.ErrBadRequest):
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not get booking"))
		}
	}

	capability := ctx.Request().Header.Get("Capability")

	var webRes = make([]dto.ProductWeb, 0, len(p))

	for _, m := range p {

		item := dto.ProductWeb{
			Id:       m.Id.String(),
			Name:     m.Name,
			Capacity: m.Capacity,
		}

		if capability == "price" {
			item.Price = m.Price
			item.Currency = m.Currency.Code.String()
		} else {
			item.Price = 0
			item.Currency = ""
		}

		webRes = append(webRes, item)
	}

	return ctx.JSON(http.StatusOK, webRes)
}

func (ctr *ProductController) GetProduct(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse id"))
	}

	p, err := ctr.services.Product.Get(ctx.Request().Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, types.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err)
		case errors.Is(err, types.ErrBadRequest):
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not get booking"))
		}
	}

	res := dto.ProductWeb{
		Id:       p.Id.String(),
		Name:     p.Name,
		Capacity: p.Capacity,
		Price:    p.Price,
		Currency: p.Currency.Code.String(),
	}

	capability := ctx.Request().Header.Get("Capability")
	if capability == "price" {
		return ctx.JSON(http.StatusOK, res)
	}

	res.Price = 0
	res.Currency = ""
	return ctx.JSON(http.StatusOK, res)
}

func (ctr *ProductController) CreateProduct(ctx echo.Context) error {
	var req dto.Product
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode booking data"))
	}

	err = ctx.Validate(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	_, err = ctr.services.Product.Create(ctx.Request().Context(), &req)

	if err != nil {
		switch {
		case errors.Is(err, types.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err)
		case errors.Is(err, types.ErrBadRequest):
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not create booking"))
		}
	}

	res := dto.ProductWeb{
		Id:       req.Id.String(),
		Name:     req.Name,
		Capacity: req.Capacity,
	}

	return ctx.JSON(http.StatusCreated, res)
}
