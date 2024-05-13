package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
	"ventrata_task/dto"
	"ventrata_task/lib/types"
	"ventrata_task/service"
)

type AvailabilityController struct {
	services *service.Manager
}

func NewAvailabilityController(services *service.Manager) *AvailabilityController {
	return &AvailabilityController{services}
}

func (ctr *AvailabilityController) Get(ctx echo.Context) error {
	var req dto.AvailabilityFilter

	if v := ctx.FormValue("productId"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid product id")
		}

		req.ProductId = id
	}

	var dateErr error

	if v := ctx.FormValue("localDate"); v != "" {
		req.LocalDate, dateErr = time.Parse("2006-01-02", v)
	}

	if v := ctx.FormValue("localDateFrom"); v != "" {
		req.LocalDateFrom, dateErr = time.Parse("2006-01-02", v)
	}

	if v := ctx.FormValue("localDateTo"); v != "" {
		req.LocalDateTo, dateErr = time.Parse("2006-01-02", v)
	}

	if req.LocalDate.IsZero() && (req.LocalDateFrom.IsZero() || req.LocalDateTo.IsZero()) {
		return echo.NewHTTPError(http.StatusBadRequest, "date fields are required")
	}

	if dateErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid date format")
	}

	data, err := ctr.services.Availability.List(ctx.Request().Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrBadRequest):
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not create booking"))
		}
	}

	dataWeb := make([]dto.AvailabilityWeb, 0, len(data))

	for _, v := range data {
		dataWeb = append(dataWeb, v.ToWeb())
	}

	return ctx.JSON(http.StatusOK, dataWeb)
}

func (ctr *AvailabilityController) Seed(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.FormValue("productId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse id"))
	}

	err = ctr.services.Availability.SeedForYear(ctx.Request().Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, types.ErrBadRequest):
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not create booking"))
		}
	}

	return ctx.NoContent(http.StatusCreated)
}
