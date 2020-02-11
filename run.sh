#!/bin/sh
./build.sh

if [ $? -eq 0 ]; then
    ./cmd/api/api $1
fi