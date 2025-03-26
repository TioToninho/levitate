#!/bin/bash

# Script para gerar chaves necessárias

# Exemplo: Geração de chave JWT
openssl genrsa -out jwt_private.key 2048
openssl rsa -in jwt_private.key -pubout -out jwt_public.key 