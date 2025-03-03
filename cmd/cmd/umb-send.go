package cmd

import (
	"github.com/devtools-qe-incubator/eventmanager/pkg/actioner/umb"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	umbSendCmdName string = "send"

	destination     string = "destination"
	messageFilePath string = "message-filepath"
)

func init() {
	umbCmd.AddCommand(umbSendCmd)
	flagSet := pflag.NewFlagSet(umbSendCmdName, pflag.ExitOnError)
	flagSet.StringP(providersFilePath, "p", "", "File containing credentials for UMB provider")
	flagSet.StringP(destination, "d", "", "UMB topic to send the message")
	flagSet.StringP(messageFilePath, "m", "", "File containing the json message to send")
	umbSendCmd.Flags().AddFlagSet(flagSet)
}

var umbSendCmd = &cobra.Command{
	Use:   umbSendCmdName,
	Short: "send messages throug umb",
	Long:  "send messages throug umb",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		return umb.Send(
			viper.GetString(providersFilePath),
			viper.GetString(destination),
			viper.GetString(messageFilePath))
	},
}
