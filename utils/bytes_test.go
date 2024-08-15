package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync"
	"testing"
)

func testRandByte(t *testing.T) {
	t.Helper()

	values := make([]string, 0)

	for i := 0; i < 10000; i++ {
		n := 32
		dst := make([]byte, n)

		RandByte(dst)

		assert.NotContains(t, string(dst), values)

		for i := range dst {
			assert.False(t, string(rune(i)) == "", fmt.Sprintf("RandBytes() invalid char '%v'", dst[i]))
		}

		assert.False(t, len(dst) != n, fmt.Sprintf("RandBytes() length '%d', want '%d'", len(dst), n))

		values = append(values, string(dst))
	}
}

func Test_RandByte(t *testing.T) {
	testRandByte(t)
}

func Test_RandByteConcurrent(t *testing.T) {
	n := 32
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			testRandByte(t)
		}()
	}

	wg.Wait()
}

func Test_Copy(t *testing.T) {
	str := []byte("cache")
	result := CopyByte(str)

	assert.False(
		t,
		reflect.ValueOf(&str).Pointer() == reflect.ValueOf(&result).Pointer(),
		fmt.Sprintf("Copy() returns the same pointer (source == %p | result == %p)", &str, &result),
	)
}

func Test_ExtendByte_Sub(t *testing.T) {
	dst := RandByte(make([]byte, 5))
	want := ExtendByte(dst, 3)

	assert.NotEqual(t, want, dst)
	assert.LessOrEqual(t, len(want), len(dst))
}

func Test_ExtendByte_More(t *testing.T) {
	dst := RandByte(make([]byte, 5))
	want := ExtendByte(dst, 9)

	assert.NotEqual(t, want, dst)
	assert.LessOrEqual(t, len(dst), len(want))
}

func Test_Equal(t *testing.T) {
	foo := []byte("foo")
	bar := []byte("bar")

	isEqual := EqualByte(foo, bar)

	assert.False(
		t,
		isEqual,
		fmt.Sprintf("Equal(%s, %s) == %v, want %v", foo, bar, isEqual, false),
	)
}

func Test_Extend(t *testing.T) {
	type args struct {
		b       []byte
		needLen int
	}

	type want struct {
		sliceLen int
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Initial 0",
			args: args{
				b:       make([]byte, 0),
				needLen: 10,
			},
			want: want{
				sliceLen: 10,
			},
		},
		{
			name: "Initial 10",
			args: args{
				b:       make([]byte, 10),
				needLen: 5,
			},
			want: want{
				sliceLen: 5,
			},
		},
		{
			name: "Initial 69",
			args: args{
				b:       make([]byte, 45),
				needLen: 69,
			},
			want: want{
				sliceLen: 69,
			},
		},
		{
			name: "Initial 30",
			args: args{
				b:       make([]byte, 45),
				needLen: 30,
			},
			want: want{
				sliceLen: 30,
			},
		},
	}

	for _, tt := range tests {
		b := tt.args.b
		needLen := tt.args.needLen
		sliceLen := tt.want.sliceLen

		t.Run(tt.name, func(t *testing.T) {
			got := ExtendByte(b, needLen)

			gotLen := len(got)

			assert.False(
				t,
				gotLen != sliceLen,
				fmt.Sprintf("ExtendByteSlice() length = %v, want = %v", gotLen, sliceLen),
			)
		})
	}
}

func Test_Prepend(t *testing.T) {
	foo := []byte("foo")
	bar := []byte("bar")

	expected := []byte("foobar")
	result := PrependByte(bar, foo...)

	assert.False(
		t,
		!EqualByte(result, expected),
		fmt.Sprintf("Prepend() == %s, want %s", result, expected),
	)
}

func Test_PrependString(t *testing.T) {
	foo := "foo"
	bar := []byte("bar")

	expected := []byte("foobar")
	result := PrependByteStr(bar, foo)

	assert.False(
		t,
		!EqualByte(result, expected),
		fmt.Sprintf("Prepend() == %s, want %s", result, expected),
	)
}

func Benchmark_Rand(b *testing.B) {
	n := 32
	dst := make([]byte, n)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			RandByte(dst)
		}
	})
}
