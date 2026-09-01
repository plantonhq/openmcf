package pulumistack

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"reflect"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/pkg/outputs"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// captureOutputs reads the just-updated stack's outputs and fills sink with
// the raw map, the flattened map, per-output sensitivity, and the kind's
// typed StackOutputs proto (honoring module-shipped transform overrides via
// the module directory).
//
// Pulumi's plain `stack output --json` carries no per-output sensitivity
// flag — it MASKS secret values instead. So capture reads the stack twice:
// the masked pass tells us WHICH outputs are secret (any output whose value
// differs between the passes), the --show-secrets pass supplies the real
// values downstream reference resolution needs. Both passes are cheap state
// reads against the same backend the update just used.
func captureOutputs(
	stackFqdn string,
	moduleRepoPath string,
	kindName string,
	extraEnv []string,
	sink *outputs.CaptureResult,
) error {
	masked, err := readStackOutputs(stackFqdn, moduleRepoPath, extraEnv, false)
	if err != nil {
		return errors.Wrap(err, "failed to read masked stack outputs")
	}

	shown, err := readStackOutputs(stackFqdn, moduleRepoPath, extraEnv, true)
	if err != nil {
		return errors.Wrap(err, "failed to read stack outputs with --show-secrets")
	}

	sink.Raw = shown
	sink.Sensitive = detectSecretOutputs(masked, shown)
	sink.Flat = outputs.Flatten(shown)

	kind := crkreflect.KindFromString(kindName)
	if kind == cloudresourcekind.CloudResourceKind_unspecified {
		return errors.Errorf("cannot resolve cloud resource kind from %q for output transformation", kindName)
	}

	typed, flat, err := outputs.TransformRaw(kind, shown, &outputs.TransformOptions{ModuleDir: moduleRepoPath})
	if err != nil {
		return errors.Wrapf(err, "output transformation failed for kind %s", kindName)
	}
	sink.Typed = typed
	if flat != nil {
		sink.Flat = flat
	}

	return nil
}

// readStackOutputs runs `pulumi stack output --json` (optionally with
// --show-secrets) and decodes the plain name->value map.
func readStackOutputs(stackFqdn, moduleRepoPath string, extraEnv []string, showSecrets bool) (map[string]interface{}, error) {
	args := []string{"stack", "output", "--stack", stackFqdn, "--json", "--non-interactive"}
	if showSecrets {
		args = append(args, "--show-secrets")
	}

	cmd := exec.Command("pulumi", args...)
	cmd.Dir = moduleRepoPath
	cmd.Env = append(os.Environ(), extraEnv...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, errors.Wrapf(err, "failed to execute pulumi %v", args)
	}

	values := map[string]interface{}{}
	if len(bytes.TrimSpace(stdout.Bytes())) > 0 {
		if err := json.Unmarshal(stdout.Bytes(), &values); err != nil {
			return nil, errors.Wrap(err, "stack output document is not a JSON object")
		}
	}
	return values, nil
}

// detectSecretOutputs marks every output whose value differs between the
// masked and --show-secrets reads. Pulumi replaces secret values with a
// placeholder in the masked form, so a difference IS the sensitivity signal;
// an output missing from the masked form entirely is treated as secret too.
func detectSecretOutputs(masked, shown map[string]interface{}) map[string]bool {
	sensitive := make(map[string]bool, len(shown))
	for name, shownValue := range shown {
		maskedValue, present := masked[name]
		sensitive[name] = !present || !reflect.DeepEqual(maskedValue, shownValue)
	}
	return sensitive
}
