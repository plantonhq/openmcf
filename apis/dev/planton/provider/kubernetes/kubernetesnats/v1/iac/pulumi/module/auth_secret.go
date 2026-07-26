package module

import (
	"fmt"

	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// authSecret generates one random password per declared user (flat and
// account users alike) and materializes them as the
// "<metadata.name>-auth" Opaque Secret — ONE KEY PER USERNAME, each key's
// value that user's password. That Secret is the ONLY place the
// credentials land: the chart values carry secretKeyRef env wiring plus
// `$NATS_PW_<i>` config references the server resolves from environment
// at load — nothing credential-bearing renders into values, the
// ConfigMap, or plan/preview output. Clients read the same Secret (the
// exported auth_secret_name handle). Returns nil when auth is not
// declared. Terraform twin: random_password.user +
// kubernetes_secret_v1.auth with the same name, keys, and contents.
func authSecret(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	if !locals.AuthEnabled {
		return nil, nil
	}

	// The generation-shape arguments are ignored after creation so an
	// IMPORTED credential never silently regenerates: rotation stays an
	// explicit verb, never plan fallout.
	generationShapeIgnores := pulumi.IgnoreChanges([]string{
		"length", "special", "upper", "lower", "numeric",
		"minLower", "minNumeric", "minSpecial", "minUpper", "overrideSpecial",
	})

	data := pulumi.StringMap{}
	allUsers := append([]authUser{}, locals.FlatUsers...)
	for _, users := range locals.AccountUsers {
		allUsers = append(allUsers, users...)
	}
	for _, u := range allUsers {
		createdPassword, err := random.NewRandomPassword(ctx,
			fmt.Sprintf("auth-password-%s", u.Username),
			&random.RandomPasswordArgs{
				Length:  pulumi.Int(24),
				Special: pulumi.Bool(false),
			},
			generationShapeIgnores)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate password for user %s", u.Username)
		}
		data[u.Username] = createdPassword.Result
	}

	createdSecret, err := kubernetescorev1.NewSecret(ctx, locals.AuthSecretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.AuthSecretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			Type:       pulumi.String("Opaque"),
			StringData: data,
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create auth secret")
	}

	return createdSecret, nil
}
