all: test main

test:
	cd core/; go test -coverprofile=coverage.out

main: main.go core/*.go utils/*.go
	go build main.go
