#!/bin/bash
set -ex

PACKAGES=$(go list ./... | grep -v /vendor/)
for package in ${PACKAGES}
do
  go test $package
done
