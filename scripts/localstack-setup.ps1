# Run this script after starting LocalStack
# Execute este script após iniciar o LocalStack

$ENDPOINT = "http://localhost:4566"
$REGION = "us-east-1"
$ACCOUNT_ID = "000000000000"
# 1. Create DLQ first
Write-Host "`n1. Creating DLQ: $DLQ_NAME" -ForegroundColor Yellow
# 2. Create SQS Queue with DLQ
Write-Host "`n2. Creating SQS Queue: $QUEUE_NAME" -ForegroundColor Yellow

Write-Host "Configurando recursos LocalStack para AlgaFood..." -ForegroundColor Green

# 1. Criar DLQ primeiro
Write-Host "`n1. Criando DLQ: $DLQ_NAME" -ForegroundColor Yellow
aws --endpoint-url=$ENDPOINT sqs create-queue --queue-name $DLQ_NAME --region $REGION 2>$null

# 2. Criar fila SQS com DLQ
Write-Host "`n2. Criando Fila SQS: $QUEUE_NAME" -ForegroundColor Yellow
$DLQ_ARN = "arn:aws:sqs:${REGION}:${ACCOUNT_ID}:${DLQ_NAME}"
$REDRIVE_POLICY = "{`"deadLetterTargetArn`":`"$DLQ_ARN`",`"maxReceiveCount`":3}"

aws --endpoint-url=$ENDPOINT sqs create-queue `
    --queue-name $QUEUE_NAME `
    --attributes "RedrivePolicy=$REDRIVE_POLICY" `
    --region $REGION 2>$null

# Get the actual queue URL from LocalStack
$QUEUE_URL_RESULT = aws --endpoint-url=$ENDPOINT sqs get-queue-url --queue-name $QUEUE_NAME --region $REGION 2>&1 | ConvertFrom-Json
$QUEUE_URL = $QUEUE_URL_RESULT.QueueUrl
$QUEUE_ARN = "arn:aws:sqs:${REGION}:${ACCOUNT_ID}:${QUEUE_NAME}"

Write-Host "   URL da Fila: $QUEUE_URL"
Write-Host "   ARN da Fila: $QUEUE_ARN"

# 3. Criar Event Bus do EventBridge
Write-Host "`n3. Criando Event Bus do EventBridge: $EVENT_BUS_NAME" -ForegroundColor Yellow
aws --endpoint-url=$ENDPOINT events create-event-bus --name $EVENT_BUS_NAME --region $REGION 2>$null

# 4. Criar regra do EventBridge para encaminhar eventos para SQS
Write-Host "`n4. Criando regra do EventBridge para encaminhar eventos para SQS" -ForegroundColor Yellow
$RULE_NAME = "algafood-to-sqs-rule"

# Criar regra que corresponde a todos os eventos da origem algafood-api
$EVENT_PATTERN = '{"source":["algafood-api"]}'
aws --endpoint-url=$ENDPOINT events put-rule `
    --name $RULE_NAME `
    --event-bus-name $EVENT_BUS_NAME `
    --event-pattern $EVENT_PATTERN `
    --state ENABLED `
    --region $REGION

# 5. Adicionar SQS como alvo da regra
Write-Host "`n5. Adicionando Fila SQS como alvo da regra" -ForegroundColor Yellow
aws --endpoint-url=$ENDPOINT events put-targets `
    --rule $RULE_NAME `
    --event-bus-name $EVENT_BUS_NAME `
    --targets "Id=sqs-target,Arn=$QUEUE_ARN" `
    --region $REGION

# 6. Definir política da fila SQS para permitir EventBridge enviar mensagens
Write-Host "`n6. Definindo política da fila SQS para EventBridge" -ForegroundColor Yellow
$POLICY = @"
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": "*",
            "Action": "sqs:SendMessage",
            "Resource": "$QUEUE_ARN"
        }
    ]
}
"@
$POLICY_ESCAPED = $POLICY -replace "`n", "" -replace "`r", "" -replace "  +", " "
aws --endpoint-url=$ENDPOINT sqs set-queue-attributes `
    --queue-url $QUEUE_URL `
    --attributes "Policy=$POLICY_ESCAPED" `
    --region $REGION 2>$null

# 7. Verify SES email identity (for sending emails)
Write-Host "`n7. Verifying SES email identity" -ForegroundColor Yellow
# 7. Verify SES email identities (for sending emails)
Write-Host "`n7. Verificando identidades de email no SES" -ForegroundColor Yellow

# Verifica o email remetente
# 8. List resources to verify
Write-Host "`n8. Verifying resources..." -ForegroundColor Yellow
    Write-Host "   Verificando destinatário: $EMAIL" -ForegroundColor Cyan
    aws --endpoint-url=$ENDPOINT ses verify-email-identity --email-address $EMAIL --region $REGION 2>$null
}

# Lista identidades verificadas
# 3. Create EventBridge Event Bus
# 4. Create EventBridge Rule to forward all events to SQS
# Create rule that matches all events from algafood-api source
aws --endpoint-url=$ENDPOINT ses list-identities --region $REGION

Write-Host "`n   SQS Queues:" -ForegroundColor Cyan
aws --endpoint-url=$ENDPOINT events list-event-buses --region $REGION

Write-Host "`n   Regras do EventBridge:" -ForegroundColor Cyan
aws --endpoint-url=$ENDPOINT events list-rules --event-bus-name $EVENT_BUS_NAME --region $REGION

Write-Host "`n   EventBridge Targets for rule ${RULE_NAME}:" -ForegroundColor Cyan
Write-Host "`n   Alvos da regra ${RULE_NAME}:" -ForegroundColor Cyan
Write-Host "`n   EventBridge Event Buses:" -ForegroundColor Cyan

Write-Host "`n`nSetup complete!" -ForegroundColor Green
Write-Host "Queue URL for config.yaml: $QUEUE_URL" -ForegroundColor Cyan
Write-Host "`nYou can now start the AlgaFood API and test the event flow." -ForegroundColor Green
Write-Host "`n`nConfiguração concluída!" -ForegroundColor Green
Write-Host "`n   EventBridge Rules:" -ForegroundColor Cyan
Write-Host "`nVocê pode agora iniciar a API AlgaFood e testar o fluxo de eventos." -ForegroundColor Green


