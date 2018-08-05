#!/bin/bash
set -e

~/semantic-release -ghr -vf
export VERSION=$(cat .version)
gox -ldflags="-s -w" -os="linux darwin windows" -arch="amd64" -output="bin/{{.Dir}}_v${VERSION}_{{.OS}}_{{.Arch}}"
ghr $(cat .ghr) bin/

cp "bin/mfaws_v${VERSION}_windows_amd64.exe" "buildassets/mfaws.exe"
cp LICENSE.txt buildassets/LICENSE.txt
sed -i "s/VERSION/$VERSION/g" buildassets/mfaws.nuspec
cd buildassets
function choco(){ sudo docker run --rm -v $(pwd):$(pwd) -w $(pwd) linuturk/mono-choco "$@" ;}
choco pack
choco push -s https://push.chocolatey.org/ -k "$CHOCO_API_KEY"
