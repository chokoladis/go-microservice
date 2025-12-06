package database

import (
	"context"
	"errors"

	"github.com/fpmoles/go-microservices/internal/dberrors"
	"github.com/fpmoles/go-microservices/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c Client) GetAllCustomers(ctx context.Context, emailAddress string)  ([]models.Customer, error){
	var customers []models.Customer
	result := c.DB.WithContext(ctx).
		Where(models.Customer{Email: emailAddress}).
		Find(&customers)

	return customers, result.Error 
}

func (c Client) AddCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	customer.CustomerId = uuid.NewString()
	result := c.DB.WithContext(ctx).
		Create(&customer)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}

		return nil, result.Error
	}

	return customer, nil
}

func (c Client) GetCustomerById(ctx context.Context, id string) (*models.Customer, error) {
	customer := &models.Customer{}
	result := c.DB.WithContext(ctx).
		Where(&models.Customer{CustomerId: id}).First(&customer)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{
				Entity: "Customer",
				ID: id,
			}
		}
		return nil, result.Error
	}

	return customer, nil
}

func (c Client) UpdateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	updatedCustomer := &models.Customer{}
	result := c.DB.WithContext(ctx).
		Model(&updatedCustomer).
		Clauses(clause.Returning{}).
		Where(&models.Customer{CustomerId: customer.CustomerId}).
		Updates(models.Customer{
			FirstName: customer.FirstName,
			LastName: customer.LastName,
			Email: customer.Email,
			Phone: customer.Phone,
			Address: customer.Address,
		}).
		Scan(updatedCustomer)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		} else if (errors.Is(result.Error, gorm.ErrRecordNotFound)) {
			return nil, &dberrors.NotFoundError{
				Entity: "Customer",
				ID: customer.CustomerId,
			}
		}
		return nil, result.Error
	}

	return updatedCustomer, nil
}

func (c Client) DeleteCustomer(ctx context.Context, id string) error {
	return c.DB.WithContext(ctx).Delete(&models.Customer{CustomerId: id}).Error
}