package models

// Order represents an order
type Order struct {
	// ItemQty is the number of items the user wants to order
	ItemQty int `json:"itemQty"`
}
