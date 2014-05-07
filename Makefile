#!/usr/bin/make -f

SHELL=/bin/bash

all: build release

build: deps
	go build

release: clean deps
	mkdir -p build
	mv cgolhuge build/

deps:
	go get

clean:
	rm -rf build
