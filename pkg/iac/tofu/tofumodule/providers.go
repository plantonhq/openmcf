package tofumodule

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/project-planton/pkg/iac/stackinput/providerenvvars"
)

// GetProviderConfigEnvVars returns provider-specific environment variables for the stack.
// It delegates to the IaC-agnostic providerenvvars package which determines the correct provider
// based on the target's api_version/kind and loads only the relevant provider configuration.
func GetProviderConfigEnvVars(stackInputYaml, fileCacheLoc, kubeContext string) ([]string, error) {
	providerConfigEnvVars, err := providerenvvars.GetEnvVarsWithOptions(stackInputYaml, providerenvvars.Options{
		FileCacheLoc: fileCacheLoc,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get provider env vars from stack input")
	}

	// Add KUBE_CTX environment variable if kube context is specified
	if kubeContext != "" {
		providerConfigEnvVars["KUBE_CTX"] = kubeContext
	}

	return mapToSlice(providerConfigEnvVars), nil
}

// mapToSlice converts a map of string to string into a slice of string slices by joining key-value pairs with an equals sign.
func mapToSlice(inputMap map[string]string) []string {
	var result []string
	for key, value := range inputMap {
		result = append(result, key+"="+value)
	}
	return result
}
