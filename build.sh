#!/usr/bin/env bash

mkdir -p ./cluster/tests/network
mkdir ./cluster/clients

cd ./tests/network
echo 'Generating network graph...'
./generate_nodes_same_ip.sh
echo -e 'Done\n'
cd ..
#echo 'Compiling tests...'
#go test -c -o ../cluster/tests/tests.test
#echo -e 'Done\n'
cd ..

cp ./tests/network/config.json ./cluster/tests/network/config.json

echo 'Generating proto files...'
protoc -I proto/ proto/node.proto --go_out=plugins=grpc:proto
echo -e 'Done\n'

echo 'Building go files...'
if [[ "$OSTYPE" == "darwin"* ]]
then
	env GOOS=darwin GOARCH=amd64 go build -o ./cluster/gossipUp ./gossipUp.go
	env GOOS=darwin GOARCH=amd64 go build -o ./cluster/snowballUp ./snowballUp.go
	env GOOS=darwin GOARCH=amd64 go build -o ./cluster/pbftUp ./pbftUp.go
	env GOOS=darwin GOARCH=amd64 go build -o ./cluster/lbftUp ./lbftUp.go
	env GOOS=darwin GOARCH=amd64 go build -o ./cluster/clients/gossipClient ./client/gossipClient.go
	env GOOS=darwin GOARCH=amd64 go build -o ./cluster/clients/snowballClient ./client/snowballClient.go
	env GOOS=darwin GOARCH=amd64 go build -o ./cluster/clients/pbftClient ./client/pbftClient.go
	env GOOS=darwin GOARCH=amd64 go build -o ./cluster/clients/lbftClient ./client/lbftClient.go

else
	env GOOS=linux GOARCH=amd64 go build -o ./cluster/gossipUp ./gossipUp.go
	env GOOS=linux GOARCH=amd64 go build -o ./cluster/snowballUp ./snowballUp.go
	env GOOS=linux GOARCH=amd64 go build -o ./cluster/pbftUp ./pbftUp.go
	env GOOS=linux GOARCH=amd64 go build -o ./cluster/lbftUp ./lbftUp.go
	env GOOS=linux GOARCH=amd64 go build -o ./cluster/clients/gossipClient ./client/gossipClient.go
	env GOOS=linux GOARCH=amd64 go build -o ./cluster/clients/snowballClient ./client/snowballClient.go
	env GOOS=linux GOARCH=amd64 go build -o ./cluster/clients/pbftClient ./client/pbftClient.go
	env GOOS=linux GOARCH=amd64 go build -o ./cluster/clients/lbftClient ./client/lbftClient.go
fi
echo -e 'Done\n'