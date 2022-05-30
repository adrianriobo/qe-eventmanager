package github

import (
	"context"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

var (
	tektonAvatar            string = "https://avatars.githubusercontent.com/u/48335577?v=4"
	commitStatusDescription string = "Tested on downstream infrastructure"
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

func CommitStatus(state, owner, repo, ref, dashboardURL string) error {
	status, response, err := _client.Repositories.CreateStatus(context.Background(),
		owner, repo, ref, &github.RepoStatus{
			State:       &state,
			Description: &commitStatusDescription,
			AvatarURL:   &tektonAvatar,
			TargetURL:   &dashboardURL})
	if err != nil {
		return err
	}
	logging.Debugf("status is %v", status)
	logging.Debugf("response is %v", response)
	return nil
}
