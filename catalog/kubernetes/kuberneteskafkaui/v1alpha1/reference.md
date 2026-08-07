# KubernetesKafkaUi

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**KubernetesKafkaUiSpec** declares a kafbat UI installation — the
Apache-2.0 web console for Kafka. One installation observes and
manages MANY clusters: browse topics and live messages, inspect
consumer groups and lag, view/register schemas through a connected
schema registry, and monitor Connect pipes — the console teams
coming off Confluent expect.

Deployed from the served `kafka-ui` Helm chart. Each `clusters`
entry wires one Kafka cluster (plus its optional schema registry
and Connect clusters) into the console — the FKs compose directly
with KubernetesKafka, KubernetesKarapace and
KubernetesKafkaConnect siblings.

EXPOSURE composes from first-class kinds (KubernetesIngress /
Gateway API) against the exported service handles — never embedded
here. Anyone who can reach the Service can act with the console's
permissions: enable `auth` (or keep the Service internal) before
any shared exposure, and set `read_only` on clusters the console
should observe but never mutate.

## Example

```yaml
# Full-surface offline-proof manifest: exercises two clusters (one plain,
# one production-postured with read_only, TLS trust from a Strimzi-style
# cluster-CA Secret, SCRAM-SHA-512 sasl with a referenced password Secret,
# a schema registry with HTTP Basic auth, and two Connect clusters — one
# authenticated, one open), login_form authentication (ONE user: the app
# supports a single LOGIN_FORM account via Spring Boot's default security
# user — see the module comments), replicas, resources, a NodePort
# Service, scheduling, the image-registry override, and the helm_values
# escape hatch (its replicaCount override must WIN over the typed replicas
# in the rendered values — the merge-last proof). Placeholder values;
# never applied to a real cluster.
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesKafkaUi
metadata:
  name: hack-kafka-ui
spec:
  namespace:
    value: hack-kafka-ui
  createNamespace: true
  chartVersion: 1.6.4
  clusters:
    - name: staging
      bootstrapServers:
        value: staging-kafka-bootstrap.kafka-staging.svc.cluster.local:9092
      properties:
        request.timeout.ms: "30000"
    - name: production
      bootstrapServers:
        value: prod-kafka-bootstrap.kafka-prod.svc.cluster.local:9093
      readOnly: true
      tls:
        caSecretName:
          value: prod-kafka-cluster-ca-cert
        caCertificate: ca.crt
      sasl:
        mechanism: SCRAM-SHA-512
        username: kafka-ui
        passwordSecret:
          secretName:
            value: prod-kafka-ui-user
          key: password
      schemaRegistry:
        url:
          value: http://prod-karapace.kafka-prod.svc.cluster.local:8081
        username: registry-ui
        passwordSecret:
          secretName:
            value: prod-karapace-ui-credentials
          key: password
      kafkaConnect:
        - name: cdc
          address:
            value: http://prod-connect-cdc.kafka-prod.svc.cluster.local:8083
          username: connect-ui
          passwordSecret:
            secretName:
              value: prod-connect-cdc-ui-credentials
            key: password
        - name: sink
          address:
            value: http://prod-connect-sink.kafka-prod.svc.cluster.local:8083
      properties:
        request.timeout.ms: "45000"
  auth:
    type: login_form
    user:
      username: admin
      password: hack-placeholder-admin
  replicas: 2
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      cpu: 500m
      memory: 512Mi
  serviceType: NodePort
  servicePort: 8080
  nodeSelector:
    kubernetes.io/os: linux
  tolerations:
    - key: dedicated
      operator: Equal
      value: consoles
      effect: NoSchedule
  imageRegistry: registry.example.com
  helmValues: |
    replicaCount: 3
    podLabels:
      example.org/owner: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.chartVersion` | `string` |  |  |  |
| `spec.clusters` | `[]KubernetesKafkaUiCluster` | yes |  |  |
| `spec.clusters[].name` | `string` | yes |  |  |
| `spec.clusters[].bootstrapServers` | `string \| valueFrom` | yes |  | KubernetesKafka (`status.outputs.internal_bootstrap_endpoint`) |
| `spec.clusters[].readOnly` | `bool` |  |  |  |
| `spec.clusters[].tls` | `KubernetesKafkaUiClusterTls` |  |  |  |
| `spec.clusters[].tls.caSecretName` | `string \| valueFrom` | yes |  | KubernetesKafka (`status.outputs.cluster_ca_cert_secret_name`) |
| `spec.clusters[].tls.caCertificate` | `string` |  | `ca.crt` |  |
| `spec.clusters[].sasl` | `KubernetesKafkaUiClusterSasl` |  |  |  |
| `spec.clusters[].sasl.mechanism` | `string` | yes |  |  |
| `spec.clusters[].sasl.username` | `string` | yes |  |  |
| `spec.clusters[].sasl.passwordSecret` | `KubernetesKafkaUiPasswordSecret` | yes |  |  |
| `spec.clusters[].sasl.passwordSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesKafkaUser (`status.outputs.secret_name`) |
| `spec.clusters[].sasl.passwordSecret.key` | `string` |  | `password` |  |
| `spec.clusters[].schemaRegistry` | `KubernetesKafkaUiSchemaRegistry` |  |  |  |
| `spec.clusters[].schemaRegistry.url` | `string \| valueFrom` | yes |  | KubernetesKarapace (`status.outputs.endpoint`) |
| `spec.clusters[].schemaRegistry.username` | `string` |  |  |  |
| `spec.clusters[].schemaRegistry.passwordSecret` | `KubernetesKafkaUiBasicAuthPasswordSecret` |  |  |  |
| `spec.clusters[].schemaRegistry.passwordSecret.secretName` | `string \| valueFrom` | yes |  |  |
| `spec.clusters[].schemaRegistry.passwordSecret.key` | `string` |  | `password` |  |
| `spec.clusters[].kafkaConnect` | `[]KubernetesKafkaUiConnectCluster` |  |  |  |
| `spec.clusters[].kafkaConnect[].name` | `string` | yes |  |  |
| `spec.clusters[].kafkaConnect[].address` | `string \| valueFrom` | yes |  | KubernetesKafkaConnect (`status.outputs.rest_api_endpoint`) |
| `spec.clusters[].kafkaConnect[].username` | `string` |  |  |  |
| `spec.clusters[].kafkaConnect[].passwordSecret` | `KubernetesKafkaUiBasicAuthPasswordSecret` |  |  |  |
| `spec.clusters[].kafkaConnect[].passwordSecret.secretName` | `string \| valueFrom` | yes |  |  |
| `spec.clusters[].kafkaConnect[].passwordSecret.key` | `string` |  | `password` |  |
| `spec.clusters[].properties` | `map<string, string>` |  |  |  |
| `spec.auth` | `KubernetesKafkaUiAuth` |  |  |  |
| `spec.auth.type` | `string` | yes |  |  |
| `spec.auth.user` | `KubernetesKafkaUiUser` | yes |  |  |
| `spec.auth.user.username` | `string` | yes |  |  |
| `spec.auth.user.password` | `string` (sensitive) | yes |  |  |
| `spec.replicas` | `int32` |  | `1` |  |
| `spec.resources` | `ContainerResources` |  |  |  |
| `spec.resources.limits` | `CpuMemory` |  |  |  |
| `spec.resources.limits.cpu` | `string` |  |  |  |
| `spec.resources.limits.memory` | `string` |  |  |  |
| `spec.resources.requests` | `CpuMemory` |  |  |  |
| `spec.resources.requests.cpu` | `string` |  |  |  |
| `spec.resources.requests.memory` | `string` |  |  |  |
| `spec.serviceType` | `string` |  | `ClusterIP` |  |
| `spec.servicePort` | `int32` |  | `80` |  |
| `spec.nodeSelector` | `map<string, string>` |  |  |  |
| `spec.tolerations` | `[]WorkloadToleration` |  |  |  |
| `spec.tolerations[].key` | `string` |  |  |  |
| `spec.tolerations[].operator` | `string` |  |  |  |
| `spec.tolerations[].value` | `string` |  |  |  |
| `spec.tolerations[].effect` | `string` |  |  |  |
| `spec.tolerations[].tolerationSeconds` | `int64` |  |  |  |
| `spec.imageRegistry` | `string` |  |  |  |
| `spec.helmValues` | `string` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace for the console. Accepts a literal namespace name or a
reference to a KubernetesNamespace resource.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.createNamespace

`bool`

When true, the namespace is created (with the standard Planton
governance labels) before installing and deleted with the
resource. When false, the namespace must already exist.

### spec.chartVersion

`string`

Helm chart version (the SERVED chart at
https://ui.charts.kafbat.io — pick versions from the repository
index, not the source tree's Chart.yaml). Empty = the module's
pinned default.

### spec.clusters

`[]KubernetesKafkaUiCluster` · required

The Kafka clusters this console connects to. Cluster names must
be unique — they are the console's display and API identifiers.

- rule: cluster names must be unique within the console
- rule: {"repeated":{"minItems":"1"}}

### spec.clusters[].name

`string` · required

Display/API name for this cluster inside the console (e.g.
"production", "staging").

- rule: cluster name may use alphanumerics, '.', '_' and '-'
- rule: {"required":true}

### spec.clusters[].bootstrapServers

`string | valueFrom` · required

Bootstrap address of the cluster, as host:port. Accepts a
literal address or a reference to a KubernetesKafka resource,
which resolves to its in-cluster bootstrap endpoint.

containment_exempt: the console OBSERVES this cluster; it is not
deployed inside it.

- references: KubernetesKafka (`status.outputs.internal_bootstrap_endpoint`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesKafka, name: <that resource's name>, fieldPath: status.outputs.internal_bootstrap_endpoint}} -- a bare string does not parse

### spec.clusters[].readOnly

`bool`

Observe-only: hide every mutating action (topic create/delete,
message produce, config edits) for this cluster. The right
posture for production clusters on a shared console.

### spec.clusters[].tls

`KubernetesKafkaUiClusterTls`

TLS trust for the cluster connection: a Secret holding the CA
certificate. For a Strimzi-managed cluster, reference the
KubernetesKafka resource to trust its cluster CA. Omitted =
plaintext.

### spec.clusters[].tls.caSecretName

`string | valueFrom` · required

Name of the Secret holding the CA certificate to trust. Accepts
a literal Secret name or a reference to a KubernetesKafka
resource, which resolves to that cluster's CA certificate
Secret. The module mounts it and renders the truststore
configuration.

containment_exempt: trust material fetched FROM the cluster —
access, never placement.

- references: KubernetesKafka (`status.outputs.cluster_ca_cert_secret_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesKafka, name: <that resource's name>, fieldPath: status.outputs.cluster_ca_cert_secret_name}} -- a bare string does not parse

### spec.clusters[].tls.caCertificate

`string` · optional (explicit presence)

The CA certificate's key within the Secret. Strimzi cluster CA
Secrets carry "ca.crt".

- default: `ca.crt`

### spec.clusters[].sasl

`KubernetesKafkaUiClusterSasl`

SASL credentials for the cluster connection. Must match the
target listener's authentication type.

### spec.clusters[].sasl.mechanism

`string` · required

SASL mechanism: "PLAIN", "SCRAM-SHA-256", or "SCRAM-SHA-512"
(Strimzi scram-sha-512 listeners pair with SCRAM-SHA-512).

- rule: mechanism must be PLAIN, SCRAM-SHA-256, or SCRAM-SHA-512
- rule: {"required":true}

### spec.clusters[].sasl.username

`string` · required

SASL username.

- rule: {"required":true}

### spec.clusters[].sasl.passwordSecret

`KubernetesKafkaUiPasswordSecret` · required

SASL password, from a key within a Secret — reference a
KubernetesKafkaUser resource (scram-sha-512 authentication) to
use its operator-generated Secret. The module wires it into the
console through a Secret-backed environment variable; it never
lands in the rendered configuration as plaintext.

- rule: {"required":true}

### spec.clusters[].sasl.passwordSecret.secretName

`string | valueFrom` · required

Name of the Secret. Accepts a literal Secret name or a reference
to a KubernetesKafkaUser resource, which resolves to its
credential Secret.

- references: KubernetesKafkaUser (`status.outputs.secret_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesKafkaUser, name: <that resource's name>, fieldPath: status.outputs.secret_name}} -- a bare string does not parse

### spec.clusters[].sasl.passwordSecret.key

`string` · optional (explicit presence)

The key within the Secret whose value is the password.
KubernetesKafkaUser credential Secrets carry it under
"password".

- default: `password`

### spec.clusters[].schemaRegistry

`KubernetesKafkaUiSchemaRegistry`

Schema registry serving this cluster — enables schema browsing
and schema-aware message rendering. Reference a
KubernetesKarapace sibling's endpoint.

### spec.clusters[].schemaRegistry.url

`string | valueFrom` · required

Registry URL. Accepts a literal URL or a reference to a
KubernetesKarapace resource, which resolves to its in-cluster
endpoint.

- references: KubernetesKarapace (`status.outputs.endpoint`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesKarapace, name: <that resource's name>, fieldPath: status.outputs.endpoint}} -- a bare string does not parse

### spec.clusters[].schemaRegistry.username

`string`

HTTP Basic username, when the registry requires it.

### spec.clusters[].schemaRegistry.passwordSecret

`KubernetesKafkaUiBasicAuthPasswordSecret`

HTTP Basic password source, when the registry requires it.

### spec.clusters[].schemaRegistry.passwordSecret.secretName

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.clusters[].schemaRegistry.passwordSecret.key

`string` · optional (explicit presence)

- default: `password`

### spec.clusters[].kafkaConnect

`[]KubernetesKafkaUiConnectCluster`

Connect clusters attached to this Kafka cluster — enables
connector monitoring and management from the console.

- rule: Connect cluster names must be unique within a cluster entry

### spec.clusters[].kafkaConnect[].name

`string` · required

Display name for this Connect cluster inside the console.

- rule: {"required":true}

### spec.clusters[].kafkaConnect[].address

`string | valueFrom` · required

Connect REST API address. Accepts a literal URL or a reference
to a KubernetesKafkaConnect resource, which resolves to its
in-cluster REST endpoint.

containment_exempt: the console MANAGES this Connect cluster
over its REST API; it is not deployed inside it.

- references: KubernetesKafkaConnect (`status.outputs.rest_api_endpoint`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesKafkaConnect, name: <that resource's name>, fieldPath: status.outputs.rest_api_endpoint}} -- a bare string does not parse

### spec.clusters[].kafkaConnect[].username

`string`

HTTP Basic username, when the Connect API sits behind one.

### spec.clusters[].kafkaConnect[].passwordSecret

`KubernetesKafkaUiBasicAuthPasswordSecret`

HTTP Basic password source, when the Connect API sits behind
one.

### spec.clusters[].kafkaConnect[].passwordSecret.secretName

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.clusters[].kafkaConnect[].passwordSecret.key

`string` · optional (explicit presence)

- default: `password`

### spec.clusters[].properties

`map<string, string>`

Additional Kafka client properties for this cluster's connection
(e.g. request timeouts). Values are configuration strings.
Security properties are operator-owned by the typed tls/sasl
blocks — never put credentials here.

### spec.auth

`KubernetesKafkaUiAuth`

Console login. Omitted = NO authentication — anyone reaching the
Service has full console access; acceptable only for
cluster-internal evaluation.

### spec.auth.type

`string` · required

Authentication type: "login_form" (ONE username/password account
— the app's form login authenticates against Spring's single
default security user; it has no multi-user store, verified in
the app source). Multi-user, OAuth2/OIDC and LDAP login ride the
chart's own Spring configuration through helm_values —
deliberately not typed here.

- rule: auth type must be login_form (OAuth2/LDAP compose through helm_values)
- rule: {"required":true}

### spec.auth.user

`KubernetesKafkaUiUser` · required

The console's single login account (login_form).

- rule: {"required":true}

### spec.auth.user.username

`string` · required

Login username.

- rule: {"required":true}

### spec.auth.user.password

`string` · required · sensitive

Login password — declared here, materialized by the module into
a Secret-backed environment variable (never plaintext in the
rendered chart values).

- rule: {"required":true}

### spec.replicas

`int32` · optional (explicit presence)

Number of console replicas. The console is stateless — replicas
are an availability measure.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.resources

`ContainerResources`

CPU/memory for the console pods.

### spec.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.resources.limits.cpu

`string`

### spec.resources.limits.memory

`string`

### spec.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.resources.requests.cpu

`string`

### spec.resources.requests.memory

`string`

### spec.serviceType

`string` · optional (explicit presence)

Service exposure: "ClusterIP" (default), "NodePort", or
"LoadBalancer". Prefer ClusterIP + a composed
Ingress/Gateway.

- default: `ClusterIP`
- rule: service_type must be ClusterIP, NodePort, or LoadBalancer

### spec.servicePort

`int32` · optional (explicit presence)

Service port the console is reachable on.

- default: `80`
- rule: {"int32":{"lte":65535,"gte":1}}

### spec.nodeSelector

`map<string, string>`

Node selector for the console pods.

### spec.tolerations

`[]WorkloadToleration`

Tolerations for the console pods.

### spec.tolerations[].key

`string`

Taint key to tolerate. Empty key with operator "Exists" tolerates every taint.

### spec.tolerations[].operator

`string`

How key/value match: "Equal" (default — value must match too) or "Exists"
(key presence alone matches).

- rule: Toleration operator must be either "Equal" or "Exists"

### spec.tolerations[].value

`string`

Taint value to match when operator is "Equal".

### spec.tolerations[].effect

`string`

Which taint effect is tolerated: "NoSchedule", "PreferNoSchedule", or
"NoExecute". Empty tolerates all effects for the key.

- rule: Toleration effect must be one of "NoSchedule", "PreferNoSchedule", or "NoExecute"

### spec.tolerations[].tolerationSeconds

`int64` · optional (explicit presence)

For "NoExecute" taints only: how many seconds already-running pods stay bound
after the taint appears. Unset means tolerate forever.

### spec.imageRegistry

`string`

Container image registry override (air-gapped mirrors). Empty =
ghcr.io.

### spec.helmValues

`string`

Escape hatch: additional chart values as a YAML document, merged
LAST over everything the typed fields render (Helm `-f`
semantics, identical on both engines). For the chart surface
beyond the typed fields (probes, security contexts, extra
volumes, OAuth2/LDAP login, ...) — never the substitute for
them. Do not put secrets here.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesKafkaUi, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace the console runs in. |
| `status.outputs.service_name` | `string` | Name of the console Service. |
| `status.outputs.endpoint` | `string` | In-cluster console endpoint (`http://<service>.<namespace>.svc.cluster.local:<port>`). |
| `status.outputs.port_forward_command` | `string` | Command to port-forward the console to localhost for quick access without any exposure. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.clusters[].bootstrapServers` | KubernetesKafka | `status.outputs.internal_bootstrap_endpoint` |
| `spec.clusters[].tls.caSecretName` | KubernetesKafka | `status.outputs.cluster_ca_cert_secret_name` |
| `spec.clusters[].sasl.passwordSecret.secretName` | KubernetesKafkaUser | `status.outputs.secret_name` |
| `spec.clusters[].schemaRegistry.url` | KubernetesKarapace | `status.outputs.endpoint` |
| `spec.clusters[].kafkaConnect[].address` | KubernetesKafkaConnect | `status.outputs.rest_api_endpoint` |

## See Also

- [Overview](../README.md)
