package s3

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type S3Repo struct {
}

func NewS3Repo() (*S3Repo, error) {
	return &S3Repo{}, initS3()
}

type FileType struct {
	File    multipart.File
	Handler *multipart.FileHeader
	Exist   bool
}

var (
	directoryPath = "./s3_storage"
	postBucket    = "/postsImg"
	commentBucket = "/commentsImg"
	maxFileSize   = 5 << 20 // 5 MB
	allowedTypes  = []string{"image/jpeg", "image/png", "image/gif"}
)

func initS3() error {
	buckets := []string{postBucket, commentBucket}

	for _, b := range buckets {
		path := filepath.Join(directoryPath, b)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
			slog.Info("bucket created", "bucket", b)
		}
	}
	return nil
}

func validateFile(f FileType) error {
	if !f.Exist {
		return errors.New("no file uploaded")
	}

	if f.Handler.Size > int64(maxFileSize) {
		return fmt.Errorf("file too large: %d bytes (limit %d)", f.Handler.Size, maxFileSize)
	}

	// Check type
	ct := f.Handler.Header.Get("Content-Type")
	valid := false
	for _, t := range allowedTypes {
		if strings.EqualFold(ct, t) {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("unsupported file type: %s", ct)
	}

	return nil
}

func (r *S3Repo) UploadObject(isPostImg bool, f FileType) (string, error) {

	if !f.Exist {
		return "", nil
	}

	if err := validateFile(f); err != nil {
		return "", err
	}

	var bucket string
	if isPostImg {
		bucket = postBucket
	} else {
		bucket = commentBucket
	}

	fileName := filepath.Base(f.Handler.Filename)
	fullPath := filepath.Join(directoryPath, bucket, fileName)

	newFile, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, f.File); err != nil {
		return "", err
	}

	return filepath.Join(bucket, fileName), nil
}

func (r *S3Repo) GetObject(key string) (string, error) {
	fullPath := filepath.Join(directoryPath, key)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", errors.New("object not found")
	}
	return fullPath, nil
}
