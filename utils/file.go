package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveImage(file *multipart.FileHeader, filename string) (string, error) {
	path := "uploads/images/"
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}
	fullPath := filepath.Join(path, filename)
	return fullPath, SaveUploadedFile(file, fullPath)
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func DeleteImage(path string) error {
	return os.Remove(path)
}
