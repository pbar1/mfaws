.DEFAULT_GOAL := build

get:
	go get -t -v ./...

build: get
	CGO_ENABLED=0 go build

fullbuild:
	goreleaser release --config build/goreleaser.yml --rm-dist --snapshot

release:
	goreleaser release --config build/goreleaser.yml --rm-dist

.PHONY: get build fullbuild release
