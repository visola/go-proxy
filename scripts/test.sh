#!/bin/bash
set -e

# Code Climate tool requires the file to be named c.out and to be in the project root
COVERAGE_OUTPUT=c.out
TEMP_COVERAGE=build/temp_cover.out
HTML_REPORT=build/coverage.html

export GO_PROXY_CERT_FILE=
export GO_PROXY_KEY_FILE=
export GO_PROXY_PORT=33080

echo "mode: set" > $COVERAGE_OUTPUT

if [ -f $TEMP_COVERAGE ]; then
  rm $TEMP_COVERAGE
fi

PACKAGES=$(go list ./...)
for package in ${PACKAGES}
do
  go test -cover -coverprofile=$TEMP_COVERAGE $package
  cat $TEMP_COVERAGE | grep -v "mode:" | sort -r >> $COVERAGE_OUTPUT
  rm $TEMP_COVERAGE
done

if [ -f $HTML_REPORT ]; then
  rm $HTML_REPORT
fi

go tool cover -html=$COVERAGE_OUTPUT -o $HTML_REPORT
