package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

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
		return fmt.Errorf("не вдалось рядок з файлу %s", err)
	}

	if hasRows {
		var skuIndex int
		var newSkuIndex int
		var priceIndex int
		var newPriceIndex int

		for i := 0; i < len(headers); i++ {
			switch headers[i] {
			case "Код/Артикль":
				skuIndex = i
			case "Новий Код/Артикль":
				newSkuIndex = i
			case "Ціна":
				priceIndex = i
			case "Нова ціна":
				newPriceIndex = i
			}
		}

		fmt.Println(skuIndex, newSkuIndex, priceIndex, newPriceIndex)

		var productsToUpdate []models.ProductUpdate

		for {
			row, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("не вдалось рядок з файлу %s", err)
			}

			changedPrice := row[newPriceIndex] != "" && row[newPriceIndex] != row[priceIndex]
			changedSku := row[newSkuIndex] != "" && row[newSkuIndex] != row[skuIndex]

			if changedPrice || changedSku {
				productsToUpdate = append(
					productsToUpdate,
					models.ProductUpdate{
						Id:    row[0],
						Sku:   row[newSkuIndex],
						Price: row[newPriceIndex],
					},
				)
			}
		}

		fmt.Println(productsToUpdate)
	}

	// products, err := apiClient.GetProductList(limit, lastId, groupId)
	// if err != nil {
	// 	return err
	// }

	return nil
}
