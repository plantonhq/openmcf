package tofumodule

import (
	"github.com/plantonhq/planton/pkg/outputs"
)

// RunOption customizes a RunCommand invocation. Options are additive by
// design: the exported RunCommand signature is a contract shared by the
// unified verbs, the standalone tofu/terraform commands, and external
// callers, so new capabilities arrive as options rather than parameters.
type RunOption func(*runConfig)

// runConfig collects the effects of the applied options.
type runConfig struct {
	// captureSink, when non-nil, asks RunCommand to read the stack's outputs
	// back after a successful apply and fill the sink. Only the apply
	// operation captures; plan, destroy, and refresh ignore the sink.
	captureSink *outputs.CaptureResult
}

// WithOutputCapture asks RunCommand to capture the stack's outputs into sink
// after a successful apply. Capture runs `<binary> output -json` in the module
// workspace while it is still alive, so remote-state credentials and provider
// overrides from the apply still hold. A capture failure never fails the
// apply — the infrastructure is already deployed; the failure is reported and
// the sink is left partially filled (Raw/Flat may be set even when the typed
// transform failed).
func WithOutputCapture(sink *outputs.CaptureResult) RunOption {
	return func(cfg *runConfig) {
		cfg.captureSink = sink
	}
}
