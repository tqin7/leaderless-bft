#!/usr/bin/env bash

serverFile=$1
size=$2
touch execTime.txt
echo "$serverFile" >> execTime.txt

# put below in a for loop for multiple sizes
cd network
./generate_symmetric_graph_same_ip.py $size
cd ..
echo "network size=$size" >> execTime.txt
cd ..
# ./$serverFile | tee -a ./tests/log.txt  # stdout still visible in terminal
./$serverFile > ./tests/log.txt  # stdout not visible in terminal


# cd tests
# ./testSendRequests $serverFile
# # ps -ef | grep $serverFile | grep -v grep | awk '{print $2}' | xargs kill
# ./calcExecTime log.txt
# # rm log.txt