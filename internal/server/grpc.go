package server

import (
    "context"
    "time"

    "github.com/mouravitor313/List-Crypto-Currency/internal/api"
    "github.com/mouravitor313/List-Crypto-Currency/internal/proto"
)

type CryptoServer struct {
    proto.UnimplementedCryptoServiceServer
}

func (s *CryptoServer) GetTopCryptos(ctx context.Context, req *proto.CryptoRequest) (*proto.CryptoResponse, error) {
    cryptos, err := api.GetTopCryptos()
    if err != nil {
        return nil, err
    }

    var response proto.CryptoResponse
	currency := req.Currency
	if currency != "" && currency != "USD" {
		exchangeRate, err := api.GetExchangeRate(currency)
		if err != nil {
			return nil, err
		}

		for _, crypto := range cryptos {
            response.Cryptos = append(response.Cryptos, &proto.Crypto{
                Name:          crypto.Name,
                Symbol:        crypto.Symbol,
                MarketCap:     crypto.MarketCap * exchangeRate,
                CurrentPrice:  crypto.CurrentPrice * exchangeRate,
                MarketCapRank: int32(crypto.MarketCapRank),
            })
        }
    } else {
        for _, crypto := range cryptos {
            response.Cryptos = append(response.Cryptos, &proto.Crypto{
                Name:          crypto.Name,
                Symbol:        crypto.Symbol,
                MarketCap:     crypto.MarketCap,
                CurrentPrice:  crypto.CurrentPrice,
                MarketCapRank: int32(crypto.MarketCapRank),
            })
		}
	}

	return &response, nil

}

func (s *CryptoServer) StreamCryptoUpdates(req *proto.CryptoRequest, stream proto.CryptoService_StreamCryptoUpdateServer) error {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        cryptos, err := api.GetTopCryptos()
        if err != nil {
            return err
        }

        var response proto.CryptoResponse
        for _, crypto := range cryptos {
            response.Cryptos = append(response.Cryptos, &proto.Crypto{
                Name:          crypto.Name,
                Symbol:        crypto.Symbol,
                MarketCap:     crypto.MarketCap,
                CurrentPrice:  crypto.CurrentPrice,
                MarketCapRank: int32(crypto.MarketCapRank),
            })
        }

        if err := stream.Send(&response); err != nil {
            return err
        }
    }

    return nil
}
