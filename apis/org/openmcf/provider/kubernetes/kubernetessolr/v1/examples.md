# SolrKubernetes API-Resource Examples

## Example 1: Basic Solr Deployment with Default Settings

This example demonstrates the most basic configuration for deploying a Solr instance within a Kubernetes cluster. It configures a single Solr pod with default resource allocations and storage size. The Zookeeper instance is also deployed to support Solr.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SolrKubernetes
metadata:
  name: solr-instance-basic
spec:
  kubernetesProviderConfigId: cluster-credential-12345
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
  kubernetesProviderConfigId: cluster-credential-67890
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
  kubernetesProviderConfigId: cluster-credential-24680
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
    isEnabled: true
    ingressClass: "nginx"
    host: "solr.example.com"
    tlsEnabled: true
```

## Example 4: Solr Deployment with Custom Garbage Collection and No Ingress

This configuration deploys a Solr cluster with custom garbage collection tuning but without ingress, relying instead on internal Kubernetes networking and port-forwarding for access.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SolrKubernetes
metadata:
  name: solr-instance-no-ingress
spec:
  kubernetesProviderConfigId: cluster-credential-112233
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

## Usage

Refer to the example section for usage instructions.
