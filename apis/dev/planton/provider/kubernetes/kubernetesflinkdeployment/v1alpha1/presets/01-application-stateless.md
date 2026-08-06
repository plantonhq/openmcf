# Application (stateless) preset

The recommended production grain: one APPLICATION cluster per
pipeline — the cluster exists to run this one job and follows its
lifecycle, with a custom image that BAKES the job jar (`local:///`
paths point inside the image; remote schemes download at submission).
Keep the Flink inside the image identical to `flinkVersion` — the
operator shapes its submission protocol from the declared version and
a mismatch fails at runtime, not at apply.

PREREQUISITE: a `KubernetesFlinkOperator` whose watch scope covers
this namespace. The pods run as the `flink` service account its chart
creates.

This preset is STATELESS on purpose: spec changes restart the job
from clean state — correct ONLY for pipelines with no state worth
carrying (enrichment, routing, filtering). The moment your pipeline
keeps state, move to the stateful preset: `upgradeMode: stateless` on
a stateful job silently discards its state on every upgrade.

Sizing truth: total task slots = TaskManager count ×
`taskmanager.numberOfTaskSlots` (a `flinkConfiguration` key, default
1), and the job needs `parallelism` slots — in native mode (the
default) Flink requests TaskManagers on demand, so this preset yields
4 TaskManagers of 1 slot each. An under-slotted cluster holds the job
in a scheduling wait, not an error.

Change first: `image` and `jarUri` to your pipeline's; `parallelism`
to its throughput.

See [01-application-stateless.yaml](./01-application-stateless.yaml)
for the manifest.
