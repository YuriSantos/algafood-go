# Script para verificar emails enviados no LocalStack
param(
    [switch]$ShowLogs,
    [switch]$ShowStats,
    [switch]$OpenBrowser
)

Write-Host "=== Verificando emails enviados no LocalStack SES ===" -ForegroundColor Green

# Configurar AWS CLI para usar LocalStack
$env:AWS_ACCESS_KEY_ID = "test"
$env:AWS_SECRET_ACCESS_KEY = "test"
$env:AWS_DEFAULT_REGION = "us-east-1"

if ($ShowStats -or (!$ShowLogs -and !$OpenBrowser)) {
    Write-Host "`n--- Estatísticas de envio SES ---" -ForegroundColor Yellow
    aws --endpoint-url=http://localhost:4566 ses get-send-statistics --query 'SendDataPoints[*].[Timestamp,DeliveryAttempts,Bounces,Complaints,Rejects]' --output table
}

if ($ShowLogs -or (!$ShowStats -and !$OpenBrowser)) {
    Write-Host "`n--- Verificando logs do LocalStack (últimos emails) ---" -ForegroundColor Yellow
    docker logs algafood-localstack-1 --since="1h" | Select-String -Pattern "(SES|email|SendEmail|Message)" | Select-Object -Last 20
}

Write-Host "`n--- Formas de visualizar emails enviados ---" -ForegroundColor Cyan
Write-Host "1. Interface Web LocalStack:" -ForegroundColor White
Write-Host "   http://localhost:4566" -ForegroundColor Gray

Write-Host "`n2. Verificar container de email (se estiver rodando):" -ForegroundColor White
Write-Host "   docker run --rm -it -p 8025:8025 mailhog/mailhog" -ForegroundColor Gray

Write-Host "`n3. Ver logs detalhados do LocalStack:" -ForegroundColor White
Write-Host "   docker logs algafood-localstack-1 -f | findstr -i ses" -ForegroundColor Gray

Write-Host "`n4. API direta para estatísticas:" -ForegroundColor White
Write-Host "   curl http://localhost:4566/_localstack/ses" -ForegroundColor Gray

Write-Host "`n5. Verificar todas as mensagens enviadas:" -ForegroundColor White
Write-Host "   aws --endpoint-url=http://localhost:4566 ses get-send-statistics" -ForegroundColor Gray

if ($OpenBrowser) {
    Write-Host "`n--- Abrindo interface web do LocalStack ---" -ForegroundColor Green
    Start-Process "http://localhost:4566"
}

Write-Host "`n--- Status dos serviços LocalStack ---" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:4566/health" -TimeoutSec 5 -ErrorAction Stop
    $response.PSObject.Properties | ForEach-Object {
        $status = if ($_.Value -eq "available") { "✓" } else { "✗" }
        Write-Host "$($_.Name): $status $($_.Value)" -ForegroundColor $(if ($_.Value -eq "available") {"Green"} else {"Red"})
    }
} catch {
    Write-Host "Erro ao verificar status do LocalStack: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n--- Como usar ---" -ForegroundColor Magenta
Write-Host ".\check-emails.ps1 -ShowLogs    # Mostra apenas logs"
Write-Host ".\check-emails.ps1 -ShowStats   # Mostra apenas estatísticas"
Write-Host ".\check-emails.ps1 -OpenBrowser # Abre interface web"
Write-Host ".\check-emails.ps1              # Mostra tudo"

