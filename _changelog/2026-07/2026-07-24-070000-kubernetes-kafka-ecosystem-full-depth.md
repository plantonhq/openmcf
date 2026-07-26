# Kubernetes Kafka ecosystem: Connect, Connector, MirrorMaker 2, Karapace schema registry, and kafbat UI — five new kinds at full depth

## What changed

- **KubernetesKafkaConnect (new, 911)** — the Strimzi `KafkaConnect`
  integration engine as a first-class kind: bootstrap wiring through a
  KubernetesKafka reference, the shared client TLS/authentication arms,
  name-derived group identity (group.id + the three storage topics,
  with the uniqueness contract documented on the fields), worker
  configuration with the operator-owned prefixes documented, and FOUR
  plugin-delivery arms — the stock image (which carries ONLY the
  MirrorMaker 2 connector classes; verified live against the workers'
  own plugin listing — Kafka's FileStream examples are NOT on the
  distribution's classpath), a prebuilt image, OCI image-volume plugin
  mounts (requires the cluster's ImageVolume feature; the mount
  mechanism is live-proven), and operator-driven image builds
  (docker/imagestream outputs, maven/url artifacts, checksum pinning).
  `image` and `build` are mutually exclusive at validation — with a
  build configured the operator deploys the image IT builds and
  silently overrides a declared one (verified in the operator source).
  Every cluster is annotated `strimzi.io/use-connector-resources`
  unconditionally, so connectors are managed declaratively.
- **KubernetesKafkaConnector (new, 912)** — one declared pipe on a
  Connect cluster: class/tasks/config/version, desired state
  (running/paused/stopped), automatic restart with back-off, and
  offset management modeled honestly — the spec declares the
  list/alter ConfigMap TARGETS while the `strimzi.io/connector-offsets`
  annotation is the VERB (the replay, skip-poison-record and
  migration-cutover mechanism). The placement contract (same namespace
  as the Connect cluster, the strimzi.io/cluster label) sits on the
  fields their authors read.
- **KubernetesKafkaMirrorMaker2 (new, 913)** — continuous, offset-aware
  replication into one TARGET cluster from one or more per-mirror
  SOURCE clusters (the current CRD shape — there is no shared clusters
  list): topic/group include+exclude patterns, per-mirror source and
  checkpoint connector tuning with automatic restarts, and the
  keep-original-topic-names migration posture documented
  (IdentityReplicationPolicy on both connectors). The target alias must
  differ from every source alias at validation (including the default
  "target"). The live proof: records produced on a source cluster are
  consumed from the target under the mirrored name.
- **KubernetesKarapace (new, 914)** — the Apache-2.0,
  Confluent-API-compatible schema registry, deployed as module-owned
  typed manifests (upstream ships no Helm chart): registry and optional
  REST-proxy Deployments/Services configured through KARAPACE_*
  environment variables exactly as upstream runs them, Kafka
  PLAINTEXT/SSL/SASL connection arms (a literal SASL password is
  materialized into a Secret and referenced — never plaintext in a pod
  spec), schemas stored Kafka-natively in a compacted topic, leader
  election across replicas with each pod advertising its own POD IP
  (a Deployment pod's bare name does not resolve in cluster DNS —
  follower-to-leader write forwarding depends on this), server TLS via
  the cert-manager seam, and authfile/OIDC HTTP authentication. The
  live proof registers a schema and reads it back through the SR API.
- **KubernetesKafkaUi (new, 915)** — the kafbat UI console from the
  served `kafka-ui` chart (1.6.4): one console observing many clusters,
  with each cluster entry composing the sibling kinds through
  references (Kafka bootstrap, Karapace registry endpoint, Connect REST
  endpoints), per-cluster read-only posture, TLS trust via a
  PEM-truststore mount, SASL credentials wired through environment
  placeholders + secret mappings so no credential ever lands in
  rendered values, and a single login-form account (the app
  authenticates against Spring's default security user — it has no
  multi-user store, verified in the app source; OAuth2/LDAP compose
  through helm_values). The live proofs: the console's own API reports
  the cluster ONLINE, and anonymous access is refused when auth is on.
- **Shared client-connection proto** — Connect, MirrorMaker 2's target
  and every mirror source declare Kafka connections through one shared
  message set (TLS trusted certificates; tls/scram-sha-512/
  scram-sha-256/plain/custom authentication with KubernetesKafkaUser
  credential-Secret references as the default wiring), so the three
  sites cannot drift.

## Validation

Spec suites green across the family (every CEL rule accept+reject
locked); dual-engine E2E green on all 22 scenario-engine lanes on the
kind cluster, including the four behavioral proofs (migration,
data-flow, register/fetch, observe) and the OCI image-volume mount
lane; blind import round-trips proven for all five kinds; offline plan
proofs full-surface and minimal per kind with type-fidelity checks;
secret-coverage, reference, import-map, outputs-conformance and
license-footer gates green.
