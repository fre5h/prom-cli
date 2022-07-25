package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	productListCmd = &cobra.Command{
		Use:     "products:list",
		Short:   "Список продуктів",
		Long:    "Список продуктів з кабінету Prom.ua",
		Args:    validateArgs,
		RunE:    getListOfProducts,
		PreRunE: preRunE,
	}
)

func init() {
	productListCmd.Flags().StringVarP(&apiKey, "apiKey", "k", "", "секретний API ключ для доступу до кабінету Prom.ua")
	productListCmd.Flags().IntVarP(&limit, "limit", "l", 20, "максимальна кількість товарів у відповіді")
	productListCmd.Flags().IntVarP(&lastId, "lastId", "i", 0, "обмежити вибірку товарів з ідентифікаторами більшими за вказаний")
	productListCmd.Flags().IntVarP(&groupId, "groupId", "g", 0, "ідентифікатор групи. по замовчуванню - ідентифікатор кореневої групи")
	productListCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "показати більше полів")

	rootCmd.AddCommand(productListCmd)
}

func getListOfProducts(cmd *cobra.Command, args []string) error {
	products, err := apiClient.GetProductList(limit, lastId, groupId)
	if err != nil {
		return err
	}

	var numberOfProducts = len(products)
	if 0 == numberOfProducts {
		return fmt.Errorf("не знайдено жодного товару")
	}

	cyan := color.New(color.FgCyan)
	yellow := color.New(color.FgYellow)

	var i = 0

	for _, product := range products {
		i++
		fmt.Println(cyan.Sprint("ID:"), yellow.Sprint(product.Id))
		fmt.Println(cyan.Sprint("Назва:"), product.Name)
		fmt.Println(cyan.Sprint("Група:"), product.Group.Name)
		fmt.Println(cyan.Sprint("Категорія:"), product.Category.Caption)
		fmt.Println(cyan.Sprint("Код/Артикул:"), product.Sku)
		fmt.Println(cyan.Sprint("Ціна:"), product.Price, product.Currency)

		if verbose {
			fmt.Println(cyan.Sprint("Опис:"), product.Description)
		}

		if i < numberOfProducts {
			fmt.Println()
		}
	}

	return nil
}
