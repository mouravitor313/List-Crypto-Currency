package api

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
	
    "github.com/mouravitor313/List-Crypto-Currency/internal/config"
    "github.com/mouravitor313/List-Crypto-Currency/internal/models"
)

const baseCoinGeckoAPIURL = "https://api.coingecko.com/api/v3/"

func GetTopCryptos() ([]models.Crypto, error) {
    url := fmt.Sprintf("%scoins/markets?vs_currency=usd&order=market_cap_desc&per_page=10&page=1", baseCoinGeckoAPIURL)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("x-cg-demo-api-key", config.CoinGeckoAPIKey)

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

    var cryptos []models.Crypto
    err = json.Unmarshal(body, &cryptos)
    if err != nil {
        return nil, err
    }

    return cryptos, nil
}
