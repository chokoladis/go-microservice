package database

import (
	"context"

	"github.com/fpmoles/go-microservices/internal/models"
)

func (c Client) GetAllVendors(ctx context.Context) ([]models.Vendor, error) {
	var items []models.Vendor
	result := c.DB.WithContext(ctx).
		Find(&items)
	return items, result.Error
}