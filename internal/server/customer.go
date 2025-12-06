package server

import (
	"fmt"
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

func (s *EchoServer) GetCustomerById(ctx echo.Context) error {
	id := ctx.Param("id")
	customer, err := s.DB.GetCustomerById(ctx.Request().Context(), id)
	if err != nil {
		switch err.(type) {
			case *dberrors.NotFoundError:
				return ctx.JSON(http.StatusNotFound, err)
			default:
				return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, customer)
}

func (s *EchoServer) UpdateCustomer(ctx echo.Context) error {
	id := ctx.Param("id")
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	} else if id != customer.CustomerId {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("%s not quals %s", id, customer.CustomerId))
	}

	if err := s.validator.Struct(customer); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok  { //todo  how work
			return ctx.JSON(http.StatusBadRequest, echo.Map{
				"error": "Ошибка валидации данных",
				"details": s.handleValidationErrors(validationErrors),
			})
		}

		return ctx.JSON(http.StatusBadRequest, err)
	}

	customer, err := s.DB.UpdateCustomer(ctx.Request().Context(), customer)
	if err != nil {
		switch err.(type) {
			case *dberrors.NotFoundError:
				return ctx.JSON(http.StatusNotFound, err)
			case *dberrors.ConflictError:
				return ctx.JSON(http.StatusConflict, err)
			default:
				return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, customer)
}

func (s *EchoServer) DeleteCustomer(ctx echo.Context) error {
	id := ctx.Param("id")

	err := s.DB.DeleteCustomer(ctx.Request().Context(), id)
	if err != nil {
		switch err.(type) {
			case *dberrors.NotFoundError:
				return ctx.JSON(http.StatusNotFound, err)
			case *dberrors.ConflictError :
				return ctx.JSON(http.StatusConflict, err)
			default:
				return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.NoContent(http.StatusNoContent)
}