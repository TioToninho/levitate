package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthStatus representa o status de saúde da API
type HealthStatus struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
}

var startTime = time.Now()

// HealthCheck verifica o status de saúde da API
// @Summary Verificar saúde da API
// @Description Verifica se a API está funcionando corretamente
// @Tags Sistema
// @Accept json
// @Produce json
// @Success 200 {object} controllers.HealthStatus
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	// Obter versão da variável de ambiente ou usar padrão
	version := os.Getenv("API_VERSION")
	if version == "" {
		version = "1.0.0"
	}

	// Calcular tempo de atividade
	uptime := time.Since(startTime).String()

	// Criar resposta
	status := HealthStatus{
		Status:    "online",
		Version:   version,
		Timestamp: time.Now(),
		Uptime:    uptime,
	}

	c.JSON(http.StatusOK, status)
}
