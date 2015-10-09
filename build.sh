#!/bin/bash

docker run --rm \
  -v "$PWD":/go/src/github.com/octoblu/journal-2-logentries \
  -w /go/src/github.com/octoblu/journal-2-logentries \
  golang go build -v
