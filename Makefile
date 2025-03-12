all: test main

DEP_DIRS=core utils
DEP_FILES=$(foreach dir, ${DEP_DIRS}, $(wildcard $(dir)/*.go))

test:
	go test ./... -coverprofile=coverage.out

main:
	go build main.go
