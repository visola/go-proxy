#!/bin/bash
set -e

BASE_PACKAGE=build/packages
mkdir -p $BASE_PACKAGE

build_and_zip() {
  # $1 -> operating system
  # $2 -> architecture
  # $3 -> OS alias, used in the output file name
  # $4 -> Optional extension with ".", e.g.: .exe
  PACKAGE_DIR=$BASE_PACKAGE/$3_$2
  PACKAGE_FILE=$PACKAGE_DIR/go-proxy$4
  mkdir $PACKAGE_DIR
  GOOS=$1 GOARCH=$2 go build -o $PACKAGE_FILE ./cmd/go-proxy
  zip -j $BASE_PACKAGE/$3_$2.zip $PACKAGE_FILE LICENSE
  rm -Rf $PACKAGE_DIR
}

build_and_zip darwin amd64 mac
build_and_zip linux amd64 linux
build_and_zip windows amd64 win .exe
