package server

import (
    "encoding/json"
    "fmt"
    "net/http"
	
    "github.com/mouravitor313/List-Crypto-Currency/internal/api"
    "github.com/mouravitor313/List-Crypto-Currency/internal/models"
)

var cryptos []models.Crypto

func VerifyIfAPIIsOnline(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "API está online")
}

func DisplayCryptos(w http.ResponseWriter, r *http.Request) {
    var err error
    cryptos, err = api.GetTopCryptos()
    if err != nil {
        http.Error(w, "Erro ao obter criptos", http.StatusInternalServerError)
        return
    }

    currency := r.URL.Query().Get("currency")
    if currency == "" || currency == "USD" {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(cryptos)
        return
    }

    exchangeRate, err := api.GetExchangeRate(currency)
    if err != nil {
        http.Error(w, "Erro ao obter taxa de câmbio", http.StatusInternalServerError)
        return
    }

    var convertedCryptos []models.Crypto
    for _, crypto := range cryptos {
        convertedCryptos = append(convertedCryptos, models.Crypto{
            Name:          crypto.Name,
            Symbol:        crypto.Symbol,
            MarketCap:     crypto.MarketCap * exchangeRate,
            CurrentPrice:  crypto.CurrentPrice * exchangeRate,
            MarketCapRank: crypto.MarketCapRank,
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(convertedCryptos)
}
