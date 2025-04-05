package service

import (
	"context"
	"fmt"
	"io"
)

// Uploader interface defines the methods for uploading files.
type Uploader interface {
	UploadFileToDB(ctx context.Context, file io.Reader, filename string) ([]byte, error)
}

// dbUploader is an implementation of the Uploader interface for database storage.
type dbUploader struct{}

// NewDBUploader creates a new instance of dbUploader.
func NewDBUploader() Uploader {
	return &dbUploader{}
}

// UploadFileToDB saves the uploaded file as binary data to the database.
func (u *dbUploader) UploadFileToDB(ctx context.Context, file io.Reader, filename string) ([]byte, error) {
	var fileData []byte
	buffer := make([]byte, 1024) // Đọc từng phần 1KB

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		if n == 0 {
			break
		}
		fileData = append(fileData, buffer[:n]...)
	}

	return fileData, nil
}
