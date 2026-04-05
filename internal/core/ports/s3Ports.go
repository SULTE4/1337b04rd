package ports

import "1337b04rd/internal/core/domain"

type S3Repository interface {
	UploadObject(isPostImg bool, f domain.UploadFile) (string, error)
	GetObject(key string) (string, error)
}
