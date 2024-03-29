#!/bin/bash
set -e

# Code Climate tool requires the file to be named c.out and to be in the project root
COVERAGE_OUTPUT=c.out
TEMP_COVERAGE=build/temp_cover.out
FINAL_COVERAGE=build/go_coverage.out
HTML_REPORT=build/coverage.html

export GO_PROXY_CERT_FILE=
export GO_PROXY_KEY_FILE=
export GO_PROXY_PORT=33080

if [ ! -d build ]; then
  mkdir build
fi

echo "mode: set" > $FINAL_COVERAGE

if [ -f $TEMP_COVERAGE ]; then
  rm $TEMP_COVERAGE
fi

PACKAGES=$(go list ./...)
for package in ${PACKAGES}
do
  package_name=$(go list -f '{{.Name}}' ${package})
  file_to_create=$(go list -f '{{.Dir}}/{{.Name}}_test.go' ${package})
  if [ ! -f "${file_to_create}" ]; then
    echo "package ${package_name}" > $file_to_create
  fi

  go test -cover -coverprofile=$TEMP_COVERAGE $package
  cat $TEMP_COVERAGE | grep -v "mode:" | sort -r >> $FINAL_COVERAGE
  rm $TEMP_COVERAGE $file_to_create
done

if [ -f $HTML_REPORT ]; then
  rm $HTML_REPORT
fi

go tool cover -html=$FINAL_COVERAGE -o $HTML_REPORT
go tool cover -func=$FINAL_COVERAGE | grep '^total:'

pushd web >> /dev/null
if [ ! -d "node_modules" ]; then
  npm install
fi
npm run test
popd >> /dev/null
