package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func SplitFile(sourcePath string, chunkSize int64) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	dir := filepath.Dir(sourcePath)
	base := filepath.Base(sourcePath)

	var partNumber int
	for {
		partNumber++
		partName := fmt.Sprintf("%s.part%03d", base, partNumber)
		partPath := filepath.Join(dir, partName)

		partFile, err := os.Create(partPath)
		if err != nil {
			return fmt.Errorf("failed to create part file: %w", err)
		}
		defer partFile.Close()

		// Copiar exactamente chunkSize bytes del origen al destino
		n, err := io.CopyN(partFile, sourceFile, chunkSize)
		if err != nil {
			if err == io.EOF {
				// Fin del archivo, el último trozo puede ser menor
				fmt.Printf("Created final part: %s (%d bytes)\n", partPath, n)
				return nil
			}
			return fmt.Errorf("failed to copy chunk: %w", err)
		}

		fmt.Printf("Created part: %s (%d bytes)\n", partPath, n)
	}
}

func JoinFiles(prefixPath string) error {
	// Crear el archivo de destino
	destFile, err := os.Create(prefixPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Buscar los archivos de los trozos
	partFiles, err := filepath.Glob(fmt.Sprintf("%s.part*", prefixPath))
	if err != nil {
		return fmt.Errorf("failed to glob part files: %w", err)
	}

	if len(partFiles) == 0 {
		return fmt.Errorf("no part files found for prefix: %s", prefixPath)
	}

	// Unir los archivos
	for _, partPath := range partFiles {
		partFile, err := os.Open(partPath)
		if err != nil {
			return fmt.Errorf("failed to open part file %s: %w", partPath, err)
		}
		defer partFile.Close()

		// Copiar el contenido del trozo al archivo de destino
		if _, err := io.Copy(destFile, partFile); err != nil {
			return fmt.Errorf("failed to copy part file %s: %w", partPath, err)
		}
		fmt.Printf("Joined part: %s\n", partPath)
	}

	return nil
}
