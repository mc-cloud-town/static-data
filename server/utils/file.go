package utils

import (
	"errors"
	"io/fs"
	"os"
	"path"
)

func HasFile(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func AutoWriteFile(file_path string, data []byte, mode fs.FileMode) error {
	os.MkdirAll(path.Dir(file_path), fs.ModePerm)
	return os.WriteFile(file_path, data, mode)
}
