# KubernetesFlinkOperator: Research and Design

## Introduction

KubernetesFlinkOperator installs the Apache Flink Kubernetes Operator
— the official ASF controller for Flink on Kubernetes — as a single
Helm release named after `metadata.name`. The operator is the ENGINE
of the Flink story in this catalog: KubernetesFlinkDeployment
declares `FlinkDeployment` custom resources, and this operator
reconciles them into running Flink clusters. FlinkSessionJob,
FlinkStateSnapshot, and FlinkBlueGreenDeployment CRs are unmodeled
here — author them directly against the same operator when needed.

## The Chart Channel

The chart is served PER VERSION from a versioned Apache downloads
directory
(https://downloads.apache.org/flink/flink-kubernetes-operator-\<version\>/)
— the version is part of the repository URL itself, not just of the
chart pin. Chart version = operator version = the image tag the
modules pin. That pin matters doubly here: the chart's own default
image tag is the unpinned `latest`, the one values.yaml default that
must never stand — the modules always render `image.tag` at the chart
version.

## The Webhook Lifecycle: cert-manager Required, Fail-Closed

The operator ships an admission webhook that validates AND defaults
Flink CRs at admission time — enabled by upstream default, kept by
this spec. The chart truth at 1.15.0:

- **cert-manager is REQUIRED when the webhook is on.** The chart
  renders cert-manager Issuer/Certificate resources UNCONDITIONALLY,
  and the webhook trusts the API server through cert-manager's CA
  injection. There is NO self-signed fallback at this version —
  KubernetesCertManager is a hard prerequisite, not a suggestion.
- **Both webhook configurations are FAIL-CLOSED** (failurePolicy
  Fail). If the webhook cannot be reached — cert-manager absent,
  operator down — EVERY flink.apache.org admission in scope is
  rejected: a policy-engine class of blast radius, not a soft
  degradation. The flip side is the value: a bad FlinkDeployment
  (JobManager standbys without HA, a stateful upgrade mode without
  its directories) fails at admission with a real message instead of
  a silent reconcile stall.
- **`webhook_enabled: false` is the honest way out**: it removes the
  webhook, the certificate machinery, and the cert-manager
  dependency. The operator still validates in its reconcile loop —
  the same rules, but failures surface on CR status instead of at
  admission. The trade is deliberate and visible.
- **The webhook scoping follows the watch fence**: with
  `watch_namespaces` set, the chart scopes its RBAC AND the admission
  webhook's namespaceSelector to exactly that list — declarations
  outside it are neither reconciled nor webhook-validated.

## The Keystore Password: Why the Module Generates a Credential

The chart's default webhook keystore Secret carries a HARDCODED
PUBLIC PASSWORD (a base64 literal in the chart's webhook templates,
behind `keystore.useDefaultPassword: true`). It must never ship. The
modules:

1. Generate a random 32-character password per install
   (letters+digits only — the credential lands in a JVM keystore env
   value, and letters-and-digits alphabets avoid whole config-parser
   bug classes; the length compensates the smaller alphabet).
2. Materialize it as a module-owned Secret
   (`<name>-webhook-keystore`).
3. Wire the chart's `webhook.keystore.passwordSecretRef` at it, with
   `useDefaultPassword: false`.
4. RE-PIN `useDefaultPassword: false` in a values document merged
   AFTER the `helm_values` escape hatch — the chart's
   hardcoded-password default cannot resurface through user values.

The generation shape is ignored after creation, so an imported
credential never silently regenerates — rotation stays an explicit
verb, never plan fallout. The webhook-disabled arm creates no
random/Secret resources at all.

## The CRD Lifecycle: crds/ Directory, Keep by Construction

The chart ships its four flink.apache.org CRDs from its `crds/`
DIRECTORY: Helm installs them once, never upgrades them on chart
upgrades, and leaves them (and every Flink declaration) on uninstall
— no release ownership metadata. That upstream posture is exactly the
keep-on-uninstall this catalog wants for workload-bearing CRDs, so
the modules neither re-own nor template them. Chart-version bumps
never touch the CRDs: apply the new release's CRD files manually when
a bump changes them, per the upstream release notes.

## The Grain: Singleton per Namespace

The chart hardcodes its webhook Service, certificate, and issuer
names (`flink-operator-webhook-service`, `flink-operator-serving-cert`,
`flink-operator-selfsigned-issuer`) — they are chart-fixed, not
fullname-derived. A second release in the same namespace collides by
construction, which makes one operator per namespace the grain; one
cluster-wide-watching operator is the normal posture.

Two naming consequences:

- The chart-fixed webhook names are EXCLUDED from the name budget —
  they don't derive from the resource name.
- The 45-character budget comes from the module's own longest derived
  child, the `<name>-webhook-keystore` Secret (17-char suffix,
  keeping one character of headroom against the 63-character cap).
  Both modules fail loudly over the budget.

## Operator Configuration: Flink's Own Format, Cluster-Wide Reach

The operator is configured through Flink's own config format —
`kubernetes.operator.*` keys rendered into a `flink-conf.yaml`
document the chart APPENDS over its built-in defaults (note the
format is `key: value`, colon-space — YAML-flavored Flink
configuration, not properties `key=value`). Entries become
CLUSTER-WIDE defaults for every FlinkDeployment this operator
manages; per-pipeline configuration belongs on each
KubernetesFlinkDeployment, not here.

Leader election is module-owned, never a spec knob: any replica count
beyond 1 REQUIRES it (the chart's own contract — it refuses
multi-replica installs without it), so the two leader-election keys
(`kubernetes.operator.leader-election.enabled`, `.lease-name`) render
exactly when `replicas > 1`. A knob could drift from the replica
count; a derivation cannot.

## The Job Service Account

Flink JOB pods run as the chart-created `jobServiceAccount` (default
`flink` — the name every FlinkDeployment references by default),
planted with reconcile RBAC in the operator's namespace and, with
`watch_namespaces` set, in each watched namespace. The chart marks it
`helm.sh/resource-policy: keep`: it survives uninstall so running
jobs never lose their identity.

## Design Decisions

- **The install is blocking.** The Helm release waits for the
  operator to become Available (atomic, 600s timeout, cleanup on
  fail) — a JVM with a 30s-initial-delay startup probe, plus (webhook
  arm) a cert-manager certificate the webhook container mounts. An
  unpullable image, an absent cert-manager, or a broken config fails
  THIS apply with a readiness timeout instead of surfacing later as
  FlinkDeployments that mysteriously never reconcile.
- **The module owns namespace creation** (`create_namespace`), never
  the Helm release.
- **`fullnameOverride` pins the chart's fullname to the resource
  name** — the catalog's Helm-kind identity convention; every
  fullname-derived fallback name hangs off it.
- **Chart-default-matching values render only on divergence** — with
  the deliberate always-rendered exceptions: `fullnameOverride`,
  `image.tag` (the anti-`latest` pin), and the keystore wiring
  whenever the webhook is on.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Chart | `flink-kubernetes-operator`, served per version from the Apache downloads directory | Pinned 1.15.0 (spec default); chart version = operator version = image tag |
| Operator image | `ghcr.io/apache/flink-kubernetes-operator` | Tag always pinned to the chart version (the chart's own default is `latest`) |
| CRD API group | `flink.apache.org` | Four CRDs, chart `crds/` directory — installed once, kept on uninstall, never upgraded by chart bumps |
| Webhook Service | `flink-operator-webhook-service` (chart-fixed) | Exported as `webhook_service`; empty when the webhook is disabled |
| Webhook certificate/issuer | `flink-operator-serving-cert` / `flink-operator-selfsigned-issuer` (chart-fixed) | The names behind the singleton-per-namespace grain |
| Keystore Secret | `<name>-webhook-keystore` (module-generated) | The 17-char suffix behind the 45-char name budget |
| Job service account | `flink` (chart default) | Exported as `job_service_account`; `helm.sh/resource-policy: keep` |

## IaC Twins

Pulumi (`module/values.go` + `module/keystore_secret.go`) and
Terraform (`main.tf` + `locals.tf`) render identical chart values,
generate the same keystore-password Secret shape, and merge the same
three documents in the same order: typed values, the `helm_values`
escape hatch, then the keystore re-pin. Keep the typed-value
rendering, the leader-election derivation, and the post-merge re-pin
in lockstep.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
