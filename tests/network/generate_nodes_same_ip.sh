#!/usr/bin/env bash

#echo "Generating network configuration..."

numOfNodes=4

maxNeighborSize=$(echo "l($numOfNodes)" | bc -l)
maxNeighborSize=$( printf "%.0f" $maxNeighborSize )
ips=()
nodesGenerated=0
curPort=30000
while [[ ${nodesGenerated} -lt ${numOfNodes} ]]
do
    if lsof -Pi :${curPort} -sTCP:LISTEN -t >/dev/null ; then
        echo "port in use, skip."
    else
#        echo "generate node $curPort"
        ips+=(127.0.0.1:$curPort)
        nodesGenerated=$((nodesGenerated+1))
    fi
    curPort=$((curPort+1))
done

#echo -e "\nnodes generated: ${ips[@]}"

jsonStr=""

for ((i=0; i<numOfNodes; i++)); do
    numNeighbors=$((1 + RANDOM % maxNeighborSize))

    peerStr="${ips[numOfNodes-1]}"
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

echo "Network graph saved to json file"

# go to https://codebeautify.org/jsonviewer to beautify the json generated
