package types

import "time"

const (
	SNOWBALL_SAMPLE_ROUNDS = 5
	MAX_SOCKETS = 25
	GRPC_TIMEOUT = time.Second * 40
	//GRPC_TIMEOUT = time.Nanosecond * 5
)

const (
	CONN_TCP_PORT = "7777"
	CONN_UDP_PORT = "7778" // TCP and UDP ports can be the same
)
