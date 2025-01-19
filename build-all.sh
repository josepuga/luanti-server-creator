#!/bin/bash
# By José Puga. 2024-2025. GPL3 License
# Compiles project to Linux & Win, 64bits.
set -e

VERSION=$(git describe --tags)
FLAGS="-w -s -X main.version=$VERSION"
cd src/

# This does not work, because fyne is not Windows 100% compatible.
#for so in linux windows; do
    #GOOS=$so go build -o ../bin/ -ldflags="$FLAGS" .
#done

go build -o ../bin/ -ldflags="$FLAGS" .

# For test purposes
#cd .. && cp bin/servercreator . && false

GOFLAGS="-ldflags=$FLAGS" fyne-cross windows \
    --app-id dev.puga.servercreator \
    --output servercreator.exe \
    .
unzip -o -x fyne-cross/dist/windows-amd64/servercreator.exe.zip -d ../bin
rm -fr fyne-cross/