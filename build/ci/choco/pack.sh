#!/bin/bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
cd "$DIR"

cp "../../../LICENSE" LICENSE.txt
sed "s/VERSION/$VERSION/g" mfaws.nuspec.template > mfaws.nuspec

# function choco(){ sudo docker run --rm -v "$(pwd)":"$(pwd)" -w "$(pwd)" linuturk/mono-choco "$@" ;}
# choco pack
