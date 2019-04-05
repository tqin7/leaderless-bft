#!/usr/bin/env bash

numOfNodes=14
maxNeighborSize=$((numOfNodes / 2))
start=2 #start IP is 127.0.0.start
end=$(( start + numOfNodes ))

ips=()
for ((i=start; i<end; i++)); do
    ips+=(127.0.0.$i)
done

jsonStr=""

for ((i=0; i<numOfNodes; i++)); do
    numNeighbors=$((1 + RANDOM % maxNeighborSize))

    peerStr=""
    if (("$i" > 0))
    then
        peerStr="${ips[i-1]}"
    fi

    for ((j=0; j<numNeighbors; j++)); do
        randomPeer=${ips[$RANDOM % numOfNodes]}
        peerStr="${peerStr:+$peerStr,}$randomPeer"
    done
    peerObj="{\"ip\":\"${ips[i]}\", \"peers\":\"$peerStr\"}"
    jsonStr="${jsonStr:+$jsonStr,}$peerObj"
done
jsonStr="{\"nodes\":[$jsonStr]}"

echo ${jsonStr} > config.json

# go to https://codebeautify.org/jsonviewer to beautify the json generated
