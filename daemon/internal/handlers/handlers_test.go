package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	// Create a new app instance
	app := &App{
		StartTime:      time.Now(),
		DataManager:    nil, // Not needed for health check
		DnsmasqManager: nil, // Not needed for health check
	}

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.HealthHandler)

	// Perform the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("could not unmarshal response: %v", err)
	}

	// Verify response fields
	if response["status"] != "healthy" {
		t.Errorf("expected status to be 'healthy', got %s", response["status"])
	}

	if response["service"] != "wild-cloud-central" {
		t.Errorf("expected service to be 'wild-cloud-central', got %s", response["service"])
	}

	// Check Content-Type header
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}
}

func TestStatusHandler(t *testing.T) {
	// Create a new app instance with test time
	testStartTime := time.Now().Add(-1 * time.Hour) // Started 1 hour ago
	app := &App{
		StartTime:      testStartTime,
		DataManager:    nil, // Not needed for status check
		DnsmasqManager: nil, // Not needed for status check
	}

	// Create a request
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.StatusHandler)

	// Perform the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("could not unmarshal response: %v", err)
	}

	// Verify response fields
	if response["status"] != "running" {
		t.Errorf("expected status to be 'running', got %v", response["status"])
	}

	if response["version"] != "1.0.0" {
		t.Errorf("expected version to be '1.0.0', got %v", response["version"])
	}

	// Check that uptime exists and is a string
	if _, ok := response["uptime"].(string); !ok {
		t.Errorf("expected uptime to be a string, got %T", response["uptime"])
	}

	// Check Content-Type header
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}
}