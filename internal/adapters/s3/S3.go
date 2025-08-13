package s3

import (
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileType struct {
	File    multipart.File
	Handler *multipart.FileHeader
	Exist   bool
}

var (
	directoryPath = "./s3_storage"
	postPath      = "/postsImg"
	commentPath   = "/commentsImg"
)

func InitS3() error {

	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err := os.MkdirAll(directoryPath, 0755)
		if err != nil {
			return err
		}
		slog.Info("s3 directory successfully created")
	}

	return nil
}

func UploadObject(isPostImg bool, f FileType) (string, error) {

	if !f.Exist {
		return "", nil
	}

	// ct := f.Handler.Header.Get("Content-Type")
	// need to implement that only image type has permission

	var basePath string
	if isPostImg {
		basePath = postPath
	} else {
		basePath = commentPath
	}

	// need to test проверка по инвалидных путей или данных

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Join(directoryPath, basePath), 0755); err != nil {
		return "", err
	}

	// Use filepath.Join for OS-safe paths
	fileName := filepath.Base(f.Handler.Filename) // prevents path traversal
	fullPath := filepath.Join(directoryPath, basePath, fileName)

	// Create file
	newFile, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	// Copy file content
	if _, err := io.Copy(newFile, f.File); err != nil {
		return "", err
	}

	// fmt.Println("File saved to:", fullPath)
	return filepath.Join(basePath, fileName), nil
}

func GetObject(imageUrl string) {

}
