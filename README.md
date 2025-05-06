# ✨ List Crypto Currency

## Introdução

Imagine ter ao alcance de um clique os **10 ativos digitais** com maior valor de mercado, sempre atualizados em tempo real e prontos para integrar seu dashboard ou aplicação.
O **List Crypto Currency** nasceu desse desafio: explorar e aprofundar meus conhecimentos em **Go** enquanto entrego uma ferramenta robusta de consulta e distribuição de dados de criptomoedas. Com ele, você:

- Acessa as top 10 moedas por Market Cap via [CoinGecko API](https://www.coingecko.com/en/api) e [CurrencyLayer API](https://currencylayer.com/documentation)
- Armazena e serve respostas em cache com **Redis**
- Atualiza automaticamente a cada minuto
- Suporta múltiplas interfaces de consumo:
  - REST HTTP
  - WebSocket (para streaming em tempo real)
  - gRPC (para integrações performáticas)

Tudo isso construído em **Go 1.24** com as seguintes bibliotecas e stacks:

- [Gorilla WebSocket](https://github.com/gorilla/websocket) – streaming bidirecional
- [go-redis/redis v9](https://github.com/redis/go-redis) – cache e pub/sub
- [joho/godotenv](https://github.com/joho/godotenv) – variáveis de ambiente
- [gRPC-Go](https://grpc.io/docs/languages/go/) + [Protobuf](https://developers.google.com/protocol-buffers) – RPC de alta performance

---

## Documentação

Este guia passo a passo vai ajudá-lo a clonar, configurar e executar o projeto **List Crypto Currency** em macOS, Linux ou Windows.

### 1. Pré-requisitos

- 🐹 **Go `1.24.2`** ou superior ([download](https://go.dev/dl/))
- 🧠 **Redis** (local, Docker ou serviço)
- 🧬 **Git**
- 🧪 **Postman** ou similar (para testes HTTP/WebSocket)
- Optional: **WSL** no Windows

---

### 2. Instalação do Go

#### macOS

```bash
brew install go@1.24
echo 'export PATH="$PATH:/usr/local/opt/go@1.24/bin"' >> ~/.zshrc
```

#### Linux (Debian/Ubuntu)

```bash
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz
echo 'export PATH="$PATH:/usr/local/go/bin"' >> ~/.bashrc
```

#### Windows

Instalar o [WSL (Windows Subsystem for Linux)](https://learn.microsoft.com/pt-br/windows/wsl/install) com uma distro Linux (ex: Ubuntu) e seguir o passo a passo para execução no Linux.

---

### 3. Clonagem do repositório

```bash
git clone https://github.com/mouravitor313/List-Crypto-Currency.git
cd List-Crypto-Currency
```

---

### 4. Variáveis de Ambiente

Crie um arquivo `.env` na raiz:

```env
COINGECKO_API_KEY=your_coin_gecko_api_key      # https://www.coingecko.com/en/api
CURRENCY_LAYER_API=your_currency_layer_api_key  # https://currencylayer.com/documentation
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
```

> ⚠️ Lembre-se de incluir `.env` no `.gitignore`.

---

### 5. Instalação das Dependências

```bash
go mod download
```

Principais libs:

- `github.com/gorilla/websocket`
- `github.com/redis/go-redis/v9`
- `github.com/joho/godotenv`

---

### 6. Configuração do Redis

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

### 7. Executar o Servidor

```bash
go run main.go
```

Saída esperada:

```text
💹 Criptomoedas carregadas: [...]
🌐 Servidor rodando em http://localhost:8000
```

---

### 8. Testes HTTP com REST

-   **Status da API**

    ```http
    GET http://localhost:8000/
    ```

    Retorna: `API está online`

-   **Listar Criptomoedas**

    ```http
    GET http://localhost:8000/cryptos?currency=EUR
    ```

    -   `currency` (opcional): código ISO (ex: `USD`, `BRL`).

Exemplos com `curl`:

```bash
curl http://localhost:8000/cryptos
curl "http://localhost:8000/cryptos?currency=BRL"
```

---

### 9. Testes via WebSocket

-   **Endpoint:** `ws://localhost:8000/ws`
-   No Postman:
    1.  `New` → `WebSocket Request`
    2.  Insira `ws://localhost:8000/ws`
    3.  `Connect` e observe JSONs de atualização em tempo real.

---

### 10. Integração gRPC

#### Instalação de Ferramentas (Linux/macOS)

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
# macOS
brew install protobuf
# Linux
sudo apt-get install protobuf-compiler
```

#### Gerar Bindings

```bash
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  internal/proto/crypto.proto
```

#### Teste com grpcurl

```bash
grpcurl -plaintext -d '{"currency":"USD"}' localhost:50051 crypto.CryptoService/GetTopCryptos
```
---

### 11. Estrutura do Projeto

```text
.
├── README.md
├── go.mod
├── go.sum
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

## 🎉 Conclusão

Com este setup, você terá uma API de criptomoedas em **Go** com:

- Cache eficiente em **Redis**
- Atualizações automáticas
- Múltiplas interfaces (REST, WebSocket, gRPC)
- Facilidade de integração em qualquer stack

Pronto para agregar **dados em tempo real** ao seu próximo projeto!
