# 🏷️ FullCycle Auction — Fechamento Automático de Leilões (Goroutines + MongoDB)

[![Go](https://img.shields.io/badge/Go-1.24-blue)](https://golang.org/)  
[![Docker](https://img.shields.io/badge/Docker-Enabled-blue)](https://www.docker.com/)  
[![MongoDB](https://img.shields.io/badge/Database-MongoDB-green)](https://www.mongodb.com/)  
[![Tests](https://img.shields.io/badge/Tests-Unit%20%26%20Integration-orange)]()

---

Sistema de leilões desenvolvido em **Go**, com a funcionalidade principal de **fechamento automático de leilões** após um tempo configurável. O fechamento é feito de forma **assíncrona** através de **goroutines**, com persistência em **MongoDB**.

## 📌 Objetivo do Projeto

Adicionar à base já existente uma rotina que:
- calcula o tempo de duração do leilão a partir de variáveis de ambiente;
- cria uma goroutine associada à criação do leilão que valida, após intervalo, se o leilão venceu;
- quando vencido, atualiza o status do leilão para Completed (fechado);
- incluir testes que garantam o fechamento automático (unitário + integração).

---

## 🧾 Requisitos do Projeto

**Objetivo**: Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.  

Clone o seguinte repositório: [clique para acessar o repositório](https://github.com/devfullcycle/labs-auction-goexpert).  
Toda rotina de criação do leilão e lances já está desenvolvida, entretanto, o projeto clonado necessita de melhoria: adicionar a rotina de fechamento automático a partir de um tempo.  
Para essa tarefa, você utilizará o go routines e deverá se concentrar no processo de criação de leilão (auction). A validação do leilão (auction) estar fechado ou aberto na rotina de novos lançes (bid) já está implementado.

**Você deverá desenvolver:**  

- Uma função que irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente;
- Uma nova go routine que validará a existência de um leilão (auction) vencido (que o tempo já se esgotou) e que deverá realizar o update, fechando o leilão (auction);
- Um teste para validar se o fechamento está acontecendo de forma automatizada;

**Dicas:**

- Concentre-se na no arquivo internal/infra/database/auction/create_auction.go, você deverá implementar a solução nesse arquivo;
- Lembre-se que estamos trabalhando com concorrência, implemente uma solução que solucione isso:
- Verifique como o cálculo de intervalo para checar se o leilão (auction) ainda é válido está sendo realizado na rotina de criação de bid;
- Para mais informações de como funciona uma goroutine, clique aqui e acesse nosso módulo de Multithreading no curso Go Expert;
 
**Entrega:**

- O código-fonte completo da implementação.
- Documentação explicando como rodar o projeto em ambiente dev.
- Utilize docker/docker-compose para podermos realizar os testes de sua aplicação.

---

## 🧩 O que foi implementado (resumo)

- Implementação da rotina de fechamento automático em `internal/infra/database/auction/create_auction.go`.
- Teste unitário que valida comportamento do `create_auction.go` (inclui cenários de tempo padrão, parsing inválido e o fluxo de fechamento).
- Teste de integração para a rota de criação de leilão (`POST /auction`) cobrindo a criação e a verificação do status no banco.
- Mecanismo `APP_MODE` para facilitar testes:
  - `APP_MODE=test` altera comportamento da rotina de criação do lilão (aceita contexto externo, da request, para permitir cancelamento do contexto no teste — ver comentários nos testes).
  - `APP_MODE=dev` insere dados iniciais ao iniciar a aplicação (útil para testes manuais).
- Arquivos padrões para Docker (`Dockerfile`, `docker-compose.yml`) e `Makefile` com targets úteis (`make up`, `make down`, `make build`, etc).

---

## ⚙️ Variáveis de Ambiente

O arquivo `.env` fica em `cmd/auction/.env`.

Antes de rodar o projeto, copie o `.env.exemplo` para gerar o arquivo `.env`:

```bash
cp cmd/auction/.env.example cmd/auction/.env
```

| Variável | Exemplo | Descrição |
|----------|----|----------------|
| `BATCH_INSERT_INTERVAL` | `5s` | Intervalo de inserção em lote para registros. |
| `MAX_BATCH_SIZE` | `4` | Número máximo de itens em um batch. |
| `AUCTION_INTERVAL` | `120s` | Tempo de duração de um leilão antes do fechamento automático. |
| `APP_MODE` | `dev` | Define o modo da aplicação: `dev`, `test`, `prod`. |
| `MONGO_INITDB_ROOT_USERNAME` | `admin` | Usuário administrador do MongoDB. |
| `MONGO_INITDB_ROOT_PASSWORD` | `admin` | Senha do administrador do MongoDB. |
| `MONGODB_URL` | `mongodb://admin:admin@mongodb:27017/auctions_db?authSource=admin` | URL de conexão principal do MongoDB. |
| `MONGODB_DB` | `auctions_db` | Nome do banco principal. |
| `MONGODB_URL_TEST` | `mongodb://admin:admin@localhost:27017/auctions_db_test?authSource=admin` | URL de conexão para o banco de testes. |
| `MONGODB_DB_TEST` | `auctions_db_test` | Nome do banco de testes. |


---

## 🧱 Como Rodar o Projeto

O projeto utiliza **Makefile** para simplificar a execução dos comandos.  

### 💾 Passo 1 - Clonar repositório
```bash
git clone git@github.com:Berchon/labs-auction-goexpert.git
cd labs-auction-goexpert
```

### ⚙️ Passo 2 - Preparar .env

Crie o `.env`  a partir do `.env.eample`em `cmd/auction/`.
```bash
cp cmd/auction/.env.example cmd/auction/.env
```

>Dica: manter APP_MODE=dev enquanto faz testes manuais. Para rodar os testes de integração/unitários, use APP_MODE=test.

### 🪄 Passo 3 — Subir containers
```bash
make up
```
> Sobe a aplicação e o MongoDB usando Docker Compose.

### 🔨 Comandos principais

| Comando | Descrição |
|----------|------------|
| `make up` | Sobe os containers da aplicação e do banco |
| `make build` | Faz o build da imagem da aplicação |
| `make down` | Remove containers, rede e volumes |
| `make stop` | Para os containers em execução |
| `make clean` | Limpa imagens, volumes e containers |
| `make status` | Mostra o status atual dos containers |
| `make logs` | Exibe os logs da aplicação e do MongoDB |

---

## 🌐 Endpoints da Aplicação

As requisições estão organizadas na pasta `./api`, onde há **3 arquivos .http**:
- `auction.http`
- `bid.http`
- `user.http`

Cada arquivo contém exemplos prontos para testar as rotas da aplicação (compatíveis com o plugin “REST Client” do VSCode ou com `curl`).

### 🏷️ auction.http — Leilões
#### Criar um novo leilão
```bash
curl -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "Casa ABCD",
    "category": "Imóveis",
    "description": "Casa da Rua ABCD, 223",
    "condition": 0
  }'
```

#### Listar todos os leilões
```bash
curl http://localhost:8080/auction
```

#### Listar leilões por id
```bash
curl http://localhost:8080/auction/44c402b6-2960-4f9f-999f-5f217f40cee8
```

#### Listar leilões por status
```bash
curl http://localhost:8080/auction?status=0   # 0 = Active, 1 = Completed
```

#### Listar leilões usando query params
```bash
curl http://localhost:8080/auction?status=0&category=Doce&productName=Mandolate
```

### 💰 bid.http — Lances
#### Criar um novo lance
```bash
curl -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d '{
    "auction_id": "<AUCTION_ID>",
    "user_id": "<USER_ID>",
    "amount": 1500
  }'
```

#### Listar os lances de um leilão específico
```bash
curl http://localhost:8080/bid/44c402b6-2960-4f9f-999f-5f217f40cee8
```

#### Listar o lance vencedor de um leilão específico
```bash
curl http://localhost:8080/auction/winner/44c402b6-2960-4f9f-999f-5f217f40cee8
```

### 👤 user.http — Usuários
#### Buscar usuário por id
```bash
curl http://localhost:8080/user/e73fce6a-ccf5-4c12-87f7-30c5f9c9a6f7
```

---

## 🧪 Testes Automatizados

Os testes estão divididos entre **unitários** e **de integração**, todos executáveis via **Makefile**. Para rodar os testes não é necessário alterar o `APP_MODE` no `.env`, o teste já faz essa configuração automaticamente. Um teste de integração foi criado para validar o fechamento automático do leilão na pasta `internal/infra/database/auction`. Ainda, foi criado um teste de integração para a rota de criação de um novo leilão na pasta `tests/integration`.

---

### 🧩 Testes Unitários
```bash
make test-unit
```

### 🔗 Testes de Integração
```bash
make test-integration
```

### 🧬 Todos os testes
```bash
make test
```

---

## 🧠 Comportamento do Fechamento Automático

1. Quando um leilão é criado (`POST /auction`), a aplicação dispara uma **goroutine**.  
2. Essa goroutine aguarda o intervalo definido em `AUCTION_INTERVAL`.  
3. Ao atingir o tempo configurado, a rotina verifica se o leilão ainda está ativo e, se sim, **atualiza seu status para “Completed”**.  
4. O valor de `APP_MODE` define se o contexto da goroutine é independente (produção) ou controlado (testes).

---

## 🧰 Tecnologias e APIs utilizadas

* **Golang 1.24**
* **Docker** / **Docker Compose**
* **MongoDB**
* **Make** (Makefile com comandos de build/start/up/down/test)

---

## 📂 Estrutura Completa do Projeto

```bash
labs-auction-goexpert/
├── api/
│   ├── auction.http
│   ├── bid.http
│   └── user.http
│
├── cmd/
│   └── auction/
│       ├── .env
│       ├── .env.example
│       ├── main.go
│       └── docker-entrypoint.sh
│
├── internal/
│   ├── entity/
│   ├── infra/
│   └── usecase/
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## 🧾 Créditos

Desenvolvido como solução para desafio **Full Cycle — Go Expert: Leilão com Fechamento Automático**.  
Repositório base/fonte: [devfullcycle/labs-auction-goexpert](https://github.com/devfullcycle/labs-auction-goexpert)

---

## 👨‍💻 Autor

Projeto desenvolvido por **Berchon** — [https://github.com/Berchon](https://github.com/Berchon)