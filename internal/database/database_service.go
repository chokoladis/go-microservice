package database

import (
	"context"

	"github.com/fpmoles/go-microservices/internal/models"
)

func (c Client) GetAllServices(ctx context.Context) ([]models.Service, error) {
	var items []models.Service
	result := c.DB.WithContext(ctx).
		Find(&items)
	return items, result.Error
}