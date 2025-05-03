package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const apiURL = "https://api.coingecko.com/api/v3/ping"
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
	req, err := http.NewRequest("GET", apiURL, nil)
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

func verifyIfAPIIsOnline(w http.ResponseWriter, r *http.Request) {
	coinGeckoResponse, err := getCoinGeckoPing()
	if err != nil {
		http.Error(w, "Erro ao obter resposta da API", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, coinGeckoResponse)
}

func main() {
	loadAPIKey()
	http.HandleFunc("/", verifyIfAPIIsOnline)
	fmt.Println("Servidor rodando em http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}