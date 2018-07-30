#!/bin/bash
set -e

~/semantic-release -ghr -vf
export VERSION=$(cat .version)
gox -ldflags="-s -w" -os="linux darwin windows" -arch="amd64" -output="bin/{{.Dir}}_v"$VERSION"_{{.OS}}_{{.Arch}}"
ghr $(cat .ghr) bin/
