#!/bin/sh
#Entry point docker container
CMD=./checker
cd /app

if [ ! -z "${CONFIG}" ];then
    CMD=$CMD" -c "$CONFIG
fi

if [ ! -z "${MIMICRY}" ];then
    CMD=$CMD" -m"
fi

if [ ! -z "${INTERVAL}" ];then
    CMD=$CMD" -i "$INTERVAL
fi

if [ ! -z "${PROFILES}" ];then
    CMD=$CMD" -p "$PROFILES
else
    CMD=$CMD" -p ./profiles"
fi

$CMD
