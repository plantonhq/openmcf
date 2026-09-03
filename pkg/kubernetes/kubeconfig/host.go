package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/plantonhq/planton/pkg/failure"
	"k8s.io/client-go/tools/clientcmd"
)

// HostKubeconfigPaths resolves the kubeconfig files a module uses when it runs
// WITHOUT a Planton connection (the local workflow): the operator's own
// machine, the same files kubectl reads. The rule is kubectl's: the KUBECONFIG
// environment variable when set (a list separated by the OS path separator,
// every named file kept in order), else the default file under the home
// directory when it exists.
//
// When neither exists there is no cluster to talk to, and the returned error
// explains that in three parts before any engine starts; the Terraform
// providers, left to themselves, would report a connection refused to
// localhost, which names the symptom and not the cause.
func HostKubeconfigPaths() ([]string, error) {
	if raw := strings.TrimSpace(os.Getenv(clientcmd.RecommendedConfigPathEnvVar)); raw != "" {
		var paths []string
		for _, p := range filepath.SplitList(raw) {
			if strings.TrimSpace(p) != "" {
				paths = append(paths, p)
			}
		}
		if len(paths) > 0 {
			return paths, nil
		}
	}

	defaultPath := defaultHostKubeconfigPath()
	if _, err := os.Stat(defaultPath); err == nil {
		return []string{defaultPath}, nil
	}

	return nil, &failure.Failure{
		Observed: fmt.Sprintf("no kubeconfig was found: the %s environment variable is not set and %s does not exist", clientcmd.RecommendedConfigPathEnvVar, defaultPath),
		Meaning:  "this deploy has no Planton connection to a cluster, so it uses this machine's kubeconfig the way kubectl does, and there is none to use",
		NextStep: fmt.Sprintf("export %s=/path/to/kubeconfig (or create %s), then re-run; pass --kube-context <name> to pick a context other than the current one", clientcmd.RecommendedConfigPathEnvVar, defaultPath),
	}
}

// defaultHostKubeconfigPath is kubectl's default file, resolved from the
// home directory at call time (clientcmd computes its copy once at process
// start, which a test cannot redirect).
func defaultHostKubeconfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return clientcmd.RecommendedHomeFile
	}
	return filepath.Join(home, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName)
}
