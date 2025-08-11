package utils

import (
	"fmt"
	"testing"
)

func TestGetMachineIdentifier(t *testing.T) {
	identifier, err := GetMachineIdentifier()
	if err != nil {
		t.Fatalf("GetMachineIdentifier() failed: %v", err)
	}

	if len(identifier) != 32 {
		t.Errorf("Expected identifier length 32, got %d", len(identifier))
	}

	// Test consistency - should return same identifier on multiple calls
	identifier2, err := GetMachineIdentifier()
	if err != nil {
		t.Fatalf("GetMachineIdentifier() second call failed: %v", err)
	}

	if string(identifier) != string(identifier2) {
		t.Error("GetMachineIdentifier() should return consistent results")
	}
}

func TestGenerateEncryptionKey(t *testing.T) {
	key, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey() failed: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key))
	}

	// Test consistency - should return same key on multiple calls
	key2, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey() second call failed: %v", err)
	}

	if string(key) != string(key2) {
		t.Error("GenerateEncryptionKey() should return consistent results")
	}
}

func TestEncryptDecryptString(t *testing.T) {
	testCases := []string{
		"",
		"simple text",
		"complex text with symbols !@#$%^&*()",
		"multi\nline\ntext",
		"unicode text: 你好世界 🌍",
		"very long text that exceeds typical block sizes and should test the encryption algorithm properly with various characters and symbols",
	}

	for _, testText := range testCases {
		t.Run("text_"+testText[:min(10, len(testText))], func(t *testing.T) {
			// Encrypt
			encrypted, err := EncryptString(testText)
			if err != nil {
				t.Fatalf("EncryptString() failed: %v", err)
			}

			// Verify encrypted data is different from original (unless empty)
			if len(testText) > 0 && string(encrypted) == testText {
				t.Error("Encrypted data should be different from original")
			}

			// Decrypt
			decrypted, err := DecryptString(encrypted)
			if err != nil {
				t.Fatalf("DecryptString() failed: %v", err)
			}

			// Verify decrypted matches original
			if decrypted != testText {
				t.Errorf("Decrypted text doesn't match original. Expected: %q, Got: %q", testText, decrypted)
			}
		})
	}
}

func TestEncryptDecryptConsistency(t *testing.T) {
	testText := "consistency test"

	// Encrypt the same text multiple times
	encrypted1, err := EncryptString(testText)
	if err != nil {
		t.Fatalf("First encryption failed: %v", err)
	}

	encrypted2, err := EncryptString(testText)
	if err != nil {
		t.Fatalf("Second encryption failed: %v", err)
	}

	// Encrypted results should be different (due to random IV)
	if string(encrypted1) == string(encrypted2) {
		t.Error("Multiple encryptions of same text should produce different ciphertext")
	}

	// But both should decrypt to the same original text
	decrypted1, err := DecryptString(encrypted1)
	if err != nil {
		t.Fatalf("First decryption failed: %v", err)
	}

	decrypted2, err := DecryptString(encrypted2)
	if err != nil {
		t.Fatalf("Second decryption failed: %v", err)
	}

	if decrypted1 != testText || decrypted2 != testText {
		t.Error("Both decryptions should return original text")
	}
}

func TestDecryptStringInvalidInput(t *testing.T) {
	testCases := []struct {
		name       string
		ciphertext []byte
		expectErr  bool
	}{
		{
			name:       "empty input",
			ciphertext: []byte{},
			expectErr:  true,
		},
		{
			name:       "too short input",
			ciphertext: []byte{1, 2, 3},
			expectErr:  true,
		},
		{
			name:       "exactly block size",
			ciphertext: make([]byte, 16), // AES block size
			expectErr:  false, // Should not error but will return empty string
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecryptString(tc.ciphertext)
			if tc.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestBase64EncodeDecode(t *testing.T) {
	testCases := [][]byte{
		{},
		{1, 2, 3, 4, 5},
		[]byte("hello world"),
		[]byte("binary data with null bytes\x00\x01\x02"),
	}

	for i, testData := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			// Encode
			encoded := EncodeBase64(testData)

			// Decode
			decoded, err := DecodeBase64(encoded)
			if err != nil {
				t.Fatalf("DecodeBase64() failed: %v", err)
			}

			// Verify
			if string(decoded) != string(testData) {
				t.Errorf("Decoded data doesn't match original. Expected: %v, Got: %v", testData, decoded)
			}
		})
	}
}

func TestDecodeBase64InvalidInput(t *testing.T) {
	invalidInputs := []string{
		"invalid base64!",
		"not base64 at all",
		"invalid==chars",
	}

	for _, input := range invalidInputs {
		t.Run("invalid_"+input, func(t *testing.T) {
			_, err := DecodeBase64(input)
			if err == nil {
				t.Error("Expected error for invalid base64 input")
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
