package model

type BioPlanetOrder struct {
	Address      Address `json:"Address"`
	PaymentId    int     `json:"PaymentId"`
	DeliveryId   int     `json:"DeliveryId"`
	DeliveryName string  `json:"DeliveryName"`
	Comment      string  `json:"Comment"`
	OrderLines   struct {
		KeyType string `json:"KeyType"`
		Lines   []struct {
			Key      string `json:"Key"`
			Quantity int    `json:"Quantity"`
		} `json:"Lines"`
	} `json:"OrderLines"`
	InpostPaczkomatCode string `json:"InpostPaczkomatCode"`
}

type Address struct {
	Name            string `json:"Name"`
	Street          string `json:"Street"`
	City            string `json:"City"`
	PostalCode      string `json:"PostalCode"`
	Phone           string `json:"Phone"`
	CountryId       int    `json:"CountryId"`
	RegionId        int    `json:"RegionId"`
	Email           string `json:"Email"`
	ApartmentNumber string `json:"ApartmentNumber"`
	HouseNumber     string `json:"HouseNumber"`
	TaxNumber       string `json:"TaxNumber"`
	OneTimeAdress   bool   `json:"OneTimeAdress"`
}
