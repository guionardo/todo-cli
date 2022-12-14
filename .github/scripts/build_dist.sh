#!/bin/env bash

# https://stackoverflow.com/questions/58177786/get-the-current-pushed-tag-in-github-actions

echo "Building distribution"
echo "====================="
BUILD_DATE=$(date -Iseconds)
BUILD_HOST=$(hostname)
gh_release=${GITHUB_REF#/refs/*/}
RELEASE=${gh_release:-develop}

echo "Build date: ${BUILD_DATE}"
echo "Build Host: ${BUILD_HOST}"
echo "Release:    ${RELEASE}"
repo="github.com/guionardo/todo-cli"
flags="-X ${repo}/utils.BuildDate=${BUILD_DATE} -X ${repo}/utils.BuildHost=${BUILD_HOST} -X ${repo}/utils.Version=${RELEASE}"
echo "Flags:      ${flags}"
go build -ldflags "$flags" .
go install -ldflags "$flags" .