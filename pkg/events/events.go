package events

import (
	"fmt"

	"github.com/adrianriobo/qe-eventmanager/pkg/events/redhat/interop/ocp"
	"github.com/adrianriobo/qe-eventmanager/pkg/events/redhat/interop/rhel"
)

func GenerateRedHatInteropTestComplete(eventSchema, dahsboardURL,
	xunitURL, duration, resultStatus, contactName, contactEmail string,
	artifactFromEvent, systemFromEvent []byte) (interface{}, error) {
	switch eventSchema {
	case RedHatInteropRHELTestComplete:
		return rhel.CreateTestComplete(dahsboardURL, xunitURL,
			duration, resultStatus, contactName, contactEmail,
			artifactFromEvent, systemFromEvent)
	case RedHatInteropOCPTestComplete:
		return ocp.CreateTestComplete(dahsboardURL, xunitURL,
			duration, resultStatus, contactName, contactEmail,
			artifactFromEvent, systemFromEvent)
	default:
		return nil, fmt.Errorf("event schema is not available")
	}
}

func GenerateRedHatInteropTestError(eventSchema, dahsboardURL, contactName, contactEmail string,
	artifactFromEvent, systemFromEvent []byte) (interface{}, error) {
	switch eventSchema {
	case RedHatInteropRHELTestError:
		return rhel.CreateTestError(dahsboardURL, contactName, contactEmail,
			artifactFromEvent, systemFromEvent)
	case RedHatInteropOCPTestError:
		return ocp.CreateTestError(dahsboardURL, contactName, contactEmail,
			artifactFromEvent, systemFromEvent)
	default:
		return nil, fmt.Errorf("event schema is not available")
	}
}
