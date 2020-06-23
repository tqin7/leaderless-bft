# Leaderless-BFT
1. [Project Overview](##Project-Overview)
2. [Instructions](##Instructions)
3. [Techincal Implementation](##Technical-Implementation)
4. [Directory Overview](##Directory-Overview)
5. [Testing Guide](##Testing-Guide)

## Project Overview
Leaderless Byzantine Fault Tolerance is a new blockchain consensus algorithm (also applicable to general distributed systems) that tolerates Byzantine faults while being leaderless and deterministic. It's based upon pBFT (practical Byzantine Fault Tolerance) and Snowball. 
* pBFT is deterministic but relies on the notion of a primary/leader to determine the total order of requests (through pre-prepare, prepare, commit phases)
* Snowball is leaderless but probabilistic as the algorithm proceeds by taking random samples of the network

We are looking to replace the primary in pBFT with Snowball algorithm, i.e. use Snowball (which pulls information from the entire network) to simulate the primary's tasks such as assigning total order, so that we arrive at a leaderless and deterministic algorithm that has a better balance between security and efficiency.

## Instructions
Install protoc: http://google.github.io/proto-lens/installing-protoc.html
Install protoc-gen-go: https://github.com/golang/protobuf

To run the system, run `build.sh` to generate a folder `cluster` that contains the necessary `[layer]Up`, `[layer]Client`, and network config files.

To run a layer, simply run its corresponding program in `cluster`.
To interact with a layer that's brought up, run the corresponding client program, which creates an interactive channel that sends messages/requests to the nodes.

<b>Example - Gossip Layer</b> 
1. Go to `./tests/network/generate_nodes_same_ip.sh` and change the top line to `numOfNodes=x` where x is your desired number of nodes in the network
2. Go back to root directory and run `./build.sh`, which should generate network graph, proto files, and build go files.
3. Go to `./cluster` and run `./gossipUp`. You should see each node being brought up.
4. To interact with the gossip nodes, go to `./cluster/clients` and run `./gossipClient`. You should see the terminal requesting for messages.
    * send message `get reqs` to see what requests each node knows
    * send any other message (e.g. `req1`) to send a request to the network

## Technical Implementation
The main programming language is Go. The peer-to-peer communication layer implements gossip protocols. The consensus layer lies on top of the p2p layer and implements Leaderless BFT.

## Directory Overview
`client/`: mainly for interactively testing and measuring the services
* `gossipClient.go`: a client that interacts with the gossip service
* `snowballClient.go`: a client that interacts with the snowball service

`node/`: service implementation
* `gossiper.go`: business logic of gossip layer, implementation of the gossip gRPC service
* `snower.go`: business logic of snowball layer, implementation of the snowball gRPC service

`proto/`: service definition
* `node.proto`: defines gRPC services (`Gossip`, `Snowball`), i.e. inter-machine communication methods

`tests/network/`: testing infrastructure
* `alias_up.sh`, `alias_down.sh`: bash scripts that bring up aliases of an IP address to simulate different machines during testing
* `generate_nodes_diff_ips.sh`: bash script that randomly generates a connected network whose nodes have different IP addresses and store its configuration in `config.json`
* `generate_nodes_same_ip.sh`: bash script that randomly generates a connected network whose nodes have the same IP address but different ports and store its configuration in `config.json`
* `config.json`: stores the network configuration in terms of each node and its neighbors. This file is read by `gossipUp.go` and `snowballUp.go`.

`types/`: initialization of system-wide variables

`utils/`: util functions

`build.sh`: generates network graph, proto files, and build go files that are necessary to run the system

`gossipUp.go`: bring up the gossip layer according to the generated network configuration

`snowballUp.go`: bring up the snowball layer according to the generated network configuration

## Testing Guide
`./build.sh` generates a `build` folder that contains all files necessary for testing. Specifically, this shell script
* generates network configuration
* generates protobuf file(s)
* builds Go files
* sets up testing framework

Testing of our system consists of two steps: 
   1. spinning up the network
   2. interacting with the network through client endpoints

For step 1, change directory to `./build/tests`. Run the command `./runNetwork.sh [networkType] [networkSize]` where
* `networkType` specifies node type and can be any of the four: `gossipUp`, `snowballUp`, `pbftUp`, `lbftUp`
* `networkSize` is an integer that specifies the number of nodes in the network and can be for example 50, 232, etc.

For step 2, change directory to `./build/clients`. Run the corresponding client endpoint: for example, if in the previous step, `snowballUp` is chosen, run `./snowballClient` in this step. Then you can easily send requests to the network through a command line interface. To test for latency, send `latency same` through the interface, a command that tells the client endpoint to send 100 requests in a row to the network. To test for throughput (only available for pbft and lbft), send `throughput` through the interface, which will keep sending requests and you can examine when and whether network has stabilized through dynamic logs from step 1. For other functionalities of client interfaces, check out `./client/[nodeType]Client.go` files, from which these client binaries are compiled from.

During testing, standard output logs will be logged in file `./build/tests/log.txt` and to calculate total execution time, change directory to `/build/tests ` and run `./calcExecTime log.txt`, which will output total execution time to `./build/tests/execTime.txt`.
