package results

import (
	"strings"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager/events"
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows"
	tektonClient "github.com/adrianriobo/qe-eventmanager/pkg/services/cicd/tekton"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/scm/github"
	"github.com/adrianriobo/qe-eventmanager/pkg/util"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/json"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	tektonUtil "github.com/adrianriobo/qe-eventmanager/pkg/util/tekton"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"golang.org/x/exp/slices"
)

const (
	eventFieldPipelineResultPrefix = "$(pipeline.results."
)

func ManageResults(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	event []byte, success flows.Success, errorFlow flows.Error) error {
	if tektonUtil.IsSuccessful(status) {
		logging.Debugf("Pipelinerun %s has finished successfully", pipelineRunName)
		return manageSuccess(status, pipelineRunName, event, success)
	}
	logging.Debugf("Pipelinerun %s has finished with errors", pipelineRunName)
	return manageError(status, pipelineRunName, event, errorFlow)
}

func manageSuccess(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	event []byte, success flows.Success) error {
	if !util.IsEmpty(success.UMB) {
		return manageSuccessUMB(status, pipelineRunName, event, success.UMB)
	}
	if !util.IsEmpty(success.Github) {
		return manageResultGithub(status, pipelineRunName, event, success.Github)
	}
	return nil
}

func manageSuccessUMB(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	event []byte, outputEvent flows.UMBEvent) error {
	if outputEvent.EventSchema == events.RedHatInteropOCPTestComplete ||
		outputEvent.EventSchema == events.RedHatInteropRHELTestComplete {
		artifactNode, systemNode := getDefaultNodes(event)
		response, err := events.GenerateRedHatInteropTestComplete(outputEvent.EventSchema,
			tektonClient.GetPipelinerunDashboardUrl(pipelineRunName),
			getEventFieldValue(outputEvent.EventFields, events.RedHatInteropFieldXunitURL, status, event),
			getEventFieldValue(outputEvent.EventFields, events.RedHatInteropFieldDuration, status, event),
			getEventFieldValue(outputEvent.EventFields, events.RedHatInteropFieldResultStatus, status, event),
			getEventFieldValue(outputEvent.EventFields, events.RedHatInteropFieldContactName, status, event),
			getEventFieldValue(outputEvent.EventFields, events.RedHatInteropFieldContactEmail, status, event),
			artifactNode, systemNode)
		if err != nil {
			return err
		}
		return umb.Send(outputEvent.Topic, response)
	}
	return nil
}

func manageError(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	event []byte, errorFLow flows.Error) error {
	if !util.IsEmpty(errorFLow.UMB) {
		return manageErrorUMB(status, pipelineRunName, event, errorFLow.UMB)
	}
	if !util.IsEmpty(errorFLow.Github) {
		return manageResultGithub(status, pipelineRunName, event, errorFLow.Github)
	}
	return nil
}

func manageErrorUMB(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	event []byte, outputEvent flows.UMBEvent) error {
	if outputEvent.EventSchema == events.RedHatInteropOCPTestError ||
		outputEvent.EventSchema == events.RedHatInteropRHELTestError {
		artifactNode, systemNode := getDefaultNodes(event)
		response, err := events.GenerateRedHatInteropTestError(outputEvent.EventSchema,
			tektonClient.GetPipelinerunDashboardUrl(pipelineRunName),
			getEventFieldValue(outputEvent.EventFields, events.RedHatInteropFieldContactName, status, event),
			getEventFieldValue(outputEvent.EventFields, events.RedHatInteropFieldContactEmail, status, event),
			artifactNode, systemNode)
		if err != nil {
			return err
		}
		return umb.Send(outputEvent.Topic, response)
	}
	return nil
}

func manageResultGithub(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	event []byte, githubInfo flows.Github) error {
	if !util.IsEmpty(githubInfo.Status) {
		return github.RepositoryStatus(
			getEventFieldValueByExpression(githubInfo.Status.Status, status, event),
			getEventFieldValueByExpression(githubInfo.Status.Owner, status, event),
			getEventFieldValueByExpression(githubInfo.Status.Repo, status, event),
			getEventFieldValueByExpression(githubInfo.Status.Commit, status, event),
			tektonClient.GetPipelinerunDashboardUrl(pipelineRunName), "", "")
	}
	return nil
}

func getEventFieldValue(source []flows.NameValuePair, item string,
	status *v1beta1.PipelineRunStatus, event []byte) string {
	itemExpression := getEventFieldExpression(source, item)
	return getEventFieldValueByExpression(itemExpression, status, event)
}

func getEventFieldExpression(source []flows.NameValuePair, item string) string {
	idx := slices.IndexFunc(source,
		func(e flows.NameValuePair) bool { return e.Name == item })
	if idx == -1 {
		return ""
	}
	return source[idx].Value
}

func getEventFieldValueByExpression(valueExpression string,
	status *v1beta1.PipelineRunStatus, event []byte) string {
	// Check if value should be picked from pipelinerun results
	if strings.HasPrefix(valueExpression, eventFieldPipelineResultPrefix) {
		resultNamePicker := strings.NewReplacer(eventFieldPipelineResultPrefix, "",
			")", "")
		resultName := resultNamePicker.Replace(valueExpression)
		return tektonUtil.GetResultValue(status.PipelineResults, resultName)
	}
	if strings.HasPrefix(valueExpression, json.JSONPathPreffix) {
		if value, err := json.GetStringValue(event, valueExpression); err != nil {
			logging.Errorf("error on picking value from event %v", err)
			return ""
		} else {
			return value
		}
	}
	// if strings.HasPrefix(valueExpression, json.
	// Otherwise value is a constant value
	return valueExpression
}

func getDefaultNodes(event []byte) (artifactNode, systemNode []byte) {
	//Default jsonpath could be passed as fields
	artifactNode, err := json.GetNodeAsByteArray([]byte(event),
		events.RedHatInteropNodeArtifactJSONPath)
	if err != nil {
		logging.Errorf("error getting artifactNode %v", err)
	}
	systemNode, err = json.GetNodeAsByteArray([]byte(event),
		events.RedHatInteropNodeSystemJSONPath)
	if err != nil {
		logging.Errorf("error getting systemNode %v", err)
	}
	return
}
