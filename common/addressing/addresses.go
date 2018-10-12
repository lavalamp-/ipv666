package addressing

import (
	"os"
	"io/ioutil"
	"errors"
	"strings"
	"fmt"
	"net"
	"io"
)

func ReadIPsFromHexFile(filePath string) ([]*net.IP, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	contentString := strings.TrimSpace(string(fileContent))
	lines := strings.Split(contentString, "\n")
	var toReturn []*net.IP
	for _, line := range lines {
		newIP := net.ParseIP(strings.TrimSpace(line))
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
