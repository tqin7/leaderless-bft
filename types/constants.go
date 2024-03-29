package types

import "time"

const (
	SNOWBALL_SAMPLE_ROUNDS = 10
	// to check # of max sockets open at once, run "ulimit -n"
	MAX_SOCKETS = 10000
	GRPC_TIMEOUT = time.Second * 40
	//GRPC_TIMEOUT = time.Nanosecond * 5
)

const (
	CONN_TCP_PORT = "7777"
	CONN_UDP_PORT = "7778" // TCP and UDP ports can be the same
)
