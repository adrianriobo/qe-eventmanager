package tekton

import (
	"github.com/adrianriobo/qe-eventmanager/pkg/util/http"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/xunit"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

const (
	resultStatusPassed string = "passed"
	resultStatusFailed string = "failed"
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
