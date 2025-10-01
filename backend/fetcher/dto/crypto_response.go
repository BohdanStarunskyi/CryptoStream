package dto

import "fetceher_service/models/crypro"

type CryptoApiResponse struct {
	Id           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	CurrentPrice float64 `json:"current_price"`
}

func (c *CryptoApiResponse) ToProto() *crypro.CryptoUpdate {
	return &crypro.CryptoUpdate{
		Id:           c.Id,
		Symbol:       c.Symbol,
		Name:         c.Name,
		Image:        c.Image,
		CurrentPrice: c.CurrentPrice,
	}
}
