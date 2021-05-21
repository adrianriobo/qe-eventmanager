package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager"
)

const (
	brokers         string = "brokers"
	certificateFile string = "certificate-file"
	privateKeyFile  string = "private-key-file"
	caCertsFile     string = "ca-certs"
	kubeconfigPath  string = "kubeconfigPath"
)

func init() {
	rootCmd.AddCommand(startCmd)
	flagSet := pflag.NewFlagSet("start", pflag.ExitOnError)
	flagSet.StringP(brokers, "b", "", "list of brokers acting on failover")
	flagSet.StringP(certificateFile, "", "", "certificate file for client auth")
	flagSet.StringP(privateKeyFile, "", "", "key file for client auth")
	flagSet.StringP(caCertsFile, "", "", "root ca for messageing server auth")
	flagSet.StringP(kubeconfigPath, "k", "", "kubeconfig file path")
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
		viper.GetString(certificateFile),
		viper.GetString(privateKeyFile),
		viper.GetString(caCertsFile),
		viper.GetString(kubeconfigPath),
		strings.Split(viper.GetString(brokers), ","))
}
