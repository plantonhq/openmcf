# Flink Deployment

A Flink cluster as one declaration: the FlinkDeployment CR the Apache
Flink Kubernetes Operator reconciles into a JobManager, its
TaskManagers, and — with `job` set — the one pipeline the cluster
exists to run (the recommended production grain; omit `job` for a
session cluster accepting external submissions). A
`KubernetesFlinkOperator` whose watch scope covers the namespace is
the prerequisite.

## Highlights

- **The operator's state rules, enforced at authoring time** —
  non-stateless upgrades need a checkpoints directory, savepoint mode
  needs a savepoints directory, JobManager standbys need HA metadata;
  the operator rejects deployments without them, and this spec tells
  you before the apply.
- **Credentials never render into config** — the S3 seam is
  Secret-native: endpoint plus credential references become pod
  environment at runtime, never entries in `flinkConfiguration`
  (which is a clear-text ConfigMap); a composed `KubernetesSeaweedFs`
  wires in through foreign-key defaults.
- **The plugin truth, told** — official Flink images ship the S3
  filesystem plugin disabled; `builtinPluginJar` activates the exact
  bundled jar, without which every `s3://` path fails at runtime with
  "unsupported filesystem scheme".
- **Version/image lockstep by construction** — the default image
  derives from `flinkVersion`, and the sizing arithmetic is spelled
  out: slots = TaskManagers × `taskmanager.numberOfTaskSlots`, and a
  job needs `parallelism` slots (an under-slotted cluster waits, it
  doesn't error).
- **Declarative operations** — suspend/resume as a field, a restart
  nonce, a savepoint-trigger nonce, and upgrade modes (`stateless`,
  `last-state`, `savepoint`) that say honestly what state they carry.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
