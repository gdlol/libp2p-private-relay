package main

import (
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peerstore"
	quic "github.com/libp2p/go-libp2p-quic-transport"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/libp2p/go-tcp-transport"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

func getSavedIdentity(ps peerstore.Peerstore) (identity libp2p.Option, ok bool) {
	for _, peerID := range ps.Peers() {
		privKey := ps.PrivKey(peerID)
		if privKey != nil {
			return libp2p.Identity(privKey), true
		}
	}
	return identity, false
}

func getHostOptions(identity libp2p.Option, ps peerstore.Peerstore, listenAddrStrings []string, acl relay.ACLFilter) []libp2p.Option {
	options := []libp2p.Option{
		identity,
		libp2p.Peerstore(ps),
		libp2p.ChainOptions(
			libp2p.Transport(tcp.NewTCPTransport),
			libp2p.Transport(quic.NewTransport)),
		libp2p.ListenAddrStrings(listenAddrStrings...),
		libp2p.AddrsFactory(func(m []multiaddr.Multiaddr) []multiaddr.Multiaddr {
			return multiaddr.FilterAddrs(m, manet.IsPublicAddr)
		}),
		libp2p.DisableRelay(),
		libp2p.EnableNATService(),
		libp2p.AutoNATServiceRateLimit(60, 6, time.Minute),
		libp2p.EnableRelayService(
			relay.WithResources(relay.DefaultResources()),
			relay.WithLimit(nil),
			relay.WithACL(acl)),
		libp2p.ForceReachabilityPublic(),
	}
	return options
}
