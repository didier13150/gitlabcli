APP_NAME := glcli

.PHONY: all build clean test

all: build

configure:
	@go install

build: configure
	@CGO_ENABLED=0 GOOS=linux go build -o $(APP_NAME) -v -ldflags="-s -w"

clean:
	@rm -f $(APP_NAME)

lint:
	@golangci-lint run

test: build
	@bash -c "glsimulator 1>glsimulator.log 2>&1 &" && go test ; killall glsimulator 1>/dev/null 2>&1 ||:

install:
	@if [ $$(id -u) -eq 0 ] ; then \
		mkdir -p /usr/local/bin ; \
		cp $(APP_NAME) /usr/local/bin/ ; \
	else \
		mkdir -p $${HOME}/.local/bin ; \
		cp $(APP_NAME) $${HOME}/.local/bin/ ; \
	fi
