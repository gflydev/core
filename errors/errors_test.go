package errors

import (
	"github.com/gflydev/core/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	tests := map[string]struct {
		format   string
		args     any
		expected string
		isEqual  bool
	}{
		"Test New Error string": {
			format:   "empty file",
			args:     nil,
			expected: "not empty file",
			isEqual:  false,
		},
		"Test New Error one argument": {
			format:   "file `%s` not found",
			args:     "my-file.pdf",
			expected: "file `my-file.pdf` not found",
			isEqual:  true,
		},
		"Test New Error multi string arguments": {
			format:   "file `%s` is not `%s` type",
			args:     []string{"my-file.pdf", "png"},
			expected: "file `my-file.pdf` is not `png` type",
			isEqual:  true,
		},
		"Test New Error multi number arguments": {
			format:   "file1's size `%v` and file2's size `%v`",
			args:     []int{123, 324},
			expected: "file1's size `123` and file2's size `324`",
			isEqual:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var err error

			// How do I determine whether an array contains
			switch tt.args.(type) {
			case []string, []int:
				err = New(tt.format, utils.UnpackArray(tt.args)...)
				break
			default:
				err = New(tt.format, tt.args)
				break
			}

			if tt.isEqual {
				require.Equal(t, tt.expected, err.Error())
			} else {
				require.NotEqual(t, tt.expected, err.Error())
			}
		})
	}
}
