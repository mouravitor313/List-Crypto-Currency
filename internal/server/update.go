package server

import (
    "fmt"
    "time"
    "github.com/mouravitor313/List-Crypto-Currency/internal/api"
)


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
