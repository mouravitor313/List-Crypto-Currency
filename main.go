package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/mouravitor313/List-Crypto-Currency/internal/config"
	"github.com/mouravitor313/List-Crypto-Currency/internal/models"
	"github.com/mouravitor313/List-Crypto-Currency/internal/proto"
	"github.com/mouravitor313/List-Crypto-Currency/internal/server"
	"github.com/mouravitor313/List-Crypto-Currency/internal/service"
	"google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

var cryptos []models.Crypto

func main() {

    var err error
    var lis net.Listener

    lis, err = net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Fail to listen at: %v", err)
    }

    grpcServer := grpc.NewServer()
    proto.RegisterCryptoServiceServer(grpcServer, &server.CryptoServer{})

    reflection.Register(grpcServer)

    go func() {
        log.Println("Server gRPC running at: 50051")
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("Fail to serve: %v", err)
        }
    }()

    if err := config.LoadAPIKey(); err != nil {
        log.Fatalf("Load API key fail: %v", err)
    }

    config.InitRedis()
    if err := config.InitRedis(); err != nil {
        log.Fatalf("Fail to init Redis: %v", err)
    }

    cryptos, err = service.FetchTopCryptos()
    if err != nil {
        fmt.Println("Erro ao obter criptos:", err)
    } else {
        fmt.Println("Criptomoedas carregadas:", cryptos)
    }

    go server.UpdateCryptosPeriodically()
    go server.BroadcastUpdates()

    http.HandleFunc("/", server.VerifyIfAPIIsOnline)
    http.HandleFunc("/cryptos", server.DisplayCryptos)
    http.HandleFunc("/ws", server.HandleWebSocketConnections)

    fmt.Println("Servidor rodando em http://localhost:8000")
    err = http.ListenAndServe(":8000", nil)
    if err != nil {
        log.Fatalf("Fail to start HTTP server: %v", err)
    }
}
