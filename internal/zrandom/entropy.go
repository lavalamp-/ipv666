package zrandom

import "math"

func GetEntropyOfBitsFromRight(bytes []byte, count int) (float64) {

	// https://stackoverflow.com/questions/990477/how-to-calculate-the-entropy-of-a-file

	bitCountMap := make(map[byte]int)
	for i := 0; i < count; i++ {
		byteOffset := i / 8
		bitOffset := (uint8)(i % 8)
		bit := (bytes[byteOffset] >> bitOffset) & 0x01
		bitCountMap[bit]++
	}
	var result float64
	for _, v := range bitCountMap {
		frequency := float64(v) / float64(count)
		result -= frequency * (math.Log(frequency) / math.Log(2))
	}
	return result
}
