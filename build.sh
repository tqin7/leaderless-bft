#!/usr/bin/env bash

rm -rf ./build

mkdir -p ./build/tests/network
mkdir -p ./build/clients

cd ./tests/network
echo 'Generating network graph...'
python3 generate_symmetric_graph_same_ip.py $1
echo -e 'Done\n'
cd ..
#echo 'Compiling tests...'
#go test -c -o ../cluster/tests/tests.test
#echo -e 'Done\n'
cd ..

cp ./tests/network/config.json ./build/tests/network/config.json

cp ./tests/network/generate_symmetric_graph_same_ip.py ./build/tests/network/generate_symmetric_graph_same_ip.py
cp ./tests/runNetwork.sh ./build/tests/runNetwork.sh
# cp ./tests/sendReqsAndCalcTime.sh ./build/tests/sendReqsAndCalcTime.sh

echo 'Generating proto files...'
protoc -I proto/ --go_out=plugins=grpc:proto --go_opt=paths=source_relative proto/node.proto
echo -e 'Done\n'

echo 'Building go files...'

# for BRC cluster use
env GOOS=linux GOARCH=amd64 go build -o ./build/gossipUp ./gossipUp.go
	env GOOS=linux GOARCH=amd64 go build -o ./build/snowballUp ./snowballUp.go
	env GOOS=linux GOARCH=amd64 go build -o ./build/pbftUp ./pbftUp.go
	# client files
	env GOOS=linux GOARCH=amd64 go build -o ./build/clients/gossipClient ./client/gossipClient.go
	env GOOS=linux GOARCH=amd64 go build -o ./build/clients/snowballClient ./client/snowballClient.go
	env GOOS=linux GOARCH=amd64 go build -o ./build/clients/pbftClient ./client/pbftClient.go
	env GOOS=linux GOARCH=amd64 go build -o ./build/clients/lbftClient ./client/lbftClient.go
	# test files
	# env GOOS=linux GOARCH=amd64 go build -o ./build/tests/testSendRequests ./tests/testSendRequests.go
	env GOOS=linux GOARCH=amd64 go build -o ./build/tests/calcExecTime ./tests/calcExecTime.go

# for local computer use
# if [[ "$OSTYPE" == "darwin"* ]]
# then
# 	# node up files
# 	env GOOS=darwin GOARCH=amd64 go build -o ./build/gossipUp ./gossipUp.go
# 	env GOOS=darwin GOARCH=amd64 go build -o ./build/snowballUp ./snowballUp.go
# 	env GOOS=darwin GOARCH=amd64 go build -o ./build/pbftUp ./pbftUp.go
# 	# client files
# 	env GOOS=darwin GOARCH=amd64 go build -o ./build/clients/gossipClient ./client/gossipClient.go
# 	env GOOS=darwin GOARCH=amd64 go build -o ./build/clients/snowballClient ./client/snowballClient.go
# 	env GOOS=darwin GOARCH=amd64 go build -o ./build/clients/pbftClient ./client/pbftClient.go
# 	env GOOS=darwin GOARCH=amd64 go build -o ./build/clients/lbftClient ./client/lbftClient.go
# 	# test files
# 	# env GOOS=darwin GOARCH=amd64 go build -o ./build//tests/testSendRequests ./tests/testSendRequests.go
# 	env GOOS=darwin GOARCH=amd64 go build -o ./build/tests/calcExecTime ./tests/calcExecTime.go
# else
# 	# node up files
# 	env GOOS=linux GOARCH=amd64 go build -o ./build/gossipUp ./gossipUp.go
# 	env GOOS=linux GOARCH=amd64 go build -o ./build/snowballUp ./snowballUp.go
# 	env GOOS=linux GOARCH=amd64 go build -o ./build/pbftUp ./pbftUp.go
# 	# client files
# 	env GOOS=linux GOARCH=amd64 go build -o ./build/clients/gossipClient ./client/gossipClient.go
# 	env GOOS=linux GOARCH=amd64 go build -o ./build/clients/snowballClient ./client/snowballClient.go
# 	env GOOS=linux GOARCH=amd64 go build -o ./build/clients/pbftClient ./client/pbftClient.go
# 	env GOOS=linux GOARCH=amd64 go build -o ./build/clients/lbftClient ./client/lbftClient.go
# 	# test files
# 	# env GOOS=linux GOARCH=amd64 go build -o ./build/tests/testSendRequests ./tests/testSendRequests.go
# 	env GOOS=linux GOARCH=amd64 go build -o ./build/tests/calcExecTime ./tests/calcExecTime.go
# fi
echo -e 'Done\n'