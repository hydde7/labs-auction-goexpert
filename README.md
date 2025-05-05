
# FullCycle Auction Service

Este serviço em Go implementa um sistema de leilões com **fechamento automático**, após um tempo configurado via variável de ambiente.

---

## Funcionalidades

- Criar leilões com dados de produto
- Fechamento automático após tempo definido
- Lances válidos apenas enquanto o leilão estiver **aberto**
- Banco de dados: MongoDB
- Observabilidade via logs (Logrus)

---

## Variáveis de Ambiente

Crie um arquivo `.env` em `cmd/auction/.env` com as variáveis:

```env
# Tempo de duração do leilão (ex: 30s, 2m, 1h)
AUCTION_INTERVAL=30s

# MongoDB
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

---

## Como rodar com Docker

1. Clone este repositório:

```bash
git clone https://github.com/devfullcycle/labs-auction-goexpert
cd labs-auction-goexpert
```

2. Ajuste o arquivo `.env` conforme mostrado acima.

3. Rode o projeto com Docker Compose:

```bash
docker compose up --build -d
```

4. Acesse a aplicação em `http://localhost:8080`

---

## Executar os testes

Para rodar os testes automatizados (inclusive do fechamento automático):

```bash
docker compose exec app go test ./... -v
```

Ou teste diretamente o pacote de leilão:

```bash
docker compose exec app go test ./internal/infra/database/auction -v
```

---

## Leilão com fechamento automático

O fechamento automático ocorre com base no tempo configurado pela variável `AUCTION_INTERVAL`.

Cada leilão criado dispara uma **goroutine** que:
- Espera o tempo configurado
- Marca o leilão como `Completed` no banco

Você pode verificar isso rodando o teste incluso:

```bash
docker compose exec app go test ./internal/infra/database/auction -v
```

---

## Finalizar

Para encerrar tudo:

```bash
docker-compose down -v
```
