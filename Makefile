all: test main

test:
	cd core/; go test -coverprofile=coverage.out

main: main.go core/core.go utils/utils.go
	go build main.go
