#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# This script creates a fake GOPATH, symlinks in the current
# directory as uber-go/atomic and verifies that tests still pass.

WORK_DIR=`mktemp -d`
function cleanup {
	rm -rf "$WORK_DIR"
}
trap cleanup EXIT


export GOPATH="$WORK_DIR"
PKG_PARENT="$WORK_DIR/src/github.com/uber-go"
PKG_DIR="$PKG_PARENT/atomic"

mkdir -p "$PKG_PARENT"
ln -s `pwd` "$PKG_DIR"
cd "$PKG_DIR"

make test