package util

import (
	"crypto/sha256"
	"math/rand"
)

func StringInArray(target string, arr []string) bool {
	for _, element := range arr {
		if element == target {
			return true
		}
	}
	return false
}

func HashBytes(bytes []byte) []byte {
	h := sha256.New()
	h.Write(bytes)
	return h.Sum(nil)
}

/* Generate a unique sample of a given size from a given array */
func UniqueRandomSample(arr []string, size int) []string {
	if size >= len(arr) {
		return arr
	}

	generated := make(map[int]bool)
	subset := make([]string, size)
	arrLen := len(arr)
	for i := 0; i < size; i++ {
		for {
			randIndex := rand.Intn(arrLen)
			if !generated[randIndex] {
				generated[randIndex] = true
				subset[i] = arr[randIndex]
				break
			}
		}
	}
	return subset
}

func RemoveOneFromArr(arr []string, toRemove int) []string {
	return append(arr[:toRemove], arr[toRemove+1:]...)
}
