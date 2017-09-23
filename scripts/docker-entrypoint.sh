#!/usr/bin/env bash
echo "Starting container services ..."

export GOPATH=/root/.go
export PATH=$PATH:$GOPATH/bin

# endless command to keep the container running
bash