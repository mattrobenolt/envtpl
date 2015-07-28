#!/bin/bash
set -x

rm -rf bin/
docker build --rm -t envtpl .
docker run --rm -v $PWD:/go/src/app envtpl
