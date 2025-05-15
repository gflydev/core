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
				hashedPwd: func() string {
					hash, err := GeneratePassword("MyPass!")
					if err != nil {
						t.Fatalf("Failed to generate password: %v", err)
					}
					return hash
				}(),
				inputPwd: "MyPass!",
			},
			want: true,
		},
		{
			name: "Compare password FALSE",
			args: args{
				hashedPwd: func() string {
					hash, err := GeneratePassword("MyPass!")
					if err != nil {
						t.Fatalf("Failed to generate password: %v", err)
					}
					return hash
				}(),
				inputPwd: "YourPass!",
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
			hash, err := GeneratePassword(tt.args.p)
			assert.NoError(t, err, "GeneratePassword(%v) should not return an error", tt.args.p)
			assert.NotEqualf(t, tt.want, hash, "GeneratePassword(%v)", tt.args.p)
		})
	}
}
