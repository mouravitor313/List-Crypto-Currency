package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mouravitor313/List-Crypto-Currency/internal/models"
	"github.com/mouravitor313/List-Crypto-Currency/internal/api"
)

func TestVerifyIfAPIIsOnline (t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	VerifyIfAPIIsOnline(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Waiting for HTTP 200, obtained %d", rec.Code)
	}

	expectedBody := "API está online\n"
	if rec.Body.String() != expectedBody {
		t.Fatalf("Incorret response. Waiting: %s. Obtained: %s", expectedBody, rec.Body.String())
	}
}

func TestGetTopCryptos_Success(t *testing.T) {
	defer func() { getTopCryptos = api.GetTopCryptos }()
	getTopCryptos = func() ([]models.Crypto, error) {
		return []models.Crypto{{Name:"X", Symbol:"X", MarketCap:1, CurrentPrice:2}}, nil
	}
}

func TestDisplayCryptos_DefaultCurrencyUSD(t *testing.T) {
	defer func() { getTopCryptos = api.GetTopCryptos } ()
	getTopCryptos = func() ([]models.Crypto, error) {
		return []models.Crypto{
			{Name: "Bitcoin", Symbol: "BTC", MarketCap: 100, CurrentPrice: 10, MarketCapRank: 1},
		}, nil
	}

	req := httptest.NewRequest("GET", "/cryptos", nil)
	rec := httptest.NewRecorder()

	DisplayCryptos(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("waiting: 200. obtained: %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("waiting: JSON. obtained: %s", ct)
	}

	var got []models.Crypto
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("fail to read JSON: %v", err)
	}
	if len(got) != 1 || got[0].Symbol != "BTC" {
		t.Fatalf("incorret response: %v", got)
	}
}

func TestDisplayCryptos_ConvertCurrencyEUR(t *testing.T) {
	defer func() { 
		getTopCryptos = api.GetTopCryptos
		getExchangeRate = api.GetExchangeRate
	} ()
	getTopCryptos = func() ([]models.Crypto, error) {
		return []models.Crypto{
			{Name: "Bitcoin", Symbol: "BTC", MarketCap: 100, CurrentPrice: 10, MarketCapRank: 1},
		}, nil
	}

	getExchangeRate = func(currency string) (float64, error) {
		if currency != "EUR" {
			t.Fatalf("waiting: EUR. obtained: %s", currency)
		}
		return 2, nil
	}

	req := httptest.NewRequest("GET", "/cryptos?currency=EUR", nil)
	rec := httptest.NewRecorder()

	DisplayCryptos(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("waiting: 200. obtained: %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("waiting: JSON. obtained: %s", ct)
	}

	var got []models.Crypto
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("fail to read JSON: %v", err)
	}
	wantPrice := 10 * 2
	if got[0].CurrentPrice != float64(wantPrice) {
		t.Errorf("waiting price: %v. obtained price: %v", wantPrice, got[0].CurrentPrice)
	}
}