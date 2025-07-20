# Makefile for ponghub Go project

BINARY=bin/ponghub.exe
SRC=cmd/main.go
CONFIG=config.yaml

.PHONY: all build run clean test

all: build

build:
	go build -o $(BINARY) $(SRC)

run: build
	$(BINARY) --config $(CONFIG)

clean:
	del $(BINARY)
