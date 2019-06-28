#!/bin/bash

echo --- Running the build ---
bin/build.sh

echo --- Tagging commit ---
git tag "v0.7.4"
git push --tags
