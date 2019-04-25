package main

import (
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/nvanbenschoten/raft-toy/peer"
	"github.com/nvanbenschoten/raft-toy/util"
	"github.com/spf13/pflag"
	"go.etcd.io/etcd/raft"
)

var raftID uint64
var raftPeers []string
var runLoad bool
var verbose bool
var recordMetrics bool

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	pflag.Uint64Var(&raftID, "id", 1, "raft.Config.ID")
	pflag.StringSliceVar(&raftPeers, "peers", []string{"localhost:1234"}, "IP addresses for raft.Peers")
	pflag.BoolVar(&runLoad, "load", false, "Propose changes to raft")
	pflag.BoolVar(&verbose, "verbose", false, "Verbose logging")
	pflag.BoolVar(&recordMetrics, "metrics", false, "Record metrics and print before exiting")
	pflag.Parse()

	// Add the set of pflags to Go's flag package so that they are usable
	// in tests and benchmarks.
	pflag.CommandLine.VisitAll(func(f *pflag.Flag) {
		switch f.Value.Type() {
		case "bool":
			def, err := strconv.ParseBool(f.DefValue)
			if err != nil {
				panic(err)
			}
			flag.Bool(f.Name, def, f.Usage)
		default:
			flag.String(f.Name, f.DefValue, f.Usage)
		}
	})

	if !verbose {
		util.DisableRaftLogging()
	}
}

func cfgFromFlags() peer.PeerConfig {
	cfg := peer.PeerConfig{ID: raftID}
	if cfg.ID == 0 {
		log.Fatalf("invalid ID (%d); must be > 0", cfg.ID)
	}
	if len(raftPeers) < int(cfg.ID) {
		log.Fatalf("missing own ID (%d) in peers (%v)", cfg.ID, raftPeers)
	}
	cfg.Peers = make([]raft.Peer, len(raftPeers))
	cfg.PeerAddrs = make(map[uint64]string, len(raftPeers))
	for i, addr := range raftPeers {
		pID := uint64(i + 1)
		cfg.Peers[i].ID = pID
		cfg.PeerAddrs[pID] = addr
	}
	cfg.SelfAddr = cfg.PeerAddrs[cfg.ID]
	return cfg
}
