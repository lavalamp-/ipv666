package addressing

import (
	"bufio"
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
	"math/rand"
)

func getFirst64BitsOfNetwork(network *net.IPNet) (uint64) {
	return GetFirst64BitsOfIP(&network.IP)
}

func GenerateRandomNetworks(toGenerate int, minMaskLen int32) ([]*net.IPNet) {
	var toReturn []*net.IPNet
	for len(toReturn) < toGenerate {
		addrBytes := zrandom.GenerateRandomBits(128)
		maskLen := uint8(rand.Int31n(128 - minMaskLen) + minMaskLen)
		newNet, _ := GetIPv6NetworkFromBytes(addrBytes, maskLen)
		toReturn = append(toReturn, newNet)
	}
	return toReturn
}

func GetNetworksFromStrings(toParse []string) ([]*net.IPNet) {
	var toReturn []*net.IPNet
	for i, curParse := range toParse {
		_, network, err := net.ParseCIDR(curParse)
		if err != nil {
			log.Printf("Error thrown when parsing string '%s' as CIDR network (index %d): %e", curParse, i, err)
		} else if network == nil {
			log.Printf("Parsing string '%s' as CIDR network returned empty CIDR network (index %d).", curParse, i)
		} else {
			toReturn = append(toReturn, network)
		}
	}
	return toReturn
}

func GetBaseAddressString(network *net.IPNet) (string) {
	ipBytes := ([]byte)(network.IP)
	maskBytes := ([]byte)(network.Mask)
	var normalized []byte
	for i := range ipBytes {
		normalized = append(normalized, ipBytes[i] & maskBytes[i])
	}
	ip := (net.IP)(normalized)
	ones, _ := network.Mask.Size()
	return fmt.Sprintf("%s/%d", ip, ones)
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

func WriteIPv6NetworksToFile(filePath string, networks []*net.IPNet) (error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	writer := bufio.NewWriter(file)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, network := range networks {
		writer.Write(network.IP)
		ones, _ := network.Mask.Size()
		writer.Write([]byte{uint8(ones)})
	}
	writer.Flush()
	return nil
}

func WriteIPv6NetworksToHexFile(filePath string, networks []*net.IPNet) (error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	writer := bufio.NewWriter(file)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, network := range networks {
		writer.WriteString(fmt.Sprintf("%s\n", network.String()))
	}
	writer.Flush()
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

func GetIPv6NetworkFromBytes(toProcess []byte, maskLength uint8) (*net.IPNet, error) {
	if len(toProcess) != 16 {
		return nil, errors.New(fmt.Sprintf("IPv6 network binary representation must be 16 bytes long (got %d).", len(toProcess)))
	}
	ipBytes := make([]byte, 16)
	copy(ipBytes, toProcess)
	byteMask := GetByteMask(maskLength)
	for i := range ipBytes {
		ipBytes[i] &= byteMask[i]
	}
	toReturn := &net.IPNet{
		IP:			ipBytes,
		Mask:		byteMask,
	}
	return toReturn, nil
}

func GetIPv6NetworkFromBytesIncLength(toProcess []byte) (*net.IPNet, error) {
	if len(toProcess) != 17 {
		return nil, errors.New(fmt.Sprintf("IPv6 network binary representation must be 17 bytes long (got %d).", len(toProcess)))
	}
	ipBytes := make([]byte, 16)
	copy(ipBytes, toProcess)
	maskLength := uint8(toProcess[16])
	return GetIPv6NetworkFromBytes(ipBytes, maskLength)
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
		newNetwork, err := GetIPv6NetworkFromBytesIncLength(buffer)
		if err != nil {
			return nil, err
		}
		toReturn = append(toReturn, newNetwork)
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

func GetNetworkFromUints(uints [2]uint64, length uint8) (*net.IPNet) {
	//TODO there has to be a better way to do this, esp with the creating a mask approach
	var addrBytes []byte
	processBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(processBytes, uints[0])
	addrBytes = append(addrBytes, processBytes...)
	binary.BigEndian.PutUint64(processBytes, uints[1])
	addrBytes = append(addrBytes, processBytes...)
	toReturn, _ := GetIPv6NetworkFromBytes(addrBytes, length)
	return toReturn
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
