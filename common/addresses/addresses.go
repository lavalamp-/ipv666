package addresses

import (
	"os"
	"io/ioutil"
	"log"
	"errors"
	"strings"
	"regexp"
	"fmt"
	"net"
)

type IPv6Address struct {
	Content [16]byte
}

type IPv6AddressList struct {
	Addresses []IPv6Address
}

func NewIPv6Address(bytes [16]byte) IPv6Address {
	return IPv6Address{bytes}
}

func NewIPv6AddressList(addresses []IPv6Address) IPv6AddressList {
	return IPv6AddressList{addresses}
}

func (address IPv6Address) GetNybble(index int) (uint8) {
	// TODO fatal error if index > 31
	byteIndex := index / 2
	addrByte := address.Content[byteIndex]
	if index % 2 == 0 {
		return addrByte >> 4
	} else {
		return addrByte & 0xf
	}
}

func (address *IPv6Address) GetIP() (*net.IP) {
	return &net.IP{
		address.Content[0],
		address.Content[1],
		address.Content[2],
		address.Content[3],
		address.Content[4],
		address.Content[5],
		address.Content[6],
		address.Content[7],
		address.Content[8],
		address.Content[9],
		address.Content[10],
		address.Content[11],
		address.Content[12],
		address.Content[13],
		address.Content[14],
		address.Content[15],
	}
}

func (address *IPv6Address) String() (string) {
	return address.GetIP().String()
}

func (list IPv6AddressList) ToAddressesFile(filePath string, updateFreq int) (error) {

	_, err := os.Stat(filePath)
	if err == nil {
		return errors.New(fmt.Sprintf("A file already exists at %s.", filePath))
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("Now writing %d IPv6 addresses to file at %s.", len(list.Addresses), filePath)

	for i, address := range(list.Addresses) {  // TODO optimize this?

		if i % updateFreq == 0 {
			log.Printf("Writing address %d out of %d.", i, len(list.Addresses))
		}

		f.WriteString(fmt.Sprintf("%s\n", address.String()))
	}

	f.Sync()

	log.Printf("Finished writing %d IPv6 addresses to file at %s.", len(list.Addresses), filePath)

	return nil

}

func (list IPv6AddressList) ToBinaryFile(filePath string, updateFreq int) (error) {

	_, err := os.Stat(filePath)
	if err == nil {
		return errors.New(fmt.Sprintf("A file already exists at %s.", filePath))
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("Now writing %d binary addresses to file at %s.", len(list.Addresses), filePath)

	for i, address := range(list.Addresses) {  // TODO optimize this?

		if i % updateFreq == 0 {
			log.Printf("Writing address %d out of %d.", i, len(list.Addresses))
		}

		f.Write(address.Content[:])
	}

	f.Sync()

	log.Printf("Finished writing binary addresses to file at %s.", filePath)

	return nil

}

func GetAddressFromBitString(bitstring string) (IPv6Address, error) {

	regex, err := regexp.Compile("^[01]{128}$")
	if err != nil {
		return IPv6Address{}, err
	}

	if !regex.MatchString(bitstring) {
		return IPv6Address{}, errors.New(fmt.Sprintf("%s is not a valid IPv6 address bit string.", bitstring))
	}

	var bytes []byte
	var curByte byte = 0
	var curLength = 0

	for pos, _ := range(bitstring) {
		curByte <<= 1
		if bitstring[pos] == "0"[0] {
			curByte |= 0
		} else {
			curByte |= 1
		}
		curLength += 1
		if curLength == 8 {
			bytes = append(bytes, curByte)
			curByte = 0
			curLength = 0
		}
	}

	var byteArray [16]byte
	copy(byteArray[:], bytes)

	return NewIPv6Address(byteArray), nil

}

func GetAddressListFromBytes(bytes []byte) (IPv6AddressList, error) {

	if len(bytes) % 16 != 0 {
		return IPv6AddressList{}, errors.New("Length of bytes did not end on a 16 byte boundary.")
	}

	var addresses []IPv6Address
	var curAddr [16]byte

	for i := 0; i < len(bytes); i += 16 {
		copy(curAddr[:], bytes[i:i+16])
		addresses = append(addresses, NewIPv6Address(curAddr))
	}

	return NewIPv6AddressList(addresses), nil

}

func GetAddressListFromAddressesFile(filePath string) (IPv6AddressList, error) {
	return IPv6AddressList{}, nil
}

func GetAddressListFromBitStringsFile(filePath string) (IPv6AddressList, error) {

	log.Printf("Checking that file exists at %s.", filePath)
	_, err := os.Stat(filePath)
	if err != nil {
		return IPv6AddressList{}, err
	}

	log.Printf("Reading Content of file at %s.", filePath)
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return IPv6AddressList{}, err
	}

	contentString := strings.TrimSpace(string(fileContent))
	lines := strings.Split(contentString, "\n")
	var cleanedLines []string
	for _, s := range(lines) {
		cleanedLines = append(cleanedLines, strings.TrimSpace(s))
	}

	var addresses []IPv6Address

	for pos, cleanedLine := range(cleanedLines) {
		address, err := GetAddressFromBitString(cleanedLine)
		if err != nil {
			return IPv6AddressList{}, errors.New(fmt.Sprintf("Error at line %s: %s.", pos, err))
		}
		addresses = append(addresses, address)
	}

	return NewIPv6AddressList(addresses), nil

}

func GetAddressListFromBinaryFile(filePath string) (IPv6AddressList, error) {

	log.Printf("Checking that file exists at %s.", filePath)
	_, err := os.Stat(filePath)
	if err != nil {
		return IPv6AddressList{}, err
	}

	log.Printf("Reading Content of file at %s.", filePath)
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return IPv6AddressList{}, err
	}

	toReturn, err := GetAddressListFromBytes(fileContent)
	if err != nil {
		return IPv6AddressList{}, err
	}

	log.Printf("A total of %d Addresses were loaded from file at %s.", len(toReturn.Addresses), filePath)

	return toReturn, nil

}
