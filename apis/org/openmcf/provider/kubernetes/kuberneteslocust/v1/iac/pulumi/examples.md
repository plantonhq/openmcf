# Example 1: Basic Locust Kubernetes Setup

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: LocustKubernetes
metadata:
  name: locust-basic
spec:
  targetCluster:
    clusterName: "my-gke-cluster"
  namespace:
    value: "locust-test"
  create_namespace: true
  master_container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    replicas: 1
  worker_container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    replicas: 2
  load_test:
    name: basic-load-test
    main_py_content: |
      from locust import HttpUser, task

      class MyUser(HttpUser):
          @task
          def my_task(self):
              self.client.get("/api/test")
  ingress:
    enabled: false
```

# Example 2: Locust Kubernetes with Custom Helm Values

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: LocustKubernetes
metadata:
  name: locust-custom
spec:
  targetCluster:
    clusterName: "production-gke-cluster"
  namespace:
    value: "locust-prod"
  create_namespace: true
  master_container:
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 2
        memory: 2Gi
    replicas: 1
  worker_container:
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 2
        memory: 2Gi
    replicas: 5
  load_test:
    name: custom-load-test
    main_py_content: |
      from locust import HttpUser, task

      class MyUser(HttpUser):
          @task
          def my_task(self):
              self.client.post("/api/test", json={"key": "value"})
    lib_files_content:
      utils.py: |
        def helper_function():
            return "Helper"
    pip_packages:
      - requests
      - locust
  ingress:
    enabled: true
    ingressClassName: "nginx"
    hosts:
      - host: locust.mydomain.com
        paths:
          - /
```

# Example 3: Locust Kubernetes with TLS and Ingress

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: LocustKubernetes
metadata:
  name: locust-tls
spec:
  targetCluster:
    clusterName: "my-gke-cluster"
  namespace:
    value: "locust-tls"
  create_namespace: true
  master_container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    replicas: 1
  worker_container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    replicas: 3
  load_test:
    name: tls-load-test
    main_py_content: |
      from locust import HttpUser, task

      class MyUser(HttpUser):
          @task
          def my_task(self):
              self.client.get("/secure-api/test")
  ingress:
    enabled: true
    ingressClassName: "nginx"
    hosts:
      - host: locust-tls.mydomain.com
        paths:
          - /
    tls:
      - secretName: locust-tls-cert
        hosts:
          - locust-tls.mydomain.com
```

# Example 4: Locust Kubernetes with External Library and PIP Packages

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: LocustKubernetes
metadata:
  name: locust-external-lib
spec:
  targetCluster:
    clusterName: "dev-cluster"
  namespace:
    value: "locust-dev"
  create_namespace: true
  master_container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    replicas: 1
  worker_container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    replicas: 2
  load_test:
    name: external-lib-load-test
    main_py_content: |
      from locust import HttpUser, task
      from utils import helper_function

      class MyUser(HttpUser):
          @task
          def my_task(self):
              result = helper_function()
              self.client.get(f"/api/test?result={result}")
    lib_files_content:
      utils.py: |
        def helper_function():
            return "Hello from helper!"
    pip_packages:
      - requests
      - locust
  ingress:
    enabled: false
```

# Example 5: Locust Kubernetes Minimal Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: LocustKubernetes
metadata:
  name: locust-minimal
spec:
  targetCluster:
    clusterName: "test-cluster"
  namespace:
    value: "locust-minimal"
  create_namespace: true
  master_container:
    resources:
      requests:
        cpu: 50m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
    replicas: 1
  worker_container:
    resources:
      requests:
        cpu: 50m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
      replicas: 1
  load_test:
    name: minimal-test
    main_py_content: |
      from locust import HttpUser, task

      class MyUser(HttpUser):
          @task
          def my_task(self):
              self.client.get("/")
  ingress:
    enabled: false
```

# Example 6: Using Existing Namespace

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: LocustKubernetes
metadata:
  name: locust-existing-ns
spec:
  targetCluster:
    clusterName: "prod-gke-cluster"
  namespace:
    value: "shared-load-testing"
  create_namespace: false  # Use existing namespace
  master_container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    replicas: 1
  worker_container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    replicas: 2
  load_test:
    name: existing-ns-test
    main_py_content: |
      from locust import HttpUser, task

      class MyUser(HttpUser):
          @task
          def my_task(self):
              self.client.get("/api/test")
  ingress:
    enabled: false
```
