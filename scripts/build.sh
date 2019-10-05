#!/bin/bash
set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

rm -Rf build
mkdir build

$SCRIPT_DIR/test.sh

npm install
npm run bundle

go get -u github.com/gobuffalo/packr/...
packr clean
packr