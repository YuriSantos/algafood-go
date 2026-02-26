# 🐳 Guia de Configuração Docker - AlgaFood

## 📁 Estrutura Criada

```
algafood-go/
├── 🐳 docker-compose.yml           # Orquestração completa dos serviços
├── 🐳 Dockerfile                   # Build da aplicação (desenvolvimento)
├── 🐳 Dockerfile.prod              # Build otimizado para produção
├── ⚙️ config.docker.yaml           # Configurações para Docker
├── 🚫 .dockerignore               # Arquivos ignorados no build
├── docker/
│   └── nginx/
│       └── nginx.conf              # Configuração do Nginx
└── scripts/
    ├── 🚀 start-dev.ps1            # Script Windows para desenvolvimento
    ├── 🚀 start-dev.sh             # Script Linux para desenvolvimento
    └── 🔧 localstack-setup.sh      # Setup automático do LocalStack
```

## 🚀 Como Usar

### 1. **Desenvolvimento Rápido**

**Windows:**
```powershell
.\scripts\start-dev.ps1
```

**Linux/Mac:**
```bash
chmod +x scripts/start-dev.sh
./scripts/start-dev.sh
```

### 2. **Comandos Manuais**

```bash
# Iniciar todos os serviços
docker-compose up -d

# Reconstruir e iniciar
docker-compose up --build -d

# Ver logs
docker-compose logs -f algafood-api

# Parar serviços
docker-compose down
```

## 🛠️ Serviços Incluídos

| Serviço | Porta | Descrição |
|---------|-------|-----------|
| **algafood-api** | 8080 | API principal Go |
| **algafood-mysql** | 13306 | Banco MySQL 8.0 |
| **algafood-redis** | 16379 | Cache Redis 7 |
| **localstack** | 4566 | Simulação AWS |
| **mailhog** | 8025 | Interface de emails |
| **nginx** | 80/443 | Proxy reverso |

## ☁️ LocalStack (AWS Local)

O LocalStack simula serviços AWS localmente:

### **Serviços Configurados:**
- ✅ **S3** - Armazenamento de arquivos
- ✅ **SES** - Envio de emails
- ✅ **SQS** - Filas de mensagens
- ✅ **SNS** - Notificações
- ✅ **EventBridge** - Eventos
- ✅ **CloudWatch** - Monitoramento

### **Recursos Criados Automaticamente:**
- Buckets S3: `algafood-files`, `algafood-fotos-produtos`
- Filas SQS: `algafood-pedido-status`
- Emails verificados: `teste@algafood.com.br`, `admin@algafood.com.br`
- EventBridge: `algafood-event-bus`

## 📊 Monitoramento

### **Health Checks:**
```bash
# Status geral
docker-compose ps

# Health check da API
curl http://localhost:8080/health

# Status do LocalStack
curl http://localhost:4566/health
```

### **Logs:**
```bash
# Logs da API
docker-compose logs -f algafood-api

# Logs do LocalStack
docker-compose logs -f localstack

# Todos os logs
docker-compose logs -f
```

## 📧 Emails de Desenvolvimento

### **MailHog (Recomendado):**
```bash
# Interface web
http://localhost:8025

# SMTP interno
mailhog:1025
```

### **LocalStack SES:**
```bash
# Verificar emails enviados
.\scripts\email-checker-simple.ps1

# Interface visual
start .\scripts\email-viewer-fixed.html
```

## ⚙️ Configurações

### **Variáveis de Ambiente (docker-compose.yml):**
```yaml
# Banco de dados
DB_HOST=algafood-mysql
DB_USER=algafood
DB_PASSWORD=algafood123

# Redis
REDIS_HOST=algafood-redis

# AWS (LocalStack)
AWS_ENDPOINT_URL=http://localstack:4566
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test

# SMTP (MailHog)
SMTP_HOST=mailhog
SMTP_PORT=1025
```

### **Configuração Personalizada:**
Edite `config.docker.yaml` para ajustes específicos.

## 🔧 Scripts de Desenvolvimento

### **Windows (PowerShell):**
```powershell
# Iniciar ambiente completo
.\scripts\start-dev.ps1

# Parar todos os serviços
.\scripts\start-dev.ps1 -Stop

# Reconstruir containers
.\scripts\start-dev.ps1 -Rebuild

# Ver logs em tempo real
.\scripts\start-dev.ps1 -Logs

# Verificar emails
.\scripts\email-checker-simple.ps1
```

### **Comandos Úteis:**
```powershell
# Status dos containers
docker-compose ps

# Entrar no container da API
docker exec -it algafood-api sh

# Backup do banco
docker exec algafood-mysql mysqldump -u root -pTeste@1992 algafood > backup.sql

# Restaurar banco
docker exec -i algafood-mysql mysql -u root -pTeste@1992 algafood < backup.sql
```

## 🏭 Produção

### **Build Otimizado:**
```bash
# Build para produção
docker build -f Dockerfile.prod -t algafood-api:latest .

# Executar em produção
docker run -p 8080:8080 algafood-api:latest
```

### **Características da Build de Produção:**
- ✅ Imagem minimal (scratch)
- ✅ Binário estático
- ✅ Usuário não-root
- ✅ Health check incluído
- ✅ Certificados SSL
- ✅ Timezone configurado

## 🚨 Troubleshooting

### **Problema: Containers não iniciam**
```bash
# Verificar logs
docker-compose logs

# Recriar volumes
docker-compose down -v
docker-compose up -d
```

### **Problema: Porta ocupada**
```bash
# Verificar portas em uso
netstat -tulpn | grep :8080

# Mudar porta no docker-compose.yml
ports:
  - "8081:8080"  # Usar porta 8081 externamente
```

### **Problema: LocalStack não responde**
```bash
# Verificar status
curl http://localhost:4566/health

# Reiniciar apenas LocalStack
docker-compose restart localstack
```

## 📋 Checklist de Verificação

Após executar `.\scripts\start-dev.ps1`:

- ✅ Todos os containers estão rodando: `docker-compose ps`
- ✅ API responde: `curl http://localhost:8080/health`
- ✅ MySQL conecta: Testar conexão na porta 13306
- ✅ Redis funcionando: `docker exec algafood-redis redis-cli ping`
- ✅ LocalStack ativo: `curl http://localhost:4566/health`
- ✅ MailHog acessível: `http://localhost:8025`

## 🎯 Resultado Final

Com essa configuração Docker você tem:

1. ✅ **Ambiente completo** de desenvolvimento
2. ✅ **LocalStack** simulando AWS
3. ✅ **MailHog** para emails de desenvolvimento
4. ✅ **Scripts automatizados** para Windows e Linux
5. ✅ **Health checks** e monitoramento
6. ✅ **Build otimizado** para produção
7. ✅ **Proxy Nginx** configurado
8. ✅ **Persistência de dados** com volumes

**Comando único para começar:**
```powershell
.\scripts\start-dev.ps1
```

🎉 **Infraestrutura pronta para desenvolvimento e produção!**
