run: coverage.out bench ydo
	./ydo ${ARGS}

all: coverage.out ydo ydo.exe

DEP_DIRS=core utils cmd
DEP_FILES=$(foreach dir, ${DEP_DIRS}, $(wildcard $(dir)/*.go)) $(wildcard cmd/webgui/*/*)
TW_OUT=cmd/webgui/static/tw-out.min.css
TW_IN=cmd/webgui/tailwind.css
HAS_TW := $(shell command -v tailwindcss 2> /dev/null)

coverage.out: ${DEP_FILES} main.go
	go test ./... -coverprofile=coverage.out || rm coverage.out

ydo: ${ICONS} ${DEP_FILES} main.go ${TW_OUT}
	go build

ydo.exe: ${ICONS} ${DEP_FILES} main.go ${TW_OUT}
	GOOS=windows go build

# thank you tailwind. we love you but don't write extra bytes in my css.
# sorry if you got tailed
${TW_OUT}: $(wildcard cmd/webgui/*/*.html ) ${TW_IN}
ifndef HAS_TW
	npm install
	npx @tailwindcss/cli -i ${TW_IN} -m | tail -n 1 > ${TW_OUT}
else
	tailwindcss -i ${TW_IN} -m | tail -n 1 > ${TW_OUT}
endif

bench:
	mv new.bench last.bench
	go test ./... -bench=. -benchtime=20s -benchmem > tee new.bench
	benchstat last.bench new.bench

clean:
	rm coverage.out ydo ydo.exe
