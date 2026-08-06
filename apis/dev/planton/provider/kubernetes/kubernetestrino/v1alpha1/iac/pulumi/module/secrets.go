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

// trinoSecrets materializes the module-owned Secrets the release depends
// on. Every properties surface in this chart renders into ConfigMaps, so
// NOTHING credential-bearing may ride values — these Secrets are the
// only place secret material lives:
//
//   - `<name>-auth` (module-generated admin arm only): key `password`
//     holds the plaintext the admin (and verifiers, and BI tools) use;
//     key `password.db` holds the htpasswd-format file
//     (`<user>:<bcrypt>`) the chart mounts through auth.passwordAuthSecret
//     and the PASSWORD authenticator reads. Both keys derive from ONE
//     random — a verified pairing by construction. KNOW THIS (import
//     semantics): the bcrypt hash re-salts when the random is imported —
//     the import map carries an import_normalized tolerance for the
//     `password.db` key (the Harbor htpasswd class).
//   - `<name>-internal` (auth on): the internal-communication shared
//     secret every node presents — reaches config.properties as
//     `${ENV:TRINO_INTERNAL_SHARED_SECRET}`, never rendered.
//
// Returns the created resources (release dependencies — the pods mount
// them at startup).
func trinoSecrets(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) ([]pulumi.Resource, error) {
	var createdResources []pulumi.Resource

	if !locals.AuthEnabled {
		return nil, nil
	}

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

	// -------------------- admin credential + password file ----------------
	if locals.ModuleOwnedPasswordDb {
		// Letters+digits only: users type this at BI-tool connection
		// forms and the trino CLI — symbol-free avoids quoting bugs;
		// the larger length compensates the smaller alphabet.
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
		// BcryptHash derives from the SAME random — the plaintext and
		// the file entry can never drift apart (the verified-pairing
		// discipline; upstream requires bcrypt in password.db).
		passwordDb := pulumi.All(adminPassword.BcryptHash).ApplyT(
			func(args []interface{}) string {
				return locals.AdminUsername + ":" + args[0].(string) + "\n"
			}).(pulumi.StringOutput)

		created, err := newSecret("auth-secret", locals.PasswordDbSecretName, pulumi.StringMap{
			vars.AdminPasswordKey: adminPassword.Result,
			vars.PasswordDbKey:    passwordDb,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create admin auth secret")
		}
		createdResources = append(createdResources, created)
	}

	// -------------------- internal-communication secret -------------------
	// Length 64, letters+digits: the value transits an env var into
	// Trino's own config parser — symbol-free sidesteps any quoting
	// ambiguity (the NATS password-parser lesson, applied
	// preemptively).
	sharedSecret, err := random.NewRandomPassword(ctx, "internal-shared-secret",
		&random.RandomPasswordArgs{
			Length:  pulumi.Int(64),
			Special: pulumi.Bool(false),
		},
		generationShapeIgnores)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate internal shared secret")
	}
	created, err := newSecret("internal-secret", locals.InternalSecretName, pulumi.StringMap{
		vars.SharedSecretKey: sharedSecret.Result,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create internal shared-secret secret")
	}
	createdResources = append(createdResources, created)

	return createdResources, nil
}
