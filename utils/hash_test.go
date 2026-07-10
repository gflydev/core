package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Sha256(t *testing.T) {
	tests := []struct {
		name    string
		one     string
		two     string
		isEqual bool
	}{
		{
			name:    "Test file path",
			one:     "/user/vinh/Avatar2023.jpeg",
			two:     "/user/vinh/hello.jpeg",
			isEqual: false,
		},
		{
			name:    "Test URL path",
			one:     "https://avatar.com/user/vinh/Avatar2023.jpeg",
			two:     "https://avatar.com/user/vinh/Avatar2023.jpeg",
			isEqual: true,
		},
		{
			name:    "Test random unique",
			one:     Token(),
			two:     Token(),
			isEqual: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isEqual {
				require.Equal(t, Sha256(tt.one), Sha256(tt.two))
			} else {
				require.NotEqual(t, Sha256(tt.one), Sha256(tt.two))
			}
		})
	}
}

func Test_MD5(t *testing.T) {
	tests := []struct {
		name    string
		one     string
		two     string
		isEqual bool
	}{
		{
			name:    "Test file path",
			one:     "/user/vinh/Avatar2023.jpeg",
			two:     "/user/vinh/hello.jpeg",
			isEqual: false,
		},
		{
			name:    "Test URL path",
			one:     "https://avatar.com/user/vinh/Avatar2023.jpeg",
			two:     "https://avatar.com/user/vinh/Avatar2023.jpeg",
			isEqual: true,
		},
		{
			name:    "Test random unique",
			one:     Token(),
			two:     Token(),
			isEqual: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isEqual {
				require.Equal(t, MD5(tt.one), MD5(tt.two))
			} else {
				require.NotEqual(t, MD5(tt.one), MD5(tt.two))
			}
		})
	}
}

// Test_Sha256_MultiArg verifies that multiple arguments are all incorporated
// into the hash (previously the slice was passed without expansion, so only the
// first verb consumed a value and the rest rendered as "%!v(MISSING)"), so
// distinct trailing arguments must produce distinct hashes.
func Test_Sha256_MultiArg(t *testing.T) {
	// Differ only in the second argument.
	h1 := Sha256("a", "b")
	h2 := Sha256("a", "c")
	require.NotEqual(t, h1, h2, "second argument must affect the hash")

	// Output must not contain fmt's missing-verb marker.
	require.NotContains(t, Sha256("a", "b", "c"), "MISSING")

	// Same inputs must be stable.
	require.Equal(t, Sha256("x", "y", "z"), Sha256("x", "y", "z"))
}
