package outputs

import (
	"google.golang.org/protobuf/proto"
)

// CaptureResult carries a stack's outputs captured after a successful apply,
// in every shape a consumer needs: the engine-decoded raw map, the dotted-key
// flattening, the kind's typed StackOutputs proto, and per-output sensitivity.
//
// Sensitivity is tracked per TOP-LEVEL output name (the granularity both
// engines report at: tofu's output envelope carries a `sensitive` flag per
// output; pulumi masks whole outputs in its plain `stack output --json` form).
// Renderers MUST consult it before printing: captured values include secrets
// by design — resolving downstream references needs the real values — so the
// leak boundary is rendering, not capture.
type CaptureResult struct {
	// Raw is the engine-decoded output map: output name -> JSON-decoded value.
	Raw map[string]interface{}

	// Flat is Raw flattened to dotted string keys (see Flatten).
	Flat map[string]string

	// Typed is the kind's StackOutputs proto populated from Raw, honoring
	// module-shipped transform overrides. Nil when the kind declares no
	// outputs message or the transform was skipped.
	Typed proto.Message

	// Sensitive marks top-level output names whose values must never render
	// in terminal or log output.
	Sensitive map[string]bool
}

// IsSensitive reports whether a flat (dotted) key belongs to a sensitive
// top-level output. A dotted key inherits its root output's sensitivity.
func (r *CaptureResult) IsSensitive(flatKey string) bool {
	if r == nil || len(r.Sensitive) == 0 {
		return false
	}
	root := flatKey
	for i := 0; i < len(flatKey); i++ {
		if flatKey[i] == '.' {
			root = flatKey[:i]
			break
		}
	}
	return r.Sensitive[root]
}
