export CGO_ENABLED     := 0
export DOCKER_BUILDKIT := 1

.DEFAULT_GOAL := build
LDFLAGS       := -ldflags="-s -w"

get:
	go get -t -v ./...

build: get
	go build $(LDFLAGS)

fullbuild:
	goreleaser release --rm-dist --snapshot

release:
	goreleaser release --rm-dist

.PHONY: get build fullbuild release
