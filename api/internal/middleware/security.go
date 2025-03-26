package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecureHeaders adiciona headers de segurança às respostas HTTP
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Strict Transport Security - força HTTPS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Evita MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Previne ataques de clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Proteção XSS
		c.Header("X-XSS-Protection", "1; mode=block")

		// Define política de origens permitidas para recursos
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")

		// Desativa cache para APIs
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Next()
	}
}

// RedirectHTTP redireciona requisições HTTP para HTTPS
func RedirectHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verifica se a requisição veio via protocolo não seguro
		// Em muitos ambientes de produção, isso é indicado por headers como X-Forwarded-Proto
		if c.Request.Header.Get("X-Forwarded-Proto") != "https" {
			// Redireciona para HTTPS
			secureURL := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(301, secureURL)
			c.Abort()
			return
		}
		c.Next()
	}
}

// CORS configura headers CORS para permitir acessos de outras origens
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-Admin-ID")
		c.Header("Access-Control-Max-Age", "86400") // 24 horas

		// Se for uma requisição OPTIONS (preflight), responda imediatamente
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
