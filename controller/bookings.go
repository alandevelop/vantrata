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

type BookingsController struct {
	services *service.Manager
}

func NewBookingsController(services *service.Manager) *BookingsController {
	return &BookingsController{services}
}

func (ctr *BookingsController) Create(ctx echo.Context) error {
	var booking dto.CreateBookingRequest
	err := ctx.Bind(&booking)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode booking data"))
	}

	err = ctx.Validate(&booking)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	createdB, err := ctr.services.Bookings.CreateBooking(ctx.Request().Context(), &booking)
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

	res := createdB.ToWeb()

	capability := ctx.Request().Header.Get("Capability")
	if capability == "price" {
		return ctx.JSON(http.StatusOK, res)
	}

	res.Price = 0
	res.Currency = ""
	return ctx.JSON(http.StatusCreated, res)
}

func (ctr *BookingsController) Get(ctx echo.Context) error {

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse id"))
	}

	b, err := ctr.services.Bookings.GetBooking(ctx.Request().Context(), id)

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

	res := b.ToWeb()

	capability := ctx.Request().Header.Get("Capability")
	if capability == "price" {
		return ctx.JSON(http.StatusOK, res)
	}

	res.Price = 0
	res.Currency = ""
	return ctx.JSON(http.StatusOK, res)
}

func (ctr *BookingsController) Confirm(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse id"))
	}

	booking, err := ctr.services.Bookings.Confirm(ctx.Request().Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, types.ErrBadRequest):
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not confirm booking"))
		}
	}

	res := booking.ToWeb()

	capability := ctx.Request().Header.Get("Capability")
	if capability == "price" {
		return ctx.JSON(http.StatusOK, res)
	}

	res.Price = 0
	res.Currency = ""
	return ctx.JSON(http.StatusOK, res)
}
