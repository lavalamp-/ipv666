package modeling

import (
	"github.com/lavalamp-/ipv666/internal"
	"github.com/lavalamp-/ipv666/internal/addressing"
	"net"
	"sync"
)

type GenClusterList struct {
	Clusters		[]*GenCluster
}

type GenCluster struct {
	Range			*GenRange
	SeedSet			AddressContainer  //TODO check to make sure this maintains state correctly
}

type GenRange struct {
	AddrNybbles		[]uint8
	WildIndices		map[int]internal.Empty
}

type clusterDistance struct {
	distance		int
	addr			*net.IP
}

func newClusterList() *GenClusterList {
	return &GenClusterList{
		Clusters: []*GenCluster{},
	}
}

func newGenRange(fromIP *net.IP) *GenRange {
	return &GenRange{
		AddrNybbles: addressing.GetNybblesFromIP(fromIP, 32),
		WildIndices: make(map[int]internal.Empty),
	}
}

func (genRange *GenRange) AddIP(toAdd *net.IP) {
	ipNybbles := addressing.GetNybblesFromIP(toAdd, 32)
	for i, curNybble := range ipNybbles {
		if genRange.AddrNybbles[i] != curNybble {
			genRange.WildIndices[i] = internal.Empty{}
		}
	}
}

func GetGenRangeFromIPs(fromIPs []*net.IP) *GenRange {
	newRange := newGenRange(fromIPs[0])
	for _, curIP := range fromIPs {
		newRange.AddIP(curIP)
	}
	return newRange
}

func newGenCluster(firstIP *net.IP, container AddressContainer) *GenCluster {
	container.AddIP(firstIP)
	return &GenCluster{
		Range: newGenRange(firstIP),
		SeedSet: container,
	}
}

func (cluster *GenCluster) distanceFromIP(toProcess *net.IP) int {
	ipNybbles := addressing.GetNybblesFromIP(toProcess, 32)
	toReturn := 0
	for i := range ipNybbles {
		if ipNybbles[i] != cluster.Range.AddrNybbles[i] {
			if _, ok := cluster.Range.WildIndices[i]; !ok {
				toReturn++
			}
		}
	}
	return toReturn
}

func (cluster *GenCluster) findCandidateSeeds(toProcess []*net.IP) (int, []*net.IP) {
	var responses = make(chan *clusterDistance)
	var wg sync.WaitGroup
	for _, curProcess := range toProcess {
		if !cluster.SeedSet.ContainsIP(curProcess) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				responses <- &clusterDistance{
					distance: cluster.distanceFromIP(curProcess),
					addr: curProcess,
				}
			}()
		}
	}
	wg.Wait()
	distance := 99999
	var toReturn []*net.IP
	for response := range responses {
		if response.distance < distance {
			distance = response.distance
			toReturn = []*net.IP{ response.addr }
		} else if response.distance == distance {
			toReturn = append(toReturn, response.addr)
		}
	}
	return distance, toReturn
}



func (clusterList *GenClusterList) doIt(toProcess []*net.IPNet) {

}
