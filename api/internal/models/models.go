package models

import (
	"time"
)

// Modelos de dados para PostgreSQL

type Donation struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Amount          float64   `json:"amount"`
	DonorID         uint      `json:"donor_id"`
	NGOID           uint      `json:"ngo_id"`
	CreatedAt       time.Time `json:"created_at"`
	Status          string    `json:"status"`
	TransactionHash string    `json:"transaction_hash,omitempty"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
}

// NGO representa uma organização não governamental
type NGO struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	CNPJ          string    `json:"cnpj"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	LogoURL       string    `json:"logo_url"`
	DocumentsIPFS string    `json:"documents_ipfs,omitempty"`
	BlockchainRef string    `json:"blockchain_ref,omitempty"`
	ResponsibleID uint      `json:"responsible_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Estrutura para request de doação
type DonationRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	DonorID       uint    `json:"donor_id" binding:"required"`
	NGOID         uint    `json:"ngo_id" binding:"required"`
	DonorDocument string  `json:"donor_document,omitempty"` // CPF ou CNPJ do doador (será anonimizado)
}

// Estrutura para resposta de doação
type DonationResponse struct {
	ID              uint   `json:"id"`
	Status          string `json:"status"`
	PaymentURL      string `json:"payment_url,omitempty"`
	TransactionHash string `json:"transaction_hash,omitempty"`
}

// Mock de Payment Gateway
type PaymentGateway struct {
	ProviderName string `json:"provider_name"`
	PaymentURL   string `json:"payment_url"`
}

// ResourceUsage representa o uso dos recursos da doação
type ResourceUsage struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	DonationID  uint      `json:"donation_id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	ReceiptIPFS string    `json:"receipt_ipfs"`
	NGOName     string    `json:"ngo_name"`
	CreatedAt   time.Time `json:"created_at"`
}

// DonationReceipt representa o comprovante de doação
type DonationReceipt struct {
	ID              uint      `json:"id"`
	DonationID      uint      `json:"donation_id"`
	DonorName       string    `json:"donor_name"`
	DonorEmail      string    `json:"donor_email"`
	NGOName         string    `json:"ngo_name"`
	Amount          float64   `json:"amount"`
	Date            time.Time `json:"date"`
	TransactionHash string    `json:"transaction_hash"`
	IPFSHash        string    `json:"ipfs_hash"`
	PdfURL          string    `json:"pdf_url"`
}

// ImpactMetrics representa as métricas de impacto de doações
type ImpactMetrics struct {
	TotalDonated      float64 `json:"total_donated"`
	DonationsCount    int     `json:"donations_count"`
	NGOsSupported     int     `json:"ngos_supported"`
	PeopleHelped      int     `json:"people_helped"`
	MealsProvided     int     `json:"meals_provided"`
	MedicinesProvided int     `json:"medicines_provided"`
}

// DonorDashboard representa os dados para o dashboard do doador
type DonorDashboard struct {
	DonorID     uint          `json:"donor_id"`
	DonorName   string        `json:"donor_name"`
	Metrics     ImpactMetrics `json:"metrics"`
	Donations   []Donation    `json:"recent_donations"`
	UsagesCount int           `json:"usages_count"`
}

// ExpenseRequest representa o pedido de registro de um gasto por uma ONG
type ExpenseRequest struct {
	DonationID  uint    `json:"donation_id" binding:"required"`
	NGOID       uint    `json:"ngo_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"required"`
	Category    string  `json:"category" binding:"required"`
}

// Expense representa um gasto registrado por uma ONG
type Expense struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	DonationID    uint      `json:"donation_id"`
	NGOID         uint      `json:"ngo_id"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	ReceiptIPFS   string    `json:"receipt_ipfs,omitempty"`
	BlockchainRef string    `json:"blockchain_ref,omitempty"`
	Status        string    `json:"status"` // pendente, aprovado, rejeitado
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ExpenseResponse representa a resposta do registro de um gasto
type ExpenseResponse struct {
	ID            uint      `json:"id"`
	DonationID    uint      `json:"donation_id"`
	NGOID         uint      `json:"ngo_id"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	ReceiptIPFS   string    `json:"receipt_ipfs,omitempty"`
	BlockchainRef string    `json:"blockchain_ref,omitempty"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// Enum para categorias de gastos
var ExpenseCategories = []string{
	"Alimentação",
	"Saúde",
	"Educação",
	"Infraestrutura",
	"Administrativo",
	"Transporte",
	"Outros",
}

// NGORegistrationRequest representa uma solicitação de registro de ONG
type NGORegistrationRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description" binding:"required"`
	Category      string `json:"category" binding:"required"`
	CNPJ          string `json:"cnpj" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone" binding:"required"`
	Address       string `json:"address" binding:"required"`
	ResponsibleID uint   `json:"responsible_id" binding:"required"`
	LogoURL       string `json:"logo_url"`
}

// NGORegistrationStatus representa o status de um registro de ONG
type NGORegistrationStatus string

const (
	NGOStatusPending    NGORegistrationStatus = "pendente"
	NGOStatusValidating NGORegistrationStatus = "validando"
	NGOStatusRejected   NGORegistrationStatus = "rejeitado"
	NGOStatusApproved   NGORegistrationStatus = "aprovado"
)

// NGORegistration representa o processo de registro de uma ONG
type NGORegistration struct {
	ID                uint                  `json:"id" gorm:"primaryKey"`
	Name              string                `json:"name"`
	Description       string                `json:"description"`
	Category          string                `json:"category"`
	CNPJ              string                `json:"cnpj"`
	CNPJValid         bool                  `json:"cnpj_valid"`
	CNPJValidationMsg string                `json:"cnpj_validation_msg,omitempty"`
	Email             string                `json:"email"`
	Phone             string                `json:"phone"`
	Address           string                `json:"address"`
	ResponsibleID     uint                  `json:"responsible_id"`
	LogoURL           string                `json:"logo_url,omitempty"`
	DocumentsIPFS     string                `json:"documents_ipfs,omitempty"`
	BlockchainRef     string                `json:"blockchain_ref,omitempty"`
	Status            NGORegistrationStatus `json:"status"`
	AdminComments     string                `json:"admin_comments,omitempty"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

// NGODocumentUploadRequest representa uma solicitação de upload de documentos
type NGODocumentUploadRequest struct {
	RegistrationID uint   `json:"registration_id" binding:"required"`
	DocumentType   string `json:"document_type" binding:"required"`
}

// AuditRequest representa uma solicitação de auditoria
type AuditRequest struct {
	EntityType string `json:"entity_type" binding:"required"` // "ngo", "donation", "expense"
	EntityID   uint   `json:"entity_id" binding:"required"`
}

// AuditResult representa o resultado de uma auditoria
type AuditResult struct {
	EntityType       string    `json:"entity_type"`
	EntityID         uint      `json:"entity_id"`
	BlockchainValid  bool      `json:"blockchain_valid"`
	IPFSValid        bool      `json:"ipfs_valid"`
	BlockchainRef    string    `json:"blockchain_ref,omitempty"`
	IPFSRef          string    `json:"ipfs_ref,omitempty"`
	ValidationDate   time.Time `json:"validation_date"`
	ValidationErrors []string  `json:"validation_errors,omitempty"`
}

// AuditLog representa um registro de auditoria
type AuditLog struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	AdminID          uint      `json:"admin_id"`
	Action           string    `json:"action"`
	EntityType       string    `json:"entity_type"`
	EntityID         uint      `json:"entity_id"`
	PreviousState    string    `json:"previous_state,omitempty"`
	NewState         string    `json:"new_state,omitempty"`
	Comments         string    `json:"comments,omitempty"`
	BlockchainValid  bool      `json:"blockchain_valid,omitempty"`
	IPFSValid        bool      `json:"ipfs_valid,omitempty"`
	ValidationErrors []string  `json:"validation_errors,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// TransactionExplorerQuery representa uma consulta para o explorador de transações
type TransactionExplorerQuery struct {
	TransactionHash string    `json:"transaction_hash,omitempty"`
	NGOID           uint      `json:"ngo_id,omitempty"`
	StartDate       time.Time `json:"start_date,omitempty"`
	EndDate         time.Time `json:"end_date,omitempty"`
	Page            int       `json:"page,omitempty"`
	PageSize        int       `json:"page_size,omitempty"`
}

// TransactionExplorerResult representa o resultado de uma busca no explorador de transações
type TransactionExplorerResult struct {
	Donations []DonationDetails `json:"donations"`
	Total     int               `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
}

// DonationDetails representa os detalhes de uma doação para o explorador
type DonationDetails struct {
	ID              uint      `json:"id"`
	Amount          float64   `json:"amount"`
	DonorName       string    `json:"donor_name"`
	NGOName         string    `json:"ngo_name"`
	NGOCategory     string    `json:"ngo_category"`
	Date            time.Time `json:"date"`
	Status          string    `json:"status"`
	TransactionHash string    `json:"transaction_hash,omitempty"`
	HasReceipt      bool      `json:"has_receipt"`
	HasExpenses     bool      `json:"has_expenses"`
	ExpensesCount   int       `json:"expenses_count,omitempty"`
}

// GlobalDashboardData representa os dados para o dashboard global
type GlobalDashboardData struct {
	TotalDonated        float64                    `json:"total_donated"`
	TotalNGOs           int                        `json:"total_ngos"`
	TotalDonors         int                        `json:"total_donors"`
	TotalTransactions   int                        `json:"total_transactions"`
	DonationsByCategory []CategorySummary          `json:"donations_by_category"`
	MonthlyDonations    []MonthlyDonationData      `json:"monthly_donations"`
	TopNGOs             []NGODonationSummary       `json:"top_ngos"`
	GeographicalData    []GeographicalDonationData `json:"geographical_data,omitempty"`
	ImpactMetrics       GlobalImpactMetrics        `json:"impact_metrics"`
}

// CategorySummary representa o resumo de doações por categoria
type CategorySummary struct {
	Category    string  `json:"category"`
	TotalAmount float64 `json:"total_amount"`
	Count       int     `json:"count"`
	Percentage  float64 `json:"percentage"`
}

// MonthlyDonationData representa dados de doações por mês
type MonthlyDonationData struct {
	Month       string  `json:"month"`
	Year        int     `json:"year"`
	TotalAmount float64 `json:"total_amount"`
	Count       int     `json:"count"`
}

// NGODonationSummary representa o resumo de doações por ONG
type NGODonationSummary struct {
	NGOID       uint    `json:"ngo_id"`
	NGOName     string  `json:"ngo_name"`
	Category    string  `json:"category"`
	TotalAmount float64 `json:"total_amount"`
	Count       int     `json:"count"`
}

// GeographicalDonationData representa dados de doações por região geográfica
type GeographicalDonationData struct {
	Region      string  `json:"region"`
	TotalAmount float64 `json:"total_amount"`
	Count       int     `json:"count"`
}

// GlobalImpactMetrics representa métricas de impacto global
type GlobalImpactMetrics struct {
	PeopleHelped      int `json:"people_helped"`
	CommunitiesServed int `json:"communities_served"`
	ProjectsCompleted int `json:"projects_completed"`
	MealsProvided     int `json:"meals_provided"`
	MedicinesProvided int `json:"medicines_provided"`
	ChildrenEducated  int `json:"children_educated"`
	HousesBuilt       int `json:"houses_built"`
	EmergenciesServed int `json:"emergencies_served"`
}
