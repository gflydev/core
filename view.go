package core

import (
	"io"
)

// ====================================================================
//                          View Structure
// ====================================================================

// IView defines the interface for a view engine that can parse templates
// into strings or write the rendered content directly to an io.Writer.
type IView interface {
	// Parse builds a string from the template `tpl` and the provided data `data`.
	//
	// Parameters:
	//   - tpl: The template string to be rendered.
	//   - data: A Data map containing key-value pairs to be injected into the template.
	//
	// Returns:
	//   - A string representation of the rendered template.
	Parse(tpl string, data Data) string

	// Writer writes the rendered content of template `tpl` and data `data`
	// to the provided io.Writer `writer`.
	//
	// Parameters:
	//   - tpl: The template string to be rendered.
	//   - data: A Data map containing key-value pairs to be injected into the template.
	//   - writer: An io.Writer where the rendered content will be written.
	//
	// Returns:
	//   - An error if writing to the writer fails, otherwise nil.
	Writer(tpl string, data Data, writer io.Writer) error
}

// ====================================================================
//                           Default View
// ====================================================================

// viewError is the error message displayed when no view engine is registered.
var viewError = "View is NULL. Please use core.RegisterView(viewEngine) to register your View Engine"

// DefaultView is a fallback view engine that panics when methods are called.
// It ensures that an appropriate view engine is registered to avoid runtime errors.
type DefaultView struct {
}

// Parse parses the given template `tpl` with the provided data `data`
// and returns the resulting string. In the case of DefaultView, it panics
// with the error message `viewError` because no view engine is registered.
//
// Parameters:
//   - tpl: The template string to be parsed.
//   - data: The data to be injected into the template.
//
// Returns:
//   - A string representation of the rendered template.
//
// Note: This function always panics with `viewError` for the DefaultView implementation.
func (v *DefaultView) Parse(tpl string, data Data) string {
	panic(viewError)
}

// Writer panics with viewError when called, indicating no view engine is registered.
//
// Parameters:
//   - tpl: The template string to be written.
//   - data: The data to be injected into the template.
//   - writer: The io.Writer where the rendered content is written.
//
// Returns:
//   - Always panics with `viewError` for DefaultView implementation.
func (v *DefaultView) Writer(tpl string, data Data, writer io.Writer) error {
	panic(viewError)
}

// view is the global instance of the IView interface used to render templates.
var view IView = &DefaultView{}

// ====================================================================
//                              Functions
// ====================================================================

// RegisterView registers a custom view engine by setting the global `view` instance to `v`.
// This should be called during application initialization to use a specific view engine.
//
// Parameters:
//   - v: An implementation of the IView interface representing the custom view engine.
func RegisterView(v IView) {
	view = v
}
