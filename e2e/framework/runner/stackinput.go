// Package runner implements the E2E test lifecycle engine.
package runner

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/iac/stackinput"
	"github.com/plantonhq/planton/pkg/iac/stackinput/stackinputproviderconfig"
)

// BuildStackInput loads a component's hack manifest and constructs the stack-input
// YAML that Pulumi modules expect. providerConfig is nil for the harness's
// default posture (ambient credentials, empty provider block) and non-nil only
// when the component ships the opt-in provider-config fixture (see
// LoadProviderConfigFixture).
//
// Returns the path to the written stack-input file.
func BuildStackInput(manifestPath string, outputDir string,
	providerConfig *stackinputproviderconfig.ProviderConfig) (string, error) {
	manifestObject, err := manifest.LoadManifest(manifestPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to load manifest from %s", manifestPath)
	}

	stackInputYaml, err := stackinput.BuildStackInputYaml(manifestObject, providerConfig)
	if err != nil {
		return "", errors.Wrap(err, "failed to build stack-input YAML")
	}

	outputPath := filepath.Join(outputDir, "stack-input.yaml")
	if err := os.WriteFile(outputPath, []byte(stackInputYaml), 0600); err != nil {
		return "", errors.Wrapf(err, "failed to write stack-input to %s", outputPath)
	}

	return outputPath, nil
}
