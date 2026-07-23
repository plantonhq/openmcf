package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createCredentialSecrets materializes every DECLARED credential in the spec
// as a deterministic Kubernetes Secret, so nothing sensitive ever appears
// inline in a rendered custom resource — the operator and the backup plugin
// only ever see secretKeyRef pointers:
//
//   - `<name>-app-provided` / `<name>-superuser-provided` /
//     `<name>-role-<role>`: kubernetes.io/basic-auth pairs CloudNativePG
//     WATCHES — rotating the value rotates the database password.
//   - `<name>-ext-<external-cluster>`: the password for an external server
//     (the operator builds a passfile from it).
//   - `<name>-backup-creds` / `<name>-recovery-creds` (+ optional
//     `-endpoint-ca`): object-store keys, created by createObjectStores'
//     helpers in backup.go because their keys depend on the backend arm.
//
// Names are deterministic (never engine-generated suffixes) so both engines
// agree byte-for-byte and the import recipes derive them blind.
func createCredentialSecrets(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) ([]pulumi.Resource, error) {
	var created []pulumi.Resource
	spec := locals.Spec

	// Application-owner password (initdb bootstrap): basic-auth with the
	// OWNER's username — CloudNativePG requires both keys and adopts this
	// secret as the application credential instead of generating one.
	if password := spec.GetBootstrap().GetInitdb().GetOwnerPassword(); password != "" {
		owner := spec.GetBootstrap().GetInitdb().GetOwner()
		if owner == "" {
			owner = spec.GetBootstrap().GetInitdb().GetDatabase()
			if owner == "" {
				owner = "app" // the upstream initdb default database/owner
			}
		}
		secret, err := createBasicAuthSecret(ctx, locals, kubernetesProvider, dependencies,
			locals.ProvidedAppSecretName, owner, password)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provided app secret")
		}
		created = append(created, secret)
	}

	// Superuser password: the operator only honors it while superuser
	// access is enabled (the spec CEL enforces the pairing).
	if password := spec.GetSuperuser().GetPassword(); password != "" {
		secret, err := createBasicAuthSecret(ctx, locals, kubernetesProvider, dependencies,
			locals.ProvidedSuperuserSecret, "postgres", password)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provided superuser secret")
		}
		created = append(created, secret)
	}

	// Managed-role passwords.
	for _, role := range spec.GetRoles() {
		if role.GetPassword() == "" {
			continue
		}
		secret, err := createBasicAuthSecret(ctx, locals, kubernetesProvider, dependencies,
			locals.ClusterName+"-role-"+role.GetName(), role.GetName(), role.GetPassword())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create password secret for role %s", role.GetName())
		}
		created = append(created, secret)
	}

	// External-cluster passwords (single `password` key; the operator
	// renders a passfile from it).
	for _, external := range spec.GetExternalClusters() {
		if external.GetPassword() == "" {
			continue
		}
		secretName := locals.ClusterName + "-ext-" + external.GetName()
		secret, err := kubernetescorev1.NewSecret(ctx, secretName,
			&kubernetescorev1.SecretArgs{
				Metadata: kubernetesmeta.ObjectMetaArgs{
					Name:      pulumi.String(secretName),
					Namespace: pulumi.String(locals.Namespace),
					Labels:    pulumi.ToStringMap(locals.Labels),
				},
				StringData: pulumi.StringMap{
					"password": pulumi.String(external.GetPassword()),
				},
			}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create password secret for external cluster %s", external.GetName())
		}
		created = append(created, secret)
	}

	return created, nil
}

// createBasicAuthSecret creates a kubernetes.io/basic-auth Secret — the
// exact shape CloudNativePG expects for role/bootstrap credentials (both
// `username` and `password` keys are mandatory in that format).
func createBasicAuthSecret(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
	secretName, username, password string,
) (pulumi.Resource, error) {
	return kubernetescorev1.NewSecret(ctx, secretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(secretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			},
			Type: pulumi.String("kubernetes.io/basic-auth"),
			StringData: pulumi.StringMap{
				"username": pulumi.String(username),
				"password": pulumi.String(password),
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
}
