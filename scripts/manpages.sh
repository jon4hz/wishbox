#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run ./cmd/wishbox/main.go man | gzip -c >manpages/wishbox.1.gz