package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_FileExt(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "Test ext from URL",
			filePath: "https://902-local.s3.us-west-1.amazonaws.com/news/Avatar2023.jpeg",
			expected: "jpeg",
		},
		{
			name:     "Test ext from No schema URL",
			filePath: "902-local.s3.us-west-1.amazonaws.com/news/Avatar2023.jpeg",
			expected: "jpeg",
		},
		{
			name:     "Test ext from Full path",
			filePath: "/user/vinh/Avatar2023.jpeg",
			expected: "jpeg",
		},
		{
			name:     "Test ext from file name",
			filePath: "Avatar2023.png",
			expected: "png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext := FileExt(tt.filePath)
			require.Equal(t, ext, tt.expected)
		})
	}
}

func Test_RenameFile(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		newName  string
		expected string
	}{
		{
			name:     "Test rename from Full path",
			filePath: "/user/vinh/Avatar2023.jpeg",
			newName:  "hello",
			expected: "/user/vinh/hello.jpeg",
		},
		{
			name:     "Test rename from Full path",
			filePath: "/user/vinh/Avatar2023.jpeg",
			newName:  "hello",
			expected: "/user/vinh/hello.jpeg",
		},
		{
			name:     "Test rename from file name",
			filePath: "Avatar2023.png",
			newName:  "hello",
			expected: "hello.png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newFilePath := RenameFile(tt.filePath, tt.newName)
			require.Equal(t, tt.expected, newFilePath)
		})
	}
}
