package generators

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/plantonhq/planton/pkg/kubernetes/helmcrds"
)

// TestHelmCRDsTFDrift holds every committed helm_crds.tf byte-identical to
// the generator. Presence-driven: the file exists only where the generator
// wrote it, so there is no allow-list to maintain. A red run means someone
// hand-edited a module's copy (regenerate it) or changed the generator
// without regenerating (run `planton tofu generate-helm-crds` into every
// module that carries the file).
func TestHelmCRDsTFDrift(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..")
	matches, err := filepath.Glob(filepath.Join(repoRoot, "catalog", "*", "*", "iac", "tf", HelmCRDsTFFileName))
	if err != nil {
		t.Fatal(err)
	}
	canonical := HelmCRDsTF()
	for _, path := range matches {
		rel, _ := filepath.Rel(repoRoot, path)
		committed, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(committed) != canonical {
			t.Errorf("%s drifted from the generator: regenerate it with `planton tofu generate-helm-crds --output-file %s`", rel, rel)
		}
	}
}

// TestHelmCRDsTFStampKeysMatchGo keeps the Terraform block's stamp literals
// equal to the Go constants the Pulumi twin writes. Both engines must read
// each other's CRDs (re-adoption, the never-downgrade check, the verifiers),
// so the keys are one contract with two spellings, held equal here.
func TestHelmCRDsTFStampKeysMatchGo(t *testing.T) {
	pairs := map[string]string{
		helmCRDsAnnotationSourceChart:   helmcrds.AnnotationSourceChart,
		helmCRDsAnnotationSourceVersion: helmcrds.AnnotationSourceVersion,
		helmCRDsLabelSource:             helmcrds.LabelSource,
		helmCRDsHelmKeepAnnotation:      helmcrds.HelmKeepAnnotation,
	}
	block := HelmCRDsTF()
	for literal, constant := range pairs {
		if literal != constant {
			t.Errorf("generator literal %q differs from helmcrds constant %q", literal, constant)
		}
		if !strings.Contains(block, `"`+literal+`"`) && !strings.Contains(block, `["`+literal+`"]`) {
			t.Errorf("the generated block does not carry the key %q", literal)
		}
	}
}
