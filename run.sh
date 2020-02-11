#!/bin/sh
./build.sh

if [ $? -eq 0 ]; then
    ./cmd/backend/backend $1
fi