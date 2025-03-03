package ack

import (
	"strings"

	"github.com/devtools-qe-incubator/eventmanager/pkg/configuration/flows"
	"github.com/devtools-qe-incubator/eventmanager/pkg/services/scm/github"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util/json"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util/logging"
)

type ACK func(event []byte) error

func CreateACK(ackInfo flows.ACK) ACK {
	if !util.IsEmpty(ackInfo.Github) {
		return createGithubStatusACK(ackInfo.Github)
	}
	return nil
}

func createGithubStatusACK(statusInfo flows.Github) ACK {
	return func(event []byte) error {
		logging.Debug("Sending ack for event")
		return github.RepositoryStatus(
			getValueByExpression(statusInfo.Status.Status, event),
			getValueByExpression(statusInfo.Status.Owner, event),
			getValueByExpression(statusInfo.Status.Repo, event),
			getValueByExpression(statusInfo.Status.Ref, event),
			"", "", "")
	}
}

func getValueByExpression(valueExpression string, event []byte) string {
	if strings.HasPrefix(valueExpression, json.JSONPathPreffix) {
		if value, err := json.GetStringValue(event, valueExpression); err != nil {
			logging.Errorf("error on picking value from event %v", err)
			return ""
		} else {
			return value
		}
	}
	return valueExpression
}
