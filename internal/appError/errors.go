package appError

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	ErrID   int
	Message error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

var (
	ErrPostNotAvailable = errors.New("post already not available")
	ErrFileTooLarge     = errors.New("file too large")

	allCustomErrors = []AppError{
		{http.StatusBadRequest, ErrPostNotAvailable},
		{http.StatusRequestEntityTooLarge, ErrFileTooLarge},
	}
)

func CustomError(e error) *AppError {
	for _, customErr := range allCustomErrors {
		if errors.Is(e, customErr.Message) {
			// return &AppError{
			// 	ErrID:   customErr.ErrID,
			// 	Message: customErr.Message,
			// }
			return &customErr
		}
	}
	return &AppError{http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError))}
}
