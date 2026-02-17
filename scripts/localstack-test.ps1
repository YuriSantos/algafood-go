# LocalStack Test Script for AlgaFood
# Use this script to test the event flow

$ENDPOINT = "http://localhost:4566"
$REGION = "us-east-1"
$EVENT_BUS_NAME = "algafood-event-bus"
$QUEUE_NAME = "algafood-pedido-status"

Write-Host "Testing AlgaFood LocalStack setup..." -ForegroundColor Green

# Get queue URL
Write-Host "`n1. Getting Queue URL..." -ForegroundColor Yellow
$QUEUE_URL_RESULT = aws --endpoint-url=$ENDPOINT sqs get-queue-url --queue-name $QUEUE_NAME --region $REGION 2>&1 | ConvertFrom-Json
$QUEUE_URL = $QUEUE_URL_RESULT.QueueUrl
Write-Host "   Queue URL: $QUEUE_URL" -ForegroundColor Cyan

# Check queue message count before
Write-Host "`n2. Checking queue message count before test..." -ForegroundColor Yellow
aws --endpoint-url=$ENDPOINT sqs get-queue-attributes `
    --queue-url $QUEUE_URL `
    --attribute-names ApproximateNumberOfMessages ApproximateNumberOfMessagesNotVisible `
    --region $REGION

# Send message directly to SQS (simulating what EventBridge would do)
Write-Host "`n3. Sending test message directly to SQS..." -ForegroundColor Yellow

$testId = Get-Random -Maximum 9999
$timestamp = (Get-Date).ToString("yyyy-MM-ddTHH:mm:ssZ")

# Create the message in EventBridge format
$sqsMessageFile = [System.IO.Path]::GetTempFileName()
$sqsMessageContent = @"
{
    "version": "0",
    "id": "test-$testId",
    "detail-type": "PedidoConfirmado",
    "source": "algafood-api",
    "account": "000000000000",
    "time": "$timestamp",
    "region": "$REGION",
    "detail": {
        "timestamp": "$timestamp",
        "pedidoCodigo": "TEST-$testId",
        "clienteId": 1,
        "clienteNome": "Cliente Teste",
        "clienteEmail": "teste@algafood.com.br",
        "restauranteId": 1,
        "restauranteNome": "Restaurante Teste",
        "valorTotal": 99.9,
        "dataConfirmacao": "$timestamp"
    }
}
"@
$sqsMessageContent | Out-File -FilePath $sqsMessageFile -Encoding utf8
Write-Host "   Message file: $sqsMessageFile" -ForegroundColor Gray

aws --endpoint-url=$ENDPOINT sqs send-message `
    --queue-url $QUEUE_URL `
    --message-body "file://$sqsMessageFile" `
    --region $REGION

Remove-Item $sqsMessageFile

# Wait a moment for message to be available
Start-Sleep -Seconds 2

# Check queue message count after
Write-Host "`n4. Checking queue message count after test..." -ForegroundColor Yellow
aws --endpoint-url=$ENDPOINT sqs get-queue-attributes `
    --queue-url $QUEUE_URL `
    --attribute-names ApproximateNumberOfMessages ApproximateNumberOfMessagesNotVisible `
    --region $REGION

# Try to receive message from queue (only if app is not running)
Write-Host "`n5. Attempting to peek at message in queue (won't delete)..." -ForegroundColor Yellow
Write-Host "   Note: If your app is running, it may have already consumed the message." -ForegroundColor Gray
aws --endpoint-url=$ENDPOINT sqs receive-message `
    --queue-url $QUEUE_URL `
    --max-number-of-messages 1 `
    --visibility-timeout 0 `
    --wait-time-seconds 2 `
    --region $REGION

Write-Host "`n`nTest complete!" -ForegroundColor Green
Write-Host "If your app is running, check the logs for message processing." -ForegroundColor Cyan



