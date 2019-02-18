package comparison

import "fmt"

func GetHammingDistance(firstBytes []byte, secondBytes []byte) (int, error) {
	if len(firstBytes) != len(secondBytes) {
		return -1, fmt.Errorf("input byte arrays were of differing lengths (%d and %d)", len(firstBytes), len(secondBytes))
	}
	toReturn := 0
	for i := range firstBytes {
		if firstBytes[i] != secondBytes[i] {
			toReturn++
		}
	}
	return toReturn, nil
}
