package manifest

import (
	"buf.build/go/protovalidate"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
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

	spec, err := ExtractSpec(manifest)
	if err != nil {
		return errors.Wrap(err, "failed to extract spec from manifest")
	}

	v, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(spec),
	)
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}

	validationErr := v.Validate(spec)
	if validationErr != nil {
		return formatValidationError(validationErr)
	}
	return nil
}

func formatValidationError(err error) error {
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

	// Display the actual validation errors (strip "validation error:" prefix if present)
	errMsg := err.Error()
	errMsg = strings.TrimPrefix(errMsg, "validation error:")
	errMsg = strings.TrimPrefix(errMsg, "validation error:\n")
	errMsg = strings.TrimSpace(errMsg)
	msg.WriteString(cyan("   "+errMsg) + "\n\n")

	// Generic guidance
	msg.WriteString(bold("💡 Next Steps:\n\n"))
	msg.WriteString("   Please review the validation error messages above and fix the issues\n")
	msg.WriteString("   in your manifest before retrying.\n\n")

	msg.WriteString(bold("📋 Helpful Commands:\n\n"))
	msg.WriteString("   • View current manifest:  " + cyan("openmcf load-manifest --kustomize-dir _kustomize --overlay prod") + "\n")
	msg.WriteString("   • Validate after fix:     " + cyan("openmcf validate-manifest --kustomize-dir _kustomize --overlay prod") + "\n")
	msg.WriteString("\n")

	msg.WriteString(bold("📚 Documentation: ") + cyan("https://github.com/plantonhq/openmcf/tree/main/apis\n"))
	msg.WriteString("\n")

	return errors.New(msg.String())
}
