#!/bin/bash

set -e

clear

go build trabandcamp.go

if [[ $* == *-d* ]]; then
  debug="DEBUG=app"
fi

if [[ $* == *-y* ]]; then
  ignore="-y"
fi

if [[ $* == *-c* ]]; then
  config="--config .trabandcamprc.sample"
fi

CMD="${debug} ./trabandcamp ${ignore} ${config} electric-moon mayhemofficial"

eval ${CMD}
