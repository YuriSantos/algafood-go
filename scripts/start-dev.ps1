# Script PowerShell para desenvolvimento - inicia todos os serviços
param(
    [switch]$Stop,
    [switch]$Rebuild,
    [switch]$Logs
)

Write-Host "🚀 AlgaFood Development Environment" -ForegroundColor Green
Write-Host "=" * 50 -ForegroundColor Gray

if ($Stop) {
    Write-Host "🛑 Parando todos os serviços..." -ForegroundColor Yellow
    docker-compose down
    Write-Host "✅ Serviços parados!" -ForegroundColor Green
    exit 0
}

if ($Logs) {
    Write-Host "📋 Mostrando logs dos serviços..." -ForegroundColor Yellow
    docker-compose logs -f
    exit 0
}

# Verificar se Docker está rodando
try {
    docker info | Out-Null
} catch {
    Write-Host "❌ Docker não está rodando. Por favor, inicie o Docker Desktop." -ForegroundColor Red
    exit 1
}

# Parar containers existentes
Write-Host "🛑 Parando containers existentes..." -ForegroundColor Yellow
docker-compose down

# Construir e iniciar serviços
if ($Rebuild) {
    Write-Host "🔨 Reconstruindo e iniciando serviços..." -ForegroundColor Yellow
    docker-compose up --build -d
} else {
    Write-Host "🔨 Iniciando serviços..." -ForegroundColor Yellow
    docker-compose up -d
}

# Aguardar serviços estarem prontos
Write-Host "⏳ Aguardando serviços iniciarem..." -ForegroundColor Yellow

# Verificar MySQL
Write-Host "🔍 Verificando MySQL..." -ForegroundColor Cyan
do {
    Start-Sleep 2
    $mysqlReady = docker exec algafood-mysql mysqladmin ping -h localhost --silent 2>$null
} while ($LASTEXITCODE -ne 0)
Write-Host "✅ MySQL pronto!" -ForegroundColor Green

# Verificar Redis
Write-Host "🔍 Verificando Redis..." -ForegroundColor Cyan
do {
    Start-Sleep 2
    docker exec algafood-redis redis-cli ping 2>$null | Out-Null
} while ($LASTEXITCODE -ne 0)
Write-Host "✅ Redis pronto!" -ForegroundColor Green

# Verificar LocalStack
Write-Host "🔍 Verificando LocalStack..." -ForegroundColor Cyan
do {
    Start-Sleep 2
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:4566/health" -TimeoutSec 5 -ErrorAction Stop
        $ready = $true
    } catch {
        $ready = $false
    }
} while (-not $ready)
Write-Host "✅ LocalStack pronto!" -ForegroundColor Green

# Mostrar status dos serviços
Write-Host "`n📊 Status dos serviços:" -ForegroundColor Yellow
docker-compose ps

Write-Host ""
Write-Host "🎉 Infraestrutura AlgaFood iniciada com sucesso!" -ForegroundColor Green
Write-Host ""
Write-Host "📋 Serviços disponíveis:" -ForegroundColor Magenta
Write-Host "   🌐 API AlgaFood:     http://localhost:8080" -ForegroundColor White
Write-Host "   🌐 Nginx (Proxy):   http://localhost:80" -ForegroundColor White
Write-Host "   📧 MailHog:          http://localhost:8025" -ForegroundColor White
Write-Host "   ☁️  LocalStack:      http://localhost:4566" -ForegroundColor White
Write-Host "   🗄️ MySQL:            localhost:13306" -ForegroundColor White
Write-Host "   🔴 Redis:            localhost:16379" -ForegroundColor White
Write-Host ""
Write-Host "🔧 Comandos úteis:" -ForegroundColor Magenta
Write-Host "   .\start-dev.ps1 -Logs         # Ver logs" -ForegroundColor Gray
Write-Host "   .\start-dev.ps1 -Stop         # Parar serviços" -ForegroundColor Gray
Write-Host "   .\start-dev.ps1 -Rebuild      # Reconstruir" -ForegroundColor Gray
Write-Host "   docker-compose ps             # Status dos containers" -ForegroundColor Gray
Write-Host ""
Write-Host "📧 Para verificar emails:" -ForegroundColor Magenta
Write-Host "   .\scripts\email-checker-simple.ps1" -ForegroundColor Gray
Write-Host "   start .\scripts\email-viewer-fixed.html" -ForegroundColor Gray
Write-Host ""
Write-Host "🌐 Abrir interfaces:" -ForegroundColor Magenta
Write-Host "   start http://localhost:8080   # API" -ForegroundColor Gray
Write-Host "   start http://localhost:8025   # MailHog" -ForegroundColor Gray
Write-Host "   start http://localhost:4566   # LocalStack" -ForegroundColor Gray
