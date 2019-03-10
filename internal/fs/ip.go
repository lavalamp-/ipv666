package fs

import (
	"encoding/hex"
	"errors"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/modeling"
	"github.com/lavalamp-/ipv666/internal/persist"
	"io/ioutil"
	"net"
	"strings"
)

func ReadIPsFromFile(filePath string) ([]*net.IP, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return ParseIPsFromBytes(bytes)
}

func ReadIPsFromAddressTreeBytes(toParse []byte) ([]*net.IP, error) {
	var tree modeling.AddressTree
	err := persist.Unmarshal(toParse, &tree)
	if err != nil {
		return nil, err
	} else {
		return tree.GetAllIPs(), nil
	}
}

func ReadIPsFromHexFileBytes(toParse []byte) []*net.IP {
	parseString := strings.TrimSpace(string(toParse))
	lines := strings.Split(parseString, "\n")
	var toReturn []*net.IP
	for _, line := range lines {
		newIP := net.ParseIP(strings.TrimSpace(line))
		if newIP == nil {
			logging.Warnf("No IP found from content '%s'.", line)
			continue
		}
		toReturn = append(toReturn, &newIP)
	}
	return toReturn
}

func fatHexStringToIP(toParse string) (*net.IP, error) {
	data, err := hex.DecodeString(toParse)
	if err != nil {
		return nil, err
	}
	ip := net.IP(data)
	return &ip, nil
}

func ReadIPsFromFatHexFileBytes(toParse []byte) []*net.IP {
	parseString := strings.TrimSpace(string(toParse))
	lines := strings.Split(parseString, "\n")
	var toReturn []*net.IP
	for _, line := range lines {
		lineStrip := strings.TrimSpace(line)
		newIp, err := fatHexStringToIP(lineStrip)
		if err != nil {
			logging.Warnf("Error thrown when processing bytes %s as fat hex: %s", line, err.Error())
		} else {
			toReturn = append(toReturn, newIp)
		}
	}
	return toReturn
}

func ReadIPsFromBinaryFileBytes(toParse []byte) []*net.IP {
	var toReturn []*net.IP
	for i := 0; i < len(toParse); i += 16 {
		ipBytes := make([]byte, 16)
		copy(ipBytes, toParse[i:i+16])
		newIP := net.IP(ipBytes)
		toReturn = append(toReturn, &newIP)
	}
	return toReturn
}

func ParseIPsFromBytes(toParse []byte) ([]*net.IP, error) {
	split := strings.Split(string(toParse), "\n")
	toCheck := split[0]
	if strings.Contains(toCheck, ":") {  // Standard ASCII hex with colons
		return ReadIPsFromHexFileBytes(toParse), nil
	} else if len(toCheck) == 32 {  // ASCII hex without colons
		return ReadIPsFromFatHexFileBytes(toParse), nil
	} else if len(toParse) % 16 == 0 {  // Binary representation
		return ReadIPsFromBinaryFileBytes(toParse), nil
	} else { // IP address tree format
		result, err := ReadIPsFromAddressTreeBytes(toParse)
		if err != nil {
			return nil, errors.New("could not determine the format of IPv6 address bytes")
		} else {
			return result, nil
		}
	}
}

func ReadIPsFromHexFile(filePath string) ([]*net.IP, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return ReadIPsFromHexFileBytes(fileContent), nil
}

