package database

import (
	"context"

	"github.com/fpmoles/go-microservices/internal/models"
)

// todo разница echo.Context и context.Context
func (c Client) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	result := c.DB.WithContext(ctx).
		Find(&products)
	return products, result.Error
}