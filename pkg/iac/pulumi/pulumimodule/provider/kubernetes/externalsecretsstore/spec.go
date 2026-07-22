// Package externalsecretsstore renders the shared ExternalSecretsStoreConfig
// proto into the external-secrets.io/v1 SecretStoreSpec JSON shape plus the
// credential Secrets the modules must materialize.
//
// SecretStore and ClusterSecretStore share an IDENTICAL spec upstream, and
// the Planton kinds share the ExternalSecretsStoreConfig proto — this
// package is the third leg of that guarantee: ONE builder renders the CR
// spec for both kinds, so their behavior can never drift. The Terraform
// modules render the same shape through their locals; keep field mappings in
// lockstep.
//
// CREDENTIAL MODEL: wherever the CRD expects a secretRef, the proto carries
// the credential VALUE (sensitive). The builder returns the Secrets to
// create (one deterministic "<resource-name>-credentials" Secret, keys fixed
// per backend) and wires the CR's secretRefs to it. Cluster-scoped stores
// read referenced Secrets from an EXPLICIT namespace, so the builder stamps
// the secrets namespace into every ref when rendering for a
// ClusterSecretStore; namespaced SecretStores omit it (refs default to the
// store's own namespace).
package externalsecretsstore

import (
	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
)

// CredentialSecret is one Kubernetes Secret a module must create in the
// store's secrets namespace before applying the CR.
type CredentialSecret struct {
	// Secret metadata.name (deterministic: "<resource-name>-credentials").
	Name string
	// stringData entries.
	Data map[string]string
}

// Result carries everything a module needs to apply the store.
type Result struct {
	// The external-secrets.io/v1 SecretStore/ClusterSecretStore spec, in
	// CRD JSON shape (the provider block plus common tuning).
	Spec map[string]interface{}
	// Credential Secrets to materialize before the CR.
	Secrets []CredentialSecret
}

// BuildSpec renders the store config. resourceName is the CR's
// metadata.name (it derives the credential Secret name); secretsNamespace is
// where declared credentials materialize; clusterScoped selects whether
// secret references carry an explicit namespace (ClusterSecretStore) or
// default to the store's own (SecretStore).
func BuildSpec(resourceName, secretsNamespace string, clusterScoped bool, config *kubernetesprovider.ExternalSecretsStoreConfig) (*Result, error) {
	result := &Result{Spec: map[string]interface{}{}}
	credentialSecretName := resourceName + "-credentials"
	credentialData := map[string]string{}

	// secretRef renders a CRD SecretKeySelector against the credential
	// Secret, carrying the explicit namespace only for cluster-scoped
	// stores.
	secretRef := func(key string) map[string]interface{} {
		ref := map[string]interface{}{
			"name": credentialSecretName,
			"key":  key,
		}
		if clusterScoped {
			ref["namespace"] = secretsNamespace
		}
		return ref
	}

	// serviceAccountRef renders a CRD ServiceAccountSelector. On a
	// ClusterSecretStore the namespace is REQUIRED by the webhook (cluster
	// scope has no home namespace); on a SecretStore it defaults to the
	// store's own namespace and an explicit override is unusual, so it is
	// carried only when the spec sets one.
	serviceAccountRef := func(name, namespace string) map[string]interface{} {
		ref := map[string]interface{}{"name": name}
		if namespace != "" {
			ref["namespace"] = namespace
		} else if clusterScoped {
			ref["namespace"] = secretsNamespace
		}
		return ref
	}

	provider := map[string]interface{}{}

	switch {
	case config.GetAws() != nil:
		aws := config.GetAws()
		awsSpec := map[string]interface{}{
			"service": aws.GetService(),
			"region":  aws.GetRegion(),
		}
		if aws.GetRole().GetValue() != "" {
			awsSpec["role"] = aws.GetRole().GetValue()
		}
		if aws.GetPrefix() != "" {
			awsSpec["prefix"] = aws.GetPrefix()
		}
		switch {
		case aws.GetAccessKeyId() != "":
			credentialData["access-key-id"] = aws.GetAccessKeyId()
			credentialData["secret-access-key"] = aws.GetSecretAccessKey()
			awsSpec["auth"] = map[string]interface{}{
				"secretRef": map[string]interface{}{
					"accessKeyIDSecretRef":     secretRef("access-key-id"),
					"secretAccessKeySecretRef": secretRef("secret-access-key"),
				},
			}
		case aws.GetServiceAccountName().GetValue() != "":
			awsSpec["auth"] = map[string]interface{}{
				"jwt": map[string]interface{}{
					"serviceAccountRef": serviceAccountRef(
						aws.GetServiceAccountName().GetValue(), aws.GetServiceAccountNamespace()),
				},
			}
		}
		// No auth block at all = the operator's ambient identity (its
		// controller ServiceAccount / node role) — upstream's documented
		// fallback.
		provider["aws"] = awsSpec

	case config.GetGcpSecretManager() != nil:
		gcp := config.GetGcpSecretManager()
		gcpSpec := map[string]interface{}{
			"projectID": gcp.GetProjectId().GetValue(),
		}
		if gcp.GetLocation() != "" {
			gcpSpec["location"] = gcp.GetLocation()
		}
		switch {
		case gcp.GetServiceAccountKeyJson() != "":
			credentialData["credentials.json"] = gcp.GetServiceAccountKeyJson()
			gcpSpec["auth"] = map[string]interface{}{
				"secretRef": map[string]interface{}{
					"secretAccessKeySecretRef": secretRef("credentials.json"),
				},
			}
		case gcp.GetServiceAccountName().GetValue() != "":
			gcpSpec["auth"] = map[string]interface{}{
				"workloadIdentity": map[string]interface{}{
					"serviceAccountRef": serviceAccountRef(
						gcp.GetServiceAccountName().GetValue(), gcp.GetServiceAccountNamespace()),
				},
			}
		}
		provider["gcpsm"] = gcpSpec

	case config.GetAzureKeyVault() != nil:
		az := config.GetAzureKeyVault()
		azSpec := map[string]interface{}{
			"vaultUrl": az.GetVaultUrl(),
		}
		if az.AuthType != nil {
			azSpec["authType"] = az.GetAuthType()
		}
		if az.GetTenantId() != "" {
			azSpec["tenantId"] = az.GetTenantId()
		}
		if az.GetIdentityId() != "" {
			azSpec["identityId"] = az.GetIdentityId()
		}
		if az.GetServiceAccountName().GetValue() != "" {
			azSpec["serviceAccountRef"] = serviceAccountRef(
				az.GetServiceAccountName().GetValue(), az.GetServiceAccountNamespace())
		}
		if az.GetClientId() != "" {
			credentialData["client-id"] = az.GetClientId()
			credentialData["client-secret"] = az.GetClientSecret()
			azSpec["authSecretRef"] = map[string]interface{}{
				"clientId":     secretRef("client-id"),
				"clientSecret": secretRef("client-secret"),
			}
		}
		provider["azurekv"] = azSpec

	case config.GetVault() != nil:
		vault := config.GetVault()
		vaultSpec := map[string]interface{}{
			"server": vault.GetServer(),
		}
		if vault.GetPath() != "" {
			vaultSpec["path"] = vault.GetPath()
		}
		if vault.Version != nil {
			vaultSpec["version"] = vault.GetVersion()
		}
		if vault.GetNamespace() != "" {
			vaultSpec["namespace"] = vault.GetNamespace()
		}
		if vault.GetCaBundle() != "" {
			// The CRD carries caBundle as base64 []byte JSON; passing the
			// PEM through Kubernetes' byte-encoding is handled by the
			// engines' JSON marshalling of string→base64 on apply — both
			// engines pass the SAME base64-encoded PEM string.
			vaultSpec["caBundle"] = vault.GetCaBundle()
		}

		auth := map[string]interface{}{}
		switch {
		case vault.GetToken() != nil:
			credentialData["vault-token"] = vault.GetToken().GetToken()
			auth["tokenSecretRef"] = secretRef("vault-token")
		case vault.GetAppRole() != nil:
			appRole := vault.GetAppRole()
			credentialData["vault-approle-secret-id"] = appRole.GetSecretId()
			approleSpec := map[string]interface{}{
				"roleId":    appRole.GetRoleId(),
				"secretRef": secretRef("vault-approle-secret-id"),
			}
			if appRole.GetPath() != "" {
				approleSpec["path"] = appRole.GetPath()
			}
			auth["appRole"] = approleSpec
		case vault.GetKubernetes() != nil:
			k8sAuth := vault.GetKubernetes()
			kubernetesSpec := map[string]interface{}{
				"role": k8sAuth.GetRole(),
			}
			if k8sAuth.GetMountPath() != "" {
				kubernetesSpec["mountPath"] = k8sAuth.GetMountPath()
			}
			if k8sAuth.GetServiceAccountName().GetValue() != "" {
				kubernetesSpec["serviceAccountRef"] = serviceAccountRef(
					k8sAuth.GetServiceAccountName().GetValue(), "")
			}
			auth["kubernetes"] = kubernetesSpec
		default:
			return nil, errors.New("vault backend must select an auth method (token, app_role, or kubernetes)")
		}
		vaultSpec["auth"] = auth
		provider["vault"] = vaultSpec

	case config.GetKubernetes() != nil:
		k8s := config.GetKubernetes()
		k8sSpec := map[string]interface{}{}
		server := map[string]interface{}{}
		if k8s.GetServerUrl() != "" {
			server["url"] = k8s.GetServerUrl()
		}
		if k8s.GetCaBundle() != "" {
			server["caBundle"] = k8s.GetCaBundle()
		}
		if len(server) > 0 {
			k8sSpec["server"] = server
		}
		if k8s.GetRemoteNamespace() != "" {
			k8sSpec["remoteNamespace"] = k8s.GetRemoteNamespace()
		}
		switch {
		case k8s.GetToken() != "":
			credentialData["bearer-token"] = k8s.GetToken()
			k8sSpec["auth"] = map[string]interface{}{
				"token": map[string]interface{}{
					"bearerToken": secretRef("bearer-token"),
				},
			}
		case k8s.GetServiceAccountName().GetValue() != "":
			k8sSpec["auth"] = map[string]interface{}{
				"serviceAccount": serviceAccountRef(k8s.GetServiceAccountName().GetValue(), ""),
			}
		}
		provider["kubernetes"] = k8sSpec

	case config.GetFake() != nil:
		entries := make([]interface{}, 0, len(config.GetFake().GetData()))
		for _, entry := range config.GetFake().GetData() {
			rendered := map[string]interface{}{
				"key":   entry.GetKey(),
				"value": entry.GetValue(),
			}
			if entry.GetVersion() != "" {
				rendered["version"] = entry.GetVersion()
			}
			entries = append(entries, rendered)
		}
		provider["fake"] = map[string]interface{}{"data": entries}

	default:
		return nil, errors.New("store config must select a backend (aws, gcp_secret_manager, azure_key_vault, vault, kubernetes, or fake)")
	}

	result.Spec["provider"] = provider

	// ---- common tuning ---------------------------------------------------
	if config.GetControllerClass() != "" {
		result.Spec["controller"] = config.GetControllerClass()
	}
	if config.GetRefreshInterval() != "" {
		result.Spec["refreshInterval"] = config.GetRefreshInterval()
	}
	if retry := config.GetRetry(); retry != nil {
		retrySettings := map[string]interface{}{}
		if retry.MaxRetries != nil {
			retrySettings["maxRetries"] = int(retry.GetMaxRetries())
		}
		if retry.GetRetryInterval() != "" {
			retrySettings["retryInterval"] = retry.GetRetryInterval()
		}
		if len(retrySettings) > 0 {
			result.Spec["retrySettings"] = retrySettings
		}
	}

	if len(credentialData) > 0 {
		result.Secrets = append(result.Secrets, CredentialSecret{
			Name: credentialSecretName,
			Data: credentialData,
		})
	}

	return result, nil
}
