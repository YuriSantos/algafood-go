# AlgaFood API - Go

API de delivery de comida implementada em Go, migrada do projeto original em Java/Spring Boot (AlgaFood).

## 🚀 Stack Tecnológica

- **Go 1.21+**
- **Gin** - Web framework
- **GORM** - ORM (Object Relational Mapper)
- **MySQL** - Banco de dados
- **JWT/JWKS** - Autenticação OAuth2
- **Viper** - Gerenciamento de configuração
- **SendGrid** - Serviço de envio de e-mails
- **AWS S3** - Armazenamento de arquivos (opcional)

## 📂 Estrutura do Projeto

```
algafood-go/
├── cmd/api/             # Ponto de entrada da aplicação (main)
├── internal/
│   ├── api/             # Camada de API (handlers, rotas, middlewares, dtos)
│   ├── config/          # Carregamento de configurações
│   ├── domain/          # Domínio (models, services, repositories, interfaces)
│   └── infrastructure/  # Implementações de infraestrutura (storage, email, db)
├── pkg/                 # Pacotes reutilizáveis e utilitários
├── migrations/          # Scripts SQL de migração
└── config.yaml          # Arquivo de configuração base
```

## 📋 Pré-requisitos

- Go 1.21 ou superior
- MySQL 8.0 ou superior
- (Opcional) Authorization Server OAuth2 rodando em `localhost:8080` (para validação de tokens JWT)

## ⚙️ Configuração

1. **Clone o repositório**

2. **Configure as variáveis de ambiente**
   Copie o arquivo de configuração de exemplo e edite conforme necessário:

   ```bash
   cp config.yaml config.local.yaml
   # Edite config.local.yaml com suas credenciais de banco, AWS, SendGrid, etc.
   ```

3. **Banco de Dados**
   Crie o banco de dados no MySQL:

   ```sql
   CREATE DATABASE algafood CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

4. **Migrations**
   Execute as migrações para criar as tabelas e popular dados iniciais (seed):

   ```bash
   # Exemplo via linha de comando
   mysql -u root -p algafood < migrations/000001_create_tables.up.sql
   mysql -u root -p algafood < migrations/000002_seed_data.up.sql
   ```

## ▶️ Executando

```bash
# Baixar dependências
go mod tidy

# Executar a aplicação
go run cmd/api/main.go
```

A API estará disponível em `http://localhost:8081`.

## 🔌 Endpoints Principais

Abaixo estão listados os principais recursos da API.

### Cadastros Básicos
- `GET /v1/estados` - Listar estados
- `GET /v1/cidades` - Listar cidades
- `GET /v1/cozinhas` - Listar cozinhas (paginado)

### Restaurantes
- `GET /v1/restaurantes` - Listar restaurantes
- `POST /v1/restaurantes` - Cadastrar restaurante
- `PUT /v1/restaurantes/:id` - Atualizar dados
- `PUT /v1/restaurantes/:id/ativo` - Ativar restaurante
- `PUT /v1/restaurantes/:id/abertura` - Abrir restaurante para pedidos

### Produtos
- `GET /v1/restaurantes/:id/produtos` - Listar produtos do restaurante
- `POST /v1/restaurantes/:id/produtos` - Adicionar produto
- `PUT /v1/restaurantes/:id/produtos/:prodId/foto` - Upload de foto do produto

### Pedidos
- `GET /v1/pedidos` - Pesquisar pedidos (com filtros)
- `POST /v1/pedidos` - Emitir novo pedido
- `PUT /v1/pedidos/:codigo/confirmacao` - Confirmar pedido
- `PUT /v1/pedidos/:codigo/entrega` - Registrar entrega
- `PUT /v1/pedidos/:codigo/cancelamento` - Cancelar pedido

### Usuários
- `GET /v1/usuarios` - Listar usuários
- `POST /v1/usuarios` - Cadastrar usuário
- `PUT /v1/usuarios/:id/senha` - Alterar senha

### Estatísticas
- `GET /v1/estatisticas/vendas-diarias` - Relatório de vendas diárias

## 🔒 Autenticação

A API suporta autenticação OAuth2 via JWT (Resource Server).

1. Configure a URL do JWKS no `config.yaml`:
   ```yaml
   jwt:
     jwks_url: "http://localhost:8080/oauth2/jwks"
   ```

2. O middleware de autenticação validará o token Bearer nas requisições protegidas.

## 📝 Exemplos de Requisições

### Criar Restaurante
```bash
curl -X POST http://localhost:8081/v1/restaurantes \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Thai Gourmet",
    "taxaFrete": 10.00,
    "cozinha": {"id": 1},
    "endereco": {
        "cep": "38400-999",
        "logradouro": "Rua João Pinheiro",
        "numero": "1000",
        "bairro": "Centro",
        "cidade": {"id": 1}
    }
  }'
```

### Criar Pedido
```bash
curl -X POST http://localhost:8081/v1/pedidos \
  -H "Content-Type: application/json" \
  -d '{
    "restaurante": {"id": 1},
    "formaPagamento": {"id": 1},
    "enderecoEntrega": {
      "cep": "38400-000",
      "logradouro": "Rua Floriano Peixoto",
      "numero": "500",
      "bairro": "Centro",
      "cidade": {"id": 1}
    },
    "itens": [
      {"produtoId": 1, "quantidade": 2, "observacao": "Sem cebola"}
    ]
  }'
```

## 📄 Licença

Este projeto foi desenvolvido para fins educacionais.
