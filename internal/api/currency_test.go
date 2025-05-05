package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"

	"github.com/mouravitor313/List-Crypto-Currency/internal/config"
)

func TestGetExchangeRate_Success(t *testing.T) {
	config.InitRedis()
	mockResponse := map[string]interface{}{
		"quotes": map[string]float64{
			"USDBRL": 5.10,
		},
	}
	mockBody, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockBody)
	}))
	defer server.Close()

	baseCurrencyAPIURL = server.URL + "?"

	rate, err := GetExchangeRate("BRL")
	if err != nil {
		t.Fatalf("Fail to obtain cryptos information: %v", err)
	}

	if rate != 5.10 {
		t.Fatalf("Expected rate: 5.10, obtained: %f", rate)
	}
}