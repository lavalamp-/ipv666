package blacklist

import (
	"bufio"
	"encoding/binary"
	"github.com/lavalamp-/ipv666/common/addressing"
	"log"
	"net"
	"os"
	"sort"
)

type ipNets struct {
	nets map[[2]uint64]struct{}
}

type NetworkBlacklist struct {
	nets			map[int]*ipNets
	masks			map[int]*[2]uint64
	maskLengths		[]int
	count			int
}

func NewNetworkBlacklist(nets []*net.IPNet) (*NetworkBlacklist) {
	toReturn := &NetworkBlacklist{
		nets:			make(map[int]*ipNets),
		masks:			make(map[int]*[2]uint64),
		maskLengths:	[]int{},
		count:			0,
	}

	// Build the per-length masks
	for l := 0; l < 129; l++ {
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

	netLen, _ := toAdd.Mask.Size()

	// New len?
	if _, ok := blacklist.nets[netLen]; !ok {
		blacklist.nets[netLen] = &ipNets{}
		blacklist.nets[netLen].nets = map[[2]uint64]struct{}{}
		blacklist.updateMaskLengths(netLen)
	}

	// Add this network
	ip := [2]uint64{}
	ip[0] = binary.BigEndian.Uint64(toAdd.IP[0:8])
	ip[1] = binary.BigEndian.Uint64(toAdd.IP[8:16])
	blacklist.nets[netLen].nets[ip] = struct{}{}

	blacklist.count++
	return true

}

func (blacklist *NetworkBlacklist) CleanIPList(toClean []*net.IP, emitFreq int) ([]*net.IP) {
	var toReturn []*net.IP
	for i, curClean := range toClean {
		if i % emitFreq == 0 && i != 0 {
			log.Printf("Cleaning entry %d out of %d.", i, len(toClean))
		}
		if !blacklist.IsIPBlacklisted(curClean) {
			toReturn = append(toReturn, curClean)
		}
	}
	return toReturn
}

func (blacklist *NetworkBlacklist) updateMaskLengths(maskLength int) () {
	blacklist.maskLengths = append(blacklist.maskLengths, maskLength)
	sort.Ints(blacklist.maskLengths)
}

func (blacklist *NetworkBlacklist) getNetworkFromAddress(toTest *net.IP) ([2]uint64, int, bool) {

	ipUints := [2]uint64{
		binary.BigEndian.Uint64((*toTest)[0:8]),
		binary.BigEndian.Uint64((*toTest)[8:16]),
	}

	// Check the IP against each network length
	for _, maskLength := range blacklist.maskLengths {
		ipMask := [2]uint64{}
		if maskLength <= 64 {
			ipMask[1] = 0
			ipMask[0] = ipUints[0] & blacklist.masks[maskLength][0]
		} else {
			ipMask[1] = ipUints[1] & blacklist.masks[maskLength][1]
			ipMask[0] = ipUints[0] & blacklist.masks[maskLength][0]
		}

		if _, ok := blacklist.nets[maskLength].nets[ipMask]; ok {
			return ipMask, maskLength, true
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

func (blacklist *NetworkBlacklist) GetBlacklistingNetworkFromIP(toTest *net.IP) (*net.IPNet) {
	uints, length, found := blacklist.getNetworkFromAddress(toTest)
	if !found {
		return nil
	} else {
		return addressing.GetNetworkFromUints(uints, uint8(length))
	}
}

func (blacklist *NetworkBlacklist) GetBlacklistingNetworkFromNetwork(toTest *net.IPNet) (*net.IPNet) {
	base, top := addressing.GetBorderAddressesFromNetwork(toTest)
	baseNetwork := blacklist.GetBlacklistingNetworkFromIP(base)
	if baseNetwork == nil {
		return nil
	}
	topNetwork := blacklist.GetBlacklistingNetworkFromIP(top)
	if topNetwork == nil {
		return nil
	} else {
		return topNetwork
	}
}

func (blacklist *NetworkBlacklist) GetCount() (int) {
	return blacklist.count
}

func (blacklist *NetworkBlacklist) GetMaskLengths() ([]int) {
	return blacklist.maskLengths
}

func (blacklist *NetworkBlacklist) GetNetworks() ([]*net.IPNet) {
	var toReturn []*net.IPNet
	for _, maskLength := range blacklist.maskLengths {
		for curNet := range blacklist.nets[maskLength].nets {
			toReturn = append(toReturn, addressing.GetNetworkFromUints(curNet, uint8(maskLength)))
		}
	}
	return toReturn
}

func (blacklist *NetworkBlacklist) Clean(emitFreq int) (int) {
	var newNetworks []*net.IPNet
	numCleaned := 0
	loopCount := 0

	for _, maskLength := range blacklist.maskLengths {
		for curNet := range blacklist.nets[maskLength].nets {
			ipNet := addressing.GetNetworkFromUints(curNet, uint8(maskLength))
			blacklistNetwork := blacklist.GetBlacklistingNetworkFromNetwork(ipNet)
			if blacklistNetwork != nil {
				ones, _ := blacklistNetwork.Mask.Size()
				if ones == maskLength {
					newNetworks = append(newNetworks, blacklistNetwork)
				} else {
					numCleaned++
				}
			}
			loopCount++
			if loopCount % emitFreq == 0 {
				log.Printf("Processing %d out of %d in blacklist cleaning.", loopCount, blacklist.count)
			}
		}
	}

	blacklist.nets = make(map[int]*ipNets)
	blacklist.count = 0
	blacklist.maskLengths = nil

	blacklist.AddNetworks(newNetworks)
	return numCleaned
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
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	defer file.Close()

	writeBytes := make([]byte, 8)
	for maskLength, ipnets := range blacklist.nets {
		for netBytes := range ipnets.nets {
			binary.BigEndian.PutUint64(writeBytes, netBytes[0])
			writer.Write(writeBytes)
			binary.BigEndian.PutUint64(writeBytes, netBytes[1])
			writer.Write(writeBytes)
			writer.Write([]byte{uint8(maskLength)})
		}
	}
	writer.Flush()
	
	return nil
}