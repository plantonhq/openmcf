# NATS on Kubernetes Examples

Below are practical examples demonstrating how to use the `NatsKubernetes` component within your OpenMCF
deployments. These examples illustrate common configurations and scenarios to help you quickly integrate a robust NATS
messaging system into your Kubernetes environments.

---

## Example 1: Basic NATS Cluster with Default Settings

Deploy a simple NATS cluster with default replicas, resources, and JetStream enabled.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: NatsKubernetes
metadata:
  name: nats-basic
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "nats-basic"
  create_namespace: true
  serverContainer:
    replicas: 3
    diskSize: "10Gi"
  disableJetStream: false
  tlsEnabled: false
  disableNatsBox: false
```

**Use Case:**

* Ideal for quick deployments, testing, or simple messaging scenarios within Kubernetes clusters, leveraging default
  optimized settings.

---

## Example 2: NATS Cluster with Authentication Enabled (Bearer Token)

Set up a secure NATS cluster using bearer token authentication to protect client connections.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: NatsKubernetes
metadata:
  name: nats-secure
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "nats-secure"
  create_namespace: true
  serverContainer:
    replicas: 5
    resources:
      limits:
        cpu: "2000m"
        memory: "4Gi"
      requests:
        cpu: "500m"
        memory: "1Gi"
    diskSize: "20Gi"
  auth:
    enabled: true
    scheme: bearer_token
  tlsEnabled: true
  disableJetStream: false
```

**Use Case:**

* Recommended for production deployments where security, reliability, and scalable message handling are critical.

---

## Example 3: NATS Cluster Exposed via Ingress for External Clients

Deploy a NATS cluster configured with ingress to allow external clients to securely access messaging services.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: NatsKubernetes
metadata:
  name: nats-external
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "nats-external"
  create_namespace: true
  serverContainer:
    replicas: 3
    diskSize: "10Gi"
  ingress:
    enabled: true
    hostname: nats.example.com
  auth:
    enabled: true
    scheme: basic_auth
  tlsEnabled: true
  disableJetStream: false
```

**Use Case:**

* Suitable for scenarios where external systems or clients require secure connectivity to internal messaging
  infrastructure via HTTPS.

---

## Example 4: Lightweight NATS Cluster Without JetStream

Set up a lightweight NATS messaging cluster without JetStream persistence, optimized for minimal resource usage.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: NatsKubernetes
metadata:
  name: nats-minimal
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "nats-minimal"
  create_namespace: true
  serverContainer:
    replicas: 1
    resources:
      limits:
        cpu: "500m"
        memory: "512Mi"
      requests:
        cpu: "100m"
        memory: "128Mi"
    diskSize: "1Gi"
  disableJetStream: true
  tlsEnabled: false
  disableNatsBox: true
```

**Use Case:**

* Best suited for lightweight deployments, non-persistent messaging needs, or development environments with limited
  resources.

---

## Example 5: High Availability NATS Cluster with Advanced Metrics

Deploy a highly available NATS cluster that includes advanced observability with Prometheus metrics collection enabled.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: NatsKubernetes
metadata:
  name: nats-ha-metrics
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "nats-ha-metrics"
  create_namespace: true
  serverContainer:
    replicas: 7
    resources:
      limits:
        cpu: "4000m"
        memory: "8Gi"
      requests:
        cpu: "1000m"
        memory: "2Gi"
    diskSize: "50Gi"
  disableJetStream: false
  tlsEnabled: true
  disableNatsBox: false
  ingress:
    enabled: true
    hostname: nats-ha.example.com
  auth:
    enabled: true
    scheme: bearer_token
```

**Use Case:**

* Recommended for production-critical workloads requiring high availability, robust persistence, secure authentication,
  and detailed observability through Prometheus monitoring.

---

## Example 6: NATS with Pre-existing Namespace

Deploy NATS into an existing namespace that you manage separately.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: NatsKubernetes
metadata:
  name: nats-existing-ns
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "shared-messaging"
  create_namespace: false
  serverContainer:
    replicas: 3
    diskSize: "10Gi"
  auth:
    enabled: true
    scheme: basic_auth
  tlsEnabled: true
```

**Use Case:**

* Suitable when the namespace already exists with specific ResourceQuotas, LimitRanges, or NetworkPolicies
* Useful in GitOps workflows where namespace creation is handled separately
* Allows sharing a namespace across multiple messaging components
