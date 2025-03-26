#!/bin/bash

echo "Instalando dependências do Swagger..."

# Limpar cache de dependências anteriores
go clean -modcache
rm -rf api/docs/*

# Instalar bibliotecas do Swagger com versões específicas
go get github.com/swaggo/swag/cmd/swag@v1.16.4
go get github.com/swaggo/gin-swagger@v1.6.0
go get github.com/swaggo/files@v1.0.1

# Instalar o executável swag globalmente
go install github.com/swaggo/swag/cmd/swag@v1.16.4

# Verificar a instalação do swag
echo "Verificando instalação do swag..."
SWAG_PATH=$(which swag || echo "$GOPATH/bin/swag")

if [ ! -f "$SWAG_PATH" ]; then
  SWAG_PATH="$HOME/go/bin/swag"
fi

if [ ! -f "$SWAG_PATH" ]; then
  echo "ERRO: swag não encontrado! Verificar instalação."
  exit 1
fi

echo "Usando swag em: $SWAG_PATH"

# Criar diretório de documentação se não existir
mkdir -p api/docs

# Gerar a documentação
echo "Gerando documentação Swagger..."
cd api && "$SWAG_PATH" init -g cmd/main.go -o docs

echo "Organizando dependências..."
cd .. && go mod tidy

echo "Documentação gerada com sucesso!"
echo "Acesse a documentação em http://localhost:8080/swagger/index.html após iniciar o servidor" 