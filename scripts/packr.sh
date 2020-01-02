#!/bin/bash
set -e

go get -u github.com/gobuffalo/packr/packr
packr clean
packr 
