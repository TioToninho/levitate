package main

import (
	"log"
	"os"
	"time"

	_ "trackable-donations/api/docs" // Importar documentação Swagger
	"trackable-donations/api/internal/middleware"
	"trackable-donations/api/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API de Doações Rastreáveis
// @version 1.0
// @description API para gerenciamento de doações rastreáveis com blockchain e IPFS

// @contact.name Suporte API
// @contact.email suporte@levitate.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey AdminAuth
// @in header
// @name X-Admin-ID
// @description Chave de autenticação para rotas administrativas

func main() {
	// Em produção, usar modo "release"
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configurar middlewares de segurança
	router.Use(middleware.CORS())
	router.Use(middleware.SecureHeaders())

	// Redirecionar HTTP para HTTPS (apenas em produção)
	if os.Getenv("ENV") == "production" {
		router.Use(middleware.RedirectHTTP())
	}

	// Aplicar rate limiting em rotas públicas
	publicRateLimiter := middleware.NewRateLimiter(100, 1*time.Minute)

	// Aplicar rate limiting mais restrito em rotas de admin
	adminRateLimiter := middleware.NewRateLimiter(30, 1*time.Minute) // 30 requisições por minuto

	// Configurar rotas com rate limiting
	routes.SetupRoutes(router, publicRateLimiter, adminRateLimiter)

	// Determinar porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configuração simplificada do Swagger - isso deve resolver o problema
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Verificar certificados SSL em produção
	certFile := os.Getenv("SSL_CERT_FILE")
	keyFile := os.Getenv("SSL_KEY_FILE")

	// Iniciar o servidor com SSL em produção ou HTTP em desenvolvimento
	log.Printf("Documentação Swagger disponível em http://localhost:%s/swagger/index.html", port)

	if os.Getenv("ENV") == "production" && certFile != "" && keyFile != "" {
		log.Printf("Servidor iniciando em modo seguro (HTTPS) na porta %s...", port)
		if err := router.RunTLS(":"+port, certFile, keyFile); err != nil {
			log.Fatalf("Falha ao iniciar servidor HTTPS: %v", err)
		}
	} else {
		log.Printf("Servidor iniciando em modo HTTP na porta %s...", port)
		log.Println("AVISO: Em produção, configure SSL_CERT_FILE e SSL_KEY_FILE para usar HTTPS")
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Falha ao iniciar servidor HTTP: %v", err)
		}
	}
}
