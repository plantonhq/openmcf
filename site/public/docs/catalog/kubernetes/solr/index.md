---
title: "Solr"
description: "Solr deployment documentation"
icon: "package"
order: 100
componentName: "kubernetessolr"
---

# Kubernetes Solr

Declares an Apache SolrCloud cluster reconciled by the Apache Solr
Operator — node count and JVM sizing, ZooKeeper wiring (an
operator-provisioned ensemble or an external one), persistent or
ephemeral storage, basic-auth security bootstrap, keystore-based TLS,
registered backup repositories (S3, GCS, or a shared volume), and
shard-aware managed rolling updates. One resource per cluster;
workloads connect through the exported common service endpoint.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **SolrCloud** (`solr.apache.org/v1beta1`, named `metadata.name`) —
  the ONLY custom resource

The operator reconciles it into the node StatefulSet
(`<name>-solrcloud`), the common Service (`<name>-solrcloud-common`),
per-node PVCs (persistent storage), the provided ZooKeeper ensemble
(via the bundled zookeeper-operator), and — with basic auth — the
security bootstrap Secrets.

## Prerequisites

- The Apache Solr Operator on the cluster (KubernetesSolrOperator) —
  with its bundled zookeeper-operator installed when clusters use
  provided ensembles
- A StorageClass for persistent storage (most managed clusters provide
  a default; or reference a KubernetesStorageClass)
- For TLS: a PKCS#12 keystore Secret and a keystore-password Secret
  (e.g. a KubernetesCertificate with a pkcs12 keystore output)
- For `volume` backup repositories: an existing ReadWriteMany PVC

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesSolr
metadata:
  name: main
spec:
  namespace:
    value: search
  create_namespace: true
  version: "9.10.0"
  replicas: 3
  storage:
    persistent:
      size: 50Gi
  security:
    authentication_type: basic
```

The operator provisions a 3-node ZooKeeper ensemble, brings up three
Solr nodes on persistent volumes, and bootstraps basic-auth security.
Workloads connect at the exported `internal_endpoint`
(`http://main-solrcloud-common.search.svc.cluster.local`).

## Configuration

### Storage

The operator default is EPHEMERAL — data is lost when a pod leaves its
node. Declare `storage.persistent` (size, StorageClass, reclaim policy
defaulting to Retain) for real workloads; keep ephemeral for throwaway
experiments only.

### Security

`authentication_type: basic` bootstraps security.json with `admin`,
`solr`, and `k8s-oper` users carrying random passwords. The admin and
solr credentials land in `<name>-solrcloud-security-bootstrap`; the
operator's own credentials land in `<name>-solrcloud-basic-auth` (the
exported handle). If you rotate a password through Solr's security API
later, update the Secret too — the operator locks itself out
otherwise. Empty security means an open cluster: development only.

### Backups

`backup_repositories` registers named targets: S3 (declared keys, or
keyless via IRSA on EKS), GCS (a service-account key, or keyless via
GKE Workload Identity), or a shared ReadWriteMany volume. Backup runs
are operational verbs — SolrBackup resources or Solr API calls —
against these names.

### Updates and scaling

The Managed strategy (default) restarts pods shard-aware, bounding
unavailable pods (25%) and unavailable shard replicas (1); scale
operations move replicas off departing pods and onto new ones by
default.

### Exposure

The `external` block models the operator's own Ingress/ExternalDNS
exposure with per-node addressability (what CloudSolrClient needs from
outside the cluster); for simple HTTP access, compose a
KubernetesIngress over the common service instead.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster runs in |
| `cluster_name` | SolrCloud resource name (= `metadata.name`) |
| `common_service_name` | Common Service fronting all nodes (`<name>-solrcloud-common`) |
| `internal_endpoint` | In-cluster base URL (http on 80, https on 443 with TLS) |
| `basic_auth_secret_name` | Operator-generated basic-auth Secret (`<name>-solrcloud-basic-auth`); empty when security is off or user-provided |
| `zookeeper_connection_string` | The ensemble the cluster uses (host:port, plus chroot) |
| `port_forward_command` | Workstation access when no exposure is composed |

## Related Components

- [KubernetesSolrOperator](/docs/catalog/kubernetes/solr-operator)
  — the engine; must be installed first
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) —
  provides the target namespace via reference
- [KubernetesCertificate](/docs/catalog/kubernetes/kubernetescertificate)
  — produces the PKCS#12 keystore Secret the TLS block references
- [KubernetesIngress](/docs/catalog/kubernetes/ingress) —
  composes external exposure over the common service

## Next Steps

Turn on basic auth before anything real touches the cluster, move to
persistent storage, and size `java_mem` together with container
resources (heap at roughly half the container memory). Register a
backup repository and schedule backup runs against it. When external
clients need to address individual nodes, declare the `external` block
— otherwise compose exposure over the common service.
