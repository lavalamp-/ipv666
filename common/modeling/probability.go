package modeling

import (
	"math/rand"
	"log"
	"github.com/lavalamp-/ipv666/common/persist"
	"github.com/lavalamp-/ipv666/common/addressing"
	"net"
	"github.com/lavalamp-/ipv666/common/config"
	"math"
)

type addrProcessFunc func(*net.IP) (bool, error)

type ProbabilisticAddressModel struct {
	Name 			string 							`json:"name"`
	DigestCount 	uint64 							`json:"digest"`
	NybbleModels 	[31]*ProbabilisticNybbleModel 	`json:"models"`
}

func NewAddressModel(name string, conf *config.Configuration) (*ProbabilisticAddressModel) {
	var models [31]*ProbabilisticNybbleModel
	for i := 0; i < 31; i++ {
		models[i] = newNybbleModel(conf)
	}
	return &ProbabilisticAddressModel{
		name,
		0,
		models,
	}
}

type ProbabilisticNybbleModel struct {
	DigestCount 			uint64 							`json:"digest"`
	Probabilities 			map[uint8]*NybbleProbabilityMap	`json:"probmap"`
	DefaultProbabilityMap 	*NybbleProbabilityMap 			`json:"defaultprobmap"`
}

func newNybbleModel(conf *config.Configuration) (*ProbabilisticNybbleModel) {
	return &ProbabilisticNybbleModel{
		0,
		make(map[uint8]*NybbleProbabilityMap),
		newProbabilityMap(conf),
	}
}

type NybbleProbabilityMap struct {
	DigestCount 	uint64 				`json:"digest"`
	LastUpdatedAt 	uint64 				`json:"updated"`
	TimesSeen 		map[uint8]uint64 	`json:"occurences"`
	Distribution 	[]uint8 			`json:"dist"`
}

func newProbabilityMap(conf *config.Configuration) (*NybbleProbabilityMap) {
	toReturn := &NybbleProbabilityMap{
		0,
		0,
		make(map[uint8]uint64),
		nil,
	}
	var i uint8
	for i = 0; i < 16; i++ {
		toReturn.TimesSeen[i] = conf.ModelDefaultWeight
	}
	toReturn.DigestCount = conf.ModelDefaultWeight * 16
	return toReturn
}

func GetProbabilisticModelFromFile(filePath string) (*ProbabilisticAddressModel, error) {
	var toReturn ProbabilisticAddressModel
	err := persist.Load(filePath, &toReturn)
	return &toReturn, err
}

func (addrModel *ProbabilisticAddressModel) Save(filePath string) (error) {
	return persist.Save(filePath, addrModel)
}

func (addrModel *ProbabilisticAddressModel) GenerateMultiIPFromNetwork(fromNetwork *net.IPNet, count int, fn addrProcessFunc) ([]*net.IP, error) {
	var toReturn []*net.IP
	netLength, _ := fromNetwork.Mask.Size()
	nybbleCount := int(math.Ceil(float64(netLength) / 4.0))
	fromNybbles := addressing.GetNybblesFromIP(&fromNetwork.IP, nybbleCount)
	totalCount := 0
	for len(toReturn) < count {
		newIP := addrModel.GenerateSingleIPFromNybbles(fromNybbles, uint(netLength))
		isFiltered, err := fn(newIP)
		if err != nil {
			return nil, err
		} else if !isFiltered {
			toReturn = append(toReturn, newIP)
		}
		totalCount++
	}
	return toReturn, nil
}

func (addrModel *ProbabilisticAddressModel) GenerateSingleIPFromNybbles(fromNybbles []byte, offset uint) (*net.IP) {
	//TODO this will throw an exception if offset < 4
	//TODO if offset does not correspond with length of fromNybbles this will error
	var mustMatch bool
	addrNybles := make([]byte, len(fromNybbles))
	copy(addrNybles, fromNybbles)
	if offset % 4 == 0 {
		mustMatch = false
	} else {
		mustMatch = true
	}
	modelOffset := (offset / 4) - 1
	curNybble := addrNybles[modelOffset]
	for len(addrNybles) < 32 {
		nybbleModel := addrModel.NybbleModels[modelOffset]
		nextNybble := nybbleModel.predictNextNybble(curNybble)
		if mustMatch {
			// Check to make sure the predicted nybble starts with the expected bits
			if (^(0xff >> (offset % 4) + 4) & nextNybble) ^ addrNybles[modelOffset + 1] == 0 {
				addrNybles[modelOffset + 1] = nextNybble
				curNybble = nextNybble
				mustMatch = false
				modelOffset++
			}
		} else {
			addrNybles = append(addrNybles, nextNybble)
			curNybble = nextNybble
			modelOffset++
		}
	}
	var addrBytes [16]byte
	for i := 0; i < 16; i++ {
		nybbleIndex := i * 2
		addrBytes[i] = (addrNybles[nybbleIndex] << 4) | addrNybles[nybbleIndex + 1]
	}
	var newIP = (net.IP)(addrBytes[:])
	return &newIP
}

func (addrModel *ProbabilisticAddressModel) GenerateMultiIPFromNybble(fromNybble uint8, count int, updateFreq int) ([]*net.IP) {
	var toReturn []*net.IP
	log.Printf("Generating %d IP addresses using model %s.", count, addrModel.Name)
	for i := 0; i < count; i++ {
		if i % updateFreq == 0 {
			log.Printf("Generating %d addresses out of %d.", i, count)
		}
		toReturn = append(toReturn, addrModel.GenerateSingleIPFromNybble(fromNybble))
	}
	log.Printf("Successfully generated %d IP addresses using model %s.", count, addrModel.Name)
	return toReturn
}

func (addrModel *ProbabilisticAddressModel) GenerateSingleIPFromNybble(fromNybble uint8) (*net.IP) {
	addrNybbles := []uint8{fromNybble}
	curNybble := fromNybble
	for _, nybbleModel := range addrModel.NybbleModels {
		nextNybble := nybbleModel.predictNextNybble(curNybble)
		addrNybbles = append(addrNybbles, nextNybble)
		curNybble = nextNybble
	}
	var addrBytes [16]byte
	for i := 0; i < 16; i++ {
		nybbleIndex := i * 2
		addrBytes[i] = (addrNybbles[nybbleIndex] << 4) | addrNybbles[nybbleIndex + 1]
	}
	var newIP = (net.IP)(addrBytes[:])
	return &newIP
}

func generateAddressModel(ips []*net.IP, name string, updateInterval int, conf *config.Configuration) (*ProbabilisticAddressModel) {
	toReturn := NewAddressModel(name, conf)
	toReturn.UpdateMultiIP(ips, updateInterval, conf)
	return toReturn
}

func (addrModel *ProbabilisticAddressModel) UpdateMultiIP(ips []*net.IP, updateInterval int, conf *config.Configuration) () {
	log.Printf("Updating model %s with %d addresses.", addrModel.Name, len(ips))
	for i, ip := range ips {
		if i % updateInterval == 0 {
			log.Printf("Processing address %d out of %d.", i, len(ips))
		}
		addrModel.UpdateSingleIP(ip, conf)
	}
	log.Printf("Successfully updated model %s with %d addresses.", addrModel.Name, len(ips))
}

func (addrModel *ProbabilisticAddressModel) UpdateSingleIP(ip *net.IP, conf *config.Configuration) () {
	fromNybble := addressing.GetNybbleFromIP(ip, 0)
	for i, nybbleModel := range addrModel.NybbleModels {
		toNybble := addressing.GetNybbleFromIP(ip, i+1)
		nybbleModel.update(fromNybble, toNybble, conf)
		fromNybble = toNybble
	}
	addrModel.DigestCount += 1
}

func (nybbleModel *ProbabilisticNybbleModel) predictNextNybble(fromNybble uint8) (uint8) {
	if val, ok := nybbleModel.Probabilities[fromNybble]; ok {
		return val.predictNextNybble()
	} else {
		return nybbleModel.DefaultProbabilityMap.predictNextNybble()
	}
}

func (nybbleModel *ProbabilisticNybbleModel) update(fromNybble uint8, toNybble uint8, conf *config.Configuration) () {
	if val, ok := nybbleModel.Probabilities[fromNybble]; ok {
		val.update(toNybble)
	} else {
		newMap := newProbabilityMap(conf)
		newMap.update(toNybble)
		nybbleModel.Probabilities[fromNybble] = newMap
	}
	nybbleModel.DefaultProbabilityMap.update(toNybble)
	nybbleModel.DigestCount += 1
}

func (probMap *NybbleProbabilityMap) update(nybble uint8) () {
	if _, ok := probMap.TimesSeen[nybble]; ok {
		probMap.TimesSeen[nybble] += 1
	} else {
		probMap.TimesSeen[nybble] = 1
	}
	probMap.DigestCount += 1
}

func (probMap *NybbleProbabilityMap) predictNextNybble() (uint8) {
	if !probMap.isModelUpdated() {
		probMap.buildDistribution()
	}
	return probMap.Distribution[rand.Intn(len(probMap.Distribution))]
}

func (probMap *NybbleProbabilityMap) isModelUpdated() (bool) {
	return probMap.DigestCount == probMap.LastUpdatedAt
}

func (probMap *NybbleProbabilityMap) buildDistribution () {
	var newDistribution []uint8
	for k, v := range probMap.TimesSeen {
		probability := float32(v) / float32(probMap.DigestCount)
		nybbleCount := int(probability * 100)
		for i := 0; i < nybbleCount; i++ {
			newDistribution = append(newDistribution, k)
		}
	}
	probMap.Distribution = newDistribution
	probMap.LastUpdatedAt = probMap.DigestCount
}

func CreateBlankModel(name string, outputPath string, conf *config.Configuration) (error) {
	log.Printf("Now creating a blank statistical model.")
	model := NewAddressModel(name, conf)
	log.Printf("Writing blank statistical model with name '%s' to file '%s'.", model.Name, outputPath)
	err := model.Save(outputPath)
	if err != nil {
		log.Printf("Error thrown when saving model '%s' to file '%s': %e", model.Name, outputPath, err)
		return err
	}
	return nil
}
