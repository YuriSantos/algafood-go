# Script para iniciar containers individuais como alternativa ao docker-compose
param(
    [switch]$Start,
    [switch]$Stop,
    [switch]$Status,
    [switch]$Logs
)

Write-Host "🐳 AlgaFood Individual Containers" -ForegroundColor Green
Write-Host "=" * 50 -ForegroundColor Gray

if ($Stop) {
    Write-Host "🛑 Parando todos os containers AlgaFood..." -ForegroundColor Yellow

    $containers = @("algafood-api-individual", "algafood-mysql-individual", "algafood-redis-individual", "algafood-localstack-individual", "algafood-mailhog-individual")

    foreach ($container in $containers) {
        Write-Host "Parando $container..." -ForegroundColor Gray
        docker stop $container 2>$null
        docker rm $container 2>$null
    }

    Write-Host "✅ Containers parados!" -ForegroundColor Green
    exit 0
}

if ($Status) {
    Write-Host "📊 Status dos containers AlgaFood:" -ForegroundColor Yellow
    docker ps --filter "name=algafood-" --format "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}"
    exit 0
}

if ($Logs) {
    Write-Host "📋 Logs da API:" -ForegroundColor Yellow
    docker logs algafood-api-individual --tail=20
    exit 0
}

if ($Start) {
    Write-Host "🚀 Iniciando containers individuais..." -ForegroundColor Cyan

    # 1. MySQL
    Write-Host "`n📦 Iniciando MySQL..." -ForegroundColor Yellow
    docker run -d --name algafood-mysql-individual `
        -p 13306:3306 `
        -e MYSQL_ROOT_PASSWORD=Teste@1992 `
        -e MYSQL_DATABASE=algafood `
        -e MYSQL_USER=algafood `
        -e MYSQL_PASSWORD=algafood123 `
        mysql:8.0

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ MySQL iniciado na porta 13306" -ForegroundColor Green
    } else {
        Write-Host "❌ Falha ao iniciar MySQL" -ForegroundColor Red
        exit 1
    }

    Start-Sleep 5

    # 2. Redis
    Write-Host "`n📦 Iniciando Redis..." -ForegroundColor Yellow
    docker run -d --name algafood-redis-individual `
        -p 16379:6379 `
        redis:7-alpine

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Redis iniciado na porta 16379" -ForegroundColor Green
    } else {
        Write-Host "❌ Falha ao iniciar Redis" -ForegroundColor Red
    }

    # 3. MailHog
    Write-Host "`n📦 Iniciando MailHog..." -ForegroundColor Yellow
    docker run -d --name algafood-mailhog-individual `
        -p 1025:1025 -p 8025:8025 `
        mailhog/mailhog

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ MailHog iniciado - Interface: http://localhost:8025" -ForegroundColor Green
    } else {
        Write-Host "❌ Falha ao iniciar MailHog" -ForegroundColor Red
    }

    # 4. LocalStack
    Write-Host "`n📦 Iniciando LocalStack..." -ForegroundColor Yellow
    docker run -d --name algafood-localstack-individual `
        -p 4566:4566 `
        -e SERVICES=s3,ses,sqs,sns,eventbridge `
        -e DEBUG=1 `
        -e AWS_DEFAULT_REGION=us-east-1 `
        -e AWS_ACCESS_KEY_ID=test `
        -e AWS_SECRET_ACCESS_KEY=test `
        localstack/localstack:3.0

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ LocalStack iniciado na porta 4566" -ForegroundColor Green
    } else {
        Write-Host "❌ Falha ao iniciar LocalStack" -ForegroundColor Red
    }

    Write-Host "`n⏳ Aguardando serviços iniciarem..." -ForegroundColor Cyan
    Start-Sleep 15

    # 5. API (usando a imagem já construída)
    Write-Host "`n📦 Iniciando API AlgaFood..." -ForegroundColor Yellow

    # Primeiro, cria um arquivo de configuração específico para containers individuais
    $configContent = @"
server:
  port: 8080

database:
  host: host.docker.internal
  port: 13306
  user: algafood
  password: algafood123
  name: algafood
  charset: utf8mb4
  parseTime: true
  loc: Local

redis:
  host: host.docker.internal
  port: 16379
  password: ""
  db: 0

jwt:
  issuer: "algafood-api"

email:
  impl: "SES"
  remetente: "noreply@algafood.com.br"
  sandbox:
    destinatario: "teste@algafood.com.br"

aws:
  region: us-east-1
  endpoint: http://host.docker.internal:4566
  credentials:
    access_key: test
    secret_key: test

sqs:
  queue:
    pedido_status: algafood-pedido-status

eventbridge:
  bus_name: algafood-event-bus

storage:
  s3:
    bucket: algafood-fotos-produtos
"@

    # Salva a configuração
    $configContent | Out-File -FilePath "config-individual.yaml" -Encoding UTF8

    # Inicia a API
    docker run -d --name algafood-api-individual `
        -p 8080:8080 `
        -v "${PWD}/config-individual.yaml:/root/config.yaml:ro" `
        algafood-go-algafood-api:latest

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ API AlgaFood iniciada na porta 8080" -ForegroundColor Green
    } else {
        Write-Host "❌ Falha ao iniciar API" -ForegroundColor Red
    }

    Write-Host "`n🎉 Serviços iniciados!" -ForegroundColor Green
    Write-Host "`n📋 URLs dos serviços:" -ForegroundColor Magenta
    Write-Host "   🌐 API AlgaFood:     http://localhost:8080" -ForegroundColor White
    Write-Host "   📧 MailHog:          http://localhost:8025" -ForegroundColor White
    Write-Host "   ☁️  LocalStack:      http://localhost:4566" -ForegroundColor White
    Write-Host "   🗄️ MySQL:            localhost:13306" -ForegroundColor White
    Write-Host "   🔴 Redis:            localhost:16379" -ForegroundColor White

    Write-Host "`n💡 Comandos úteis:" -ForegroundColor Magenta
    Write-Host "   .\individual-containers.ps1 -Status  # Ver status" -ForegroundColor Gray
    Write-Host "   .\individual-containers.ps1 -Logs    # Ver logs da API" -ForegroundColor Gray
    Write-Host "   .\individual-containers.ps1 -Stop    # Parar tudo" -ForegroundColor Gray

} else {
    Write-Host "💡 Uso:" -ForegroundColor Yellow
    Write-Host "  .\individual-containers.ps1 -Start   # Iniciar todos os containers" -ForegroundColor Gray
    Write-Host "  .\individual-containers.ps1 -Status  # Ver status" -ForegroundColor Gray
    Write-Host "  .\individual-containers.ps1 -Logs    # Ver logs da API" -ForegroundColor Gray
    Write-Host "  .\individual-containers.ps1 -Stop    # Parar todos os containers" -ForegroundColor Gray
}
