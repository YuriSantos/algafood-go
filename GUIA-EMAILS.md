# 📧 Guia Completo: Como Visualizar Emails no LocalStack

## 🎯 Soluções Funcionais Criadas

### 1. **Interface Web Corrigida**
```powershell
start .\scripts\email-viewer-fixed.html
```
- Interface HTML totalmente funcional, sem erros
- Copia comandos para clipboard com um clique
- Links diretos para LocalStack e MailHog

### 2. **Script PowerShell Funcional**
```powershell
.\scripts\email-checker-simple.ps1        # Ver tudo
.\scripts\email-checker-simple.ps1 -Stats # Apenas estatísticas
.\scripts\email-checker-simple.ps1 -Logs  # Apenas logs
.\scripts\email-checker-simple.ps1 -MailHog # Configurar MailHog
```

## 📊 Estatísticas de Envio (Método Principal)

O LocalStack mantém estatísticas de envio que você pode verificar:

```powershell
# Configurar AWS CLI
$env:AWS_ACCESS_KEY_ID = "test"
$env:AWS_SECRET_ACCESS_KEY = "test"
$env:AWS_DEFAULT_REGION = "us-east-1"

# Ver estatísticas de envio
aws --endpoint-url=http://localhost:4566 ses get-send-statistics --output table
```

**Resultado atual**: 33 tentativas de envio, 6 rejeições

## 🔍 Métodos de Visualização

### 1. **Interface Web Corrigida** (RECOMENDADO)
```powershell
start .\scripts\email-viewer-fixed.html
```
- ✅ Sem erros de JavaScript
- ✅ Copia comandos automaticamente
- ✅ Links funcionais
- ✅ Design responsivo

### 2. **Logs do LocalStack** 
```powershell
# Ver logs em tempo real
docker logs algafood-localstack-1 -f

# Ver logs das últimas horas
docker logs algafood-localstack-1 --since="2h"

# Filtrar por SES
docker logs algafood-localstack-1 --since="1h" | findstr -i "ses\|email"
```

### 3. **Captura de Emails com MailHog** (Para Ver Conteúdo Real)

```powershell
# Executar MailHog
docker run --rm -d --name mailhog -p 1025:1025 -p 8025:8025 mailhog/mailhog

# Acessar interface web
start http://localhost:8025
```

**Para usar com sua aplicação:**
1. Configure o SMTP para `localhost:1025` 
2. Todos os emails aparecerão na interface http://localhost:8025

### 4. **Interface Web LocalStack**
- Abrir: http://localhost:4566
- Pode não ter interface visual na versão Community

## ✅ Verificação Rápida

Execute este comando único para verificar tudo:

```powershell
Write-Host "📊 Estatísticas:"; $env:AWS_ACCESS_KEY_ID="test"; $env:AWS_SECRET_ACCESS_KEY="test"; aws --endpoint-url=http://localhost:4566 ses get-send-statistics --output table; Write-Host "`n🐳 Containers:"; docker ps --format "table {{.Names}}\t{{.Status}}" | findstr "localstack\|mailhog\|algafood"
```

## 🚨 Status Atual

- ✅ **LocalStack funcionando**: 33 emails enviados
- ⚠️ **6 emails rejeitados**: Endereços não verificados no SES
- ✅ **Interface corrigida**: `email-viewer-fixed.html`
- ✅ **Scripts funcionais**: `email-checker-simple.ps1`

## 🛠️ Solução de Problemas

### Se não está enviando emails:
1. Verificar se LocalStack está rodando: `docker ps`
2. Verificar endpoint: `curl http://localhost:4566/health`
3. Verificar logs da aplicação Go

### Para desenvolvimento com emails reais:
1. **Use MailHog** (recomendado): `.\scripts\email-checker-simple.ps1 -MailHog`
2. Configure SMTP real (Gmail, SendGrid, etc.)
3. Use serviços de email de teste (Mailtrap, etc.)

## 📱 Scripts e Interfaces Criados

| Arquivo | Descrição | Status |
|---------|-----------|---------|
| `email-viewer-fixed.html` | Interface web corrigida | ✅ Funcional |
| `email-checker-simple.ps1` | Script PowerShell simples | ✅ Funcional |
| `GUIA-EMAILS.md` | Este guia | ✅ Atualizado |

## 🚀 Uso Rápido

**1. Ver estatísticas:**
```powershell
.\scripts\email-checker-simple.ps1 -Stats
```

**2. Configurar MailHog para desenvolvimento:**
```powershell
.\scripts\email-checker-simple.ps1 -MailHog
start http://localhost:8025
```

**3. Interface web completa:**
```powershell
start .\scripts\email-viewer-fixed.html
```

## 💡 Dica para Produção

Em produção, use:
- **AWS SES real** para envios
- **CloudWatch** para monitoramento  
- **SNS** para notificações de bounce/complaint

## 🎉 Resultado

Agora você tem **4 formas funcionais** de visualizar emails:
1. Interface HTML sem erros
2. Script PowerShell funcional
3. Comandos diretos no terminal
4. MailHog para emails reais

**Problema resolvido!** ✅

