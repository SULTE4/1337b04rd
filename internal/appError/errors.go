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
	ErrPostNotAvailable        = errors.New("post already not available")
	ErrFileTooLarge            = errors.New("file too large")
	ErrTitleOutOfRange         = errors.New("title must be less than 50 characters")
	ErrTitleTooShort           = errors.New("title must be at least 3 characters")
	ErrContentShouldNotBeEmpty = errors.New("content should not be empty")
	ErrContentOutOfRange       = errors.New("content must be less than 500 characters")

	allCustomErrors = []AppError{
		{http.StatusBadRequest, ErrPostNotAvailable},
		{http.StatusRequestEntityTooLarge, ErrFileTooLarge},
		{http.StatusBadRequest, ErrTitleOutOfRange},
		{http.StatusBadRequest, ErrTitleTooShort},
		{http.StatusBadRequest, ErrContentShouldNotBeEmpty},
		{http.StatusBadRequest, ErrContentOutOfRange},
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
