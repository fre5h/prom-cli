package cmd

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/fre5h/prom-cli/internal/models"
)

var (
	productExportCmd = &cobra.Command{
		Use:     "products:export",
		Short:   "Експорт списку товарів в CSV файл",
		Args:    validateArgs,
		RunE:    exportListOfProducts,
		PreRunE: preRunE,
	}
)

func init() {
	productExportCmd.Flags().StringVarP(&apiKey, "apiKey", "k", "", "секретний API ключ для доступу до кабінету Prom.ua")
	productExportCmd.Flags().IntVarP(&limit, "limit", "l", math.MaxInt32, "максимальна кількість товарів у відповіді")
	productExportCmd.Flags().IntVarP(&lastId, "lastId", "i", 0, "обмежити вибірку товарів з ідентифікаторами більшими за вказаний")
	productExportCmd.Flags().IntVarP(&groupId, "groupId", "g", 0, "ідентифікатор групи. по замовчуванню - ідентифікатор кореневої групи")
	productExportCmd.Flags().StringVarP(&fileName, "fileName", "f", "data.csv", "назва файлу, в який буде збережено отримані товари у форматі CSV")

	rootCmd.AddCommand(productExportCmd)
}

func exportListOfProducts(_ *cobra.Command, _ []string) error {
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("не вдалось створити файл %s", fileName)
	}

	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	fmt.Printf(green.Sprint("Відправка запиту до Prom.ua на отримання списку товарів\n"))
	fmt.Printf(green.Sprint("Зачекайте, будь ласка...\n"))

	products, err := apiClient.GetProductList(limit, lastId, groupId)
	if err != nil {
		return err
	}

	var numberOfProducts = len(products)
	if 0 == numberOfProducts {
		return fmt.Errorf("не знайдено жодного товару")
	}

	sort.Sort(models.ProductsArray(products))

	w := csv.NewWriter(file)
	defer w.Flush()

	fmt.Printf(green.Sprintf("Кількість знайдених товарів: "))
	fmt.Printf(yellow.Sprintf("%d\n\n", numberOfProducts))
	fmt.Printf(green.Sprint("Експорт товарів в файл...\n"))

	w.Write([]string{
		"ID",
		"Назва",
		"Група",
		"Категорія",
		"Статус",
		"Код/Артикль",
		"Ціна",
		"Нова ціна",
		"Остання дата редагування",
	})

	for _, product := range products {
		w.Write([]string{
			strconv.Itoa(product.ID),
			product.Name,
			product.Group.Name,
			product.Category.Caption,
			product.GetTranslatedStatus(),
			product.Sku,
			fmt.Sprintf("%.2f", product.Price),
			"",
			product.DateModified.Format("02.01.2006 15:04:05"),
		})

		if err := w.Error(); err != nil {
			return fmt.Errorf("помилка запису CSV: %s", err)
		}
	}

	fmt.Printf(green.Sprintf("Товари експортовано у файл "))
	fmt.Printf(yellow.Sprintf("%s\n", fileName))

	fmt.Printf(green.Sprintf("Кількість експортованих товарів у файл: "))
	fmt.Printf(yellow.Sprintf("%d\n", numberOfProducts))

	return nil
}
