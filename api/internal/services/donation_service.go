package services

import (
	"errors"
	"fmt"
	"log"
	"time"
	"trackable-donations/api/internal/models"
)

// DonationService gerencia operações relacionadas a doações
type DonationService struct {
	// Em um sistema real, teríamos repositórios para acesso ao banco de dados
	// Aqui usaremos dados em memória para demonstração
	donations      []models.Donation
	ngos           []models.NGO
	users          []models.User
	resourceUsages []models.ResourceUsage
	receipts       []models.DonationReceipt
}

// NewDonationService cria uma nova instância do serviço
func NewDonationService() *DonationService {
	// Inicializa com algumas ONGs para demonstração
	ngos := []models.NGO{
		{ID: 1, Name: "Alimentando Esperança", Description: "Distribuição de alimentos para pessoas em situação de vulnerabilidade", Category: "Alimentação", LogoURL: "https://example.com/logo1.png"},
		{ID: 2, Name: "Saúde para Todos", Description: "Fornecimento de medicamentos e atendimento médico gratuito", Category: "Saúde", LogoURL: "https://example.com/logo2.png"},
		{ID: 3, Name: "Educação é Futuro", Description: "Apoio educacional para crianças de baixa renda", Category: "Educação", LogoURL: "https://example.com/logo3.png"},
	}

	// Inicializa com alguns usuários para demonstração
	users := []models.User{
		{ID: 1, Name: "João Silva", Email: "joao@example.com", CreatedAt: time.Now()},
		{ID: 2, Name: "Maria Oliveira", Email: "maria@example.com", CreatedAt: time.Now()},
	}

	return &DonationService{
		donations:      []models.Donation{},
		ngos:           ngos,
		users:          users,
		resourceUsages: []models.ResourceUsage{},
		receipts:       []models.DonationReceipt{},
	}
}

// GetAllNGOs retorna todas as ONGs disponíveis
func (s *DonationService) GetAllNGOs() []models.NGO {
	return s.ngos
}

// GetNGOByID busca uma ONG pelo ID
func (s *DonationService) GetNGOByID(id uint) (models.NGO, error) {
	for _, ngo := range s.ngos {
		if ngo.ID == id {
			return ngo, nil
		}
	}
	return models.NGO{}, errors.New("ONG não encontrada")
}

// GetUserByID busca um usuário pelo ID
func (s *DonationService) GetUserByID(id uint) (models.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, errors.New("Usuário não encontrado")
}

// ProcessDonation processa uma nova doação
func (s *DonationService) ProcessDonation(req models.DonationRequest) (models.DonationResponse, error) {
	// Verificar se a ONG existe
	_, err := s.GetNGOByID(req.NGOID)
	if err != nil {
		return models.DonationResponse{}, err
	}

	// Verificar se o doador existe
	_, err = s.GetUserByID(req.DonorID)
	if err != nil {
		return models.DonationResponse{}, err
	}

	// Criar nova doação
	donationID := uint(len(s.donations) + 1) // Em um banco real, seria auto-incremento
	donation := models.Donation{
		ID:        donationID,
		Amount:    req.Amount,
		DonorID:   req.DonorID,
		NGOID:     req.NGOID,
		CreatedAt: time.Now(),
		Status:    "pending", // Inicialmente pendente
	}

	// Adicionar à lista (em um sistema real, seria salvo no banco)
	s.donations = append(s.donations, donation)

	// Simular url de pagamento
	paymentURL := fmt.Sprintf("https://payment-gateway-mock.com/pay?donationId=%d&amount=%.2f", donation.ID, donation.Amount)

	return models.DonationResponse{
		ID:         donation.ID,
		Status:     donation.Status,
		PaymentURL: paymentURL,
	}, nil
}

// MockPaymentConfirmation simula a confirmação de pagamento pelo gateway
func (s *DonationService) MockPaymentConfirmation(donationID uint) (models.DonationResponse, error) {
	// Encontrar a doação
	var donation models.Donation
	var donorID uint
	var ngoID uint
	found := false

	for i, d := range s.donations {
		if d.ID == donationID {
			donorID = d.DonorID
			ngoID = d.NGOID
			donation = d
			// Atualizar o status
			s.donations[i].Status = "completed"
			// Gerar hash fictício para simulação de blockchain
			s.donations[i].TransactionHash = generateMockTransactionHash()
			donation = s.donations[i]
			found = true
			break
		}
	}

	if !found {
		return models.DonationResponse{}, errors.New("doação não encontrada")
	}

	// Simular registro na blockchain (em um sistema real, registraríamos na blockchain)
	log.Printf("Registrando doação na blockchain: %v", donation)

	// Gerar comprovante de doação
	s.generateDonationReceipt(donation, donorID, ngoID)

	// Gerar uso dos recursos (mockado)
	s.mockResourceUsage(donation)

	return models.DonationResponse{
		ID:              donation.ID,
		Status:          donation.Status,
		TransactionHash: donation.TransactionHash,
	}, nil
}

// generateDonationReceipt gera um comprovante de doação
func (s *DonationService) generateDonationReceipt(donation models.Donation, donorID, ngoID uint) models.DonationReceipt {
	donor, _ := s.GetUserByID(donorID)
	ngo, _ := s.GetNGOByID(ngoID)

	// Simular um hash IPFS para o comprovante
	ipfsHash := fmt.Sprintf("Qm%s", generateMockHash(46))

	receiptID := uint(len(s.receipts) + 1)
	receipt := models.DonationReceipt{
		ID:              receiptID,
		DonationID:      donation.ID,
		DonorName:       donor.Name,
		DonorEmail:      donor.Email,
		NGOName:         ngo.Name,
		Amount:          donation.Amount,
		Date:            donation.CreatedAt,
		TransactionHash: donation.TransactionHash,
		IPFSHash:        ipfsHash,
		PdfURL:          fmt.Sprintf("https://ipfs.example.com/ipfs/%s", ipfsHash),
	}

	s.receipts = append(s.receipts, receipt)
	return receipt
}

// mockResourceUsage simula o uso dos recursos da doação
func (s *DonationService) mockResourceUsage(donation models.Donation) {
	ngo, _ := s.GetNGOByID(donation.NGOID)
	amount := donation.Amount

	// Simular diferentes tipos de uso de recursos baseados na categoria da ONG
	var descriptions []string
	var usageAmounts []float64
	var percentages []float64

	switch ngo.Category {
	case "Alimentação":
		descriptions = []string{
			"Compra de alimentos não perecíveis",
			"Transporte de doações",
			"Material para empacotamento",
		}
		percentages = []float64{0.7, 0.2, 0.1}
	case "Saúde":
		descriptions = []string{
			"Compra de medicamentos",
			"Equipamentos médicos descartáveis",
			"Apoio a atendimentos médicos",
		}
		percentages = []float64{0.6, 0.3, 0.1}
	case "Educação":
		descriptions = []string{
			"Material escolar",
			"Livros didáticos",
			"Equipamentos para escola",
		}
		percentages = []float64{0.4, 0.4, 0.2}
	default:
		descriptions = []string{
			"Uso principal dos recursos",
			"Custos administrativos",
		}
		percentages = []float64{0.8, 0.2}
	}

	// Calcular os valores baseados nos percentuais
	usageAmounts = make([]float64, len(percentages))
	for i, percentage := range percentages {
		usageAmounts[i] = amount * percentage
	}

	// Criar os registros de uso
	for i, description := range descriptions {
		usageDate := donation.CreatedAt.Add(time.Duration(i*24) * time.Hour) // Cada uso alguns dias depois
		ipfsHash := fmt.Sprintf("Qm%s", generateMockHash(46))

		usage := models.ResourceUsage{
			ID:          uint(len(s.resourceUsages) + i + 1),
			DonationID:  donation.ID,
			Description: description,
			Amount:      usageAmounts[i],
			Date:        usageDate,
			ReceiptIPFS: ipfsHash,
			NGOName:     ngo.Name,
			CreatedAt:   time.Now(),
		}

		s.resourceUsages = append(s.resourceUsages, usage)
	}
}

// GetDonationsByDonorID retorna todas as doações de um doador
func (s *DonationService) GetDonationsByDonorID(donorID uint) ([]models.Donation, error) {
	// Verificar se o doador existe
	_, err := s.GetUserByID(donorID)
	if err != nil {
		return nil, err
	}

	var donorDonations []models.Donation
	for _, donation := range s.donations {
		if donation.DonorID == donorID {
			donorDonations = append(donorDonations, donation)
		}
	}

	return donorDonations, nil
}

// GetDonationReceipt retorna o comprovante de uma doação
func (s *DonationService) GetDonationReceipt(donationID uint) (models.DonationReceipt, error) {
	for _, receipt := range s.receipts {
		if receipt.DonationID == donationID {
			return receipt, nil
		}
	}
	return models.DonationReceipt{}, errors.New("comprovante não encontrado")
}

// GetResourceUsagesByDonationID retorna os usos dos recursos de uma doação
func (s *DonationService) GetResourceUsagesByDonationID(donationID uint) ([]models.ResourceUsage, error) {
	// Verificar se a doação existe
	found := false
	for _, donation := range s.donations {
		if donation.ID == donationID {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("doação não encontrada")
	}

	var usages []models.ResourceUsage
	for _, usage := range s.resourceUsages {
		if usage.DonationID == donationID {
			usages = append(usages, usage)
		}
	}

	return usages, nil
}

// GetDonorDashboard retorna o dashboard de um doador
func (s *DonationService) GetDonorDashboard(donorID uint) (models.DonorDashboard, error) {
	// Verificar se o doador existe
	donor, err := s.GetUserByID(donorID)
	if err != nil {
		return models.DonorDashboard{}, err
	}

	// Obter doações do doador
	donations, err := s.GetDonationsByDonorID(donorID)
	if err != nil {
		return models.DonorDashboard{}, err
	}

	// Calcular métricas
	var totalDonated float64
	var ngosMap = make(map[uint]bool)

	for _, donation := range donations {
		totalDonated += donation.Amount
		ngosMap[donation.NGOID] = true
	}

	// Calcular métricas fictícias de impacto
	peopleHelped := int(totalDonated / 50)      // Estima 1 pessoa ajudada a cada R$ 50
	mealsProvided := int(totalDonated / 10)     // Estima 1 refeição a cada R$ 10
	medicinesProvided := int(totalDonated / 30) // Estima 1 medicamento a cada R$ 30

	metrics := models.ImpactMetrics{
		TotalDonated:      totalDonated,
		DonationsCount:    len(donations),
		NGOsSupported:     len(ngosMap),
		PeopleHelped:      peopleHelped,
		MealsProvided:     mealsProvided,
		MedicinesProvided: medicinesProvided,
	}

	// Contar todos os usos de recursos relacionados às doações do usuário
	var usagesCount int
	for _, usage := range s.resourceUsages {
		for _, donation := range donations {
			if usage.DonationID == donation.ID {
				usagesCount++
				break
			}
		}
	}

	return models.DonorDashboard{
		DonorID:     donorID,
		DonorName:   donor.Name,
		Metrics:     metrics,
		Donations:   donations,
		UsagesCount: usagesCount,
	}, nil
}
