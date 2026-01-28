# KubernetesStrimziKafkaOperator Examples

## Example 1: Basic Operator Deployment

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStrimziKafkaOperator
metadata:
  name: kafka-operator-basic
spec:
  targetCluster:
    clusterName: "my-gke-cluster"
  namespace:
    value: "strimzi-kafka-operator"
  createNamespace: true
  container: {}
```

## Example 2: Production Operator (Custom Resources)

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStrimziKafkaOperator
metadata:
  name: kafka-operator-prod
spec:
  targetCluster:
    clusterName: "production-gke-cluster"
  namespace:
    value: "strimzi-kafka-operator"
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 2000m
        memory: 2Gi
```

## Example 3: Specific Cluster Deployment

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStrimziKafkaOperator
metadata:
  name: kafka-operator-gke
spec:
  targetCluster:
    clusterName: "prod-gke-cluster-01"
  namespace:
    value: "strimzi-kafka-operator"
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 100Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

## Example 4: Using Existing Namespace

This example demonstrates using a pre-existing namespace. The namespace must already exist in the cluster before deploying the operator.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStrimziKafkaOperator
metadata:
  name: kafka-operator-existing-ns
spec:
  targetCluster:
    clusterName: "my-gke-cluster"
  namespace:
    value: "platform-operators"
  createNamespace: false
  container:
    resources:
      requests:
        cpu: 50m
        memory: 100Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

**Note:** When `createNamespace: false`, ensure the namespace exists before deployment. This is useful when:
- The namespace has pre-configured quotas, labels, or annotations
- Multiple operators share the same namespace
- Namespace lifecycle is managed by a separate GitOps process
```

## Post-Deployment: Kafka Cluster Examples

### Example: Production Kafka Cluster (3 Brokers + ZooKeeper)

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
  name: production-kafka
  namespace: kafka-production
spec:
  kafka:
    version: 3.6.0
    replicas: 3
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
      - name: tls
        port: 9093
        type: internal
        tls: true
      - name: external
        port: 9094
        type: loadbalancer
        tls: true
    config:
      offsets.topic.replication.factor: 3
      transaction.state.log.replication.factor: 3
      transaction.state.log.min.isr: 2
      default.replication.factor: 3
      min.insync.replicas: 2
      log.retention.hours: 168
    storage:
      type: jbod
      volumes:
        - id: 0
          type: persistent-claim
          size: 500Gi
          deleteClaim: false
  zookeeper:
    replicas: 3
    storage:
      type: persistent-claim
      size: 50Gi
      deleteClaim: false
  entityOperator:
    topicOperator: {}
    userOperator: {}
```

### Example: Dev Kafka Cluster (Minimal)

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
  name: dev-kafka
  namespace: kafka-dev
spec:
  kafka:
    version: 3.6.0
    replicas: 1
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
    storage:
      type: ephemeral
  zookeeper:
    replicas: 1
    storage:
      type: ephemeral
  entityOperator:
    topicOperator: {}
```

### Example: Kafka Topic

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaTopic
metadata:
  name: orders
  namespace: kafka-production
  labels:
    strimzi.io/cluster: production-kafka
spec:
  partitions: 12
  replicas: 3
  config:
    retention.ms: 604800000  # 7 days
    segment.bytes: 1073741824  # 1GB
```

### Example: Kafka User (SCRAM-SHA-512 Auth)

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaUser
metadata:
  name: order-service
  namespace: kafka-production
  labels:
    strimzi.io/cluster: production-kafka
spec:
  authentication:
    type: scram-sha-512
  authorization:
    type: simple
    acls:
      - resource:
          type: topic
          name: orders
        operations:
          - Read
          - Write
      - resource:
          type: group
          name: order-processors
        operations:
          - Read
```

## Troubleshooting

```bash
# Check operator status
kubectl get pods -n strimzi-kafka-operator

# View operator logs
kubectl logs -n strimzi-kafka-operator -l name=strimzi-cluster-operator

# List Kafka clusters
kubectl get kafka --all-namespaces

# Check specific Kafka cluster status
kubectl describe kafka production-kafka -n kafka-production
```

## Additional Resources

- [README.md](README.md) - Component overview
- [docs/README.md](docs/README.md) - Deep-dive research (32KB)
- https://strimzi.io/quickstarts/ - Official quickstarts

