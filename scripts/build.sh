#!/bin/bash
set -e

echo "Building Go Lambda for Linux/AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 /usr/local/go/bin/go build -o bootstrap ./cmd/authorizer/main.go

echo "Packaging Lambda into function.zip..."
zip function.zip bootstrap

echo "Build and packaging complete: function.zip"
rm bootstrap
