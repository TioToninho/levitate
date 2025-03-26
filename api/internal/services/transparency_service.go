package services

import (
	"fmt"
	"sort"
	"time"
)

// TransparencyService gerencia operações relacionadas à transparência pública
type TransparencyService struct {
	donationService *DonationService
	expenseService  *ExpenseService
}

// TransparencyDonation representa uma doação para exibição pública
type TransparencyDonation struct {
	ID              uint      `json:"id"`
	Amount          float64   `json:"amount"`
	NGOName         string    `json:"ngo_name"`
	NGOCategory     string    `json:"ngo_category"`
	Date            time.Time `json:"date"`
	Status          string    `json:"status"`
	TransactionHash string    `json:"transaction_hash,omitempty"`
}

// TransparencyExpense representa uma despesa para exibição pública
type TransparencyExpense struct {
	ID            uint      `json:"id"`
	DonationID    uint      `json:"donation_id"`
	NGOName       string    `json:"ngo_name"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Date          time.Time `json:"date"`
	ReceiptIPFS   string    `json:"receipt_ipfs,omitempty"`
	BlockchainRef string    `json:"blockchain_ref,omitempty"`
	Status        string    `json:"status"`
}

// TransparencyNGOSummary representa o resumo de uma ONG para transparência
type TransparencyNGOSummary struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	Category         string  `json:"category"`
	TotalReceived    float64 `json:"total_received"`
	TotalSpent       float64 `json:"total_spent"`
	DonationsCount   int     `json:"donations_count"`
	ExpensesCount    int     `json:"expenses_count"`
	AvailableBalance float64 `json:"available_balance"`
}

// TransparencyDashboard representa o resumo geral de transparência
type TransparencyDashboard struct {
	TotalDonations  float64                  `json:"total_donations"`
	TotalExpenses   float64                  `json:"total_expenses"`
	DonationsCount  int                      `json:"donations_count"`
	ExpensesCount   int                      `json:"expenses_count"`
	NGOsCount       int                      `json:"ngos_count"`
	RecentDonations []TransparencyDonation   `json:"recent_donations"`
	RecentExpenses  []TransparencyExpense    `json:"recent_expenses"`
	NGOsSummary     []TransparencyNGOSummary `json:"ngos_summary"`
}

// NewTransparencyService cria uma nova instância do serviço de transparência
func NewTransparencyService(donationSvc *DonationService, expenseSvc *ExpenseService) *TransparencyService {
	return &TransparencyService{
		donationService: donationSvc,
		expenseService:  expenseSvc,
	}
}

// GetPublicDonations retorna todas as doações públicas
func (s *TransparencyService) GetPublicDonations() []TransparencyDonation {
	var publicDonations []TransparencyDonation

	// Filtrar apenas doações que foram completadas
	for _, donation := range s.donationService.donations {
		if donation.Status == "completed" {
			ngo, _ := s.donationService.GetNGOByID(donation.NGOID)

			publicDonation := TransparencyDonation{
				ID:              donation.ID,
				Amount:          donation.Amount,
				NGOName:         ngo.Name,
				NGOCategory:     ngo.Category,
				Date:            donation.CreatedAt,
				Status:          donation.Status,
				TransactionHash: donation.TransactionHash,
			}

			publicDonations = append(publicDonations, publicDonation)
		}
	}

	// Ordenar por data (mais recentes primeiro)
	sort.Slice(publicDonations, func(i, j int) bool {
		return publicDonations[i].Date.After(publicDonations[j].Date)
	})

	return publicDonations
}

// GetPublicExpenses retorna todas as despesas públicas
func (s *TransparencyService) GetPublicExpenses() []TransparencyExpense {
	var publicExpenses []TransparencyExpense

	// Filtrar apenas despesas aprovadas
	for _, expense := range s.expenseService.expenses {
		if expense.Status == "aprovado" {
			ngo, _ := s.donationService.GetNGOByID(expense.NGOID)

			publicExpense := TransparencyExpense{
				ID:            expense.ID,
				DonationID:    expense.DonationID,
				NGOName:       ngo.Name,
				Amount:        expense.Amount,
				Description:   expense.Description,
				Category:      expense.Category,
				Date:          expense.CreatedAt,
				ReceiptIPFS:   expense.ReceiptIPFS,
				BlockchainRef: expense.BlockchainRef,
				Status:        expense.Status,
			}

			publicExpenses = append(publicExpenses, publicExpense)
		}
	}

	// Ordenar por data (mais recentes primeiro)
	sort.Slice(publicExpenses, func(i, j int) bool {
		return publicExpenses[i].Date.After(publicExpenses[j].Date)
	})

	return publicExpenses
}

// GetDonationsByNGO retorna todas as doações recebidas por uma ONG específica
func (s *TransparencyService) GetDonationsByNGO(ngoID uint) ([]TransparencyDonation, error) {
	// Verificar se a ONG existe
	ngo, err := s.donationService.GetNGOByID(ngoID)
	if err != nil {
		return nil, fmt.Errorf("ONG não encontrada: %v", err)
	}

	var ngoDonations []TransparencyDonation

	// Filtrar doações da ONG
	for _, donation := range s.donationService.donations {
		if donation.NGOID == ngoID && donation.Status == "completed" {
			publicDonation := TransparencyDonation{
				ID:              donation.ID,
				Amount:          donation.Amount,
				NGOName:         ngo.Name,
				NGOCategory:     ngo.Category,
				Date:            donation.CreatedAt,
				Status:          donation.Status,
				TransactionHash: donation.TransactionHash,
			}

			ngoDonations = append(ngoDonations, publicDonation)
		}
	}

	// Ordenar por data (mais recentes primeiro)
	sort.Slice(ngoDonations, func(i, j int) bool {
		return ngoDonations[i].Date.After(ngoDonations[j].Date)
	})

	return ngoDonations, nil
}

// GetExpensesByNGO retorna todas as despesas de uma ONG específica
func (s *TransparencyService) GetExpensesByNGO(ngoID uint) ([]TransparencyExpense, error) {
	// Verificar se a ONG existe
	ngo, err := s.donationService.GetNGOByID(ngoID)
	if err != nil {
		return nil, fmt.Errorf("ONG não encontrada: %v", err)
	}

	var ngoExpenses []TransparencyExpense

	// Filtrar despesas da ONG
	for _, expense := range s.expenseService.expenses {
		if expense.NGOID == ngoID && expense.Status == "aprovado" {
			publicExpense := TransparencyExpense{
				ID:            expense.ID,
				DonationID:    expense.DonationID,
				NGOName:       ngo.Name,
				Amount:        expense.Amount,
				Description:   expense.Description,
				Category:      expense.Category,
				Date:          expense.CreatedAt,
				ReceiptIPFS:   expense.ReceiptIPFS,
				BlockchainRef: expense.BlockchainRef,
				Status:        expense.Status,
			}

			ngoExpenses = append(ngoExpenses, publicExpense)
		}
	}

	// Ordenar por data (mais recentes primeiro)
	sort.Slice(ngoExpenses, func(i, j int) bool {
		return ngoExpenses[i].Date.After(ngoExpenses[j].Date)
	})

	return ngoExpenses, nil
}

// GetNGOSummary retorna um resumo dos dados de transparência de uma ONG
func (s *TransparencyService) GetNGOSummary(ngoID uint) (TransparencyNGOSummary, error) {
	// Verificar se a ONG existe
	ngo, err := s.donationService.GetNGOByID(ngoID)
	if err != nil {
		return TransparencyNGOSummary{}, fmt.Errorf("ONG não encontrada: %v", err)
	}

	var totalReceived float64
	var donationsCount int

	// Calcular total recebido
	for _, donation := range s.donationService.donations {
		if donation.NGOID == ngoID && donation.Status == "completed" {
			totalReceived += donation.Amount
			donationsCount++
		}
	}

	var totalSpent float64
	var expensesCount int

	// Calcular total gasto
	for _, expense := range s.expenseService.expenses {
		if expense.NGOID == ngoID && expense.Status == "aprovado" {
			totalSpent += expense.Amount
			expensesCount++
		}
	}

	// Calcular saldo disponível
	availableBalance := totalReceived - totalSpent

	return TransparencyNGOSummary{
		ID:               ngo.ID,
		Name:             ngo.Name,
		Category:         ngo.Category,
		TotalReceived:    totalReceived,
		TotalSpent:       totalSpent,
		DonationsCount:   donationsCount,
		ExpensesCount:    expensesCount,
		AvailableBalance: availableBalance,
	}, nil
}

// GetAllNGOsSummary retorna um resumo de todas as ONGs
func (s *TransparencyService) GetAllNGOsSummary() []TransparencyNGOSummary {
	var summaries []TransparencyNGOSummary

	for _, ngo := range s.donationService.ngos {
		summary, err := s.GetNGOSummary(ngo.ID)
		if err == nil {
			summaries = append(summaries, summary)
		}
	}

	// Ordenar por valor total recebido (maior primeiro)
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalReceived > summaries[j].TotalReceived
	})

	return summaries
}

// GetTransparencyDashboard retorna o dashboard geral de transparência
func (s *TransparencyService) GetTransparencyDashboard() TransparencyDashboard {
	var totalDonations float64
	var donationsCount int

	// Contar doações completadas
	for _, donation := range s.donationService.donations {
		if donation.Status == "completed" {
			totalDonations += donation.Amount
			donationsCount++
		}
	}

	var totalExpenses float64
	var expensesCount int

	// Contar despesas aprovadas
	for _, expense := range s.expenseService.expenses {
		if expense.Status == "aprovado" {
			totalExpenses += expense.Amount
			expensesCount++
		}
	}

	// Obter doações recentes (limitado a 5)
	recentDonations := s.GetPublicDonations()
	if len(recentDonations) > 5 {
		recentDonations = recentDonations[:5]
	}

	// Obter despesas recentes (limitado a 5)
	recentExpenses := s.GetPublicExpenses()
	if len(recentExpenses) > 5 {
		recentExpenses = recentExpenses[:5]
	}

	// Obter resumo de todas as ONGs
	ngosSummary := s.GetAllNGOsSummary()

	return TransparencyDashboard{
		TotalDonations:  totalDonations,
		TotalExpenses:   totalExpenses,
		DonationsCount:  donationsCount,
		ExpensesCount:   expensesCount,
		NGOsCount:       len(s.donationService.ngos),
		RecentDonations: recentDonations,
		RecentExpenses:  recentExpenses,
		NGOsSummary:     ngosSummary,
	}
}
