package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/fre5h/prom-cli/internal/models"
)

var (
	productUpdateCmd = &cobra.Command{
		Use:     "products:update",
		Short:   "Зміни ціни товарів (з CSV файла)",
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

		green := color.New(color.FgGreen)
		red := color.New(color.FgRed)

		productsUpdateResult, err := apiClient.UpdateProduct(productsToUpdate)
		if err != nil {
			return fmt.Errorf("помилка на оновленні товарів: %s", err)
		}

		numberOfProcessedProducts := len(productsUpdateResult.ProcessedIds)
		if numberOfProcessedProducts > 0 {
			fmt.Println(green.Sprint("Успішно оновлені товари:"))

			for _, id := range productsUpdateResult.ProcessedIds {
				fmt.Println(green.Sprintf(" %s", id))
			}

			fmt.Printf(green.Sprintf("\nКількість успішно оновлених товарів: %d\n", numberOfProcessedProducts))
		}

		numberOfErrors := len(productsUpdateResult.Errors)
		if numberOfErrors > 0 {
			fmt.Println(red.Sprint("Помилки:"))

			for _, errorMessage := range productsUpdateResult.Errors {
				fmt.Println(red.Sprintf(" %s", errorMessage))
			}

			fmt.Printf(red.Sprintf("\nКількість помилок: %d\n", numberOfErrors))
		}
	}

	return nil
}
