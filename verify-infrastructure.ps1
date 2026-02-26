param(
    [switch]$Full,
    [switch]$Quick
)

Write-Host "🔍 Verificação Infraestrutura AlgaFood" -ForegroundColor Green
Write-Host "=" * 50 -ForegroundColor Gray

# Função para testar serviço
function Test-Service {
    param($Name, $Url, $Port)

    try {
        $response = Invoke-WebRequest -Uri $Url -TimeoutSec 5 -UseBasicParsing -ErrorAction Stop
        Write-Host "✅ $Name ($Port): OK - Status $($response.StatusCode)" -ForegroundColor Green
        return $true
    } catch {
        Write-Host "❌ $Name ($Port): FALHOU - $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Função para verificar porta
function Test-Port {
    param($Port, $Service)

    try {
        $connection = Test-NetConnection -ComputerName localhost -Port $Port -WarningAction SilentlyContinue
        if ($connection.TcpTestSucceeded) {
            Write-Host "✅ Porta $Port ($Service): ABERTA" -ForegroundColor Green
            return $true
        } else {
            Write-Host "❌ Porta $Port ($Service): FECHADA" -ForegroundColor Red
            return $false
        }
    } catch {
        Write-Host "❌ Porta $Port ($Service): ERRO - $_" -ForegroundColor Red
        return $false
    }
}

if ($Quick) {
    Write-Host "`n📊 Teste Rápido das Portas:" -ForegroundColor Yellow

    $ports = @(
        @{Port=8080; Service="API AlgaFood"},
        @{Port=13306; Service="MySQL"},
        @{Port=16379; Service="Redis"},
        @{Port=4566; Service="LocalStack"},
        @{Port=8025; Service="MailHog"}
    )

    foreach ($portInfo in $ports) {
        Test-Port -Port $portInfo.Port -Service $portInfo.Service
    }
    exit 0
}

# Verificação completa
Write-Host "`n📊 Status dos Containers:" -ForegroundColor Yellow
try {
    $containers = docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}" 2>$null
    if ($containers) {
        Write-Host $containers
    } else {
        Write-Host "❌ Nenhum container em execução ou Docker não responde" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ Erro ao verificar containers: $_" -ForegroundColor Red
}

Write-Host "`n🔍 Teste de Conectividade dos Serviços:" -ForegroundColor Yellow

$services = @(
    @{Name="API AlgaFood"; Url="http://localhost:8080/health"; Port=8080},
    @{Name="MailHog Web"; Url="http://localhost:8025"; Port=8025},
    @{Name="LocalStack"; Url="http://localhost:4566/health"; Port=4566}
)

$successCount = 0
foreach ($service in $services) {
    if (Test-Service -Name $service.Name -Url $service.Url -Port $service.Port) {
        $successCount++
    }
}

Write-Host "`n📊 Teste de Portas TCP:" -ForegroundColor Yellow

$ports = @(
    @{Port=8080; Service="API AlgaFood"},
    @{Port=13306; Service="MySQL"},
    @{Port=16379; Service="Redis"},
    @{Port=4566; Service="LocalStack"},
    @{Port=8025; Service="MailHog"},
    @{Port=1025; Service="MailHog SMTP"}
)

$openPorts = 0
foreach ($portInfo in $ports) {
    if (Test-Port -Port $portInfo.Port -Service $portInfo.Service) {
        $openPorts++
    }
}

Write-Host "`n📋 Logs da API (se disponível):" -ForegroundColor Yellow
try {
    $apiLogs = docker logs algafood-api --tail=5 2>$null
    if ($apiLogs) {
        Write-Host $apiLogs -ForegroundColor Gray
    } else {
        Write-Host "❌ Logs da API não disponíveis" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ Erro ao obter logs da API" -ForegroundColor Red
}

Write-Host "`n📊 Resumo:" -ForegroundColor Magenta
Write-Host "   Serviços Web respondendo: $successCount/3" -ForegroundColor White
Write-Host "   Portas TCP abertas: $openPorts/6" -ForegroundColor White

if ($successCount -ge 2 -and $openPorts -ge 4) {
    Write-Host "`n🎉 Infraestrutura funcionando bem!" -ForegroundColor Green
} elseif ($successCount -ge 1 -or $openPorts -ge 3) {
    Write-Host "`n⚠️  Infraestrutura parcialmente funcional" -ForegroundColor Yellow
} else {
    Write-Host "`n❌ Infraestrutura com problemas" -ForegroundColor Red
}

Write-Host "`n💡 Comandos úteis:" -ForegroundColor Magenta
Write-Host "   .\verify-infrastructure.ps1 -Quick    # Teste rápido de portas" -ForegroundColor Gray
Write-Host "   docker-compose ps                     # Status containers" -ForegroundColor Gray
Write-Host "   docker-compose logs -f algafood-api   # Logs da API" -ForegroundColor Gray
Write-Host "   docker-compose down && docker-compose up -d  # Restart completo" -ForegroundColor Gray
