#!/bin/bash

if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "running"
else
    echo "not running"
fi
