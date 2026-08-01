package module

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

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

// newAlnumPassword generates a letters+digits password (no symbols —
// every consumer here embeds these into connection URLs and env values
// where shell/URL-structural characters invite quoting bugs; the larger
// length compensates the smaller alphabet).
func newAlnumPassword(ctx *pulumi.Context, name string, length int) (*random.RandomPassword, error) {
	return random.NewRandomPassword(ctx, name,
		&random.RandomPasswordArgs{
			Length:     pulumi.Int(length),
			Special:    pulumi.Bool(false),
			MinUpper:   pulumi.Int(2),
			MinLower:   pulumi.Int(2),
			MinNumeric: pulumi.Int(2),
		},
		generationShapeIgnores)
}

// readSecretValue reads one key from an EXISTING Secret in the install
// namespace — the referenced database/broker credential the module
// composes into Airflow's connection Secrets. The referenced Secret is
// created by another component (e.g. the KubernetesPostgres app Secret),
// so it exists before this program runs; the read is gated on DryRun so
// OFFLINE previews (no cluster) never dial the API server — during
// preview the composed values render as opaque secret placeholders,
// during apply the real value is read (a read-only GetSecret resource)
// and composed. Terraform twin: a kubernetes_secret_v1 data source
// deferred behind a module-created resource.
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
			"failed to read the referenced credential Secret %s/%s — it must exist in Airflow's own namespace (a secretKeyRef and this module can only consume Secrets from the namespace they live in; co-locate Airflow with its database or replicate the Secret)",
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

// urlEncode mirrors the chart's own `urlquery` treatment of userinfo
// segments so passwords with reserved characters survive URI composition.
func urlEncode(in pulumi.StringInput) pulumi.StringOutput {
	return pulumi.ToOutput(in).ApplyT(func(v interface{}) string {
		return url.QueryEscape(v.(string))
	}).(pulumi.StringOutput)
}

// airflowSecrets generates every module-owned credential and materializes
// the Secrets the chart's *SecretName contracts point at. NOTHING
// credential-bearing lands in rendered Helm values — only Secret NAMES do.
//
//   - `<name>-fernet-key` (key fernet-key): 32 random bytes, URL-safe
//     base64 — the exact shape Fernet requires. BYO respected.
//   - `<name>-api-secret-key` (key api-secret-key) and
//     `<name>-webserver-secret-key` (key webserver-secret-key): the API
//     session/signing keys. The chart's own path would REGENERATE these
//     on every upgrade render (randAlphaNum in the template), logging
//     out every session — module-generated values are stable.
//   - `<name>-jwt-secret` (key jwt-secret): Airflow 3's internal token
//     signing secret. Same render-regeneration class.
//   - `<name>-admin-auth` (key password): the bootstrap admin password
//     (generated arm only) — the exported credential handle; reaches the
//     create-user Job as an env var, never a rendered pod argument.
//   - `<name>-metadata-conn` / `<name>-result-backend-conn` (key
//     connection): SQLAlchemy URIs composed at apply time from the
//     referenced database credential (result-backend carries the
//     chart's `db+` scheme prefix and only exists on Celery arms).
//     With PgBouncer the URIs point at the pooler and the ini aliases,
//     mirroring the chart's own rewrite exactly. The metadata Secret
//     also carries `kedaConnection` (the direct-database form the
//     chart's KEDA autoscalers read) on mysql and pgbouncer arms.
//   - `<name>-redis-password` + `<name>-broker-url` (bundled arm): the
//     chart's own template draws a NEW random password on every render
//     behind a pre-install hook (upstream's admitted hack) — the module
//     owns both instead. The composed arm reuses `<name>-broker-url`
//     with the referenced Valkey credential.
//   - `<name>-log-read-conn` (key connection): the Elasticsearch/
//     OpenSearch read-path URI.
//   - `<name>-pgbouncer-config` (keys pgbouncer.ini + users.txt): the
//     pooler's config — the chart's own rendering path would embed the
//     database password in Helm values, so the module composes it.
//
// Returns the created resources (release dependencies — the chart reads
// them at install time).
func airflowSecrets(ctx *pulumi.Context,
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

	// ------------------------------ fernet key ----------------------------
	if locals.FernetKeyModuleOwned {
		// Fernet requires EXACTLY 32 random bytes in URL-SAFE base64;
		// random bytes give the full keyspace (a password-charset string
		// would not decode to 32 bytes). Terraform twin: random_bytes +
		// the same +/→-_ substitution.
		fernetBytes, err := random.NewRandomBytes(ctx, "fernet-key",
			&random.RandomBytesArgs{Length: pulumi.Int(32)},
			pulumi.IgnoreChanges([]string{"length"}))
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate fernet key bytes")
		}
		fernetKey := fernetBytes.Base64.ApplyT(func(b64 string) string {
			return base64StdToUrlSafe(b64)
		}).(pulumi.StringOutput)
		created, err := newSecret("fernet-key-secret", locals.FernetKeySecretName, pulumi.StringMap{
			vars.FernetKeyKey: pulumi.ToSecret(fernetKey).(pulumi.StringOutput),
			// The companion STANDARD-base64 form of the same bytes —
			// the import handle for the random_bytes resource, whose
			// import format is standard base64 while Fernet requires
			// the URL-safe alphabet. Never consumed by the chart.
			vars.FernetKeyStdB64Key: pulumi.ToSecret(fernetBytes.Base64).(pulumi.StringOutput),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create fernet key secret")
		}
		createdResources = append(createdResources, created)
	}

	// --------------------------- api/session keys -------------------------
	if locals.ApiSecretKeyModuleOwned {
		apiKey, err := newAlnumPassword(ctx, "api-secret-key", 32)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate api secret key")
		}
		created, err := newSecret("api-secret-key-secret", locals.ApiSecretKeySecretName, pulumi.StringMap{
			vars.ApiSecretKeyKey: apiKey.Result,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create api secret key secret")
		}
		createdResources = append(createdResources, created)
	}

	webserverKey, err := newAlnumPassword(ctx, "webserver-secret-key", 32)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate webserver secret key")
	}
	createdWebserverKey, err := newSecret("webserver-secret-key-secret", locals.WebserverSecretKeyName, pulumi.StringMap{
		vars.WebserverSecretKeyKey: webserverKey.Result,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create webserver secret key secret")
	}
	createdResources = append(createdResources, createdWebserverKey)

	if locals.JwtSecretModuleOwned {
		jwtKey, err := newAlnumPassword(ctx, "jwt-secret", 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate jwt secret")
		}
		created, err := newSecret("jwt-secret-secret", locals.JwtSecretName, pulumi.StringMap{
			vars.JwtSecretKey: jwtKey.Result,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create jwt secret")
		}
		createdResources = append(createdResources, created)
	}

	// ----------------------------- admin password -------------------------
	if locals.AdminCreate && locals.AdminSecretModuleOwned {
		adminPassword, err := newAlnumPassword(ctx, "admin-password", 24)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate admin password")
		}
		created, err := newSecret("admin-auth-secret", locals.AdminSecretName, pulumi.StringMap{
			vars.AdminPasswordKey: adminPassword.Result,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create admin auth secret")
		}
		createdResources = append(createdResources, created)
	}

	// -------------------- database connection composition -----------------
	dbPassword, err := readSecretValue(ctx, kubernetesProvider, "db-password-read",
		locals.Namespace, locals.DbPasswordSecret, locals.DbPasswordSecretKey)
	if err != nil {
		return nil, err
	}
	encodedDbPassword := urlEncode(dbPassword)
	encodedDbUser := url.QueryEscape(locals.DbUser)

	// sslmode rides the query string only on the postgresql protocol —
	// the chart's own urlJoin does the same.
	query := ""
	if locals.DbProtocol == "postgresql" {
		query = fmt.Sprintf("?sslmode=%s", locals.DbSslMode)
	}

	metadataUri := pulumi.Sprintf("%s://%s:%s@%s:%d/%s%s",
		locals.DbProtocol, encodedDbUser, encodedDbPassword,
		locals.EffectiveDbHost, locals.EffectiveDbPort, locals.EffectiveMetadataDbName, query)
	metadataConnData := pulumi.StringMap{
		vars.ConnectionKey: pulumi.ToSecret(metadataUri).(pulumi.StringOutput),
	}
	// The chart's KEDA autoscalers read env KEDA_DB_CONN from THIS
	// Secret's `kedaConnection` key whenever the trigger cannot ride
	// the normal `connection` URI: always on mysql (KEDA's mysql scaler
	// wants the Go-DSN `user:pass@tcp(host:port)/db` form, never a URI)
	// and on the postgres pgbouncer-BYPASS posture (the scaler dials
	// the real database while Airflow rides the pooler). The chart
	// gates the key on KEDA being enabled; the module renders it
	// whenever it COULD be needed — the extra key is inert otherwise
	// (nothing else reads this module-owned Secret, and it carries the
	// same credential material as `connection`), and escape-hatch KEDA
	// configs (triggerer.keda, workers.keda.usePgbouncer=false) then
	// work without a spec change. Terraform twin: the same merge on
	// keda_conn_needed.
	if locals.DbProtocol == "mysql" {
		kedaUri := pulumi.Sprintf("%s:%s@tcp(%s:%d)/%s",
			encodedDbUser, encodedDbPassword, locals.DbHost, locals.DbPort, locals.DbName)
		metadataConnData[vars.KedaConnectionKey] = pulumi.ToSecret(kedaUri).(pulumi.StringOutput)
	} else if locals.PgbouncerEnabled {
		// The DIRECT postgres URI (real host and database — bypassing
		// the pooler), mirroring the chart's own kedaConnection
		// rendering.
		kedaUri := pulumi.Sprintf("postgresql://%s:%s@%s:%d/%s%s",
			encodedDbUser, encodedDbPassword, locals.DbHost, locals.DbPort, locals.DbName, query)
		metadataConnData[vars.KedaConnectionKey] = pulumi.ToSecret(kedaUri).(pulumi.StringOutput)
	}
	createdMetadataConn, err := newSecret("metadata-conn-secret", locals.MetadataConnSecretName, metadataConnData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create metadata connection secret")
	}
	createdResources = append(createdResources, createdMetadataConn)

	if locals.CeleryEnabled {
		// Celery's result backend needs SQLAlchemy's `db+` scheme prefix
		// (the chart's result-backend template renders exactly this).
		resultBackendUri := pulumi.Sprintf("db+%s://%s:%s@%s:%d/%s%s",
			locals.DbProtocol, encodedDbUser, encodedDbPassword,
			locals.EffectiveDbHost, locals.EffectiveDbPort, locals.EffectiveResultBackendDb, query)
		created, err := newSecret("result-backend-conn-secret", locals.ResultBackendConnSecretName, pulumi.StringMap{
			vars.ConnectionKey: pulumi.ToSecret(resultBackendUri).(pulumi.StringOutput),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create result backend connection secret")
		}
		createdResources = append(createdResources, created)
	}

	// --------------------------- broker composition -----------------------
	if locals.BundledRedisEnabled {
		redisPassword, err := newAlnumPassword(ctx, "redis-password", 40)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate redis password")
		}
		createdRedisPassword, err := newSecret("redis-password-secret", locals.RedisPasswordSecretName, pulumi.StringMap{
			"password": redisPassword.Result,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create redis password secret")
		}
		createdResources = append(createdResources, createdRedisPassword)

		// redis://:<password>@<name>-redis:6379/0 — the chart's own URL
		// shape for the bundled broker (alphanumeric password: no
		// encoding surprises by construction).
		brokerUrl := pulumi.Sprintf("redis://:%s@%s:%d/%d",
			redisPassword.Result, locals.BrokerHost, locals.BrokerPort, locals.BrokerDb)
		createdBrokerUrl, err := newSecret("broker-url-secret", locals.BrokerUrlSecretName, pulumi.StringMap{
			vars.ConnectionKey: pulumi.ToSecret(brokerUrl).(pulumi.StringOutput),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create broker url secret")
		}
		createdResources = append(createdResources, createdBrokerUrl)
	} else if locals.BrokerUrlSecretModuleOwned {
		userinfo := pulumi.Sprintf("%s:", url.QueryEscape(locals.BrokerUser))
		if locals.BrokerPasswordSecret != "" {
			brokerPassword, err := readSecretValue(ctx, kubernetesProvider, "broker-password-read",
				locals.Namespace, locals.BrokerPasswordSecret, locals.BrokerPasswordSecretKey)
			if err != nil {
				return nil, err
			}
			userinfo = pulumi.Sprintf("%s:%s", url.QueryEscape(locals.BrokerUser), urlEncode(brokerPassword))
		}
		brokerUrl := pulumi.Sprintf("redis://%s@%s:%d/%d",
			userinfo, locals.BrokerHost, locals.BrokerPort, locals.BrokerDb)
		created, err := newSecret("broker-url-secret", locals.BrokerUrlSecretName, pulumi.StringMap{
			vars.ConnectionKey: pulumi.ToSecret(brokerUrl).(pulumi.StringOutput),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create broker url secret")
		}
		createdResources = append(createdResources, created)
	}

	// -------------------------- log read connection -----------------------
	if locals.LogBackend != "" {
		var logUri pulumi.StringOutput
		if locals.LogBackendUser != "" {
			logPassword, err := readSecretValue(ctx, kubernetesProvider, "log-backend-password-read",
				locals.Namespace, locals.LogBackendPasswordSecret, locals.LogBackendPasswordSecretKey)
			if err != nil {
				return nil, err
			}
			logUri = pulumi.Sprintf("%s://%s:%s@%s:%d",
				locals.LogBackendScheme, url.QueryEscape(locals.LogBackendUser), urlEncode(logPassword),
				locals.LogBackendHost, locals.LogBackendPort)
		} else {
			logUri = pulumi.Sprintf("%s://%s:%d",
				locals.LogBackendScheme, locals.LogBackendHost, locals.LogBackendPort)
		}
		created, err := newSecret("log-read-conn-secret", locals.LogReadConnSecretName, pulumi.StringMap{
			vars.ConnectionKey: pulumi.ToSecret(logUri).(pulumi.StringOutput),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create log read connection secret")
		}
		createdResources = append(createdResources, created)
	}

	// ---------------------------- pgbouncer config ------------------------
	if locals.PgbouncerEnabled {
		ini := pgbouncerIni(locals)
		usersTxt := pulumi.Sprintf("%q %q\n%q %q\n",
			locals.DbUser, dbPassword, locals.DbUser, dbPassword)
		created, err := newSecret("pgbouncer-config-secret", locals.PgbouncerConfigSecretName, pulumi.StringMap{
			"pgbouncer.ini": pulumi.String(ini),
			"users.txt":     pulumi.ToSecret(usersTxt).(pulumi.StringOutput),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create pgbouncer config secret")
		}
		createdResources = append(createdResources, created)
	}

	return createdResources, nil
}

// pgbouncerIni renders the pooler configuration BYTE-FAITHFUL to the
// chart's own pgbouncer_config helper at the pin: the [databases]
// aliases (`<release>-metadata`, `<release>-result-backend`) map to the
// REAL database host/name, and Airflow's connection URIs (composed
// above) dial the aliases through the pooler.
func pgbouncerIni(locals *Locals) string {
	pgb := locals.Spec.GetPgbouncer()
	metadataPoolSize := int32(10)
	if pgb.MetadataPoolSize != nil {
		metadataPoolSize = pgb.GetMetadataPoolSize()
	}
	resultBackendPoolSize := int32(5)
	if pgb.ResultBackendPoolSize != nil {
		resultBackendPoolSize = pgb.GetResultBackendPoolSize()
	}
	maxClientConn := int32(100)
	if pgb.MaxClientConnections != nil {
		maxClientConn = pgb.GetMaxClientConnections()
	}

	return fmt.Sprintf(`[databases]
%s-metadata = host=%s dbname=%s port=%d pool_size=%d
%s-result-backend = host=%s dbname=%s port=%d pool_size=%d

[pgbouncer]
pool_mode = transaction
listen_port = %d
listen_addr = *
auth_type = scram-sha-256
auth_file = /etc/pgbouncer/users.txt
stats_users = %s
ignore_startup_parameters = extra_float_digits
max_client_conn = %d
verbose = 0
log_disconnections = 0
log_connections = 0

server_tls_sslmode = prefer
server_tls_ciphers = normal
`,
		locals.ReleaseName, locals.DbHost, locals.DbName, locals.DbPort, metadataPoolSize,
		locals.ReleaseName, locals.DbHost, locals.DbName, locals.DbPort, resultBackendPoolSize,
		vars.PgbouncerPort,
		locals.DbUser,
		maxClientConn)
}

// base64StdToUrlSafe converts standard base64 to the URL-safe alphabet
// Fernet requires (+ → -, / → _; padding kept).
func base64StdToUrlSafe(in string) string {
	return strings.NewReplacer("+", "-", "/", "_").Replace(in)
}
