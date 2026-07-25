package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// adminAuthSecret materializes the console credentials as the
// "<metadata.name>-admin-auth" Opaque Secret (keys user/password, user
// "admin", random 24-char password held in the Pulumi state as a secret).
// The chart consumes it via admin.secret.existingSecret + userKey/pwKey —
// its OWN adminPassword value path is never used, so the credential never
// transits chart values. Created only when the console is enabled without
// an existing Secret; returns nil otherwise. Terraform twin:
// random_password + kubernetes_secret_v1.admin_auth with the same name and
// keys.
func adminAuthSecret(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	if !locals.CreateAdminSecret {
		return nil, nil
	}

	// The generation-shape arguments are ignored after creation so an
	// IMPORTED credential never silently regenerates: random_password's
	// import carries only the VALUE, the importer assumes the provider's
	// own generation defaults (special=true — verified live), and every
	// argument forces replacement — without this, the first update after
	// an import proposes replacing (rotating) the live credential.
	// Rotation stays an explicit verb, never plan fallout.
	createdPassword, err := random.NewRandomPassword(ctx, "admin-auth-password",
		&random.RandomPasswordArgs{
			Length:  pulumi.Int(24),
			Special: pulumi.Bool(false),
		},
		pulumi.IgnoreChanges([]string{
			"length", "special", "upper", "lower", "numeric",
			"minLower", "minNumeric", "minSpecial", "minUpper", "overrideSpecial",
		}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate admin console password")
	}

	createdSecret, err := kubernetescorev1.NewSecret(ctx, locals.AdminAuthSecretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.AdminAuthSecretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			Type: pulumi.String("Opaque"),
			StringData: pulumi.StringMap{
				"user":     pulumi.String("admin"),
				"password": createdPassword.Result,
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create admin auth secret")
	}

	return createdSecret, nil
}
