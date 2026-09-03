package kubeconfig

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/plantonhq/planton/pkg/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
)

// The local workflow reads the operator's kubeconfig exactly as kubectl does:
// KUBECONFIG first (a list, in order), else the default home file, else a
// three-part refusal rather than an engine's connection error.
func TestHostKubeconfigPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	defaultFile := filepath.Join(home, ".kube", "config")

	t.Run("KUBECONFIG wins, as a list in order", func(t *testing.T) {
		a := filepath.Join(home, "a.yaml")
		b := filepath.Join(home, "b.yaml")
		t.Setenv(clientcmd.RecommendedConfigPathEnvVar, a+string(os.PathListSeparator)+b)
		paths, err := HostKubeconfigPaths()
		require.NoError(t, err)
		assert.Equal(t, []string{a, b}, paths)
	})

	t.Run("blank entries in KUBECONFIG are dropped", func(t *testing.T) {
		a := filepath.Join(home, "a.yaml")
		t.Setenv(clientcmd.RecommendedConfigPathEnvVar, string(os.PathListSeparator)+a+string(os.PathListSeparator))
		paths, err := HostKubeconfigPaths()
		require.NoError(t, err)
		assert.Equal(t, []string{a}, paths)
	})

	t.Run("the default home file when KUBECONFIG is unset", func(t *testing.T) {
		t.Setenv(clientcmd.RecommendedConfigPathEnvVar, "")
		require.NoError(t, os.MkdirAll(filepath.Dir(defaultFile), 0o700))
		require.NoError(t, os.WriteFile(defaultFile, []byte("apiVersion: v1\nkind: Config\n"), 0o600))
		t.Cleanup(func() { _ = os.Remove(defaultFile) })
		paths, err := HostKubeconfigPaths()
		require.NoError(t, err)
		assert.Equal(t, []string{defaultFile}, paths)
	})

	t.Run("neither: a three-part refusal naming both places", func(t *testing.T) {
		t.Setenv(clientcmd.RecommendedConfigPathEnvVar, "")
		_, err := HostKubeconfigPaths()
		require.Error(t, err)
		var f *failure.Failure
		require.True(t, errors.As(err, &f), "the error must be a three-part Failure, got %T", err)
		assert.Contains(t, f.Observed, clientcmd.RecommendedConfigPathEnvVar)
		assert.Contains(t, f.Observed, defaultFile)
		assert.Contains(t, f.NextStep, "--kube-context")
	})
}
