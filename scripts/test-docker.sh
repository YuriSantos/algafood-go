#!/bin/bash

echo "🔧 Teste de Infraestrutura Docker"
echo "================================="

echo "📋 Verificando Docker..."
docker --version
docker info

echo ""
echo "🐳 Testando container simples..."
docker run --rm hello-world

echo ""
echo "📦 Listando containers em execução..."
docker ps

echo ""
echo "📋 Listando todos os containers..."
docker ps -a

echo ""
echo "🌐 Testando conectividade de rede..."
docker network ls

echo ""
echo "📊 Uso de recursos..."
docker system df

echo ""
echo "✅ Teste concluído!"
