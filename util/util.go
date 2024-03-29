package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	if toRemove == -1 {
		return arr
	}
	return append(arr[:toRemove], arr[toRemove+1:]...)
}

func Digest(object interface{}) (string, error) {
	msg, err := json.Marshal(object)

	if err != nil {
		return "", err
	}

	return Hash(msg), nil
}

func Hash(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}