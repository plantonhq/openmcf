# Default preset

The standard operator install: the pinned version (1.15.0, the signed
Apache per-version chart channel) into its own `flink-system`
namespace, watching every namespace, the admission webhook ON.

PREREQUISITE: a `KubernetesCertManager` on the cluster. With the
webhook enabled (the upstream default this component keeps), the
chart's webhook certificate is issued and rotated by cert-manager —
there is NO self-signed fallback at this version — and BOTH webhook
configurations are FAIL-CLOSED: if the webhook cannot be reached,
every `flink.apache.org` admission in scope is rejected. That is the
policy-engine class of blast radius; it is also what makes bad Flink
declarations fail at admission with a real message instead of a
silent reconcile stall. `webhookEnabled: false` removes the webhook,
the certificates, and the cert-manager dependency — validation then
happens in the reconcile loop.

The credential truth: the chart's default webhook keystore password is
a HARDCODED public value — it never ships from this component; the
module generates a random password Secret (`<name>-webhook-keystore`)
per install.

What the install owns: the four `flink.apache.org` CRDs ride the
chart's `crds/` directory — installed once, never upgraded by chart
bumps, KEPT on uninstall. The `flink` job service account (the
runner identity every `KubernetesFlinkDeployment` references) is
chart-created and survives uninstall by upstream design.

One operator per NAMESPACE by construction (the webhook artifact names
are chart-fixed); one cluster-wide-watching operator is the normal
posture. Declare Flink clusters with `KubernetesFlinkDeployment`.

Change first: nothing, usually. Reach for `watchNamespaces` to fence
Flink to team namespaces, `replicas: 2` for a leader-elected standby.

See [01-default.yaml](./01-default.yaml) for the manifest.
