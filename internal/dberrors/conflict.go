package dberrors

type ConflictError struct{}

func (e *ConflictError) Error() string {
	return "error create record"
}
