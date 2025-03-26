package controllers

import (
	"io"
	"net/http"
	"strconv"
	"trackable-donations/api/internal/models"
	"trackable-donations/api/internal/services"

	"github.com/gin-gonic/gin"
)

// ExpenseService é a instância do serviço de despesas
var ExpenseService *services.ExpenseService

// SetupExpenseService configura o serviço de despesas
func SetupExpenseService(donationService *services.DonationService) {
	ExpenseService = services.NewExpenseService(donationService)
}

// RegisterExpense registra uma nova despesa
// @Summary Registrar despesa
// @Description Registra uma nova despesa vinculada a uma doação
// @Tags Despesas
// @Accept json
// @Produce json
// @Param despesa body models.ExpenseRequest true "Dados da despesa"
// @Success 201 {object} models.ExpenseResponse
// @Failure 400 {object} map[string]string "Erro nos dados da despesa"
// @Router /expenses [post]
func RegisterExpense(ctx *gin.Context) {
	var expenseReq models.ExpenseRequest
	if err := ctx.ShouldBindJSON(&expenseReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar dados da despesa"})
		return
	}

	response, err := ExpenseService.RegisterExpense(expenseReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// UploadReceipt faz upload de um comprovante para uma despesa
// @Summary Fazer upload de comprovante
// @Description Envia o comprovante/recibo de uma despesa
// @Tags Despesas
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID da despesa"
// @Param receipt formData file true "Arquivo do comprovante (PDF, JPG, PNG)"
// @Success 200 {object} models.ExpenseResponse
// @Failure 400 {object} map[string]string "ID de despesa inválido ou erro no arquivo"
// @Router /expenses/{id}/receipt [post]
func UploadReceipt(ctx *gin.Context) {
	expenseID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de despesa inválido"})
		return
	}

	// Limite o upload para 10MB
	file, _, err := ctx.Request.FormFile("receipt")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar arquivo: " + err.Error()})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler o arquivo: " + err.Error()})
		return
	}

	response, err := ExpenseService.UploadReceipt(uint(expenseID), fileBytes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetExpensesByDonation retorna as despesas relacionadas a uma doação específica
// @Summary Listar despesas por doação
// @Description Retorna todas as despesas relacionadas a uma doação específica
// @Tags Despesas
// @Accept json
// @Produce json
// @Param donationId path int true "ID da doação"
// @Success 200 {array} models.Expense
// @Failure 400 {object} map[string]string "ID de doação inválido"
// @Router /expenses/donation/{donationId} [get]
func GetExpensesByDonation(ctx *gin.Context) {
	donationID, err := strconv.ParseUint(ctx.Param("donationId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de doação inválido"})
		return
	}

	expenses, err := ExpenseService.GetExpensesByDonation(uint(donationID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, expenses)
}

// GetExpensesByNGO retorna todas as despesas de uma ONG
// @Summary Listar despesas por ONG
// @Description Retorna todas as despesas registradas por uma ONG específica
// @Tags Despesas
// @Accept json
// @Produce json
// @Param ngoId path int true "ID da ONG"
// @Success 200 {array} models.Expense
// @Failure 400 {object} map[string]string "ID de ONG inválido"
// @Router /expenses/ngo/{ngoId} [get]
func GetExpensesByNGO(ctx *gin.Context) {
	ngoID, err := strconv.ParseUint(ctx.Param("ngoId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de ONG inválido"})
		return
	}

	expenses, err := ExpenseService.GetExpensesByNGO(uint(ngoID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, expenses)
}
