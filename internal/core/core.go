package core

import (
	"fmt"
	"go-file-splitter/internal/hasher"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SplitFile divide un archivo grande en trozos más pequeños.
func SplitFile(sourcePath string, chunkSize int64) error {
	// Paso 2: Calcular el hash del archivo original
	originalHash, err := hasher.CalculateSHA256(sourcePath)
	if err != nil {
		return fmt.Errorf("error al calcular el hash original: %w", err)
	}
	fmt.Printf("Hash original calculado: %s\n", originalHash)

	// Crear la estructura de directorios
	baseName := filepath.Base(sourcePath)
	projectDir := fmt.Sprintf("%s_split", strings.TrimSuffix(baseName, filepath.Ext(baseName)))
	splitDir := filepath.Join(projectDir, "split")

	// Paso 3: Crear las carpetas
	if err := os.MkdirAll(splitDir, 0755); err != nil {
		return fmt.Errorf("error al crear las carpetas del proyecto: %w", err)
	}

	// Paso 4: Guardar el hash en un archivo de texto
	hashFilePath := filepath.Join(projectDir, "original_hash.txt")
	if err := os.WriteFile(hashFilePath, []byte(originalHash), 0644); err != nil {
		return fmt.Errorf("error al guardar el hash en el archivo: %w", err)
	}
	fmt.Printf("Hash original guardado en: %s\n", hashFilePath)

	// Abrir el archivo de origen para la división
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo de origen: %w", err)
	}
	defer sourceFile.Close()

	// Lógica de división
	var partNumber int
	for {
		partNumber++
		partName := fmt.Sprintf("%s.part%03d", baseName, partNumber)
		// Paso 5: Los archivos de las partes se guardan en la nueva subcarpeta
		partPath := filepath.Join(splitDir, partName)

		partFile, err := os.Create(partPath)
		if err != nil {
			return fmt.Errorf("error al crear el archivo de la parte: %w", err)
		}
		defer partFile.Close()

		n, err := io.CopyN(partFile, sourceFile, chunkSize)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Parte final creada: %s (%d bytes)\n", partPath, n)
				break
			}
			return fmt.Errorf("error al copiar el trozo del archivo: %w", err)
		}

		fmt.Printf("Parte creada: %s (%d bytes)\n", partPath, n)
	}

	return nil
}

// JoinFiles une los trozos de un archivo y verifica su integridad con un hash.
func JoinFiles(sourcePath string) error {
	baseName := filepath.Base(sourcePath)
	projectDir := fmt.Sprintf("%s_split", strings.TrimSuffix(baseName, filepath.Ext(baseName)))
	splitDir := filepath.Join(projectDir, "split")

	// 1. Leer el hash original del archivo
	hashFilePath := filepath.Join(projectDir, "original_hash.txt")
	originalHashBytes, err := os.ReadFile(hashFilePath)
	if err != nil {
		return fmt.Errorf("error al leer el hash original: %w", err)
	}
	originalHash := string(originalHashBytes)

	// 2. Buscar las partes dentro de la carpeta 'split'
	parts, err := filepath.Glob(filepath.Join(splitDir, fmt.Sprintf("%s.part*", baseName)))
	if err != nil {
		return fmt.Errorf("error al buscar las partes: %w", err)
	}
	if len(parts) == 0 {
		return fmt.Errorf("no se encontraron archivos de partes para el prefijo: %s", baseName)
	}

	sort.Strings(parts) // Es crucial ordenar las partes para unirlas correctamente

	// 3. Unir las partes
	destFile, err := os.Create(filepath.Join(projectDir, baseName))
	if err != nil {
		return fmt.Errorf("error al crear el archivo de destino: %w", err)
	}
	defer destFile.Close()

	for _, part := range parts {
		fmt.Printf("Uniendo parte: %s\n", part)
		partFile, err := os.Open(part)
		if err != nil {
			return fmt.Errorf("error al abrir la parte del archivo: %w", err)
		}
		defer partFile.Close()

		if _, err := io.Copy(destFile, partFile); err != nil {
			return fmt.Errorf("error al copiar los datos de la parte: %w", err)
		}
	}

	// 4. Calcular el hash del archivo reconstruido
	newHash, err := hasher.CalculateSHA256(filepath.Join(projectDir, baseName))
	if err != nil {
		return fmt.Errorf("error al calcular el hash del archivo reconstruido: %w", err)
	}

	// 5. Comparar los hashes
	fmt.Println("\nVerificando integridad...")
	fmt.Printf("Hash original: %s\n", originalHash)
	fmt.Printf("Hash reconstruido: %s\n", newHash)

	if newHash == originalHash {
		fmt.Println("\n¡Integridad del archivo verificada! Los hashes coinciden.")
	} else {
		fmt.Println("\n¡ERROR DE INTEGRIDAD! Los hashes no coinciden. El archivo podría estar corrupto.")
	}

	return nil
}
