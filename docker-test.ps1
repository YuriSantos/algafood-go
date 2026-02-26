param(
    [switch]$Build,
    [switch]$Up,
    [switch]$Down,
    [switch]$Status
)

Write-Host "🐳 AlgaFood Docker Manager" -ForegroundColor Green
Write-Host "=" * 40 -ForegroundColor Gray

if ($Down) {
    Write-Host "🛑 Parando containers..." -ForegroundColor Yellow
    docker-compose down
    exit 0
}

if ($Status) {
    Write-Host "📊 Status dos containers:" -ForegroundColor Yellow
    docker-compose ps
    exit 0
}

if ($Build) {
    Write-Host "🔨 Construindo containers..." -ForegroundColor Yellow
    docker-compose build --no-cache
}

if ($Up -or $Build) {
    Write-Host "🚀 Iniciando containers..." -ForegroundColor Yellow
    if ($Build) {
        docker-compose up --build -d
    } else {
        docker-compose up -d
    }

    Write-Host "⏳ Aguardando serviços..." -ForegroundColor Cyan
    Start-Sleep 10

    Write-Host "`n📊 Status final:" -ForegroundColor Green
    docker-compose ps

    Write-Host "`n🌐 Serviços disponíveis:" -ForegroundColor Magenta
    Write-Host "• API: http://localhost:8080" -ForegroundColor White
    Write-Host "• MailHog: http://localhost:8025" -ForegroundColor White
    Write-Host "• LocalStack: http://localhost:4566" -ForegroundColor White
    Write-Host "• MySQL: localhost:13306" -ForegroundColor White
    Write-Host "• Redis: localhost:16379" -ForegroundColor White
} else {
    Write-Host "💡 Uso:" -ForegroundColor Yellow
    Write-Host "  .\docker-test.ps1 -Up      # Iniciar"
    Write-Host "  .\docker-test.ps1 -Build   # Construir e iniciar"
    Write-Host "  .\docker-test.ps1 -Status  # Ver status"
    Write-Host "  .\docker-test.ps1 -Down    # Parar"
}
