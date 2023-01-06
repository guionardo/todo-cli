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
flags="-X ${repo}/cmd.BuildDate=${BUILD_DATE} -X ${repo}/cmd.BuildHost=${BUILD_HOST} -X ${repo}/cmd.Version=${RELEASE}"
echo "Flags:      ${flags}"
go build -ldflags "$flags" .
go install -ldflags "$flags" .
go_bin="$(go env GOPATH)/bin"
mv "${go_bin}/todo-cli" "${go_bin}/todo"
