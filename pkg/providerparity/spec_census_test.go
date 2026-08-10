//go:build !codegen
// +build !codegen

package providerparity

import (
	"reflect"
	"testing"

	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// TestCollectSpecPaths_HermeticFixture proves the descriptor walk against the
// permanent testcloudresourcegeneric fixture at its REGISTERED version
// (v1alpha2 -- the registry serves a kind's declared version, so the walk
// must be asserted against the compiled contract, not a .proto read by eye):
// scalars and optional scalars are leaves, a nested message contributes its
// own leaves (never itself), a repeated message recurses into its element
// (steps -> steps.command), a StringValueOrRef is ONE leaf (the author
// configures one value slot), and a map field is one leaf.
func TestCollectSpecPaths_HermeticFixture(t *testing.T) {
	msg, err := crkreflect.NewInstance(cloudresourcekind.CloudResourceKind_TestCloudResourceGeneric)
	if err != nil {
		t.Fatalf("new instance: %v", err)
	}
	specField := msg.ProtoReflect().Descriptor().Fields().ByName("spec")
	if specField == nil {
		t.Fatal("fixture has no spec field")
	}

	got := CollectSpecPaths(specField.Message(), "spec")
	want := []string{
		"spec.annotated_ref",
		"spec.bool_field",
		"spec.display_name",
		"spec.double_field",
		"spec.float_field",
		"spec.int32_field",
		"spec.int64_field",
		"spec.labels",
		"spec.nested.nested_int",
		"spec.nested.nested_string",
		"spec.optional_ref",
		"spec.replicas",
		"spec.required_ref",
		"spec.sensitive_ref",
		"spec.sensitive_string",
		"spec.steps.command",
		"spec.uint32_field",
		"spec.uint64_field",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("paths mismatch:\n got %v\nwant %v", got, want)
	}
}

// TestCollectSpecPaths_StructIsOneLeaf proves the well-known-value-type
// rule against a live spec that carries a google.protobuf.Struct field
// (AwsEventBridgeRule's event_pattern): the Struct is ONE leaf the author
// writes as YAML -- the walk must never descend into its descriptor
// internals (fields.bool_value and siblings are protobuf plumbing, not
// configurable surface).
func TestCollectSpecPaths_StructIsOneLeaf(t *testing.T) {
	msg, err := crkreflect.NewInstance(cloudresourcekind.CloudResourceKind_AwsEventBridgeRule)
	if err != nil {
		t.Fatalf("new instance: %v", err)
	}
	specField := msg.ProtoReflect().Descriptor().Fields().ByName("spec")
	if specField == nil {
		t.Fatal("AwsEventBridgeRule has no spec field")
	}

	got := CollectSpecPaths(specField.Message(), "spec")
	var sawLeaf bool
	for _, p := range got {
		if p == "spec.event_pattern" {
			sawLeaf = true
			continue
		}
		if len(p) > len("spec.event_pattern") && p[:len("spec.event_pattern.")] == "spec.event_pattern." {
			t.Errorf("walk descended into the Struct's internals: %s", p)
		}
	}
	if !sawLeaf {
		t.Error("spec.event_pattern is missing from the census -- the Struct field should be one leaf")
	}
}

// TestSpecCensusGcp is the live-catalog smoke test: every implemented GCP
// kind is censused and none has an empty spec surface (a kind an author
// cannot configure at all would be a registry error, not a real component).
func TestSpecCensusGcp(t *testing.T) {
	census := SpecCensus(cloudresourcekind.CloudResourceProvider_gcp)
	if len(census) == 0 {
		t.Fatal("GCP spec census is empty -- the registry walk is broken")
	}
	for _, k := range census {
		if len(k.SpecFieldPaths) == 0 {
			t.Errorf("%s: zero spec leaf fields", k.Kind)
		}
	}
}
