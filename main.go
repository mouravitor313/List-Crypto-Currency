package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type Crypto struct {
	Name          string  `json:"name"`
	Symbol        string  `json:"symbol"`
	MarketCap     float64 `json:"market_cap"`
	CurrentPrice  float64 `json:"current_price"`
	MarketCapRank int     `json:"market_cap_rank"`
}

const baseCoinGeckoAPIURL = "https://api.coingecko.com/api/v3/"
const baseCurrencyAPIURL = "https://api.currencylayer.com/live?"

var (
	coingeckoAPIKey string
	currencylayerAPIKey string
	cryptos []Crypto
	clients = make(map[*websocket.Conn]bool)
	broadcast = make(chan []Crypto)
	mu sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func loadAPIKey() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar .env:", err)
		return
	}

	coingeckoAPIKey = os.Getenv("COINGECKO_API_KEY")
	if coingeckoAPIKey == "" {
		fmt.Println("API_KEY não encontrada no .env!")
		os.Exit(1)
	}

	currencylayerAPIKey = os.Getenv("CURRENCY_LAYER_API")
	if currencylayerAPIKey == "" {
		fmt.Println("API_KEY não encontrada no .env!")
		os.Exit(1)
	}
}

func getCoinGeckoPing() (string, error) {
	ping_rote := "ping"
	req, err := http.NewRequest("GET", baseCoinGeckoAPIURL+ping_rote, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("x-cg-demo-api-key", coingeckoAPIKey)

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
	req, err := http.NewRequest("GET", baseCoinGeckoAPIURL+marketRote, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-cg-demo-api-key", coingeckoAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
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

func updateCryptosPeriodically() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Atualizando lista de criptos...")
			newCryptos, err := getTopCryptos()
			if err == nil {
				mu.Lock()
				cryptos = newCryptos
				mu.Unlock()
				broadcast <- cryptos
			} else {
				fmt.Println("Erro ao atualizar criptomoedas:", err)
			}

		}
	}
}

func handleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao conectar Websocket:", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	for {
		select {
		case data := <- broadcast:
			conn.WriteJSON(data)
		}
	}
}

func broadcastUpdates() {
	for {
		data := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(data)
			if err != nil {
				fmt.Println("Erro ao enviar atualização:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func verifyIfAPIIsOnline(w http.ResponseWriter, r *http.Request) {
	coinGeckoResponse, err := getCoinGeckoPing()
	if err != nil {
		http.Error(w, "Erro ao obter resposta da API", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, coinGeckoResponse)
}

func getExchangeRate(currency string) (float64, error) {
	url := fmt.Sprintf("%saccess_key=%s&currencies=%s", baseCurrencyAPIURL, currencylayerAPIKey, currency)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	quotes, ok := result["quotes"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("erro ao obter taxa de câmbio")
	}

	rateKey := "USD" + currency
	rate, ok := quotes[rateKey].(float64)
	if !ok {
		return 0, fmt.Errorf("moeda não encontrada")
	}

	return rate, nil
}

func displayCryptos(w http.ResponseWriter, r *http.Request) {
	currency := r.URL.Query().Get("currency")
	if currency == "" || currency == "USD" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cryptos)
		return
	}

	exchangeRate, err := getExchangeRate(currency)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter taxa de câmbio para %s", currency), http.StatusInternalServerError)
		return
	}

	var convertedCryptos []Crypto
	for _, crypto := range cryptos {
		convertedCryptos = append(convertedCryptos, Crypto{
			Name: crypto.Name,
			Symbol: crypto.Symbol,
			MarketCap: crypto.MarketCap * exchangeRate,
			CurrentPrice: crypto.CurrentPrice * exchangeRate,
			MarketCapRank: crypto.MarketCapRank,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convertedCryptos)
}

func main() {
	loadAPIKey()
	initialCryptos, err := getTopCryptos()
	if err == nil {
		cryptos = initialCryptos
	} else {
		fmt.Println("Erro ao obter criptos:", err)
	}

	go updateCryptosPeriodically()
	go broadcastUpdates()

	http.HandleFunc("/", verifyIfAPIIsOnline)
	http.HandleFunc("/cryptos", displayCryptos)
	http.HandleFunc("/ws", handleWebSocketConnections)

	fmt.Println("Servidor rodando em http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
