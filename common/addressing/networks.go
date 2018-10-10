package addressing

import (
	"net"
	"os"
	"errors"
	"fmt"
	"io"
	"github.com/lavalamp-/ipv666/common/zrandom"
)

func GenerateRandomAddressesInNetwork(network *net.IPNet, addrCount int) ([]*net.IP) {
	var existsMap = make(map[string]bool)
	var toReturn []*net.IP
	for len(toReturn) < addrCount {
		newAddr := GenerateRandomAddressInNetwork(network)
		if _, ok := existsMap[newAddr.String()]; !ok {
			toReturn = append(toReturn, newAddr)
			existsMap[newAddr.String()] = true
		}
	}
	return toReturn
}

func GenerateRandomAddressInNetwork(network *net.IPNet) (*net.IP) {
	ones, _ := network.Mask.Size()
	randomBytes := zrandom.GenerateHostBits(128 - ones)
	var newBytes []byte
	for i := range network.IP {
		newBytes = append(newBytes, (network.IP[i] & network.Mask[i]) | randomBytes[i])
	}
	var genIP = net.IP(newBytes)
	return &genIP
}

func GetUniqueNetworks(networks []*net.IPNet) ([]*net.IPNet) {
	var toReturn []*net.IPNet
	for i, curNet := range networks {
		found := false
		for _, checkNet := range toReturn {
			if CheckNetworkEquality(curNet, checkNet) {
				found = true
				break
			}
		}
		if !found {
			toReturn = append(toReturn, networks[i])
		}
	}
	return toReturn
}

func CheckNetworkEquality(first *net.IPNet, second *net.IPNet) (bool) {
	for i := range first.IP {
		if first.IP[i] != second.IP[i] {
			return false
		}
	}
	for i := range first.Mask {
		if first.Mask[i] != second.Mask[i] {
			return false
		}
	}
	return true
}

func WriteIPv6NetworksToFile(filePath string, networks []*net.IPNet) (error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, network := range networks {
		file.Write(network.IP)
		ones, _ := network.Mask.Size()
		foo := uint8(ones)
		file.Write([]byte{foo})
	}
	return nil
}

func ReadIPv6NetworksFromFile(filePath string) ([]*net.IPNet, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()
	if fileSize % 17 != 0 {
		return nil, errors.New(fmt.Sprintf("Expected file size to be a multiple of 17 (got %d).", fileSize))
	}
	buffer := make([]byte, 17)
	var toReturn []*net.IPNet
	for {
		_, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		maskLength := uint8(buffer[16])
		byteMask := GetByteMask(maskLength)
		ipBytes := make([]byte, 16)
		copy(ipBytes, buffer)
		toReturn = append(toReturn, &net.IPNet{
			IP:			ipBytes,
			Mask:		byteMask,
		})
	}
	return toReturn, nil
}

func GetByteMask(maskLength uint8) ([]byte) {
	var toReturn []byte
	byteCount := maskLength / 8
	var i uint8
	for i = 0; i < byteCount; i++ {
		toReturn = append(toReturn, 0xff)
	}
	bitOff := (uint)(maskLength % 8)
	if bitOff != 0 {
		toReturn = append(toReturn, GetByteWithBitsMasked(bitOff))
	}
	for len(toReturn) < 16 {
		toReturn = append(toReturn,0x00)
	}
	return toReturn
}

func GetByteWithBitsMasked(bitMaskLength uint) (byte) {
	return (byte)(^(0xff >> bitMaskLength))
}

