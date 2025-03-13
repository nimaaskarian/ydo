all: coverage.out ydo

DEP_DIRS=core utils cmd
DEP_FILES=$(foreach dir, ${DEP_DIRS}, $(wildcard $(dir)/*.go))

coverage.out: ${DEP_FILES} main.go
	go test ./... -coverprofile=coverage.out

ydo: ${DEP_FILES} main.go
	go build
