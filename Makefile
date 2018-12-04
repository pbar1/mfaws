.DEFAULT_GOAL := build

get:
	go get -t -v ./...

build: get
	CGO_ENABLED=0 go build

fullbuild:
	goreleaser release --rm-dist --snapshot

release:
	goreleaser release --rm-dist

.PHONY: get build fullbuild release
