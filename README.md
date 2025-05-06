# List Crypto Currency

Este documento descreve como configurar e executar o projeto List Crypto Currency, desde a instalação do Go em diferentes sistemas operacionais, configuração de dependências, uso do Redis como cache, até acesso às rotas HTTP e conexão WebSocket pelo Postman.

---

## 1. Pré-requisitos

*   Sistema operacional: macOS, Linux ou Windows
*   Go versão `1.24.2`
*   Redis (você pode usar localmente ou através de container)
*   `git`
*   Postman (para testes de HTTP e WebSocket)

---

## 2. Instalação do Go

### macOS

```bash
brew install go@1.24
```

> Adicione ao `PATH` (no `.zshrc` ou `.bash_profile`):
> ```bash
> export PATH="$PATH:/usr/local/opt/go@1.24/bin"
> ```

### Linux (Debian/Ubuntu)

```bash
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz
```

> No `~/.bashrc`:
> ```bash
> export PATH="$PATH:/usr/local/go/bin"
> ```

### Windows

1.  Baixe o instalador MSI em [https://go.dev/dl/](https://go.dev/dl/)
2.  Execute o instalador e siga as instruções
3.  Verifique no PowerShell:

```powershell
go version
```

---

## 3. Clonar o repositório

```bash
git clone https://github.com/mouravitor313/List-Crypto-Currency.git
cd List-Crypto-Currency
```

---

## 4. Configuração das variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto com as chaves de API:

```dotenv
COINGECKO_API_KEY=your_coin_gecko_api_key
CURRENCY_LAYER_API=your_currency_layer_api_key
REDIS_ADDR=localhost:6379         # opcional, caso seu init leu esta chave
REDIS_PASSWORD=                   # opcional
```

O código em `internal/config/env.go` faz o carregamento dessas variáveis usando `github.com/joho/godotenv`:

```go
err := godotenv.Load()
CoinGeckoAPIKey = os.Getenv("COINGECKO_API_KEY")
CurrencyLayerAPIKey = os.Getenv("CURRENCY_LAYER_API")
```

> **Atenção**: Certifique-se de que o `.env` está no `.gitignore` para não vazar suas chaves.

---

## 5. Instalação de dependências Go

No diretório raiz do projeto, execute:

```bash
go mod download
```

Isso instalará:

*   `github.com/gorilla/websocket`
*   `github.com/joho/godotenv`
*   `github.com/redis/go-redis/v9`

---

## 6. Configuração e inicialização do Redis

### Opção 1: Instalar localmente

#### macOS (Homebrew)

```bash
brew install redis
brew services start redis
```

#### Linux (Debian/Ubuntu)

```bash
sudo apt-get install redis-server
sudo systemctl enable redis-server
sudo systemctl start redis-server
```

### Opção 2: Usar Docker

```bash
docker run -d --name redis -p 6379:6379 redis:latest
```

No código, a inicialização do Redis é feita em `config.InitRedis()` (em `internal/config/config.go`), que deve usar as variáveis:

(O ideal)
```go
client := redis.NewClient(&redis.Options{
    Addr:     os.Getenv("REDIS_ADDR"),
    Password: os.Getenv("REDIS_PASSWORD"),
    DB:       0,
})
err := client.Ping(context.Background()).Err()
```

(Presente no meu código / Setup mais rápido para rodar localmente)
```go
RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := RedisClient.Ping(context.Background()).Result()
```

---

## 7. Rodando a API

No terminal:

```bash
go run main.go
```

**Saída esperada:**

```text
Criptomoedas carregadas: [{...}]
Servidor rodando em http://localhost:8000
```

---

## 8. Rotas HTTP

### 8.1. Verificar status da API

**Método:** `GET`

**URL:** <http://localhost:8000/>

**Descrição:** Retorna um texto simples indicando que a API está online.

**Resposta esperada:**

```text
API está online
```

### 8.2. Listar criptomoedas

**Método:** `GET`

**URL padrão:** <http://localhost:8000/cryptos>

**Parâmetros de query:**

*   `currency` (opcional): código da moeda destino (ex: `EUR`, `BRL`). Default = `USD`.

#### 8.2.1. Exemplo sem currency

```bash
curl http://localhost:8000/cryptos
```

**Resposta JSON:**

```json
[
    {
        "name": "Bitcoin",
        "symbol": "BTC",
        "market_cap": 500000000000,
        "current_price": 27000,
        "market_cap_rank": 1
    },
    ...
]
```

#### 8.2.2. Exemplo com currency=EUR

```bash
curl "http://localhost:8000/cryptos?currency=EUR"
```

**Resposta JSON:** valores convertidos para Euro.

---

## 9. Conexão WebSocket

**Endpoint:** `ws://localhost:8000/ws`

**Descrição:** Recebe atualizações periódicas de preços.

**Como testar no Postman:**

1.  Abra o Postman.
2.  Clique em `New` > `WebSocket Request`.
3.  Insira a URL `ws://localhost:8000/ws`.
4.  Clique em `Connect`.
5.  Você verá mensagens `JSON` de atualização de criptomoedas.

---

## 10. Estrutura de Pastas
.
├── README.md
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── coingecko.go
│   │   ├── coingecko_test.go
│   │   ├── currency.go
│   │   └── currency_test.go
│   ├── config
│   │   ├── config.go
│   │   └── env.go
│   ├── models
│   │   └── crypto.go
│   ├── proto
│   │   └── crypto.proto
│   ├── server
│   │   ├── routes.go
│   │   ├── routes_test.go
│   │   ├── update.go
│   │   └── websocket.go
│   └── service
│       └── crypto.go
└── main.go
---
