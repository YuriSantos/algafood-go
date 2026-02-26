param(
    [switch]$Stats,
    [switch]$Logs,
    [switch]$MailHog
)

Write-Host "📧 LocalStack Email Checker" -ForegroundColor Green
Write-Host ("=" * 50) -ForegroundColor Gray

# Configurar AWS CLI
$env:AWS_ACCESS_KEY_ID = "test"
$env:AWS_SECRET_ACCESS_KEY = "test"
$env:AWS_DEFAULT_REGION = "us-east-1"

if ($Stats -or (-not $Logs -and -not $MailHog)) {
    Write-Host "`n📊 Estatísticas SES:" -ForegroundColor Yellow
    aws --endpoint-url=http://localhost:4566 ses get-send-statistics --output table
}

if ($Logs) {
    Write-Host "`n📋 Logs recentes:" -ForegroundColor Yellow
    docker logs algafood-localstack-1 --since="1h" | Select-String -Pattern "email|ses" | Select-Object -Last 10
}

if ($MailHog) {
    Write-Host "`n📧 Configurando MailHog:" -ForegroundColor Yellow
    docker run --rm -d --name mailhog -p 1025:1025 -p 8025:8025 mailhog/mailhog
    Write-Host "✅ MailHog iniciado! Acesse: http://localhost:8025" -ForegroundColor Green
}

Write-Host "`n🔗 Links úteis:" -ForegroundColor Magenta
Write-Host "• LocalStack: http://localhost:4566"
Write-Host "• MailHog: http://localhost:8025"
Write-Host "• Interface: .\scripts\email-viewer-fixed.html"

Write-Host "`n🚀 Status dos containers:" -ForegroundColor Yellow
docker ps --format "table {{.Names}}\t{{.Status}}" | Select-String "localstack|mailhog|algafood"
