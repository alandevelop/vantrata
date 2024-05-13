package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
	"ventrata_task/config"
	"ventrata_task/controller"
	libError "ventrata_task/lib/error"
	"ventrata_task/lib/validator"
	"ventrata_task/service"
	"ventrata_task/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	cfg := config.Get()

	store, err := store.New(ctx)
	if err != nil {
		return errors.Wrap(err, "store.New failed")
	}

	serviceManager, err := service.NewManager(store)
	if err != nil {
		return errors.Wrap(err, "manager.New failed")
	}

	availabilityController := controller.NewAvailabilityController(serviceManager)
	bookingController := controller.NewBookingsController(serviceManager)
	productController := controller.NewProductController(serviceManager)

	// Initialize Echo instance
	e := echo.New()
	e.Validator = validator.NewValidator()
	e.HTTPErrorHandler = libError.Error
	// Disable Echo JSON logger in debug mode
	if cfg.LogLevel == "debug" {
		if l, ok := e.Logger.(*echoLog.Logger); ok {
			l.SetHeader("${time_rfc3339} | ${level} | ${short_file}:${line}")
		}
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/products", productController.CreateProduct)
	e.GET("/products", productController.GetProducts)
	e.GET("/products/:id", productController.GetProduct)

	e.POST("/availability/seed", availabilityController.Seed)
	e.POST("/availability", availabilityController.Get)

	e.POST("/bookings", bookingController.Create)
	e.POST("/bookings/:id/confirm", bookingController.Confirm)
	e.GET("/bookings/:id", bookingController.Get)

	// Start server
	s := &http.Server{
		Addr:         cfg.HttpAddr,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))

	return nil
}
