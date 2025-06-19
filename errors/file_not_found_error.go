package errors

// FileNotFound is used when a requested file is not found in the system
type FileNotFound struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
}

// Error returns an error message indicating that the file was not found
func (f FileNotFound) Error() string {
	return ToStr("File %v not found at location %v", f.FileName, f.Path)
}
