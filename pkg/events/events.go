package events

import (
	"github.com/adrianriobo/qe-eventmanager/pkg/events/redhat/interop/ocp"
	"github.com/adrianriobo/qe-eventmanager/pkg/events/redhat/interop/rhel"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

const (
	eventSchemaRedHatInteropRHELTestcomplete = "interop-rhel-testcomplete"
	eventSchemaRedHatInteropOCPTestcomplete  = "interop-ocp-testcomplete"
)

func GenerateRedHatInteropTestComplete(eventSchema, dahsboardURL,
	xunitURL, duration, resultStatus string, artifact interface{}) interface{} {
	switch eventSchema {
	case eventSchemaRedHatInteropRHELTestcomplete:
		return rhel.CreateTestComplete(eventSchema, dahsboardURL, xunitURL,
			duration, resultStatus, artifact.(rhel.Artifact))
	case eventSchemaRedHatInteropOCPTestcomplete:
		return ocp.CreateTestComplete(eventSchema, dahsboardURL, xunitURL,
			duration, resultStatus, artifact.(ocp.Artifact))
	default:
		logging.Error("event schema is not available")
		return nil
	}
}
