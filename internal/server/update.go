package server

import (
    "fmt"
    "time"
    "github.com/mouravitor313/List-Crypto-Currency/internal/api"
    "github.com/mouravitor313/List-Crypto-Currency/internal/models"
)

func removeClient(clientList []ClientInfo, clientToRemove ClientInfo) []ClientInfo {
    var newClients []ClientInfo
    for _, c := range clientList {
        if c.Conn != clientToRemove.Conn {
            newClients = append(newClients, c)
        }
    }
    return newClients
}

func UpdateCryptosPeriodically() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            fmt.Println("Atualizando lista de criptos...")
            newCryptos, err := api.GetTopCryptos()
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

func BroadcastUpdates() {
    for {
        data := <-broadcast

        mu.Lock()
        for _, client := range clients {

            exchangeRate, err := api.GetExchangeRate(client.Currency)
            if err != nil {
                fmt.Println("Erro ao obter taxa de câmbio:", err)
                exchangeRate = 1.0
            }

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

            err = client.Conn.WriteJSON(convertedCryptos)
            if err != nil {
                fmt.Println("Erro ao enviar atualização:", err)
                client.Conn.Close()
                clients = removeClient(clients, client)
            }
        }
        mu.Unlock()
    }
}
