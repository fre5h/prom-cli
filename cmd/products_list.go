package cmd

import (
	"fmt"
	"math"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	productListCmd = &cobra.Command{
		Use:     "products:list",
		Short:   "Список товарів",
		Long:    "Список товарів з кабінету Prom.ua",
		Args:    validateArgs,
		RunE:    getListOfProducts,
		PreRunE: preRunE,
	}
)

func init() {
	productListCmd.Flags().StringVarP(&apiKey, "apiKey", "k", "", "секретний API ключ для доступу до кабінету Prom.ua")
	productListCmd.Flags().IntVarP(&limit, "limit", "l", math.MaxInt32, "максимальна кількість товарів у відповіді")
	productListCmd.Flags().IntVarP(&lastId, "lastId", "i", 0, "обмежити вибірку товарів з ідентифікаторами більшими за вказаний")
	productListCmd.Flags().IntVarP(&groupId, "groupId", "g", 0, "ідентифікатор групи. по замовчуванню - ідентифікатор кореневої групи")
	productListCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "показати більше полів")

	rootCmd.AddCommand(productListCmd)
}

func getListOfProducts(_ *cobra.Command, _ []string) error {
	cyan := color.New(color.FgCyan)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)

	fmt.Printf(green.Sprint("Відправка запиту до Prom.ua на отримання списку товарів...\n"))
	fmt.Printf(green.Sprint("Зачекайте, будь ласка...\n\n"))

	products, err := apiClient.GetProductList(limit, lastId, groupId)
	if err != nil {
		return err
	}

	numberOfProducts := len(products)
	if 0 == numberOfProducts {
		return fmt.Errorf("не знайдено жодного товару")
	}

	var i = 0

	for _, product := range products {
		i++
		fmt.Println(cyan.Sprint("ID:"), yellow.Sprint(product.ID))
		fmt.Println(cyan.Sprint("Назва:"), product.Name)
		fmt.Println(cyan.Sprint("Група:"), product.Group.Name, blue.Sprintf("(ID: %d)", product.Group.ID))
		fmt.Println(cyan.Sprint("Категорія:"), product.Category.Caption, blue.Sprintf("(ID: %d)", product.Category.ID))
		fmt.Println(cyan.Sprint("Код/Артикул:"), product.Sku)
		fmt.Println(cyan.Sprint("Ціна:"), fmt.Sprintf("%.2f", product.Price), product.Currency)

		if verbose {
			fmt.Println(cyan.Sprint("Опис:"), product.Description)
		}

		if i < numberOfProducts {
			fmt.Println()
		}
	}

	fmt.Printf(green.Sprintf("\nКількість знайдених товарів: %d\n", numberOfProducts))

	return nil
}
