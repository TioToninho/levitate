package controllers

import (
	"net/http"
	"strconv"
	"trackable-donations/api/internal/services"

	"github.com/gin-gonic/gin"
)

// TransparencyService é a instância do serviço de transparência
var TransparencyService *services.TransparencyService

// SetupTransparencyService configura o serviço de transparência
func SetupTransparencyService(donationService *services.DonationService, expenseService *services.ExpenseService) {
	TransparencyService = services.NewTransparencyService(donationService, expenseService)
}

// GetPublicDashboard retorna o dashboard público de transparência
func GetPublicDashboard(ctx *gin.Context) {
	dashboard := TransparencyService.GetTransparencyDashboard()
	ctx.JSON(http.StatusOK, dashboard)
}

// GetPublicDonations retorna todas as doações públicas
func GetPublicDonations(ctx *gin.Context) {
	donations := TransparencyService.GetPublicDonations()
	ctx.JSON(http.StatusOK, donations)
}

// GetPublicExpenses retorna todas as despesas públicas
func GetPublicExpenses(ctx *gin.Context) {
	expenses := TransparencyService.GetPublicExpenses()
	ctx.JSON(http.StatusOK, expenses)
}

// GetPublicNGOsSummary retorna um resumo de todas as ONGs
func GetPublicNGOsSummary(ctx *gin.Context) {
	summaries := TransparencyService.GetAllNGOsSummary()
	ctx.JSON(http.StatusOK, summaries)
}

// GetPublicNGOSummary retorna um resumo de uma ONG específica
func GetPublicNGOSummary(ctx *gin.Context) {
	ngoID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de ONG inválido"})
		return
	}

	summary, err := TransparencyService.GetNGOSummary(uint(ngoID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// GetPublicNGODonations retorna todas as doações de uma ONG específica
func GetPublicNGODonations(ctx *gin.Context) {
	ngoID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de ONG inválido"})
		return
	}

	donations, err := TransparencyService.GetDonationsByNGO(uint(ngoID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, donations)
}

// GetPublicNGOExpenses retorna todas as despesas de uma ONG específica
func GetPublicNGOExpenses(ctx *gin.Context) {
	ngoID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de ONG inválido"})
		return
	}

	expenses, err := TransparencyService.GetExpensesByNGO(uint(ngoID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, expenses)
}
