#!/bin/sh

go build --ldflags "-X main.GitCommit=$(git rev-parse HEAD) -X main.Version=$(git symbolic-ref --short -q HEAD) " -o ./metrics-forwarder
