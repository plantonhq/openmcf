package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// clickHouseAuthSecret generates the bundled ClickHouse admin password and
// materializes it as the "<metadata.name>-clickhouse-auth" Opaque Secret
// (keys username/password, username "admin" — the chart's bundled-arm
// user). The chart itself offers no Secret for the bundled arm (it renders
// the password into the installation object and container env — the
// upstream grain); this module-owned Secret is what downstream kinds and
// operators reference. The chart's publicly-documented default password
// NEVER ships. Created only on the bundled arm; returns the password
// output for the values injection. Terraform twin: random_password +
// kubernetes_secret_v1.clickhouse_auth with the same name and keys.
func clickHouseAuthSecret(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) (pulumi.StringOutput, pulumi.Resource, error) {
	// The generation-shape arguments are ignored after creation so an
	// IMPORTED credential never silently regenerates: random_password's
	// import carries only the VALUE, the importer assumes the provider's
	// own generation defaults (special=true — verified live), and every
	// argument forces replacement — without this, the first update after
	// an import proposes replacing (rotating) the live credential.
	// Rotation stays an explicit verb, never plan fallout.
	createdPassword, err := random.NewRandomPassword(ctx, "clickhouse-auth-password",
		&random.RandomPasswordArgs{
			Length:  pulumi.Int(24),
			Special: pulumi.Bool(false),
		},
		pulumi.IgnoreChanges([]string{
			"length", "special", "upper", "lower", "numeric",
			"minLower", "minNumeric", "minSpecial", "minUpper", "overrideSpecial",
		}))
	if err != nil {
		return pulumi.StringOutput{}, nil, errors.Wrap(err, "failed to generate clickhouse admin password")
	}

	opts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)

	createdSecret, err := kubernetescorev1.NewSecret(ctx, locals.ClickHouseAuthSecretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.ClickHouseAuthSecretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			Type: pulumi.String("Opaque"),
			StringData: pulumi.StringMap{
				"username": pulumi.String(vars.BundledClickHouseUser),
				"password": createdPassword.Result,
			},
		}, opts...)
	if err != nil {
		return pulumi.StringOutput{}, nil, errors.Wrap(err, "failed to create clickhouse auth secret")
	}

	return createdPassword.Result, createdSecret, nil
}
