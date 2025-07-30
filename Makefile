APP_NAME := glvars
BUILD_FILE := Containerfile

.PHONY: all build clean

all: clean build

build:
	@go build -o $(APP_NAME)

clean:
	@rm -f $(APP_NAME)
