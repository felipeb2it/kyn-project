// KynService
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type EncryptDecryptResponse struct {
	EncryptedValue []byte `json:"encrypted_value,omitempty"`
	DecryptedValue string `json:"decrypted_value,omitempty"`
}

var (
	key       = randBytes(256 / 8)
	gcm       cipher.AEAD
	nonceSize int
)

// GCM initialization for encrypt and decrypt on program start.
func init() {
	log.Println("Initializing Kyndryl Service")
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Error reading key: %s\n", err.Error())
		os.Exit(1)
	}

	gcm, err = cipher.NewGCM(block)
	if err != nil {
		log.Printf("Error initializing AEAD: %s\n", err.Error())
		os.Exit(1)
	}

	nonceSize = gcm.NonceSize()
}

func randBytes(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

func encrypt(plaintext []byte) (ciphertext []byte, err error) {
	nonce := randBytes(nonceSize)
	c := gcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, c...), nil
}

func decrypt(ciphertext []byte) (plaintext []byte, err error) {
	if len(ciphertext) < nonceSize {
		errMsg := fmt.Sprintf("Ciphertext too short: %v", err)
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	nonce := ciphertext[0:nonceSize]
	msg := ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, msg, nil)
}

func handleEncryptRequest(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var value string
	var jsonData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonData); err == nil {
		if len(jsonData) != 1 {
			http.Error(w, "JSON object must contain exactly one key-value pair", http.StatusBadRequest)
			return
		}

		value = string(bodyBytes)

	} else if err := json.Unmarshal(bodyBytes, &value); err != nil {
		http.Error(w, "Invalid JSON: Expected a single string value or an object with a single value", http.StatusBadRequest)
		return
	}

	encryptedValue, err := encrypt([]byte(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(EncryptDecryptResponse{
		EncryptedValue: encryptedValue,
	})
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func handleDecryptRequest(w http.ResponseWriter, r *http.Request) {
	var req string
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(req)
	if err != nil {
		log.Println("Error decoding string:", err)
		return
	}

	decryptedValue, err := decrypt(decodedBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(EncryptDecryptResponse{
		DecryptedValue: string(decryptedValue),
	})
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/encrypt", handleEncryptRequest)
	r.HandleFunc("/api/decrypt", handleDecryptRequest)

	log.Println("Service ready for queries!")
	log.Fatal(http.ListenAndServe(":8080", r))
}
