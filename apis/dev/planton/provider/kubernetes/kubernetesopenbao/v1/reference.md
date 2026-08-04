# KubernetesOpenBao

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1`

**KubernetesOpenBaoSpec** installs OpenBao — the open-source,
Linux Foundation-governed secrets manager (MPL-2.0 fork of Vault) —
from the official `openbao` chart
(https://openbao.github.io/openbao-helm, chart 0.28.x = OpenBao 2.6.x).

THE SEAL LIFECYCLE — the fact everything else follows from: a fresh
OpenBao server starts UNINITIALIZED and SEALED. Initialization
(`bao operator init` — generates the unseal key shares and the
initial root token) and unsealing are RUNTIME operations performed
against the API after deploy; no Kubernetes deployment tool can do
them declaratively, and this component deliberately does not try.
Until a server is initialized and unsealed, its pod reports
NotReady BY DESIGN (the readiness probe is `bao status`, which
exits non-zero for sealed servers) — the chart keeps sealed pods
addressable through its Services (publishNotReadyAddresses), so
`kubectl port-forward` and the DNS names work for the init/unseal
calls. A deployment that never becomes "ready" until you initialize
it is the designed behavior, not a failure. Auto-unseal (below)
removes the UNSEAL step from restarts, but the one-time
initialization is always yours.

ONE SERVER MODE at a time: dev XOR standalone XOR ha (Raft). When
no mode is declared, standalone is used — the chart's own default.
The server always runs as a StatefulSet with an OnDelete update
strategy (config changes never roll pods automatically; delete pods
to pick up config).

## Example

```yaml
# Full-surface development manifest — exercises every module-rendered arm
# so the offline plan/preview proofs cover what the kind-cluster lanes
# exclude (HA Raft with synthesized retry_join, TLS listener wiring, a
# declared-credential auto-unseal seal, the injector, metrics +
# ServiceMonitor, audit storage, the snapshot agent).
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesOpenBao
metadata:
  name: bao-dev
spec:
  namespace:
    value: openbao-dev
  createNamespace: true
  server:
    ha:
      replicas: 3
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 512Mi
    dataStorage:
      size: 10Gi
    auditStorage:
      size: 5Gi
    logLevel: info
    logFormat: json
    scheduling:
      nodeSelector:
        kubernetes.io/os: linux
      tolerations:
        - key: dedicated
          operator: Equal
          value: secrets
          effect: NoSchedule
  tls:
    enabled: true
    certSecretName:
      value: bao-dev-tls
  autoUnseal:
    awsKms:
      region: us-west-2
      kmsKeyId: alias/openbao-unseal
      accessKeyId: AKIAEXAMPLEDEVONLY
      secretAccessKey: dev-only-placeholder-secret-key
  injector:
    enabled: true
    replicas: 1
    failurePolicy: Ignore
    resources:
      requests:
        cpu: 50m
        memory: 64Mi
      limits:
        cpu: 250m
        memory: 128Mi
  uiEnabled: true
  networkPolicyEnabled: true
  metrics:
    enabled: true
    serviceMonitorEnabled: true
  snapshotAgent:
    enabled: true
    schedule: "0 */6 * * *"
    s3Host:
      value: s3.us-west-2.amazonaws.com
    s3Bucket: bao-dev-snapshots
    s3ExpireDays: 7
    s3CredentialsSecretName: bao-dev-s3-credentials
    baoRole: snapshot
    baoAuthPath: kubernetes
  serviceAccount:
    annotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::111122223333:role/bao-dev-unseal
    authDelegatorEnabled: true
  helmValues: |
    server:
      annotations:
        example.planton.ai/full-surface: "true"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.chartVersion` | `string` |  | `0.28.6` |  |
| `spec.server` | `KubernetesOpenBaoServer` |  |  |  |
| `spec.server.dev` | `KubernetesOpenBaoDevMode` |  |  |  |
| `spec.server.standalone` | `KubernetesOpenBaoStandaloneMode` |  |  |  |
| `spec.server.ha` | `KubernetesOpenBaoHaMode` |  |  |  |
| `spec.server.ha.replicas` | `int32` |  | `3` |  |
| `spec.server.resources` | `ContainerResources` |  |  |  |
| `spec.server.resources.limits` | `CpuMemory` |  |  |  |
| `spec.server.resources.limits.cpu` | `string` |  |  |  |
| `spec.server.resources.limits.memory` | `string` |  |  |  |
| `spec.server.resources.requests` | `CpuMemory` |  |  |  |
| `spec.server.resources.requests.cpu` | `string` |  |  |  |
| `spec.server.resources.requests.memory` | `string` |  |  |  |
| `spec.server.dataStorage` | `KubernetesOpenBaoStorage` |  |  |  |
| `spec.server.dataStorage.size` | `string` |  | `10Gi` |  |
| `spec.server.dataStorage.storageClass` | `string \| valueFrom` |  |  | KubernetesStorageClass (`status.outputs.storage_class_name`) |
| `spec.server.auditStorage` | `KubernetesOpenBaoStorage` |  |  |  |
| `spec.server.auditStorage.size` | `string` |  | `10Gi` |  |
| `spec.server.auditStorage.storageClass` | `string \| valueFrom` |  |  | KubernetesStorageClass (`status.outputs.storage_class_name`) |
| `spec.server.logLevel` | `string` |  | `info` |  |
| `spec.server.logFormat` | `string` |  | `standard` |  |
| `spec.server.scheduling` | `KubernetesOpenBaoScheduling` |  |  |  |
| `spec.server.scheduling.nodeSelector` | `map<string, string>` |  |  |  |
| `spec.server.scheduling.tolerations` | `[]WorkloadToleration` |  |  |  |
| `spec.server.scheduling.tolerations[].key` | `string` |  |  |  |
| `spec.server.scheduling.tolerations[].operator` | `string` |  |  |  |
| `spec.server.scheduling.tolerations[].value` | `string` |  |  |  |
| `spec.server.scheduling.tolerations[].effect` | `string` |  |  |  |
| `spec.server.scheduling.tolerations[].tolerationSeconds` | `int64` |  |  |  |
| `spec.tls` | `KubernetesOpenBaoTls` |  |  |  |
| `spec.tls.enabled` | `bool` |  |  |  |
| `spec.tls.certSecretName` | `string \| valueFrom` |  |  | KubernetesCertificate (`status.outputs.secret_name`) |
| `spec.autoUnseal` | `KubernetesOpenBaoAutoUnseal` |  |  |  |
| `spec.autoUnseal.awsKms` | `KubernetesOpenBaoAwsKmsSeal` |  |  |  |
| `spec.autoUnseal.awsKms.region` | `string` | yes |  |  |
| `spec.autoUnseal.awsKms.kmsKeyId` | `string` | yes |  |  |
| `spec.autoUnseal.awsKms.accessKeyId` | `string` |  |  |  |
| `spec.autoUnseal.awsKms.secretAccessKey` | `string` (sensitive) |  |  |  |
| `spec.autoUnseal.gcpKms` | `KubernetesOpenBaoGcpKmsSeal` |  |  |  |
| `spec.autoUnseal.gcpKms.project` | `string \| valueFrom` | yes |  | GcpProject (`status.outputs.project_id`) |
| `spec.autoUnseal.gcpKms.region` | `string` | yes |  |  |
| `spec.autoUnseal.gcpKms.keyRing` | `string \| valueFrom` | yes |  | GcpKmsKeyRing (`status.outputs.key_ring_name`) |
| `spec.autoUnseal.gcpKms.cryptoKey` | `string \| valueFrom` | yes |  | GcpKmsKey (`status.outputs.key_name`) |
| `spec.autoUnseal.gcpKms.workloadIdentityServiceAccount` | `string \| valueFrom` |  |  | GcpServiceAccount (`status.outputs.email`) |
| `spec.autoUnseal.azureKeyVault` | `KubernetesOpenBaoAzureKeyVaultSeal` |  |  |  |
| `spec.autoUnseal.azureKeyVault.vaultName` | `string` | yes |  |  |
| `spec.autoUnseal.azureKeyVault.keyName` | `string` | yes |  |  |
| `spec.autoUnseal.azureKeyVault.tenantId` | `string` | yes |  |  |
| `spec.autoUnseal.azureKeyVault.clientId` | `string` |  |  |  |
| `spec.autoUnseal.azureKeyVault.clientSecret` | `string` (sensitive) |  |  |  |
| `spec.autoUnseal.transit` | `KubernetesOpenBaoTransitSeal` |  |  |  |
| `spec.autoUnseal.transit.address` | `string` | yes |  |  |
| `spec.autoUnseal.transit.keyName` | `string` | yes |  |  |
| `spec.autoUnseal.transit.mountPath` | `string` |  | `transit/` |  |
| `spec.autoUnseal.transit.token` | `string` (sensitive) |  |  |  |
| `spec.injector` | `KubernetesOpenBaoInjector` |  |  |  |
| `spec.injector.enabled` | `bool` |  |  |  |
| `spec.injector.replicas` | `int32` |  | `1` |  |
| `spec.injector.failurePolicy` | `string` |  | `Ignore` |  |
| `spec.injector.resources` | `ContainerResources` |  |  |  |
| `spec.injector.resources.limits` | `CpuMemory` |  |  |  |
| `spec.injector.resources.limits.cpu` | `string` |  |  |  |
| `spec.injector.resources.limits.memory` | `string` |  |  |  |
| `spec.injector.resources.requests` | `CpuMemory` |  |  |  |
| `spec.injector.resources.requests.cpu` | `string` |  |  |  |
| `spec.injector.resources.requests.memory` | `string` |  |  |  |
| `spec.uiEnabled` | `bool` |  | `true` |  |
| `spec.networkPolicyEnabled` | `bool` |  |  |  |
| `spec.metrics` | `KubernetesOpenBaoMetrics` |  |  |  |
| `spec.metrics.enabled` | `bool` |  |  |  |
| `spec.metrics.serviceMonitorEnabled` | `bool` |  |  |  |
| `spec.snapshotAgent` | `KubernetesOpenBaoSnapshotAgent` |  |  |  |
| `spec.snapshotAgent.enabled` | `bool` |  |  |  |
| `spec.snapshotAgent.schedule` | `string` |  | `*/15 * * * *` |  |
| `spec.snapshotAgent.s3Host` | `string \| valueFrom` | yes |  | KubernetesSeaweedFs (`status.outputs.s3_endpoint`) |
| `spec.snapshotAgent.s3Bucket` | `string` | yes |  |  |
| `spec.snapshotAgent.s3ExpireDays` | `int32` |  | `14` |  |
| `spec.snapshotAgent.s3CredentialsSecretName` | `string` | yes |  |  |
| `spec.snapshotAgent.baoRole` | `string` |  | `snapshot` |  |
| `spec.snapshotAgent.baoAuthPath` | `string` |  | `kubernetes` |  |
| `spec.serviceAccount` | `KubernetesOpenBaoServiceAccount` |  |  |  |
| `spec.serviceAccount.annotations` | `map<string, string>` |  |  |  |
| `spec.serviceAccount.authDelegatorEnabled` | `bool` |  | `true` |  |
| `spec.helmValues` | `string` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace to install into (conventionally "openbao"). Accepts a
literal namespace name or a reference to a KubernetesNamespace
resource.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.createNamespace

`bool`

When true, the namespace is created (with the standard Planton
governance labels) before installing and deleted with the
resource. When false, the namespace must already exist.

### spec.chartVersion

`string` · optional (explicit presence)

Helm chart version to install (e.g. "0.28.6" = OpenBao v2.6.1 —
the chart's appVersion pins the server image). Versions must
exist in the SERVED index at https://openbao.github.io/openbao-helm.

- default: `0.28.6`

### spec.server

`KubernetesOpenBaoServer`

The OpenBao server: mode, sizing, storage, and logging.

### spec.server.dev

`KubernetesOpenBaoDevMode`

Dev mode: in-memory, auto-initialized, auto-unsealed, root
token literally "root". NEVER for real secrets — all data is
lost on every restart, and the root token is plaintext in the
pod spec (readable by anyone who can get pods). Exists so the
component can be evaluated and composed against without the
init/unseal ceremony. No PVC is created in dev mode, and
workload-identity ServiceAccount annotations are NOT applied
(a chart behavior — dev mode drops them).

### spec.server.standalone

`KubernetesOpenBaoStandaloneMode`

Standalone: one instance, `storage "file"` on the data PVC.
The production shape for single-instance installs.

### spec.server.ha

`KubernetesOpenBaoHaMode`

High availability with integrated Raft storage: every replica
persists to its own data PVC and the cluster elects a leader.
This module renders `retry_join` stanzas for every peer (the
chart alone ships NONE — without them a multi-replica Raft
install never forms a cluster and each pod sits uninitialized
and independent). Bootstrap: initialize pod-0 and unseal every
pod; joins then happen automatically through retry_join.

### spec.server.ha.replicas

`int32` · optional (explicit presence)

Number of server replicas (Raft peers). Odd counts (3, 5)
tolerate minority loss; 3 is the standard production shape. A
single replica is a legal Raft cluster of one (useful in labs).
Remember the chart's default required anti-affinity: replicas
beyond the node count stay Pending (see scheduling).

- default: `3`
- rule: {"int32":{"lte":11,"gte":1}}

### spec.server.resources

`ContainerResources`

CPU and memory for the server container. The chart ships no
defaults; these are modest laboratory defaults — size real
installs to the workload.

### spec.server.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.server.resources.limits.cpu

`string`

### spec.server.resources.limits.memory

`string`

### spec.server.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.server.resources.requests.cpu

`string`

### spec.server.resources.requests.memory

`string`

### spec.server.dataStorage

`KubernetesOpenBaoStorage`

The data volume (file storage in standalone, Raft storage in HA).
Ignored in dev mode (in-memory). One PVC per replica, mounted at
/openbao/data.

### spec.server.dataStorage.size

`string` · optional (explicit presence)

Volume size (e.g. "10Gi").

- default: `10Gi`
- rule: {"string":{"pattern":"^\\d+(\\.\\d+)?(Ki|Mi|Gi|Ti|Pi|Ei|K|M|G|T|P|E)$"}}

### spec.server.dataStorage.storageClass

`string | valueFrom`

StorageClass name. Empty uses the cluster's default class.
Accepts a literal name or a reference to a
KubernetesStorageClass resource.

- references: KubernetesStorageClass (`status.outputs.storage_class_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesStorageClass, name: <that resource's name>, fieldPath: status.outputs.storage_class_name}} -- a bare string does not parse

### spec.server.auditStorage

`KubernetesOpenBaoStorage`

Optional dedicated volume for file audit logs, mounted at
/openbao/audit. Creating the volume does NOT enable auditing —
after initialization run
`bao audit enable file file_path=/openbao/audit/audit.log`.

### spec.server.auditStorage.size

`string` · optional (explicit presence)

Volume size (e.g. "10Gi").

- default: `10Gi`
- rule: {"string":{"pattern":"^\\d+(\\.\\d+)?(Ki|Mi|Gi|Ti|Pi|Ei|K|M|G|T|P|E)$"}}

### spec.server.auditStorage.storageClass

`string | valueFrom`

StorageClass name. Empty uses the cluster's default class.
Accepts a literal name or a reference to a
KubernetesStorageClass resource.

- references: KubernetesStorageClass (`status.outputs.storage_class_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesStorageClass, name: <that resource's name>, fieldPath: status.outputs.storage_class_name}} -- a bare string does not parse

### spec.server.logLevel

`string` · optional (explicit presence)

Server log verbosity: trace, debug, info (default), warn, error.

- default: `info`
- rule: Log level must be one of: trace, debug, info, warn, error.

### spec.server.logFormat

`string` · optional (explicit presence)

Server log format: standard (default) or json.

- default: `standard`
- rule: Log format must be either "standard" or "json".

### spec.server.scheduling

`KubernetesOpenBaoScheduling`

Pod scheduling constraints for the server pods. NOTE the chart
ships a REQUIRED pod anti-affinity on hostname by default, so an
HA cluster needs as many schedulable nodes as replicas; relax it
through `helm_values` (`server.affinity: ""`) when running
multiple replicas on fewer nodes (labs only — co-located Raft
replicas share their node's fate).

### spec.server.scheduling.nodeSelector

`map<string, string>`

Schedule onto nodes carrying these labels.

### spec.server.scheduling.tolerations

`[]WorkloadToleration`

Tolerations for tainted nodes.

### spec.server.scheduling.tolerations[].key

`string`

Taint key to tolerate. Empty key with operator "Exists" tolerates every taint.

### spec.server.scheduling.tolerations[].operator

`string`

How key/value match: "Equal" (default — value must match too) or "Exists"
(key presence alone matches).

- rule: Toleration operator must be either "Equal" or "Exists"

### spec.server.scheduling.tolerations[].value

`string`

Taint value to match when operator is "Equal".

### spec.server.scheduling.tolerations[].effect

`string`

Which taint effect is tolerated: "NoSchedule", "PreferNoSchedule", or
"NoExecute". Empty tolerates all effects for the key.

- rule: Toleration effect must be one of "NoSchedule", "PreferNoSchedule", or "NoExecute"

### spec.server.scheduling.tolerations[].tolerationSeconds

`int64` · optional (explicit presence)

For "NoExecute" taints only: how many seconds already-running pods stay bound
after the taint appears. Unset means tolerate forever.

### spec.tls

`KubernetesOpenBaoTls`

End-to-end TLS for the OpenBao listener. When unset, the server
listens on plain HTTP inside the cluster (the chart default).
Enabling TLS is a COMPOSITE change this module owns end to end:
the listener gets tls_cert_file/tls_key_file, the certificate
Secret is mounted, and every derived URL and probe switches to
https. (The chart's own `global.tlsDisable` value alone does NOT
configure the listener — flipping it without listener changes
produces a plaintext server addressed as https, an instant
outage. This module renders all the pieces together.)

- rule: TLS needs certificate material: set cert_secret_name to a kubernetes.io/tls Secret (or reference a KubernetesCertificate) when TLS is enabled.

### spec.tls.enabled

`bool`

Enable TLS on the listener (port 8200) and the cluster port
(8201 always uses OpenBao's own cluster TLS regardless).

### spec.tls.certSecretName

`string | valueFrom`

Name of a kubernetes.io/tls Secret in the install namespace
carrying tls.crt / tls.key (and optionally ca.crt) for the
server. Accepts a literal name or a reference to a
KubernetesCertificate resource (cert-manager) — the natural
issuer: point the certificate's dnsNames at
`<name>.<namespace>.svc` and this component's derived DNS names.
Required when enabled.

- references: KubernetesCertificate (`status.outputs.secret_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesCertificate, name: <that resource's name>, fieldPath: status.outputs.secret_name}} -- a bare string does not parse

### spec.autoUnseal

`KubernetesOpenBaoAutoUnseal`

Auto-unseal: delegate master-key protection to an external KMS so
servers unseal themselves on startup (initialization is still a
one-time manual step; it produces RECOVERY keys instead of unseal
keys). Exactly one seal backend may be declared.

VERSION HORIZON (verified at OpenBao v2.6.1): the cloud KMS seal
mechanisms (awskms, gcpckms, azurekeyvault) are built in but
DEPRECATED — upstream moves them to external KMS plugins in
v2.7.0. This module renders the seal stanza for the pinned 2.6.x
line; expect the rendering to gain a plugin declaration when the
chart pin crosses 2.7.

### spec.autoUnseal.awsKms

`KubernetesOpenBaoAwsKmsSeal`

AWS KMS. Natural on EKS (IRSA for keyless auth).

### spec.autoUnseal.awsKms.region

`string` · required

AWS region of the KMS key (e.g. "us-west-2").

- rule: {"string":{"minLen":"1"}}

### spec.autoUnseal.awsKms.kmsKeyId

`string` · required

KMS key ID or full ARN of a SYMMETRIC encrypt/decrypt key.

- rule: {"string":{"minLen":"1"}}

### spec.autoUnseal.awsKms.accessKeyId

`string`

Static credentials — only when keyless (IRSA / instance profile)
is unavailable. The module materializes them into a Secret
(`<name>-seal-credentials`) delivered as environment variables;
nothing credential-bearing lands in the config ConfigMap.

### spec.autoUnseal.awsKms.secretAccessKey

`string` · sensitive

The secret access key paired with access_key_id.

### spec.autoUnseal.gcpKms

`KubernetesOpenBaoGcpKmsSeal`

GCP Cloud KMS. Natural on GKE (Workload Identity).

### spec.autoUnseal.gcpKms.project

`string | valueFrom` · required

GCP project containing the KMS key ring.

containment_exempt: names where the unseal key lives — the server
runs in the cluster, not the project.

- references: GcpProject (`status.outputs.project_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpProject, name: <that resource's name>, fieldPath: status.outputs.project_id}} -- a bare string does not parse

### spec.autoUnseal.gcpKms.region

`string` · required

KMS key ring region (e.g. "global", "us-central1").

- rule: {"string":{"minLen":"1"}}

### spec.autoUnseal.gcpKms.keyRing

`string | valueFrom` · required

Key ring name.

containment_exempt: an unseal-key source the server calls out to —
access, never placement.

- references: GcpKmsKeyRing (`status.outputs.key_ring_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpKmsKeyRing, name: <that resource's name>, fieldPath: status.outputs.key_ring_name}} -- a bare string does not parse

### spec.autoUnseal.gcpKms.cryptoKey

`string | valueFrom` · required

Crypto key (symmetric encrypt/decrypt) used to wrap the master
key. The identity running OpenBao needs
roles/cloudkms.cryptoKeyEncrypterDecrypter on it.

- references: GcpKmsKey (`status.outputs.key_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_name}} -- a bare string does not parse

### spec.autoUnseal.gcpKms.workloadIdentityServiceAccount

`string | valueFrom`

GKE Workload Identity: the GCP service account email to annotate
the server ServiceAccount with (iam.gke.io/gcp-service-account).
Leave empty to rely on node/ambient credentials. NOTE dev mode
drops ServiceAccount annotations (chart behavior) — auto-unseal
with workload identity requires standalone or ha mode.

- references: GcpServiceAccount (`status.outputs.email`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpServiceAccount, name: <that resource's name>, fieldPath: status.outputs.email}} -- a bare string does not parse

### spec.autoUnseal.azureKeyVault

`KubernetesOpenBaoAzureKeyVaultSeal`

Azure Key Vault. Natural on AKS (Workload Identity / MSI).

### spec.autoUnseal.azureKeyVault.vaultName

`string` · required

Key Vault name (the vault, not the key).

- rule: {"string":{"minLen":"1"}}

### spec.autoUnseal.azureKeyVault.keyName

`string` · required

Name of the key inside the vault.

- rule: {"string":{"minLen":"1"}}

### spec.autoUnseal.azureKeyVault.tenantId

`string` · required

Entra (Azure AD) tenant ID.

- rule: {"string":{"minLen":"1"}}

### spec.autoUnseal.azureKeyVault.clientId

`string`

Service-principal client ID — only when keyless (AKS Workload
Identity / Managed Identity) is unavailable.

### spec.autoUnseal.azureKeyVault.clientSecret

`string` · sensitive

Service-principal client secret paired with client_id. Delivered
as environment variables from a module-owned Secret; never lands
in the config ConfigMap.

### spec.autoUnseal.transit

`KubernetesOpenBaoTransitSeal`

Transit engine of another OpenBao/Vault instance.

### spec.autoUnseal.transit.address

`string` · required

Address of the central instance (e.g. "https://bao.example.com:8200").
The satellite depends on the central instance being reachable
and unsealed at every startup.

- rule: {"string":{"minLen":"1"}}

### spec.autoUnseal.transit.keyName

`string` · required

Transit key name used to wrap the master key.

- rule: {"string":{"minLen":"1"}}

### spec.autoUnseal.transit.mountPath

`string` · optional (explicit presence)

Transit engine mount path.

- default: `transit/`

### spec.autoUnseal.transit.token

`string` · sensitive

Token authorized for encrypt/decrypt on the transit key.
Delivered as environment variables from a module-owned Secret;
never lands in the config ConfigMap.

### spec.injector

`KubernetesOpenBaoInjector`

The OpenBao Agent Injector: a MutatingWebhookConfiguration that
intercepts pod creation CLUSTER-WIDE and injects secret-fetching
agent sidecars into annotated pods.

OFF by default here — a deliberate divergence from the chart
(whose default installs the webhook for every pod create/update
in the cluster). Enable it only when workloads will actually use
agent injection annotations. The webhook fails OPEN by default
(failure_policy Ignore), so injector downtime never blocks pod
creation — it silently skips injection instead.

### spec.injector.enabled

`bool`

Deploy the injector. See the spec-level comment for the
cluster-wide webhook blast radius.

### spec.injector.replicas

`int32` · optional (explicit presence)

Injector replicas. Above 1, leader election activates and the
chart creates a HARD-CODED Secret `openbao-injector-certs` —
only one multi-replica injector can exist per namespace.
The injector pods also carry a required anti-affinity, so
replicas need distinct nodes.

- default: `1`
- rule: {"int32":{"lte":5,"gte":1}}

### spec.injector.failurePolicy

`string` · optional (explicit presence)

Webhook failure policy: "Ignore" (default — injector downtime
skips injection, pods still schedule) or "Fail" (pod creation
BLOCKS while the injector is down; only for clusters that treat
missing secrets as worse than blocked deploys).

- default: `Ignore`
- rule: Webhook failure policy must be either "Ignore" (fail open) or "Fail" (fail closed).

### spec.injector.resources

`ContainerResources`

CPU and memory for the injector container.

### spec.injector.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.injector.resources.limits.cpu

`string`

### spec.injector.resources.limits.memory

`string`

### spec.injector.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.injector.resources.requests.cpu

`string`

### spec.injector.resources.requests.memory

`string`

### spec.uiEnabled

`bool` · optional (explicit presence)

Serve the built-in web UI and create the `<name>-ui` Service for
it. Defaults to true — the UI is part of the product experience.
(Exposure beyond the cluster composes from Gateway API kinds
referencing the exported service handles; this component never
creates ingress resources.)

- default: `true`

### spec.networkPolicyEnabled

`bool`

Render the chart's NetworkPolicy for the server pods (ingress on
8200/8201 from all namespaces by default). Off by default; most
clusters manage network policy through first-class
KubernetesNetworkPolicy resources instead.

### spec.metrics

`KubernetesOpenBaoMetrics`

Prometheus metrics. Enabling this renders the telemetry stanza
into the server config (prometheus_retention_time plus
unauthenticated_metrics_access on the listener — the metrics
endpoint is UNAUTHENTICATED when enabled, scoped to
/v1/sys/metrics) and optionally a ServiceMonitor.

- rule: A ServiceMonitor without the telemetry stanza scrapes an endpoint that rejects every request — enable metrics too.

### spec.metrics.enabled

`bool`

Render the telemetry stanza (prometheus_retention_time = "30s",
disable_hostname) and allow UNAUTHENTICATED access to
/v1/sys/metrics on the listener. Without this, the metrics
endpoint requires an OpenBao token and Prometheus cannot scrape.

### spec.metrics.serviceMonitorEnabled

`bool`

Also create a ServiceMonitor (requires the Prometheus Operator
CRDs — a KubernetesKubePrometheusStack — on the cluster; the
install FAILS without them). In HA mode the chart scrapes only
the active node.

### spec.snapshotAgent

`KubernetesOpenBaoSnapshotAgent`

Scheduled Raft snapshots shipped to an S3-compatible object store
(a CronJob running the openbao-snapshot-agent). This is the
disaster-recovery story for Raft-mode servers: snapshot files in
the bucket outlive the cluster and restore with
`bao operator raft snapshot restore`.

PREREQUISITE the module cannot create: the agent authenticates to
OpenBao through the Kubernetes auth method using `bao_role` —
that auth method and role are RUNTIME configuration you create
inside OpenBao after initialization (the docs carry the exact
recipe). Until the role exists the CronJob pods fail their login.

### spec.snapshotAgent.enabled

`bool`

Enable the snapshot CronJob. Meaningful only for Raft (HA) mode —
the agent calls `bao operator raft snapshot`.

### spec.snapshotAgent.schedule

`string` · optional (explicit presence)

Cron schedule.

- default: `*/15 * * * *`

### spec.snapshotAgent.s3Host

`string | valueFrom` · required

S3(-compatible) endpoint HOST (e.g. "s3.us-east-1.amazonaws.com",
or an in-cluster KubernetesSeaweedFs S3 endpoint host).

- references: KubernetesSeaweedFs (`status.outputs.s3_endpoint`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesSeaweedFs, name: <that resource's name>, fieldPath: status.outputs.s3_endpoint}} -- a bare string does not parse

### spec.snapshotAgent.s3Bucket

`string` · required

Bucket name to write snapshots into.

- rule: {"string":{"minLen":"1"}}

### spec.snapshotAgent.s3ExpireDays

`int32` · optional (explicit presence)

Days after which snapshots in the bucket expire (agent-side
cleanup).

- default: `14`
- rule: {"int32":{"gte":1}}

### spec.snapshotAgent.s3CredentialsSecretName

`string` · required

Name of an existing Secret with the s3cmd-style credentials the
agent image expects (keys access_key / secret_key). Required —
snapshot pods crash-loop without it.

- rule: {"string":{"minLen":"1"}}

### spec.snapshotAgent.baoRole

`string` · optional (explicit presence)

The OpenBao Kubernetes-auth ROLE the agent logs in with
(runtime configuration — create it inside OpenBao; see the
component docs for the recipe).

- default: `snapshot`

### spec.snapshotAgent.baoAuthPath

`string` · optional (explicit presence)

Kubernetes auth method mount path inside OpenBao.

- default: `kubernetes`

### spec.serviceAccount

`KubernetesOpenBaoServiceAccount`

Server ServiceAccount identity: cloud workload-identity
annotations and the Kubernetes-auth delegation binding.

### spec.serviceAccount.annotations

`map<string, string>`

Annotations for the server ServiceAccount — the cloud
workload-identity seam (eks.amazonaws.com/role-arn,
iam.gke.io/gcp-service-account, azure.workload.identity/client-id).
NOTE: dev mode drops these (chart behavior).

### spec.serviceAccount.authDelegatorEnabled

`bool` · optional (explicit presence)

Bind the ServiceAccount to the cluster's system:auth-delegator
role (a ClusterRoleBinding). Required for OpenBao's Kubernetes
AUTH METHOD to validate workload tokens via TokenReview — leave
on unless the cluster forbids the binding and you will not use
Kubernetes auth.

- default: `true`

### spec.helmValues

`string`

Advanced escape hatch: raw Helm values merged LAST (Helm `-f`
semantics) over everything this spec renders — later keys win.
Use it for the chart surfaces deliberately not modeled (CSI
provider, injector webhook selectors, extra volumes, affinity
overrides). The module re-pins `fullnameOverride` after the
merge, so resource naming cannot be overridden. YAML document as
a string.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesOpenBao, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace the server runs in. |
| `status.outputs.service` | `string` | The main client Service name (round-robins ALL server pods, including sealed/not-ready ones — by design, so init/unseal can reach them). |
| `status.outputs.internal_service` | `string` | The headless Service (`<name>-internal`) used for peer discovery and Raft cluster addresses. |
| `status.outputs.active_service` | `string` | The active-leader Service (`<name>-active`) — HA mode only, empty otherwise. Points at exactly the elected leader; the right target for write-heavy clients. |
| `status.outputs.ui_service` | `string` | The UI Service name (`<name>-ui`) when ui_enabled, empty otherwise. |
| `status.outputs.api_endpoint` | `string` | In-cluster API endpoint, scheme included (e.g. "http://bao.openbao.svc.cluster.local:8200" — https when TLS is enabled). What secret-consuming addons (external-secrets ClusterSecretStore, cert-manager Vault issuers) should point at. |
| `status.outputs.port` | `string` | API port (8200). |
| `status.outputs.service_account_name` | `string` | The server ServiceAccount name — the identity to bind cloud IAM (auto-unseal KMS access) and OpenBao Kubernetes-auth trust to. |
| `status.outputs.port_forward_command` | `string` | Copy-paste command for reaching the API from a workstation. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.server.dataStorage.storageClass` | KubernetesStorageClass | `status.outputs.storage_class_name` |
| `spec.server.auditStorage.storageClass` | KubernetesStorageClass | `status.outputs.storage_class_name` |
| `spec.tls.certSecretName` | KubernetesCertificate | `status.outputs.secret_name` |
| `spec.autoUnseal.gcpKms.project` | GcpProject | `status.outputs.project_id` |
| `spec.autoUnseal.gcpKms.keyRing` | GcpKmsKeyRing | `status.outputs.key_ring_name` |
| `spec.autoUnseal.gcpKms.cryptoKey` | GcpKmsKey | `status.outputs.key_name` |
| `spec.autoUnseal.gcpKms.workloadIdentityServiceAccount` | GcpServiceAccount | `status.outputs.email` |
| `spec.snapshotAgent.s3Host` | KubernetesSeaweedFs | `status.outputs.s3_endpoint` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
