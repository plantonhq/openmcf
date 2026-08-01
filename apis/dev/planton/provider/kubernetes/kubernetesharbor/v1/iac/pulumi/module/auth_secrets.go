package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// generatedCredentials carries the module-generated secret material that
// must reach the chart OUTSIDE the module-owned Secrets: the internal
// database superuser password is the one credential the chart accepts
// only as a value (`database.internal.password` — no existingSecret
// site exists for it at chart 1.19.1). It rides the release values as a
// SECRET output (redacted in previews); everything else travels by
// Secret name only.
type generatedCredentials struct {
	// nil unless the internal database arm is active.
	InternalDatabasePassword pulumi.StringInput
}

// The generation-shape arguments are ignored after creation so an
// IMPORTED credential never silently regenerates: rotation stays an
// explicit verb, never plan fallout. Terraform twin: random_password
// lifecycle ignore_changes on the same arguments.
var generationShapeIgnores = pulumi.IgnoreChanges([]string{
	"length", "special", "upper", "lower", "numeric",
	"minLower", "minNumeric", "minSpecial", "minUpper", "overrideSpecial",
})

// newAlnumPassword generates a letters+digits password (no symbols —
// several consumers embed these into URLs, htpasswd lines, and env
// values where shell/URL-structural characters invite quoting bugs;
// the larger length compensates the smaller alphabet).
func newAlnumPassword(ctx *pulumi.Context, name string, length int) (*random.RandomPassword, error) {
	return random.NewRandomPassword(ctx, name,
		&random.RandomPasswordArgs{
			Length:  pulumi.Int(length),
			Special: pulumi.Bool(false),
			// Harbor enforces password complexity on the admin login
			// (upper + lower + digit); pinning minimums keeps every
			// generated credential valid by construction.
			MinUpper:   pulumi.Int(2),
			MinLower:   pulumi.Int(2),
			MinNumeric: pulumi.Int(2),
		},
		generationShapeIgnores)
}

// authSecrets generates every module-owned credential and materializes
// the Secrets the chart's existingSecret sites point at:
//
//   - `<name>-admin-auth` (generated arm only): HARBOR_ADMIN_PASSWORD —
//     the exported credential handle.
//   - `<name>-internal-auth`: every generated inter-component
//     credential, one Secret with each chart site's CONTRACT KEY —
//     `secretKey` (the 16-char encryption key), `secret` (core's
//     inter-component secret), CSRF_KEY, JOBSERVICE_SECRET,
//     REGISTRY_HTTP_SECRET, and the registry basic-auth pair
//     REGISTRY_PASSWD + REGISTRY_HTPASSWD. The htpasswd line uses
//     random_password's STABLE bcrypt hash — the chart's own `htpasswd`
//     template function re-salts on every render, which would rotate
//     the credential on every apply (the chart's values comment itself
//     recommends a pre-computed line for CD tools).
//   - `<name>-redis-auth` (declared external-redis password only):
//     REDIS_PASSWORD (+ REDIS_USERNAME when a username is declared).
//   - `<name>-storage-auth` (declared s3/gcs/azure credentials only):
//     the storage driver's contract keys.
//
// Returns the Secrets created (release dependencies — the chart reads
// them at install time) plus the values-borne generated credentials.
func authSecrets(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) ([]pulumi.Resource, *generatedCredentials, error) {
	var createdResources []pulumi.Resource
	gen := &generatedCredentials{}

	newSecret := func(name string, data pulumi.StringMap) (pulumi.Resource, error) {
		created, err := kubernetescorev1.NewSecret(ctx, name,
			&kubernetescorev1.SecretArgs{
				Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
					Name:      pulumi.String(name),
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

	// ---------------------------- admin password --------------------------
	if locals.AdminGenerated {
		adminPassword, err := newAlnumPassword(ctx, "admin-password", 32)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to generate admin password")
		}
		created, err := newSecret(locals.AdminSecretName, pulumi.StringMap{
			vars.AdminPasswordSecretKey: adminPassword.Result,
		})
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create admin auth secret")
		}
		createdResources = append(createdResources, created)
	}

	// ------------------------ inter-component secrets ---------------------
	// Chart contract lengths: secretKey MUST be exactly 16 characters
	// (AES key), the CSRF key exactly 32; the component secrets follow
	// the chart's own randAlphaNum(16) shape.
	encryptionKey, err := newAlnumPassword(ctx, "encryption-key", 16)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate encryption key")
	}
	coreSecret, err := newAlnumPassword(ctx, "core-secret", 16)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate core secret")
	}
	csrfKey, err := newAlnumPassword(ctx, "csrf-key", 32)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate csrf key")
	}
	jobserviceSecret, err := newAlnumPassword(ctx, "jobservice-secret", 16)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate jobservice secret")
	}
	registryHttpSecret, err := newAlnumPassword(ctx, "registry-http-secret", 16)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate registry http secret")
	}
	registryPassword, err := newAlnumPassword(ctx, "registry-credential-password", 32)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate registry credential password")
	}

	createdInternalAuth, err := newSecret(locals.InternalAuthSecretName, pulumi.StringMap{
		"secretKey":            encryptionKey.Result,
		"secret":               coreSecret.Result,
		"CSRF_KEY":             csrfKey.Result,
		"JOBSERVICE_SECRET":    jobserviceSecret.Result,
		"REGISTRY_HTTP_SECRET": registryHttpSecret.Result,
		"REGISTRY_PASSWD":      registryPassword.Result,
		// The stable bcrypt hash (computed once, kept in state) —
		// never the chart's per-render htpasswd function.
		"REGISTRY_HTPASSWD": pulumi.Sprintf("%s:%s",
			vars.RegistryCredentialUsername, registryPassword.BcryptHash),
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create internal auth secret")
	}
	createdResources = append(createdResources, createdInternalAuth)

	// -------------------------- internal database -------------------------
	if locals.InternalDatabase {
		databasePassword, err := newAlnumPassword(ctx, "internal-database-password", 24)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to generate internal database password")
		}
		gen.InternalDatabasePassword = databasePassword.Result
	}

	// ------------------------ declared redis password ---------------------
	if locals.RedisAuthSecretName != "" {
		ext := locals.Spec.GetCache().GetExternal()
		data := pulumi.StringMap{
			"REDIS_PASSWORD": pulumi.ToSecret(pulumi.String(ext.GetPassword())).(pulumi.StringOutput),
		}
		if ext.GetUsername() != "" {
			data["REDIS_USERNAME"] = pulumi.String(ext.GetUsername()).ToStringOutput()
		}
		created, err := newSecret(locals.RedisAuthSecretName, data)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create redis auth secret")
		}
		createdResources = append(createdResources, created)
	}

	// ---------------------- declared storage credentials ------------------
	if locals.StorageAuthSecretName != "" {
		data := pulumi.StringMap{}
		if s3 := locals.Spec.GetStorage().GetS3(); s3 != nil && s3.GetCredentials().GetAccessKey() != "" {
			data["REGISTRY_STORAGE_S3_ACCESSKEY"] = pulumi.String(s3.GetCredentials().GetAccessKey()).ToStringOutput()
			data["REGISTRY_STORAGE_S3_SECRETKEY"] = pulumi.ToSecret(pulumi.String(s3.GetCredentials().GetSecretKey())).(pulumi.StringOutput)
		}
		if gcs := locals.Spec.GetStorage().GetGcs(); gcs != nil && gcs.GetKeyData() != "" {
			data["GCS_KEY_DATA"] = pulumi.ToSecret(pulumi.String(gcs.GetKeyData())).(pulumi.StringOutput)
		}
		if azure := locals.Spec.GetStorage().GetAzure(); azure != nil && azure.GetAccountKey() != "" {
			data["AZURE_STORAGE_ACCESS_KEY"] = pulumi.ToSecret(pulumi.String(azure.GetAccountKey())).(pulumi.StringOutput)
		}
		created, err := newSecret(locals.StorageAuthSecretName, data)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create storage auth secret")
		}
		createdResources = append(createdResources, created)
	}

	return createdResources, gen, nil
}
