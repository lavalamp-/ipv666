package addressing

import (
	"net"
	"os"
	"errors"
	"fmt"
	"io"
)

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
	defer file.Close()
	if err != nil {
		return err
	}
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

