package cli

import (
	"fmt"
	"go-file-splitter/internal/core"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// Define una variable para el tamaño del trozo como un string
var chunkSizeStr string

var rootCmd = &cobra.Command{
	Use:   "file-splitter",
	Short: "Herramienta para dividir y unir archivos",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Por favor, utilice 'file-splitter split' o 'file-splitter join'.")
	},
}

var splitCmd = &cobra.Command{
	Use:   "split [filename]",
	Short: "Dividir un archivo grande en partes más pequeñas",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Convierte el tamaño del trozo a bytes
		chunkSize, err := parseSize(chunkSizeStr)
		if err != nil {
			fmt.Println("Error: Formato de tamaño inválido.")
			return
		}

		if chunkSize <= 0 {
			// Valor por defecto: 100 MB
			chunkSize = 1024 * 1024 * 100
		}

		if err := core.SplitFile(args[0], chunkSize); err != nil {
			fmt.Println("Error:", err)
		}
	},
}

var joinCmd = &cobra.Command{
	Use:   "join [filename]",
	Short: "Unir archivos más pequeños al archivo original",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// La función JoinFiles ahora necesita el nombre del archivo original, no la carpeta.
		if err := core.JoinFiles(args[0]); err != nil {
			fmt.Println("Error:", err)
		}
	},
}

// parseSize convierte una cadena como "100M" o "1G" a un valor en bytes
func parseSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToUpper(sizeStr)
	if strings.HasSuffix(sizeStr, "G") {
		val, err := strconv.ParseInt(strings.TrimSuffix(sizeStr, "G"), 10, 64)
		if err != nil {
			return 0, err
		}
		return val * 1024 * 1024 * 1024, nil
	} else if strings.HasSuffix(sizeStr, "M") {
		val, err := strconv.ParseInt(strings.TrimSuffix(sizeStr, "M"), 10, 64)
		if err != nil {
			return 0, err
		}
		return val * 1024 * 1024, nil
	} else if strings.HasSuffix(sizeStr, "K") {
		val, err := strconv.ParseInt(strings.TrimSuffix(sizeStr, "K"), 10, 64)
		if err != nil {
			return 0, err
		}
		return val * 1024, nil
	}

	// Si no tiene sufijo, lo trata como un número en bytes.
	return strconv.ParseInt(sizeStr, 10, 64)
}

func init() {
	// Agrega los subcomandos al comando raíz
	rootCmd.AddCommand(splitCmd, joinCmd)

	// Define el flag para el tamaño como un string
	splitCmd.Flags().StringVarP(&chunkSizeStr, "size", "s", "100M", "Size of each chunk (e.g., 100M, 1G, 2048K)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
