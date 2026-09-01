package tofumodule

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/pkg/outputs"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// tofuOutputEnvelope is one entry of `tofu output -json` / `terraform output -json`:
// every output rides in an envelope carrying its value, type, and sensitivity.
type tofuOutputEnvelope struct {
	Value     json.RawMessage `json:"value"`
	Sensitive bool            `json:"sensitive"`
}

// captureOutputs reads the just-applied stack's outputs back with
// `<binary> output -json`, unwraps the envelopes, and fills sink with the raw
// map, the flattened map, per-output sensitivity, and the kind's typed
// StackOutputs proto (honoring module-shipped transform overrides via the
// module directory).
//
// It must run while the module workspace is still alive: `output -json`
// re-reads state through the configured backend, so it needs the same
// provider env vars and the provider-override file the apply used.
func captureOutputs(
	ctx context.Context,
	binaryName string,
	modulePath string,
	kindName string,
	providerConfigEnvVars []string,
	sink *outputs.CaptureResult,
) error {
	cmd := newReapableCommand(ctx, binaryName, "output", "-json")
	cmd.Dir = modulePath
	cmd.Env = append(os.Environ(), providerConfigEnvVars...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to read stack outputs with `%s output -json`", binaryName)
	}

	raw, sensitive, err := unwrapTofuOutputs(stdout.Bytes())
	if err != nil {
		return errors.Wrapf(err, "failed to decode `%s output -json`", binaryName)
	}

	sink.Raw = raw
	sink.Sensitive = sensitive
	sink.Flat = outputs.Flatten(raw)

	kind := crkreflect.KindFromString(kindName)
	if kind == cloudresourcekind.CloudResourceKind_unspecified {
		return errors.Errorf("cannot resolve cloud resource kind from %q for output transformation", kindName)
	}

	typed, flat, err := outputs.TransformRaw(kind, raw, &outputs.TransformOptions{ModuleDir: modulePath})
	if err != nil {
		return errors.Wrapf(err, "output transformation failed for kind %s", kindName)
	}
	sink.Typed = typed
	if flat != nil {
		sink.Flat = flat
	}

	return nil
}

// unwrapTofuOutputs decodes the `output -json` document into the plain
// name->value map plus the per-output sensitivity the envelopes declare.
// An empty document (a module with no outputs prints `{}`) is not an error.
func unwrapTofuOutputs(outputJson []byte) (map[string]interface{}, map[string]bool, error) {
	envelopes := map[string]tofuOutputEnvelope{}
	if len(bytes.TrimSpace(outputJson)) > 0 {
		if err := json.Unmarshal(outputJson, &envelopes); err != nil {
			return nil, nil, errors.Wrap(err, "output document is not the expected {name: {value, sensitive}} shape")
		}
	}

	values := make(map[string]interface{}, len(envelopes))
	sensitive := make(map[string]bool, len(envelopes))
	for name, envelope := range envelopes {
		var value interface{}
		if len(envelope.Value) > 0 {
			if err := json.Unmarshal(envelope.Value, &value); err != nil {
				return nil, nil, errors.Wrapf(err, "output %q carries an undecodable value", name)
			}
		}
		values[name] = value
		sensitive[name] = envelope.Sensitive
	}
	return values, sensitive, nil
}
