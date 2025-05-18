package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
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

// FileExists checks if a file exists and is not a directory.
//
// Parameters:
//   - path (string): The path to the file to check.
//
// Returns:
//   - bool: True if the file exists and is not a directory, false otherwise.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists.
//
// Parameters:
//   - path (string): The path to the directory to check.
//
// Returns:
//   - bool: True if the directory exists, false otherwise.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// FileSize returns the size of a file in bytes.
//
// Parameters:
//   - path (string): The path to the file.
//
// Returns:
//   - int64: The size of the file in bytes, or -1 if the file doesn't exist or there's an error.
func FileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return -1
	}
	return info.Size()
}

// ReadFileAsString reads the contents of a file as a string.
//
// Parameters:
//   - path (string): The path to the file to read.
//
// Returns:
//   - string: The contents of the file as a string.
//   - error: An error if the file couldn't be read.
func ReadFileAsString(path string) (string, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteStringToFile writes a string to a file.
//
// Parameters:
//   - path (string): The path to the file to write to.
//   - content (string): The content to write to the file.
//   - perm (os.FileMode): The file permissions to use if the file is created.
//
// Returns:
//   - error: An error if the file couldn't be written to.
func WriteStringToFile(path, content string, perm os.FileMode) error {
	return os.WriteFile(path, []byte(content), perm)
}

// CopyFile copies a file from src to dst.
//
// Parameters:
//   - src (string): The path to the source file.
//   - dst (string): The path to the destination file.
//   - perm (os.FileMode): The file permissions to use for the destination file.
//
// Returns:
//   - error: An error if the file couldn't be copied.
func CopyFile(src, dst string, perm os.FileMode) error {
	srcFile, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer func(srcFile *os.File) {
		_ = srcFile.Close()
	}(srcFile)

	// #nosec G304
	dstFile, err := os.OpenFile(filepath.Clean(dst), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer func(dstFile *os.File) {
		_ = dstFile.Close()
	}(dstFile)

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// MkdirIfNotExists creates a directory if it doesn't exist.
//
// Parameters:
//   - path (string): The path to the directory to create.
//   - perm (os.FileMode): The directory permissions to use if the directory is created.
//
// Returns:
//   - error: An error if the directory couldn't be created.
func MkdirIfNotExists(path string, perm os.FileMode) error {
	if !DirExists(path) {
		return os.MkdirAll(path, perm)
	}
	return nil
}

// FileModTime returns the modification time of a file.
//
// Parameters:
//   - path (string): The path to the file.
//
// Returns:
//   - time.Time: The modification time of the file, or the zero time if there's an error.
//   - error: An error if the file information couldn't be retrieved.
func FileModTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// IsFileNewer checks if file1 is newer than file2.
//
// Parameters:
//   - file1 (string): The path to the first file.
//   - file2 (string): The path to the second file.
//
// Returns:
//   - bool: True if file1 is newer than file2, false otherwise or if there's an error.
//   - error: An error if the file information couldn't be retrieved.
func IsFileNewer(file1, file2 string) (bool, error) {
	time1, err := FileModTime(file1)
	if err != nil {
		return false, err
	}

	time2, err := FileModTime(file2)
	if err != nil {
		return false, err
	}

	return time1.After(time2), nil
}
