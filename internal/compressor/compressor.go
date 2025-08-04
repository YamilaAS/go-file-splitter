// internal/compressor/compressor.go
package compressor

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// CompressFile comprime un archivo con gzip y devuelve el nuevo nombre del archivo.
func CompressFile(sourcePath string) (string, error) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destPath := fmt.Sprintf("%s.gz", sourcePath)
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	gzipWriter := gzip.NewWriter(destFile)
	defer gzipWriter.Close()

	if _, err := io.Copy(gzipWriter, sourceFile); err != nil {
		return "", fmt.Errorf("failed to compress file: %w", err)
	}

	return destPath, nil
}
