package core

import (
	crand "crypto/rand"
	"github.com/gflydev/core/utils"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

var randBytesPool = bytebufferpool.Pool{}

const (
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetIdxBits = 6                     // 6 bits to represent a charset index
	charsetIdxMask = 1<<charsetIdxBits - 1 // All 1-bits, as many as charsetIdxBits
)

// ExtendByte extends b to needLen bytes.
func extendByte(b []byte, needLen int) []byte {
	b = b[:cap(b)]
	if n := needLen - cap(b); n > 0 {
		b = append(b, make([]byte, n)...)
	}

	return b[:needLen]
}

// randByte returns dst with a cryptographically secure string random bytes.
//
// NOTE: Make sure that dst has the length you need.
func randByte(dst []byte) []byte {
	buf := randBytesPool.Get()
	buf.B = extendByte(buf.B, len(dst))

	if _, err := crand.Read(buf.B); err != nil {
		panic(err)
	}

	size := len(dst)

	for i, j := 0, 0; i < size; j++ {
		// Mask bytes to get an index into the character slice.
		if idx := int(buf.B[j%size] & charsetIdxMask); idx < len(charset) {
			dst[i] = charset[idx]
			i++
		}
	}

	randBytesPool.Put(buf)

	return dst
}

// cleanPath removes the '.' if it is the last character of the route
func cleanPath(path string) string {
	return strings.TrimSuffix(path, ".")
}

// validatePath validate path before add to router
func validatePath(path string) {
	if path == "" || !strings.HasPrefix(path, "/") {
		panic("path must begin with '/' in path '" + path + "'")
	}
}

// getOptionalPaths returns all possible paths when the original path
// has optional arguments
func getOptionalPaths(path string) []string {
	paths := make([]string, 0)

	start := 0
walk:
	for {
		if start >= len(path) {
			return paths
		}

		c := path[start]
		start++

		if c != '{' {
			continue
		}

		newPath := ""
		hasRegex := false
		questionMarkIndex := -1

		brackets := 0

		for end, c := range []byte(path[start:]) {
			switch c {
			case '{':
				brackets++

			case '}':
				if brackets > 0 {
					brackets--
					continue
				} else if questionMarkIndex == -1 {
					continue walk
				}

				end++
				newPath += path[questionMarkIndex+1 : start+end]

				path = path[:questionMarkIndex] + path[questionMarkIndex+1:] // remove '?'
				paths = append(paths, newPath)
				start += end - 1

				continue walk

			case ':':
				hasRegex = true

			case '?':
				if hasRegex {
					continue
				}

				questionMarkIndex = start + end
				newPath += path[:questionMarkIndex]

				if path[:start-2] == "" {
					// include the root slash because the param is in the first segment
					paths = append(paths, "/")

				} else if !utils.IncludeStr(paths, path[:start-2]) {
					// include the path without the wildcard
					// -2 due to remove the '/' and '{'
					paths = append(paths, path[:start-2])
				}
			}
		}
	}
}

// token generate unique token
func token(object ...string) string {
	// Make random data
	currentTime := time.Now().Format("20060102150405")
	randomNum := utils.RandInt64(20)
	// Token
	return utils.Sha256(object, currentTime, randomNum)
}

// quoteStr escape special characters in a given string
func quoteStr(raw string) string {
	bb := bytebufferpool.Get()
	quoted := utils.UnsafeStr(fasthttp.AppendQuotedArg(bb.B, utils.UnsafeBytes(raw)))
	bytebufferpool.Put(bb)

	return quoted
}
