package s3

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileType struct {
	File    multipart.File
	Handler *multipart.FileHeader
}

var (
	directoryPath = "./s3_storage"
	postPath      = "./s3_storage/postsImg"
	commentPath   = "./s3_storage/commentsImg"
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

	// if f.Handler.Size == 0 {
	// 	return "", nil
	// }

	var basePath string
	if isPostImg {
		basePath = postPath
	} else {
		basePath = commentPath
	}

	// Ensure directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		fmt.Println("fdfsd")
		return "", err
	}
	fmt.Println("haha")
	// Use filepath.Join for OS-safe paths
	fileName := filepath.Base(f.Handler.Filename) // prevents path traversal
	fullPath := filepath.Join(basePath, fileName)

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

	fmt.Println("File saved to:", fullPath)
	return fullPath, nil
}

func GetObject(imageUrl string) {

}
