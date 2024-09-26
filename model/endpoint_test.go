package model

import (
	"os"
	"testing"
)

// Test ReadEndpoints with a valid file
func TestReadEndpoints(t *testing.T) {
	// Step 1: Create a temporary file with valid data
	file, err := os.CreateTemp("", "mockzz-endpoints.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name()) // Clean up the file after the test

	// Step 2: Write valid content to the file
	content := "GET /api1:response1.json\nPOST /api2:response2.json\n"
	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Close the file to ensure it's flushed to disk
	file.Close()

	// Step 3: Call the function and check the result
	endpoints, err := ReadEndpoints(file.Name())
	if err != nil {
		t.Fatalf("ReadEndpoints failed: %v", err)
	}

	// Step 4: Validate the returned map
	expected := map[string]string{
		"GET /api1":  "response1.json",
		"POST /api2": "response2.json",
	}

	for key, value := range expected {
		if endpoints[key] != value {
			t.Errorf("Expected %s -> %s, got %s -> %s", key, value, key, endpoints[key])
		}
	}
}

// Test ReadEndpoints with a non-existent file
func TestReadEndpointsFileNotFound(t *testing.T) {
	_, err := ReadEndpoints("non-existent-file.txt")
	if err == nil {
		t.Error("Expected an error for a non-existent file, got nil")
	}
}

// Test ReadEndpoints with a mix of valid and invalid content
func TestReadEndpointsInvalidContent(t *testing.T) {
	// Step 1: Create a temporary file with both valid and invalid content
	file, err := os.CreateTemp("", "mockzz-invalid-endpoints.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name()) // Clean up the file after the test

	// Step 2: Write a mix of valid and invalid content to the file
	content := `
GET /api1:response1.json
INVALIDLINE
POST /api2 response2.json
GET /api3:response3.json
 : invalid-entry
`
	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	file.Close()

	// Step 3: Call ReadEndpoints to process the file
	endpoints, err := ReadEndpoints(file.Name())
	if err != nil {
		t.Fatalf("ReadEndpoints failed: %v", err)
	}

	// Step 4: Validate the returned map - it should only include valid lines
	expected := map[string]string{
		"GET /api1": "response1.json",
		"GET /api3": "response3.json",
	}

	if len(endpoints) != len(expected) {
		t.Errorf("Expected %d valid endpoints, got %d", len(expected), len(endpoints))
	}

	for key, value := range expected {
		if endpoints[key] != value {
			t.Errorf("Expected %s -> %s, got %s -> %s", key, value, key, endpoints[key])
		}
	}
}
