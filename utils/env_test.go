package utils

import (
	"fmt"
	"testing"
)

func TestGetenv(t *testing.T) {
	got := fmt.Sprintf(
		"%s:%d",
		Getenv("SERVER_HOST", "0.0.0.0"),
		Getenv("SERVER_PORT", 7789),
	)
	want := "0.0.0.0:7789"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
