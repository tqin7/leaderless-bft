#!/bin/bash

numOfNodes=4
start=2
end=$(( start + numOfNodes ))

for ((i=start; i<end; i++)); do
    ifconfig lo0 -alias "127.0.0.$i"
done
