#!/bin/sh

go build --ldflags "-X main.Version=v0.0.1 -X main.GitCommit=$(git rev-parse HEAD)" -o ./metrics-forwarder

