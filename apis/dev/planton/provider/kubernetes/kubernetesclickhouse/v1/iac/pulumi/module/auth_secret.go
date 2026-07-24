package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// authSecret creates the `<metadata.name>-clickhouse-auth` Secret with one
// key per provisioned user holding that user's password (resolved literal
// from the spec's value-or-ref). The CHI users section references it via
// valueFrom.secretKeyRef — the password never appears in the CHI itself.
// Returns nil when spec.users is empty.
func authSecret(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (*kubernetescorev1.Secret, error) {
	users := locals.Spec.GetUsers()
	if len(users) == 0 {
		return nil, nil
	}

	passwords := pulumi.StringMap{}
	for _, user := range users {
		passwords[user.GetName()] = pulumi.String(user.GetPassword().GetValue())
	}

	createdSecret, err := kubernetescorev1.NewSecret(ctx,
		locals.AuthSecretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(
				&kubernetesmeta.ObjectMetaArgs{
					Name:      pulumi.String(locals.AuthSecretName),
					Namespace: pulumi.String(locals.Namespace),
					Labels:    pulumi.ToStringMap(locals.Labels),
				}),
			// StringData lets the apiserver base64-encode; ToSecret keeps
			// the plaintext passwords out of rendered previews and state.
			StringData: pulumi.ToSecret(passwords).(pulumi.StringMapOutput),
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create %s secret", locals.AuthSecretName)
	}

	return createdSecret, nil
}
