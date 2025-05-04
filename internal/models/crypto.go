package models

type Crypto struct {
    Name          string  `json:"name"`
    Symbol        string  `json:"symbol"`
    MarketCap     float64 `json:"market_cap"`
    CurrentPrice  float64 `json:"current_price"`
    MarketCapRank int     `json:"market_cap_rank"`
}