all: options test main

DEP_DIRS=core utils
DEP_FILES=$(foreach dir, ${DEP_DIRS}, $(wildcard $(dir)/*.go))

options:
	@echo Dirs:  ${DEP_DIRS}
	@echo Files: ${DEP_FILES}

test:
	@$(foreach dir, ${DEP_DIRS}, cd $(dir); go test -coverprofile=coverage.out || exit 1; cd ..;)

main: main.go ${DEP_FILES}
	go build main.go
