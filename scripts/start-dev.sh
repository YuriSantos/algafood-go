#!/bin/bash

# Script para desenvolvimento - inicia todos os serviços
echo "🚀 Iniciando infraestrutura AlgaFood..."

# Verificar se Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker não está rodando. Por favor, inicie o Docker Desktop."
    exit 1
fi

# Parar containers existentes
echo "🛑 Parando containers existentes..."
docker-compose down

# Construir e iniciar serviços
echo "🔨 Construindo e iniciando serviços..."
docker-compose up --build -d

# Aguardar serviços estarem prontos
echo "⏳ Aguardando serviços iniciarem..."

# Verificar MySQL
echo "🔍 Verificando MySQL..."
until docker exec algafood-mysql mysqladmin ping -h localhost --silent; do
    echo "   Aguardando MySQL..."
    sleep 2
done
echo "✅ MySQL pronto!"

# Verificar Redis
echo "🔍 Verificando Redis..."
until docker exec algafood-redis redis-cli ping > /dev/null 2>&1; do
    echo "   Aguardando Redis..."
    sleep 2
done
echo "✅ Redis pronto!"

# Verificar LocalStack
echo "🔍 Verificando LocalStack..."
until curl -s http://localhost:4566/health > /dev/null 2>&1; do
    echo "   Aguardando LocalStack..."
    sleep 2
done
echo "✅ LocalStack pronto!"

# Executar migrações (se necessário)
echo "🗄️ Executando migrações..."
# docker exec algafood-api ./main migrate

# Mostrar status dos serviços
echo "📊 Status dos serviços:"
docker-compose ps

echo ""
echo "🎉 Infraestrutura AlgaFood iniciada com sucesso!"
echo ""
echo "📋 Serviços disponíveis:"
echo "   🌐 API AlgaFood:     http://localhost:8080"
echo "   🌐 Nginx (Proxy):   http://localhost:80"
echo "   📧 MailHog:          http://localhost:8025"
echo "   ☁️  LocalStack:      http://localhost:4566"
echo "   🗄️ MySQL:            localhost:13306"
echo "   🔴 Redis:            localhost:16379"
echo ""
echo "🔧 Comandos úteis:"
echo "   docker-compose logs -f algafood-api    # Ver logs da API"
echo "   docker-compose logs -f localstack      # Ver logs do LocalStack"
echo "   docker-compose down                    # Parar todos os serviços"
echo "   docker-compose up -d                   # Reiniciar serviços"
echo ""
echo "📧 Para verificar emails:"
echo "   .\scripts\email-checker-simple.ps1"
echo "   start .\scripts\email-viewer-fixed.html"
