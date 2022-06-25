package main

import (
	"context"
	"fmt"
	"log"
	"net"

	leveldb "github.com/ipfs/go-ds-leveldb"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-peerstore/pstoreds"
)

func loadACL(config Config) (*aclFilter, error) {
	idList := make([]peer.ID, len(config.WhitelistPeers))
	for _, s := range config.WhitelistPeers {
		id, err := peer.IDFromString(s)
		if err != nil {
			return nil, fmt.Errorf("error parsing peer id: %w", err)
		}
		idList = append(idList, id)
	}
	ipList := make([]net.IP, len(config.WhitelistAddrs))
	for _, s := range config.WhitelistAddrs {
		ip := net.ParseIP(s)
		if ip == nil {
			return nil, fmt.Errorf("error parsing IP address: %s", s)
		}
		ipList = append(ipList, ip)
	}
	filter := &aclFilter{
		idList: idList,
		ipList: ipList,
	}
	return filter, nil
}

func main() {
	logging.SetLogLevel("relay", "DEBUG")

	config, err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Config: %+v\n", config)

	acl, err := loadACL(config)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	dbPath := "/root/data/datastore.db"
	ds, err := leveldb.NewDatastore(dbPath, nil)
	if err != nil {
		log.Fatalln(fmt.Errorf("error creating DataStore: %w", err))
	}
	defer ds.Close()
	ps, err := pstoreds.NewPeerstore(ctx, ds, pstoreds.DefaultOpts())
	if err != nil {
		log.Fatalln(fmt.Errorf("error creating PeerStore: %w", err))
	}
	defer ps.Close()

	identity, ok := getSavedIdentity(ps)
	if ok {
		log.Println("Found saved identity.")
	} else {
		log.Println("Creating new identity.")
		identity = libp2p.RandomIdentity
	}

	hostOptions := getHostOptions(identity, ps, config.ListenAddrStrings, acl)
	host, err := libp2p.New(hostOptions...)
	if err != nil {
		log.Fatalf("Error creating host: %v\n", err)
	}
	defer host.Close()

	selfAddrs, err := peer.AddrInfoToP2pAddrs(&peer.AddrInfo{
		ID: host.ID(),
	})
	if err != nil {
		log.Fatalln(fmt.Errorf("error getting self addresses: %w", err))
	}
	log.Printf("Self addresses: %v\n", selfAddrs)

	<-ctx.Done()
}
