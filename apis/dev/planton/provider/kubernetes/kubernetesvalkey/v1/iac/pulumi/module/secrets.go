package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// authSecret materializes the declared ACL passwords as the
// "<metadata.name>-auth" Opaque Secret — ONE KEY PER USERNAME, each key's
// value that user's password. That layout is the chart's contract for
// auth.usersExistingSecret: its init script reads
// /valkey-users-secret/<passwordKey> where passwordKey defaults to the
// username (the module leaves passwordKey unset), and its metrics exporter
// reads the "default" key the same way. Because the rendered aclUsers carry
// no inline passwords, the chart renders no auth Secret of its own — this
// Secret is the only place the credentials land, and it never transits
// chart values. Returns nil when auth is not declared. Terraform twin:
// kubernetes_secret_v1.auth with the same name, keys, and contents.
func authSecret(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	if !locals.AuthEnabled {
		return nil, nil
	}

	data := map[string]string{}
	for _, user := range locals.Spec.GetAuth().GetUsers() {
		data[user.GetName()] = user.GetPassword()
	}

	createdSecret, err := kubernetescorev1.NewSecret(ctx, locals.AuthSecretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.AuthSecretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			Type:       pulumi.String("Opaque"),
			StringData: pulumi.ToStringMap(data),
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create auth secret")
	}

	return createdSecret, nil
}
