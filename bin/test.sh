#!/bin/bash
set -ex

PACKAGES=$(go list ./... | grep -v /vendor/)
go test $PACKAGES
