package domain

import "mime/multipart"

type UploadFile struct {
	File    multipart.File
	Handler *multipart.FileHeader
	Exist   bool
}
