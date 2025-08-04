// test/core_test.go (Versión corregida del test)
package test

import (
	"go-file-splitter/internal/core"
	"os"
	"path/filepath"
	"testing"
)

func TestSplitAndJoinFile(t *testing.T) {
	// Crear un archivo de prueba temporal
	testFilePath := filepath.Join(os.TempDir(), "testfile.txt")
	testData := []byte("This is a test file with some data.")
	if err := os.WriteFile(testFilePath, testData, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(testFilePath)

	// Dividir el archivo
	chunkSize := int64(10)
	if err := core.SplitFile(testFilePath, chunkSize); err != nil {
		t.Fatalf("SplitFile failed: %v", err)
	}

	// Limpiar los archivos de los trozos al final del test
	// La función glob buscará los trozos en el mismo directorio que el archivo original
	parts, _ := filepath.Glob(testFilePath + ".part*")
	for _, part := range parts {
		defer os.Remove(part)
	}

	// Unir los trozos
	// La función JoinFiles ahora encontrará los archivos correctamente
	if err := core.JoinFiles(testFilePath); err != nil {
		t.Fatalf("JoinFiles failed: %v", err)
	}

	// Verificar el archivo unido
	joinedData, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("failed to read joined file: %v", err)
	}

	if string(joinedData) != string(testData) {
		t.Errorf("joined file content does not match original")
	}
}
