
# Example 1: Basic Elasticsearch and Kibana Deployment

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ElasticsearchKubernetes
metadata:
  name: logging-cluster
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: logging
  create_namespace: true
  elasticsearch:
    container:
      replicas: 1
      resources:
        requests:
          cpu: 500m
          memory: 1Gi
        limits:
          cpu: 1000m
          memory: 2Gi
      persistence_enabled: true
      disk_size: 10Gi
  kibana:
    enabled: true
    container:
      replicas: 1
      resources:
        requests:
          cpu: 200m
          memory: 512Mi
        limits:
          cpu: 500m
          memory: 1Gi
```

---

# Example 2: Elasticsearch with Ingress Enabled

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ElasticsearchKubernetes
metadata:
  name: search-service
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: search
  create_namespace: true
  elasticsearch:
    container:
      replicas: 3
      resources:
        requests:
          cpu: 1
          memory: 2Gi
        limits:
          cpu: 2
          memory: 4Gi
      persistence_enabled: true
      disk_size: 50Gi
    ingress:
      enabled: true
      hostname: search.example.com
  kibana:
    enabled: true
    container:
      replicas: 1
      resources:
        requests:
          cpu: 200m
          memory: 512Mi
        limits:
          cpu: 500m
          memory: 1Gi
    ingress:
      enabled: true
      hostname: search-kibana.example.com
```

---

# Example 3: Elasticsearch with Multiple Replicas

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ElasticsearchKubernetes
metadata:
  name: logging-app
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: logging-app
  create_namespace: true
  elasticsearch:
    container:
      replicas: 5
      resources:
        requests:
          cpu: 500m
          memory: 1Gi
        limits:
          cpu: 1000m
          memory: 2Gi
      persistence_enabled: true
      disk_size: 20Gi
  kibana:
    enabled: true
    container:
      replicas: 2
      resources:
        requests:
          cpu: 200m
          memory: 512Mi
        limits:
          cpu: 500m
          memory: 1Gi
```

---

# Example 4: Using Existing Namespace (create_namespace: false)

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ElasticsearchKubernetes
metadata:
  name: shared-elasticsearch
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: shared-services
  create_namespace: false
  elasticsearch:
    container:
      replicas: 1
      resources:
        requests:
          cpu: 500m
          memory: 1Gi
        limits:
          cpu: 1000m
          memory: 2Gi
      persistence_enabled: true
      disk_size: 10Gi
  kibana:
    enabled: true
    container:
      replicas: 1
      resources:
        requests:
          cpu: 200m
          memory: 512Mi
        limits:
          cpu: 500m
          memory: 1Gi
```

---

# Example 4: Minimal Elasticsearch Deployment with Defaults

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ElasticsearchKubernetes
metadata:
  name: minimal-elasticsearch
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: elasticsearch
```
