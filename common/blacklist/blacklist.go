package blacklist

import (
  "encoding/binary"
	"net"
	"github.com/lavalamp-/ipv666/common/addressing"
	"log"
)

type blacklistPlaceholder struct {}

type ipNets struct {
	nets map[[2]uint64]struct{}
}

type NetworkBlacklist struct {
	Networks			map[string]*net.IPNet
	Nets          map[int]*ipNets
	Masks         map[int]*[2]uint64
	checks				map[uint64]blacklistPlaceholder
}

func NewNetworkBlacklist(nets []*net.IPNet) (*NetworkBlacklist) {
	toReturn := &NetworkBlacklist{
		Networks:	make(map[string]*net.IPNet),
		Nets:     make(map[int]*ipNets),
		checks:		make(map[uint64]blacklistPlaceholder),
	}

  // Build the per-length masks
  toReturn.Masks = map[int]*[2]uint64{}
  for l, _ := range nets {
    toReturn.Masks[l] = &[2]uint64{}
    if l <= 64 {
      toReturn.Masks[l][1] = 0
      for x := 0; x < l; x++ {
        toReturn.Masks[l][0] |= (1 << (63 - uint64(x)))
      }
    } else {
      for x := 0; x < l; x++ {
        toReturn.Masks[l][0] |= 1 << (63 - uint64(x))
      }
      for x := 64; x < l; x++ {
        toReturn.Masks[l][1] |= 1 << (127 - uint64(x))
      }
    }
  }

	for _, curNet := range nets {
		toReturn.AddNetwork(*curNet)
	}
	return toReturn
}

func (blacklist *NetworkBlacklist) AddNetwork(toAdd net.IPNet) {

	// Compute the length of the network
	netLen := 0
	for x := 15; x >= 0; x-- {
		if toAdd.Mask[x] > 0 {
			netLen = (x-1)*8
			for y := 0; y < 8; y++ {
				mask := byte(1 << uint(y))
				if toAdd.Mask[x] & mask == mask {
					netLen += (8-y)
					break
				}
			}
			break
		}
	}

  // New len?
  if _, ok := blacklist.Nets[netLen]; !ok {
    blacklist.Nets[netLen] = &ipNets{}
    blacklist.Nets[netLen].nets = map[[2]uint64]struct{}{}
  }

  // Add this network
  ip := [2]uint64{}
  ip[0] = binary.LittleEndian.Uint64(toAdd.IP[0:8])
  ip[1] = binary.LittleEndian.Uint64(toAdd.IP[8:16])
  blacklist.Nets[netLen].nets[ip] = struct{}{}

 //  // Legacy logic
	// networkString := addressing.GetBaseAddressString(&toAdd)
	// if _, ok := blacklist.Networks[networkString]; !ok {
	// 	blacklist.Networks[networkString] = &toAdd
	// 	netBytes := addressing.GetFirst64BitsOfNetwork(&toAdd)
	// 	blacklist.checks[netBytes] = blacklistPlaceholder{}
	// }
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

  ipUints := [2]uint64{}
  ipUints[0] = binary.LittleEndian.Uint64((*toTest)[0:8])
  ipUints[1] = binary.LittleEndian.Uint64((*toTest)[8:16])

  // Check the IP against each network length
  for l, _ := range blacklist.Nets {
    ipMask := [2]uint64{}
    if l <= 64 {
      ipMask[1] = 0
      ipMask[0] = ipUints[0] & blacklist.Masks[l][0]
    } else {
      ipMask[1] = ipUints[1] & blacklist.Masks[l][1]
      ipMask[0] = ipUints[0] & blacklist.Masks[l][0]
    }

    if _, ok := blacklist.Nets[l].nets[ipMask]; ok {
      return true
    }
  }

  return false

	// first64 := addressing.GetFirst64BitsOfIP(toTest)
	// _, ok := blacklist.checks[first64]
	// return ok
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
