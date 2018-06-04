package modeling

import (
	"math/rand"
	"github.com/lavalamp-/ipv666/common"
	"log"
)

type ProbabilisticAddressModel struct {
	Name string `json:"name"`
	DigestCount uint64 `json:"digest"`
	NybbleModels [31]*ProbabilisticNybbleModel `json:"models"`
}

func newAddressModel(name string) (*ProbabilisticAddressModel) {
	var models [31]*ProbabilisticNybbleModel
	for i := 0; i < 31; i++ {
		models[i] = newNybbleModel()
	}
	return &ProbabilisticAddressModel{
		name,
		0,
		models,
	}
}

type ProbabilisticNybbleModel struct {
	DigestCount uint64 `json:"digest"`
	Probabilities map[uint8]*NybbleProbabilityMap `json:"probmap"`
	DefaultProbabilityMap *NybbleProbabilityMap `json:"defaultprobmap"`
}

func newNybbleModel() (*ProbabilisticNybbleModel) {
	return &ProbabilisticNybbleModel{
		0,
		make(map[uint8]*NybbleProbabilityMap),
		newProbabilityMap(),
	}
}

type NybbleProbabilityMap struct {
	DigestCount uint64 `json:"digest"`
	LastUpdatedAt uint64 `json:"updated"`
	TimesSeen map[uint8]uint64 `json:"occurences"`
	Distribution []uint8 `json:"dist"`
}

func newProbabilityMap() (*NybbleProbabilityMap) {
	return &NybbleProbabilityMap{
		0,
		0,
		make(map[uint8]uint64),
		nil,
	}
}

func GetProbablisticModelFromFile(filePath string) (ProbabilisticAddressModel, error) {
	var toReturn ProbabilisticAddressModel
	err := common.Load(filePath, &toReturn)
	return toReturn, err
}

func (addrModel *ProbabilisticAddressModel) Save(filePath string) (error) {
	return common.Save(filePath, addrModel)
}

func (addrModel *ProbabilisticAddressModel) GenerateMulti(fromNybble uint8, count int) (common.IPv6AddressList) {
	var addresses []common.IPv6Address

	log.Printf("Generating %s IP addresses using model %s.", count, addrModel.Name)

	for i := 0; i < count; i++ {

		if i % 10000 == 0 {
			log.Printf("Generating %d addresses out of %d.", i, count)
		}

		addresses = append(addresses, addrModel.GenerateSingle(fromNybble))
	}

	log.Printf("Successfully generated %d IP addresses using model %s.", count, addrModel.Name)

	return common.NewIPv6AddressList(addresses)
}

func (addrModel *ProbabilisticAddressModel) GenerateSingle(fromNybble uint8) (common.IPv6Address) {
	addrNybbles := []uint8{fromNybble}
	curNybble := fromNybble
	for _, nybbleModel := range(addrModel.NybbleModels) {
		nextNybble := nybbleModel.predictNextNybble(curNybble)
		addrNybbles = append(addrNybbles, nextNybble)
		curNybble = nextNybble
	}
	var addrBytes [16]byte
	for i := 0; i < 16; i++ {
		nybbleIndex := i * 2
		addrBytes[i] = (addrNybbles[nybbleIndex] << 4) | addrNybbles[nybbleIndex + 1]
	}
	return common.NewIPv6Address(addrBytes)
}

func GenerateAddressModel(addresses common.IPv6AddressList, name string) (*ProbabilisticAddressModel) {
	toReturn := newAddressModel(name)
	toReturn.UpdateMulti(addresses)
	return toReturn
}

func (addrModel *ProbabilisticAddressModel) UpdateMulti(addresses common.IPv6AddressList) () {

	log.Printf("Updating model %s with %d addresses.", addrModel.Name, len(addresses.Addresses))

	for i, address := range(addresses.Addresses) {

		if i % 10000 == 0 {
			log.Printf("Processing address %d out of %d.", i, len(addresses.Addresses))
		}

		addrModel.UpdateSingle(address)
	}

	log.Printf("Successfully updated model %s with %d addresses.", addrModel.Name, len(addresses.Addresses))

}

func (addrModel *ProbabilisticAddressModel) UpdateSingle(address common.IPv6Address) () {
	fromNybble := address.GetNybble(0)
	for i, nybbleModel := range(addrModel.NybbleModels) {
		toNybble := address.GetNybble(i+1)
		nybbleModel.update(fromNybble, toNybble)
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

func (nybbleModel *ProbabilisticNybbleModel) update(fromNybble uint8, toNybble uint8) () {
	if val, ok := nybbleModel.Probabilities[fromNybble]; ok {
		val.update(toNybble)
	} else {
		newMap := newProbabilityMap()
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
	if probMap.DigestCount != probMap.LastUpdatedAt {
		return false
	} else {
		return true
	}
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
