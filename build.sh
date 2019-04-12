#!/usr/bin/env bash

cd ./tests/network
./generate_nodes.sh
cd ../..
go build ./main.go
cd ./client
go build ./client.go
cd ..
cp ./main ./cluster
cp ./client/client ./cluster