package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// EncryptAES encrypts plaintext using AES-256-CBC with a passphrase
// Compatible with CryptoJS encryption format used by Supabase Studio
func EncryptAES(plaintext, passphrase string) (string, error) {
	// Generate random salt (8 bytes)
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key and IV using EVP_BytesToKey (MD5-based, compatible with CryptoJS)
	key, iv := deriveKeyAndIV(passphrase, salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Pad plaintext to block size
	plaintextBytes := pkcs7Pad([]byte(plaintext), aes.BlockSize)

	// Encrypt using CBC mode
	ciphertext := make([]byte, len(plaintextBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintextBytes)

	// Format: "Salted__" + salt + ciphertext, then base64 encode
	// This matches CryptoJS format
	result := make([]byte, 8+len(salt)+len(ciphertext))
	copy(result[0:8], []byte("Salted__"))
	copy(result[8:16], salt)
	copy(result[16:], ciphertext)

	return base64.StdEncoding.EncodeToString(result), nil
}

// deriveKeyAndIV derives encryption key and IV using MD5 (EVP_BytesToKey algorithm)
// This matches CryptoJS's key derivation method
func deriveKeyAndIV(passphrase string, salt []byte) (key []byte, iv []byte) {
	// We need 32 bytes for key (AES-256) + 16 bytes for IV = 48 bytes total
	// MD5 produces 16 bytes, so we need 3 rounds
	var derivedKey []byte
	var block []byte

	for len(derivedKey) < 48 {
		hash := md5.New()
		hash.Write(block)
		hash.Write([]byte(passphrase))
		hash.Write(salt)
		block = hash.Sum(nil)
		derivedKey = append(derivedKey, block...)
	}

	return derivedKey[:32], derivedKey[32:48]
}

// pkcs7Pad adds PKCS#7 padding to data
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

// DecryptAES decrypts ciphertext using AES-256-CBC with a passphrase
// Compatible with CryptoJS decryption format
func DecryptAES(ciphertext, passphrase string) (string, error) {
	// Decode base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Check for "Salted__" prefix
	if len(data) < 16 || string(data[:8]) != "Salted__" {
		return "", fmt.Errorf("invalid ciphertext format")
	}

	// Extract salt
	salt := data[8:16]
	encrypted := data[16:]

	// Derive key and IV
	key, iv := deriveKeyAndIV(passphrase, salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Decrypt using CBC mode
	if len(encrypted)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of block size")
	}

	plaintext := make([]byte, len(encrypted))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, encrypted)

	// Remove PKCS#7 padding
	plaintext, err = pkcs7Unpad(plaintext)
	if err != nil {
		return "", fmt.Errorf("failed to unpad: %w", err)
	}

	return string(plaintext), nil
}

// pkcs7Unpad removes PKCS#7 padding from data
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding: empty data")
	}

	padding := int(data[length-1])
	if padding > length || padding > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding size")
	}

	// Verify padding
	for i := length - padding; i < length; i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding bytes")
		}
	}

	return data[:length-padding], nil
}
