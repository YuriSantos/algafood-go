# 🐳 Troubleshooting Docker - AlgaFood

## 🚨 Problemas Identificados e Soluções

### Problema 1: Comandos Docker não retornam output
**Causa**: PowerShell pode estar suprimindo output ou Docker Desktop não está respondendo

**Soluções:**
```powershell
# 1. Verificar se Docker Desktop está rodando
Get-Process "Docker Desktop" -ErrorAction SilentlyContinue

# 2. Reiniciar Docker Desktop
Stop-Process -Name "Docker Desktop" -Force
Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"

# 3. Verificar serviços do Docker
Get-Service docker
```

### Problema 2: Docker Compose falha silenciosamente
**Causa**: Configuração do docker-compose.yml ou problemas de rede

**Soluções:**
```powershell
# 1. Verificar sintaxe do docker-compose.yml
docker-compose config

# 2. Executar com verbose
docker-compose --verbose up

# 3. Verificar logs
docker-compose logs --follow
```

### Problema 3: Containers não iniciam
**Possíveis Causas:**
- Portas já em uso
- Permissões insuficientes
- Recursos insuficientes
- Configuração inválida

**Verificações:**
```powershell
# Verificar portas em uso
netstat -ano | findstr :8080
netstat -ano | findstr :3306
netstat -ano | findstr :6379
netstat -ano | findstr :4566

# Verificar recursos
docker system df
docker system info
```

## 🔧 Script de Verificação

Execute este script para diagnosticar problemas:

```powershell
# Verificação completa do ambiente Docker
Write-Host "🔍 Diagnóstico Docker AlgaFood" -ForegroundColor Green

# 1. Verificar Docker Desktop
Write-Host "📋 Status Docker Desktop:" -ForegroundColor Yellow
try {
    $dockerProcess = Get-Process "Docker Desktop" -ErrorAction Stop
    Write-Host "✅ Docker Desktop está rodando (PID: $($dockerProcess.Id))" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker Desktop não está rodando" -ForegroundColor Red
    Write-Host "👉 Inicie o Docker Desktop e tente novamente" -ForegroundColor Yellow
    exit 1
}

# 2. Verificar serviços
Write-Host "`n📋 Serviços Docker:" -ForegroundColor Yellow
Get-Service docker | Format-Table -AutoSize

# 3. Verificar conectividade
Write-Host "`n📋 Teste de conectividade:" -ForegroundColor Yellow
try {
    $dockerVersion = docker version --format json | ConvertFrom-Json
    Write-Host "✅ Docker Engine: $($dockerVersion.Client.Version)" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker Engine não está respondendo" -ForegroundColor Red
}

# 4. Verificar portas
Write-Host "`n📋 Verificando portas:" -ForegroundColor Yellow
$ports = @(8080, 3306, 6379, 4566, 1025, 8025)
foreach ($port in $ports) {
    $inUse = netstat -ano | findstr ":$port"
    if ($inUse) {
        Write-Host "⚠️  Porta $port em uso" -ForegroundColor Yellow
    } else {
        Write-Host "✅ Porta $port disponível" -ForegroundColor Green
    }
}

# 5. Verificar recursos
Write-Host "`n📋 Recursos do sistema:" -ForegroundColor Yellow
docker system df 2>$null
```

## 🛠️ Soluções Rápidas

### Solução 1: Reset Completo
```powershell
# Parar todos os containers
docker stop $(docker ps -aq) 2>$null

# Remover todos os containers
docker rm $(docker ps -aq) 2>$null

# Limpar redes não utilizadas
docker network prune -f

# Limpar volumes órfãos
docker volume prune -f
```

### Solução 2: Restart Serviços
```powershell
# Reiniciar serviço Docker (como Admin)
Restart-Service docker

# Ou reiniciar Docker Desktop
Stop-Process -Name "Docker Desktop" -Force
Start-Sleep 5
Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
```

### Solução 3: Configuração Alternativa
Se o Docker Compose não funcionar, use containers individuais:

```powershell
# MySQL
docker run -d --name algafood-mysql `
  -p 13306:3306 `
  -e MYSQL_ROOT_PASSWORD=Teste@1992 `
  -e MYSQL_DATABASE=algafood `
  -e MYSQL_USER=algafood `
  -e MYSQL_PASSWORD=algafood123 `
  mysql:8.0

# Redis
docker run -d --name algafood-redis `
  -p 16379:6379 `
  redis:7-alpine

# LocalStack
docker run -d --name algafood-localstack `
  -p 4566:4566 `
  -e SERVICES=s3,ses,sqs,sns,eventbridge `
  localstack/localstack:3.0

# MailHog
docker run -d --name algafood-mailhog `
  -p 1025:1025 -p 8025:8025 `
  mailhog/mailhog
```

## 🔍 Verificações Finais

Após aplicar as soluções, execute:

```powershell
# 1. Verificar containers
docker ps

# 2. Testar conectividade
curl http://localhost:8025  # MailHog
curl http://localhost:4566/health  # LocalStack

# 3. Testar aplicação
./algafood-api.exe  # Se tudo estiver funcionando
```

## 📞 Se Nada Funcionar

1. **Reinicie o computador** - Às vezes resolve problemas de rede/serviços
2. **Reinstale Docker Desktop** - Download da versão mais recente
3. **Verifique antivírus/firewall** - Podem estar bloqueando Docker
4. **Execute como Administrador** - Alguns comandos precisam de privilégios
5. **Use WSL2** - Se disponível, pode resolver problemas no Windows

## 🎯 Teste Rápido

Execute este comando para teste rápido:
```powershell
docker run --rm nginx:alpine echo "Docker está funcionando!"
```

Se este comando não funcionar, o problema é com o Docker básico, não com a configuração do projeto.
