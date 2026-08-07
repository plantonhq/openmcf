# KubernetesKafkaUser

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**KubernetesKafkaUserSpec** declares a Kafka client identity on the
Strimzi `KafkaUser` custom resource. The target cluster's USER
OPERATOR (enabled by default on KubernetesKafka) reconciles it into
a real principal: it generates credentials into a Secret named after
the user (exported as `secret_name`), and — when the cluster runs
`simple` authorization — applies the declared ACLs.

PLACEMENT CONTRACT (verified against the Strimzi operator): the
KafkaUser must live in the SAME NAMESPACE as its Kafka cluster, and
it binds to the cluster through the strimzi.io/cluster label
(rendered from `kafka_cluster`). A user in another namespace is
accepted by the API server and then silently never reconciled.

MATCH THE LISTENER: a user's authentication type must match a
listener's authentication type on the target cluster — a
scram-sha-512 user cannot authenticate on a tls-auth listener and
vice versa. The generated Secret carries `password` (and a ready
`sasl.jaas.config`) for scram-sha-512 users, or `user.crt` /
`user.key` (plus keystore forms) for tls users.

## Example

```yaml
# Full-surface development manifest: SCRAM authentication, a
# producer+consumer ACL set spanning literal and prefix patterns plus a
# consumer group, and every quota field.
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesKafkaUser
metadata:
  name: analytics-pipeline
spec:
  namespace:
    value: kafka-hack
  kafkaCluster:
    value: kafka-hack
  authentication:
    type: scram-sha-512
  authorization:
    type: simple
    acls:
      - resource:
          type: topic
          name: orders.v1_events
        operations:
          - Read
          - Describe
        host: "*"
      - resource:
          type: topic
          name: analytics.
          patternType: prefix
        operations:
          - Write
          - Create
          - Describe
          - IdempotentWrite
      - resource:
          type: group
          name: analytics-pipeline
        operations:
          - Read
      - resource:
          type: cluster
        operations:
          - DescribeConfigs
  quotas:
    producerByteRate: 1048576
    consumerByteRate: 2097152
    requestPercentage: 55
    controllerMutationRate: 10
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.kafkaCluster` | `string \| valueFrom` | yes |  | KubernetesKafka (`status.outputs.cluster_name`) |
| `spec.authentication` | `KubernetesKafkaUserAuthentication` |  |  |  |
| `spec.authentication.type` | `string` | yes |  |  |
| `spec.authorization` | `KubernetesKafkaUserAuthorization` |  |  |  |
| `spec.authorization.type` | `string` |  | `simple` |  |
| `spec.authorization.acls` | `[]KubernetesKafkaUserAcl` | yes |  |  |
| `spec.authorization.acls[].resource` | `KubernetesKafkaUserAclResource` | yes |  |  |
| `spec.authorization.acls[].resource.type` | `string` | yes |  |  |
| `spec.authorization.acls[].resource.name` | `string` |  |  |  |
| `spec.authorization.acls[].resource.patternType` | `string` |  | `literal` |  |
| `spec.authorization.acls[].operations` | `[]string` | yes |  |  |
| `spec.authorization.acls[].host` | `string` |  |  |  |
| `spec.quotas` | `KubernetesKafkaUserQuotas` |  |  |  |
| `spec.quotas.producerByteRate` | `int32` |  |  |  |
| `spec.quotas.consumerByteRate` | `int32` |  |  |  |
| `spec.quotas.requestPercentage` | `int32` |  |  |  |
| `spec.quotas.controllerMutationRate` | `double` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace of the KafkaUser — MUST be the Kafka cluster's own
namespace (the user operator watches only there; see the
placement contract above). Accepts a literal namespace name or a
reference to a KubernetesNamespace resource.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.kafkaCluster

`string | valueFrom` · required

The Kafka cluster this user belongs to. Accepts a literal cluster
name (the KubernetesKafka resource's metadata.name) or a
reference to a KubernetesKafka resource. Rendered as the
strimzi.io/cluster label.

- references: KubernetesKafka (`status.outputs.cluster_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesKafka, name: <that resource's name>, fieldPath: status.outputs.cluster_name}} -- a bare string does not parse

### spec.authentication

`KubernetesKafkaUserAuthentication`

How this user authenticates. Required for credential-bearing
users; omit only for a principal that exists purely to carry
ACLs for an externally-authenticated identity (custom listener
mechanisms).

### spec.authentication.type

`string` · required

Authentication type:
"scram-sha-512" (username/password — the user operator generates
the password into the user's Secret; pairs with scram-sha-512
listeners),
"tls" (mutual TLS — the user operator issues a client
certificate from the cluster's clients CA into the user's
Secret; pairs with tls-auth listeners), or
"tls-external" (mutual TLS with certificates issued OUTSIDE the
cluster — no Secret is generated; the principal is the
certificate's subject).

- rule: authentication type must be scram-sha-512, tls, or tls-external
- rule: {"required":true}

### spec.authorization

`KubernetesKafkaUserAuthorization`

ACLs for this user — requires the cluster's authorization to be
`simple`. On a cluster WITHOUT simple authorization the user
operator REJECTS the resource outright ("Simple authorization ACL
rules are configured but not supported in the Kafka cluster
configuration", the KafkaUser reports NotReady and no credentials
are generated — verified in the operator source). Declare ACLs
only against clusters that enforce them.

### spec.authorization.type

`string` · optional (explicit presence)

Authorization declaration type — "simple" is the only Strimzi
type (Kafka's built-in ACL model).

- default: `simple`
- rule: authorization type must be simple

### spec.authorization.acls

`[]KubernetesKafkaUserAcl` · required

The ACL rules granting this user access.

- rule: {"repeated":{"minItems":"1"}}

### spec.authorization.acls[].resource

`KubernetesKafkaUserAclResource` · required

What the rule applies to.

- rule: {"required":true}
- rule: resource name is required for topic, group and transactionalId ACLs (only cluster ACLs are nameless)

### spec.authorization.acls[].resource.type

`string` · required

Resource type: "topic", "group" (consumer group),
"cluster", or "transactionalId".

- rule: resource type must be topic, group, cluster, or transactionalId
- rule: {"required":true}

### spec.authorization.acls[].resource.name

`string`

Resource name (topic name, group id, transactional id). With
pattern_type "prefix", the rule covers every resource whose name
starts with this value. Not used for type "cluster".

### spec.authorization.acls[].resource.patternType

`string` · optional (explicit presence)

Name matching: "literal" (exact, the default) or "prefix".

- default: `literal`
- rule: pattern_type must be literal or prefix

### spec.authorization.acls[].operations

`[]string` · required

Operations granted: Read, Write, Create, Delete, Alter, Describe,
ClusterAction, AlterConfigs, DescribeConfigs, IdempotentWrite, or
All. A typical producer needs Write + Describe (+ IdempotentWrite
for idempotent producers); a consumer needs Read + Describe on
the topic and Read on its consumer group.

- rule: each operation must be one of Read, Write, Create, Delete, Alter, Describe, ClusterAction, AlterConfigs, DescribeConfigs, IdempotentWrite, All
- rule: {"repeated":{"minItems":"1"}}

### spec.authorization.acls[].host

`string`

Host restriction for the rule (Kafka ACL host). Empty = "*"
(any host).

### spec.quotas

`KubernetesKafkaUserQuotas`

Client quotas for this user, enforced by the brokers.

### spec.quotas.producerByteRate

`int32` · optional (explicit presence)

Producer throughput cap, bytes/second (before acknowledgement
throttling).

- rule: {"int32":{"gte":0}}

### spec.quotas.consumerByteRate

`int32` · optional (explicit presence)

Consumer throughput cap, bytes/second.

- rule: {"int32":{"gte":0}}

### spec.quotas.requestPercentage

`int32` · optional (explicit presence)

Share of broker request-handler time this user may consume,
as a percentage (can exceed 100 — it is per-thread-group
percentage in Kafka's model).

- rule: {"int32":{"gte":0}}

### spec.quotas.controllerMutationRate

`double` · optional (explicit presence)

Cap on partition-mutation operations per second (create/delete
partitions) attributable to this user.

- rule: {"double":{"gte":0}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesKafkaUser, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace the KafkaUser resource lives in (the Kafka cluster's namespace). |
| `status.outputs.username` | `string` | Kafka principal name (`metadata.name`) — combine with the authentication type for ACL/super-user principal form ("User:<name>" for scram-sha-512; "User:CN=<name>" for tls). |
| `status.outputs.secret_name` | `string` | Name of the credentials Secret the user operator generates (equal to the user name). Keys: `password` + `sasl.jaas.config` for scram-sha-512 users; `user.crt` / `user.key` / `user.p12` / `user.password` for tls users. Empty semantics for tls-external users (no Secret is generated). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.kafkaCluster` | KubernetesKafka | `status.outputs.cluster_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| KubernetesKafkaConnect | `spec.authentication.certificateAndKey.secretName` | `status.outputs.secret_name` |
| KubernetesKafkaConnect | `spec.authentication.passwordSecret.secretName` | `status.outputs.secret_name` |
| KubernetesKafkaMirrorMaker2 | `spec.target.authentication.certificateAndKey.secretName` | `status.outputs.secret_name` |
| KubernetesKafkaMirrorMaker2 | `spec.target.authentication.passwordSecret.secretName` | `status.outputs.secret_name` |
| KubernetesKafkaMirrorMaker2 | `spec.mirrors[].source.authentication.certificateAndKey.secretName` | `status.outputs.secret_name` |
| KubernetesKafkaMirrorMaker2 | `spec.mirrors[].source.authentication.passwordSecret.secretName` | `status.outputs.secret_name` |
| KubernetesKafkaUi | `spec.clusters[].sasl.passwordSecret.secretName` | `status.outputs.secret_name` |
| KubernetesKarapace | `spec.kafka.tls.clientCertSecretName` | `status.outputs.secret_name` |
| KubernetesKarapace | `spec.kafka.sasl.passwordSecret.secretName` | `status.outputs.secret_name` |

## See Also

- [Overview](../README.md)
