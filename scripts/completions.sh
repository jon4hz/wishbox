#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run ./cmd/wishbox/main.go completion "$sh" >"completions/wishbox.$sh"
done