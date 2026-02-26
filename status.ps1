param(
    [switch]$Logs,
    [switch]$Restart,
    [switch]$Status
)

Write-Host "🚀 AlgaFood Docker Status" -ForegroundColor Green

if ($Restart) {
    Write-Host "🔄 Reiniciando containers..." -ForegroundColor Yellow
    docker-compose down 2>$null
    Start-Sleep 3
    docker-compose up -d 2>$null
    Start-Sleep 10
}

if ($Status) {
    Write-Host "`n📊 Status dos containers:" -ForegroundColor Cyan
    $containers = docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" 2>$null
    if ($containers) {
        $containers
    } else {
        Write-Host "❌ Nenhum container em execução ou Docker não está respondendo" -ForegroundColor Red
    }
}

if ($Logs) {
    Write-Host "`n📋 Logs da API (últimas 20 linhas):" -ForegroundColor Yellow
    $apiLogs = docker-compose logs algafood-api --tail=20 2>$null
    if ($apiLogs) {
        $apiLogs
    } else {
        Write-Host "❌ Não foi possível obter logs da API" -ForegroundColor Red
    }
}

# Status automático se nenhuma flag foi especificada
if (-not ($Logs -or $Restart -or $Status)) {
    Write-Host "`n📊 Status dos containers:" -ForegroundColor Cyan
    docker ps --format "table {{.Names}}\t{{.Status}}" 2>$null

    Write-Host "`n📋 Logs da API (últimas 10 linhas):" -ForegroundColor Yellow
    docker-compose logs algafood-api --tail=10 2>$null

    Write-Host "`n💡 Uso:" -ForegroundColor Magenta
    Write-Host "  .\status.ps1 -Status    # Ver status detalhado" -ForegroundColor Gray
    Write-Host "  .\status.ps1 -Logs      # Ver logs da API" -ForegroundColor Gray
    Write-Host "  .\status.ps1 -Restart   # Reiniciar containers" -ForegroundColor Gray
}
