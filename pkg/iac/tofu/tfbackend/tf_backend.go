package tfbackend

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/shared/iac/terraform"
)

// WriteBackendFile creates a `backend.tf` file in projectDir declaring the state backend.
// The enum value's NAME is the HCL backend name (`backend "s3"`, `backend "remote"`), so
// TerraformBackendType entries must be spelled exactly as the engine spells the backend.
//
// body lines, when given, are rendered inside the backend block, indented one level.
// Most backends need no body -- their settings ride `-backend-config` key=value flags.
// The exception is the remote (TFE-protocol) backend: its `workspaces { name = "..." }`
// is an HCL BLOCK, and blocks cannot be expressed as -backend-config flags (the engine's
// override mechanism looks up attributes only), so the workspace identity belongs here,
// in the declaration itself. With no body the output is byte-identical to the historical
// `backend "<type>" {}` form.
func WriteBackendFile(projectDir string, tofuBackendType terraform.TerraformBackendType, body ...string) error {
	backendName := tofuBackendType.String()

	inner := ""
	if len(body) > 0 {
		inner = "\n"
		for _, line := range body {
			inner += "    " + line + "\n"
		}
		inner += "  "
	}

	backendContent := fmt.Sprintf(`terraform {
  backend "%s" {%s}
}
`, backendName, inner)

	backendFilePath := filepath.Join(projectDir, "backend.tf")
	if err := os.WriteFile(backendFilePath, []byte(backendContent), 0644); err != nil {
		return errors.Wrap(err, "failed to write backend file")
	}

	return nil
}

// WorkspacesNameBody renders the body lines for a remote backend pinned to a single
// named workspace -- the standard TFE state-storage shape. Callers pass the result to
// WriteBackendFile (or through tofumodule.Init) as the backend body.
func WorkspacesNameBody(workspaceName string) []string {
	return []string{
		"workspaces {",
		fmt.Sprintf("  name = %q", workspaceName),
		"}",
	}
}

func BackendTypeFromString(backendTypeStr string) terraform.TerraformBackendType {
	switch backendTypeStr {
	case "local":
		return terraform.TerraformBackendType_local
	case "s3":
		return terraform.TerraformBackendType_s3
	case "gcs":
		return terraform.TerraformBackendType_gcs
	case "azurerm":
		return terraform.TerraformBackendType_azurerm
	case "remote":
		return terraform.TerraformBackendType_remote
	default:
		return terraform.TerraformBackendType_terraform_backend_type_unspecified
	}
}
