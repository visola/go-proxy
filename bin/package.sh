#!/bin/bash
set -ex

build_and_zip() {
  PACKAGE_DIR=build/$3_$2
  mkdir $PACKAGE_DIR
  GOOS=$1 GOARCH=$2 go build -o $PACKAGE_DIR/go-proxy
  zip -j build/$3_$2.zip $PACKAGE_DIR/go-proxy
  rm -Rf $PACKAGE_DIR
}

build_and_zip darwin amd64 mac
build_and_zip linux amd64 linux
build_and_zip windows amd64 win
