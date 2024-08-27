package utils

import (
	"fmt"
	"path/filepath"
)

// FileExt Extract extension of file.
//
//	Eg Get `jpeg` from `https://902-local.s3.us-west-1.amazonaws.com/news/Avatar2023.jpeg`
//		Or `63e85ba1.png` from `storage/tmp/63e85ba1.png`
func FileExt(fileName string) string {
	filePart := filepath.Ext(fileName)

	if filePart != "" {
		return filePart[1:]
	}

	return filePart
}

// RenameFile Extract new file path.
//
//	Eg Get `hello.jpeg` from `Avatar2023.jpeg` and `hello`
//		Or `storage/tmp/hello.png` from `storage/tmp/63e85ba1.png` and `hello`
func RenameFile(fileName, newName string) string {
	// Create new file base
	newBase := fmt.Sprintf("%s%s", newName, filepath.Ext(fileName))
	// Get file path
	filePath := filepath.Dir(fileName)

	return filepath.Join(filePath, newBase)
}
