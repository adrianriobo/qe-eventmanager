package github

import (
	"context"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

var _client *github.Client

func CreateClient(token string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
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

func test(owner string, repo string, ref int) error {
	// for PR need action synchronize or opened
	// repoStatus := &github.RepoStatus{}
	// repoStatus, _, err := _client.Repositories.CreateStatus(context.Background(), owner, repo, ref, repoStatus)
	// _client.Checks.CreateCheckRun()
	return nil
}
