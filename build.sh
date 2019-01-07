#!/bin/sh

go build xmcd2cue.go

os=''
binext=''
case "$(uname | tr '[:upper:]' '[:lower:]' 2>/dev/null)" in
    windows*|mingw*|msys*|cygwin*) os='win'; binext='.exe';;
    linux*)                        os='linux';;
    darwin*)                       os='osx';;
    bsd*)                          os='bsd';;
    solaris*)                      os='solaris';;
esac

if [ -n "${os}" ]; then
    rm -f "xmcd2cue-${os}.zip"
    zip -X -9 -o "xmcd2cue-${os}" LICENSE.md README.md "xmcd2cue${binext}"
fi
