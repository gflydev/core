package utils

import (
	"fmt"
	"path/filepath"
)

// FileExt extracts the extension of a file from its name or path.
//
// Parameters:
//   - fileName (string): The file name or path from which to extract the extension.
//
// Returns:
//   - string: The file extension without the leading dot, or an empty string if none is found.
//
// Examples:
//   - Input: "https://example.com/file.jpeg", Output: "jpeg"
//   - Input: "path/to/file.png", Output: "png"
//   - Input: "no_extension_file", Output: ""
func FileExt(fileName string) string {
	// Extract the file extension including the leading dot.
	filePart := filepath.Ext(fileName)

	// Return the extension without the leading dot if it exists.
	if filePart != "" {
		return filePart[1:]
	}

	return filePart
}

// RenameFile generates a new file path by replacing the file name with a new name while preserving the extension.
//
// Parameters:
//   - fileName (string): The original file name or full file path.
//   - newName (string): The new base name to replace the old file name.
//
// Returns:
//   - string: The new file path with the same directory and extension as the original.
//
// Examples:
//   - Input: "Avatar2023.jpeg", "hello", Output: "hello.jpeg"
//   - Input: "path/to/file.png", "newname", Output: "path/to/newname.png"
//   - Input: "file_without_extension", "name", Output: "name"
func RenameFile(fileName, newName string) string {
	// Extract the file extension and create a new file base name.
	newBase := fmt.Sprintf("%s%s", newName, filepath.Ext(fileName))

	// Extract the directory of the original file.
	filePath := filepath.Dir(fileName)

	// Combine the directory and new base name to form the new file path.
	return filepath.Join(filePath, newBase)
}
