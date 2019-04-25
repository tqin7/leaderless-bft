#!/usr/bin/env bash

mkdir -p ./cluster/tests/network
mkdir ./cluster/clients

cd ./tests/network
./generate_nodes_same_ip.sh
cd ../..

cp ./tests/network/config.json ./cluster/tests/network/config.json

protoc -I proto/ proto/node.proto --go_out=plugins=grpc:proto

go build -o ./cluster/gossipUp ./gossipUp.go
go build -o ./cluster/snowballUp ./snowballUp.go
go build -o ./cluster/clients/gossipClient ./client/gossipClient.go
go build -o ./cluster/clients/snowballClient ./client/snowballClient.go
