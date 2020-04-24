#!/usr/bin/env bash

serverFile=$1

echo "sending requests"
./testSendRequests $serverFile
echo "done"
echo "calculating execution time"
./calcExecTime log.txt
echo "done"
# rm log.txt