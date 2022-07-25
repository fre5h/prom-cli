package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/fre5h/prom-cli/internal/api"
)

var (
	// Flags
	apiKey   string
	limit    int
	lastId   int
	groupId  int
	verbose  bool
	fileName string

	apiClient *api.Client

	apiClient *api.Client

	usageTemplate = `Використання:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [команда]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Приклади:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Доступні команди:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Параметри:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Глобальні параметри:

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Спробуйте "{{.CommandPath}} [команда] --help" для отримання більшої інформації про команду.{{end}}
`

	rootCmd = &cobra.Command{
		Use:   "prom-cli",
		Short: "Консольна програма для роботи з кабінетом Prom.ua",
		Long: `Консольна програма для роботи з кабінетом Prom.ua.
Ця програма дозволить вам використовувати деякий функціонал Prom.ua, який доступний через існуючі API методи.`,
		Hidden:       true,
		SilenceUsage: true,
		Version:      "1.0.0",
		Args:         validateArgs,
	}
)

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "довідка про будь-яку команду")
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// for _, c := range append(rootCmd.Commands(), rootCmd) {
	// 	c.Flags().BoolP("help", "h", false, "Help for "+c.Name())
	// }
	//
	// for _, c := range rootCmd.Commands() {
	// 	c.Flags().BoolP("help", "h", false, "Help for "+c.Name())
	// }
}

func Execute() error {
	return rootCmd.Execute()
}

func preRunE(cmd *cobra.Command, args []string) error {
	if apiKey == "" {
		godotenv.Load(".env")
		apiKey = os.Getenv("PROM_UA_API_KEY")
	}

	if apiKey == "" {
		return fmt.Errorf(`Відсутній API ключ, ви можете вказати його за допомогою параметра -k або --apiKey,
або через змінну оточення PROM_UA_API_KEY. Для отримання допоміжної інформації запустіть:
  prom-cli help
`)
	}

	apiClient = api.NewClient(apiKey)

	return nil
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("для цієї команди не передбачено додаткових аргументів")
	}

	return nil
}
