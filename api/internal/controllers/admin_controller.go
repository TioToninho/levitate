package controllers

import (
	"net/http"
	"strconv"
	"trackable-donations/api/internal/models"
	"trackable-donations/api/internal/services"

	"github.com/gin-gonic/gin"
)

// AdminService é a instância do serviço de administração
var AdminService *services.AdminService

// SetupAdminService configura o serviço de administração
func SetupAdminService(donationService *services.DonationService, expenseService *services.ExpenseService) {
	AdminService = services.NewAdminService(donationService, expenseService)
}

// RegisterNGO processa o registro de uma nova ONG
func RegisterNGO(ctx *gin.Context) {
	var req models.NGORegistrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar dados do registro de ONG"})
		return
	}

	registration, err := AdminService.RegisterNGO(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, registration)
}

// ValidateCNPJ valida o CNPJ de um registro de ONG
func ValidateCNPJ(ctx *gin.Context) {
	regID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de registro inválido"})
		return
	}

	registration, err := AdminService.ValidateCNPJOnline(uint(regID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, registration)
}

// UploadNGODocuments processa o upload de documentos de uma ONG
func UploadNGODocuments(ctx *gin.Context) {
	regID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de registro inválido"})
		return
	}

	// Limite o upload para 10MB
	file, _, err := ctx.Request.FormFile("documents")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar arquivo: " + err.Error()})
		return
	}
	defer file.Close()

	fileBytes := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			fileBytes = append(fileBytes, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	registration, err := AdminService.UploadNGODocuments(uint(regID), fileBytes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, registration)
}

// ApproveNGO aprova o registro de uma ONG
func ApproveNGO(ctx *gin.Context) {
	regID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de registro inválido"})
		return
	}

	type ApprovalRequest struct {
		AdminID  uint   `json:"admin_id" binding:"required"`
		Comments string `json:"comments"`
	}

	var req ApprovalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar dados da aprovação"})
		return
	}

	ngo, err := AdminService.ApproveNGO(uint(regID), req.AdminID, req.Comments)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ngo)
}

// RejectNGO rejeita o registro de uma ONG
func RejectNGO(ctx *gin.Context) {
	regID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de registro inválido"})
		return
	}

	type RejectionRequest struct {
		AdminID uint   `json:"admin_id" binding:"required"`
		Reason  string `json:"reason" binding:"required"`
	}

	var req RejectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar dados da rejeição"})
		return
	}

	registration, err := AdminService.RejectNGO(uint(regID), req.AdminID, req.Reason)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, registration)
}

// GetNGORegistrations retorna todos os registros de ONGs
func GetNGORegistrations(ctx *gin.Context) {
	registrations := AdminService.GetNGORegistrations()
	ctx.JSON(http.StatusOK, registrations)
}

// GetNGORegistrationByID retorna um registro de ONG pelo ID
func GetNGORegistrationByID(ctx *gin.Context) {
	regID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de registro inválido"})
		return
	}

	registration, err := AdminService.GetNGORegistrationByID(uint(regID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, registration)
}

// GetNGORegistrationsByCNPJ retorna registros de ONGs pelo CNPJ
func GetNGORegistrationsByCNPJ(ctx *gin.Context) {
	cnpj := ctx.Query("cnpj")
	if cnpj == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "CNPJ não fornecido"})
		return
	}

	registrations := AdminService.GetNGORegistrationsByCNPJ(cnpj)
	ctx.JSON(http.StatusOK, registrations)
}

// AuditEntity realiza auditoria em uma entidade
func AuditEntity(ctx *gin.Context) {
	var req models.AuditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar dados da auditoria"})
		return
	}

	// Obter ID do administrador dos headers (em um sistema real, validaria o token)
	adminIDStr := ctx.GetHeader("X-Admin-ID")
	adminID := uint(0)
	if adminIDStr != "" {
		id, err := strconv.ParseUint(adminIDStr, 10, 32)
		if err == nil {
			adminID = uint(id)
		}
	}

	result, err := AdminService.AuditEntity(req, adminID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetAuditLogs retorna os logs de auditoria
func GetAuditLogs(ctx *gin.Context) {
	entityType := ctx.Query("entity_type")
	entityIDStr := ctx.Query("entity_id")

	if entityType != "" && entityIDStr != "" {
		entityID, err := strconv.ParseUint(entityIDStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de entidade inválido"})
			return
		}
		logs := AdminService.GetAuditLogsByEntityID(entityType, uint(entityID))
		ctx.JSON(http.StatusOK, logs)
		return
	}

	if entityType != "" {
		logs := AdminService.GetAuditLogsByEntityType(entityType)
		ctx.JSON(http.StatusOK, logs)
		return
	}

	logs := AdminService.GetAuditLogs()
	ctx.JSON(http.StatusOK, logs)
}
