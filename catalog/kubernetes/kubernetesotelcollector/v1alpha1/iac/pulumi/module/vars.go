package module

var vars = struct {
	// ApiVersion / Kind of the custom resource the OpenTelemetry Operator
	// serves (v1beta1 is the storage version; its `config` is a
	// structured object, not the v1alpha1 string).
	ApiVersion string
	Kind       string

	// NameBudget bounds metadata.name: the operator derives child names
	// by suffixing — "-collector-networkpolicy" is the longest at 25
	// characters (rendered only under the operand.networkpolicy feature
	// gate, but an operator-side gate flip must never break existing
	// collector names; "-collector-monitoring" at 21 is the longest
	// default-path suffix) — verified in the operator's naming source at
	// the pin. Kubernetes caps names at 63.
	NameBudget int
}{
	ApiVersion: "opentelemetry.io/v1beta1",
	Kind:       "OpenTelemetryCollector",
	NameBudget: 38,
}
