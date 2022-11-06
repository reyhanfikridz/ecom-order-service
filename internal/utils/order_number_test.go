/*
Package utils containing utilities function

This package cannot have import from another package except for config package
*/
package utils

import (
	"testing"
)

// TestGetRandomOrderNumber test GetRandomOrderNumber
func TestGetRandomOrderNumber(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orderNumber := GetRandomOrderNumber()
		if len(orderNumber) != 15 {
			t.Errorf("Expected orderNumber length 15, but got %d", len(orderNumber))
		}
	}
}
