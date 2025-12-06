package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fpmoles/go-microservices/internal/database"
	"github.com/fpmoles/go-microservices/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Server interface {
	Start() error
	Readiness(ctx echo.Context) error
	Liveness(ctx echo.Context) error

	GetAllCustomers(ctx echo.Context) error
	GetCustomerById(ctx echo.Context) error 
	AddCustomer(ctx echo.Context) error

	GetAllProducts(ctx echo.Context) error
	GetAllVendors(ctx echo.Context) error
	GetAllServices(ctx echo.Context) error
}

type EchoServer struct {
	echo *echo.Echo
	DB   database.DatabaseClient
	validator *validator.Validate
}

var validate *validator.Validate

func init() {
    validate = validator.New()
}

func NewEchoServer(db database.DatabaseClient) Server {
	server := &EchoServer{
		echo: echo.New(),
		DB:   db,
		validator: validate,
	}
	server.registerRoutes()
	return server
}

func (s *EchoServer) Start() error {
	if err := s.echo.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server shotdown, reason - %s", err)
		return err
	}

	return nil
}
func (s *EchoServer) registerRoutes() {
	s.echo.GET("/readness", s.Readiness)
	s.echo.GET("/liveness", s.Liveness)

	cg := s.echo.Group("/customers")
	cg.GET("", s.GetAllCustomers)
	cg.GET("/:id", s.GetCustomerById)
	cg.POST("", s.AddCustomer)
	cg.PUT("/:id", s.UpdateCustomer)
	cg.DELETE("/:id", s.DeleteCustomer)

	productGroup := s.echo.Group("/products")
	productGroup.GET("", s.GetAllProducts)

	vendorsGroup := s.echo.Group("/vendors")
	vendorsGroup.GET("", s.GetAllVendors)

	servicesGroup := s.echo.Group("/services")
	servicesGroup.GET("", s.GetAllServices)
}

func (s *EchoServer) Readiness(ctx echo.Context) error {
	ready := s.DB.Ready()
	if ready {
		return ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
	}
	return ctx.JSON(http.StatusInternalServerError, models.Health{Status: "Failed"})
}

func (s *EchoServer) Liveness(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
}


var validationTranslations = map[string]string{
    "required": "обязательно для заполнения",
	"min":      "маленькая длинна",
    "max":      "превышена максимальная длинна",
    "email":    "имеет неверный формат почты",
    "e164":     "должен быть в формате E.164 (например, +79123456789)",
	// todo products, vendors, services
}

var fieldNames = map[string]string{
    "FirstName": "Имя",
    "LastName":  "Фамилия",
    "Email":     "Email-адрес",
    "Phone":     "Номер телефона",
    "Address":   "Адрес",
}

func (s *EchoServer) handleValidationErrors(errs validator.ValidationErrors) []echo.Map {
	var results []echo.Map
	for _, err := range errs {
		fieldName := err.Field()
		friendlyFieldName := fieldNames[fieldName]
		if friendlyFieldName == "" {
			friendlyFieldName = fieldName
		}

		tag := err.Tag()
		message := validationTranslations[tag]
		if message == "" {
			message = fmt.Sprintf("не проходит проверку '%s'", tag)
		}

		finalMessage := fmt.Sprintf("Поле '%s' %s", friendlyFieldName, message)
        
		results = append(results, echo.Map{
			"field":   fieldName,
			"message": finalMessage,
			"rule":    tag,
		})
	}
	return results
}