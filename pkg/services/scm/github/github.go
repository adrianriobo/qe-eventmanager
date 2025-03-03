package github

import (
	"context"
	"log"
	"net/http"
	"strconv"

	ghinstallation "github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util/logging"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

var (
	defaultStatusDescription string = "Tested on downstream infrastructure"
	defaultStatusContext     string = "eventmanager"
)

var _client *github.Client

func CreateClientForUser(pat string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)
	_client = github.NewClient(tc)
	rateLimits, _, err := _client.RateLimits(ctx)
	if err != nil {
		return err
	}
	logging.Debugf("Github client initialized with rate limit %v", rateLimits.GetCore().Limit)
	return nil
}

func CreateClientForApp(appID, appInstallationID string, appKey []byte) error {
	tr := http.DefaultTransport
	appIDAsInt, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		return err
	}
	appInstallationAsInt, err := strconv.ParseInt(appInstallationID, 10, 64)
	if err != nil {
		return err
	}
	itr, err := ghinstallation.New(tr, appIDAsInt, appInstallationAsInt, appKey)
	if err != nil {
		log.Fatal(err)
	}
	_client = github.NewClient(&http.Client{Transport: itr})
	rateLimits, _, err := _client.RateLimits(context.Background())
	if err != nil {
		return err
	}
	logging.Debugf("Github client initialized as github app with id %s with rate limit %v",
		appID, rateLimits.GetCore().Limit)
	return nil
}

func RepositoryStatus(state, owner, repo, ref, targetURL, statusContext, description string) error {
	if len(statusContext) == 0 {
		statusContext = defaultStatusContext
	}
	if len(description) == 0 {
		description = defaultStatusDescription
	}
	_, _, err := _client.Repositories.CreateStatus(context.Background(),
		owner, repo, ref, &github.RepoStatus{
			State:       &state,
			Description: &description,
			TargetURL:   &targetURL,
			Context:     &statusContext})
	if err != nil {
		return err
	}
	logging.Debugf("Sending repository status with state %s to repo %s and ref %s", state, repo, ref)
	return nil
}
