//go:build integration

package test_integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestHealthEndpoint(t *testing.T) {
	healthURL := "http://localhost:8080/health"
	
	time.Sleep(3 * time.Second)

	resp, err := http.Get(healthURL)
	if err != nil {
		t.Fatalf("Health endpoint unavailable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var healthStatus map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&healthStatus)

	if healthStatus["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%v'", healthStatus["status"])
	}

	t.Logf("✅ Health check passed: %+v", healthStatus)
}
