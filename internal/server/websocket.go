package server

import (
    "fmt"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
    "github.com/mouravitor313/List-Crypto-Currency/internal/api"
    "github.com/mouravitor313/List-Crypto-Currency/internal/models"
)

type ClientInfo struct {
    Conn *websocket.Conn
    Currency string
}

var (
    clients   []ClientInfo
    broadcast = make(chan []models.Crypto)
    mu        sync.Mutex
    upgrader  = websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool { return true },
    }
)

func HandleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Erro ao conectar Websocket:", err)
        return
    }
    defer conn.Close()

    currency := r.URL.Query().Get("currency")
    if currency == "" {
        currency = "USD"
    }

    exchangeRate, err := api.GetExchangeRate(currency)
    if err != nil {
        fmt.Println("Erro ao obter taxa de câmbio:", err)
        exchangeRate = 1.0
    }

    mu.Lock()
    clients = append(clients, ClientInfo{Conn: conn, Currency: currency})
    mu.Unlock()

    for {
        select {
        case data := <-broadcast:
            var convertedCryptos []models.Crypto
            for _, crypto := range data {
                convertedCryptos = append(convertedCryptos, models.Crypto{
                    Name:          crypto.Name,
                    Symbol:        crypto.Symbol,
                    MarketCap:     crypto.MarketCap * exchangeRate,
                    CurrentPrice:  crypto.CurrentPrice * exchangeRate,
                    MarketCapRank: crypto.MarketCapRank,
                })
            }
            conn.WriteJSON(convertedCryptos)
        }
    }
}