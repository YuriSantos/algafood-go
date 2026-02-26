#!/bin/bash

echo "Inicializando LocalStack..."
sleep 5

export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1

# Criar buckets S3
aws --endpoint-url=http://localhost:4566 s3 mb s3://algafood-files 2>/dev/null || true
aws --endpoint-url=http://localhost:4566 s3 mb s3://algafood-fotos-produtos 2>/dev/null || true

# Configurar SES
aws --endpoint-url=http://localhost:4566 ses verify-email-identity --email-address teste@algafood.com.br 2>/dev/null || true
aws --endpoint-url=http://localhost:4566 ses verify-email-identity --email-address admin@algafood.com.br 2>/dev/null || true

# Criar filas SQS
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name algafood-pedido-status 2>/dev/null || true

# Criar event bus
aws --endpoint-url=http://localhost:4566 events create-event-bus --name algafood-event-bus 2>/dev/null || true

echo "LocalStack inicializado!"


