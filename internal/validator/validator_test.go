/*
Package validator containing any validator function
*/
package validator

import (
	"fmt"
	"testing"

	"github.com/reyhanfikridz/ecom-order-service/internal/model"
)

// TestIsOrderInfoValid test IsOrderInfoValid
func TestIsOrderInfoValid(t *testing.T) {
	// initialize testing table
	testTable := []struct {
		TestName       string
		Order          model.Order
		ExpectedResult error
	}{
		{
			TestName: "Test Form Complete",
			Order: model.Order{
				Status:        "in-cart",
				Qty:           2,
				TotalPrice:    2000000.50,
				ProductName:   "Product 1",
				ProductPrice:  1000000.50,
				ProductWeight: 1.5,
			},
			ExpectedResult: nil,
		},
		{
			TestName: "Test Form Incomplete 1",
			Order: model.Order{
				Status:        "",
				Qty:           2,
				TotalPrice:    2000000.50,
				ProductName:   "Product 1",
				ProductPrice:  1000000.50,
				ProductWeight: 1.5,
			},
			ExpectedResult: fmt.Errorf("status empty/not found"),
		},
		{
			TestName: "Test Form Incomplete 2",
			Order: model.Order{
				Status:        "in-cart",
				Qty:           0,
				TotalPrice:    2000000.50,
				ProductName:   "Product 1",
				ProductPrice:  1000000.50,
				ProductWeight: 1.5,
			},
			ExpectedResult: fmt.Errorf("qty empty/not found"),
		},
		{
			TestName: "Test Form Incomplete 3",
			Order: model.Order{
				Status:        "in-cart",
				Qty:           2,
				TotalPrice:    0,
				ProductName:   "Product 1",
				ProductPrice:  1000000.50,
				ProductWeight: 1.5,
			},
			ExpectedResult: fmt.Errorf("total_price empty/not found"),
		},
		{
			TestName: "Test Form Incomplete 4",
			Order: model.Order{
				Status:        "in-cart",
				Qty:           2,
				TotalPrice:    2000000.50,
				ProductName:   "",
				ProductPrice:  1000000.50,
				ProductWeight: 1.5,
			},
			ExpectedResult: fmt.Errorf("product_name empty/not found"),
		},
		{
			TestName: "Test Form Incomplete 5",
			Order: model.Order{
				Status:        "in-cart",
				Qty:           2,
				TotalPrice:    2000000.50,
				ProductName:   "Product 1",
				ProductPrice:  0,
				ProductWeight: 1.5,
			},
			ExpectedResult: fmt.Errorf("product_price empty/not found"),
		},
		{
			TestName: "Test Form Incomplete 6",
			Order: model.Order{
				Status:        "in-cart",
				Qty:           2,
				TotalPrice:    2000000.50,
				ProductName:   "Product 1",
				ProductPrice:  1000000.50,
				ProductWeight: 0,
			},
			ExpectedResult: fmt.Errorf("product_weight empty/not found"),
		},
	}

	// loop test in test table
	for _, test := range testTable {
		err := IsOrderValid(test.Order)
		if test.ExpectedResult == nil && err != nil {
			t.Errorf("Expected order valid, but got invalid => %s", err.Error())
		} else if test.ExpectedResult != nil {
			if err == nil {
				t.Errorf("Expected order invalid, but got valid")
			} else if test.ExpectedResult.Error() != err.Error() {
				t.Errorf("Expected error '" +
					test.ExpectedResult.Error() + "', but got '" + err.Error() + "'")
			}
		}
	}
}
