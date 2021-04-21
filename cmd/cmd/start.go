package cmd

import (
	"github.com/adrianriobo/qe-eventmanager/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(startCmd)

	flagSet := pflag.NewFlagSet("start", pflag.ExitOnError)
	// flagSet.StringArrayP("brokers", "b", nil, "list of brokers acting on failover")
	// flagSet.StringP("certificate-file", "c", "", "certificate file for client auth")
	// flagSet.StringP("certificate-file", "c", "", "certificate file for client auth")
	// flagSet.StringP("certificate-file", "c", "", "certificate file for client auth")

	startCmd.Flags().AddFlagSet(flagSet)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start event manager",
	Long:  "Start event manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runStart(); err != nil {
			return err
		}
		return nil
	},
}

func runStart() error {
	if err := validateStartFlags(); err != nil {
		return err
	}
	logging.Infof("Testing the start")
	return nil
}

func validateStartFlags() error {
	return nil
}
