run: fmt coverage.out bench ydo
	./ydo ${ARGS}

all: fmt coverage.out ydo

DEP_DIRS=core utils cmd
DEPS=$(foreach dir, ${DEP_DIRS}, $(wildcard $(dir)/*.go)) main.go
ANDROID_NDK_HOME:=/opt/android-sdk/ndk/27.0.12077973

coverage.out: ${DEPS} main.go
	go test ./... -coverprofile=coverage.out || rm coverage.out

fmt:
	go fmt ./...

ydo: ${DEPS}
	go build -tags "$(TAGS)"

ydo.exe: ${BIN_DEPS}
	GOOS=windows go build -tags "$(TAGS)"

ydo.termux: ${BIN_DEPS}
	GOARCH=arm64 CC=${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang\
				 GOOS=android CGO_ENABLED=1 go build -o ydo.termux -tags "$(TAGS)"

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
