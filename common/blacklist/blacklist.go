package blacklist

import (
	"net"
	"github.com/lavalamp-/ipv666/common/addressing"
)

type NetworkBlacklist struct {
	Networks		map[string]*net.IPNet
}

func NewNetworkBlacklist(nets []*net.IPNet) (*NetworkBlacklist) {
	blacklistMap := make(map[string]*net.IPNet)
	for _, curNet := range nets {
		blacklistMap[curNet.String()] = curNet
	}
	return &NetworkBlacklist{
		Networks:	blacklistMap,
	}
}

func (blacklist *NetworkBlacklist) CleanIPList(toClean []*net.IP) ([]*net.IP) {
	var toReturn []*net.IP
	for _, curClean := range toClean {
		if !blacklist.IsIPBlacklisted(curClean) {
			toReturn = append(toReturn, curClean)
		}
	}
	return toReturn
}

func (blacklist *NetworkBlacklist) IsIPBlacklisted(toTest *net.IP) (bool) {
	for _, v := range blacklist.Networks {
		if v.Contains(*toTest) {
			return true
		}
	}
	return false
}

func (blacklist *NetworkBlacklist) Update(toAdd *net.IPNet) {
	blacklist.Networks[toAdd.String()] = toAdd
}

func ReadNetworkBlacklistFromFile(filePath string) (*NetworkBlacklist, error) {
	networks, err := addressing.ReadIPv6NetworksFromFile(filePath)
	if err != nil {
		return nil, err
	}
	return NewNetworkBlacklist(networks), nil
}

func WriteNetworkBlacklistToFile(filePath string, blacklist *NetworkBlacklist) (error) {
	toWrite := make([]*net.IPNet, len(blacklist.Networks))
	for _, v := range blacklist.Networks {
		toWrite = append(toWrite, v)
	}
	return addressing.WriteIPv6NetworksToFile(filePath, toWrite)
}
