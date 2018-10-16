package blacklist

import (
	"net"
	"github.com/lavalamp-/ipv666/common/addressing"
	"log"
)

type blacklistPlaceholder struct {}

type NetworkBlacklist struct {
	Networks			map[string]*net.IPNet
	checks				map[uint64]blacklistPlaceholder
}

func NewNetworkBlacklist(nets []*net.IPNet) (*NetworkBlacklist) {
	toReturn := &NetworkBlacklist{
		Networks:	make(map[string]*net.IPNet),
		checks:		make(map[uint64]blacklistPlaceholder),
	}
	for _, curNet := range nets {
		toReturn.AddNetwork(*curNet)
	}
	return toReturn
}

func (blacklist *NetworkBlacklist) AddNetwork(toAdd net.IPNet) {
	networkString := addressing.GetBaseAddressString(&toAdd)
	if _, ok := blacklist.Networks[networkString]; !ok {
		blacklist.Networks[networkString] = &toAdd
		netBytes := addressing.GetFirst64BitsOfNetwork(&toAdd)
		blacklist.checks[netBytes] = blacklistPlaceholder{}
	}
}

func (blacklist *NetworkBlacklist) CleanIPList(toClean []*net.IP, emitFreq int) ([]*net.IP) {
	var toReturn []*net.IP
	for i, curClean := range toClean {
		if i % emitFreq == 0 {
			log.Printf("Cleaning entry %d out of %d.", i, len(toClean))
		}
		if !blacklist.IsIPBlacklisted(curClean) {
			toReturn = append(toReturn, curClean)
		}
	}
	return toReturn
}

func (blacklist *NetworkBlacklist) IsIPBlacklisted(toTest *net.IP) (bool) {
	first64 := addressing.GetFirst64BitsOfIP(toTest)
	_, ok := blacklist.checks[first64]
	return ok
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
