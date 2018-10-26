package filtering

import (
	"github.com/willf/bloom"
	"os"
	"errors"
	"fmt"
	"github.com/lavalamp-/ipv666/common/config"
)

func NewFromConfig(conf *config.Configuration) (*bloom.BloomFilter) {
	return bloom.New(conf.AddressFilterSize, conf.AddressFilterHashCount)
}

func GetBloomFilterFromFile(filePath string, filterSize uint, keyCount uint) (*bloom.BloomFilter, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	filter := bloom.New(filterSize, keyCount)
	readCount, err := filter.ReadFrom(file)
	if readCount <= 0 {
		return nil, errors.New(fmt.Sprintf("Read %d from file for Bloom filter initialization.", readCount))
	} else if err != nil {
		return nil, err
	}
	return filter, nil
}

func WriteBloomFilterToFile(filePath string, filter *bloom.BloomFilter) (error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	writeCount, err := filter.WriteTo(file)
	if writeCount <= 0 {
		return errors.New(fmt.Sprintf("Wrote %d bytes to file at path '%s'.", writeCount, filePath))
	} else if err != nil {
		return err
	}
	return nil
}
