BINARY := mfaws
HASH := $(shell git rev-parse --short HEAD)
VERSION ?= $(HASH)
PLATFORMS := linux darwin windows
OS = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p bin
	GOOS=$(OS) GOARCH=amd64 go build -o bin/$(BINARY)-$(VERSION)-$(OS)-amd64

.PHONY: build
build: linux darwin windows

.PHONY: choco
choco:
	@VERSION=$(VERSION) sh -c './build/ci/choco/deploy.sh'

.PHONY: clean
clean:
	rm -rf bin
	rm build/ci/choco/LICENSE.txt
	rm build/ci/choco/mfaws.nuspec
	rm build/ci/choco/mfaws.exe

.PHONY: version
version:
	@echo $(VERSION)
