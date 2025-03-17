run: coverage.out ydo
	./ydo ${ARGS}

all: coverage.out ydo ydo.exe

DEP_DIRS=core utils cmd
DEP_FILES=$(foreach dir, ${DEP_DIRS}, $(wildcard $(dir)/*.go))

coverage.out: ${DEP_FILES} main.go
	go test ./... -coverprofile=coverage.out || rm coverage.out

ydo: ${DEP_FILES} main.go
	go build

ydo.exe: ${DEP_FILES} main.go
	GOOS=windows go build

clean:
	rm coverage.out ydo ydo.exe
