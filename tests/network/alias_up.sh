#!/usr/bin/env bash

# creates aliases of the loopback interface for testing purposes
numOfNodes=3
start=2
end=$(( start + numOfNodes ))

for ((i=start; i<end; i++)); do
    ifconfig lo0 alias "127.0.0.$i"
done
