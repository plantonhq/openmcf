# SolrKubernetes API-Resource Examples

## Example 1: Basic Solr Deployment with Default Settings

This example demonstrates the most basic configuration for deploying a Solr instance within a Kubernetes cluster. It configures a single Solr pod with default resource allocations and storage size. The Zookeeper instance is also deployed to support Solr.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SolrKubernetes
metadata:
  name: solr-instance-basic
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: solr-instance-basic
  createNamespace: true
  solrContainer:
    replicas: 1
    image:
      repo: solr
      tag: 8.7.0
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    diskSize: "1Gi"
  zookeeperContainer:
    replicas: 1
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    diskSize: "1Gi"
```

## Example 2: Solr Deployment with Custom JVM and Persistent Storage

This example configures a Solr deployment with custom JVM memory settings and a larger persistent volume attached to each Solr pod. Additionally, the Zookeeper container is customized with specific resource limits and a larger storage size.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SolrKubernetes
metadata:
  name: solr-instance-custom
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: solr-instance-custom
  createNamespace: true
  solrContainer:
    replicas: 3
    image:
      repo: solr
      tag: 8.7.0
    resources:
      requests:
        cpu: 100m
        memory: 512Mi
      limits:
        cpu: 2
        memory: 2Gi
    diskSize: "5Gi"
    config:
      javaMem: "-Xms2g -Xmx4g"
      opts: "-Dsolr.autoSoftCommit.maxTime=5000"
      garbageCollectionTuning: "-XX:SurvivorRatio=6 -XX:MaxTenuringThreshold=10"
  zookeeperContainer:
    replicas: 3
    resources:
      requests:
        cpu: 100m
        memory: 512Mi
      limits:
        cpu: 2
        memory: 2Gi
    diskSize: "5Gi"
```

## Example 3: Solr Deployment with Ingress Enabled

This example demonstrates how to deploy Solr with ingress enabled. This allows external access to the Solr instance through a Kubernetes ingress resource, useful for exposing Solr services to external clients.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SolrKubernetes
metadata:
  name: solr-instance-ingress
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: solr-instance-ingress
  createNamespace: true
  solrContainer:
    replicas: 2
    image:
      repo: solr
      tag: 8.7.0
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    diskSize: "2Gi"
  zookeeperContainer:
    replicas: 2
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    diskSize: "2Gi"
  ingress:
    enabled: true
    hostname: "solr.example.com"
```

## Example 4: Solr Deployment with Custom Garbage Collection and No Ingress

This configuration deploys a Solr cluster with custom garbage collection tuning but without ingress, relying instead on internal Kubernetes networking and port-forwarding for access.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SolrKubernetes
metadata:
  name: solr-instance-no-ingress
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: solr-instance-no-ingress
  createNamespace: true
  solrContainer:
    replicas: 1
    image:
      repo: solr
      tag: 8.7.0
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    diskSize: "1Gi"
    config:
      garbageCollectionTuning: "-XX:SurvivorRatio=4 -XX:TargetSurvivorRatio=85 -XX:MaxTenuringThreshold=6"
  zookeeperContainer:
    replicas: 1
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    diskSize: "1Gi"
```

## Example 5: Using Existing Namespace

This example shows how to deploy Solr into an existing namespace that's managed separately. This is useful when multiple components share a namespace or when namespace management is centralized.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SolrKubernetes
metadata:
  name: solr-shared-namespace
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: shared-services
  createNamespace: false  # Don't create namespace, use existing one
  solrContainer:
    replicas: 1
    image:
      repo: solr
      tag: 8.7.0
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    diskSize: "1Gi"
  zookeeperContainer:
    replicas: 1
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    diskSize: "1Gi"
```

## Namespace Management

The Solr Kubernetes component provides flexible namespace management through the `createNamespace` flag:

- **createNamespace: true** (default) - The module creates the namespace with proper labels and manages its lifecycle. The namespace will be created by the Pulumi module and all child resources will depend on it.

- **createNamespace: false** - The module uses an existing namespace. You're responsible for ensuring the namespace exists before deploying Solr. This is useful when:
  - Multiple components share a namespace
  - Namespace is managed by a separate infrastructure module
  - Organization policies require centralized namespace management

When using an existing namespace (createNamespace: false), ensure the namespace exists before running `pulumi up`, otherwise the deployment will fail.

## Usage

Refer to the example section for usage instructions.
