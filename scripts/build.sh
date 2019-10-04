#!/bin/bash
set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

rm -Rf build
mkdir build

$SCRIPT_DIR/test.sh
$SCRIPT_DIR/package.sh
