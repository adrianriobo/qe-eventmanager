package tekton

import (
	"strings"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager/events"
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows"
	tektonClient "github.com/adrianriobo/qe-eventmanager/pkg/services/cicd/tekton"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/json"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	tektonUtil "github.com/adrianriobo/qe-eventmanager/pkg/util/tekton"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"golang.org/x/exp/slices"
)

const (
	eventFieldPipelineResultPrefix = "$(pipeline.results."
)

func manageResults(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	event []byte, success flows.Success, errorFlow flows.Error) error {
	if tektonUtil.IsSuccessful(status) {
		return manageSuccess(status, pipelineRunName, event, success)
	}
	return manageError(status, pipelineRunName, event, errorFlow)
}

func manageSuccess(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	event []byte, success flows.Success) error {
	if success.UMB.EventSchema == events.RedHatInteropOCPTestComplete ||
		success.UMB.EventSchema == events.RedHatInteropRHELTestComplete {
		artifactNode, systemNode := getDefaultNodes(event)
		response, err := events.GenerateRedHatInteropTestComplete(success.UMB.EventSchema,
			tektonClient.GetPipelinerunDashboardUrl(pipelineRunName),
			getEventFieldValue(success.UMB.EventFields, events.RedHatInteropFieldXunitURL, status),
			getEventFieldValue(success.UMB.EventFields, events.RedHatInteropFieldDuration, status),
			getEventFieldValue(success.UMB.EventFields, events.RedHatInteropFieldResultStatus, status),
			getEventFieldValue(success.UMB.EventFields, events.RedHatInteropFieldContactName, status),
			getEventFieldValue(success.UMB.EventFields, events.RedHatInteropFieldContactEmail, status),
			artifactNode, systemNode)
		if err != nil {
			return err
		}
		return umb.Send(success.UMB.Topic, response)
	}
	return nil
}

func manageError(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	event []byte, errorFLow flows.Error) error {
	if errorFLow.UMB.EventSchema == events.RedHatInteropOCPTestError ||
		errorFLow.UMB.EventSchema == events.RedHatInteropRHELTestError {
		artifactNode, systemNode := getDefaultNodes(event)
		response, err := events.GenerateRedHatInteropTestError(errorFLow.UMB.EventSchema,
			tektonClient.GetPipelinerunDashboardUrl(pipelineRunName),
			getEventFieldValue(errorFLow.UMB.EventFields, events.RedHatInteropFieldContactName, status),
			getEventFieldValue(errorFLow.UMB.EventFields, events.RedHatInteropFieldContactEmail, status),
			artifactNode, systemNode)
		if err != nil {
			return err
		}
		return umb.Send(errorFLow.UMB.Topic, response)
	}
	return nil
}

func getEventFieldValue(source []flows.NameValuePair, item string,
	status *v1beta1.PipelineRunStatus) string {
	itemExpression := getEventFieldExpression(source, item)
	return getEventFieldValueByExpression(itemExpression, status)
}

func getEventFieldExpression(source []flows.NameValuePair, item string) string {
	idx := slices.IndexFunc(source,
		func(e flows.NameValuePair) bool { return e.Name == item })
	if idx == -1 {
		return ""
	}
	return source[idx].Value
}

func getEventFieldValueByExpression(valueExpression string, status *v1beta1.PipelineRunStatus) string {
	// Check if value should be picked from pipelinerun results
	if strings.HasPrefix(valueExpression, eventFieldPipelineResultPrefix) {
		resultNamePicker := strings.NewReplacer(eventFieldPipelineResultPrefix, "",
			")", "")
		resultName := resultNamePicker.Replace(valueExpression)
		return tektonUtil.GetResultValue(status.PipelineResults, resultName)
	}
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
