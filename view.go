package core

import (
	"io"
)

// ==================================================================================
//                                   View Structure
// ==================================================================================

type IView interface {
	// Parse build string from template `tpl` and data `data`
	Parse(tpl string, data Data) string
	// Writer set writer `w` from template `tpl` and data `data`
	Writer(tpl string, data Data, writer io.Writer) error
}

// ==================================================================================
//                                   Default View
// ==================================================================================

var viewError = "View is NULL. Please use core.RegisterView(viewEngine) to register your View Engine"

type DefaultView struct {
}

func (v *DefaultView) Parse(tpl string, data Data) string {
	panic(viewError)
}

func (v *DefaultView) Writer(tpl string, data Data, writer io.Writer) error {
	panic(viewError)
}

var view IView = &DefaultView{}

// ==================================================================================
//                                   Functions
// ==================================================================================

// RegisterView inject View
func RegisterView(v IView) {
	view = v
}
