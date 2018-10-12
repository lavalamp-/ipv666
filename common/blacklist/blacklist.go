package blacklist

import (
	"net"
	"github.com/lavalamp-/ipv666/common/addressing"
	"encoding/binary"
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

	// Blacklist net.IPNet's to uint64's
	blnets := []uint64{}
	for _, ipnet:= range blacklist.Networks {
		n := binary.LittleEndian.Uint64((*ipnet).IP[:8])
		blnets = append(blnets, n)
	}

	for _, curClean := range toClean {
		n := binary.LittleEndian.Uint64((*curClean)[:8])
		found := false
		for _, v := range blnets {
			if v == n {
				found = true
				break
			}
		}

		// if !blacklist.IsIPBlacklisted(curClean) {
		if !found {
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
	toWrite := []*net.IPNet{}
	for _, v := range blacklist.Networks {
		toWrite = append(toWrite, v)
	}
	return addressing.WriteIPv6NetworksToFile(filePath, toWrite)
}
