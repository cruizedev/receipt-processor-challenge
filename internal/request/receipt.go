package request

type Receipt struct {
	Retailer     string `json:"retailer,omitempty"`
	PurchaseDate string `json:"purchaseDate,omitempty"`
	PurchaseTime string `json:"purchaseTime,omitempty"`
	Total        string `json:"total,omitempty"`
	Items        []Item `json:"items,omitempty"`
}

type Item struct {
	ShortDescription string `json:"shortDescription,omitempty"`
	Price            string `json:"price,omitempty"`
}
