package core

import (
	"io"
)

type FnPerformViewWriter func(template string, data Data, writer io.Writer) error

var fnPerformViewWriter FnPerformViewWriter

// RegisterViewWriter inject View Writer
func RegisterViewWriter(fn FnPerformViewWriter) {
	fnPerformViewWriter = fn
}
