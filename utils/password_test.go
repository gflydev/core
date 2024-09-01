package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComparePasswords(t *testing.T) {
	type args struct {
		hashedPwd string
		inputPwd  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Compare password TRUE",
			args: args{
				hashedPwd: GeneratePassword("MyPass!"),
				inputPwd:  "MyPass!",
			},
			want: true,
		},
		{
			name: "Compare password FALSE",
			args: args{
				hashedPwd: GeneratePassword("MyPass!"),
				inputPwd:  "YourPass!",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ComparePasswords(tt.args.hashedPwd, tt.args.inputPwd), "ComparePasswords(%v, %v)", tt.args.hashedPwd, tt.args.inputPwd)
		})
	}
}

func TestGeneratePassword(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Generate password",
			args: args{
				p: "YourPass!",
			},
			want: "YourPass!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEqualf(t, tt.want, GeneratePassword(tt.args.p), "GeneratePassword(%v)", tt.args.p)
		})
	}
}
