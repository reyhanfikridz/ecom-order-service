/*
Package validator containing any validator function
*/
package validator

import (
	"fmt"
	"strings"

	"github.com/reyhanfikridz/ecom-order-service/internal/model"
)

// IsOrderValid check if order data is valid
//
// return error nil if it's valid
func IsOrderValid(o model.Order) error {
	if strings.TrimSpace(o.Status) == "" {
		return fmt.Errorf("status empty/not found")
	}

	if o.Qty == 0 {
		return fmt.Errorf("qty empty/not found")
	}

	if o.TotalPrice == 0 {
		return fmt.Errorf("total_price empty/not found")
	}

	if strings.TrimSpace(o.ProductName) == "" {
		return fmt.Errorf("product_name empty/not found")
	}

	if o.ProductPrice == 0 {
		return fmt.Errorf("product_price empty/not found")
	}

	if o.ProductWeight == 0 {
		return fmt.Errorf("product_weight empty/not found")
	}

	return nil
}
