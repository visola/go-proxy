#!/bin/bash
set -e

go get -u github.com/gobuffalo/packr/v2/packr2
packr2 clean
packr2
