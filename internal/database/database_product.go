package database

import (
	"context"
	"net/url"

	"github.com/fpmoles/go-microservices/internal/models"
)

func (c Client) GetAllProducts(ctx context.Context, params url.Values) ([]models.Product, error) {

	var products []models.Product

	query := c.DB.WithContext(ctx)

	if vendorId := params.Get("vendorId"); vendorId != "" {
		query = query.Where("vendor_id = ?",vendorId)
	}

	if min_price := params.Get("min_price"); min_price != "" {
		query = query.Where("price > ?",min_price)
	}

	result := query.Find(&products)
	return products, result.Error
}
