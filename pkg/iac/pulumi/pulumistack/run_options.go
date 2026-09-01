package pulumistack

import (
	"github.com/plantonhq/planton/pkg/outputs"
)

// RunOption customizes a Run invocation. Options are additive by design: the
// exported Run signature is a contract shared by the unified verbs and the
// standalone pulumi commands, so new capabilities arrive as options rather
// than parameters.
type RunOption func(*runConfig)

// runConfig collects the effects of the applied options.
type runConfig struct {
	// captureSink, when non-nil, asks Run to read the stack's outputs back
	// after a successful update and fill the sink. Only the update operation
	// (and never a preview) captures.
	captureSink *outputs.CaptureResult

	// backendUrl, when non-empty, is the caller-resolved backend URL
	// (typically from a CLI flag). It participates in the backend-URL
	// precedence chain: this option > manifest annotation > environment.
	backendUrl string
}

// WithOutputCapture asks Run to capture the stack's outputs into sink after a
// successful update. Capture reads the stack twice — once masked to learn
// which outputs are secret, once with --show-secrets for the real values —
// because pulumi's plain JSON output carries no per-output sensitivity flag.
// A capture failure never fails the update — the infrastructure is already
// deployed; the failure is reported and the sink is left partially filled.
func WithOutputCapture(sink *outputs.CaptureResult) RunOption {
	return func(cfg *runConfig) {
		cfg.captureSink = sink
	}
}

// WithBackendURL supplies a caller-resolved pulumi backend URL (s3://...,
// gs://..., azblob://..., file://..., or a Pulumi Cloud URL). It wins over
// the manifest annotation and the environment variable, mirroring the
// flags-beat-annotations-beat-env precedence the tofu backend follows.
func WithBackendURL(url string) RunOption {
	return func(cfg *runConfig) {
		cfg.backendUrl = url
	}
}
