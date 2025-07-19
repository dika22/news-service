package utils

import (
	"os"
	"path/filepath"
)

func GetMappingPath(fileName string) string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "elastic_doc", fileName)
}
