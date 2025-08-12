APP_NAME := glcli
SIMU_NAME := glsimulator
BUILD_FILE := Containerfile

.PHONY: all build clean test

all: build

configure:
	@go install
	@cd test/glsimulator && go install

build: configure
	@CGO_ENABLED=0 GOOS=linux go build -o $(APP_NAME) -v -ldflags="-s -w"
	@cd test/glsimulator && CGO_ENABLED=0 GOOS=linux go build -o $(SIMU_NAME) -v -ldflags="-s -w"

clean:
	@rm -f $(APP_NAME) test/glsimulator/glsimulator

lint:
	@golangci-lint run

test: build
	@bash -c "test/glsimulator/$(SIMU_NAME) 1>glsimulator.log 2>&1 &" && cd gitlabcli && go test ; killall $(SIMU_NAME) 1>/dev/null 2>&1 ||:
