package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// tokenSecret materializes the runner token as the `<name>-token` Secret
// the chart reads through its existingSecret form — so the token never
// rides rendered chart values (Helm stores rendered values in its release
// Secret, where an inline token would be readable by anyone who can read
// release history). Created BEFORE the release so the pod's secretKeyRef
// resolves on first schedule.
//
// The token authorizes JOINING and is never the runner's identity: the
// runner presents it once per enrollment, registers itself, and persists
// the identity it receives on its own volume. Rotating the token here
// updates the Secret; running pods keep serving on their minted identity
// (the token is only read at join), and the next pod recreation joins with
// the new token.
func tokenSecret(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (*kubernetescorev1.Secret, error) {
	opts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)

	createdSecret, err := kubernetescorev1.NewSecret(ctx,
		locals.TokenSecretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(
				&kubernetesmeta.ObjectMetaArgs{
					Name:      pulumi.String(locals.TokenSecretName),
					Namespace: pulumi.String(locals.Namespace),
					Labels:    pulumi.ToStringMap(locals.Labels),
				}),
			Type: pulumi.String("Opaque"),
			StringData: pulumi.StringMap{
				vars.TokenSecretKey: pulumi.String(locals.Spec.GetToken()),
			},
		}, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create token secret")
	}

	return createdSecret, nil
}
