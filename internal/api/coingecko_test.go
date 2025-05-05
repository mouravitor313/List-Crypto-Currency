package api

import (
	"testing"
	"github.com/mouravitor313/List-Crypto-Currency/internal/config"
)

func TestGetTopCryptos_Success(t *testing.T) {
	config.InitRedis()

	cryptos, err := GetTopCryptos()
	if err != nil {
		t.Fatalf("Fail to obtain cryptos information: %v", err)
	}

	if len(cryptos) == 0 {
		t.Fatalf("Nenhuma criptomoeda retornada")
	}
}

func TestGetTopCryptos_CacheUsage(t *testing.T) {
	config.InitRedis()

	GetTopCryptos()
	cryptos, err := GetTopCryptos()
	if err != nil {
		t.Fatalf("Fail to obtain data from cache: %v", err)
	}

	if len(cryptos) == 0 {
		t.Fatalf("Cache returns 0 cryptos")
	}
}