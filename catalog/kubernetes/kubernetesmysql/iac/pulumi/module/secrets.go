package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createCredentialSecrets materializes every DECLARED credential in the spec
// as a deterministic Opaque Kubernetes Secret, so nothing sensitive ever
// appears inline in the rendered custom resource — the operator and its
// backup jobs only ever see Secret NAMES:
//
//   - `<name>-user-<username>` (key `password`): a declared user password.
//     The operator watches the Secret the user's passwordSecretRef points
//     at — rotating the value rotates the database password. Users WITHOUT
//     a declared password get no Secret here: the operator generates one.
//   - `<name>-backup-<storage>`: object-store keys for a backup storage.
//     The key names are the operator's own contract — S3 jobs read
//     AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY; Azure jobs read
//     AZURE_STORAGE_ACCOUNT_NAME / AZURE_STORAGE_ACCOUNT_KEY (the account
//     name has NO field in the CR's azure block — it travels only in this
//     Secret). PVC storages need no credentials.
//
// Names are deterministic (never engine-generated suffixes) so both engines
// agree byte-for-byte and the import recipes derive them blind.
func createCredentialSecrets(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) ([]pulumi.Resource, error) {
	var created []pulumi.Resource

	for _, user := range locals.Spec.GetUsers() {
		if user.GetPassword() == "" {
			continue
		}
		secretName := locals.ClusterName + "-user-" + user.GetName()
		// The operator CO-OWNS this Secret's annotations: after applying
		// a password it stamps a `percona.com/<cluster>-<user>-hash`
		// marker (its rotation detector) onto the object. The module
		// owns the data; the annotations are the operator's — never
		// fight them (Terraform twin: lifecycle ignore_changes).
		secret, err := createOpaqueSecret(ctx, locals, kubernetesProvider,
			append(dependencies, pulumi.IgnoreChanges([]string{"metadata.annotations"})),
			secretName, map[string]string{"password": user.GetPassword()})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create password secret for user %s", user.GetName())
		}
		created = append(created, secret)
	}

	for _, storage := range locals.Spec.GetBackup().GetStorages() {
		var data map[string]string
		switch {
		case storage.GetS3() != nil:
			accessKeys := storage.GetS3().GetAccessKeys()
			data = map[string]string{
				"AWS_ACCESS_KEY_ID":     accessKeys.GetAccessKeyId(),
				"AWS_SECRET_ACCESS_KEY": accessKeys.GetSecretAccessKey(),
			}
		case storage.GetAzure() != nil:
			azure := storage.GetAzure()
			data = map[string]string{
				"AZURE_STORAGE_ACCOUNT_NAME": azure.GetStorageAccount(),
				"AZURE_STORAGE_ACCOUNT_KEY":  azure.GetAccessKey(),
			}
		default:
			// PVC storage — backups land on a volume, no credentials.
			continue
		}
		secretName := locals.ClusterName + "-backup-" + storage.GetName()
		secret, err := createOpaqueSecret(ctx, locals, kubernetesProvider, dependencies,
			secretName, data)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create credentials secret for backup storage %s", storage.GetName())
		}
		created = append(created, secret)
	}

	return created, nil
}

// createOpaqueSecret creates a plain Opaque Secret with the given string
// data — the shape the operator's credentialsSecret / passwordSecretRef
// references expect.
func createOpaqueSecret(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
	secretName string, data map[string]string,
) (pulumi.Resource, error) {
	return kubernetescorev1.NewSecret(ctx, secretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(secretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			},
			StringData: pulumi.ToStringMap(data),
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
}
