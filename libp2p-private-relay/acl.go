package main

import (
	"net"

	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

type aclFilter struct {
	idList []peer.ID
	ipList []net.IP
}

func (acl *aclFilter) isWhitelistedPeer(peerID peer.ID) bool {
	peers := acl.idList
	if len(peers) > 0 {
		idText := peer.Encode(peerID)
		for _, whitelistedPeer := range peers {
			if peer.Encode(whitelistedPeer) == idText {
				return true
			}
		}
	}
	return false
}

func (filter *aclFilter) isWhitelistedAddr(addr ma.Multiaddr) bool {
	addrs := filter.ipList
	if len(addrs) > 0 {
		ip, err := manet.ToIP(addr)
		if err != nil {
			return false
		}
		for _, whitelistedAddr := range addrs {
			if whitelistedAddr.Equal(ip) {
				return true
			}
		}
	}
	return false
}

func (filter *aclFilter) AllowReserve(p peer.ID, a ma.Multiaddr) bool {
	return filter.isWhitelistedPeer(p) || filter.isWhitelistedAddr(a)
}

func (filter *aclFilter) AllowConnect(src peer.ID, srcAddr ma.Multiaddr, dest peer.ID) bool {
	return true
}
