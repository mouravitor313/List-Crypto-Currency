package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Crypto struct {
	Name          string  `json:"name"`
	Symbol        string  `json:"symbol"`
	MarketCap     float64 `json:"market_cap"`
	CurrentPrice  float64 `json:"current_price"`
	MarketCapRank int     `json:"market_cap_rank"`
}

const baseAPIURL = "https://api.coingecko.com/api/v3/"

var apiKey string

func loadAPIKey() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar .env:", err)
		return
	}

	apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("API_KEY não encontrada no .env!")
		os.Exit(1)
	}
}

func getCoinGeckoPing() (string, error) {
	ping_rote := "ping"
	req, err := http.NewRequest("GET", baseAPIURL+ping_rote, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("x-cg-demo-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getTopCryptos() ([]Crypto, error) {
	marketRote := "coins/markets?vs_currency=usd&order=market_cap_desc&per_page=10&page=1"
	req, err := http.NewRequest("GET", baseAPIURL+marketRote, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-cg-demo-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("status_code:",resp.StatusCode)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("API response:", string(body))

	var cryptos []Crypto
	err = json.Unmarshal(body, &cryptos)
	if err != nil {
		fmt.Println("Error Unmarshall", err)
		return nil, err
	}

	return cryptos, nil
}

func verifyIfAPIIsOnline(w http.ResponseWriter, r *http.Request) {
	coinGeckoResponse, err := getCoinGeckoPing()
	if err != nil {
		http.Error(w, "Erro ao obter resposta da API", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, coinGeckoResponse)
}

func displayCryptos(w http.ResponseWriter, r *http.Request) {
	cryptos, err := getTopCryptos()
	if err != nil {
		fmt.Println("Erro:", err)
		http.Error(w, "Erro ao obter criptomoedas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cryptos)
}

func main() {
	loadAPIKey()
	http.HandleFunc("/", verifyIfAPIIsOnline)
	http.HandleFunc("/cryptos", displayCryptos)
	fmt.Println("Servidor rodando em http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
