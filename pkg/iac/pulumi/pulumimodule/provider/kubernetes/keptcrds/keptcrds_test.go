package keptcrds

import (
	"testing"

	"github.com/plantonhq/planton/pkg/kubernetes/helmcrds"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// The live read and the resource registration are proven by the e2e lanes;
// what can be proven offline is the mapping from what the cluster returns to
// what the never-downgrade check compares.
func TestToExistingReadsTheVersionStamp(t *testing.T) {
	stamped := unstructured.Unstructured{}
	stamped.SetName("widgets.example")
	stamped.SetAnnotations(map[string]string{helmcrds.AnnotationSourceVersion: "0.120.0"})
	unstamped := unstructured.Unstructured{}
	unstamped.SetName("legacy.example")

	existing := toExisting([]unstructured.Unstructured{stamped, unstamped})
	if len(existing) != 2 {
		t.Fatalf("expected 2, got %d", len(existing))
	}
	if existing[0].Name != "widgets.example" || existing[0].Version != "0.120.0" {
		t.Fatalf("stamped CRD mapped wrong: %+v", existing[0])
	}
	if existing[1].Version != "" {
		t.Fatalf("an unstamped CRD must report no version, got %q", existing[1].Version)
	}
	// An unstamped CRD never counts as a downgrade; a stamped higher one does.
	if err := helmcrds.CheckNoDowngrade(existing, "0.120.3"); err != nil {
		t.Fatal(err)
	}
	if err := helmcrds.CheckNoDowngrade(existing, "0.119.0"); err == nil {
		t.Fatal("expected the downgrade refusal")
	}
}

func TestApplyWithoutInstallRegistersNothing(t *testing.T) {
	resources, err := Apply(nil, Args{Install: false})
	if err != nil || resources != nil {
		t.Fatalf("install=false must be a no-op, got %v %v", resources, err)
	}
}
