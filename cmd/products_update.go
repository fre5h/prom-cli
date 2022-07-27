package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/fre5h/prom-cli/internal/models"
)

var (
	productUpdateCmd = &cobra.Command{
		Use:     "products:update",
		Short:   "Внесення змін в товари взятих з CSV файла",
		Long:    "Внесення змін в товари взятих з CSV файла, що відправляться в кабінеті Prom.ua",
		Args:    validateArgs,
		RunE:    updateListOfProducts,
		PreRunE: preRunE,
	}
)

func init() {
	productUpdateCmd.Flags().StringVarP(&apiKey, "apiKey", "k", "", "секретний API ключ для доступу до кабінету Prom.ua")
	productUpdateCmd.Flags().StringVarP(&fileName, "fileName", "f", "data.csv", "назва файлу, з якого будуть зчитані товари в форматі CSV")

	rootCmd.AddCommand(productUpdateCmd)
}

func updateListOfProducts(_ *cobra.Command, _ []string) error {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("не вдалось прочитати файл %s", fileName)
	}

	var hasRows = true

	r := csv.NewReader(file)
	headers, err := r.Read()
	if err == io.EOF {
		hasRows = false
	}
	if err != nil {
		return fmt.Errorf("не вдалось прочитати рядок з файлу %s", err)
	}

	if hasRows {
		var priceIndex int
		var newPriceIndex int

		for i := 0; i < len(headers); i++ {
			switch headers[i] {
			case "Ціна":
				priceIndex = i
			case "Нова ціна":
				newPriceIndex = i
			}
		}

		var productsToUpdate []models.ProductUpdate

		for {
			row, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("не вдалось прочитати рядок з файлу %s", err)
			}

			if row[newPriceIndex] != "" && row[newPriceIndex] != row[priceIndex] {
				id, _ := strconv.Atoi(row[0])
				price, _ := strconv.ParseFloat(row[newPriceIndex], 64)
				productsToUpdate = append(productsToUpdate, models.ProductUpdate{ID: id, Price: price})
			}
		}

		return apiClient.UpdateProduct(productsToUpdate)
	}

	return nil
}
