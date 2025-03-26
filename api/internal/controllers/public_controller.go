package controllers

import (
	"net/http"
	"strconv"
	"time"
	"trackable-donations/api/internal/models"
	"trackable-donations/api/internal/services"

	"github.com/gin-gonic/gin"
)

// ExplorerService é a instância do serviço do explorador de transações
var ExplorerService *services.ExplorerService

// DashboardService é a instância do serviço de dashboard
var DashboardService *services.DashboardService

// SetupPublicServices configura os serviços públicos
func SetupPublicServices(donationService *services.DonationService, expenseService *services.ExpenseService) {
	ExplorerService = services.NewExplorerService(donationService, expenseService)
	DashboardService = services.NewDashboardService(donationService, expenseService)
}

// SearchDonations processa a busca de doações
// @Summary Buscar doações
// @Description Busca doações com filtros por hash, ONG e período
// @Tags Explorador
// @Accept json
// @Produce json
// @Param hash query string false "Hash da transação na blockchain"
// @Param ngo_id query int false "ID da ONG"
// @Param start_date query string false "Data inicial (formato: YYYY-MM-DD)"
// @Param end_date query string false "Data final (formato: YYYY-MM-DD)"
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} models.TransactionExplorerResult
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /explorer/search [get]
func SearchDonations(ctx *gin.Context) {
	// Criar objeto de consulta
	var query models.TransactionExplorerQuery

	// Obter parâmetros de consulta
	if hash := ctx.Query("hash"); hash != "" {
		query.TransactionHash = hash
	}

	if ngoIDStr := ctx.Query("ngo_id"); ngoIDStr != "" {
		ngoID, err := strconv.ParseUint(ngoIDStr, 10, 32)
		if err == nil {
			query.NGOID = uint(ngoID)
		}
	}

	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			query.StartDate = startDate
		}
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			// Definir para o final do dia
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			query.EndDate = endDate
		}
	}

	// Obter parâmetros de paginação
	if pageStr := ctx.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil && page > 0 {
			query.Page = page
		}
	}

	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err == nil && pageSize > 0 {
			query.PageSize = pageSize
		}
	}

	// Executar a busca
	result, err := ExplorerService.SearchDonations(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetDonationByHash obtém os detalhes de uma doação pelo hash
// @Summary Obter doação por hash
// @Description Retorna os detalhes de uma doação pelo hash da transação na blockchain
// @Tags Explorador
// @Accept json
// @Produce json
// @Param hash path string true "Hash da transação"
// @Success 200 {object} models.DonationDetails
// @Failure 400 {object} map[string]string "Hash não fornecido"
// @Failure 404 {object} map[string]string "Doação não encontrada"
// @Router /explorer/donations/hash/{hash} [get]
func GetDonationByHash(ctx *gin.Context) {
	hash := ctx.Param("hash")
	if hash == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Hash não fornecido"})
		return
	}

	donation, err := ExplorerService.GetDonationByHash(hash)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, donation)
}

// GetDonationByID obtém os detalhes de uma doação pelo ID
// @Summary Obter doação por ID
// @Description Retorna os detalhes de uma doação pelo ID
// @Tags Explorador
// @Accept json
// @Produce json
// @Param id path int true "ID da doação"
// @Success 200 {object} models.DonationDetails
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Doação não encontrada"
// @Router /explorer/donations/{id} [get]
func GetDonationByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	donation, err := ExplorerService.GetDonationByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, donation)
}

// GetDonationsByNGO obtém as doações de uma ONG específica
// @Summary Listar doações por ONG
// @Description Retorna todas as doações recebidas por uma ONG específica
// @Tags Explorador
// @Accept json
// @Produce json
// @Param ngo_id path int true "ID da ONG"
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} models.TransactionExplorerResult
// @Failure 400 {object} map[string]string "ID de ONG inválido"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /explorer/donations/ngo/{ngo_id} [get]
func GetDonationsByNGO(ctx *gin.Context) {
	ngoIDStr := ctx.Param("ngo_id")
	ngoID, err := strconv.ParseUint(ngoIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de ONG inválido"})
		return
	}

	// Obter parâmetros de paginação
	page := 1
	pageSize := 10

	if pageStr := ctx.Query("page"); pageStr != "" {
		pageVal, err := strconv.Atoi(pageStr)
		if err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		pageSizeVal, err := strconv.Atoi(pageSizeStr)
		if err == nil && pageSizeVal > 0 {
			pageSize = pageSizeVal
		}
	}

	result, err := ExplorerService.GetDonationsByNGO(uint(ngoID), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetRecentDonations obtém as doações mais recentes
// @Summary Listar doações recentes
// @Description Retorna as doações mais recentes
// @Tags Explorador
// @Accept json
// @Produce json
// @Param limit query int false "Limite de resultados (padrão: 10)"
// @Success 200 {array} models.DonationDetails
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /explorer/donations/recent [get]
func GetRecentDonations(ctx *gin.Context) {
	limit := 10
	if limitStr := ctx.Query("limit"); limitStr != "" {
		limitVal, err := strconv.Atoi(limitStr)
		if err == nil && limitVal > 0 {
			limit = limitVal
		}
	}

	donations, err := ExplorerService.GetRecentDonations(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, donations)
}

// GetGlobalDashboard obtém os dados do dashboard global
// @Summary Obter dashboard global
// @Description Retorna todas as métricas e dados consolidados do dashboard global
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} models.GlobalDashboardData
// @Router /dashboard/global [get]
func GetGlobalDashboard(ctx *gin.Context) {
	dashboard := DashboardService.GetGlobalDashboard()
	ctx.JSON(http.StatusOK, dashboard)
}

// GetDashboardByDateRange obtém os dados do dashboard para um intervalo de datas
// @Summary Obter dashboard por período
// @Description Retorna dados do dashboard filtrados por período de tempo
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param start_date query string true "Data inicial (formato: YYYY-MM-DD)"
// @Param end_date query string true "Data final (formato: YYYY-MM-DD)"
// @Success 200 {object} models.GlobalDashboardData
// @Failure 400 {object} map[string]string "Formato de data inválido"
// @Router /dashboard/by-date-range [get]
func GetDashboardByDateRange(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datas de início e fim são obrigatórias"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido para data inicial"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido para data final"})
		return
	}

	// Definir o fim do dia para a data final
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	dashboard := DashboardService.GetDashboardByDateRange(startDate, endDate)
	ctx.JSON(http.StatusOK, dashboard)
}

// GetDashboardByCategory obtém os dados do dashboard para uma categoria específica
// @Summary Obter dashboard por categoria
// @Description Retorna dados do dashboard filtrados por categoria de ONG
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param category path string true "Categoria (ex: Educação, Saúde, etc.)"
// @Success 200 {object} models.GlobalDashboardData
// @Failure 400 {object} map[string]string "Categoria não fornecida"
// @Router /dashboard/by-category/{category} [get]
func GetDashboardByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	if category == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Categoria não fornecida"})
		return
	}

	dashboard := DashboardService.GetDashboardByCategory(category)
	ctx.JSON(http.StatusOK, dashboard)
}
