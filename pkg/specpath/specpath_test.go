package specpath

import (
	"testing"

	componentv1 "github.com/plantonhq/planton/iac/componentpermissions/v1"
)

// The tests walk a real committed descriptor (ComponentPermissionsSpec) so
// they exercise the same protoreflect surfaces production paths do, with no
// synthetic fixture to maintain.
func TestValidate(t *testing.T) {
	desc := (&componentv1.ComponentPermissionsSpec{}).ProtoReflect().Descriptor()

	valid := []string{
		"aws",                             // message leaf
		"aws.statements",                  // repeated message leaf
		"aws.statements.actions",          // traversal THROUGH a repeated field to a repeated scalar leaf
		"kubernetes.rules.cluster_scoped", // scalar leaf
	}
	for _, path := range valid {
		if err := Validate(desc, path); err != nil {
			t.Errorf("Validate(%q) = %v, want nil", path, err)
		}
	}

	invalid := []string{
		"",                             // empty
		"nope",                         // unknown root
		"aws.nope",                     // unknown nested
		"aws.statements.sid.deeper",    // continues past a scalar
		"aws.statements.actions.après", // continues past a repeated scalar
	}
	for _, path := range invalid {
		if err := Validate(desc, path); err == nil {
			t.Errorf("Validate(%q) = nil, want error", path)
		}
	}
}
