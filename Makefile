.PHONY: run debug build

build:
	go build

run:
	./dvml2 -dir conf

debug:
	./dvml2 -dir conf -debug

dev: build debug