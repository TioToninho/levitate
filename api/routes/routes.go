package routes

import (
	"trackable-donations/api/internal/controllers"
	"trackable-donations/api/internal/middleware"
	"trackable-donations/api/internal/services"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware middleware para verificar se o usuário é um administrador
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Em um sistema real, verificaria o token JWT para confirmar se é um administrador
		// Aqui, apenas verificamos se existe um header específico
		adminID := c.GetHeader("X-Admin-ID")
		if adminID == "" {
			c.JSON(401, gin.H{"error": "Acesso não autorizado"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SetupRoutes configura todas as rotas da API
func SetupRoutes(router *gin.Engine, publicRateLimiter, adminRateLimiter *middleware.RateLimiter) {
	// Configurar serviços
	donationService := services.NewDonationService()
	controllers.SetupExpenseService(donationService)
	controllers.SetupTransparencyService(donationService, controllers.ExpenseService)
	controllers.SetupAdminService(donationService, controllers.ExpenseService)
	controllers.SetupPublicServices(donationService, controllers.ExpenseService)

	// Rota de verificação de saúde sem rate limiting
	router.GET("/health", controllers.HealthCheck)

	// Rotas públicas com rate limiting
	publicRoutes := router.Group("/")
	publicRoutes.Use(publicRateLimiter.RateLimit())
	{
		// Rotas para ONGs
		publicRoutes.GET("/ngos", controllers.ListNGOs)
		publicRoutes.GET("/ngos/:id", controllers.GetNGOByID)

		// Rotas para doações
		publicRoutes.POST("/donations", controllers.CreateDonation)
		publicRoutes.POST("/donations/:id/confirm-payment", controllers.ConfirmPayment)

		// Rotas para rastreamento de doações
		publicRoutes.GET("/donations/:id/receipt", controllers.GetDonationReceipt)
		publicRoutes.GET("/donations/:id/usages", controllers.GetResourceUsagesByDonation)

		// Rotas para doadores
		publicRoutes.GET("/donors/:id/donations", controllers.GetDonationsByDonor)
		publicRoutes.GET("/donors/:id/dashboard", controllers.GetDonorDashboard)

		// Rotas para despesas
		publicRoutes.POST("/expenses", controllers.RegisterExpense)
		publicRoutes.POST("/expenses/:id/receipt", controllers.UploadReceipt)
		publicRoutes.GET("/expenses/donation/:donationId", controllers.GetExpensesByDonation)
		publicRoutes.GET("/expenses/ngo/:ngoId", controllers.GetExpensesByNGO)

		// Rotas para transparência pública
		publicRoutes.GET("/transparency", controllers.GetPublicDashboard)
		publicRoutes.GET("/transparency/donations", controllers.GetPublicDonations)
		publicRoutes.GET("/transparency/expenses", controllers.GetPublicExpenses)
		publicRoutes.GET("/transparency/ngos", controllers.GetPublicNGOsSummary)
		publicRoutes.GET("/transparency/ngos/:id", controllers.GetPublicNGOSummary)
		publicRoutes.GET("/transparency/ngos/:id/donations", controllers.GetPublicNGODonations)
		publicRoutes.GET("/transparency/ngos/:id/expenses", controllers.GetPublicNGOExpenses)

		// Rotas para explorador de transações
		publicRoutes.GET("/explorer/search", controllers.SearchDonations)
		publicRoutes.GET("/explorer/donations/hash/:hash", controllers.GetDonationByHash)
		publicRoutes.GET("/explorer/donations/:id", controllers.GetDonationByID)
		publicRoutes.GET("/explorer/donations/ngo/:ngo_id", controllers.GetDonationsByNGO)
		publicRoutes.GET("/explorer/donations/recent", controllers.GetRecentDonations)

		// Rotas para dashboard global
		publicRoutes.GET("/dashboard/global", controllers.GetGlobalDashboard)
		publicRoutes.GET("/dashboard/by-date-range", controllers.GetDashboardByDateRange)
		publicRoutes.GET("/dashboard/by-category/:category", controllers.GetDashboardByCategory)

		// Rotas para teste do Swagger
		publicRoutes.GET("/swagger-test", controllers.SwaggerUITest)
	}

	// Rotas para administração (protegidas por middleware e com rate limiting mais restrito)
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(AdminMiddleware())
	adminRoutes.Use(adminRateLimiter.RateLimit())
	{
		// Cadastro e gestão de ONGs
		adminRoutes.POST("/ngos/register", controllers.RegisterNGO)
		adminRoutes.POST("/ngos/registration/:id/validate-cnpj", controllers.ValidateCNPJ)
		adminRoutes.POST("/ngos/registration/:id/upload-documents", controllers.UploadNGODocuments)
		adminRoutes.POST("/ngos/registration/:id/approve", controllers.ApproveNGO)
		adminRoutes.POST("/ngos/registration/:id/reject", controllers.RejectNGO)
		adminRoutes.GET("/ngos/registrations", controllers.GetNGORegistrations)
		adminRoutes.GET("/ngos/registrations/:id", controllers.GetNGORegistrationByID)
		adminRoutes.GET("/ngos/registrations/by-cnpj", controllers.GetNGORegistrationsByCNPJ)

		// Auditoria
		adminRoutes.POST("/audit", controllers.AuditEntity)
		adminRoutes.GET("/audit/logs", controllers.GetAuditLogs)
	}
}
