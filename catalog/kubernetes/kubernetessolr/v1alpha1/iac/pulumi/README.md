# KubernetesSolr Pulumi Module

Deploys one Apache Solr Operator-managed SolrCloud cluster: the optional
namespace and the `solr.apache.org/v1beta1` SolrCloud resource. The custom
resource renders through typed crd2pulumi SDK bindings pinned to the Solr
Operator CRDs — field or structure drift against the pinned CRD fails at
COMPILE time, not at apply time.

Prerequisites at deploy time: the Apache Solr Operator
(`KubernetesSolrOperator`) on the cluster — with its bundled
zookeeper-operator when the ZooKeeper block is empty or `provided` (the
default posture).

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance labels
   when `create_namespace` is true
2. **SolrCloud** — the cluster itself; the operator creates everything
   else (node StatefulSet, provided ZooKeeper ensemble, services, PVCs,
   the basic-auth bootstrap Secret, Ingress exposure) from this one
   resource. Unset optionals are omitted entirely so the apiserver
   applies the CRD's own defaults.

## The Operator's Naming Contract

Every object the operator creates derives from `metadata.name` — the
module computes these in `locals.go` and exports them (never read back
from the cluster):

| Object | Name |
|---|---|
| Node StatefulSet | `<name>-solrcloud` |
| Common Service (all nodes) | `<name>-solrcloud-common` |
| Headless Service | `<name>-solrcloud-headless` |
| Generated basic-auth Secret | `<name>-solrcloud-basic-auth` |
| Provided ZooKeeper client service | `<name>-solrcloud-zookeeper-client:2181` |

## CR Mapping (spec → SolrCloud)

| Spec field | SolrCloud field |
|---|---|
| `replicas` (default 3) | `spec.replicas` |
| `version` / `image_repository` (default `solr`) | `spec.solrImage{tag, repository}` (pull policy omitted) |
| `java_mem` / `solr_opts` / `log_level` (default INFO) / `gc_tune` | `spec.solrJavaMem` / `solrOpts` / `solrLogLevel` / `solrGCTune` |
| `zookeeper.provided{replicas, chroot, persistence, resources}` | `spec.zookeeperRef.provided{replicas, chroot, persistence.spec{resources.requests.storage, storageClassName}, zookeeperPodPolicy.resources}` |
| `zookeeper.external{connection_string, chroot}` | `spec.zookeeperRef.connectionInfo{internalConnectionString, chroot}` |
| `zookeeper` empty | NOTHING rendered — the operator defaults to a provided 3-node ensemble |
| `storage.persistent{size, storage_class, reclaim_policy}` | `spec.dataStorage.persistent{reclaimPolicy, pvcTemplate.spec{resources.requests.storage, storageClassName}}` |
| `storage.ephemeral{size_limit}` | `spec.dataStorage.ephemeral.emptyDir{sizeLimit}` (`emptyDir: {}` when no limit) |
| `resources` / `node_selector` / `tolerations` | `spec.customSolrKubeOptions.podOptions{resources, nodeSelector, tolerations}` — only when any is set |
| `pod_port` (default 8983) / `external{...}` | `spec.solrAddressability{podPort, external{method, domainName, useExternalAddress, hideCommon, hideNodes}}` |
| `update_strategy{method, max_pods_unavailable, max_shard_replicas_unavailable, restart_schedule}` | `spec.updateStrategy{method, managed{maxPodsUnavailable, maxShardReplicasUnavailable}, restartSchedule}` — the managed budgets keep the CRD's int-or-string semantics (`"2"` renders as the number 2, `"25%"` stays a string) |
| `availability.pdb_enabled` | `spec.availability.podDisruptionBudget.enabled` — rendered only when EXPLICITLY set (absence already means enabled upstream) |
| `scaling{vacate_pods_on_scale_down, populate_pods_on_scale_up}` | `spec.scaling{vacatePodsOnScaleDown, populatePodsOnScaleUp}` — rendered only when explicitly set |
| `tls{pkcs12_secret, keystore_password_secret, truststore_secret, truststore_password_secret, client_auth (default None), verify_client_hostname}` | `spec.solrTLS{pkcs12Secret, keyStorePasswordSecret, trustStoreSecret, trustStorePasswordSecret, clientAuth, verifyClientHostname}` |
| `security.authentication_type: basic` | `spec.solrSecurity{authenticationType: "Basic"` (CRD value is capitalized)`, basicAuthSecret, probesRequireAuth, bootstrapSecurityJson}` — nothing renders when basic auth is not enabled |
| `backup_repositories[]` | `spec.backupRepositories[]{name}` + one arm: `s3{region, bucket, baseLocation, endpoint, credentials{accessKeyIdSecret, secretAccessKeySecret}}` \| `gcs{bucket, gcsCredentialSecret, baseLocation}` \| `volume{source.persistentVolumeClaim.claimName, directory}` |
| `solr_modules` / `additional_libs` | `spec.solrModules` / `additionalLibs` |

## Usage

```shell
planton pulumi up --manifest ../hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster runs in |
| `cluster_name` | Name of the SolrCloud resource (equals `metadata.name`) |
| `common_service_name` | Common Service fronting all nodes (`<name>-solrcloud-common`) |
| `internal_endpoint` | In-cluster base URL through the common service — `http://…` (port 80) without TLS, `https://…` (port 443) with TLS |
| `basic_auth_secret_name` | Operator-generated credential Secret (`<name>-solrcloud-basic-auth`) — empty when security is disabled or a user-provided `basic_auth_secret` is in play |
| `zookeeper_connection_string` | What the cluster uses: the external connection string, or `<name>-solrcloud-zookeeper-client:2181`, plus the chroot when it diverges from `/` |
| `port_forward_command` | Port-forward command for workstation access (8983 → 80, or 443 with TLS) |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → SolrCloud → output exports
- `module/locals.go`: the operator's naming contract (common service,
  basic-auth secret, provided-ZooKeeper client service) — kept in
  lockstep with the Terraform module's `locals.tf`
- `module/solr_cloud.go`: the SolrCloud resource (image, ZooKeeper
  wiring, storage, pod options, addressability, update strategy,
  availability, scaling, TLS, security, backup repositories)
- `module/namespace.go`: the optional namespace
