package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// credentialSecrets materializes the selected provider's DECLARED credentials
// as deterministically-named Kubernetes Secrets, so the credential itself
// never appears in chart values or pod specs — the chart wires env/volume
// references to these Secrets (values.go). Providers running keyless (or the
// webhook/in-memory arms) materialize nothing. Terraform twin:
// kubernetes_secret_v1 resources with the same names, keys, and contents.
func credentialSecrets(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) ([]pulumi.Resource, error) {
	var created []pulumi.Resource

	newSecret := func(secretName string, data map[string]string) error {
		secret, err := kubernetescorev1.NewSecret(ctx, secretName,
			&kubernetescorev1.SecretArgs{
				Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
					Name:      pulumi.String(secretName),
					Namespace: pulumi.String(locals.Namespace),
					Labels:    pulumi.ToStringMap(locals.Labels),
				}),
				StringData: pulumi.ToStringMap(data),
			}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)...)
		if err != nil {
			return err
		}
		created = append(created, secret)
		return nil
	}

	spec := locals.Spec

	// Cloudflare: token consumed via CF_API_TOKEN (values.go env wiring).
	if cf := spec.GetCloudflare(); cf != nil {
		if err := newSecret(locals.CloudflareSecretName, map[string]string{
			"api-token": cf.GetApiToken(),
		}); err != nil {
			return nil, errors.Wrap(err, "failed to create cloudflare credential secret")
		}
	}

	// AWS static keys: consumed via AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY.
	// Keyless installs (workload identity / node role) leave both empty and
	// materialize nothing.
	if aws := spec.GetAwsRoute53(); aws != nil && aws.GetAccessKeyId() != "" {
		if err := newSecret(locals.AwsSecretName, map[string]string{
			"access-key-id":     aws.GetAccessKeyId(),
			"secret-access-key": aws.GetSecretAccessKey(),
		}); err != nil {
			return nil, errors.Wrap(err, "failed to create aws credential secret")
		}
	}

	// GCP service-account key: mounted as a file with
	// GOOGLE_APPLICATION_CREDENTIALS pointing at it (ADC's file path form).
	if gcp := spec.GetGoogleCloudDns(); gcp != nil && gcp.GetServiceAccountKeyJson() != "" {
		if err := newSecret(locals.GcpSecretName, map[string]string{
			"credentials.json": gcp.GetServiceAccountKeyJson(),
		}); err != nil {
			return nil, errors.Wrap(err, "failed to create gcp credential secret")
		}
	}

	// Azure: the controller reads EVERYTHING (including identity mode) from
	// a mounted azure.json — materialize it from the typed fields. Mounted
	// at the controller's default config path, so no --azure-config-file
	// override is needed.
	if az := spec.GetAzureDns(); az != nil {
		azureJson, err := renderAzureJson(locals)
		if err != nil {
			return nil, err
		}
		if err := newSecret(locals.AzureSecretName, map[string]string{
			"azure.json": azureJson,
		}); err != nil {
			return nil, errors.Wrap(err, "failed to create azure config secret")
		}
	}

	return created, nil
}

// renderAzureJson builds the controller's azure.json from the typed spec.
// Identity selection order (mirrored in the Terraform twin and documented on
// the spec): service principal when client_id/client_secret are set;
// otherwise Workload Identity when spec.workload_identity.aks is set;
// otherwise managed identity (user-assigned when
// managed_identity_client_id is set, else system-assigned).
func renderAzureJson(locals *Locals) (string, error) {
	az := locals.Spec.GetAzureDns()

	config := map[string]interface{}{
		"subscriptionId": az.GetSubscriptionId(),
		"resourceGroup":  az.GetResourceGroup(),
	}
	if az.GetTenantId() != "" {
		config["tenantId"] = az.GetTenantId()
	}

	switch {
	case az.GetClientId() != "":
		config["aadClientId"] = az.GetClientId()
		config["aadClientSecret"] = az.GetClientSecret()
	case locals.Spec.GetWorkloadIdentity().GetAks() != nil:
		config["useWorkloadIdentityExtension"] = true
	default:
		config["useManagedIdentityExtension"] = true
		if az.GetManagedIdentityClientId() != "" {
			config["userAssignedIdentityID"] = az.GetManagedIdentityClientId()
		}
	}

	rendered, err := json.Marshal(config)
	if err != nil {
		return "", errors.Wrap(err, "failed to render azure.json")
	}
	return string(rendered), nil
}
