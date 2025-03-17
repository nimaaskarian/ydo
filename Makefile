run: coverage.out ydo
	./ydo ${ARGS}

all: coverage.out ydo ydo.exe

DEP_DIRS=core utils cmd
DEP_FILES=$(foreach dir, ${DEP_DIRS}, $(wildcard $(dir)/*.go)) $(wildcard cmd/webgui/*/*)
TW_OUT=cmd/webgui/static/tw-out.min.css
TW_IN=cmd/webgui/tailwind.css

coverage.out: ${DEP_FILES} main.go
	go test ./... -coverprofile=coverage.out || rm coverage.out

ydo: ${DEP_FILES} main.go ${TW_OUT}
	go build

ydo.exe: ${DEP_FILES} main.go ${TW_OUT}
	GOOS=windows go build

${TW_OUT}: ${TW_IN}
	which tailwindcss || npm install
	# thank you tailwind. we love you but don't write extra bytes in my css
	npx @tailwindcss/cli -i ${TW_IN} -m | tail -n 1 > ${TW_OUT}

clean:
	rm coverage.out ydo ydo.exe
