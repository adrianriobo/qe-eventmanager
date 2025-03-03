package github

import (
	"fmt"
	"os"

	"github.com/devtools-qe-incubator/eventmanager/pkg/configuration/providers"
	"github.com/devtools-qe-incubator/eventmanager/pkg/services/scm/github"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util/logging"
)

func CreateStatus(providersFilePath, ref, owner, repo,
	state, targetURL, statusContext, description string) error {
	if err := setupGithubClient(providersFilePath); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	return github.RepositoryStatus(state, owner, repo,
		ref, targetURL, statusContext, description)
}

func setupGithubClient(providersFilePath string) error {
	providersInfo, err := providers.LoadFile(providersFilePath)
	if err != nil {
		return err
	}
	if util.IsEmpty(providersInfo.Github) {
		return fmt.Errorf("github provider configuration is required")
	}
	if len(providersInfo.Github.AppKey) > 0 {
		appKey, err := providers.ParseGithubFiles(providersInfo.Github)
		if err != nil {
			return err
		}
		return github.CreateClientForApp(providersInfo.Github.AppID,
			providersInfo.Github.AppInstallationID, appKey)
	}
	return nil
}
