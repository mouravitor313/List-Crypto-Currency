package main

import (
	"fmt"
	"net/http"

	"github.com/mouravitor313/List-Crypto-Currency/internal/api"
	"github.com/mouravitor313/List-Crypto-Currency/internal/config"
	"github.com/mouravitor313/List-Crypto-Currency/internal/models"
	"github.com/mouravitor313/List-Crypto-Currency/internal/server"
)

var cryptos []models.Crypto

func main() {
    config.LoadAPIKey()

    var err error
    cryptos, err = api.GetTopCryptos()
    if err != nil {
        fmt.Println("Erro ao obter criptos:", err)
    } else {
        fmt.Println("Criptomoedas carregadas:", cryptos) // Debug
    }

    http.HandleFunc("/", server.VerifyIfAPIIsOnline)
    http.HandleFunc("/cryptos", server.DisplayCryptos)
    http.HandleFunc("/ws", server.HandleWebSocketConnections)

    fmt.Println("Servidor rodando em http://localhost:8000")
    http.ListenAndServe(":8000", nil)
}
