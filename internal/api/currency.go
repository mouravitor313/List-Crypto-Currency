package api

import (
    "encoding/json"
    "fmt"
    "net/http"
	
    "github.com/mouravitor313/List-Crypto-Currency/internal/config"
)

const baseCurrencyAPIURL = "https://api.currencylayer.com/live?"

func GetExchangeRate(currency string) (float64, error) {
    url := fmt.Sprintf("%saccess_key=%s&currencies=%s", baseCurrencyAPIURL, config.CurrencyLayerAPIKey, currency)

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
