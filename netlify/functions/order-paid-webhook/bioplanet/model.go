package bioplanet

type Order struct {
	Address             Address    `json:"Address"`
	PaymentId           int        `json:"PaymentId"`
	DeliveryId          int        `json:"DeliveryId,omitempty"`
	DeliveryName        string     `json:"DeliveryName"`
	Comment             string     `json:"Comment"`
	OrderLines          OrderLines `json:"OrderLines"`
	InpostPaczkomatCode string     `json:"InpostPaczkomatCode,omitempty"`
}

type Address struct {
	Name            string `json:"Name"`
	Street          string `json:"Street"`
	City            string `json:"City"`
	PostalCode      string `json:"PostalCode"`
	Phone           string `json:"Phone"`
	CountryId       int    `json:"CountryId,omitempty"`
	RegionId        int    `json:"RegionId,omitempty"`
	Email           string `json:"Email"`
	ApartmentNumber string `json:"ApartmentNumber,omitempty"`
	HouseNumber     string `json:"HouseNumber,omitempty"`
	TaxNumber       string `json:"TaxNumber,omitempty"`
	OneTimeAdress   bool   `json:"OneTimeAdress"`
}

type OrderLines struct {
	KeyType string `json:"KeyType"`
	Lines   []Line `json:"Lines"`
}

type Line struct {
	Key      string `json:"Key"`
	Quantity int    `json:"Quantity"`
}

type ApiTokenPost struct {
	Hash      string `json:"Hash"`
	ClientId  int    `json:"ClientId"`
	Timestamp string `json:"Timestamp"`
}

type ApiToken struct {
	AccessToken string `json:"AccessToken"`
	TokenType   string `json:"TokenType"`
	ExpiresIn   int    `json:"ExpiresIn"`
}

type OrderConfirmation struct {
	OrderId int `json:"OrderId"`
}
