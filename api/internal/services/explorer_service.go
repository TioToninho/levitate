package services

import (
	"errors"
	"strings"
	"time"
	"trackable-donations/api/internal/models"
)

// ExplorerService gerencia a busca e exploração de transações
type ExplorerService struct {
	donationService *DonationService
	expenseService  *ExpenseService
}

// NewExplorerService cria uma nova instância do serviço explorador
func NewExplorerService(donationSvc *DonationService, expenseSvc *ExpenseService) *ExplorerService {
	return &ExplorerService{
		donationService: donationSvc,
		expenseService:  expenseSvc,
	}
}

// SearchDonations busca doações com base nos critérios fornecidos
func (s *ExplorerService) SearchDonations(query models.TransactionExplorerQuery) (models.TransactionExplorerResult, error) {
	result := models.TransactionExplorerResult{
		Donations: []models.DonationDetails{},
		Page:      query.Page,
		PageSize:  query.PageSize,
	}

	// Definir valores padrão para paginação se não fornecidos
	if result.Page <= 0 {
		result.Page = 1
	}
	if result.PageSize <= 0 {
		result.PageSize = 10
	}

	// Filtrar doações com base nos critérios
	var filteredDonations []models.Donation
	for _, donation := range s.donationService.donations {
		// Filtrar apenas doações completadas
		if donation.Status != "completed" {
			continue
		}

		// Filtrar por hash de transação
		if query.TransactionHash != "" && !strings.EqualFold(donation.TransactionHash, query.TransactionHash) {
			continue
		}

		// Filtrar por ONG
		if query.NGOID != 0 && donation.NGOID != query.NGOID {
			continue
		}

		// Filtrar por período
		if !query.StartDate.IsZero() && donation.CreatedAt.Before(query.StartDate) {
			continue
		}
		if !query.EndDate.IsZero() && donation.CreatedAt.After(query.EndDate) {
			continue
		}

		filteredDonations = append(filteredDonations, donation)
	}

	// Calcular total
	result.Total = len(filteredDonations)

	// Aplicar paginação
	startIndex := (result.Page - 1) * result.PageSize
	endIndex := startIndex + result.PageSize
	if startIndex >= len(filteredDonations) {
		return result, nil
	}
	if endIndex > len(filteredDonations) {
		endIndex = len(filteredDonations)
	}

	// Processar doações selecionadas
	for _, donation := range filteredDonations[startIndex:endIndex] {
		donationDetail, err := s.getDonationDetails(donation)
		if err != nil {
			continue
		}
		result.Donations = append(result.Donations, donationDetail)
	}

	return result, nil
}

// GetDonationByHash obtém os detalhes de uma doação pelo hash de transação
func (s *ExplorerService) GetDonationByHash(hash string) (models.DonationDetails, error) {
	for _, donation := range s.donationService.donations {
		if strings.EqualFold(donation.TransactionHash, hash) {
			return s.getDonationDetails(donation)
		}
	}
	return models.DonationDetails{}, errors.New("doação não encontrada")
}

// GetDonationByID obtém os detalhes de uma doação pelo ID
func (s *ExplorerService) GetDonationByID(id uint) (models.DonationDetails, error) {
	for _, donation := range s.donationService.donations {
		if donation.ID == id {
			return s.getDonationDetails(donation)
		}
	}
	return models.DonationDetails{}, errors.New("doação não encontrada")
}

// getDonationDetails obtém os detalhes de uma doação
func (s *ExplorerService) getDonationDetails(donation models.Donation) (models.DonationDetails, error) {
	// Obter nome do doador
	donor, err := s.donationService.GetUserByID(donation.DonorID)
	if err != nil {
		return models.DonationDetails{}, err
	}

	// Obter nome e categoria da ONG
	ngo, err := s.donationService.GetNGOByID(donation.NGOID)
	if err != nil {
		return models.DonationDetails{}, err
	}

	// Verificar se tem recibo
	hasReceipt := false
	for _, receipt := range s.donationService.receipts {
		if receipt.DonationID == donation.ID {
			hasReceipt = true
			break
		}
	}

	// Verificar se tem despesas e contar
	hasExpenses := false
	expensesCount := 0
	for _, expense := range s.expenseService.expenses {
		if expense.DonationID == donation.ID {
			hasExpenses = true
			expensesCount++
		}
	}

	// Criar detalhes da doação
	details := models.DonationDetails{
		ID:              donation.ID,
		Amount:          donation.Amount,
		DonorName:       donor.Name,
		NGOName:         ngo.Name,
		NGOCategory:     ngo.Category,
		Date:            donation.CreatedAt,
		Status:          donation.Status,
		TransactionHash: donation.TransactionHash,
		HasReceipt:      hasReceipt,
		HasExpenses:     hasExpenses,
		ExpensesCount:   expensesCount,
	}

	return details, nil
}

// GetDonationsByNGO obtém todas as doações para uma ONG específica
func (s *ExplorerService) GetDonationsByNGO(ngoID uint, page, pageSize int) (models.TransactionExplorerResult, error) {
	query := models.TransactionExplorerQuery{
		NGOID:    ngoID,
		Page:     page,
		PageSize: pageSize,
	}
	return s.SearchDonations(query)
}

// GetDonationsByPeriod obtém todas as doações em um período específico
func (s *ExplorerService) GetDonationsByPeriod(startDate, endDate time.Time, page, pageSize int) (models.TransactionExplorerResult, error) {
	query := models.TransactionExplorerQuery{
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		PageSize:  pageSize,
	}
	return s.SearchDonations(query)
}

// GetRecentDonations obtém as doações mais recentes
func (s *ExplorerService) GetRecentDonations(limit int) ([]models.DonationDetails, error) {
	if limit <= 0 {
		limit = 10
	}

	// Filtrar apenas doações completadas
	var completedDonations []models.Donation
	for _, donation := range s.donationService.donations {
		if donation.Status == "completed" {
			completedDonations = append(completedDonations, donation)
		}
	}

	// Ordenar por data (mais recentes primeiro)
	// Em um sistema real, usaríamos ORDER BY na consulta SQL
	for i := 0; i < len(completedDonations)-1; i++ {
		for j := i + 1; j < len(completedDonations); j++ {
			if completedDonations[i].CreatedAt.Before(completedDonations[j].CreatedAt) {
				completedDonations[i], completedDonations[j] = completedDonations[j], completedDonations[i]
			}
		}
	}

	// Limitar ao número solicitado
	if len(completedDonations) > limit {
		completedDonations = completedDonations[:limit]
	}

	// Processar doações
	var details []models.DonationDetails
	for _, donation := range completedDonations {
		donationDetail, err := s.getDonationDetails(donation)
		if err != nil {
			continue
		}
		details = append(details, donationDetail)
	}

	return details, nil
}
