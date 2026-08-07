# SCRAM app cluster preset

A quorum-safe application cluster with exactly one credential story:
a single SCRAM-SHA-512-over-TLS listener plus simple ACL
authorization. Every client — including the admin — is a
KubernetesKafkaUser resource: the user operator generates the
password into a Secret, and the user's `authorization` block declares
the ACLs the authorizer enforces. No anonymous path exists on this
cluster.

For application teams who want per-service credentials and
least-privilege topic access without operating a certificate
hierarchy for clients — SCRAM passwords rotate by rotating the
Secret, and services consume them as env-from references. The
trade-off against mutual TLS: identity is a password, not a
certificate, so transport TLS (already on) is doing the
eavesdropping protection while SCRAM does authentication. The super
user is the SCRAM principal `User:kafka-admin` (SCRAM principals are
`User:<username>`; TLS principals are `User:CN=<name>`) — declare a
KubernetesKafkaUser named `kafka-admin` with `scram-sha-512`
authentication to materialize its credentials before anything needs
admin access.

See [04-scram-app-cluster.yaml](./04-scram-app-cluster.yaml) for the
manifest.
