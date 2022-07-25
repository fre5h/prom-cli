package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	productExportCmd = &cobra.Command{
		Use:     "products:export",
		Short:   "Експорт списку продуктів в CSV файл",
		Long:    "Експорт списку продуктів з кабінету Prom.ua в CSV файл",
		Args:    validateArgs,
		RunE:    exportListOfProducts,
		PreRunE: preRunE,
	}
)

func init() {
	productExportCmd.Flags().StringVarP(&apiKey, "apiKey", "k", "", "секретний API ключ для доступу до кабінету Prom.ua")
	productExportCmd.Flags().IntVarP(&limit, "limit", "l", 20, "максимальна кількість товарів у відповіді")
	productExportCmd.Flags().IntVarP(&lastId, "lastId", "i", 0, "обмежити вибірку товарів з ідентифікаторами більшими за вказаний")
	productExportCmd.Flags().IntVarP(&groupId, "groupId", "g", 0, "ідентифікатор групи. по замовчуванню - ідентифікатор кореневої групи")
	productExportCmd.Flags().StringVarP(&outputFileName, "outputFileName", "o", "data.csv", "назва файлу, в який буде збережено отриманий результат")

	rootCmd.AddCommand(productExportCmd)
}

func exportListOfProducts(cmd *cobra.Command, args []string) error {
	file, err := os.Create(outputFileName)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("не вдалось створити файл %s", outputFileName)
	}

	products, err := apiClient.GetProductList(limit, lastId, groupId)
	if err != nil {
		return err
	}

	var numberOfProducts = len(products)
	if 0 == numberOfProducts {
		return fmt.Errorf("не знайдено жодного товару")
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write([]string{"ID", "Назва", "Група", "Категорія", "Статус", "Код/Артикль", "Новий Код/Артикль", "Ціна", "Нова Ціна"})

	for _, product := range products {
		w.Write([]string{
			strconv.Itoa(product.Id),
			product.Name,
			product.Group.Name,
			product.Category.Caption,
			product.GetTranslatedStatus(),
			product.Sku,
			"",
			fmt.Sprintf("%.2f", product.Price),
			"",
		})

		if err := w.Error(); err != nil {
			return fmt.Errorf("помилка запису CSV: %s", err)
		}
	}

	return nil
}
