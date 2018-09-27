#!/bin/bash

echo --- Running the build ---
bin/run.sh

echo --- Tagging commit ---
git tag "v0.7.0"
git push --tags
