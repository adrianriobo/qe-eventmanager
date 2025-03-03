package tekton

import (
	"github.com/devtools-qe-incubator/eventmanager/pkg/util"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util/http"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util/logging"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util/xunit"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"knative.dev/pkg/apis"
)

const (
	resultStatusPassed string = "passed"
	resultStatusFailed string = "failed"
)

var (
	succeededConditions []string = []string{string(v1beta1.PipelineRunReasonSuccessful),
		string(v1beta1.PipelineRunReasonCompleted)}
	failedConditions []string = []string{string(v1beta1.PipelineRunReasonFailed),
		string(v1beta1.PipelineRunReasonCancelled),
		string(v1beta1.PipelineRunReasonTimedOut)}
)

// TODO make general vailable
func GetResultValue(results []v1beta1.PipelineRunResult, resultParamID string) string {
	for _, result := range results {
		if result.Name == resultParamID {
			return result.Value
		}
	}
	return ""
}

// TODO this should be moved to result parameter from the pipeline
func GetResultState(url string) string {
	file, err := http.GetFile(url)
	if err != nil {
		logging.Error(err)
		return ""
	}
	count, err := xunit.CountFailures(file)
	if err != nil {
		logging.Error(err)
		return ""
	}
	if count == 0 {
		return resultStatusPassed
	}
	return resultStatusFailed
}

func IsSuccessful(status *v1beta1.PipelineRunStatus) bool {
	return util.SliceContains(succeededConditions,
		status.GetCondition(apis.ConditionSucceeded).Reason)
}

func IsFailed(status *v1beta1.PipelineRunStatus) bool {
	return util.SliceContains(failedConditions,
		status.GetCondition(apis.ConditionSucceeded).Reason)
}
