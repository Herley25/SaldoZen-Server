# 📊 SALDOZEN API

API RESTful desenvolvida em **Go** com arquitetura **MVC**, persistência em **PostgreSQL** e autenticação baseada em UUID. Ideal para controle e gestão de despesas com foco em escalabilidade e boas práticas.

---

## 🚀 Tecnologias Utilizadas

- [Golang](https://go.dev/) — Backend principal
- [PostgreSQL](https://www.postgresql.org/) — Banco de dados relacional
- [Gorilla Mux](https://github.com/gorilla/mux) — Roteador HTTP
- [joho/godotenv](https://github.com/joho/godotenv) — Carregamento de variáveis de ambiente
- [Docker + Docker Compose](https://docs.docker.com/compose/) — Containerização e ambiente isolado
- [uuid](https://pkg.go.dev/github.com/google/uuid) — Identificadores únicos para entidades

---

## 📁 Estrutura do Projeto

```shell
.
├── cmd/                # Entry point da aplicação
│   └── main.go
├── config/             # Configurações e variáveis de ambiente
│   └── config.go
├── controllers/        # Controladores (camada de negócio e HTTP)
│   └── expense_controller.go
├── db/                 # Conexão e inicialização do banco
│   └── database.go
├── models/             # Modelos das entidades (ORM-like)
│   ├── user.go
│   └── expense.go
├── routes/             # Definição das rotas da API
│   └── routes.go
├── migrations/         # Scripts SQL para criação de tabelas
│   └── 001_init.sql
├── .env                # Variáveis de ambiente (não versionado)
├── .gitignore
└── go.mod / go.sum     # Dependências do projeto
```
## 🧪 Funcionalidades Implementadas

✅ Cadastro e gerenciamento de despesas

✅ Vinculação de despesas a usuários

✅ Validações de valores e vencimentos

✅ Cálculo de status da despesa (A vencer, Vencida, Paga)

✅ Criação automática de timestamp

✅ Organização em camadas (MVC)

✅ Migrations versionadas em SQL

✅ Estrutura pronta para deploy com Docker

## ⚙️ Como Rodar o Projeto

1. Clone o repositório
```shell
git clone https://github.com/seu-usuario/expense-manager-go.git
cd expense-manager-go
```

2. Configure seu .env

Crie um arquivo .env com as seguintes variáveis:
```shell
DATABASE_URL=postgres://usuario:senha@localhost:5432/expense_db?sslmode=disable
```

3. Suba os serviços com Docker Compose
```shell
docker-compose up --build
```
Ou rode localmente (com banco configurado):
```shell
go run cmd/main.go
```

## 🛠️ Migrations
As tabelas podem ser criadas executando os scripts em ``` /migrations/001_init.sql. ```
