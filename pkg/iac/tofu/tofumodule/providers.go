package tofumodule

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/iac/stackinput/providerenvvars"
)

// GetProviderConfigEnvVars returns provider-specific environment variables for the stack.
// It delegates to the IaC-agnostic providerenvvars package which determines the correct provider
// based on the target's api_version/kind and loads only the relevant provider configuration.
func GetProviderConfigEnvVars(stackInputYaml, fileCacheLoc, kubeContext string) ([]string, error) {
	// ResolveAwsWebIdentity is set here because this is the tofu/terraform execution boundary:
	// the AWS modules' `provider "aws" {}` block is empty, so keyless connections must have
	// their web-identity JWT exchanged for temporary credentials and injected as AWS_* env
	// vars. The pulumi path calls GetEnvVarsWithOptions directly and leaves this false.
	// KubeContext rides the same options: the loader exports it as KUBE_CTX
	// beside the kubeconfig it resolves (a connection's rendered file, or the
	// operator's own kubeconfig when there is no connection), so one seam owns
	// every name the kubernetes and helm providers read.
	providerConfigEnvVars, err := providerenvvars.GetEnvVarsWithOptions(stackInputYaml, providerenvvars.Options{
		FileCacheLoc:          fileCacheLoc,
		ResolveAwsWebIdentity: true,
		KubeContext:           kubeContext,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get provider env vars from stack input")
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
