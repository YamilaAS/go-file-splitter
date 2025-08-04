// internal/hasher/hasher.go
package hasher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// CalculateSHA256 calcula el checksum SHA-256 de un archivo.
func CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file for hashing: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
