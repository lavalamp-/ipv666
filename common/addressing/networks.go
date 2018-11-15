package addressing

import (
	"net"
	"os"
	"errors"
	"fmt"
	"io"
	"github.com/lavalamp-/ipv666/common/zrandom"
	"log"
	"io/ioutil"
	"strings"
	"encoding/binary"
)

func getFirst64BitsOfNetwork(network *net.IPNet) (uint64) {
	return GetFirst64BitsOfIP(&network.IP)
}

func GetBaseAddressString(network *net.IPNet) (string) {
	ipBytes := ([]byte)(network.IP)
	maskBytes := ([]byte)(network.Mask)
	var normalized []byte
	for i := range ipBytes {
		normalized = append(normalized, ipBytes[i] & maskBytes[i])
	}
	ip := (net.IP)(normalized)
	return ip.String()
}

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

func GetUniqueNetworks(networks []*net.IPNet, updateFreq int) ([]*net.IPNet) {
	checkMap := make(map[string]bool)
	var toReturn []*net.IPNet
	for i, curNet := range networks {
		if i % updateFreq == 0 {
			log.Printf("Processing %d out of %d for unique networks.", i, len(networks))
		}
		netString := GetBaseAddressString(curNet)
		if _, ok := checkMap[netString]; !ok {
			checkMap[netString] = true
			toReturn = append(toReturn, curNet)
		}
	}
	return toReturn
}

func checkNetworkEquality(first *net.IPNet, second *net.IPNet) (bool) {
	//TODO fix this - checking IPs can fail in this way as the IPNet struct doesn't mask out the address
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

func ReadIPv6NetworksFromHexFile(filePath string) ([]*net.IPNet, error) {
	log.Printf("Reading IPv6 networks from file at path '%s'.", filePath)
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fileContent := strings.TrimSpace(string(fileBytes))
	fileLines := strings.Split(fileContent, "\n")
	var networks []*net.IPNet
	for _, fileLine := range fileLines {
		trimmedLine := strings.TrimSpace(fileLine)
		_, network, err := net.ParseCIDR(trimmedLine)
		if err != nil {
			log.Printf("Error thrown when parsing '%s' as CIDR: %e", trimmedLine, err)
			continue
		}
		networks = append(networks, network)
	}
	log.Printf("Read %d networks from file '%s'.", len(networks), filePath)
	return networks, nil
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

func GetNetworkFromUints(uints [2]uint64, length int) (*net.IPNet) {
	//TODO there has to be a better way to do this, esp with the creating a mask approach
	var addrBytes []byte
	processBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(processBytes, uints[0])
	addrBytes = append(addrBytes, processBytes...)
	binary.BigEndian.PutUint64(processBytes, uints[1])
	addrBytes = append(addrBytes, processBytes...)
	maskBytes := GetByteMask(uint8(length))
	for i := range addrBytes {
		addrBytes[i] = addrBytes[i] & maskBytes[i]
	}
	return &net.IPNet{
		IP:		addrBytes,
		Mask:	maskBytes,
	}
}

func GetBorderAddressesFromNetwork(network *net.IPNet) (*net.IP, *net.IP) {
	var baseAddrBytes []byte
	var topAddrBytes []byte
	for i := range network.IP {
		baseAddrBytes = append(baseAddrBytes, network.IP[i] & network.Mask[i])
		topAddrBytes = append(topAddrBytes, baseAddrBytes[i] | ^network.Mask[i])
	}
	baseAddr := net.IP(baseAddrBytes)
	topAddr := net.IP(topAddrBytes)
	return &baseAddr, &topAddr
}
