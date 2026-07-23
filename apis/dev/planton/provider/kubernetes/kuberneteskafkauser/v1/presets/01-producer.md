# Producer (SCRAM)

This preset declares a producer identity: SCRAM-SHA-512 authentication with an operator-generated password, and a prefix ACL granting write access to every topic starting with `orders-`. The user operator generates the credentials into a Secret named after the user — `password` plus a ready-made `sasl.jaas.config` the client consumes directly; nothing sensitive lives in this manifest.

## When to Use

- A service that publishes to one family of topics and reads none
- SCRAM listeners — username/password is the simplest client wiring (one JAAS line), the common choice for in-cluster application traffic

## Key Configuration Choices

- **`type: scram-sha-512`** -- must match a scram-sha-512 listener on the target cluster; a SCRAM user cannot authenticate on a tls-auth listener
- **Prefix ACL** (`patternType: prefix`, name `orders-`) -- one rule covers every current and future topic in the family; switch to `literal` (the default) to pin a single topic
- **`Write` + `Describe`** -- the canonical producer grant
- **`Create`** -- lets the producer create missing topics; drop it on clusters where topics are strictly declared (the KubernetesKafkaTopic posture)
- Idempotent producers need no extra grant on modern Kafka -- `Write`
  on the topic covers them (the legacy `IdempotentWrite` operation
  applies to the cluster resource, not topics)
- **ACLs enforced only with cluster-side `simple` authorization** -- declared against a cluster without it, the rules are rejected at reconcile and the resource reports NotReady

## Values to Adapt

| Value | Description | Where to Find |
|---|---|---|
| `kafka` (namespace) | The Kafka cluster's own namespace | The KubernetesKafka resource's `namespace` |
| `my-kafka` (kafkaCluster) | The Kafka cluster's name | The KubernetesKafka resource's `metadata.name` or its `cluster_name` output |
| `orders-` | The topic-name prefix this producer may write to | Your topic naming convention / KubernetesKafkaTopic resources |

## Related Presets

- **02-consumer** -- The read-side counterpart: topic Read/Describe plus consumer-group Read
- **03-mtls-service** -- Certificate-based identity with quotas, for tls-auth listeners
