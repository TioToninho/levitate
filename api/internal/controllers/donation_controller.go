package controllers

import (
	"net/http"
	"strconv"
	"trackable-donations/api/internal/models"
	"trackable-donations/api/internal/services"
	"trackable-donations/api/internal/utils"

	"github.com/gin-gonic/gin"
)

var donationService = services.NewDonationService()

// ListNGOs lista todas as ONGs disponíveis
// @Summary Listar ONGs
// @Description Retorna a lista de todas as ONGs cadastradas
// @Tags NGOs
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]models.NGO
// @Router /ngos [get]
func ListNGOs(c *gin.Context) {
	ngos := donationService.GetAllNGOs()
	c.JSON(http.StatusOK, gin.H{
		"data": ngos,
	})
}

// GetNGOByID retorna uma ONG específica pelo ID
// @Summary Obter ONG por ID
// @Description Retorna os detalhes de uma ONG específica pelo ID
// @Tags NGOs
// @Accept json
// @Produce json
// @Param id path int true "ID da ONG"
// @Success 200 {object} map[string]models.NGO
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "ONG não encontrada"
// @Router /ngos/{id} [get]
func GetNGOByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	ngo, err := donationService.GetNGOByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ngo})
}

// CreateDonation processa uma nova doação
// @Summary Criar doação
// @Description Registra uma nova doação no sistema
// @Tags Doações
// @Accept json
// @Produce json
// @Param doacao body models.DonationRequest true "Dados da doação"
// @Success 201 {object} map[string]models.DonationResponse
// @Failure 400 {object} map[string]string "Erro nos dados ou documento inválido"
// @Router /donations [post]
func CreateDonation(c *gin.Context) {
	var req models.DonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar se a requisição contém o documento do doador (CPF/CNPJ)
	if req.DonorDocument != "" {
		// Validar o formato do documento
		if len(req.DonorDocument) == 11 || len(req.DonorDocument) == 14 ||
			utils.ValidateCPF(req.DonorDocument) || utils.ValidateCNPJ(req.DonorDocument) {
			// Anonimizar o documento usando hash SHA-256
			req.DonorDocument = utils.HashSensitiveData(req.DonorDocument, false)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de documento inválido"})
			return
		}
	}

	// Se tiver outros dados sensíveis, anonimizar aqui também

	response, err := donationService.ProcessDonation(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// ConfirmPayment simula a confirmação de um pagamento
// @Summary Confirmar pagamento
// @Description Confirma o pagamento de uma doação e gera comprovante
// @Tags Doações
// @Accept json
// @Produce json
// @Param id path int true "ID da doação"
// @Success 200 {object} map[string]models.DonationResponse
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Doação não encontrada"
// @Router /donations/{id}/confirm-payment [post]
func ConfirmPayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	response, err := donationService.MockPaymentConfirmation(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetDonationsByDonor retorna todas as doações de um doador
// @Summary Listar doações por doador
// @Description Retorna todas as doações realizadas por um doador específico
// @Tags Doações
// @Accept json
// @Produce json
// @Param id path int true "ID do doador"
// @Success 200 {object} map[string][]models.Donation
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Doador não encontrado"
// @Router /donors/{id}/donations [get]
func GetDonationsByDonor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	donations, err := donationService.GetDonationsByDonorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": donations})
}

// GetDonationReceipt retorna o comprovante de uma doação
// @Summary Obter comprovante de doação
// @Description Retorna o comprovante de uma doação específica
// @Tags Doações
// @Accept json
// @Produce json
// @Param id path int true "ID da doação"
// @Success 200 {object} map[string]models.DonationReceipt
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Comprovante não encontrado"
// @Router /donations/{id}/receipt [get]
func GetDonationReceipt(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	receipt, err := donationService.GetDonationReceipt(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": receipt})
}

// GetResourceUsagesByDonation retorna os usos dos recursos de uma doação
// @Summary Obter usos dos recursos de doação
// @Description Retorna os registros de uso dos recursos de uma doação específica
// @Tags Doações
// @Accept json
// @Produce json
// @Param id path int true "ID da doação"
// @Success 200 {object} map[string][]models.ResourceUsage
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Doação não encontrada"
// @Router /donations/{id}/usages [get]
func GetResourceUsagesByDonation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	usages, err := donationService.GetResourceUsagesByDonationID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": usages})
}

// GetDonorDashboard retorna o dashboard de um doador
// @Summary Obter dashboard do doador
// @Description Retorna dados consolidados do dashboard de um doador específico
// @Tags Doações
// @Accept json
// @Produce json
// @Param id path int true "ID do doador"
// @Success 200 {object} map[string]models.DonorDashboard
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Doador não encontrado"
// @Router /donors/{id}/dashboard [get]
func GetDonorDashboard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	dashboard, err := donationService.GetDonorDashboard(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dashboard})
}
