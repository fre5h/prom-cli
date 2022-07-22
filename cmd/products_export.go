package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/fre5h/promua-helper/internal/service"
)

var (
	productExportCmd = &cobra.Command{
		Use:   "products:export",
		Short: "Експорт списку продуктів в CSV файл",
		Long:  "Експорт списку продуктів в CSV файл",
		Run: func(cmd *cobra.Command, args []string) {
			exportListOfProducts(apiKey)
		},
	}
)

func init() {
	productExportCmd.Flags().IntVarP(&limit, "limit", "l", 20, "максимальна кількість товарів у відповіді")
	productExportCmd.Flags().IntVarP(&lastId, "lastId", "i", 0, "обмежити вибірку товарів з ідентифікаторами більшими за вказаний")
	productExportCmd.Flags().IntVarP(&groupId, "groupId", "g", 0, "ідентифікатор групи. по замовчуванню - ідентифікатор кореневої групи")
	productExportCmd.Flags().StringVarP(&outputFileName, "outputFileName", "o", "data.csv", "назва файлу, в який буде збережено отриманий результат")

	// @todo Check that file can be readable

	rootCmd.AddCommand(productExportCmd)
}

func exportListOfProducts(apiKey string) {
	products, err := service.GetProductList(apiKey, limit, lastId, groupId)
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.Create(outputFileName)
	defer file.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write([]string{"ID", "Name", "Price"})

	for _, product := range products {
		w.Write([]string{
			strconv.Itoa(product.Id),
			product.Name,
			fmt.Sprintf("%.2f", product.Price),
		})

		if err := w.Error(); err != nil {
			log.Fatalln("error writing csv:", err)
		}
	}
}
