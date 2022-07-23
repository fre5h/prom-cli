package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/fre5h/prom-cli/internal/service"
)

var (
	productListCmd = &cobra.Command{
		Use:   "products:list",
		Short: "Список продуктів",
		Long:  "Список продуктів компанії на Prom.ua",
		Run: func(cmd *cobra.Command, args []string) {
			getListOfProducts(apiKey)
		},
	}
)

func init() {
	productListCmd.Flags().IntVarP(&limit, "limit", "l", 20, "максимальна кількість товарів у відповіді")
	productListCmd.Flags().IntVarP(&lastId, "lastId", "i", 0, "обмежити вибірку товарів з ідентифікаторами більшими за вказаний")
	productListCmd.Flags().IntVarP(&groupId, "groupId", "g", 0, "ідентифікатор групи. по замовчуванню - ідентифікатор кореневої групи")
	productListCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "показати більше полів")

	rootCmd.AddCommand(productListCmd)
}

func getListOfProducts(apiKey string) {
	products, err := service.GetProductList(apiKey, limit, lastId, groupId)
	if err != nil {
		fmt.Println(err)
	}

	cyan := color.New(color.FgCyan)
	yellow := color.New(color.FgYellow)

	var numberOfGroups = len(products)
	var i = 0

	for _, product := range products {
		i++
		fmt.Println(cyan.Sprint("ID:"), yellow.Sprint(product.Id))
		fmt.Println(cyan.Sprint("Name:"), product.Name)

		if i < numberOfGroups {
			fmt.Println()
		}
	}
}
