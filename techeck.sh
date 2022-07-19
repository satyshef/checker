#!/bin/bash

NAME=techeck
CONFIG=data/collector.toml
DATA=$(pwd)/data
PROFILES=/home/boba/profiles/tg/strainer
MIMICRY=true
INTERVAL=0


docker run -ti --rm \
        --name $NAME \
        -e CONFIG=$CONFIG \
        -e MIMICRY=$MIMICRY \
        -e INTERVAL=$INTERVAL \
        -v $PROFILES:/app/profiles \
        -v $DATA:/app/data \
        techeck
