test.try:
	go test -v ./try

test.utils:
	go test -v ./utils

test.log:
	go test -v ./log

test: test.try test.log test.utils