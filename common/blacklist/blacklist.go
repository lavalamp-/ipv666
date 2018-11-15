package blacklist

import (
	"encoding/binary"
	"net"
	"github.com/lavalamp-/ipv666/common/addressing"
	"log"
	"os"
)

type blacklistPlaceholder struct {}

type ipNets struct {
	nets map[[2]uint64]struct{}
}

type NetworkBlacklist struct {
	Networks map[string]*net.IPNet
	nets     map[int]*ipNets
	masks    map[int]*[2]uint64
	checks   map[uint64]blacklistPlaceholder
}

func NewNetworkBlacklist(nets []*net.IPNet) (*NetworkBlacklist) {
	toReturn := &NetworkBlacklist{
		Networks: make(map[string]*net.IPNet),
		nets:     make(map[int]*ipNets),
		checks:   make(map[uint64]blacklistPlaceholder),
		masks:    make(map[int]*[2]uint64),
	}

	// Build the per-length masks
	for l := 0; l < 128; l++ {
		toReturn.masks[l] = &[2]uint64{}
		if l <= 64 {
			toReturn.masks[l][1] = 0
			for x := 0; x < l; x++ {
				toReturn.masks[l][0] |= 1 << (63 - uint64(x))
			}
		} else {
			for x := 0; x < l; x++ {
				toReturn.masks[l][0] |= 1 << (63 - uint64(x))
			}
			for x := 64; x < l; x++ {
				toReturn.masks[l][1] |= 1 << (127 - uint64(x))
			}
		}
	}

	toReturn.AddNetworks(nets)

	return toReturn
}

func (blacklist *NetworkBlacklist) AddNetworks(toAdd []*net.IPNet) (int, int) {
	addedCount, skippedCount := 0, 0
	for _, curNet := range toAdd {
		added := blacklist.AddNetwork(curNet)
		if added {
			addedCount++
		} else {
			skippedCount++
		}
	}
	return addedCount, skippedCount
}

func (blacklist *NetworkBlacklist) AddNetwork(toAdd *net.IPNet) (bool) {

	if blacklist.IsNetworkBlacklisted(toAdd) {
		return false
	}

	networkString := addressing.GetBaseAddressString(toAdd)
	//TODO now that we have network membership check (above) I think we can strip out this string set check
	if _, ok := blacklist.Networks[networkString]; !ok {
		blacklist.Networks[networkString] = toAdd

		// Compute the length of the network
		netLen := 0
		for x := 15; x >= 0; x-- {
			if toAdd.Mask[x] > 0 {
				netLen = x*8
				for y := 0; y < 8; y++ {
					mask := byte(1 << uint(y))
					if toAdd.Mask[x] & mask == mask {
						netLen += 8 - y
						break
					}
				}
				break
			}
		}

		// New len?
		if _, ok := blacklist.nets[netLen]; !ok {
			blacklist.nets[netLen] = &ipNets{}
			blacklist.nets[netLen].nets = map[[2]uint64]struct{}{}
		}

		// Add this network
		ip := [2]uint64{}
		ip[0] = binary.BigEndian.Uint64(toAdd.IP[0:8])
		ip[1] = binary.BigEndian.Uint64(toAdd.IP[8:16])
		blacklist.nets[netLen].nets[ip] = struct{}{}

	}

	return true

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

func (blacklist *NetworkBlacklist) getNetworkFromAddress(toTest *net.IP) ([2]uint64, int, bool) {

	ipUints := [2]uint64{
		binary.BigEndian.Uint64((*toTest)[0:8]),
		binary.BigEndian.Uint64((*toTest)[8:16]),
	}

	// Check the IP against each network length
	for l := range blacklist.nets {
		ipMask := [2]uint64{}
		if l <= 64 {
			ipMask[1] = 0
			ipMask[0] = ipUints[0] & blacklist.masks[l][0]
		} else {
			ipMask[1] = ipUints[1] & blacklist.masks[l][1]
			ipMask[0] = ipUints[0] & blacklist.masks[l][0]
		}

		if _, ok := blacklist.nets[l].nets[ipMask]; ok {
			return ipMask, l, true
		}
	}

	return [2]uint64{0,0}, -1, false

}

func (blacklist *NetworkBlacklist) IsNetworkBlacklisted(toTest *net.IPNet) (bool) {
	//TODO make sure this logic isn't flawed. I'm fairly certain that if both the top and bottom of the network
	// are blacklisted then the network is, in its entirety, blacklisted as well.
	top, bottom := addressing.GetBorderAddressesFromNetwork(toTest)
	return blacklist.IsIPBlacklisted(top) && blacklist.IsIPBlacklisted(bottom)
}

func (blacklist *NetworkBlacklist) IsIPBlacklisted(toTest *net.IP) (bool) {
	_, _, found := blacklist.getNetworkFromAddress(toTest)
	return found
}

func (blacklist *NetworkBlacklist) GetBlacklistingNetwork(toTest *net.IP) (*net.IPNet) {
	uints, length, found := blacklist.getNetworkFromAddress(toTest)
	if !found {
		return nil
	} else {
		return addressing.GetNetworkFromUints(uints, length)
	}
}

func ReadNetworkBlacklistFromFile(filePath string) (*NetworkBlacklist, error) {
	log.Printf("Loading blacklist from path '%s'.", filePath)
	networks, err := addressing.ReadIPv6NetworksFromFile(filePath)
	log.Printf("Read %d networks from file '%s'.", len(networks), filePath)
	if err != nil {
		return nil, err
	}
	log.Printf("Creating blacklist from %d networks.", len(networks))
	return NewNetworkBlacklist(networks), nil
}

func WriteNetworkBlacklistToFile(filePath string, blacklist *NetworkBlacklist) (error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, network := range blacklist.Networks {
		file.Write(network.IP)
		ones, _ := network.Mask.Size()
		length := uint8(ones)
		file.Write([]byte{length})
	}
	return nil
}