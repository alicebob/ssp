.PHONY: all test build static release

all: test build

test:
	go test ./...

build:
	go build -i ./...
	go build -o ./my_first_ssp

static:
	go get github.com/mjibson/esc
	esc -o static.go --prefix static/ static/
	$(MAKE) build

release:
	go get -v github.com/goreleaser/goreleaser
	goreleaser --rm-dist # --snapshot
