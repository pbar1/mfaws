#!/bin/bash
set -e

cd "$TRAVIS_BUILD_DIR/build/ci/aur"

cat sshconfig >> "$HOME/.ssh/config"
openssl aes-256-cbc -K $encrypted_5db6e1fb054e_key -iv $encrypted_5db6e1fb054e_iv -in aur.enc -out ~/.ssh/aur -d
chmod 400 ~/.ssh/aur

git clone ssh://aur@aur.archlinux.org/mfaws-bin.git

curl -s https://api.github.com/repos/pbar1/mfaws/releases/$TRAVIS_TAG \
| grep "https://github.com/pbar1/mfaws/releases/download" \
| grep "linux" \
| cut -d : -f 2,3 \
| tr -d \" \
| wget -qi - -O mfaws

export VERSION=$(echo "$TRAVIS_TAG" | tr -d v)
sed -i "s/VERSION/$VERSION/g" PKGBUILD .SRCINFO

export CHECKSUM=$(sha256sum mfaws | cut -d ' ' -f 1)
sed -i "s/CHECKSUM/$CHECKSUM/g" PKGBUILD .SRCINFO

cp PKGBUILD mfaws-bin
cp .SRCINFO mfaws-bin
cd mfaws-bin

git add PKGBUILD .SRCINFO
git commit -m "Deployed $TRAVIS_TAG from Travis CI"
git push
