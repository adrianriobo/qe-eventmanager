package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	githubCmdName string = "github"
)

func init() {
	rootCmd.AddCommand(githubCmd)
	flagSet := pflag.NewFlagSet(githubCmdName, pflag.ExitOnError)
	githubCmd.Flags().AddFlagSet(flagSet)
}

var githubCmd = &cobra.Command{
	Use:   githubCmdName,
	Short: "github operations",
	Long:  "github operations",
}
