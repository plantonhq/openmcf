# KubernetesKafkaUi Pulumi Module

Deploys the kafbat UI console from the served `kafka-ui` Helm chart
(https://ui.charts.kafbat.io, default pin 1.6.4): the optional
namespace, the optional module-materialized console Secret, and the
Helm release with the typed spec rendered into chart values and the
`helm_values` escape hatch merged LAST. Every rendering has an exact
twin in the Terraform module's `locals.tf`/`main.tf`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true
2. **Console Secret** (`<name>-secrets`, key `console-user-password`)
   — ONLY when `auth` is declared; holds the one literal credential
   the spec carries (the console login password)
3. **The Helm release** — release name = `metadata.name`, chart
   fullname pinned to it (`fullnameOverride`), so the Service IS the
   resource name and several consoles coexist per cluster

## The Placeholder / secretMappings Mechanism

The chart writes the app config verbatim into a ConfigMap, so the
module NEVER renders a credential value into it. Every password
position carries a `${ENV_VAR}` placeholder, and the chart's
`envs.secretMappings` wires each env var to a Secret key (rendered as
secretKeyRef). Spring resolves the placeholders inside the running
container. Referenced credentials (cluster SASL, schema-registry and
Connect passwords) map straight at their source Secrets; the console
login password maps at the module-materialized Secret. Env-var names
are deterministic and index-based (`KAFKA_CLUSTER_<i>_PASSWORD`,
`KAFKA_CLUSTER_<i>_SCHEMA_REGISTRY_PASSWORD`,
`KAFKA_CLUSTER_<i>_CONNECT_<j>_PASSWORD`, `KAFKA_UI_USER_PASSWORD`)
so both engines emit identical placeholders.

## Rendering Notes

- **`yamlApplicationConfig`** carries the app's cluster list
  (bootstrap, readOnly only when true, properties, schemaRegistry +
  auth, kafkaConnect entries) plus the auth tree: `LOGIN_FORM` with
  Spring's single default security user when `auth` is declared
  (the app registers no user store of its own — the reason the spec
  models ONE user), `DISABLED` explicitly otherwise.
- **PEM truststores, mounted as-is** — one secret volume per TLS
  cluster at `/etc/kafkaui/cluster-<i>-ca`; the rendered client
  properties set `ssl.truststore.type=PEM` and point at the mounted
  key, so Strimzi cluster-CA Secrets work directly. The module owns
  the derived `security.protocol` and the JAAS line
  (ScramLoginModule / PlainLoginModule).
- **Always-rendered service/replica values** — resolved defaults
  (ClusterIP/80, 1 replica) so both engines emit identical documents
  whether or not the platform's defaulting middleware ran.
- **`helm_values` merges LAST** with Helm `-f` semantics (maps
  deep-merge with the override winning, lists replace) — the exact
  semantic twin of the Terraform provider's native two-document
  values list.
- **Wait-for-ready, atomic** — SkipAwait false + Atomic +
  CleanupOnFail, 600s timeout: a console that never starts fails the
  deploy, not the first browser hit.

## Usage

```shell
planton pulumi up --manifest hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the console runs in |
| `service_name` | Console Service (`<name>`, pinned via fullnameOverride) |
| `endpoint` | In-cluster endpoint (`http://<name>.<namespace>.svc.cluster.local:<port>`) |
| `port_forward_command` | Workstation access without any exposure (local side pinned to 8080 — the Service port is often 80, unprivileged locally) |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → console Secret → Helm release →
  output exports
- `module/helm_release.go`: the chart values (app config,
  placeholders, secretMappings, TLS volumes/mounts), the helm_values
  merge, the release
- `module/secret.go`: the `<name>-secrets` materialization
- `module/locals.go`: naming, chart-version/default resolution,
  endpoint and port-forward outputs — kept in lockstep with
  `locals.tf`
- `module/vars.go`: chart name/repo and the 1.6.4 default pin
  (verified against the served repository index)
