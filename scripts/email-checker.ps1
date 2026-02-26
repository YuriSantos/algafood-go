param(
    [switch]$Stats,
    [switch]$Logs,
    [switch]$MailHog
)

Write-Host "📧 LocalStack Email Checker" -ForegroundColor Green
Write-Host "=" * 50 -ForegroundColor Gray

# Configurar AWS CLI
$env:AWS_ACCESS_KEY_ID = "test"
$env:AWS_SECRET_ACCESS_KEY = "test"
$env:AWS_DEFAULT_REGION = "us-east-1"

if ($Stats -or (-not $Logs -and -not $MailHog)) {
    Write-Host "`n📊 Estatísticas SES:" -ForegroundColor Yellow
    try {
        aws --endpoint-url=http://localhost:4566 ses get-send-statistics --output table
    } catch {
        Write-Host "❌ Erro ao obter estatísticas. Verifique se LocalStack está rodando." -ForegroundColor Red
    }
}

if ($Logs -or (-not $Stats -and -not $MailHog)) {
    Write-Host "`n📋 Logs recentes (última hora):" -ForegroundColor Yellow
    try {
        $logs = docker logs algafood-localstack-1 --since="1h" 2>$null
        if ($logs) {
            $emailLogs = $logs | Select-String -Pattern "(email|ses|SendEmail)" -CaseSensitive:$false
            if ($emailLogs) {
                $emailLogs | Select-Object -Last 10
            } else {
                Write-Host "Nenhum log de email encontrado na última hora." -ForegroundColor Gray
            }
        } else {
            Write-Host "Nenhum log encontrado. Verifique se o container está rodando." -ForegroundColor Gray
        }
    } catch {
        Write-Host "❌ Erro ao verificar logs." -ForegroundColor Red
    }
}

if ($MailHog) {
    Write-Host "`n📧 Configurando MailHog para capturar emails:" -ForegroundColor Yellow
    try {
        $mailhogRunning = docker ps --filter "name=mailhog" --format "{{.Names}}" 2>$null
        if ($mailhogRunning -eq "mailhog") {
            Write-Host "✅ MailHog já está rodando!" -ForegroundColor Green
        } else {
            Write-Host "🚀 Iniciando MailHog..." -ForegroundColor Cyan
            docker run --rm -d --name mailhog -p 1025:1025 -p 8025:8025 mailhog/mailhog
            Start-Sleep 2
            Write-Host "✅ MailHog iniciado com sucesso!" -ForegroundColor Green
        }
        Write-Host "🌐 Interface web: http://localhost:8025" -ForegroundColor Cyan
        Write-Host "📧 SMTP: localhost:1025" -ForegroundColor Cyan
    } catch {
        Write-Host "❌ Erro ao configurar MailHog: $_" -ForegroundColor Red
    }
}

Write-Host "`n🔗 Links úteis:" -ForegroundColor Magenta
Write-Host "• LocalStack: http://localhost:4566" -ForegroundColor Gray
Write-Host "• MailHog: http://localhost:8025" -ForegroundColor Gray
Write-Host "• Interface HTML: .\scripts\email-viewer-fixed.html" -ForegroundColor Gray

Write-Host "`n💡 Comandos:" -ForegroundColor Magenta
Write-Host "• .\email-checker.ps1 -Stats     # Apenas estatísticas" -ForegroundColor Gray
Write-Host "• .\email-checker.ps1 -Logs      # Apenas logs" -ForegroundColor Gray
Write-Host "• .\email-checker.ps1 -MailHog   # Configurar MailHog" -ForegroundColor Gray
Write-Host "• .\email-checker.ps1            # Mostrar tudo" -ForegroundColor Gray

Write-Host "`n🚀 Status dos containers:" -ForegroundColor Yellow
try {
    $containers = docker ps --format "table {{.Names}}\t{{.Status}}" 2>$null
    if ($containers) {
        $containers | Where-Object { $_ -like "*localstack*" -or $_ -like "*mailhog*" -or $_ -like "*algafood*" }
    } else {
        Write-Host "Nenhum container encontrado." -ForegroundColor Gray
    }
} catch {
    Write-Host "❌ Erro ao verificar containers." -ForegroundColor Red
}
