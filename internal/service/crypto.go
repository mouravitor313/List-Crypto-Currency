package service

import (
	"github.com/mouravitor313/List-Crypto-Currency/internal/api"
	"github.com/mouravitor313/List-Crypto-Currency/internal/models"
)

func FetchTopCryptos() ([]models.Crypto, error) {
	return api.GetTopCryptos()
}