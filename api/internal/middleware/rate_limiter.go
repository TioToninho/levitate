package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implementa limitação de requisições por IP
type RateLimiter struct {
	sync.Mutex
	ipLimits     map[string][]time.Time
	maxRequests  int
	windowLength time.Duration
	enabled      bool
}

// NewRateLimiter cria um novo limitador de requisições
func NewRateLimiter(maxRequests int, windowLength time.Duration) *RateLimiter {
	return &RateLimiter{
		ipLimits:     make(map[string][]time.Time),
		maxRequests:  maxRequests,
		windowLength: windowLength,
		enabled:      true,
	}
}

// RateLimit retorna um middleware Gin para limitar requisições
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Se o limitador estiver desativado, apenas continue
		if !rl.enabled {
			c.Next()
			return
		}

		ip := c.ClientIP()

		rl.Lock()
		defer rl.Unlock()

		// Remover requisições antigas do período de janela
		now := time.Now()
		validTime := now.Add(-rl.windowLength)

		if _, exists := rl.ipLimits[ip]; exists {
			var validRequests []time.Time
			for _, t := range rl.ipLimits[ip] {
				if t.After(validTime) {
					validRequests = append(validRequests, t)
				}
			}
			rl.ipLimits[ip] = validRequests
		} else {
			rl.ipLimits[ip] = []time.Time{}
		}

		// Verificar limite
		if len(rl.ipLimits[ip]) >= rl.maxRequests {
			// Adicionar headers para informar cliente sobre limites
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.maxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			resetTime := validTime.Add(rl.windowLength)
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Limite de requisições excedido. Tente novamente mais tarde.",
			})
			return
		}

		// Registrar requisição
		rl.ipLimits[ip] = append(rl.ipLimits[ip], now)

		// Adicionar headers informativos
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.maxRequests-len(rl.ipLimits[ip])))

		c.Next()
	}
}

// SetEnabled ativa ou desativa o limitador (útil para ambientes de desenvolvimento)
func (rl *RateLimiter) SetEnabled(enabled bool) {
	rl.Lock()
	defer rl.Unlock()
	rl.enabled = enabled
}

// GetLimits retorna informações sobre os limites (útil para debugging)
func (rl *RateLimiter) GetLimits() map[string]int {
	rl.Lock()
	defer rl.Unlock()

	limits := make(map[string]int)
	for ip, requests := range rl.ipLimits {
		limits[ip] = len(requests)
	}

	return limits
}
