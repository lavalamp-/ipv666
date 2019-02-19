package modeling

import (
	"github.com/lavalamp-/ipv666/internal/addressing"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/spf13/viper"
	"net"
)

type BinaryAddressContainer struct {
	addresses			map[uint64][]uint64
	sortedHighKeys		[]uint64
}

func ContainerFromAddrs(toProcess []*net.IP) *BinaryAddressContainer {
	toReturn := &BinaryAddressContainer{
		addresses:			make(map[uint64][]uint64),
		sortedHighKeys:		[]uint64{},
	}
	toReturn.AddIPs(toProcess, viper.GetInt("LogLoopEmitFreq"))
	return toReturn
}

func EmptyContainer() *BinaryAddressContainer {
	return &BinaryAddressContainer{
		addresses:			make(map[uint64][]uint64),
		sortedHighKeys:		[]uint64{},
	}
}

func (container *BinaryAddressContainer) Size() int { //TODO too computationally expensive
	return len(container.GetAllIPs())
}

func (container *BinaryAddressContainer) AddIP(toAdd *net.IP) bool {
	first, second := addressing.AddressToUints(*toAdd)
	var added = false
	var secondAdded = false
	if _, ok := container.addresses[first]; !ok {
		container.addresses[first] = []uint64{}
		container.sortedHighKeys, _ = insert(container.sortedHighKeys, first)
		added = true
	}
	container.addresses[first], secondAdded = insert(container.addresses[first], second)
	return added || secondAdded
}

func (container *BinaryAddressContainer) AddIPs(toAdd []*net.IP, emitFreq int) (int, int) { //TODO get rid of emit freq
	added, skipped := 0, 0
	for i, curAdd := range toAdd {
		if i % emitFreq == 0 {
			logging.Infof("Adding address %d out of %d to BinaryAddressContainer.", i, len(toAdd))
		}
		wasAdded := container.AddIP(curAdd)
		if wasAdded {
			added += 1
		} else {
			skipped += 1
		}
	}
	return added, skipped
}

func (container *BinaryAddressContainer) GetAllIPs() []*net.IP {
	var toReturn []*net.IP
	processed := 0
	for k, v := range container.addresses {
		for _, curLower := range v {
			processed++
			if processed % viper.GetInt("LogLoopEmitFreq") == 0 {
				logging.Infof("Dumping address %d from BinaryAddressContainer.", processed)
			}
			toReturn = append(toReturn, addressing.UintsToAddress(k, curLower))
		}
	}
	return toReturn
}

func (container *BinaryAddressContainer) ContainsIP(toCheck *net.IP) bool {
	first, second := addressing.AddressToUints(*toCheck)
	if val, ok := container.addresses[first]; !ok {
		return false
	} else {
		_, found := seek(val, second)
		return found
	}
}

func (container *BinaryAddressContainer) GetIPsInRange(fromRange *net.IPNet) ([]*net.IP, error) {
	ones, _ := fromRange.Mask.Size()
	lowerFirst, lowerSecond, upperFirst, upperSecond := addressing.NetworkToUints(fromRange)
	var toReturn []*net.IP
	if ones == 0 {
		return container.GetAllIPs(), nil
	} else if ones == 128 {
		if container.ContainsIP(&fromRange.IP) {
			return []*net.IP{ &fromRange.IP }, nil
		} else {
			return []*net.IP{}, nil
		}
	} else if ones == 64 {
		if val, ok := container.addresses[lowerFirst]; !ok {
			return []*net.IP{}, nil
		} else {
			for _, curSecond := range val {
				toReturn = append(toReturn, addressing.UintsToAddress(lowerFirst, curSecond))
			}
			return toReturn, nil
		}
	} else if ones < 64 {
		parentRanges := seekRange(container.sortedHighKeys, lowerFirst, upperFirst)
		toReturn = []*net.IP{}
		for _, curRange := range parentRanges {
			for _, curLower := range container.addresses[curRange] {
				toReturn = append(toReturn, addressing.UintsToAddress(curRange, curLower))
			}
		}
		return toReturn, nil
	} else {
		if val, ok := container.addresses[lowerFirst]; !ok {
			return []*net.IP{}, nil
		} else {
			toReturn = []*net.IP{}
			seconds := seekRange(val, lowerSecond, upperSecond)
			for _, curSecond := range seconds {
				toReturn = append(toReturn, addressing.UintsToAddress(lowerFirst, curSecond))
			}
			return toReturn, nil
		}
	}
}

func (container *BinaryAddressContainer) CountIPsInRange(fromRange *net.IPNet) (uint32, error) { //TODO this is computationally expensive, is there a way we can pull it out?
	found, err := container.GetIPsInRange(fromRange)
	if err != nil {
		return 0, err
	} else {
		return uint32(len(found)), err
	}
}

func (container *BinaryAddressContainer) GetIPsInGenRange(fromRange *GenRange) []*net.IP {
	if len(fromRange.WildIndices) == 0 {
		checkIP := fromRange.GetIP()
		if container.ContainsIP(checkIP) {
			return []*net.IP{ checkIP }
		}
	}
	var toReturn []*net.IP
	rangeMask := fromRange.GetMask()
	firstCandidates := seekRange(container.sortedHighKeys, rangeMask.FirstMin, rangeMask.FirstMax)
	if len(firstCandidates) == 0 {
		return []*net.IP{}
	}
	firstCandidates = filterByMask(firstCandidates, rangeMask.FirstMask, rangeMask.FirstExpected)
	if len(firstCandidates) == 0 {
		return []*net.IP{}
	}
	for _, curFirst := range firstCandidates {
		secondCandidates := seekRange(container.addresses[curFirst], rangeMask.SecondMin, rangeMask.SecondMax)
		if len(secondCandidates) == 0 {
			continue
		}
		secondCandidates = filterByMask(secondCandidates, rangeMask.SecondMask, rangeMask.SecondExpected)
		if len(secondCandidates) == 0 {
			continue
		}
		for _, curSecond := range secondCandidates {
			toReturn = append(toReturn, addressing.UintsToAddress(curFirst, curSecond))
		}
	}
	return toReturn
}

func (container *BinaryAddressContainer) CountIPsInGenRange(fromRange *GenRange) int {
	return len(container.GetIPsInGenRange(fromRange))
}

func filterByMask(toFilter []uint64, mask uint64, expected uint64) []uint64 {
	var toReturn []uint64
	for _, curCheck := range toFilter {
		if (curCheck & mask) == expected {
			toReturn = append(toReturn, curCheck)
		}
	}
	return toReturn
}

func insert(into []uint64, toInsert uint64) ([]uint64, bool) {
	if len(into) == 0 {
		return []uint64{ toInsert }, true
	} else if len(into) == 1 {
		if into[0] < toInsert {
			return []uint64{ into[0], toInsert }, true
		} else {
			return []uint64{ toInsert, into[0] }, true
		}
	} else {
		index, found := seek(into, toInsert)
		if found {
			return into, false
		} else {
			into = append(into, uint64(0))
			copy(into[index+1:], into[index:])
			into[index] = toInsert
			return into, true
		}
	}
}

func seekRange(source []uint64, lowerBound uint64, upperBound uint64) []uint64 {
	if len(source) == 0 {
		return []uint64{}
	} else if lowerBound > source[len(source) - 1] || upperBound < source[0] {
		return []uint64{}
	} else if lowerBound == upperBound {
		_, found := seek(source, lowerBound)
		if found {
			return []uint64{ lowerBound }
		} else {
			return []uint64{}
		}
	} else if len(source) == 1 && lowerBound <= source[0] && upperBound >= source[0] {
		return source
	}
	lowerIndex, _ := seek(source, lowerBound)
	upperIndex, upperFound := seek(source, upperBound)
	if upperFound {
		upperIndex++
	}
	if upperIndex - lowerIndex == 1 {
		return []uint64{}
	} else {
		return source[lowerIndex:upperIndex]
	}
}

func seek(source []uint64, sought uint64) (int, bool) {
	if len(source) == 0 {
		return 0, false
	} else if sought < source[0] {
		return 0, false
	} else if sought == source[0] {
		return 0, true
	}
	curLower := 0
	curUpper := len(source)
	for {
		middle := curLower + (curUpper - curLower) / 2
		if source[middle] == sought {
			return middle, true
		} else if source[middle] < sought {
			curLower = middle
		} else {
			curUpper = middle
		}
		if curUpper - curLower == 1 {
			return curUpper, false
		}
	}
}