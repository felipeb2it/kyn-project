// KynService_test.go
package main

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test for randBytes function.
func TestRandBytes(t *testing.T) {
	length := 10
	bytes := randBytes(length)
	if len(bytes) != length {
		t.Errorf("Expected length %d, but got %d", length, len(bytes))
	}
}

// Test for encrypt and decrypt functions.
func TestEncryptDecrypt(t *testing.T) {
	plaintext := []byte("Hello, world!")
	ciphertext, err := encrypt(plaintext)
	if err != nil {
		t.Errorf("Error during encryption: %v", err)
	}

	decrypted, err := decrypt(ciphertext)
	if err != nil {
		t.Errorf("Error during decryption: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Expected %s, but got %s", plaintext, decrypted)
	}
}

// Test for handleEncryptRequest function.
func TestHandleEncryptRequest(t *testing.T) {
	validJSON := `{"data":"Hello, world!"}`
	invalidJSON := `{"data1":"Hello, world!", "data2":"Hello again!"}`
	req, _ := http.NewRequest("POST", "/api/encrypt", bytes.NewBufferString(validJSON))
	rec := httptest.NewRecorder()
	handleEncryptRequest(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	req, _ = http.NewRequest("POST", "/api/encrypt", bytes.NewBufferString(invalidJSON))
	rec = httptest.NewRecorder()
	handleEncryptRequest(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

// Test for handleDecryptRequest function.
func TestHandleDecryptRequest(t *testing.T) {
	data := "Hello, world!"
	encryptedData, _ := encrypt([]byte(data))
	encodedData := base64.StdEncoding.EncodeToString(encryptedData)
	req, _ := http.NewRequest("POST", "/api/decrypt", bytes.NewBufferString(encodedData))
	rec := httptest.NewRecorder()
	handleDecryptRequest(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	invalidData := base64.StdEncoding.EncodeToString([]byte("Invalid data"))
	req, _ = http.NewRequest("POST", "/api/decrypt", bytes.NewBufferString(invalidData))
	rec = httptest.NewRecorder()
	handleDecryptRequest(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
