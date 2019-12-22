#!/bin/bash
set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

WEB_DIR=$SCRIPT_DIR/../web

pushd $WEB_DIR >> /dev/null

npm install
npm run build

popd >> /dev/null
