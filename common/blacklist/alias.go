package blacklist

import (
	"net"
	"math"
	"errors"
	"fmt"
	"github.com/lavalamp-/ipv666/common/addressing"
	"log"
)

type Empty struct {}

type AliasCheckState struct {
	baseAddress			*net.IP
	leftPosition		uint8
	rightPosition		uint8
	found				bool
	testAddr			*net.IP
}

func NewAliasCheckState(addr *net.IP, left uint8, right uint8) (*AliasCheckState, error) {
	if right > 127 {
		return nil, errors.New(fmt.Sprintf("Right must be less than 128 (got %d).", right))
	}
	if right < left {
		return nil, errors.New(fmt.Sprintf("Right must be greater than or equal to left (got %d, %d).", left, right))
	}
	toReturn := &AliasCheckState{
		baseAddress:	addr,
		leftPosition:	left,
		rightPosition:	right,
		found:			false,
		testAddr:		nil,
	}
	return toReturn, nil
}

// Get the left-most checking index
func (state *AliasCheckState) GetLeft() (uint8) {
	return state.leftPosition
}

// Get the right-mount checking index
func (state *AliasCheckState) GetRight() (uint8) {
	return state.rightPosition
}

// Get whether or not the aliased network length has been found
func (state *AliasCheckState) GetFound() (bool) {
	return state.found
}

// Get the IPv6 address being used to test against for this alias check
func (state *AliasCheckState) GetTestAddr() (*net.IP) {
	return state.testAddr
}

// Get the base IPv6 address that is being permuted against for this alias check
func (state *AliasCheckState) GetBaseAddress() (*net.IP) {
	return state.baseAddress
}

// Get the left index (inclusive) for the next test range
func (state *AliasCheckState) GetLeftTestIndex() (uint8) {
	testDistance := state.GetTestDistance()
	if testDistance <= 1 {
		return state.rightPosition
	} else if testDistance % 2 == 0 {
		return state.leftPosition + (testDistance / 2)
	} else {
		return state.leftPosition + (testDistance / 2) + 1
	}
}

// Get the right index (inclusive) for the next test range
func (state *AliasCheckState) GetRightTestIndex() (uint8) {
	testDistance := state.GetTestDistance()
	if testDistance <= 1 {
		return state.rightPosition
	} else {
		return state.rightPosition - 1
	}
}

// Get the distance between the left and right positions
func (state *AliasCheckState) GetTestDistance() (uint8) {
	return state.rightPosition - state.leftPosition
}

// Get the number of bits that will be tested against in the next round of alias checking
func (state *AliasCheckState) GetTestBitCount() (uint8) {
	leftIndex := state.GetLeftTestIndex()
	rightIndex := state.GetRightTestIndex()
	return rightIndex - leftIndex + 1
}

// Get the total number of potential addresses in the next test range for this alias check
// Note that if there are over 16 bits to test against a boolean value of true will be returned,
// indicating that there are at least 65535 possible addresses
func (state *AliasCheckState) GetPossibleTestAddressCount() (uint64, bool) {
	testBitCount := state.GetTestBitCount()
	if testBitCount > 16 {
		return 0, true
	} else {
		return uint64(math.Pow(2, float64(testBitCount))), false
	}
}

func (state *AliasCheckState) GenerateTestAddress() () {
	state.testAddr = addressing.FlipBitsInAddress(state.baseAddress, state.GetLeftTestIndex(), state.GetRightTestIndex())
}

func (state *AliasCheckState) Update(foundAddrs map[string]*Empty) () {
	//// TODO for set membership checks, i'm guessing strings are expensive. how about 128bit int?
	//// TODO by only checking for a single address, we risk marking ranges as aliased when they aren't. small amount of error, but could be a lot of effort to fix.

	if _, ok := foundAddrs[state.testAddr.String()]; ok {
		// The bit flipped address responded, meaning the range is aliased
		state.rightPosition = state.GetLeftTestIndex()
	} else {
		// The bit flipped address did not respond, meaning the range is not aliased
		state.leftPosition = state.GetLeftTestIndex()
	}

	// If the distance between left and right is 1 then we've found the aliased network threshold
	if state.GetTestDistance() == 1 {
		state.found = true
	}

	// Empty out the list of test addresses to preserve memory
	state.testAddr = nil

}

func (state *AliasCheckState) GetAliasedNetwork() (*net.IPNet, error) {
	if !state.found {
		return nil, errors.New("attempted to get aliased network from AliasCheckState that hadn't found the network yet")
	} else {
		cidrString := fmt.Sprintf("%s/%d", state.baseAddress, state.rightPosition)
		_, ipnet, err := net.ParseCIDR(cidrString)
		return ipnet, err
	}
}

type AliasCheckStates struct {
	checks				[]*AliasCheckState
}

func NewAliasCheckStates(addrs []*net.IP, left uint8, right uint8) (*AliasCheckStates, error) {
	var checkStates []*AliasCheckState
	for _, addr := range addrs {
		newState, err := NewAliasCheckState(addr, left, right)
		if err != nil {
			return nil, err
		}
		checkStates = append(checkStates, newState)
	}
	toReturn := &AliasCheckStates{
		checks:		checkStates,
	}
	return toReturn, nil
}

func (states *AliasCheckStates) GetChecksCount() (int) {
	return len(states.checks)
}

func (states *AliasCheckStates) GetFoundCount() (int) {
	toReturn := 0
	for _, check := range states.checks {
		if check.GetFound() {
			toReturn++
		}
	}
	return toReturn
}

func (states *AliasCheckStates) GenerateTestAddresses() () {
	for _, check := range states.checks {
		check.GenerateTestAddress()
	}
}

func (states *AliasCheckStates) GetAllFound() (bool) {
	return states.GetChecksCount() == states.GetFoundCount()
}

func (states *AliasCheckStates) GetAliasedNetworks() ([]*net.IPNet, error) {
	if !states.GetAllFound() {
		return nil, errors.New("attempted to get aliased networks from an AliasCheckStates when not all checks were found")
	}
	var toReturn []*net.IPNet
	for _, check := range states.checks {
		aliasedNet, err := check.GetAliasedNetwork()
		if err != nil {
			return nil, err
		}
		toReturn = append(toReturn, aliasedNet)
	}
	return toReturn, nil
}

func (states *AliasCheckStates) Update(foundAddrs map[string]*Empty) () {
	for _, check := range states.checks {
		check.Update(foundAddrs)
	}
}

func (states *AliasCheckStates) PrintAliasedNetworks() (error) {
	networks, err := states.GetAliasedNetworks()
	if err != nil {
		return err
	}
	for i, network := range networks {
		log.Printf("Network %d: %s", i, network)
	}
	return nil
}
