package manifest

import (
	"fmt"
	"strings"

	"buf.build/go/protovalidate"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func Validate(manifestPath string) error {
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		// Preserve ManifestLoadError type for beautiful error display
		if IsManifestLoadError(err) {
			return err
		}
		return errors.Wrap(err, "failed to load manifest")
	}

	return ValidateLoaded(manifest)
}

// ValidateLoaded runs protovalidate over the complete manifest message: the
// envelope (apiVersion and kind constants, metadata presence) and every rule
// inside spec. The envelope constants are part of each kind's schema — a
// manifest that declares the wrong apiVersion for its kind is invalid input,
// not decoration — so they are enforced here, on every path that validates a
// manifest. Callers that load manifests from memory use this directly;
// Validate wraps it for the file-path flow.
func ValidateLoaded(manifest proto.Message) error {
	v, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(manifest),
	)
	if err != nil {
		return errors.Wrap(err, "failed to initialize manifest validator")
	}

	validationErr := v.Validate(manifest)
	if validationErr == nil {
		return nil
	}
	return formatValidationError(describeViolations(manifest, validationErr))
}

// describeViolations turns a protovalidate error into user-facing lines.
// Envelope violations (apiVersion, kind) are rephrased into plain language
// naming the provided value, the kind, and the exact fix — the raw rule text
// ("value must equal ...") does not tell the user which line of their manifest
// to change. All other violations keep protovalidate's field-path attribution.
func describeViolations(manifest proto.Message, err error) string {
	var valErr *protovalidate.ValidationError
	if !errors.As(err, &valErr) {
		msg := strings.TrimPrefix(err.Error(), "validation error:")
		return strings.TrimSpace(msg)
	}

	kindName := string(manifest.ProtoReflect().Descriptor().Name())
	var lines []string
	for _, violation := range valErr.Violations {
		path := protovalidate.FieldPathString(violation.Proto.GetField())
		switch path {
		case "api_version":
			expected := envelopeConst(manifest, "api_version")
			got := envelopeValue(manifest, "api_version")
			if got == "" {
				lines = append(lines, fmt.Sprintf(
					"manifest is missing apiVersion: kind %s requires 'apiVersion: %s'",
					kindName, expected))
			} else {
				lines = append(lines, fmt.Sprintf(
					"apiVersion '%s' does not match kind %s: this kind requires 'apiVersion: %s'",
					got, kindName, expected))
			}
		case "kind":
			expected := envelopeConst(manifest, "kind")
			got := envelopeValue(manifest, "kind")
			if got == "" {
				lines = append(lines, fmt.Sprintf(
					"manifest is missing kind: write 'kind: %s'", expected))
			} else {
				lines = append(lines, fmt.Sprintf(
					"kind '%s' is not the canonical name: write 'kind: %s'", got, expected))
			}
		default:
			lines = append(lines, fmt.Sprintf("%s: %s", path, violation.Proto.GetMessage()))
		}
	}
	return strings.Join(lines, "\n   ")
}

func formatValidationError(violationText string) error {
	// Create colored output functions
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	// Build the error message
	var msg strings.Builder

	msg.WriteString("\n")
	msg.WriteString(red("╔═══════════════════════════════════════════════════════════════════════════════╗") + "\n")
	msg.WriteString(red("║") + bold("                    ❌  MANIFEST VALIDATION FAILED                             ") + red("║") + "\n")
	msg.WriteString(red("╚═══════════════════════════════════════════════════════════════════════════════╝") + "\n\n")

	msg.WriteString(yellow("⚠️  Validation Errors:\n\n"))
	msg.WriteString(cyan("   "+violationText) + "\n\n")

	// Generic guidance
	msg.WriteString(bold("💡 Next Steps:\n\n"))
	msg.WriteString("   Please review the validation error messages above and fix the issues\n")
	msg.WriteString("   in your manifest before retrying.\n\n")

	msg.WriteString(bold("📋 Helpful Commands:\n\n"))
	msg.WriteString("   • View current manifest:  " + cyan("planton load-manifest --kustomize-dir _kustomize --overlay prod") + "\n")
	msg.WriteString("   • Validate after fix:     " + cyan("planton validate-manifest --kustomize-dir _kustomize --overlay prod") + "\n")
	msg.WriteString("\n")

	msg.WriteString(bold("📚 Documentation: ") + cyan("https://github.com/plantonhq/planton/tree/main/apis\n"))
	msg.WriteString("\n")

	return errors.New(msg.String())
}
