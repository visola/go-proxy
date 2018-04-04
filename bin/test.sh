#!/bin/bash
set -ex

if [ -f cc-test-reporter ]; then
  ./cc-test-reporter before-build
fi

COVERAGE_OUTPUT=build/coverage.out
TEMP_COVERAGE=build/temp_cover.out

echo "mode: set" > $COVERAGE_OUTPUT

if [ -f $TEMP_COVERAGE ]; then
  rm $TEMP_COVERAGE
fi

PACKAGES=$(go list ./... | grep -v /vendor/)
for package in ${PACKAGES}
do
  go test -coverprofile=$TEMP_COVERAGE $package
  if [ -f $TEMP_COVERAGE ]; then
    cat $TEMP_COVERAGE | grep -v "mode:" | sort -r >> $COVERAGE_OUTPUT
    rm $TEMP_COVERAGE
  fi
done

if [ -f cc-test-reporter ]; then
  ./cc-test-reporter after-build
fi
