package utils

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func testToken(t *testing.T) {
	t.Helper()

	values := make([]string, 0)

	for i := 0; i < 10000; i++ {
		token := Token()

		assert.NotContains(t, token, values)

		values = append(values, token)
	}
}

func Test_Token(t *testing.T) {
	testToken(t)
}

func Test_TokenConcurrent(t *testing.T) {
	n := 32
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			testToken(t)
		}()
	}

	wg.Wait()
}
