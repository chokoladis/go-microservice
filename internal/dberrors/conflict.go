package dberrors

type ConflictError struct{}

func Error(e ConflictError) string {
	return "error create record"
}