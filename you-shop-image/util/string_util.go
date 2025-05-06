package util

import (
	"math/rand"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)

}

// Check if arr2 contains any element of arr1
func ContainsAny(arr1, arr2 []string) bool {
	elements := make(map[string]bool)
	for _, arr2Element := range arr2 {
		elements[strings.Trim(arr2Element, " ")] = true
	}
	for _, arr1Element := range arr1 {
		if elements[strings.Trim(arr1Element, " ")] {
			return true
		}
	}
	return false
}

// Check if arr2 contains all elements of arr1
func ContainsAll(arr1, arr2 []string) bool {
	elements := make(map[string]bool)
	for _, arr2Element := range arr2 {
		elements[strings.Trim(arr2Element, " ")] = true
	}
	for _, arr1Element := range arr1 {
		if !elements[strings.Trim(arr1Element, " ")] {
			return false
		}
	}
	return true
}
