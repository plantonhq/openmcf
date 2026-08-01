package module

import (
	"encoding/base64"
	"fmt"

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

// readSecretValue reads one key from an EXISTING Secret in the install
// namespace — the referenced database credential the module copies into
// the hub's existing-secret. The referenced Secret is created by another
// component (e.g. the KubernetesPostgres app Secret), so it exists before
// this program runs; the read is gated on DryRun so OFFLINE previews (no
// cluster) never dial the API server — during preview the composed value
// renders as an opaque secret placeholder, during apply the real value is
// read (a read-only GetSecret resource). Terraform twin: a
// kubernetes_secret_v1 data source deferred behind a module-created
// resource.
func readSecretValue(ctx *pulumi.Context,
	kubernetesProvider pulumi.ProviderResource,
	readName, namespace, secretName, secretKey string,
) (pulumi.StringOutput, error) {
	if ctx.DryRun() {
		return pulumi.ToSecret(pulumi.String(fmt.Sprintf(
			"(known after apply: %s/%s key %s)", namespace, secretName, secretKey,
		))).(pulumi.StringOutput), nil
	}
	got, err := kubernetescorev1.GetSecret(ctx, readName,
		pulumi.ID(fmt.Sprintf("%s/%s", namespace, secretName)), nil,
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return pulumi.StringOutput{}, errors.Wrapf(err,
			"failed to read the referenced credential Secret %s/%s — it must exist in JupyterHub's own namespace (this module can only consume Secrets from the namespace it installs into; co-locate JupyterHub with its database or replicate the Secret)",
			namespace, secretName)
	}
	value := got.Data.ApplyT(func(data map[string]string) (string, error) {
		encoded, ok := data[secretKey]
		if !ok {
			return "", errors.Errorf(
				"the referenced credential Secret %s/%s has no key %q — set the password secret's secret_key to the key that actually holds the password",
				namespace, secretName, secretKey)
		}
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return "", errors.Wrapf(err, "failed to decode key %q in Secret %s/%s", secretKey, namespace, secretName)
		}
		return string(decoded), nil
	}).(pulumi.StringOutput)
	return pulumi.ToSecret(value).(pulumi.StringOutput), nil
}

// jupyterhubSecrets materializes the module-owned Secrets the release
// depends on. NOTHING credential-bearing lands in rendered Helm values —
// this matters doubly for this chart, whose hub Secret embeds the ENTIRE
// rendered values document readable to anyone who can read Secrets in the
// namespace:
//
//   - `<name>-hub-secret` (key `hub.db.password`, external database arms
//     only): the hub mounts it through the chart's hub.existingSecret
//     seam and exports PGPASSWORD/MYSQL_PWD from it at startup — the
//     password never rides hub.db.url or any rendered value. Composed at
//     apply time from the referenced credential Secret.
//   - `<name>-auth` (key `password`, shared-password arm with no BYO
//     Secret): the module-generated shared sign-in password — the
//     exported credential handle. Reaches the hub as an env var the
//     module's extraConfig snippet reads; never rendered into values.
//
// The chart's three internal auth materials (proxy token, cookie secret,
// crypt keys) stay CHART-MANAGED deliberately: the chart generates them
// lookup-stable (reuses the existing Secret's values on upgrade), and
// pushing module values through hub.config would land them readable in
// rendered values — strictly worse.
//
// Returns the created resources (release dependencies — the hub mounts
// them at startup).
func jupyterhubSecrets(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) ([]pulumi.Resource, error) {
	var createdResources []pulumi.Resource

	newSecret := func(resourceName, secretName string, data pulumi.StringMap) (pulumi.Resource, error) {
		created, err := kubernetescorev1.NewSecret(ctx, resourceName,
			&kubernetescorev1.SecretArgs{
				Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
					Name:      pulumi.String(secretName),
					Namespace: pulumi.String(locals.Namespace),
					Labels:    pulumi.ToStringMap(locals.Labels),
				}),
				Type:       pulumi.String("Opaque"),
				StringData: data,
			}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)...)
		if err != nil {
			return nil, err
		}
		return created, nil
	}

	// -------------------- hub existing-secret (database) ------------------
	if locals.HubSecretEnabled {
		dbPassword, err := readSecretValue(ctx, kubernetesProvider, "db-password-read",
			locals.Namespace, locals.DbPasswordSecret, locals.DbPasswordSecretKey)
		if err != nil {
			return nil, err
		}
		created, err := newSecret("hub-secret", locals.HubSecretName, pulumi.StringMap{
			vars.HubDbPasswordKey: dbPassword,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create hub existing-secret")
		}
		createdResources = append(createdResources, created)
	}

	// ------------------------ shared sign-in password ---------------------
	if locals.SharedPasswordModuleOwned {
		// Letters+digits only: the password reaches JupyterHub through
		// an env var and users type it at the login form — symbol-free
		// avoids both quoting bugs and login-typing friction; the
		// larger length compensates the smaller alphabet.
		sharedPassword, err := random.NewRandomPassword(ctx, "shared-password",
			&random.RandomPasswordArgs{
				Length:     pulumi.Int(24),
				Special:    pulumi.Bool(false),
				MinUpper:   pulumi.Int(2),
				MinLower:   pulumi.Int(2),
				MinNumeric: pulumi.Int(2),
			},
			generationShapeIgnores)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate shared sign-in password")
		}
		created, err := newSecret("auth-secret", locals.SharedPasswordSecretName, pulumi.StringMap{
			vars.SharedPasswordKey: sharedPassword.Result,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create shared sign-in password secret")
		}
		createdResources = append(createdResources, created)
	}

	return createdResources, nil
}
