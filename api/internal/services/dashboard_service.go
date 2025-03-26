package services

import (
	"fmt"
	"math"
	"sort"
	"time"
	"trackable-donations/api/internal/models"
)

// DashboardService gerencia as operações relacionadas ao dashboard global
type DashboardService struct {
	donationService *DonationService
	expenseService  *ExpenseService
}

// NewDashboardService cria uma nova instância do serviço de dashboard
func NewDashboardService(donationSvc *DonationService, expenseSvc *ExpenseService) *DashboardService {
	return &DashboardService{
		donationService: donationSvc,
		expenseService:  expenseSvc,
	}
}

// GetGlobalDashboard obtém os dados para o dashboard global
func (s *DashboardService) GetGlobalDashboard() models.GlobalDashboardData {
	dashboard := models.GlobalDashboardData{}

	// Filtrar apenas doações completadas
	var completedDonations []models.Donation
	donorMap := make(map[uint]struct{}) // Para contar doadores únicos
	for _, donation := range s.donationService.donations {
		if donation.Status == "completed" {
			completedDonations = append(completedDonations, donation)
			donorMap[donation.DonorID] = struct{}{}
			dashboard.TotalDonated += donation.Amount
		}
	}

	// Calcular totais
	dashboard.TotalTransactions = len(completedDonations)
	dashboard.TotalDonors = len(donorMap)
	dashboard.TotalNGOs = len(s.donationService.ngos)

	// Calcular doações por categoria
	dashboard.DonationsByCategory = s.calculateDonationsByCategory(completedDonations)

	// Calcular doações mensais
	dashboard.MonthlyDonations = s.calculateMonthlyDonations(completedDonations)

	// Calcular top ONGs
	dashboard.TopNGOs = s.calculateTopNGOs(completedDonations, 5)

	// Gerar dados geográficos simulados
	dashboard.GeographicalData = s.generateGeographicalData()

	// Calcular métricas de impacto
	dashboard.ImpactMetrics = s.calculateImpactMetrics(dashboard.TotalDonated)

	return dashboard
}

// calculateDonationsByCategory calcula as doações por categoria
func (s *DashboardService) calculateDonationsByCategory(donations []models.Donation) []models.CategorySummary {
	categoryMap := make(map[string]models.CategorySummary)

	for _, donation := range donations {
		ngo, err := s.donationService.GetNGOByID(donation.NGOID)
		if err != nil {
			continue
		}

		category := ngo.Category
		summary, exists := categoryMap[category]
		if exists {
			summary.TotalAmount += donation.Amount
			summary.Count++
		} else {
			summary = models.CategorySummary{
				Category:    category,
				TotalAmount: donation.Amount,
				Count:       1,
			}
		}
		categoryMap[category] = summary
	}

	// Converter mapa para slice
	var categorySummaries []models.CategorySummary
	totalAmount := 0.0
	for _, summary := range categoryMap {
		categorySummaries = append(categorySummaries, summary)
		totalAmount += summary.TotalAmount
	}

	// Calcular percentagens
	if totalAmount > 0 {
		for i := range categorySummaries {
			categorySummaries[i].Percentage = math.Round((categorySummaries[i].TotalAmount/totalAmount)*100) / 100
		}
	}

	// Ordenar por valor total (maior primeiro)
	sort.Slice(categorySummaries, func(i, j int) bool {
		return categorySummaries[i].TotalAmount > categorySummaries[j].TotalAmount
	})

	return categorySummaries
}

// calculateMonthlyDonations calcula as doações mensais
func (s *DashboardService) calculateMonthlyDonations(donations []models.Donation) []models.MonthlyDonationData {
	monthMap := make(map[string]models.MonthlyDonationData)

	// Processar cada doação
	for _, donation := range donations {
		// Formatar chave como "YYYY-MM"
		key := fmt.Sprintf("%d-%02d", donation.CreatedAt.Year(), donation.CreatedAt.Month())

		// Construir nome do mês em português
		monthName := s.getMonthName(int(donation.CreatedAt.Month()))

		data, exists := monthMap[key]
		if exists {
			data.TotalAmount += donation.Amount
			data.Count++
		} else {
			data = models.MonthlyDonationData{
				Month:       monthName,
				Year:        donation.CreatedAt.Year(),
				TotalAmount: donation.Amount,
				Count:       1,
			}
		}
		monthMap[key] = data
	}

	// Converter mapa para slice
	var monthlyData []models.MonthlyDonationData
	for _, data := range monthMap {
		monthlyData = append(monthlyData, data)
	}

	// Ordenar por data (mais antigo primeiro)
	sort.Slice(monthlyData, func(i, j int) bool {
		if monthlyData[i].Year != monthlyData[j].Year {
			return monthlyData[i].Year < monthlyData[j].Year
		}
		// Compare month names (not ideal, but works for our simulated data)
		return s.getMonthIndex(monthlyData[i].Month) < s.getMonthIndex(monthlyData[j].Month)
	})

	return monthlyData
}

// getMonthName retorna o nome do mês em português
func (s *DashboardService) getMonthName(month int) string {
	months := []string{
		"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho",
		"Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro",
	}
	if month >= 1 && month <= 12 {
		return months[month-1]
	}
	return ""
}

// getMonthIndex retorna o índice do mês a partir do nome
func (s *DashboardService) getMonthIndex(monthName string) int {
	months := map[string]int{
		"Janeiro": 1, "Fevereiro": 2, "Março": 3, "Abril": 4,
		"Maio": 5, "Junho": 6, "Julho": 7, "Agosto": 8,
		"Setembro": 9, "Outubro": 10, "Novembro": 11, "Dezembro": 12,
	}
	return months[monthName]
}

// calculateTopNGOs calcula as ONGs com mais doações
func (s *DashboardService) calculateTopNGOs(donations []models.Donation, limit int) []models.NGODonationSummary {
	ngoMap := make(map[uint]models.NGODonationSummary)

	// Processar cada doação
	for _, donation := range donations {
		ngo, err := s.donationService.GetNGOByID(donation.NGOID)
		if err != nil {
			continue
		}

		summary, exists := ngoMap[ngo.ID]
		if exists {
			summary.TotalAmount += donation.Amount
			summary.Count++
		} else {
			summary = models.NGODonationSummary{
				NGOID:       ngo.ID,
				NGOName:     ngo.Name,
				Category:    ngo.Category,
				TotalAmount: donation.Amount,
				Count:       1,
			}
		}
		ngoMap[ngo.ID] = summary
	}

	// Converter mapa para slice
	var ngoSummaries []models.NGODonationSummary
	for _, summary := range ngoMap {
		ngoSummaries = append(ngoSummaries, summary)
	}

	// Ordenar por valor total (maior primeiro)
	sort.Slice(ngoSummaries, func(i, j int) bool {
		return ngoSummaries[i].TotalAmount > ngoSummaries[j].TotalAmount
	})

	// Limitar ao número solicitado
	if len(ngoSummaries) > limit {
		ngoSummaries = ngoSummaries[:limit]
	}

	return ngoSummaries
}

// generateGeographicalData gera dados geográficos simulados
func (s *DashboardService) generateGeographicalData() []models.GeographicalDonationData {
	// Em um sistema real, estes dados viriam do banco de dados
	// Aqui estamos simulando com regiões do Brasil
	regions := []string{
		"Norte", "Nordeste", "Centro-Oeste", "Sudeste", "Sul",
	}

	// Criar dados simulados
	var geoData []models.GeographicalDonationData
	totalDonations := float64(0)

	// Contabilizar doações totais para calcular proporções realistas
	for _, donation := range s.donationService.donations {
		if donation.Status == "completed" {
			totalDonations += donation.Amount
		}
	}

	// Distribuir proporcionalmente com base em uma distribuição simulada
	distribution := []float64{0.1, 0.15, 0.15, 0.4, 0.2} // 10%, 15%, 15%, 40%, 20%

	for i, region := range regions {
		amount := totalDonations * distribution[i]
		count := int(float64(len(s.donationService.donations)) * distribution[i])

		geoData = append(geoData, models.GeographicalDonationData{
			Region:      region,
			TotalAmount: math.Round(amount*100) / 100, // Arredondar para 2 casas decimais
			Count:       count,
		})
	}

	return geoData
}

// calculateImpactMetrics calcula métricas de impacto simuladas
func (s *DashboardService) calculateImpactMetrics(totalDonated float64) models.GlobalImpactMetrics {
	// Em um sistema real, esses dados seriam baseados em relatórios reais de impacto
	// Aqui estamos simulando com base no valor total doado

	// Fatores de conversão simulados (por exemplo, R$ 100 = 10 refeições)
	mealsPerMoney := 0.1         // 1 refeição por R$ 10
	medicinesPerMoney := 0.02    // 1 medicamento por R$ 50
	educationPerMoney := 0.01    // 1 criança educada por R$ 100
	housesPerMoney := 0.0005     // 1 casa por R$ 2.000
	emergenciesPerMoney := 0.005 // 1 emergência por R$ 200

	metrics := models.GlobalImpactMetrics{
		PeopleHelped:      int(totalDonated * 0.5),   // Cada R$ 2 ajuda 1 pessoa
		CommunitiesServed: int(totalDonated * 0.005), // Cada R$ 200 atende 1 comunidade
		ProjectsCompleted: int(totalDonated * 0.02),  // Cada R$ 50 completa 1 projeto
		MealsProvided:     int(totalDonated * mealsPerMoney),
		MedicinesProvided: int(totalDonated * medicinesPerMoney),
		ChildrenEducated:  int(totalDonated * educationPerMoney),
		HousesBuilt:       int(totalDonated * housesPerMoney),
		EmergenciesServed: int(totalDonated * emergenciesPerMoney),
	}

	return metrics
}

// GetDashboardByDateRange obtém dados do dashboard para um intervalo de datas específico
func (s *DashboardService) GetDashboardByDateRange(startDate, endDate time.Time) models.GlobalDashboardData {
	// Filtrar doações pelo intervalo de datas
	var filteredDonations []models.Donation
	for _, donation := range s.donationService.donations {
		if donation.Status == "completed" &&
			(startDate.IsZero() || !donation.CreatedAt.Before(startDate)) &&
			(endDate.IsZero() || !donation.CreatedAt.After(endDate)) {
			filteredDonations = append(filteredDonations, donation)
		}
	}

	// Calcular dashboard com as doações filtradas
	dashboard := models.GlobalDashboardData{}
	donorMap := make(map[uint]struct{})

	for _, donation := range filteredDonations {
		dashboard.TotalDonated += donation.Amount
		donorMap[donation.DonorID] = struct{}{}
	}

	dashboard.TotalTransactions = len(filteredDonations)
	dashboard.TotalDonors = len(donorMap)
	dashboard.TotalNGOs = len(s.donationService.ngos)
	dashboard.DonationsByCategory = s.calculateDonationsByCategory(filteredDonations)
	dashboard.MonthlyDonations = s.calculateMonthlyDonations(filteredDonations)
	dashboard.TopNGOs = s.calculateTopNGOs(filteredDonations, 5)
	dashboard.ImpactMetrics = s.calculateImpactMetrics(dashboard.TotalDonated)

	return dashboard
}

// GetDashboardByCategory obtém dados do dashboard para uma categoria específica
func (s *DashboardService) GetDashboardByCategory(category string) models.GlobalDashboardData {
	// Filtrar doações pela categoria da ONG
	var filteredDonations []models.Donation
	for _, donation := range s.donationService.donations {
		if donation.Status != "completed" {
			continue
		}

		ngo, err := s.donationService.GetNGOByID(donation.NGOID)
		if err != nil {
			continue
		}

		if ngo.Category == category {
			filteredDonations = append(filteredDonations, donation)
		}
	}

	// Calcular dashboard com as doações filtradas
	dashboard := models.GlobalDashboardData{}
	donorMap := make(map[uint]struct{})

	for _, donation := range filteredDonations {
		dashboard.TotalDonated += donation.Amount
		donorMap[donation.DonorID] = struct{}{}
	}

	dashboard.TotalTransactions = len(filteredDonations)
	dashboard.TotalDonors = len(donorMap)

	// Contar ONGs nesta categoria
	var ngosInCategory int
	for _, ngo := range s.donationService.ngos {
		if ngo.Category == category {
			ngosInCategory++
		}
	}
	dashboard.TotalNGOs = ngosInCategory

	dashboard.DonationsByCategory = s.calculateDonationsByCategory(filteredDonations)
	dashboard.MonthlyDonations = s.calculateMonthlyDonations(filteredDonations)
	dashboard.TopNGOs = s.calculateTopNGOs(filteredDonations, 5)
	dashboard.ImpactMetrics = s.calculateImpactMetrics(dashboard.TotalDonated)

	return dashboard
}
