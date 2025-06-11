package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

// Guarda la imagen en ./uploads/images/ y devuelve el path local
func SaveImage(file *multipart.FileHeader, filename string) (string, error) {
	uploadDir := "./uploads/images"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("error creando directorio: %v", err)
		}
	}

	dst := filepath.Join(uploadDir, filename)
	if err := saveUploadedFile(file, dst); err != nil {
		return "", err
	}

	return dst, nil
}

// Usa Gin para guardar el archivo
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	data := make([]byte, file.Size)
	_, err = src.Read(data)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}

// Elimina imagen del disco si existe
func DeleteImage(path string) error {
	if path == "" {
		return nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // no hay nada que borrar
	}
	return os.Remove(path)
}
