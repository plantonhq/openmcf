package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The generation-shape arguments are ignored after creation so an
// IMPORTED credential never silently regenerates: rotation stays an
// explicit verb, never plan fallout. Terraform twin: random_password
// lifecycle ignore_changes on the same arguments.
var generationShapeIgnores = pulumi.IgnoreChanges([]string{
	"length", "special", "upper", "lower", "numeric",
	"minLower", "minNumeric", "minSpecial", "minUpper", "overrideSpecial",
})

// webhookKeystoreSecret generates the webhook keystore password and
// materializes the module-owned Secret the chart's
// webhook.keystore.passwordSecretRef points at — the replacement for the
// chart's default Secret, which carries a HARDCODED PUBLIC PASSWORD
// ("password1234", base64 in templates/webhook/secret.yaml behind
// keystore.useDefaultPassword=true) and must never ship.
//
// Webhook-enabled installs only — the disabled arm creates no
// random/Secret resources at all (returns nil). Terraform equivalent:
// random_password + kubernetes_secret_v1, both count-gated on
// webhook_enabled.
func webhookKeystoreSecret(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	deps []pulumi.Resource,
) (pulumi.Resource, error) {
	if !locals.WebhookEnabled {
		return nil, nil
	}

	// Letters+digits only: letters-only-safe alphabets avoid whole
	// config-parser bug classes (the credential lands in a JVM keystore
	// env value); the 32-char length compensates the smaller alphabet.
	keystorePassword, err := random.NewRandomPassword(ctx, "webhook-keystore-password",
		&random.RandomPasswordArgs{
			Length:  pulumi.Int(32),
			Special: pulumi.Bool(false),
		},
		generationShapeIgnores)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate webhook keystore password")
	}

	createdSecret, err := kubernetescorev1.NewSecret(ctx, locals.WebhookKeystoreSecretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.WebhookKeystoreSecretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			Type: pulumi.String("Opaque"),
			StringData: pulumi.StringMap{
				// "password" is the key the chart's passwordSecretRef
				// values point at — kept in lockstep with values.go.
				"password": keystorePassword.Result,
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)},
			dependsOn(deps)...)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create webhook keystore secret")
	}

	return createdSecret, nil
}
