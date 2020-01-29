#!/bin/bash

# ./cc-test-reporter after-build --exit-code $TRAVIS_TEST_RESULT

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR=$SCRIPT_DIR/..
BUILD_DIR=$SCRIPT_DIR/../build

CC_REPORTER=$ROOT_DIR/cc-test-reporter

GO_FORMATTED_REPORT=$BUILD_DIR/go_coverage.json
FE_FORMATTED_REPORT=$BUILD_DIR/fe_coverage.json

SUMMED_REPORT=$BUILD_DIR/summed_report.json

echo "Formatting reports..."
$CC_REPORTER format-coverage -t gocov $BUILD_DIR/go_coverage.out -o $GO_FORMATTED_REPORT --prefix github.com/visola/go-proxy

pushd web > /dev/null
$CC_REPORTER format-coverage -t lcov $BUILD_DIR/lcov.info -o $FE_FORMATTED_REPORT
popd > /dev/null

echo "Generating combine report..."
$CC_REPORTER sum-coverage -o $SUMMED_REPORT $GO_FORMATTED_REPORT $FE_FORMATTED_REPORT

echo "Uploading report..."
$CC_REPORTER upload-coverage -i $SUMMED_REPORT
