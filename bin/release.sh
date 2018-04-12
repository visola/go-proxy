#!/bin/bash

echo --- Running the build ---
bin/run.sh

echo --- Tagging commit ---
git tag "v0.6.0"
git push --tags
