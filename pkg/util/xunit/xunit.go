package xunit

import (
	"github.com/devtools-qe-incubator/eventmanager/pkg/util/logging"
	"github.com/joshdk/go-junit"
)

func CountFailures(xunit []byte) (int, error) {
	suites, err := junit.Ingest(xunit)
	if err != nil {
		logging.Error(err)
		return 0, err
	}
	failures := 0
	for _, suite := range suites {
		failures += suite.Totals.Failed
	}
	return failures, nil
}
