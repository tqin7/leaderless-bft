package tests

import (
	"testing"
	"google.golang.org/grpc"
	"context"
	"github.com/spockqin/leaderless-bft/tests/network"
	"time"
	pb "github.com/spockqin/leaderless-bft/proto"
	"math/rand"
	"strconv"
)

func TestSendAndRcvData(t *testing.T) {

}

// throughput if all requests are sent through the same connection
func TestGossipThroughputSameConn(t *testing.T) {
	t.Log("Measuring Gossip Layer Throughput...")

	var nodes network.Nodes
	network.ReadNetworkConfig(&nodes, "./network/config.json")

	var gossipers []string
	for _, node := range nodes.Nodes {
		gossipers = append(gossipers, node.Ip)
	}

	t.Log("Timestamp right before first dialing: ",
		time.Now().Format("2006-01-01 15:04:05 .000"))

	mainIp := gossipers[0]
	mainConn, err := grpc.Dial(mainIp, grpc.WithInsecure())
	if err != nil {
		t.Error("Cannot establish TCP connection with gossiper")
		return
	}
	defer mainConn.Close()

	mainClient := pb.NewGossipClient(mainConn)

	for i := 0; i < 10; i++ {
		req := []byte(strconv.Itoa(i))
		mainClient.Push(context.Background(), &pb.ReqBody{Body: req})
	}
}

// throughput if a new connection is established for each request
func TestGossipThroughputRandomConn(t *testing.T) {
	t.Log("Measuring Gossip Layer Throughput...")

	var nodes network.Nodes
	network.ReadNetworkConfig(&nodes, "./network/config.json")

	var gossipers []string
	for _, node := range nodes.Nodes {
		gossipers = append(gossipers, node.Ip)
	}
	n := len(gossipers)

	now := time.Now()
	t.Log("Timestamp right before first dialing: ", now.Format("2006-01-01 15:04:05"))

	i := 0
	for i < 1000 {
		conn, err := grpc.Dial(gossipers[rand.Intn(n)], grpc.WithInsecure())
		if err != nil {
			continue
		}
		client := pb.NewGossipClient(conn)
		req := []byte(string(i))
		client.Push(context.Background(), &pb.ReqBody{Body: req})
		conn.Close()
		i++
	}
}
