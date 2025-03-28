run: coverage.out bench ydo
	./ydo ${ARGS}

all: coverage.out ydo ydo.exe

DEP_DIRS=core utils cmd
DEP_FILES=$(foreach dir, ${DEP_DIRS}, $(wildcard $(dir)/*.go)) $(wildcard cmd/webgui/*/*)
TW_OUT=cmd/webgui/static/tw-out.min.css
TW_IN=cmd/webgui/tailwind.css
HAS_TW := $(shell command -v tailwindcss 2> /dev/null)
ANDROID_NDK_HOME:=/opt/android-sdk/ndk/27.0.12077973

coverage.out: ${DEP_FILES} main.go
	go test ./... -coverprofile=coverage.out || rm coverage.out

ydo: ${ICONS} ${DEP_FILES} main.go ${TW_OUT}
	go build

ydo.exe: ${ICONS} ${DEP_FILES} main.go ${TW_OUT}
	GOOS=windows go build

ydo.termux:
	GOARCH=arm64 CC=${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang\
				 GOOS=android CGO_ENABLED=1 go build -o ydo.termux

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
ifneq (,$(wildcard new.bench))
	mv new.bench last.bench
endif
	go test ./... -bench=. -benchtime=20s -benchmem > new.bench
ifneq (,$(wildcard last.bench))
	benchstat last.bench new.bench
endif

clean:
	rm coverage.out ydo ydo.exe
