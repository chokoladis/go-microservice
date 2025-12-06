package server

import (
	"net/http"

	"github.com/fpmoles/go-microservices/internal/dberrors"
	"github.com/fpmoles/go-microservices/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (s *EchoServer) GetAllCustomers(ctx echo.Context) error {
	email := ctx.QueryParam("email")

	customers, err := s.DB.GetAllCustomers(ctx.Request().Context(), email)
	if (err != nil) {
		return ctx.JSON(http.StatusInternalServerError, err)
	} 

	return ctx.JSON(http.StatusOK, customers)
}

func (s *EchoServer) AddCustomer(ctx echo.Context) error {
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	if err := s.validator.Struct(customer); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok  { // err.() ?что за конструкция
			return ctx.JSON(http.StatusBadRequest, echo.Map{
				"error": "Ошибка валидации данных",
				"details": s.handleValidationErrors(validationErrors),
			})
		}

		return ctx.JSON(http.StatusBadRequest, err)
	}

	customer, err := s.DB.AddCustomer(ctx.Request().Context(), customer)
	if err != nil {
		switch err.(type) {
			case *dberrors.ConflictError:
				return ctx.JSON(http.StatusConflict, err)
			default:
				return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, customer)
}