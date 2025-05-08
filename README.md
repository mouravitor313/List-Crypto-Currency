# ✨ List Crypto Currency

## Introduction

List Crypto Currency is a Go API that provides, in real-time, the top 10 cryptocurrencies by market capitalization, with support for multiple interfaces (REST, WebSocket, and gRPC) and Redis caching. This API is capable of:

- Accessing the top 10 coins by Market Cap via [CoinGecko API](https://www.coingecko.com/en/api) and [CurrencyLayer API](https://currencylayer.com/documentation)
- Storing and serving cached responses with **Redis**
- Automatically updating every minute
- Supporting multiple consumption interfaces:
  - REST HTTP
  - WebSocket (for real-time streaming)
  - gRPC (for high-performance integrations)

All of this built in **Go 1.24** with the following libraries and stacks:

- [Gorilla WebSocket](https://github.com/gorilla/websocket) – bidirectional streaming
- [go-redis/redis v9](https://github.com/redis/go-redis) – cache and pub/sub
- [joho/godotenv](https://github.com/joho/godotenv) – environment variables
- [gRPC-Go](https://grpc.io/docs/languages/go/) + [Protobuf](https://developers.google.com/protocol-buffers) – high-performance RPC

---

## Documentation

This step-by-step guide will help you clone, configure, and run the **List Crypto Currency** project on macOS, Linux, or Windows (via WSL).

### 1. Prerequisites

- 🐹 **Go `1.24.2`** or higher ([download](https://go.dev/dl/))
- 🧠 **Redis** (local, Docker, or service)
- 🧬 **Git**
- 🧪 **Postman** or similar (for HTTP/WebSocket testing)
- **CoinGecko** API Key ([here](https://www.coingecko.com/en/api))
- **CurrencyLayer** API Key ([here](https://currencylayer.com/documentation))
- Optional: **WSL** on Windows

---

### 2. Go Installation

#### Windows

Install [WSL (Windows Subsystem for Linux)](https://learn.microsoft.com/pt-br/windows/wsl/install) with a Linux distro (recommended: Ubuntu) and follow the steps for Linux execution.

#### macOS

```bash
brew install go@1.24
echo 'export PATH="$PATH:/usr/local/opt/go@1.24/bin"' >> ~/.zshrc
```

#### Linux (Debian/Ubuntu)

```bash
rm -rf /usr/local/go
sudo apt-get update
sudo apt install golang-go
echo 'export PATH="$PATH:/usr/local/go/bin"' >> ~/.bashrc
```

**or if that doesn't work:**

```bash
rm -rf /usr/local/go
sudo apt-get update
sudo apt install golang-go
echo 'export PATH="$PATH:/usr/local/go/bin"' >> ~/.zshrc
```

---

### 3. Cloning the repository

```bash
git clone https://github.com/mouravitor313/List-Crypto-Currency.git
cd List-Crypto-Currency
```

---

### 4. Environment Variables

Create a `.env` file in the root directory:

```env
COINGECKO_API_KEY=your_coin_gecko_api_key      # get your API key at: https://www.coingecko.com/en/api
CURRENCY_LAYER_API=your_currency_layer_api_key  # get your API key at: https://currencylayer.com/documentation
# Optional if you want to configure Redis further, ensuring security:
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
```

> ⚠️ Remember to include your API keys by replacing "your_coin_gecko_api_key" and "your_currency_layer_api_key".

---

### 5. Installing Dependencies

```bash
go mod download
```

Main libs:

- `github.com/gorilla/websocket`
- `github.com/redis/go-redis/v9`
- `github.com/joho/godotenv`

---

### 6. Redis Configuration

#### macOS

```bash
brew install redis
brew services start redis
```

#### Linux

```bash
sudo apt-get install redis-server
sudo systemctl enable --now redis-server
```

#### Docker

```bash
docker run -d --name redis -p 6379:6379 redis:latest
```

---

### 7. Running the Server

```bash
go run main.go
```

Expected output:

```text
💹 Cryptocurrencies loaded: [...]
🌐 Server running on http://localhost:8000
```

---

### 8. HTTP Tests with REST

-   **API Status**

    ```http
    GET http://localhost:8000/
    ```

    Returns: `API is online`

-   **List Cryptocurrencies**

    ```http
    GET http://localhost:8000/cryptos?currency=EUR
    ```

    -   `currency` (optional): ISO code (e.g., `USD`, `BRL`).

Examples with `curl`:

```bash
curl http://localhost:8000/cryptos
curl "http://localhost:8000/cryptos?currency=BRL"
```

---

### 9. WebSocket Tests

-   **Endpoint:** `ws://localhost:8000/ws?currency=USD` (tip): replace `USD` with the ISO code of your preferred currency (e.g., `EUR`, `BRL`)
-   In Postman:
    1.  `New` → `WebSocket Request`
    2.  Enter `ws://localhost:8000/ws?currency=USD`
    3.  `Connect` and observe real-time update JSONs.

---

### 10. gRPC Integration

#### Tool Installation (Linux/macOS)

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```
#### macOS
```bash
brew install protobuf
```
#### Linux
```bash
sudo apt-get install protobuf-compiler
```

#### Run the server:

```bash
go run main.go
```

 #### Open another terminal instance and add protoc to PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

#### Execute in the second instance to generate Bindings
```bash
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  internal/proto/crypto.proto
```

#### Execute in the second instance to test with grpcurl

```bash
grpcurl -plaintext -d '{"currency":"USD"}' localhost:50051 crypto.CryptoService/GetTopCryptos
```
---

### 11. Project Structure

```text
.
├── README.md
├── go.mod
├── go.sum
├── .env
├── internal
│   ├── api
│   ├── config
│   ├── models
│   ├── proto
│   ├── server
│   └── service
└── main.go
```

---

## 🎉 Conclusion

With this setup, you will have a cryptocurrency API in **Go** with:

- Efficient caching in **Redis**
- Automatic updates
- Multiple interfaces (REST, WebSocket, gRPC)
- Ease of integration into any stack

## Contacts:

[Email](dev.vitormoura@gmail.com)

[LinkedIn](https://www.linkedin.com/in/v%C3%ADtor-moura/)

Ready to add **real-time data** to your next project!
