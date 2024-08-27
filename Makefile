mod:
	go list -m --versions

test.try:
	go test -v ./try

test.utils:
	go test -v ./utils

test.log:
	go test -v ./log

test.errors:
	go test -v ./errors

test: test.try test.log test.utils test.errors

test.cover:
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -html=cover.out

critic:
	gocritic check -enableAll -disable=unnamedResult,unlabelStmt,hugeParam,singleCaseSwitch,builtinShadow,typeAssertChain ./...

security:
	gosec -exclude-dir=core -exclude=G103,G401,G501 ./...

vulncheck:
	govulncheck ./...

lint:
	golangci-lint run ./...

all: critic security vulncheck lint test