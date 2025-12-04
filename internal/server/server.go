package server

import (
	"log"
	"net/http"

	"github.com/fpmoles/go-microservices/internal/database"
	"github.com/fpmoles/go-microservices/internal/models"
	"github.com/labstack/echo/v4"
)

type Server interface {
	Start() error
	Readiness(ctx echo.Context) error
	Liveness(ctx echo.Context) error
	GetAllCustomers(ctx echo.Context) error
	GetAllProducts(ctx echo.Context) error
}

type EchoServer struct {
	echo *echo.Echo
	DB database.DatabaseClient
}

func NewEchoServer(db database.DatabaseClient) Server {
	server := &EchoServer{
		echo: echo.New(),
		DB: db,
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

	productGroup := s.echo.Group("/products")
	productGroup.GET("", s.GetAllProducts)
}

func (s *EchoServer) Readiness(ctx echo.Context) error {
	ready := s.DB.Ready()
	if (ready) {
		return ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
	}
	return ctx.JSON(http.StatusInternalServerError, models.Health{Status: "Failed"})
}

func (s *EchoServer) Liveness(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
}