package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	cpfRegex  = regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	cnpjRegex = regexp.MustCompile(`^\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}$`)
)

// HashSensitiveData aplica SHA-256 com salt em dados sensíveis como CPF/CNPJ
func HashSensitiveData(data string, prefixaConsulta bool) string {
	// Se o dado estiver vazio, retorne vazio
	if data == "" {
		return ""
	}

	// Remover caracteres não numéricos para uniformização
	cleanData := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(data, ".", ""), "-", ""), "/", "")

	// Obter salt da variável de ambiente ou usar valor padrão
	salt := os.Getenv("HASH_SALT")
	if salt == "" {
		salt = "levitate-default-salt" // Em produção, usar um valor mais seguro
		log.Println("AVISO: HASH_SALT não está definido, usando salt padrão. NÃO use em produção!")
	}

	// Se for para manter o prefixo (útil para consultas parciais)
	prefix := ""
	if prefixaConsulta {
		// Vamos manter apenas os 3 primeiros dígitos para CPF ou 4 para CNPJ
		if cpfRegex.MatchString(data) || len(cleanData) == 11 {
			prefix = cleanData[:3]
		} else if cnpjRegex.MatchString(data) || len(cleanData) == 14 {
			prefix = cleanData[:4]
		}
	}

	// Combina os dados com salt para evitar ataques de tabela arco-íris
	combined := cleanData + salt

	// Aplica o algoritmo SHA-256
	hash := sha256.Sum256([]byte(combined))

	// Converte para string hexadecimal
	hashString := hex.EncodeToString(hash[:])

	// Retorna o prefixo + hash, ou apenas o hash
	if prefixaConsulta && prefix != "" {
		return prefix + "-" + hashString
	}
	return hashString
}

// ValidateCPF verifica se o formato do CPF está correto antes de anonimizar
func ValidateCPF(cpf string) bool {
	return cpfRegex.MatchString(cpf) ||
		(len(cpf) == 11 && regexp.MustCompile(`^\d{11}$`).MatchString(cpf))
}

// ValidateCNPJ verifica se o formato do CNPJ está correto antes de anonimizar
func ValidateCNPJ(cnpj string) bool {
	return cnpjRegex.MatchString(cnpj) ||
		(len(cnpj) == 14 && regexp.MustCompile(`^\d{14}$`).MatchString(cnpj))
}
