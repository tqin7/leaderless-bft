# Leaderless-BFT
## Overview
Leaderless Byzantine Fault Tolerance is a new blockchain consensus algorithm (also applicable to general distributed systems) that tolerates Byzantine faults while being leaderless and deterministic. It's based upon pBFT (practical Byzantine Fault Tolerance) and Snowball. 
* pBFT is deterministic but relies on the notion of a primary/leader to determine the total order of requests (through pre-prepare, prepare, commit phases)
* Snowball is leaderless but probabilistic as the algorithm proceeds by taking random samples of the network

We are looking to replace the primary in pBFT with Snowball algorithm, i.e. use Snowball (which pulls information from the entire network) to simulate the primary's tasks such as assigning total order, so that we arrive at a leaderless and deterministic algorithm that has a better balance between security and efficiency.

## Instructions
To run the system, run `build.sh` to generate a folder `cluster` that contains the necessary `[layer]Up`, `[layer]Client`, and network config files.

To run a layer, simply run its corresponding program in `cluster`.
To interact with a layer that's brought up, run the corresponding client program, which creates an interactive channel that sends messages/requests to the nodes.

## Technical Implementation
The main programming language is Go. The peer-to-peer communication layer implements gossip protocols. The consensus layer lies on top of the p2p layer and implements Leaderless BFT.

## Directories
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
