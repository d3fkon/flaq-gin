#!/usr/bin/env bash
# Stops the process if something fails
set -xe

# All of the dependencies needed/fetched for your project.
# FOR EXAMPLE:
go mod tidy

# create the application binary that eb uses
go build -o bin/app -ldflags="-s -w"
