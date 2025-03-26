<div align="center">
<pre>
 ██╗     ███████╗██╗   ██╗██╗████████╗ █████╗ ████████╗███████╗
 ██║     ██╔════╝██║   ██║██║╚══██╔══╝██╔══██╗╚══██╔══╝██╔════╝
 ██║     █████╗  ██║   ██║██║   ██║   ███████║   ██║   █████╗  
 ██║     ██╔══╝  ╚██╗ ██╔╝██║   ██║   ██╔══██║   ██║   ██╔══╝  
 ███████╗███████╗ ╚████╔╝ ██║   ██║   ██║  ██║   ██║   ███████╗
 ╚══════╝╚══════╝  ╚═══╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝   ╚═╝   ╚══════╝
-------------------------------------------------------------------
Blockchain-based donation tracking platform
</pre>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org/)
[![Framework](https://img.shields.io/badge/Framework-Gin-green.svg)](https://gin-gonic.com/)
[![Status](https://img.shields.io/badge/Status-Work%20In%20Progress-orange.svg)]()

</div>

# Levitate API Documentation

> [!NOTE]  
> **Work in Progress**: Este projeto está atualmente em desenvolvimento ativo. Recursos e endpoints podem mudar conforme o desenvolvimento continua.

Levitate is a blockchain-based donation tracking platform that enables complete transparency between donors and NGOs. The platform allows donors to track how their donations are used and enables NGOs to provide proof of resource allocation.

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Installation](#installation)
- [Data Storage](#data-storage)
- [Security Features](#security-features)
- [API Endpoints](#api-endpoints)
  - [Health Check](#health-check)
  - [NGOs](#ngos)
  - [Donations](#donations)
  - [Expenses](#expenses)
  - [Transaction Explorer](#transaction-explorer)
  - [Global Dashboard](#global-dashboard)
  - [Transparency](#transparency)
  - [Admin Functions](#admin-functions)
- [Models](#models)
- [Technical Challenges](#technical-challenges)

## Overview

Levitate connects donors with NGOs in a transparent ecosystem where all transactions are recorded on a blockchain. Each donation is tracked from the moment it's made until it's fully utilized, with complete audit capabilities.

## Key Features

- **Complete Donation Traceability**: From donor to final beneficiary
- **Blockchain Verification**: Immutable records of all transactions
- **IPFS Document Storage**: Decentralized storage for receipts and proofs
- **Data Anonymization**: Privacy-preserving hashed personal data
- **Transaction Explorer**: Public search engine for all donations
- **Global Dashboard**: Visual analytics of impact and distribution
- **Administrative Tools**: NGO validation and auditing capabilities

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/your-organization/levitate.git
   cd levitate
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Configure environment variables:
   ```
   cp configs/.env.example configs/.env
   # Edit the .env file with your settings
   ```

4. Run the application:
   ```
   go run api/cmd/main.go
   ```

## Data Storage

| Data Type | Storage Location | Example |
|-----------|------------------|---------|
| NGO Metadata | PostgreSQL | Name, hashed CNPJ, description |
| Donation Transactions | Blockchain | Amount, recipient NGO, receipt hash |
| Receipts & Documents | IPFS | Unique CID (e.g., QmXYZ...) |
| Metrics Cache | Redis | Today's total, top 5 NGOs |

## Security Features

- **Data Anonymization**: CPF/CNPJ are hashed (SHA-256) before storage
- **Input Validation**: Checks for negative values, non-existent NGOs, and data format
- **Authentication**: JWT for administrators and NGOs
- **Data Protection**: All endpoints use HTTPS and rate limiting
- **Headers Security**: HSTS, CSP, XSS protection headers

## API Endpoints

### Health Check

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| GET | `/health` | Check API status | None |

**Example Request:**
```
GET /health
```

**Example Response:**
```json
{
  "status": "online",
  "version": "1.0.0",
  "timestamp": "2025-03-26T01:21:55.123Z",
  "uptime": "3h24m12s"
}
```

### NGOs

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| GET | `/ngos` | List all NGOs | None |
| GET | `/ngos/:id` | Get NGO details | None |

**Example Request:**
```
GET /ngos/2
```

**Example Response:**
```json
{
  "data": {
    "id": 2,
    "name": "Saúde para Todos",
    "description": "Fornecimento de medicamentos e atendimento médico gratuito",
    "category": "Saúde",
    "email": "contato@saudeparatodos.org",
    "phone": "+55 11 3333-4444",
    "address": "Av. Paulista, 1500, São Paulo - SP",
    "logo_url": "https://example.com/logo2.png"
  }
}
```

### Donations

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| POST | `/donations` | Create a new donation | None |
| POST | `/donations/:id/confirm-payment` | Confirm payment | None |
| GET | `/donations/:id/receipt` | Get donation receipt | None |
| GET | `/donations/:id/usages` | Get resource usage details | None |
| GET | `/donors/:id/donations` | List donor's donations | None |
| GET | `/donors/:id/dashboard` | Get donor's dashboard | None |

**Example Request:**
```
POST /donations
Content-Type: application/json

{
  "amount": 100.00,
  "donor_id": 1,
  "ngo_id": 2,
  "donor_document": "123.456.789-00"
}
```

**Example Response:**
```json
{
  "data": {
    "id": 42,
    "status": "pending",
    "payment_url": "https://payment-gateway-mock.com/pay?donationId=42&amount=100.00"
  }
}
```

### Expenses

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| POST | `/expenses` | Register an expense | None |
| POST | `/expenses/:id/receipt` | Upload expense receipt | None |
| GET | `/expenses/donation/:donationId` | Get expenses by donation | None |
| GET | `/expenses/ngo/:ngoId` | Get expenses by NGO | None |

**Example Request:**
```
POST /expenses
Content-Type: application/json

{
  "donation_id": 42,
  "ngo_id": 2,
  "amount": 50.00,
  "description": "Purchase of medicines for community clinic",
  "category": "Saúde"
}
```

**Example Response:**
```json
{
  "id": 15,
  "donation_id": 42,
  "ngo_id": 2,
  "amount": 50.00,
  "description": "Purchase of medicines for community clinic",
  "category": "Saúde",
  "status": "pending",
  "created_at": "2025-03-26T14:22:36Z"
}
```

### Transaction Explorer

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| GET | `/explorer/search` | Search donations with filters | None |
| GET | `/explorer/donations/hash/:hash` | Get donation by transaction hash | None |
| GET | `/explorer/donations/:id` | Get donation by ID | None |
| GET | `/explorer/donations/ngo/:ngo_id` | Get donations by NGO | None |
| GET | `/explorer/donations/recent` | Get recent donations | None |

**Example Request:**
```
GET /explorer/search?ngo_id=2&start_date=2025-01-01&end_date=2025-03-31&page=1&page_size=10
```

**Example Response:**
```json
{
  "donations": [
    {
      "id": 42,
      "amount": 100.00,
      "donor_name": "João Silva",
      "ngo_name": "Saúde para Todos",
      "ngo_category": "Saúde",
      "date": "2025-03-20T10:30:00Z",
      "status": "completed",
      "transaction_hash": "0x7a69f3ee23e5e4d888c3ef3b2e91c9077064a2f59ee8f937776b18f7602d13km",
      "has_receipt": true,
      "has_expenses": true,
      "expenses_count": 2
    }
  ],
  "total": 5,
  "page": 1,
  "page_size": 10
}
```

### Global Dashboard

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| GET | `/dashboard/global` | Get global dashboard data | None |
| GET | `/dashboard/by-date-range` | Get dashboard for date range | None |
| GET | `/dashboard/by-category/:category` | Get dashboard for category | None |

**Example Request:**
```
GET /dashboard/global
```

**Example Response:**
```json
{
  "total_amount": 250000.00,
  "ngos_count": 15,
  "donors_count": 1200,
  "transactions_count": 1850,
  "donations_by_category": [
    {
      "category": "Saúde",
      "amount": 85000.00,
      "count": 650,
      "percentage": 34
    },
    {
      "category": "Educação",
      "amount": 75000.00,
      "count": 580,
      "percentage": 30
    }
  ],
  "monthly_donations": [
    {
      "month": "January",
      "amount": 45000.00,
      "count": 320
    }
  ],
  "top_ngos": [
    {
      "ngo_id": 2,
      "ngo_name": "Saúde para Todos",
      "amount": 42000.00,
      "count": 320
    }
  ],
  "geographical_data": {
    "regions": [
      {
        "name": "Southeast",
        "amount": 120000.00,
        "percentage": 48
      }
    ]
  },
  "impact_metrics": {
    "people_helped": 25000,
    "communities_served": 120,
    "meals_provided": 75000,
    "medical_services": 12000,
    "educational_resources": 15000
  }
}
```

### Transparency

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| GET | `/transparency` | Get public dashboard | None |
| GET | `/transparency/donations` | Get public donations | None |
| GET | `/transparency/expenses` | Get public expenses | None |
| GET | `/transparency/ngos` | Get NGOs summary | None |
| GET | `/transparency/ngos/:id` | Get specific NGO summary | None |
| GET | `/transparency/ngos/:id/donations` | Get NGO donations | None |
| GET | `/transparency/ngos/:id/expenses` | Get NGO expenses | None |

**Example Request:**
```
GET /transparency/ngos/2
```

**Example Response:**
```json
{
  "id": 2,
  "name": "Saúde para Todos",
  "category": "Saúde",
  "total_donations": 42000.00,
  "total_expenses": 35000.00,
  "donation_count": 320,
  "expense_count": 48,
  "impact_summary": {
    "people_helped": 5000,
    "medical_services": 12000
  }
}
```

### Admin Functions

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| POST | `/admin/ngos/register` | Register new NGO | Admin |
| POST | `/admin/ngos/registration/:id/validate-cnpj` | Validate CNPJ | Admin |
| POST | `/admin/ngos/registration/:id/upload-documents` | Upload NGO documents | Admin |
| POST | `/admin/ngos/registration/:id/approve` | Approve NGO | Admin |
| POST | `/admin/ngos/registration/:id/reject` | Reject NGO | Admin |
| GET | `/admin/ngos/registrations` | List NGO registrations | Admin |
| GET | `/admin/ngos/registrations/:id` | Get registration details | Admin |
| GET | `/admin/ngos/registrations/by-cnpj` | Search registrations by CNPJ | Admin |
| POST | `/admin/audit` | Audit entity | Admin |
| GET | `/admin/audit/logs` | Get audit logs | Admin |

**Example Request:**
```
POST /admin/ngos/register
Content-Type: application/json
X-Admin-ID: admin-123

{
  "name": "Educação é Futuro",
  "description": "Apoio educacional para crianças de baixa renda",
  "category": "Educação",
  "cnpj": "12.345.678/0001-90",
  "email": "contato@educacaoefuturo.org",
  "phone": "+55 11 5555-6666",
  "address": "Rua Augusta, 500, São Paulo - SP",
  "responsible_id": 3,
  "logo_url": "https://example.com/logo3.png"
}
```

**Example Response:**
```json
{
  "id": 5,
  "name": "Educação é Futuro",
  "description": "Apoio educacional para crianças de baixa renda",
  "category": "Educação",
  "cnpj": "12.345.678/0001-90",
  "cnpj_valid": false,
  "email": "contato@educacaoefuturo.org",
  "phone": "+55 11 5555-6666",
  "address": "Rua Augusta, 500, São Paulo - SP",
  "responsible_id": 3,
  "logo_url": "https://example.com/logo3.png",
  "status": "pendente",
  "created_at": "2025-03-26T14:30:00Z"
}
```

## Models

### Donation Flow

1. **NGO Registration (via admin)**
   - NGO data stored in PostgreSQL and its ID registered on blockchain

2. **Donation Simulation**
   - Donor enters amount and selects NGO → API simulates payment and generates receipt

3. **Blockchain Registration**
   - IPFS receipt hash and donation metadata added to blockchain

4. **Resource Usage by NGO**
   - NGO records expenses with new receipts, linked to original transaction

5. **Public Audit**
   - Anyone can verify hashes on blockchain and compare with IPFS files

### Core Models

- **NGO**: Organization that receives donations
- **Donation**: Transaction record between donor and NGO
- **Expense**: Record of how donation funds were used
- **Receipt**: Proof of donation or expense
- **ResourceUsage**: Record of resource allocation
- **Audit**: Verification of blockchain and IPFS data integrity

## Technical Challenges

| Challenge | Solution |
|-----------|----------|
| Blockchain Performance | Proof-of-Authority consensus algorithm |
| File Storage | IPFS + replication to prevent data loss |
| Anonymity vs. Transparency | Hashing sensitive data + public metadata |
| Frontend/Blockchain Integration | Well-documented REST API + WebSocket for real-time updates |

## Contribution

Para contribuir com este projeto:

1. Faça um fork do repositório
2. Crie seu branch de funcionalidade: `git checkout -b feature/funcionalidade-incrivel`
3. Faça commit das suas alterações: `git commit -m 'Adiciona alguma funcionalidade incrível'`
4. Faça push para o branch: `git push origin feature/funcionalidade-incrivel`
5. Abra um Pull Request

### Ambiente de Desenvolvimento

- Go 1.18+
- Docker e Docker Compose para ambiente de desenvolvimento local
- PostgreSQL para armazenamento de dados
- IPFS para armazenamento de documentos
- Blockchain privada para registros imutáveis

### Considerações sobre o .gitignore

O projeto inclui um arquivo `.gitignore` configurado para projetos Go que exclui:

- Binários compilados
- Arquivos de ambiente (.env)
- Arquivos de chaves privadas e certificados
- Arquivos temporários e de cache
- Arquivos específicos de IDEs

**Importante**: Nunca comite chaves privadas, tokens JWT ou credenciais de banco de dados no repositório. Utilize o arquivo `.env.example` como modelo e crie seu próprio `.env` local.

## License

Distributed under the MIT License. See `LICENSE` for more information. 