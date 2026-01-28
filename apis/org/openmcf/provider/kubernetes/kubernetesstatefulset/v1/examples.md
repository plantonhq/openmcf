# Create using CLI

Create a yaml using the example shown below. After the yaml is created, use the below command to apply.

```shell
planton apply -f <yaml-path>
```

# Basic Example - PostgreSQL Database

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStatefulSet
metadata:
  name: postgres-db
spec:
  namespace: database
  createNamespace: true
  container:
    app:
      image:
        repo: postgres
        tag: "15"
      ports:
        - appProtocol: tcp
          containerPort: 5432
          name: postgres
          networkProtocol: TCP
          servicePort: 5432
      resources:
        requests:
          cpu: 100m
          memory: 256Mi
        limits:
          cpu: 1000m
          memory: 2Gi
      volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
    - name: data
      size: 10Gi
      accessModes:
        - ReadWriteOnce
```

# Example with Environment Variables and Secrets (Direct String Values)

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStatefulSet
metadata:
  name: postgres-db
spec:
  namespace: database
  createNamespace: true
  container:
    app:
      image:
        repo: postgres
        tag: "15"
      ports:
        - appProtocol: tcp
          containerPort: 5432
          name: postgres
          networkProtocol: TCP
          servicePort: 5432
      resources:
        requests:
          cpu: 100m
          memory: 256Mi
        limits:
          cpu: 1000m
          memory: 2Gi
      env:
        variables:
          POSTGRES_DB: myapp
          PGDATA: /var/lib/postgresql/data/pgdata
        secrets:
          POSTGRES_PASSWORD:
            value: supersecretpassword
          POSTGRES_USER:
            value: admin
      volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
    - name: data
      size: 20Gi
      storageClass: ssd
      accessModes:
        - ReadWriteOnce
```

# Example with Kubernetes Secret References (Production)

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStatefulSet
metadata:
  name: postgres-db
spec:
  namespace: database
  createNamespace: true
  container:
    app:
      image:
        repo: postgres
        tag: "15"
      ports:
        - appProtocol: tcp
          containerPort: 5432
          name: postgres
          networkProtocol: TCP
          servicePort: 5432
      resources:
        requests:
          cpu: 100m
          memory: 256Mi
        limits:
          cpu: 1000m
          memory: 2Gi
      env:
        variables:
          POSTGRES_DB: myapp
          PGDATA: /var/lib/postgresql/data/pgdata
        secrets:
          # Reference existing Kubernetes Secret for production credentials
          POSTGRES_PASSWORD:
            secretRef:
              name: postgres-credentials  # Name of existing K8s Secret
              key: password               # Key within the Secret
          POSTGRES_USER:
            secretRef:
              name: postgres-credentials
              key: username
      volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
    - name: data
      size: 20Gi
      storageClass: ssd
      accessModes:
        - ReadWriteOnce
```

# Example with Mixed Secret Types

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStatefulSet
metadata:
  name: postgres-db
spec:
  namespace: database
  createNamespace: true
  container:
    app:
      image:
        repo: postgres
        tag: "15"
      ports:
        - appProtocol: tcp
          containerPort: 5432
          name: postgres
          networkProtocol: TCP
          servicePort: 5432
      resources:
        requests:
          cpu: 100m
          memory: 256Mi
        limits:
          cpu: 1000m
          memory: 2Gi
      env:
        variables:
          POSTGRES_DB: myapp
          PGDATA: /var/lib/postgresql/data/pgdata
        secrets:
          # Dev secret - direct value (for development/testing)
          DEBUG_TOKEN:
            value: debug-only-token
          # Production secrets - external references (for production)
          POSTGRES_PASSWORD:
            secretRef:
              name: postgres-credentials
              key: password
          POSTGRES_USER:
            secretRef:
              name: postgres-credentials
              key: username
      volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
    - name: data
      size: 20Gi
      storageClass: ssd
      accessModes:
        - ReadWriteOnce
```

# Example - Redis Cluster with Multiple Replicas

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStatefulSet
metadata:
  name: redis-cluster
spec:
  namespace: cache
  createNamespace: true
  container:
    app:
      image:
        repo: redis
        tag: "7-alpine"
      ports:
        - appProtocol: tcp
          containerPort: 6379
          name: redis
          networkProtocol: TCP
          servicePort: 6379
        - appProtocol: tcp
          containerPort: 16379
          name: cluster-bus
          networkProtocol: TCP
          servicePort: 16379
      resources:
        requests:
          cpu: 100m
          memory: 128Mi
        limits:
          cpu: 500m
          memory: 512Mi
      command:
        - redis-server
      args:
        - --appendonly
        - "yes"
        - --cluster-enabled
        - "yes"
      volumeMounts:
        - name: data
          mountPath: /data
  availability:
    replicas: 6
  volumeClaimTemplates:
    - name: data
      size: 5Gi
      accessModes:
        - ReadWriteOnce
  podManagementPolicy: Parallel
```

# Example - MongoDB with Pod Disruption Budget

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStatefulSet
metadata:
  name: mongodb
spec:
  namespace: mongodb
  createNamespace: true
  container:
    app:
      image:
        repo: mongo
        tag: "6"
      ports:
        - appProtocol: tcp
          containerPort: 27017
          name: mongodb
          networkProtocol: TCP
          servicePort: 27017
      resources:
        requests:
          cpu: 200m
          memory: 512Mi
        limits:
          cpu: 2000m
          memory: 4Gi
      volumeMounts:
        - name: data
          mountPath: /data/db
  availability:
    replicas: 3
    podDisruptionBudget:
      enabled: true
      minAvailable: "2"
  volumeClaimTemplates:
    - name: data
      size: 50Gi
      storageClass: fast-ssd
      accessModes:
        - ReadWriteOnce
```

# Example with Ingress Enabled

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStatefulSet
metadata:
  name: elasticsearch
spec:
  namespace: search
  createNamespace: true
  container:
    app:
      image:
        repo: elasticsearch
        tag: "8.11.1"
      ports:
        - appProtocol: http
          containerPort: 9200
          name: http
          networkProtocol: TCP
          servicePort: 9200
          isIngressPort: true
        - appProtocol: tcp
          containerPort: 9300
          name: transport
          networkProtocol: TCP
          servicePort: 9300
      resources:
        requests:
          cpu: 500m
          memory: 2Gi
        limits:
          cpu: 2000m
          memory: 4Gi
      env:
        variables:
          discovery.type: single-node
          ES_JAVA_OPTS: "-Xms2g -Xmx2g"
      volumeMounts:
        - name: data
          mountPath: /usr/share/elasticsearch/data
  ingress:
    enabled: true
    hostname: elasticsearch.example.com
  volumeClaimTemplates:
    - name: data
      size: 100Gi
      accessModes:
        - ReadWriteOnce
```

# Example - Kafka with Multiple Volume Mounts

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesStatefulSet
metadata:
  name: kafka
spec:
  namespace: messaging
  createNamespace: true
  container:
    app:
      image:
        repo: confluentinc/cp-kafka
        tag: "7.5.0"
      ports:
        - appProtocol: tcp
          containerPort: 9092
          name: kafka
          networkProtocol: TCP
          servicePort: 9092
        - appProtocol: tcp
          containerPort: 9093
          name: kafka-internal
          networkProtocol: TCP
          servicePort: 9093
      resources:
        requests:
          cpu: 500m
          memory: 2Gi
        limits:
          cpu: 2000m
          memory: 4Gi
      env:
        variables:
          KAFKA_BROKER_ID: "0"
          KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
          KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-0.kafka-headless.messaging.svc.cluster.local:9092
      volumeMounts:
        - name: data
          mountPath: /var/lib/kafka/data
        - name: logs
          mountPath: /var/log/kafka
  availability:
    replicas: 3
  volumeClaimTemplates:
    - name: data
      size: 100Gi
      storageClass: fast-ssd
      accessModes:
        - ReadWriteOnce
    - name: logs
      size: 10Gi
      accessModes:
        - ReadWriteOnce
  podManagementPolicy: OrderedReady
```
