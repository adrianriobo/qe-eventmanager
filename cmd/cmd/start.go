package cmd

import (
	"strings"

	"github.com/go-stomp/stomp/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/adrianriobo/qe-eventmanager/pkg/logging"
	"github.com/adrianriobo/qe-eventmanager/pkg/messaging"
)

const (
	flagBrokers         string = "brokers"
	flagCertificateFile string = "certificate-file"
	privateKeyFile      string = "private-key-file"
	caCerts             string = "ca-certs"
)

func init() {
	rootCmd.AddCommand(startCmd)

	flagSet := pflag.NewFlagSet("start", pflag.ExitOnError)
	flagSet.StringP(flagBrokers, "b", "", "list of brokers acting on failover")
	flagSet.StringP(flagCertificateFile, "", "", "certificate file for client auth")
	flagSet.StringP(privateKeyFile, "", "", "key file for client auth")
	flagSet.StringP(caCerts, "", "", "root ca for messageing server auth")

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
	connection := messaging.NewUMBConnection(
		viper.GetString(flagCertificateFile),
		viper.GetString(privateKeyFile),
		viper.GetString(privateKeyFile),
		strings.Split(viper.GetString(flagBrokers), ","))
	if err := connection.Connect(); err != nil {
		return err
	}
	sub, _ := connection.FailoverSubscribe("Consumer.psi-crcqe-openstack.test5.VirtualTopic.qe.ci.product-scenario.crcqe.test4",
		stomp.AckAuto)
	msg, err := sub.Read()
	if err == nil {
		logging.Debugf(string(msg.Body[:]))
	}
	connection.Disconnect()
	return nil
}

func validateStartFlags() error {
	return nil
}
