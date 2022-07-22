package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/fre5h/promua-helper/internal/service"
)

var (
	groupListCmd = &cobra.Command{
		Use:   "groups:list",
		Short: "Список груп",
		Long:  "Список груп компанії на Prom.ua",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("для цієї команди не передбачено додаткових аргументів")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			getListOfCompanyGroups(apiKey)
		},
	}
)

func init() {
	groupListCmd.Flags().IntVarP(&limit, "limit", "l", 20, "максимальна кількість груп у відповіді")
	groupListCmd.Flags().IntVarP(&lastId, "lastId", "i", 0, "обмежити вибірку груп з ідентифікаторами більшими за вказаний")
	groupListCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "показати більше полів")

	rootCmd.AddCommand(groupListCmd)
}

func getListOfCompanyGroups(apiKey string) {
	fmt.Println("groups command executions")
	fmt.Println("api key flag", apiKey)
	fmt.Println("api key flag", os.Getenv("PROM_UA_API_KEY"))

	groups, err := service.GetGroupList(apiKey, limit, lastId)
	if err != nil {
		fmt.Println(err)
	}

	cyan := color.New(color.FgCyan)
	yellow := color.New(color.FgYellow)

	var numberOfGroups = len(groups)
	var i = 0

	for _, group := range groups {
		i++
		fmt.Println(cyan.Sprint("ID:"), yellow.Sprint(group.Id))
		fmt.Println(cyan.Sprint("Назва:"), group.Name)
		fmt.Println(cyan.Sprint("Батьківська група:"), group.ParentGroupId)

		if verbose {
			fmt.Println(cyan.Sprint("Картинка:"), group.Image)
			fmt.Println(cyan.Sprint("Опис:"), group.Description)
		}

		if i < numberOfGroups {
			fmt.Println()
		}
	}
}
