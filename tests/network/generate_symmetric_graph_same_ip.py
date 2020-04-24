#!/usr/bin/env python

import socket
import argparse
import random
import math
from collections import defaultdict

IP = "127.0.0.1"
NODE_STR = "{{\"ip\": \"{ip}\", \"peers\": \"{peers}\"}}"

def port_in_use(port):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        return s.connect_ex((IP, port)) == 0

parser = argparse.ArgumentParser(description="generate symmetric network graph")
parser.add_argument("size", help="size of network", type=int, default=10)

args = parser.parse_args()
numOfNodes = args.size
maxNeighborSize = int(math.floor(math.sqrt(numOfNodes)))

addrs = []
port = 29999
while len(addrs) < numOfNodes:
    port += 1
    addrs.append(IP + ":" + str(port))
    
    # if port_in_use(port):
    #     continue
    # else:
    #     addrs.append(IP + ":" + str(port))

network = defaultdict(list)
for addr in addrs:
    neighbors = random.sample(addrs, maxNeighborSize)
    network[addr].extend(neighbors)
    for neighbor in neighbors:
        network[neighbor].append(addr)

def node_to_json(addr):
    return NODE_STR.format(ip=addr, peers=','.join(network[addr]))
jsonStr = ','.join(map(node_to_json, addrs))

jsonStr = "{{\"nodes\":[{}]}}".format(jsonStr)

with open("config.json", 'w') as f:
    f.write(jsonStr)
