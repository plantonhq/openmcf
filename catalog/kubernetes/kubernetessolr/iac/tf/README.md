# KubernetesSolr Terraform Module

Deploys one Apache Solr Operator-managed SolrCloud cluster: the optional
namespace and the `solr.apache.org/v1beta1` SolrCloud CR. The custom
resource applies through `kubectl_manifest` (alekc/kubectl provider,
server-side apply), which needs no cluster connection at plan time — the
cluster can be planned before the Solr operator's CRDs exist.

Prerequisites at apply time: the Apache Solr Operator
(`KubernetesSolrOperator`) on the cluster — with its bundled
zookeeper-operator when the ZooKeeper block is empty or `provided` (the
default posture).

## Module Behavior

- **The naming contract flows from `metadata.name`** — every object the
  operator creates derives from it, and the module derives the same names
  blind (never read back from the cluster):

  | Object | Name |
  |---|---|
  | Node StatefulSet | `<name>-solrcloud` |
  | Common Service (all nodes) | `<name>-solrcloud-common` |
  | Headless Service | `<name>-solrcloud-headless` |
  | Generated basic-auth Secret | `<name>-solrcloud-basic-auth` |
  | Provided ZooKeeper client service | `<name>-solrcloud-zookeeper-client:2181` |

- **One CR, nothing else** — the operator owns the StatefulSet, the
  provided ZooKeeper ensemble, services, PVCs, the security bootstrap,
  and (when the `external` block is declared) its own Ingress /
  ExternalDNS exposure. The module renders no Ingress resources of its
  own; composing a KubernetesIngress over the common service handle is
  the alternative exposure path.
- **Unset optionals are omitted** — the rendered CR body is null-pruned
  so the apiserver applies the CRD's own defaults; presence-sensitive
  fields (`availability.pdb_enabled`, the two `scaling` flags) render
  only when explicitly set, and plain booleans render only when true.
- **An empty `zookeeper` block renders NOTHING** — the operator then
  defaults to a provided 3-node ensemble by itself.
- **Int-or-string budgets** — the managed update-strategy budgets keep
  the CRD's int-or-string semantics with `try(tonumber(x), x)`: `"2"`
  reaches the API as the number 2, `"25%"` stays a string.
- **The CRD's authenticationType is capitalized** — the spec's
  `security.authentication_type: basic` renders as `Basic`; nothing
  renders when basic auth is not enabled.
- **The module (not the operator) owns namespace creation** —
  `create_namespace` drives a `kubernetes_namespace_v1` resource carrying
  the standard governance labels.

## CR Mapping (spec → SolrCloud)

| Spec field | SolrCloud field |
|---|---|
| `replicas` (default 3) | `spec.replicas` |
| `version` / `image_repository` (default `solr`) | `spec.solrImage{tag, repository}` (pull policy omitted) |
| `java_mem` / `solr_opts` / `log_level` (default INFO) / `gc_tune` | `spec.solrJavaMem` / `solrOpts` / `solrLogLevel` / `solrGCTune` |
| `zookeeper.provided{replicas, chroot, persistence, resources}` | `spec.zookeeperRef.provided{replicas, chroot, persistence.spec{resources.requests.storage, storageClassName}, zookeeperPodPolicy.resources}` |
| `zookeeper.external{connection_string, chroot}` | `spec.zookeeperRef.connectionInfo{internalConnectionString, chroot}` |
| `storage.persistent{size, storage_class, reclaim_policy}` | `spec.dataStorage.persistent{reclaimPolicy, pvcTemplate.spec{resources.requests.storage, storageClassName}}` |
| `storage.ephemeral{size_limit}` | `spec.dataStorage.ephemeral.emptyDir{sizeLimit}` (`emptyDir: {}` when no limit) |
| `resources` / `node_selector` / `tolerations` | `spec.customSolrKubeOptions.podOptions{resources, nodeSelector, tolerations}` — only when any is set |
| `pod_port` (default 8983) / `external{...}` | `spec.solrAddressability{podPort, external{method, domainName, useExternalAddress, hideCommon, hideNodes}}` |
| `update_strategy{...}` | `spec.updateStrategy{method, managed{maxPodsUnavailable, maxShardReplicasUnavailable}, restartSchedule}` |
| `availability.pdb_enabled` | `spec.availability.podDisruptionBudget.enabled` (explicit values only) |
| `scaling{...}` | `spec.scaling{vacatePodsOnScaleDown, populatePodsOnScaleUp}` (explicit values only) |
| `tls{...}` | `spec.solrTLS{pkcs12Secret, keyStorePasswordSecret, trustStoreSecret, trustStorePasswordSecret, clientAuth, verifyClientHostname}` |
| `security{...}` | `spec.solrSecurity{authenticationType: "Basic", basicAuthSecret, probesRequireAuth, bootstrapSecurityJson}` |
| `backup_repositories[]` | `spec.backupRepositories[]{name}` + one arm: `s3{region, bucket, baseLocation, endpoint, credentials{accessKeyIdSecret, secretAccessKeySecret}}` \| `gcs{bucket, gcsCredentialSecret, baseLocation}` \| `volume{source.persistentVolumeClaim.claimName, directory}` |
| `solr_modules` / `additional_libs` | `spec.solrModules` / `additionalLibs` |

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.namespace` | `spec.create_namespace` |
| `kubectl_manifest.solr_cloud` | always |

## Usage

```bash
planton tofu apply --manifest kubernetes-solr.yaml
```

## Local Development

```bash
tofu init -backend=false
tofu validate
tofu plan -var-file=terraform.tfvars.json
```

The full-surface hack manifest lives in `../../e2e/manifest.yaml`.

## Inputs

See `variables.tf` for the full variable specification (generated from
the spec proto). The spec arrives from the proto→tfvars converter in
snake_case with every `StringValueOrRef` foreign key — `namespace`
(KubernetesNamespace), the two `storage_class` references
(KubernetesStorageClass), and `security.basic_auth_secret`
(KubernetesSecret) — resolved to a literal string before Terraform runs.

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Namespace the cluster runs in |
| `cluster_name` | Name of the SolrCloud resource (equals `metadata.name`) |
| `common_service_name` | Common Service fronting all nodes (`<name>-solrcloud-common`) |
| `internal_endpoint` | In-cluster base URL — `http://…` (80) without TLS, `https://…` (443) with TLS |
| `basic_auth_secret_name` | Operator-generated credential Secret — empty when security is disabled or a user-provided secret is in play |
| `zookeeper_connection_string` | Connection string the cluster uses (host:port/chroot) |
| `port_forward_command` | Port-forward command for workstation access |

## Parity

This module is the behavioral twin of the Pulumi module
(`../pulumi/module/`): same rendered CR body (null-pruned,
presence-sensitive fields on explicit values only, int-or-string budgets),
same derived operator names, same outputs — keep them in lockstep.
