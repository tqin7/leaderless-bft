#!/usr/bin/env bash

mkdir -p ./cluster/tests/network
mkdir ./cluster/clients

echo 'Generating network graph...'
cd ./tests/network
./generate_nodes_same_ip.sh
cd ../..
echo -e 'Done\n'

cp ./tests/network/config.json ./cluster/tests/network/config.json

echo 'Generating proto files...'
protoc -I proto/ proto/node.proto --go_out=plugins=grpc:proto
echo -e 'Done\n'

echo 'Building go files...'
env GOOS=linux GOARCH=amd64 go build -o ./cluster/gossipUp ./gossipUp.go
env GOOS=linux GOARCH=amd64 go build -o ./cluster/snowballUp ./snowballUp.go
env GOOS=linux GOARCH=amd64 go build -o ./cluster/clients/gossipClient ./client/gossipClient.go
env GOOS=linux GOARCH=amd64 go build -o ./cluster/clients/snowballClient ./client/snowballClient.go
echo -e 'Done\n'