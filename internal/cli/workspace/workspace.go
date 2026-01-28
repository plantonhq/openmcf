package workspace

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/openmcf/pkg/fileutil"
	"os"
	"path/filepath"
)

const (
	OpenMCFDir = ".openmcf"
)

// GetWorkspaceDir returns the path of the openmcf cli workspace directory.
func GetWorkspaceDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get home directory")
	}
	//base directory will always be ${HOME}/.openmcf
	cliWorkspaceDirectory := filepath.Join(homeDir, OpenMCFDir)
	if !fileutil.IsDirExists(cliWorkspaceDirectory) {
		if err := os.MkdirAll(cliWorkspaceDirectory, os.ModePerm); err != nil {
			return "", errors.Wrapf(err, "failed to ensure %s dir", cliWorkspaceDirectory)
		}
	}
	return cliWorkspaceDirectory, nil
}

func GetManifestDownloadDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get home directory")
	}
	dir := filepath.Join(homeDir, OpenMCFDir, "downloads")
	if !fileutil.IsDirExists(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", errors.Wrapf(err, "failed to ensure %s dir", dir)
		}
	}
	return filepath.Join(homeDir, OpenMCFDir, "downloads"), nil
}
