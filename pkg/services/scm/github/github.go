package github

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	ghinstallation "github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

var (
	commitStatusDescription string = "Tested on downstream infrastructure"
	commitStatusContext     string = "qe-eventmanager"
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
	// Shared transport to reuse TCP connections.
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

func CommitStatus(state, owner, repo, ref, dashboardURL string) error {
	status, response, err := _client.Repositories.CreateStatus(context.Background(),
		owner, repo, ref, &github.RepoStatus{
			State:       &state,
			Description: &commitStatusDescription,
			TargetURL:   &dashboardURL,
			Context:     &commitStatusContext})
	if err != nil {
		return err
	}
	logging.Debugf("status is %v", status)
	logging.Debugf("response is %v", response)
	return nil
}
