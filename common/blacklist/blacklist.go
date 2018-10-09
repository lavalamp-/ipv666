package blacklist

import "net"

type IPv6NetworkBlacklist struct {
	NetworkCount	int
	Networks		[]*net.IPNet
}

func NewNetworkBlacklist() (*IPv6NetworkBlacklist) {
	return &IPv6NetworkBlacklist{
		NetworkCount:	0,
	}
}

