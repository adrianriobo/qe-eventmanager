package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	umbCmdName string = "umb"
)

func init() {
	rootCmd.AddCommand(umbCmd)
	flagSet := pflag.NewFlagSet(umbCmdName, pflag.ExitOnError)
	umbCmd.Flags().AddFlagSet(flagSet)
}

var umbCmd = &cobra.Command{
	Use:   umbCmdName,
	Short: "umb operations",
	Long:  "umb operations",
}
