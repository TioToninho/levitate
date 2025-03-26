package services

import (
	"errors"
	"fmt"
	"time"
	"trackable-donations/api/internal/models"
)

// ExpenseService gerencia operações relacionadas a gastos das ONGs
type ExpenseService struct {
	// Em um sistema real, teríamos repositórios para acesso ao banco de dados
	expenses    []models.Expense
	donationSvc *DonationService
}

// NewExpenseService cria uma nova instância do serviço de gastos
func NewExpenseService(donationSvc *DonationService) *ExpenseService {
	return &ExpenseService{
		expenses:    []models.Expense{},
		donationSvc: donationSvc,
	}
}

// RegisterExpense registra um novo gasto relacionado a uma doação
func (s *ExpenseService) RegisterExpense(req models.ExpenseRequest) (models.ExpenseResponse, error) {
	// Verificar se a doação existe
	found := false
	var donation models.Donation

	for _, d := range s.donationSvc.donations {
		if d.ID == req.DonationID {
			donation = d
			found = true
			break
		}
	}

	if !found {
		return models.ExpenseResponse{}, errors.New("doação não encontrada")
	}

	// Verificar se a ONG é a mesma da doação
	if donation.NGOID != req.NGOID {
		return models.ExpenseResponse{}, errors.New("esta ONG não está associada a esta doação")
	}

	// Verificar status da doação
	if donation.Status != "completed" {
		return models.ExpenseResponse{}, errors.New("só é possível registrar gastos para doações confirmadas")
	}

	// Verificar se o valor do gasto não excede o total disponível
	totalExpenses := float64(0)
	for _, e := range s.expenses {
		if e.DonationID == req.DonationID {
			totalExpenses += e.Amount
		}
	}

	remainingAmount := donation.Amount - totalExpenses

	if req.Amount > remainingAmount {
		return models.ExpenseResponse{}, fmt.Errorf("valor excede o saldo disponível da doação (%.2f)", remainingAmount)
	}

	// Criar novo gasto
	expenseID := uint(len(s.expenses) + 1) // Em um banco real, seria auto-incremento

	expense := models.Expense{
		ID:          expenseID,
		DonationID:  req.DonationID,
		NGOID:       req.NGOID,
		Amount:      req.Amount,
		Description: req.Description,
		Category:    req.Category,
		Status:      "pendente", // Inicialmente pendente até upload de comprovante
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Adicionar à lista (em um sistema real, seria salvo no banco)
	s.expenses = append(s.expenses, expense)

	return models.ExpenseResponse{
		ID:          expense.ID,
		DonationID:  expense.DonationID,
		NGOID:       expense.NGOID,
		Amount:      expense.Amount,
		Description: expense.Description,
		Category:    expense.Category,
		Status:      expense.Status,
		CreatedAt:   expense.CreatedAt,
	}, nil
}

// UploadReceipt faz upload do comprovante para o IPFS e atualiza o gasto
func (s *ExpenseService) UploadReceipt(expenseID uint, fileContent []byte) (models.ExpenseResponse, error) {
	// Encontrar o gasto
	found := false
	var index int

	for i, e := range s.expenses {
		if e.ID == expenseID {
			index = i
			found = true
			break
		}
	}

	if !found {
		return models.ExpenseResponse{}, errors.New("gasto não encontrado")
	}

	// Em um sistema real, faríamos o upload para o IPFS
	// Por ora, simularemos com um hash
	ipfsHash := fmt.Sprintf("Qm%s", generateMockHash(46))

	// Em um sistema real, registraríamos na blockchain
	blockchainRef := generateMockTransactionHash()

	// Atualizar o gasto
	s.expenses[index].ReceiptIPFS = ipfsHash
	s.expenses[index].BlockchainRef = blockchainRef
	s.expenses[index].Status = "aprovado"
	s.expenses[index].UpdatedAt = time.Now()

	// Retornar o gasto atualizado
	return models.ExpenseResponse{
		ID:            s.expenses[index].ID,
		DonationID:    s.expenses[index].DonationID,
		NGOID:         s.expenses[index].NGOID,
		Amount:        s.expenses[index].Amount,
		Description:   s.expenses[index].Description,
		Category:      s.expenses[index].Category,
		ReceiptIPFS:   s.expenses[index].ReceiptIPFS,
		BlockchainRef: s.expenses[index].BlockchainRef,
		Status:        s.expenses[index].Status,
		CreatedAt:     s.expenses[index].CreatedAt,
	}, nil
}

// GetExpensesByDonation obtém todos os gastos relacionados a uma doação
func (s *ExpenseService) GetExpensesByDonation(donationID uint) ([]models.ExpenseResponse, error) {
	var expenseResponses []models.ExpenseResponse

	for _, e := range s.expenses {
		if e.DonationID == donationID {
			expenseResponses = append(expenseResponses, models.ExpenseResponse{
				ID:            e.ID,
				DonationID:    e.DonationID,
				NGOID:         e.NGOID,
				Amount:        e.Amount,
				Description:   e.Description,
				Category:      e.Category,
				ReceiptIPFS:   e.ReceiptIPFS,
				BlockchainRef: e.BlockchainRef,
				Status:        e.Status,
				CreatedAt:     e.CreatedAt,
			})
		}
	}

	return expenseResponses, nil
}

// GetExpensesByNGO obtém todos os gastos relacionados a uma ONG
func (s *ExpenseService) GetExpensesByNGO(ngoID uint) ([]models.ExpenseResponse, error) {
	var expenseResponses []models.ExpenseResponse

	for _, e := range s.expenses {
		if e.NGOID == ngoID {
			expenseResponses = append(expenseResponses, models.ExpenseResponse{
				ID:            e.ID,
				DonationID:    e.DonationID,
				NGOID:         e.NGOID,
				Amount:        e.Amount,
				Description:   e.Description,
				Category:      e.Category,
				ReceiptIPFS:   e.ReceiptIPFS,
				BlockchainRef: e.BlockchainRef,
				Status:        e.Status,
				CreatedAt:     e.CreatedAt,
			})
		}
	}

	return expenseResponses, nil
}
