package controlprofile

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	controlprofilev1 "github.com/plantonhq/planton/compliance/componentcontrolprofile/v1"
	"github.com/plantonhq/planton/pkg/compliance/controlcatalog"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/pkg/specpath"
)

// TestControlProfileConformance holds every authored control profile to its
// contract, offline:
//
//  1. The profile parses strictly against its proto schema and names its
//     component (metadata.name equals the component directory).
//  2. Every referenced control id exists in the central catalog, and every
//     catalog control appears EXACTLY ONCE -- completeness is the product
//     claim ("every control was examined"), so an omitted control fails the
//     same way an invented one does. not_applicable with a reason counts as
//     examined.
//  3. Posture claims carry their proof: enforced_by_default and
//     configurable require typed evidence; configurable requires the spec
//     field that makes the choice; not_applicable requires the reason.
//  4. Every field path resolves against the served version's compiled
//     descriptors -- a schema rename that orphans a claim fails CI loudly.
func TestControlProfileConformance(t *testing.T) {
	root := repoRoot(t)
	if _, err := os.Stat(filepath.Join(root, "catalog")); err != nil {
		t.Skip("catalog source tree not present (bazel sandbox); runs under go test and the lint.catalog-data lane")
	}

	catalog, err := controlcatalog.Load(root)
	if err != nil {
		t.Fatalf("control catalog: %v", err)
	}
	catalogIDs := controlcatalog.ControlIDs(catalog)

	discovered, err := Discover(root)
	if err != nil {
		t.Fatalf("discovering control profiles: %v", err)
	}
	if len(discovered) == 0 {
		t.Skip("no control profiles authored yet")
	}

	for provider, components := range discovered {
		for _, component := range components {
			component := component
			t.Run(provider+"/"+component, func(t *testing.T) {
				profile, err := Load(root, provider, component)
				if err != nil {
					t.Fatalf("control profile: %v", err)
				}
				if profile.GetKind() != "ComponentControlProfile" {
					t.Fatalf("kind is %q, want ComponentControlProfile", profile.GetKind())
				}
				if profile.GetMetadata().GetName() != component {
					t.Errorf("metadata.name is %q, want %q", profile.GetMetadata().GetName(), component)
				}

				specDescriptor := kindSpecDescriptor(t, component)
				seen := map[string]bool{}
				for _, posture := range profile.GetSpec().GetControls() {
					id := posture.GetControlId()
					if !catalogIDs[id] {
						t.Errorf("control %q does not exist in the control catalog -- profiles never invent controls", id)
						continue
					}
					if seen[id] {
						t.Errorf("control %q appears more than once", id)
					}
					seen[id] = true

					switch posture.GetStatus() {
					case controlprofilev1.Status_status_unspecified:
						t.Errorf("control %s has no status", id)
					case controlprofilev1.Status_enforced_by_default, controlprofilev1.Status_configurable:
						evidence := posture.GetEvidence()
						if evidence.GetType() == controlprofilev1.Type_type_unspecified ||
							strings.TrimSpace(evidence.GetReference()) == "" {
							t.Errorf("control %s (%s) carries no evidence -- an evidence-free posture claim is marketing, not data",
								id, posture.GetStatus())
						}
						if posture.GetStatus() == controlprofilev1.Status_configurable &&
							strings.TrimSpace(posture.GetFieldPath()) == "" {
							t.Errorf("control %s is configurable but names no field_path -- which spec choice makes it true?", id)
						}
					case controlprofilev1.Status_not_applicable:
						if strings.TrimSpace(posture.GetNotes()) == "" {
							t.Errorf("control %s is not_applicable without a reason in notes -- unexplained exclusion is indistinguishable from an unexamined control", id)
						}
					}

					if path := posture.GetFieldPath(); path != "" {
						if err := specpath.Validate(specDescriptor, path); err != nil {
							t.Errorf("control %s field_path %q: %v", id, path, err)
						}
					}
				}

				for id := range catalogIDs {
					if !seen[id] {
						t.Errorf("catalog control %q was not examined -- add it with an honest status (not_applicable with a reason counts)", id)
					}
				}
			})
		}
	}
}

// kindSpecDescriptor resolves a component directory name to its kind's spec
// message descriptor via the kind registry.
func kindSpecDescriptor(t *testing.T, component string) protoreflect.MessageDescriptor {
	t.Helper()
	kind := crkreflect.KindFromString(component)
	apiMessage, err := crkreflect.NewInstance(kind)
	if err != nil {
		t.Fatalf("NewInstance(%s): %v", component, err)
	}
	return specDescriptor(t, apiMessage)
}

// specDescriptor returns the descriptor of the kind's spec message (the api
// envelope's `spec` field).
func specDescriptor(t *testing.T, apiMessage proto.Message) protoreflect.MessageDescriptor {
	t.Helper()
	specField := apiMessage.ProtoReflect().Descriptor().Fields().ByName("spec")
	if specField == nil || specField.Kind() != protoreflect.MessageKind {
		t.Fatalf("%s has no spec message field", apiMessage.ProtoReflect().Descriptor().FullName())
	}
	return specField.Message()
}

// repoRoot walks up from this test file to the directory containing go.mod.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found walking up from test file")
		}
		dir = parent
	}
}
