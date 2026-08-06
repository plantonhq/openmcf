package module

import (
	kubernetesopenbaov1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesopenbao/v1alpha1"
)

// Environment variable names the seal wrappers read their credentials
// from (the cloud SDKs' standard variables; transit follows the wrapper's
// token env). These pair with sealPlainEnv/sealSecretData below and with
// the extraSecretEnvironmentVars wiring in values.go.
const (
	envAwsAccessKeyID     = "AWS_ACCESS_KEY_ID"
	envAwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	envAzureClientSecret  = "AZURE_CLIENT_SECRET"
	envTransitToken       = "VAULT_TOKEN"
)

// sealSecretData returns the credential material a declared auto-unseal
// arm carries, keyed by the ENV VAR name it must reach the server as —
// or nil for the keyless postures (workload identity / ambient
// credentials), which need no Secret at all. The module materializes
// this into `<name>-seal-credentials` and wires each key through the
// chart's extraSecretEnvironmentVars, so credential material never
// touches the config ConfigMap or rendered values.
func sealSecretData(spec *kubernetesopenbaov1alpha1.KubernetesOpenBaoSpec) map[string]string {
	seal := spec.GetAutoUnseal()
	if seal == nil {
		return nil
	}
	data := map[string]string{}
	switch {
	case seal.GetAwsKms() != nil:
		if seal.GetAwsKms().GetSecretAccessKey() != "" {
			data[envAwsSecretAccessKey] = seal.GetAwsKms().GetSecretAccessKey()
		}
	case seal.GetAzureKeyVault() != nil:
		if seal.GetAzureKeyVault().GetClientSecret() != "" {
			data[envAzureClientSecret] = seal.GetAzureKeyVault().GetClientSecret()
		}
	case seal.GetTransit() != nil:
		if seal.GetTransit().GetToken() != "" {
			// The transit seal wrapper authenticates to the central
			// instance with a token read from the standard token env
			// var. The server itself never self-authenticates, so the
			// variable is inert beyond the seal wrapper.
			data[envTransitToken] = seal.GetTransit().GetToken()
		}
	}
	if len(data) == 0 {
		return nil
	}
	return data
}

// sealPlainEnv returns NON-secret environment variables a declared seal
// arm needs (identifiers, not credentials) — rendered through the
// chart's plain extraEnvironmentVars.
func sealPlainEnv(spec *kubernetesopenbaov1alpha1.KubernetesOpenBaoSpec) map[string]string {
	seal := spec.GetAutoUnseal()
	if seal == nil {
		return nil
	}
	env := map[string]string{}
	if aws := seal.GetAwsKms(); aws != nil {
		if aws.GetAccessKeyId() != "" {
			// An access-key ID is a public identifier; only its paired
			// secret key is credential material (which rides the
			// Secret).
			env[envAwsAccessKeyID] = aws.GetAccessKeyId()
		}
		env["AWS_REGION"] = aws.GetRegion()
	}
	if len(env) == 0 {
		return nil
	}
	return env
}
