package services

import (
	"math/rand"
	"time"
)

// Função auxiliar para gerar um hash de transação fictício
func generateMockTransactionHash() string {
	const charset = "abcdef0123456789"
	rand.Seed(time.Now().UnixNano())

	hash := "0x"
	for i := 0; i < 64; i++ {
		hash += string(charset[rand.Intn(len(charset))])
	}

	return hash
}

// Função auxiliar para gerar um hash fictício genérico
func generateMockHash(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	hash := ""
	for i := 0; i < length; i++ {
		hash += string(charset[rand.Intn(len(charset))])
	}

	return hash
}
