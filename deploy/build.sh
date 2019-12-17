#!/bin/sh
env GOOS=linux GOARCH=amd64 go build -o deploy-linux
env GOOS=windows GOARCH=amd64 go build
env GOOS=darwin GOARCH=amd64 go build