.DEFAULT_GOAL := build

build:
	goreleaser release --config build/goreleaser.yml --rm-dist --snapshot

release:
	goreleaser release --config build/goreleaser.yml --rm-dist

.PHONY: build
