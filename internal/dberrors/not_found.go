package dberrors

import "fmt"

type NotFoundError struct {
	ID     string
	Entity string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not found error with ID - %s , Entity - %s", e.ID, e.Entity)
}