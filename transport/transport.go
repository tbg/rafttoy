package transport

import "go.etcd.io/etcd/raft/raftpb"

// Transport handles RPC messages for Raft coordination.
type Transport interface {
	Init(addr string, peers map[uint64]string)
	Serve(RaftHandler)
	Send(epoch int32, msgs []raftpb.Message)
	Close()
}

// RaftHandler is an object capable of accepting incoming Raft messages.
type RaftHandler interface {
	HandleMessage(epoch int32, msg *raftpb.Message)
}
