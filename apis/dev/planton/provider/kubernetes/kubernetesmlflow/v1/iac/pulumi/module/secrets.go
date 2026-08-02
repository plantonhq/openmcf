package module

import (
	"encoding/base64"
	"fmt"
	"net/url"

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
// namespace — the referenced database credential the module composes into
// the backend-store URI. The referenced Secret is created by another
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
			"failed to read the referenced credential Secret %s/%s — it must exist in MLflow's own namespace (this module can only consume Secrets from the namespace it installs into; co-locate MLflow with its database or replicate the Secret)",
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

// mlflowSecrets materializes the module-owned Secrets the server mounts
// or reads at startup. NOTHING credential-bearing lands in any rendered
// pod spec — the backend URI reaches the server as an env var from a
// Secret, the admin password lives only inside Secrets, and the
// basic-auth ini is Secret-mounted:
//
//   - `<name>-backend-uri` (key `uri`, database arms only): the full
//     SQLAlchemy URI (`postgresql://user:pass@host:port/db`) composed AT
//     APPLY TIME from the referenced credential Secret — the exported
//     handle other tools can mount to reach the same tracking database.
//   - `<name>-admin-auth` (key `password`, generated arm only): the
//     bootstrap admin password — the exported credential handle.
//     Upstream's admin/password1234 default never ships.
//   - `<name>-auth-config` (keys `basic_auth.ini` + `flask_secret_key`):
//     the basic-auth app's configuration — admin credential, default
//     permission, the auth database (the backend database on database
//     arms; a sqlite file beside the tracking data on the sqlite arm),
//     and the Flask CSRF signing key the auth app REFUSES to start
//     without (MLFLOW_FLASK_SERVER_SECRET_KEY, env-wired; verified live
//     at the pin). The server reads the ini through
//     MLFLOW_AUTH_CONFIG_PATH; the key rides env from the same Secret,
//     so every replica shares one consistent value — upstream's own
//     multi-server requirement.
//
// Returns the created resources plus the composed backend URI (the
// deployment wires it as an env var).
func mlflowSecrets(ctx *pulumi.Context,
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

	// --------------------- backend URI composition ------------------------
	// The composed URI doubles as the auth database URI on database
	// arms (the basic-auth app keeps its user/permission tables in the
	// same database — upstream-supported, and it keeps auth state as
	// durable as tracking state).
	var backendUri pulumi.StringOutput
	if locals.BackendType != "sqlite" {
		dbPassword, err := readSecretValue(ctx, kubernetesProvider, "db-password-read",
			locals.Namespace, locals.DbPasswordSecret, locals.DbPasswordSecretKey)
		if err != nil {
			return nil, err
		}
		encodedUser := url.QueryEscape(locals.DbUser)
		encodedPassword := dbPassword.ApplyT(func(v string) string {
			return url.QueryEscape(v)
		}).(pulumi.StringOutput)
		backendUri = pulumi.Sprintf("%s://%s:%s@%s:%d/%s",
			locals.DbProtocol, encodedUser, encodedPassword,
			locals.DbHost, locals.DbPort, locals.DbName)
		created, err := newSecret("backend-uri-secret", locals.BackendUriSecretName, pulumi.StringMap{
			vars.BackendUriKey: pulumi.ToSecret(backendUri).(pulumi.StringOutput),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create backend uri secret")
		}
		createdResources = append(createdResources, created)
	}

	if !locals.AuthEnabled {
		return createdResources, nil
	}

	// ----------------------------- admin password -------------------------
	var adminPassword pulumi.StringOutput
	if locals.AdminSecretModuleOwned {
		// Letters+digits only: the password lands inside an ini file
		// and in users' MLFLOW_TRACKING_PASSWORD env vars — symbol-free
		// avoids quoting bugs in both; the larger length compensates
		// the smaller alphabet.
		generated, err := random.NewRandomPassword(ctx, "admin-password",
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
		adminPassword = generated.Result
		created, err := newSecret("admin-auth-secret", locals.AdminSecretName, pulumi.StringMap{
			vars.AdminPasswordKey: adminPassword,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create admin auth secret")
		}
		createdResources = append(createdResources, created)
	} else {
		byoPassword, err := readSecretValue(ctx, kubernetesProvider, "admin-password-read",
			locals.Namespace, locals.AdminSecretName, locals.AdminSecretKey)
		if err != nil {
			return nil, err
		}
		adminPassword = byoPassword
	}

	// ------------------------------ auth config ----------------------------
	// The ini shape is the server's own contract
	// (mlflow/server/auth/basic_auth.ini at the pin): default
	// permission, auth database, admin bootstrap credential, and the
	// authorization function left at the server's default.
	var authDatabaseUri pulumi.StringInput
	if locals.BackendType == "sqlite" {
		authDatabaseUri = pulumi.String(locals.AuthDatabaseUriSqlite)
	} else {
		authDatabaseUri = backendUri
	}
	authIni := pulumi.Sprintf(`[mlflow]
default_permission = %s
database_uri = %s
admin_username = %s
admin_password = %s
authorization_function = mlflow.server.auth:authenticate_request_basic_auth
`, locals.DefaultPermission, authDatabaseUri, locals.AdminUsername, adminPassword)

	// The Flask CSRF signing key: the auth app raises at create_app
	// without MLFLOW_FLASK_SERVER_SECRET_KEY (server source at the
	// pin). Letters+digits only (env-consumed); module-generated —
	// there is no user-facing reason to model a knob for a purely
	// internal signing key.
	flaskSecretKey, err := random.NewRandomPassword(ctx, "flask-secret-key",
		&random.RandomPasswordArgs{
			Length:     pulumi.Int(40),
			Special:    pulumi.Bool(false),
			MinUpper:   pulumi.Int(2),
			MinLower:   pulumi.Int(2),
			MinNumeric: pulumi.Int(2),
		},
		generationShapeIgnores)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate flask secret key")
	}

	created, err := newSecret("auth-config-secret", locals.AuthConfigSecretName, pulumi.StringMap{
		vars.AuthConfigFileName: pulumi.ToSecret(authIni).(pulumi.StringOutput),
		vars.FlaskSecretKeyKey:  flaskSecretKey.Result,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create auth config secret")
	}
	createdResources = append(createdResources, created)

	return createdResources, nil
}
