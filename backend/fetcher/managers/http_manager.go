package managers

import (
	"fmt"
	"io"
	"net/http"
)

type NetworkManager struct {
	Client *http.Client
}

func (nm NetworkManager) MakeRequest() ([]byte, error) {
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd", nil)
	if err != nil {
		return nil, err
	}

	resp, err := nm.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	return body, err
}
