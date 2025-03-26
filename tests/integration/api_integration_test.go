package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Fatalf("Erro ao chamar o endpoint: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode, "Status code deve ser 200")
}
