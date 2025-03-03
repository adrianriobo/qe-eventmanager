package cmd

import (
	"github.com/devtools-qe-incubator/eventmanager/pkg/actioner/github"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	githubStatusCmdName string = "status"

	ref           string = "ref"
	owner         string = "owner"
	repo          string = "repo"
	state         string = "state"
	targetURL     string = "target-url"
	statusContext string = "context"
	description   string = "description"
)

func init() {
	githubCmd.AddCommand(githubStatusCmd)
	flagSet := pflag.NewFlagSet(githubStatusCmdName, pflag.ExitOnError)
	flagSet.StringP(providersFilePath, "p", "", "File containing credentials for Github provider")
	flagSet.StringP(ref, "", "", "Ref for the status; can be a SHA, a branch name, or a tag name")
	flagSet.StringP(owner, "o", "", "Owner for the repository")
	flagSet.StringP(repo, "", "", "Repository")
	flagSet.StringP(state, "s", "",
		"State is the current state of the repository. Possible values are: pending, success, error, or failure")
	flagSet.StringP(targetURL, "u", "", "TargetURL is the URL of the page representing this status")
	flagSet.StringP(statusContext, "c", "", "A string label to differentiate this status from the statuses of other systems")
	flagSet.StringP(description, "d", "", "Description is a short high level summary of the status.")
	githubStatusCmd.Flags().AddFlagSet(flagSet)
}

var githubStatusCmd = &cobra.Command{
	Use:   githubStatusCmdName,
	Short: "create / update repository status",
	Long:  "create / update repository status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		return github.CreateStatus(
			viper.GetString(providersFilePath),
			viper.GetString(ref),
			viper.GetString(owner),
			viper.GetString(repo),
			viper.GetString(state),
			viper.GetString(targetURL),
			viper.GetString(statusContext),
			viper.GetString(description))
	},
}
