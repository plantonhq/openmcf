package kubeconfig

import (
	"os"

	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/catalog/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeContextEnvVar selects a context from the host kubeconfig when a module
// runs without a Planton connection (the local workflow). The Pulumi provider
// getter honours the same variable, so an in-process client and the provider
// always talk to the same cluster.
const KubeContextEnvVar = "KUBE_CTX"

// RESTConfig resolves a client-go REST config for the cluster a module is
// about to deploy to. It mirrors the two branches the Pulumi provider getter
// takes: a nil provider config means the host kubeconfig (with KUBE_CTX
// context selection), anything else is rendered by Build. Modules that need
// to READ the cluster in-process before registering resources (a never-
// downgrade check on kept CRDs, for example) use this so the read and the
// provider can never disagree about which cluster they address.
func RESTConfig(config *kubernetesprovider.KubernetesProviderConfig, credentialCommand string) (*rest.Config, error) {
	if config == nil {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		overrides := &clientcmd.ConfigOverrides{}
		if kubeContext := os.Getenv(KubeContextEnvVar); kubeContext != "" {
			overrides.CurrentContext = kubeContext
		}
		restConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides).ClientConfig()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load the host kubeconfig")
		}
		return restConfig, nil
	}

	kubeConfigString, err := Build(config, credentialCommand)
	if err != nil {
		return nil, err
	}
	restConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfigString))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the rendered kubeconfig")
	}
	return restConfig, nil
}
