#!/bin/bash
set -ex

if ! which dep >> /dev/null; then
  echo Installing dep...
  curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
fi


echo Updating dependencies...
dep ensure
