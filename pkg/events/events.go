package events

import (
	"fmt"

	"github.com/adrianriobo/qe-eventmanager/pkg/events/redhat/interop/ocp"
	"github.com/adrianriobo/qe-eventmanager/pkg/events/redhat/interop/rhel"
)

const (
	RedHatInteropRHELTestComplete = "interop-rhel-test-complete"
	RedHatInteropOCPTestComplete  = "interop-ocp-test-complete"
	RedHatInteropRHELTestError    = "interop-rhel-test-error"
	RedHatInteropOCPTestError     = "interop-ocp-test-error"
	RedHatInteropXunitURL         = "xunitURL"
	RedHatInteropDuration         = "duration"
	RedHatInteropResultStatus     = "resultStatus"
)

func GenerateRedHatInteropTestComplete(eventSchema, dahsboardURL,
	xunitURL, duration, resultStatus string, artifactFromEvent, systemFromEvent []byte) (interface{}, error) {
	switch eventSchema {
	case RedHatInteropRHELTestComplete:
		return rhel.CreateTestComplete(dahsboardURL, xunitURL,
			duration, resultStatus, artifactFromEvent, systemFromEvent)
	case RedHatInteropOCPTestComplete:
		return ocp.CreateTestComplete(dahsboardURL, xunitURL,
			duration, resultStatus, artifactFromEvent, systemFromEvent)
	default:
		return nil, fmt.Errorf("event schema is not available")
	}
}

func GenerateRedHatInteropTestError(eventSchema, dahsboardURL string,
	artifactFromEvent, systemFromEvent []byte) (interface{}, error) {
	switch eventSchema {
	case RedHatInteropRHELTestError:
		return rhel.CreateTestError(dahsboardURL, artifactFromEvent, systemFromEvent)
	case RedHatInteropOCPTestError:
		return ocp.CreateTestError(dahsboardURL, artifactFromEvent, systemFromEvent)
	default:
		return nil, fmt.Errorf("event schema is not available")
	}
}
