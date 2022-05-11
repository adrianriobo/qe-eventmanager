package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager"
)

const (
	providersFilePath string = "providers-filepath"
	flowsFilePath     string = "flows-filepath"
)

func init() {
	rootCmd.AddCommand(startCmd)
	flagSet := pflag.NewFlagSet("start", pflag.ExitOnError)
	flagSet.StringP(providersFilePath, "p", "", "Credentials and defaults for integrated providers")
	flagSet.StringP(flowsFilePath, "r", "", "List of comma separated file paths of rules")
	startCmd.Flags().AddFlagSet(flagSet)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start event manager",
	Long:  "Start event manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		runStart()
		return nil
	},
}

func runStart() {
	manager.Initialize(
		viper.GetString(providersFilePath),
		strings.Split(viper.GetString(flowsFilePath), ","))
}
