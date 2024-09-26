package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCreateHandlerFunc(t *testing.T) {
	// Step 1: Create a temporary file to simulate a response file
	tempFile, err := os.CreateTemp("", "response.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the file after the test

	// Step 2: Write some test data to the file
	expectedResponse := `{"message": "Hello, World!"}`
	if _, err := tempFile.WriteString(expectedResponse); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Close the file so it can be read in the handler
	tempFile.Close()

	// Step 3: Create the handler
	handler := createHandlerFunc(tempFile.Name())

	// Step 4: Create a test request (GET request in this case)
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Step 5: Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Step 6: Call the handler
	handler.ServeHTTP(rr, req)

	// Step 7: Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Step 8: Check the Content-Type header
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("handler returned wrong Content-Type: got %v want %v", contentType, "application/json")
	}

	// Step 9: Check the response body
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned wrong body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}

func TestCreateHandlerFuncFileNotFound(t *testing.T) {
	// Step 1: Create a handler with a non-existent file path
	handler := createHandlerFunc("/non-existent-file.json")

	// Step 2: Create a test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Step 3: Create a response recorder
	rr := httptest.NewRecorder()

	// Step 4: Call the handler
	handler.ServeHTTP(rr, req)

	// Step 5: Check the status code (expecting 500 Internal Server Error)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}
