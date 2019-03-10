package modeling

import (
	"fmt"
	"github.com/lavalamp-/ipv666/internal"
	"github.com/lavalamp-/ipv666/internal/addressing"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/persist"
	"github.com/spf13/viper"
	"math"
	"math/rand"
	"net"
	"sort"
	"strings"
)

type ClusterModel struct {
	ClusterSet			*ClusterSet					`msgpack:"c"`
	NybbleCounts		[]map[uint8]int				`msgpack:"n"`
	normalizedCounts	[][]uint8
}

type ClusterSet struct {
	Clusters			[]*GenCluster				`msgpack:"c"`
	Captured			int							`msgpack:"a"`
	RangeSize			int							`msgpack:"s"`
	Density				float64						`msgpack:"d"`
}

type GenCluster struct {
	Range				*GenRange					`msgpack:"r"`
	Captured			int							`msgpack:"c"`
	Density				float64						`msgpack:"d"`
	Size				int							`msgpack:"s"`
}

type GenRange struct {
	AddrNybbles			[]uint8						`msgpack:"n"`
	WildIndices			map[int]internal.Empty		`msgpack:"w"`
}

type GenRangeMask struct {
	FirstMask      uint64
	FirstExpected  uint64
	FirstMin       uint64
	FirstMax       uint64
	SecondMask     uint64
	SecondExpected uint64
	SecondMin      uint64
	SecondMax      uint64
}

type clusterList []*GenCluster

type addrProcessFunc func(*net.IP) (bool, error)

// Model

func (clusterModel *ClusterModel) GenerateAddresses(generateCount int, jitter float64) []*net.IP {
	addrTree := newAddressTree()
	var toReturn []*net.IP
	iteration := 0
	for {
		if iteration % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Infof("Generating new candidate address %d using clustering model. Unique count size is %d.", iteration, addrTree.Size())
		}
		newAddr := clusterModel.GenerateAddress(jitter)
		if addrTree.AddIP(newAddr) {
			toReturn = append(toReturn, newAddr)
			if len(toReturn) >= generateCount {
				break
			}
		}
		iteration++
	}
	logging.Infof("Successfully generated %d addresses in %d iterations.", len(toReturn), iteration)
	return toReturn
}

func (clusterModel *ClusterModel) GenerateAddressesFromNetwork(generateCount int, jitter float64, network *net.IPNet) ([]*net.IP, error) {
	ones, _ := network.Mask.Size()
	if ones % 4 != 0 {
		return nil, fmt.Errorf("generating addresses in a network requires a network length that is divisible by 4 (got length of %d)", ones)
	}
	networkNybbles := addressing.GetNybblesFromNetwork(network)
	var toReturn []*net.IP
	iteration := 0
	addrTree := newAddressTree()
	for {
		if iteration % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Infof("Generating new candidate address %d using clustering model. Unique count size is %d.", iteration, addrTree.Size())
		}
		newAddr := clusterModel.generateAddressFromNybbles(jitter, networkNybbles)
		if addrTree.AddIP(newAddr) {
			toReturn = append(toReturn, newAddr)
			if len(toReturn) >= generateCount {
				break
			}
		}
		iteration++
	}
	logging.Infof("Successfully generated %d addresses in %d iterations.", len(toReturn), iteration)
	return toReturn, nil
}

func (clusterModel *ClusterModel) GenerateAddressesFromNetworkWithCallback(generateCount int, jitter float64, network *net.IPNet, fn addrProcessFunc) ([]*net.IP, error) {
	ones, _ := network.Mask.Size()
	if ones % 4 != 0 {
		return nil, fmt.Errorf("generating addresses in a network requires a network length that is divisible by 4 (got length of %d)", ones)
	}
	networkNybbles := addressing.GetNybblesFromNetwork(network)
	var toReturn []*net.IP
	iteration := 0
	for len(toReturn) < generateCount {
		newIP := clusterModel.generateAddressFromNybbles(jitter, networkNybbles)
		isFiltered, err := fn(newIP)
		if err != nil {
			return nil, err
		} else if !isFiltered {
			toReturn = append(toReturn, newIP)
		}
		iteration++
	}
	return toReturn, nil
}

func (clusterModel *ClusterModel) GenerateAddress(jitter float64) *net.IP {
	if len(clusterModel.normalizedCounts) == 0 {
		clusterModel.generateNormalizedCounts()
	}
	index := rand.Int63n(int64(len(clusterModel.ClusterSet.Clusters)))
	cluster := clusterModel.ClusterSet.Clusters[index]
	var nybbles []uint8
	for i := range cluster.Range.AddrNybbles {
		if _, ok := cluster.Range.WildIndices[i]; ok {
			nybbles = append(nybbles, uint8(rand.Int31n(16)))
		} else if float64(rand.Int31n(10000)) / 100.0 <= jitter * 100 {
			index := rand.Int63n(int64(len(clusterModel.normalizedCounts[i])))
			nybbles = append(nybbles, clusterModel.normalizedCounts[i][index])
		} else {
			nybbles = append(nybbles, cluster.Range.AddrNybbles[i])
		}
	}
	return addressing.NybblesToIP(nybbles)
}

func (clusterModel *ClusterModel) generateAddressFromNybbles(jitter float64, fromNybbles []uint8) *net.IP {
	if len(clusterModel.normalizedCounts) == 0 {
		clusterModel.generateNormalizedCounts()
	}
	index := rand.Int63n(int64(len(clusterModel.ClusterSet.Clusters)))
	cluster := clusterModel.ClusterSet.Clusters[index]
	var nybbles []uint8
	for i := len(fromNybbles); i < 32; i++ {
		if _, ok := cluster.Range.WildIndices[i]; ok {
			nybbles = append(nybbles, uint8(rand.Int31n(16)))
		} else if float64(rand.Int31n(10000)) / 100.0 <= jitter * 100 {
			index := rand.Int63n(int64(len(clusterModel.normalizedCounts[i])))
			nybbles = append(nybbles, clusterModel.normalizedCounts[i][index])
		} else {
			nybbles = append(nybbles, cluster.Range.AddrNybbles[i])
		}
	}
	return addressing.NybblesToIP(append(fromNybbles, nybbles...))
}

func (clusterModel *ClusterModel) generateNormalizedCounts() {
	var normalizedCounts [][]uint8
	for _, curCounts := range clusterModel.NybbleCounts {
		if len(curCounts) == 0 {
			normalizedCounts = append(normalizedCounts, getCountsEvenDist())
		} else if len(curCounts) < 16 {
			normalizedCounts = append(normalizedCounts, getCountsWithMinDist(curCounts, viper.GetFloat64("ModelMinNybblePercent")))
		} else {
			normalizedCounts = append(normalizedCounts, normalizeCountsToMinPercent(curCounts, viper.GetFloat64("ModelMinNybblePercent")))
		}
	}
	clusterModel.normalizedCounts = normalizedCounts
}

func getCountsEvenDist() []uint8 {
	var addPercents = make(map[uint8]float64)
	var i uint8
	for i = 0; i < 16; i++ {
		addPercents[i] = 1.0 / 16.0
	}
	return percentsToNybbleCounts(addPercents, viper.GetInt("ModelDistributionSize"))
}

func getCountsWithMinDist(counts map[uint8]int, minPercent float64) []uint8 {
	existingCount := 0.0
	for _, v := range counts {
		existingCount += float64(v)
	}
	var i uint8
	populateIndices := make(map[uint8]*internal.Empty)
	calculateCount := 0
	for i = 0; i < 16; i++ {
		if val, ok := counts[i]; ok {
			if float64(val) / existingCount < minPercent {
				populateIndices[i] = &internal.Empty{}
			} else {
				calculateCount += val
			}
		} else {
			populateIndices[i] = &internal.Empty{}
		}
	}
	totalSize := float64(calculateCount) / (1.0 - (minPercent * float64(len(populateIndices))))
	var addPercents = make(map[uint8]float64)
	for i = 0; i < 16; i++ {
		if _, ok := populateIndices[i]; ok {
			addPercents[i] = minPercent
		} else {
			val, _ := counts[i]
			addPercents[i] = float64(val) / totalSize
		}
	}
	return percentsToNybbleCounts(addPercents, viper.GetInt("ModelDistributionSize"))
}

func normalizeCountsToMinPercent(counts map[uint8]int, minPercent float64) []uint8 {
	var minFound = 999999999999999.0
	var maxFound = -1.0
	var totalFound = 0
	for _, v := range counts {
		minFound = math.Min(minFound, float64(v))
		maxFound = math.Max(maxFound, float64(v))
		totalFound += v
	}
	var addPercents = make(map[uint8]float64)
	if minFound / float64(totalFound) < minPercent {
		middle := float64(totalFound) / 16.0
		toGrow := (minPercent * float64(totalFound)) - minFound
		middleDistance := middle - minFound
		middleShrinkPercent := (middleDistance - toGrow) / middleDistance
		for k, v := range counts {
			newVal := middle - ((middle - float64(v)) * middleShrinkPercent)
			addPercents[k] = newVal / float64(totalFound)
		}
	} else {
		for k, v := range counts {
			addPercents[k] = float64(v) / float64(totalFound)
		}
	}
	return percentsToNybbleCounts(addPercents, viper.GetInt("ModelDistributionSize"))
}

func percentsToNybbleCounts(fromPercents map[uint8]float64, distSize int) []uint8 {
	var toReturn []uint8
	for k, v := range fromPercents {
		for i := 0; i < int(math.Ceil(v * float64(distSize))); i++ {
			toReturn = append(toReturn, k)
		}
	}
	return toReturn
}

func CreateClusteringModel(fromAddrs []*net.IP) *ClusterModel {

	// Convert all IPs to clusters and create corpus

	logging.Infof("Preparing initial data structures from %d addresses.", len(fromAddrs))

	clusters := newGenClusters(fromAddrs)
	corpus := CreateFromAddresses(fromAddrs, viper.GetInt("LogLoopEmitFreq"))

	logging.Infof("Done creating data structures from %d addresses.", len(fromAddrs))

	// Calculate upgrade potential and create stardust from poor performers

	logging.Infof("Now reviewing %d initial cluster candidates for poor performers.", len(fromAddrs))
	var dustAddrs []*net.IP
	var clusterMap = make(map[string]*internal.Empty)
	var modelCandidates clusterList
	empty := &internal.Empty{}

	for i, cluster := range clusters {
		if i % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Infof("Processing cluster %d out of %d for poor performers.", i, len(clusters))
		}
		upgradeDensity, upgradeCount, upgradeIndices := cluster.getBestUpgradeOptions(corpus)
		if len(upgradeIndices) == 32 {  // If all upgrades are equivalent then all upgrades are bad
			dustAddrs = append(dustAddrs, cluster.Range.GetIP())
		} else {
			for _, curIndex := range upgradeIndices {
				newRange := cluster.Range.CopyWithIndices([]int{ curIndex })
				newCluster := GenCluster{
					Range: 		newRange,
					Captured:	upgradeCount,
					Density:	upgradeDensity,
					Size:		int(newRange.Size()), // TODO undo silly typecasting
				}
				sig := newCluster.signature()
				if _, ok := clusterMap[sig]; ok {
					continue
				} else {
					modelCandidates = append(modelCandidates, &newCluster)
					clusterMap[sig] = empty
				}
			}
		}
	}

	initialModelCandidateSize := len(modelCandidates)
	logging.Infof("Reviewed %d initial cluster candidates. %d are now stardust, %d are model candidates.", len(fromAddrs), len(dustAddrs), initialModelCandidateSize)

	// Generate upgrade candidates

	logging.Infof("Processing %d model candidates into upgrade candidates.", len(modelCandidates))
	skipped := 0
	var upgradeMap = make(map[string]*internal.Empty)
	var upgradeCandidates clusterList

	for i, cluster := range modelCandidates {
		if i % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Infof("Processing model candidate %d out of %d for upgrade candidates.", i, len(modelCandidates))
		}
		upgradeDensity, upgradeCount, upgradeIndices := cluster.getBestUpgradeOptions(corpus)
		if len(upgradeIndices) == 31 {  // Thee case where all upgrades are the same is not an upgrade
			skipped++
		} else {
			for _, curIndex := range upgradeIndices {
				newRange := cluster.Range.CopyWithIndices([]int{ curIndex })
				newCluster := GenCluster{
					Range: 		newRange,
					Captured:	upgradeCount,
					Density:	upgradeDensity,
					Size:		int(newRange.Size()), // TODO undo silly typecasting
				}
				sig := newCluster.signature()
				if _, ok := upgradeMap[sig]; !ok {
					upgradeCandidates = append(upgradeCandidates, &newCluster)
					upgradeMap[sig] = empty
				} else {
					skipped++
				}
			}
		}
	}

	logging.Infof("Processed %d model candidates into %d upgrade candidates (skipped %d). Sorting results now.", len(modelCandidates), len(upgradeCandidates), skipped)
	sort.Slice(upgradeCandidates, func(i, j int) bool {
		return upgradeCandidates[i].Density > upgradeCandidates[j].Density
	})
	logging.Infof("Upgrade candidates sorted by density.")

	// Enter into processing loop

	iteration := 0
	var lastScore float64
	var lastClusterSet *ClusterSet

	for {

		// Take the next upgrade candidate in line

		candidate := upgradeCandidates[0]
		upgradeCandidates = upgradeCandidates[1:]
		sig := candidate.signature()

		// Test to see if this cluster has already been added

		if _, ok := clusterMap[sig]; ok {
			continue
		}

		// Calculate its upgrade candidates and insert them into the upgrade candidate list

		upgradeDensity, upgradeCount, upgradeIndices := candidate.getBestUpgradeOptions(corpus)
		if len(upgradeIndices) + len(candidate.Range.WildIndices) != 32 {  // So long as upgrade is not worst case scenario, we add them to upgrade candidates
			for _, curIndex := range upgradeIndices {
				newRange := candidate.Range.CopyWithIndices([]int{ curIndex })
				newCluster := &GenCluster{
					Range: 		newRange,
					Captured:	upgradeCount,
					Density:	upgradeDensity,
					Size:		int(newRange.Size()), // TODO undo silly typecasting
				}
				sig := newCluster.signature()
				if _, ok := upgradeMap[sig]; !ok {
					upgradeCandidates = upgradeCandidates.insertByDensity(newCluster)
					upgradeMap[sig] = empty
				}
			}
		}

		// Add the upgrade candidate to the list of model candidates

		modelCandidates = append(modelCandidates, candidate)

		// Add to the iteration counter and check to see if checkpoint should be evaluated

		iteration++
		if iteration % viper.GetInt("ModelCheckCount") == 0 {

			logging.Infof("Evaluating %d model candidates for cluster model.", len(modelCandidates))

			// Remove all redundant candidates

			logging.Infof("Removing redundant clusters from model candidates...")
			modelCandidates = modelCandidates.removeRedundant()
			logging.Infof("Resulting model candidate list is down to length of %d.", len(modelCandidates))

			// Create new cluster set from candidates

			toCheck := newClusterSetFromClusters(modelCandidates)
			toCheck.ResetCounts(corpus)

			// Calculate cluster set score

			capturedScore := float64(toCheck.Captured) / float64(len(fromAddrs))
			densityScore := toCheck.Density
			clusterCountScore := float64(initialModelCandidateSize - len(toCheck.Clusters)) / float64(initialModelCandidateSize)
			cumulativeScore := (0.5 * densityScore) + (0.333333333 * capturedScore) + (0.16666666667 * clusterCountScore)

			logging.Infof("Cluster set has %f density, %d captured, and %d size with %d clusters. Cluster set score is %f.", toCheck.Density, toCheck.Captured, toCheck.RangeSize, len(toCheck.Clusters), cumulativeScore)

			if lastClusterSet == nil {
				logging.Infof("No existing cluster set to compare against yet...")
				lastScore = cumulativeScore
				lastClusterSet = toCheck
			} else if cumulativeScore >= lastScore {
				logging.Infof("New cluster set has better score than last cluster set (%f vs %f). Continuing analysis.", cumulativeScore, lastScore)
				lastScore = cumulativeScore
				lastClusterSet = toCheck
			} else {
				logging.Infof("New cluster set's score of %f is worse than last score (%f). No more improvements to be made.", cumulativeScore, lastScore)
				break
			}
		}
	}

	return &ClusterModel{
		ClusterSet:	lastClusterSet,
		NybbleCounts: addrsToNybbleCounts(dustAddrs),
	}
}

func addrsToNybbleCounts(toProcess []*net.IP) []map[uint8]int {
	var toReturn []map[uint8]int
	for i := 0; i < 32; i++ {
		toReturn = append(toReturn, make(map[uint8]int))
	}
	for _, curAddr := range toProcess {
		for i, curNybble := range addressing.GetNybblesFromIP(curAddr, 32) {
			if _, ok := toReturn[i][curNybble]; !ok {
				toReturn[i][curNybble] = 0
			}
			toReturn[i][curNybble]++
		}
	}
	return toReturn
}

func (clusterModel *ClusterModel) Save(filePath string) error {
	return persist.Save(filePath, clusterModel)
}

func LoadModelFromFile(filePath string) (*ClusterModel, error) {
	var toReturn ClusterModel
	err := persist.Load(filePath, &toReturn)
	return &toReturn, err
}

func LoadModelFromBytes(fromBytes []byte) (*ClusterModel, error) {
	var toReturn ClusterModel
	err := persist.Unmarshal(fromBytes, &toReturn)
	return &toReturn, err
}

// ClusterSet

func (clusterSet *ClusterSet) GenerateAddresses(generateCount int, jitter float64) []*net.IP {
	toReturn := newAddressTree()
	iteration := 0
	for {
		if iteration % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Infof("Generating new candidate address %d using clustering model. Unique count size is %d.", iteration, toReturn.Size())
		}
		cluster := clusterSet.Clusters[rand.Int63n(int64(len(clusterSet.Clusters)))]
		newAddr := cluster.generateAddr(jitter)
		toReturn.AddIP(newAddr)
		iteration++
		if toReturn.Size() >= generateCount {
			break
		}
	}
	returnIPs := toReturn.GetAllIPs()
	logging.Infof("Successfully generated %d addresses in %d iterations.", len(returnIPs), iteration)
	return returnIPs
}

func (clusterSet *ClusterSet) Save(filePath string) error {
	return persist.Save(filePath, clusterSet)
}

func LoadClusterSetFromFile(filePath string) (*ClusterSet, error) {
	var toReturn ClusterSet
	err := persist.Load(filePath, &toReturn)
	return &toReturn, err
}

func (clusterSet *ClusterSet) ResetCounts(corpus AddressContainer) {
	total := newAddressTree()
	rangeSize := 0
	for _, cluster := range clusterSet.Clusters {
		covered := corpus.GetIPsInGenRange(cluster.Range)
		total.AddIPs(covered, viper.GetInt("LogLoopEmitFreq"))
		rangeSize += int(cluster.Range.Size()) // TODO wtf is with this casting
	}
	clusterSet.Captured = total.Size()
	clusterSet.RangeSize = rangeSize
	clusterSet.Density = float64(clusterSet.Captured) / float64(clusterSet.RangeSize)
}

func (clusterSet *ClusterSet) GetUpgrade(corpus AddressContainer, densityThreshold float64) *ClusterSet {
	var clusterMap = make(map[string]*GenCluster)
	changed := false
	for i, cluster := range clusterSet.Clusters {
		if i % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Infof("Processing cluster %d out of %d for upgrade candidates.", i, len(clusterSet.Clusters))
		}
		upgradeDensity, upgradeCount, upgradeIndices := cluster.getBestUpgradeOptions(corpus)
		if len(upgradeIndices) == 0 && cluster.Density != 1.0 {
			clusterMap[cluster.signature()] = cluster
		} else if len(upgradeIndices) + len(cluster.Range.WildIndices) == 32 {
			clusterMap[cluster.signature()] = cluster // Bottomed out, every next step is the worst case scenario
		} else {
			changed = true
			for _, curIndex := range upgradeIndices {
				newRange := cluster.Range.CopyWithIndices([]int{ curIndex })
				newCluster := &GenCluster{
					Range: 		newRange,
					Captured:	upgradeCount,
					Density:	upgradeDensity,
					Size:		int(newRange.Size()), // TODO undo silly typecasting
				}
				clusterMap[newCluster.signature()] = newCluster
			}
		}
	}
	var candidateClusters []*GenCluster
	for _, v := range clusterMap {
		candidateClusters = append(candidateClusters, v)
	}
	// TODO remove clusters that are captured by other clusters
	logging.Infof("A total of %d candidate clusters were identified.", len(candidateClusters))
	sort.Slice(candidateClusters, func(i, j int) bool {
		return candidateClusters[i].Density > candidateClusters[j].Density
	})
	logging.Infof("Checking to see if new cluster set covers all candidates...")
	var nClusterSet = newClusterSetFromClusters(candidateClusters)
	nClusterSet.ResetCounts(corpus)
	if nClusterSet.Captured < corpus.Size() {
		logging.Infof("We haven't yet covered the corpus (%d out of %d). Including all growth candidates.", nClusterSet.Captured, corpus.Size())
	} else {
		var newClusters []*GenCluster
		curCaptured, curSize, curDensity, capacityMet := 0, 0, 0.0, false
		capacityMetAt, additionalCount := 0, 0
		for i := 0; i < len(candidateClusters); i++ {
			if i % viper.GetInt("LogLoopEmitFreq") == 0 {
				logging.Infof("Processing candidate cluster %d out of %d for new cluster set.", i, len(candidateClusters))
			}
			if !capacityMet {
				newClusters = append(newClusters, candidateClusters[i])
				curCaptured += candidateClusters[i].Captured
				curSize += candidateClusters[i].Size
				if curCaptured >= clusterSet.Captured {
					curDensity = float64(curCaptured) / float64(curSize)
					capacityMet = true
					capacityMetAt = i
					logging.Infof("Met previous coverage capacity of %d with only %d clusters (down from %d).", clusterSet.Captured, i, len(clusterSet.Clusters))
				}
			} else {
				newCaptured := curCaptured + candidateClusters[i].Captured
				newSize := curSize + candidateClusters[i].Size
				newDensity := float64(newCaptured) / float64(newSize)
				if newDensity >= densityThreshold {
					curDensity = newDensity
					curSize = newSize
					curCaptured = newCaptured
					newClusters = append(newClusters, candidateClusters[i])
					additionalCount++
				} else {
					break
				}
			}
		}
		if !capacityMet {
			logging.Infof("Previous capacity was not met despite adding all upgrade candidates. No upgrade possible.")
			return nil
		}
		logging.Infof("The new candidate cluster set is %d long (capacity met at %d, additionals ended at %d). Captured is %d and density is %f.", len(newClusters), capacityMetAt, additionalCount, curCaptured, curDensity)
		nClusterSet = newClusterSetFromClusters(newClusters)
		nClusterSet.ResetCounts(corpus)
		logging.Infof("After re-processing, density went from %f -> %f.", curDensity, nClusterSet.Density)
	}
	// TODO calculate actual density of new cluster set and see how far off we are
	if clusterSet.Density == 1.0 {
		logging.Infof("Existing cluster was perfect match (ergo, first iteration). New set of capacity %d and density %f is better.", nClusterSet.Captured, nClusterSet.Density)
		return nClusterSet
	} else if !changed {
		logging.Infof("No change was detected since the previous iteration. Further upgrade not possible.")
		return nil
	} else if nClusterSet.Captured < corpus.Size() && nClusterSet.Density >= densityThreshold { //TODO put this in a variable
		logging.Infof("The cluster does not yet cover the whole corpus (%d out of %d).", nClusterSet.Captured, corpus.Size())
		return nClusterSet
	} else if nClusterSet.Density > clusterSet.Density {
		logging.Infof("New density of %f beats old density of %f.", nClusterSet.Density, clusterSet.Density)
		return nClusterSet
	} else if nClusterSet.Density == clusterSet.Density && nClusterSet.Captured > clusterSet.Captured {
		logging.Infof("New density is the same but has more captured (%d vs %d).", nClusterSet.Captured, clusterSet.Captured)
		return nClusterSet
	} else {
		logging.Infof("The new cluster set is not an upgrade. Old captured and density: %d, %f, new captured and density: %d, %f.", clusterSet.Captured, clusterSet.Density, nClusterSet.Captured, nClusterSet.Density)
		return nil
	}
}

func FindGoodSeeds(find int, pick int, from int, corpus AddressContainer) []*net.IP {
	var toReturn []*net.IP
	addrs := corpus.GetAllIPs()
	size := len(addrs)
	picked := EmptyContainer()
	iteration := 0
	logging.Infof("Finding a total of %d good seeds from a corpus of size %d. Picking %d from %d.", find, size, pick, from)
	for {
		logging.Infof("On iteration %d.", iteration)
		var candidates []*net.IP
		for {
			nextAddr := addrs[rand.Int63n(int64(len(addrs)))]
			if !picked.ContainsIP(nextAddr) {
				picked.AddIP(nextAddr)
				candidates = append(candidates, nextAddr)
			}
			if len(candidates) >= from {
				break
			}
		}
		logging.Infof("Picked a total of %d candidates at random out of %d.", len(candidates), size)
		var results []*GenCluster
		for _, curCandidate := range candidates {
			cluster := newGenCluster(curCandidate)
			density, count, indices := cluster.getBestUpgradeOptions(corpus)
			if len(indices) == 0 && cluster.Density != 1.0 {
				results = append(results, cluster)
			} else {
				for _, curIndex := range indices {
					newRange := cluster.Range.CopyWithIndices([]int{ curIndex })
					newCluster := &GenCluster{
						Range: 		newRange,
						Captured:	count,
						Density:	density,
						Size:		int(newRange.Size()), // TODO undo silly typecasting
					}
					results = append(results, newCluster)
				}
			}
		}
		logging.Infof("Done evaluating %d candidates. Now sorting and picking the top %d addresses.", len(candidates), pick)
		sort.Slice(results, func(i, j int) bool {
			return results[i].Density > results[j].Density
		})
		resultContainer := EmptyContainer()
		for _, cluster := range results {
			resultContainer.AddIP(cluster.Range.GetIP())
			if resultContainer.Size() >= pick {
				break
			}
		}
		logging.Infof("Successfully picked %d high-performing candidates.", resultContainer.Size())
		toReturn = append(toReturn, resultContainer.GetAllIPs()...)
		if len(toReturn) >= find {
			break
		} else {
			iteration++
		}
	}
	toReturn = toReturn[:find]
	logging.Infof("It took a total of %d iterations to find %d seed candidates.", iteration, len(toReturn))
	return toReturn
}

func GetBestClusterSetFromIPs(toParse []*net.IP, modelSize int, pickCount int, pickSize int, threshold float64) *ClusterSet {  //TODO add ability to specify container type at command line
	corpus := CreateFromAddresses(toParse, viper.GetInt("LogLoopEmitFreq")) //TODO get rid of this sillymess
	logging.Infof("Generating a best cluster set based on %d addresses.", corpus.ChildrenCount)
	bestSeeds := FindGoodSeeds(modelSize, pickCount, pickSize, corpus)
	logging.Infof("Instantiating a new cluster with the best candidate seeds.")
	curSet := newClusterSetFromIPs(bestSeeds)
	logging.Infof("Initial coverage is %d and density is %f.", curSet.Captured, curSet.Density)
	for {
		nextSet := curSet.GetUpgrade(corpus, threshold)
		if nextSet == nil {
			logging.Infof("Upgrade was not possible. It looks like we've found our match.")
			break
		} else {
			curSet = nextSet
			logging.Infof("Upgrading to new cluster (size %d) with coverage of %d and density of %f.", len(curSet.Clusters), curSet.Captured, curSet.Density)
		}
	}
	return curSet
}

func newClusterSetFromIPs(addrs []*net.IP) *ClusterSet {
	toReturn := &ClusterSet{
		Clusters: 	[]*GenCluster{},
		Captured:	0,
		RangeSize:	0,
		Density:	-1.0,
	}
	for i, addr := range addrs {
		if i % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Infof("Processing address %d out of %d.", i, len(addrs))
		}
		toReturn.AddCluster(newGenCluster(addr), false)
	}
	toReturn.Captured = toReturn.getCumulativeCaptured()
	toReturn.RangeSize = toReturn.getCumulativeRangeSize()
	toReturn.Density = float64(toReturn.Captured) / float64(toReturn.RangeSize)
	return toReturn
}

func newClusterSetFromClusters(clusters []*GenCluster) *ClusterSet {
	toReturn := ClusterSet{
		Clusters: 	[]*GenCluster{},
		Captured:	0,
		RangeSize:	0,
		Density:	-1.0,
	}
	toReturn.AddClusters(clusters)
	return &toReturn
}

func (clusterSet *ClusterSet) getCumulativeCaptured() int {
	var toReturn = 0
	for _, cluster := range clusterSet.Clusters {
		toReturn += cluster.Captured
	}
	return toReturn
}

func (clusterSet *ClusterSet) getCumulativeRangeSize() int {
	var toReturn = 0
	for _, cluster := range clusterSet.Clusters {
		toReturn += int(cluster.Size)
	}
	return toReturn
}

func (clusterSet *ClusterSet) AddCluster(toAdd *GenCluster, withUpdate bool) {
	clusterSet.Clusters = append(clusterSet.Clusters, toAdd)
	if withUpdate {
		clusterSet.Captured = clusterSet.getCumulativeCaptured()
		clusterSet.RangeSize = clusterSet.getCumulativeRangeSize()
		clusterSet.Density = float64(clusterSet.Captured) / float64(clusterSet.RangeSize)
	}
}

func (clusterSet *ClusterSet) AddClusters(toAdd []*GenCluster) {
	for i, curAdd := range toAdd {
		if i % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Infof("Processing cluster %d out of %d.", i, len(toAdd))
		}
		clusterSet.AddCluster(curAdd, false)
	}
	clusterSet.Captured = clusterSet.getCumulativeCaptured()
	clusterSet.RangeSize = clusterSet.getCumulativeRangeSize()
	clusterSet.Density = float64(clusterSet.Captured) / float64(clusterSet.RangeSize)
}

// Cluster

func newGenClusters(addrs []*net.IP) []*GenCluster {
	var toReturn []*GenCluster
	for _, addr := range addrs {
		toReturn = append(toReturn, newGenCluster(addr))
	}
	return toReturn
}

func newGenCluster(firstIP *net.IP) *GenCluster {
	return &GenCluster{
		Range:			newGenRange(firstIP),
		Captured:		1,
		Density:		1.0,
		Size:			1,
	}
}

func (list clusterList) removeRedundant() clusterList {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Size > list[j].Size
	})
	rangeTree := NewRangeTree()
	var toReturn clusterList
	for _, curCluster := range list {
		if rangeTree.AddRange(curCluster.Range) {
			toReturn = append(toReturn, curCluster)
		}
	}
	return toReturn
}

func (list clusterList) insertByDensity(toInsert *GenCluster) clusterList {
	if len(list) == 0 {
		return []*GenCluster{ toInsert }
	} else if len(list) == 1 {
		if list[0].Density > toInsert.Density {
			return []*GenCluster{ list[0], toInsert }
		} else {
			return []*GenCluster{ toInsert, list[0] }
		}
	} else {
		index, _ := list.seekDensity(toInsert.Density)
		list = append(list, &GenCluster{})
		copy(list[index+1:], list[index:])
		list[index] = toInsert
		return list
	}
}

func (list clusterList) seekDensity(density float64) (int, bool) {
	if len(list) == 0 {
		return 0, false
	} else if density > list[0].Density {
		return 0, false
	} else if density == list[0].Density {
		return 0, true
	} else if density < list[len(list) - 1].Density {
		return len(list), false
	} else if density == list[len(list) - 1].Density {
		return list.findFirstOfDensity(len(list) - 1, density), true
	}
	curLower := 0
	curUpper := len(list)
	for {
		middle := curLower + (curUpper - curLower) / 2
		if list[middle].Density == density {
			return list.findFirstOfDensity(middle, density), true
		} else if list[middle].Density > density {
			curLower = middle
		} else {
			curUpper = middle
		}
		if curUpper - curLower == 1 {
			return curUpper, false
		}
	}
}

func (list clusterList) findFirstOfDensity(searchStart int, density float64) int {
	for i := searchStart; i >= 0; i-- {
		if list[i].Density > density {
			return i + 1
		}
	}
	return 0
}

func (cluster *GenCluster) generateAddr(jitter float64) *net.IP {
	var addrNybbles []uint8
	for i := range cluster.Range.AddrNybbles {
		if _, ok := cluster.Range.WildIndices[i]; ok {
			addrNybbles = append(addrNybbles, uint8(rand.Int31n(16)))
		} else if float64(rand.Int31n(10000)) / 100.0 <= jitter * 100 {
			addrNybbles = append(addrNybbles, uint8(rand.Int31n(16)))
		} else {
			addrNybbles = append(addrNybbles, cluster.Range.AddrNybbles[i])
		}
	}
	return addressing.NybblesToIP(addrNybbles)
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

func (cluster *GenCluster) signature() string {
	var toReturn []string
	for i := range cluster.Range.AddrNybbles {
		if _, ok := cluster.Range.WildIndices[i]; ok {
			toReturn = append(toReturn, "?")
		} else {
			toReturn = append(toReturn, fmt.Sprintf("%02x", cluster.Range.AddrNybbles[i])[1:])

		}
	}
	return strings.Join(toReturn, "")
}

func (cluster *GenCluster) getBestUpgradeOptions(corpus AddressContainer) (float64, int, []int) {
	var toReturn []int
	var density = -1.0
	var count = -1
	for i := 0; i < 32; i++ {
		if _, ok := cluster.Range.WildIndices[i]; !ok {
			newRange := cluster.Range.CopyWithIndices([]int{ i })
			capturedCount := corpus.CountIPsInGenRange(newRange)
			capturedDensity := float64(capturedCount) / newRange.Size()
			if capturedDensity > density {
				density = capturedDensity
				count = capturedCount
				toReturn = []int { i }
			} else if capturedDensity == density {
				toReturn = append(toReturn, i)
			}
		}
	}
	return density, count, toReturn
}

// Range

func (genRange *GenRange) GetTreeNybbles() []uint16 {
	var toReturn []uint16
	for i := range genRange.AddrNybbles {
		if _, ok := genRange.WildIndices[i]; ok {
			toReturn = append(toReturn, 16)
		} else {
			toReturn = append(toReturn, uint16(genRange.AddrNybbles[i]))
		}
	}
	return toReturn
}

func (genRange *GenRange) Contains(toCheck *GenRange) bool {
	for i := range genRange.AddrNybbles {
		if _, ok := genRange.WildIndices[i]; !ok {
			if genRange.AddrNybbles[i] != toCheck.AddrNybbles[i] {
				return false
			}
		}
	}
	return true
}

func (genRange *GenRange) GetIP() *net.IP {
	return addressing.NybblesToIP(genRange.AddrNybbles)
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

func (genRange *GenRange) Equals(otherRange *GenRange) bool {
	for i := range genRange.AddrNybbles {
		if _, ok :=  genRange.WildIndices[i]; !ok {
			if genRange.AddrNybbles[i] != otherRange.AddrNybbles[i] {
				return false
			}
		}
	}
	return true
}

func (genRange *GenRange) AddIPs(toAdd []*net.IP) {
	for _, curAdd := range toAdd {
		genRange.AddIP(curAdd)
	}
}

func (genRange *GenRange) Copy() *GenRange {
	toReturn := GenRange{
		AddrNybbles: []uint8{},
		WildIndices: make(map[int]internal.Empty),
	}
	for _, nybble := range genRange.AddrNybbles {
		toReturn.AddrNybbles = append(toReturn.AddrNybbles, nybble)
	}
	for k, v := range genRange.WildIndices {
		toReturn.WildIndices[k] = v
	}
	return &toReturn
}

func (genRange *GenRange) Size() float64 {
	return math.Pow(16, float64(len(genRange.WildIndices)))
}

func (genRange *GenRange) GetMask() *GenRangeMask {
	firstMask := uint64(0)
	firstExpected := uint64(0)
	firstMin := uint64(0)
	firstMax := uint64(0)
	secondMask := uint64(0)
	secondExpected := uint64(0)
	secondMin := uint64(0)
	secondMax := uint64(0)
	for i := range genRange.AddrNybbles {
		if i < 16 {
			shiftAmount := uint((15 - i) * 4)
			if _, ok := genRange.WildIndices[i]; ok {
				firstExpected ^= uint64(0) << shiftAmount
				firstMask ^= uint64(0) << shiftAmount
				firstMin ^= uint64(0) << shiftAmount
				firstMax ^= uint64(0x0f) << shiftAmount
			} else {
				firstExpected ^= uint64(genRange.AddrNybbles[i]) << shiftAmount
				firstMask ^= uint64(0x0f) << shiftAmount
				firstMin ^= uint64(genRange.AddrNybbles[i]) << shiftAmount
				firstMax ^= uint64(genRange.AddrNybbles[i]) << shiftAmount
			}
		} else {
			shiftAmount := uint((31 - i) * 4)
			if _, ok := genRange.WildIndices[i]; ok {
				secondExpected ^= uint64(0) << shiftAmount
				secondMask ^= uint64(0) << shiftAmount
				secondMin ^= uint64(0) << shiftAmount
				secondMax ^= uint64(0x0f) << shiftAmount
			} else {
				secondExpected ^= uint64(genRange.AddrNybbles[i]) << shiftAmount
				secondMask ^= uint64(0x0f) << shiftAmount
				secondMin ^= uint64(genRange.AddrNybbles[i]) << shiftAmount
				secondMax ^= uint64(genRange.AddrNybbles[i]) << shiftAmount
			}
		}
	}
	return &GenRangeMask{
		FirstMask:      firstMask,
		FirstExpected:  firstExpected,
		FirstMin:       firstMin,
		FirstMax:       firstMax,
		SecondMask:     secondMask,
		SecondExpected: secondExpected,
		SecondMin:      secondMin,
		SecondMax:      secondMax,
	}
}

func (genRange *GenRange) CopyWithIPs(newIPs []*net.IP) *GenRange {
	toReturn := genRange.Copy()
	toReturn.AddIPs(newIPs)
	return toReturn
}

func (genRange *GenRange) CopyWithIndices(newIndices []int) *GenRange {
	toReturn := genRange.Copy()
	for _, curIndex := range newIndices {
		toReturn.WildIndices[curIndex] = internal.Empty{}
	}
	return toReturn
}

func GetGenRangeFromIPs(fromIPs []*net.IP) *GenRange {
	newRange := newGenRange(fromIPs[0])
	newRange.AddIPs(fromIPs[1:])
	return newRange
}
