package model

import "time"

type Receipt struct {
	Retailer     string    `json:"retailer,omitempty"`
	PurchaseTime time.Time `json:"purchase_time,omitempty"`
	Total        float64   `json:"total,omitempty"`
	Items        []Item    `json:"items,omitempty"`
}

type Item struct {
	ShortDescription string  `json:"short_description,omitempty"`
	Price            float64 `json:"price,omitempty"`
}
