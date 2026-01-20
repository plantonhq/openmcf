package manifest

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantonhq/project-planton/internal/cli/iacflags"
	"github.com/plantonhq/project-planton/internal/cli/workspace"
	"github.com/plantonhq/project-planton/pkg/clipboard"
	"github.com/plantonhq/project-planton/pkg/ulidgen"
	"github.com/spf13/cobra"
)

// resolveFromClipboard checks for clipboard flags (--clipboard, --clip, --cb, -c) and reads manifest from clipboard.
// Returns empty string if flag not provided.
// When clipboard content is read, it is written to a file in the downloads directory.
func resolveFromClipboard(cmd *cobra.Command) (manifestPath string, isTemp bool, err error) {
	useClipboard, err := iacflags.IsClipboardFlagSet(cmd)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to get clipboard flag")
	}
	if !useClipboard {
		return "", false, nil
	}

	content, err := clipboard.Read()
	if err != nil {
		return "", false, err
	}

	manifestPath, err = writeClipboardContent(content)
	if err != nil {
		return "", false, err
	}

	return manifestPath, true, nil
}

// writeClipboardContent writes content to a file in the downloads directory.
// Follows the same pattern as extract_manifest.go for consistent temp file handling.
func writeClipboardContent(content []byte) (string, error) {
	downloadDir, err := workspace.GetManifestDownloadDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get manifest download directory")
	}

	fileName := ulidgen.NewGenerator().Generate().String() + "-clipboard-manifest.yaml"
	manifestPath := filepath.Join(downloadDir, fileName)

	if err := os.WriteFile(manifestPath, content, 0600); err != nil {
		return "", errors.Wrapf(err, "failed to write manifest to %s", manifestPath)
	}

	return manifestPath, nil
}
