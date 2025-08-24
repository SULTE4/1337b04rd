package ports

import "1337b04rd/internal/adapters/s3"

type S3Repository interface {
	UploadObject(isPostImg bool, f s3.FileType) (string, error)
	GetObject(key string) (string, error)
}
