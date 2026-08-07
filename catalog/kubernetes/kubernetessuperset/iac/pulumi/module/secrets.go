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

// supersetSecrets materializes the module-owned Secrets the release
// depends on. The chart consumes ALL runtime credentials through ONE
// environment Secret (every component envFroms it); the chart's own copy
// is turned OFF (secretEnv.create=false) and the module composes
// `<name>-env` itself:
//
//   - NON-SECRET connection facts (DB_HOST/PORT/USER/NAME, REDIS_*,
//     the admin identity) render as plain keys.
//   - Module-GENERATED material (SUPERSET_SECRET_KEY, ADMIN_PASSWORD,
//     the websocket JWT) renders into this Secret from the randoms —
//     plus dedicated handle Secrets (`<name>-secret-key`,
//     `<name>-admin-auth`) so compositions and operators have stable,
//     single-purpose references.
//   - REFERENCED material (the database/cache passwords, bring-your-own
//     keys) is NEVER copied here — it arrives in the pods as
//     extraEnvRaw secretKeyRef entries (values.go), the chart's own
//     mechanism; explicit env beats envFrom.
//
// SECRET_KEY STABILITY (load-bearing): the Flask session-signing key
// also encrypts datasource credentials stored in the metadata database —
// the random is generated ONCE and shape-ignored; rotating it without
// Superset's own re-encrypt procedure orphans every stored connection.
//
// Returns the created resources (release dependencies — the pods envFrom
// the env Secret at startup).
func supersetSecrets(ctx *pulumi.Context,
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

	envData := pulumi.StringMap{}
	for key, value := range locals.EnvPlain {
		envData[key] = pulumi.String(value)
	}
	if ssl := locals.Spec.GetMetadataDatabase().GetSsl(); ssl.GetEnabled() {
		envData["DB_SSL_MODE"] = pulumi.String(defaultString(ssl.GetMode(), "require"))
	}

	// --------------------------- SECRET_KEY -------------------------------
	if locals.SecretKeyModuleOwned {
		secretKey, err := random.NewRandomPassword(ctx, "secret-key",
			&random.RandomPasswordArgs{
				Length:  pulumi.Int(42),
				Special: pulumi.Bool(false),
			},
			generationShapeIgnores)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate SECRET_KEY")
		}
		envData["SUPERSET_SECRET_KEY"] = secretKey.Result
		created, err := newSecret("secret-key-secret", locals.SecretKeySecretName, pulumi.StringMap{
			vars.SecretKeyKey: secretKey.Result,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create secret-key secret")
		}
		createdResources = append(createdResources, created)
	}

	// -------------------------- admin password ----------------------------
	if locals.AdminModuleOwned {
		// Letters+digits only: operators type this at the login form —
		// symbol-free avoids quoting bugs and typing friction; the
		// larger length compensates the smaller alphabet.
		adminPassword, err := random.NewRandomPassword(ctx, "admin-password",
			&random.RandomPasswordArgs{
				Length:     pulumi.Int(24),
				Special:    pulumi.Bool(false),
				MinUpper:   pulumi.Int(2),
				MinLower:   pulumi.Int(2),
				MinNumeric: pulumi.Int(2),
			},
			generationShapeIgnores)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate admin password")
		}
		envData["ADMIN_PASSWORD"] = adminPassword.Result
		created, err := newSecret("admin-auth-secret", locals.AdminPasswordSecret, pulumi.StringMap{
			vars.AdminPasswordKey: adminPassword.Result,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create admin auth secret")
		}
		createdResources = append(createdResources, created)
	}

	// -------------------------- websocket JWT -----------------------------
	// One key serves both sides: the ws server reads JWT_SECRET from
	// its environment natively (env beats config.json), and the
	// module's configOverrides snippet points Superset's
	// GLOBAL_ASYNC_QUERIES_JWT_SECRET at the same variable.
	if locals.WebsocketsEnabled {
		jwtSecret, err := random.NewRandomPassword(ctx, "ws-jwt-secret",
			&random.RandomPasswordArgs{
				Length:  pulumi.Int(48),
				Special: pulumi.Bool(false),
			},
			generationShapeIgnores)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate websocket JWT secret")
		}
		envData[vars.JwtSecretEnvVar] = jwtSecret.Result
	}

	// ------------------------- the env Secret -----------------------------
	created, err := newSecret("env-secret", locals.EnvSecretName, envData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create environment secret")
	}
	createdResources = append(createdResources, created)

	return createdResources, nil
}
