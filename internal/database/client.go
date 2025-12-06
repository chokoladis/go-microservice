package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/fpmoles/go-microservices/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DatabaseClient interface {
	Ready() bool

	GetAllCustomers(ctx context.Context, emailAddress string)  ([]models.Customer, error)
	AddCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) 

	GetAllProducts(ctx context.Context, params url.Values) ([]models.Product, error)
	GetAllVendors(ctx context.Context) ([]models.Vendor, error)
	GetAllServices(ctx context.Context) ([]models.Service, error)
}

type Client struct {
	DB *gorm.DB
}

func NewDatabaseClient() (DatabaseClient, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		"localhost",
		"postgres",
		"root",
		"postgres",
		5432,
		"disable",
	)

	db, error := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "wisdom.",
		},
		NowFunc: func () time.Time {
			return time.Now().UTC()
		},
		QueryFields: true,
	})

	if (error != nil) {
		return nil, error
	}

	client := Client {
		DB: db,
	}

	return client, nil;
}

func (c Client) Ready() bool {
	var ready string
	tx := c.DB.Raw("SELECT 1 as ready").Scan(&ready)
	if tx.Error != nil {
		return false
	}

	if ready == "1" {
		return true
	}

	return false
}