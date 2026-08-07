# KubernetesKarapace Pulumi Module

Deploys a Karapace schema registry as module-owned native manifests —
upstream ships no Helm chart or operator, so the module IS the
deployment: Deployments and Services per role, `KARAPACE_*`
environment-variable configuration on upstream's own image, and
Secret-mounted file material. Every rendering has an exact twin in
the Terraform module's `main.tf`/`locals.tf`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true
2. **SASL password Secret** (`<name>-sasl`, key `password`) — ONLY
   when the spec declares a literal `kafka.sasl.password`; a
   referenced `password_secret` is used in place and nothing is
   created
3. **Registry Deployment + Service** (`<name>`) — the schema-registry
   role (`python3 -m karapace`,
   `KARAPACE_KARAPACE_REGISTRY=true` / `KARAPACE_KARAPACE_REST=false`)
4. **REST-proxy Deployment + Service** (`<name>-rest`, optional,
   `rest_proxy.enabled`) — the same image with the role flags flipped
   (`python3 -m karapace.kafka_rest_apis`), wired to the registry
   Service and the same Kafka cluster

## Key Rendering Contracts

- **Image pin**: `ghcr.io/aiven-open/karapace:6.2.1` when
  `spec.image` is empty (`vars.KarapaceImage`; bump in lockstep with
  the Terraform module).
- **Per-pod advertised identity**: `KARAPACE_ADVERTISED_HOSTNAME`
  injects each pod's IP via the downward API (`status.podIP`) — the
  leader publishes `advertised_protocol://advertised_hostname:port`
  and followers forward writes to it, so a shared Service name would
  make followers forward to themselves; a Deployment pod's bare name
  does not resolve in cluster DNS. Must be explicit: the engine's
  fallback is `host`, which the module sets to `0.0.0.0`.
- **Secret-mounted file material at fixed paths**: Kafka CA at
  `/etc/karapace/kafka-ca`, mutual-TLS client identity at
  `/etc/karapace/kafka-cert`, server TLS at
  `/etc/karapace/server-tls`, the basic-auth authfile at
  `/etc/karapace/auth` — env vars point inside (paths byte-identical
  with the Terraform module).
- **The SASL password always arrives via secretKeyRef** — either the
  referenced Secret or the module-materialized one; never a plaintext
  env value (pod specs are world-readable to get-pod RBAC, Secret
  values have their own ACL).
- **server_tls flips the coupled trio**:
  `KARAPACE_ADVERTISED_PROTOCOL=https`, the cert/key file paths, and
  the probe scheme (HTTPS probes skip certificate verification —
  cert-manager Service-SAN certificates work). The REST proxy always
  serves plain HTTP.
- **Probes hit `/_health`** — the engine's unauthenticated health
  path (skip-auth list), with an initial delay and generous failure
  threshold so schemas-topic replay at startup is not punished.
- **Per-role selector identity** — both Deployments run the same
  image in one namespace; the role-specific `app` label
  (`<name>` / `<name>-rest`) keeps each Service from selecting the
  other role's pods.
- **Scheduling is registry-scoped** — `node_selector`/`tolerations`
  apply to the registry pods per the spec contract; the REST-proxy
  role carries only replicas/port/resources.
- **Registry behavior env** — topic name, at-creation replication
  factor, compatibility, leader-election group id (default
  `metadata.name`), election strategy — resolved in `locals.go` with
  the spec defaults applied.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the registry runs in |
| `service_name` | Registry Service (`<name>`) |
| `endpoint` | In-cluster endpoint (`http(s)://<name>.<namespace>.svc.cluster.local:<port>`) — the `schema.registry.url` value |
| `rest_proxy_endpoint` | REST-proxy endpoint (`<name>-rest`); empty when the role is off |
| `schemas_topic` | The Kafka topic storing the schemas |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → SASL Secret → registry Deployment +
  Service → REST-proxy Deployment + Service → output exports
- `module/deployment.go`: both Deployments — role flags, common env
  (advertised identity, Kafka connection, SASL secretKeyRef), TLS
  volume/mount rendering, the authfile/OIDC arms, probes
- `module/service.go`: one Service per role
- `module/secret.go`: the `<name>-sasl` materialization (literal
  passwords only)
- `module/locals.go`: naming, defaults, scheme/endpoint resolution,
  SASL source resolution — kept in lockstep with `locals.tf`
- `module/vars.go`: the image pin, entrypoints, mount paths, the
  health path
