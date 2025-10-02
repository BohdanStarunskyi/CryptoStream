package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"fetceher_service/dto"
	"fetceher_service/managers"
	"fetceher_service/models/crypto"
)

func main() {
	gatewayHost := os.Getenv("GATEWAY_HOST")
	if gatewayHost == "" {
		gatewayHost = "localhost"
	}
	gatewayPort := os.Getenv("GATEWAY_PORT")
	if gatewayPort == "" {
		gatewayPort = "50051"
	}

	gatewayAddress := fmt.Sprintf("%s:%s", gatewayHost, gatewayPort)
	connector, err := managers.NewGRPCConnector(gatewayAddress)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer connector.Close()

	httpClient := managers.NetworkManager{
		Client: &http.Client{},
	}

	cryptosChan := make(chan []dto.CryptoApiResponse)

	go startFetching(httpClient, cryptosChan, 10*time.Second)

	processUpdates(connector, cryptosChan)
}

func startFetching(client managers.NetworkManager, ch chan<- []dto.CryptoApiResponse, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		data, err := fetchData(client)
		if err != nil {
			fmt.Println("Error fetching data:", err)
		} else {
			ch <- data
		}

		<-ticker.C
	}
}

func fetchData(client managers.NetworkManager) ([]dto.CryptoApiResponse, error) {
	respBytes, err := client.MakeRequest()
	if err != nil {
		return nil, err
	}

	var cryptos []dto.CryptoApiResponse
	if err := json.Unmarshal(respBytes, &cryptos); err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	return cryptos, nil
}

func processUpdates(connector *managers.GRPCConnector, ch <-chan []dto.CryptoApiResponse) {
	for cryptos := range ch {
		list := []*crypto.CryptoUpdate{}
		for _, c := range cryptos {
			list = append(list, c.ToProto())
		}
		connector.SendUpdates(list)
		fmt.Println("Sent updates to server")
	}
}
