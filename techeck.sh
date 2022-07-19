#!/bin/bash

NAME=techeck
CONFIG=data/collector.toml
DATA=$(pwd)/data
PROFILES=/home/boba/profiles/tg/strainer
MIMICRY=true
REPEATE=true
INTERVAL=0


docker run -ti --rm \
        --name $NAME \
        -e CONFIG=$CONFIG \
        -e MIMICRY=$MIMICRY \
        -e REPEATE=$REPEATE \
        -e INTERVAL=$INTERVAL \
        -v $PROFILES:/app/profiles \
        -v $DATA:/app/data \
        techeck
