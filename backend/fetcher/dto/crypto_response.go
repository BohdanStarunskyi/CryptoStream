package dto

import "fetceher_service/models/crypto"

type CryptoApiResponse struct {
	Id           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	CurrentPrice float64 `json:"current_price"`
	PriceChange  float64 `json:"price_change_24h"`
}

func (c *CryptoApiResponse) ToProto() *crypto.CryptoUpdate {
	return &crypto.CryptoUpdate{
		Id:              c.Id,
		Symbol:          c.Symbol,
		Name:            c.Name,
		Image:           c.Image,
		CurrentPrice:    c.CurrentPrice,
		PriceChange_24H: c.PriceChange,
	}
}
