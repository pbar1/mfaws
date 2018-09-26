#!/bin/sh
set -e

OS_ARCH="$(uname -sm)"

if [ "$OS_ARCH" = "Linux x86_64" ]; then
  OS_NAME="linux"
elif [ "$OS_ARCH" = "Darwin x86_64" ]; then
  OS_NAME="darwin"
else
  echo "Your OS/Architecture not supported. Unless you're on Windows, in which case this script ain't gonna work!"
  exit 1
fi

curl -s https://api.github.com/repos/pbar1/mfaws/releases/latest \
| grep "https://github.com/pbar1/mfaws/releases/download" \
| grep "$OS_NAME" \
| cut -d : -f 2,3 \
| tr -d \" \
| wget -qi - -O mfaws

chmod +x mfaws
mkdir -p "$HOME/.local/bin"
mv mfaws "$HOME/.local/bin/mfaws"

echo "Installed to: $HOME/.local/bin"
