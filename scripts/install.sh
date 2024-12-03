#!/bin/bash
mkdir -p ./build/gohttpd-linux

go mod tidy

if [ $? -ne 0 ]; then
    echo "Error: go mod tidy failed."
    exit 1
else
    echo "go mod tidy successful."
fi

go build -o ./build/gohttpd-linux/gohttpd ./cmd/main.go
if [ $? -ne 0 ]; then
    echo "Error: build failed."
    exit 1
else
    echo "build successful."
fi

if [ -d ./conf ] && [ -d ./html ] && [ -f ./banner.txt ]; then
    cp -r ./conf/ ./html/ ./banner.txt ./build/gohttpd-linux
    echo "cp successful."
else
    echo "Error: Required files or directories are missing."
    exit 1
fi

echo "Build and installation completed successfully!"
