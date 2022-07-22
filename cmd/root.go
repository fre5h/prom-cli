package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Flags
	apiKey         string
	limit          int
	lastId         int
	groupId        int
	verbose        bool
	outputFileName string

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

Global Flags:

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Спробуйте "{{.CommandPath}} [команда] --help" для отримання більшої інформації про команду.{{end}}
`

	rootCmd = &cobra.Command{
		Use:   "promua-cli",
		Short: "Консольна програма для роботи з кабінетом Prom.ua",
		Long: `Консольна програма для роботи з кабінетом Prom.ua.
Ця програма дозволить вам використовувати деякий функціонал Prom.ua
який доступний через існуючі API методи.`,
		Hidden:  true,
		Version: "1.0.0",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&apiKey, "apiKey", "k", "", "секретний API ключ для доступу до кабінету Prom.ua")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetUsageTemplate(usageTemplate)
}

func Execute() error {
	return rootCmd.Execute()
}
