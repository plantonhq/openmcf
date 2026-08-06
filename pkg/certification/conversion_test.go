package certification

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sigs.k8s.io/yaml"

	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/conversion"
)

// The version-conversion certification cases. The corpus in pkg/conversion
// proves the engine's document mechanics; these cases prove the promises a
// USER experiences: an old manifest upgrades into a manifest the platform
// actually accepts, and nothing is ever lost without being said.

func tortureSpec(t *testing.T) *conversion.Spec {
	t.Helper()
	spec, err := conversion.LoadSpec(filepath.Join(TortureKindRoot(t), "conversions", "v1alpha1_to_v1alpha2.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return spec
}

func corpusInput(t *testing.T, name string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(TortureKindRoot(t), "conversions", "testdata", name, "input.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	return doc
}

// Case: an upgraded manifest is not merely shape-correct -- it passes the
// SAME offline validation an apply would run. Conversion output that fails
// validation is a spec defect, and it fails here, not at a server.
func TestCertify_UpgradedManifestPassesRealValidation(t *testing.T) {
	spec := tortureSpec(t)
	for _, fixture := range []string{"full-shape", "minimal"} {
		upgraded, _, err := conversion.Apply(spec, conversion.Upgrade, corpusInput(t, fixture))
		if err != nil {
			t.Fatalf("%s: upgrade failed: %v", fixture, err)
		}
		out, err := yaml.Marshal(upgraded)
		if err != nil {
			t.Fatal(err)
		}
		loaded, err := manifest.LoadManifestBytes(out, fixture+" (upgraded)")
		if err != nil {
			t.Errorf("%s: the upgraded manifest does not load against the served version: %v", fixture, err)
			continue
		}
		if err := manifest.ValidateLoaded(loaded); err != nil {
			t.Errorf("%s: the upgraded manifest does not validate against the served version: %v", fixture, err)
		}
	}
}

// Case: loss is never silent. The full-shape fixture carries the one field
// v1alpha2 removed; upgrading MUST report exactly that loss, with a reason a
// user can act on.
func TestCertify_LossIsDeclaredNeverSilent(t *testing.T) {
	spec := tortureSpec(t)
	_, losses, err := conversion.Apply(spec, conversion.Upgrade, corpusInput(t, "full-shape"))
	if err != nil {
		t.Fatal(err)
	}
	if len(losses) != 1 {
		t.Fatalf("expected exactly one declared loss (spec.stringNoDefault), got %d: %+v", len(losses), losses)
	}
	if losses[0].Path != "spec.stringNoDefault" {
		t.Errorf("wrong loss path: %s", losses[0].Path)
	}
	if !strings.Contains(losses[0].Reason, "removed") {
		t.Errorf("the loss reason must explain WHY the value is gone; got: %s", losses[0].Reason)
	}
}

// Case: the minimal fixture upgrades with ZERO losses -- declared loss fires
// only when a value actually existed to lose.
func TestCertify_NoLossWhenNothingToLose(t *testing.T) {
	spec := tortureSpec(t)
	_, losses, err := conversion.Apply(spec, conversion.Upgrade, corpusInput(t, "minimal"))
	if err != nil {
		t.Fatal(err)
	}
	if len(losses) != 0 {
		t.Errorf("the minimal manifest has nothing to lose, yet losses were declared: %+v", losses)
	}
}
