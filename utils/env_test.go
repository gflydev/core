package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Getenv(t *testing.T) {
	got := fmt.Sprintf(
		"%s:%d",
		Getenv("SERVER_HOST", "0.0.0.0"),
		Getenv("SERVER_PORT", 7789),
	)
	want := "0.0.0.0:7789"

	assert.Equal(t, want, got)
}
