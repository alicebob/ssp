.PHONY: all test build

all: test build

test:
	go test

build:
	go build -o ./my_first_ssp
