/*
Package utils containing utilities function

This package cannot have import from another package except for config package
*/
package utils

import (
	"math/rand"
	"time"
)

// GetRandomOrderNumber get random order number
func GetRandomOrderNumber() string {
	// change rand seed so the result is different everytime program running
	rand.Seed(time.Now().UnixNano())

	// set order number length
	n := 15

	// get random order number
	base := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	runeOfOrderNumber := make([]rune, n)
	for i := range runeOfOrderNumber {
		runeOfOrderNumber[i] = base[rand.Intn(len(base))]
	}

	return string(runeOfOrderNumber)
}
