package errors

import "fmt"

// FileNotFound is used when a requested file is not found in the system
type FileNotFound struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
	baseError
}

// NewFileNotFound creates a new FileNotFound error
func NewFileNotFound(fileName, path string) Error {
	err := &FileNotFound{
		FileName: fileName,
		Path:     path,
	}

	// Initialize the baseError fields
	err.code = CodeNotFound
	err.message = fmt.Sprintf("File %v not found at location %v", fileName, path)
	err.stackTrace = captureStackTrace(2)

	return err
}

// Error returns an error message indicating that the file was not found
func (f FileNotFound) Error() string {
	return f.message
}

// WithMetadata adds metadata to the error and returns a new error
func (f FileNotFound) WithMetadata(key string, value interface{}) Error {
	// Create a copy of the error to avoid modifying the original
	newErr := f
	if newErr.metadata == nil {
		newErr.metadata = make(map[string]interface{})
	}
	newErr.metadata[key] = value
	return &newErr
}

// WithContext adds context information to the error and returns a new error
func (f FileNotFound) WithContext(key string, value interface{}) Error {
	// Create a copy of the error to avoid modifying the original
	newErr := f
	if newErr.context == nil {
		newErr.context = make(map[string]interface{})
	}
	newErr.context[key] = value
	return &newErr
}
