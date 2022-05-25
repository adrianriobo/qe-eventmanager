package events

const (
	RedHatInteropRHELTestComplete = "interop-rhel-test-complete"
	RedHatInteropOCPTestComplete  = "interop-ocp-test-complete"
	RedHatInteropRHELTestError    = "interop-rhel-test-error"
	RedHatInteropOCPTestError     = "interop-ocp-test-error"

	RedHatInteropFieldXunitURL     = "xunit-url"
	RedHatInteropFieldDuration     = "duration"
	RedHatInteropFieldResultStatus = "result-status"
	RedHatInteropFieldContactName  = "contact-name"
	RedHatInteropFieldContactEmail = "contact-email"

	RedHatInteropNodeArtifactJSONPath = "$.artifact"
	RedHatInteropNodeSystemJSONPath   = "$.system"
)
