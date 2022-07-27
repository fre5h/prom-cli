package cmd

import (
	"fmt"
	"math"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	groupListCmd = &cobra.Command{
		Use:     "groups:list",
		Short:   "Список груп",
		Args:    validateArgs,
		RunE:    getListOfCompanyGroups,
		PreRunE: preRunE,
	}
)

func init() {
	groupListCmd.Flags().StringVarP(&apiKey, "apiKey", "k", "", "секретний API ключ для доступу до кабінету Prom.ua")
	groupListCmd.Flags().IntVarP(&limit, "limit", "l", math.MaxInt32, "максимальна кількість груп у відповіді")
	groupListCmd.Flags().IntVarP(&lastId, "lastId", "i", 0, "обмежити вибірку груп з ідентифікаторами більшими за вказаний")
	groupListCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "показати більше полів")

	rootCmd.AddCommand(groupListCmd)
}

func getListOfCompanyGroups(_ *cobra.Command, _ []string) error {
	cyan := color.New(color.FgCyan)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	groups, err := apiClient.GetGroupList(limit, lastId)
	if err != nil {
		return err
	}

	fmt.Printf(green.Sprint("Відправка запиту до Prom.ua на отримання списку груп...\n"))
	fmt.Printf(green.Sprint("Зачекайте, будь ласка...\n\n"))

	var numberOfGroups = len(groups)
	if 0 == numberOfGroups {
		return fmt.Errorf("не знайдено жодної групи")
	}

	var i = 0

	for _, group := range groups {
		i++
		fmt.Println(cyan.Sprint("ID:"), yellow.Sprint(group.ID))
		fmt.Println(cyan.Sprint("Назва:"), group.Name)
		fmt.Println(cyan.Sprint("Батьківська група:"), group.ParentGroupId)

		if verbose {
			fmt.Println(cyan.Sprint("Зображення:"), group.Image)
			fmt.Println(cyan.Sprint("Опис:"), group.Description)
		}

		if i < numberOfGroups {
			fmt.Println()
		}
	}

	fmt.Printf(green.Sprintf("\nКількість знайдених груп: %d\n", numberOfGroups))

	return nil
}
