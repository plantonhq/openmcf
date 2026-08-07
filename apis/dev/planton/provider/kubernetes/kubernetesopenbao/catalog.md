# OpenBao

Deploys OpenBao -- the open-source secrets manager forked from HashiCorp Vault -- from the official `openbao` chart at `openbao.github.io/openbao-helm`. Three server modes cover the whole lifecycle: dev (in-memory, auto-unsealed, evaluation only), standalone (one server, file storage on a persistent volume -- the chart default), and HA (integrated Raft storage with leader election; the module synthesizes the `retry_join` stanzas the chart alone ships without, so multi-replica clusters actually form). Auto-unseal delegates master-key protection to AWS KMS, GCP Cloud KMS, Azure Key Vault, or a central instance's transit engine. Uses a Kubernetes Provider Connection for cluster access.

Know the seal lifecycle before you deploy: a fresh OpenBao server starts UNINITIALIZED and SEALED, and reports NotReady BY DESIGN until you run `bao operator init` and unseal it -- initialization is a runtime operation no deployment tool can perform declaratively. The readiness probe is `bao status`, and the chart keeps sealed pods addressable so init and unseal can reach them.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kubernetes Namespace** -- created only when `createNamespace` is `true`; otherwise deploys into an existing namespace
- **Helm Release** -- the `openbao` chart, creating:
  - StatefulSet for the server (one replica in standalone; `server.ha.replicas` Raft peers in HA, each with its own data PVC at `/openbao/data`; dev mode runs in-memory with NO PVC)
  - The client Service (round-robins ALL server pods including sealed ones -- by design), the headless `-internal` Service for Raft peer discovery, an `-active` Service pointing at the elected leader (HA only), and a `-ui` Service when the UI is on
  - ConfigMap carrying the rendered HCL server configuration -- listener, storage backend, synthesized `retry_join` stanzas, and the optional seal and telemetry stanzas
  - ServiceAccount with your workload-identity annotations and, unless disabled, the `system:auth-delegator` ClusterRoleBinding OpenBao's Kubernetes auth method needs for TokenReview
- **Optional audit volume** -- a second PVC at `/openbao/audit` when `server.auditStorage` is declared; auditing itself is enabled at runtime (`bao audit enable file ...`) after initialization
- **Seal credentials Secret** -- created only when an auto-unseal arm declares static credentials; delivered as environment variables, never written into the config ConfigMap
- **Agent Injector** -- created only when `injector.enabled` is `true` (OFF by default here -- a deliberate divergence from the chart's cluster-wide-webhook default); a mutating webhook that injects agent sidecars into annotated pods
- **Snapshot CronJob** -- created only when `snapshotAgent.enabled` is `true`; ships `bao operator raft snapshot` output to an S3-compatible bucket on a schedule
- **ServiceMonitor** -- created only when `metrics.serviceMonitorEnabled` is `true`; requires the Prometheus Operator CRDs on the cluster
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target Kubernetes cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **A storage class** for the data (and optional audit) persistent volumes -- standalone and HA modes; dev mode is in-memory and needs none.
- **As many schedulable nodes as HA replicas** -- the chart ships a REQUIRED hostname anti-affinity, so a three-replica Raft cluster needs three nodes (relaxable through `helmValues` in labs only).
- **A cloud KMS key or a central transit engine** (only when using auto-unseal) -- e.g. a GCP Cloud KMS crypto key with `roles/cloudkms.cryptoKeyEncrypterDecrypter` granted to the identity OpenBao runs as. ValueFromRef can resolve the project, key ring, and crypto key from other Cloud Resources.
- **Prometheus Operator CRDs** (only when enabling the ServiceMonitor) -- the install FAILS without them.

## Deploy

### Console

Open the deployment store, find **OpenBao**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Dev Mode** preset for evaluation, **Production HA** for a Raft cluster with manual unsealing, or **Production HA GCP Auto-Unseal** for the restart-toil-free shape, in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesOpenBao
metadata:
  name: openbao
  org: acme-corp
  env: prod
spec:
  namespace:
    value: "openbao"
  create_namespace: true
  server:
    ha:
      replicas: 3
    resources:
      requests:
        cpu: 250m
        memory: 512Mi
      limits:
        cpu: "1"
        memory: 1Gi
    data_storage:
      size: 10Gi
    audit_storage:
      size: 10Gi
  metrics:
    enabled: true
```

```shell
planton apply -f openbao.yaml
```

This deploys a three-node Raft cluster where each replica persists to its own 10Gi PVC and an audit volume waits at `/openbao/audit`. The bootstrap is yours by design: initialize once through pod 0 (`bao operator init` -- custody of the unseal key shares and root token is the whole point of a secrets manager), then unseal every pod; peers join automatically through the synthesized `retry_join`. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire OpenBao to dependencies managed by other Cloud Resources:

```yaml
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: openbao-namespace
      fieldPath: spec.name
  create_namespace: false
  auto_unseal:
    gcp_kms:
      project:
        valueFrom:
          kind: GcpProject
          name: infra-project
          fieldPath: status.outputs.project_id
      region: global
      key_ring:
        valueFrom:
          kind: GcpKmsKeyRing
          name: openbao-keyring
          fieldPath: status.outputs.key_ring_name
      crypto_key:
        valueFrom:
          kind: GcpKmsKey
          name: openbao-unseal-key
          fieldPath: status.outputs.key_name
      workload_identity_service_account:
        valueFrom:
          kind: GcpServiceAccount
          name: openbao-sa
          fieldPath: status.outputs.email
```

The InfraPipeline resolves the dependency graph, deploys the GCP project, KMS key ring, KMS key, and service account first, then provisions OpenBao with the resolved values. The KMS references name where the unseal key LIVES -- access, never placement: the server runs in this cluster, not the GCP project.

## Key Configuration

These are the most important decisions when configuring OpenBao. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Server mode** -- `server.dev`, `server.standalone`, and `server.ha` are one choice (leave all unset for the chart's standalone default). Dev mode is evaluation only: in-memory data lost on every restart, auto-initialized and auto-unsealed, the root token literally `root` in plain text in the pod spec, no PVC -- and the chart DROPS ServiceAccount annotations, so cloud workload identity does not apply there. HA mode (`server.ha.replicas`, default 3, range 1-11) runs integrated Raft: odd counts tolerate minority loss, and the module synthesizes `retry_join` for every peer -- without it a multi-replica install never forms a cluster.

**Auto-unseal** -- By default every restarted pod waits SEALED for a human with unseal key shares (Shamir-mode reality). Declaring one `auto_unseal` arm (`aws_kms`, `gcp_kms`, `azure_key_vault`, or `transit`) wraps the master key with an external KMS so servers unseal THEMSELVES at startup. Initialization stays a one-time manual step -- with auto-unseal it produces RECOVERY keys instead of unseal keys. Keyless-first: prefer ambient workload identity (IRSA / GKE Workload Identity / Azure MSI) and leave the static-credential fields empty; declared credentials are org-secret references materialized into a module-owned Secret and delivered as environment variables. Version horizon: at the pinned OpenBao 2.6.x the cloud KMS seals are built in but deprecated upstream -- v2.7 moves them to external KMS plugins.

**TLS is a composite this module owns end to end** -- `tls.enabled` with `tls.cert_secret_name` switches the listener's certificate files, the Secret mount, every derived URL, and the probe scheme TOGETHER. (The chart's lone `global.tlsDisable` flag alone produces a plaintext server addressed as https -- an instant outage; the module renders all the pieces coherently.) A KubernetesCertificate reference is the natural issuer.

**Storage split** -- `server.data_storage` (default 10Gi; one PVC per replica) holds file storage in standalone and Raft data in HA; dev mode ignores it. `server.audit_storage` optionally mounts a second volume at `/openbao/audit` -- creating the volume does NOT enable auditing; run `bao audit enable file file_path=/openbao/audit/audit.log` after initialization.

**Agent Injector** -- OFF by default here, a deliberate divergence from the chart (whose default installs the MutatingWebhookConfiguration for every pod create/update cluster-wide). When on, `injector.failure_policy` chooses `Ignore` (fail open -- injector downtime skips injection; the default) or `Fail`; above 1 replica, leader election creates the hard-coded `openbao-injector-certs` Secret -- one multi-replica injector per namespace.

**Raft snapshots** -- `snapshot_agent` runs a CronJob (default `*/15 * * * *`) shipping `bao operator raft snapshot` to an S3-compatible bucket (an in-cluster KubernetesSeaweedFs endpoint composes naturally); `s3_expire_days` (default 14) prunes old snapshots agent-side. PREREQUISITE the module cannot create: the Kubernetes auth method and the `bao_role` it logs in with are runtime configuration inside OpenBao -- until they exist, snapshot pods fail their login. `s3_credentials_secret_name` names an existing Secret with s3cmd-style `access_key` / `secret_key` -- pods crash-loop without it.

**Metrics are unauthenticated when enabled** -- `metrics.enabled` renders the telemetry stanza AND opens `/v1/sys/metrics` without a token; anything that can reach the Service can read operational telemetry. `metrics.service_monitor_enabled` additionally requires the Prometheus Operator CRDs, and in HA scrapes only the active node.

**Workload identity and Kubernetes auth** -- `service_account.annotations` is the cloud identity seam (`eks.amazonaws.com/role-arn`, `iam.gke.io/gcp-service-account`, `azure.workload.identity/client-id`) -- the keyless path for auto-unseal KMS access. `service_account.auth_delegator_enabled` (default true) binds `system:auth-delegator`, which OpenBao's Kubernetes auth method needs to validate workload tokens via TokenReview.

**The escape hatch** -- `helm_values` merges LAST over everything the typed fields render (Helm `-f` semantics): the CSI provider, injector webhook selectors, extra volumes, the lab-only anti-affinity relaxation. The module re-pins `fullnameOverride` after the merge so the exported Service names stay stable. Never put secret material in this document.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **KubernetesNamespace** | `namespace` | `spec.name` |
| **KubernetesStorageClass** (optional) | `server.data_storage.storage_class` / `server.audit_storage.storage_class` | `status.outputs.storage_class_name` |
| **KubernetesCertificate** (optional) | `tls.cert_secret_name` | `status.outputs.secret_name` |
| **GcpProject** (optional) | `auto_unseal.gcp_kms.project` | `status.outputs.project_id` |
| **GcpKmsKeyRing** (optional) | `auto_unseal.gcp_kms.key_ring` | `status.outputs.key_ring_name` |
| **GcpKmsKey** (optional) | `auto_unseal.gcp_kms.crypto_key` | `status.outputs.key_name` |
| **GcpServiceAccount** (optional) | `auto_unseal.gcp_kms.workload_identity_service_account` | `status.outputs.email` |
| **KubernetesSeaweedFs** (optional) | `snapshot_agent.s3_host` | `status.outputs.s3_endpoint` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace the server runs in | Locating the install for diagnostics |
| `service` | The main client Service (round-robins ALL pods, sealed included -- by design, so init/unseal can reach them) | General client traffic |
| `internal_service` | The headless `-internal` Service | Raft peer discovery and cluster addresses |
| `active_service` | The `-active` Service pointing at the elected leader (HA only, empty otherwise) | Write-heavy clients |
| `ui_service` | The `-ui` Service (empty when the UI is disabled) | Exposing the web UI via Gateway API kinds |
| `api_endpoint` | In-cluster API endpoint, scheme included (https when TLS is on) | external-secrets ClusterSecretStore, cert-manager Vault issuers |
| `port` | API port (8200) | Connection configuration |
| `service_account_name` | The server ServiceAccount | Binding cloud IAM (auto-unseal KMS access) and Kubernetes-auth trust |
| `port_forward_command` | Ready-to-run `kubectl port-forward` command | Workstation access |

Deliberately absent: root tokens and unseal/recovery keys are produced by the runtime `bao operator init` call and are never known to (or owned by) the deployment -- capture them from the init output and store them outside the cluster.

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Dev mode** -- Zero ceremony: auto-initialized, auto-unsealed, root token `root`, all data in memory. Evaluation and API integration work only -- never real secrets. Start from the **Dev Mode** preset.

**Production HA (Raft)** -- Three servers with integrated Raft storage, per-replica data PVCs, an audit volume, and metrics on. Initialization and unsealing are yours; after every pod restart the affected server waits sealed. Start from the **Production HA** preset.

**Production HA + GCP auto-unseal** -- The HA shape with the restart toil removed: the master key wrapped by a Cloud KMS crypto key via GKE Workload Identity -- no static credential anywhere. Start from the **Production HA GCP Auto-Unseal** preset.

## Works With

- [**Kubernetes Namespace**](/cloud-catalog/kubernetes-namespace) -- provides the namespace for the OpenBao install
- [**Kubernetes Storage Class**](/cloud-catalog/kubernetes-storage-class) -- backs the data and audit persistent volumes
- [**Kubernetes Certificate**](/cloud-catalog/kubernetes-certificate) -- issues the server TLS certificate Secret
- [**GCP KMS Key**](/cloud-catalog/gcp-kms-key) -- wraps the master key for GCP Cloud KMS auto-unseal
- [**GCP Service Account**](/cloud-catalog/gcp-service-account) -- the workload identity for keyless KMS access
- [**Kubernetes SeaweedFS**](/cloud-catalog/kubernetes-seaweed-fs) -- an in-cluster S3 endpoint for Raft snapshots
- [**Kubernetes External Secrets**](/cloud-catalog/kubernetes-external-secrets) -- consumes the `api_endpoint` output as a ClusterSecretStore backend
