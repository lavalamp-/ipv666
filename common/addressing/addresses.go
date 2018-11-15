package addressing

import (
	"os"
	"io/ioutil"
	"errors"
	"strings"
	"fmt"
	"net"
	"io"
	"log"
	"encoding/binary"
	"github.com/lavalamp-/ipv666/common/zrandom"
	"github.com/lavalamp-/ipv666/common"
)

func GetIPSet(ips []*net.IP) (map[string]*common.Empty) {
	toReturn := make(map[string]*common.Empty)
	blacklistEntry := &common.Empty{}
	for _, ip := range ips {
		toReturn[ip.String()] = blacklistEntry
	}
	return toReturn
}

func GetFirst64BitsOfIP(ip *net.IP) (uint64) {
	ipBytes := ([]byte)(*ip)
	return binary.LittleEndian.Uint64(ipBytes[:8])
}

func GetUniqueIPs(ips []*net.IP, updateFreq int) ([]*net.IP) {
	checkMap := make(map[string]bool)
	var toReturn []*net.IP
	for i, ip := range ips {
		if i % updateFreq == 0 {
			log.Printf("Processing %d out of %d for unique IPs.", i, len(ips))
		}
		if _, ok := checkMap[ip.String()]; !ok {
			checkMap[ip.String()] = true
			toReturn = append(toReturn, ip)
		}
	}
	return toReturn
}

func ReadIPsFromHexFile(filePath string) ([]*net.IP, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	contentString := strings.TrimSpace(string(fileContent))
	lines := strings.Split(contentString, "\n")
	var toReturn []*net.IP
	for i, line := range lines {
		newIP := net.ParseIP(strings.TrimSpace(line))
		if newIP == nil {
			log.Printf("No IP found from content '%s' (line %d in file '%s').", line, i, filePath)
			continue
		}
		toReturn = append(toReturn, &newIP)
	}
	return toReturn, nil
}

func WriteIPsToHexFile(filePath string, addrs []*net.IP) (error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, addr := range addrs {
		file.WriteString(fmt.Sprintf("%s\n", addr.String()))
	}
	return nil
}

func GetTextLinesFromIPs(addrs []*net.IP) (string) {
	var toReturn []string
	for _, addr := range addrs {
		toReturn = append(toReturn, fmt.Sprintf("%s\n", addr.String()))
	}
	return strings.Join(toReturn, "")
}

func ReadIPsFromBinaryFile(filePath string) ([]*net.IP, error) {
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
	if fileSize % 16 != 0 {
		return nil, errors.New(fmt.Sprintf("Expected file size to be a multiple of 16 (got %d).", fileSize))
	}
	buffer := make([]byte, 16)
	var toReturn []*net.IP
	for {
		_, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		ipBytes := make([]byte, 16)
		copy(ipBytes, buffer)
		newIP := net.IP(ipBytes)
		toReturn = append(toReturn, &newIP)
	}
	return toReturn, nil
}

func WriteIPsToBinaryFile(filePath string, addrs []*net.IP) (error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, addr := range addrs {
		file.Write(*addr)
	}
	return nil
}

func GetNybbleFromIP(ip *net.IP, index int) (uint8) {
	// TODO fatal error if index > 31
	byteIndex := index / 2
	addrBytes := ([]byte)(*ip)
	addrByte := addrBytes[byteIndex]
	if index % 2 == 0 {
		return addrByte >> 4
	} else {
		return addrByte & 0xf
	}
}

func GenerateRandomAddress() (*net.IP) {
	bytes := zrandom.GenerateHostBits(128)
	toReturn := net.IP(bytes)
	return &toReturn
}

func FlipBitsInAddress(toFlip *net.IP, startIndex uint8, endIndex uint8) (*net.IP) {
	toFlipBytes := *toFlip
	endIndex++
	startByte := startIndex / 8
	startOffset := startIndex % 8
	endByte := endIndex / 8
	endOffset := endIndex % 8
	var maskBytes []byte
	var flipBytes []byte
	var i uint8

	if startByte == endByte {
		for i = 0; i < 16; i++ {
			if i == startByte {
				firstHalf := byte(^(0xff >> startOffset))
				secondHalf := byte(0xff >> endOffset)
				maskBytes = append(maskBytes, firstHalf | secondHalf)
			} else {
				maskBytes = append(maskBytes, 0xff)
			}
		}
	} else {
		for i = 0; i < 16; i++ {
			if i < startByte {
				maskBytes = append(maskBytes, 0xff)
			} else if i == startByte {
				maskBytes = append(maskBytes, byte(^(0xff >> startOffset)))
			} else if i < endByte {
				maskBytes = append(maskBytes, 0x00)
			} else if i == endByte {
				maskBytes = append(maskBytes, byte(0xff >> endOffset))
			} else {
				maskBytes = append(maskBytes, 0xff)
			}
		}
	}

	for i = 0; i < 16; i++ {
		flippedBits := ^toFlipBytes[i] & ^maskBytes[i]
		flipBytes = append(flipBytes, toFlipBytes[i] & maskBytes[i] | flippedBits)
	}

	toReturn := net.IP(flipBytes)
	return &toReturn

}
