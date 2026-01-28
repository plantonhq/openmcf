# Create using CLI

Create a YAML file using the examples shown below. After the YAML is created, use the following command to apply:

```shell
planton apply -f <yaml-path>
```

# Basic Example - Single ConfigMap

Deploy a simple ConfigMap:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesManifest
metadata:
  name: my-configmap
spec:
  namespace: my-namespace
  create_namespace: true
  manifest_yaml: |
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: app-config
    data:
      database_url: postgres://localhost:5432/mydb
      log_level: info
```

# Example - Multiple Resources

Deploy multiple resources in a single manifest:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesManifest
metadata:
  name: complete-app
spec:
  namespace: production
  create_namespace: true
  manifest_yaml: |
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: app-config
    data:
      environment: production
    ---
    apiVersion: v1
    kind: Secret
    metadata:
      name: app-secrets
    type: Opaque
    stringData:
      api-key: my-secret-key
    ---
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: my-app
    spec:
      replicas: 3
      selector:
        matchLabels:
          app: my-app
      template:
        metadata:
          labels:
            app: my-app
        spec:
          containers:
          - name: app
            image: nginx:latest
            ports:
            - containerPort: 80
            envFrom:
            - configMapRef:
                name: app-config
            - secretRef:
                name: app-secrets
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: my-app
    spec:
      selector:
        app: my-app
      ports:
      - port: 80
        targetPort: 80
```

# Example - Custom Resource Definition (CRD)

Deploy a CRD and its Custom Resource:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesManifest
metadata:
  name: my-operator-resources
spec:
  namespace: operators
  create_namespace: true
  manifest_yaml: |
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      name: myresources.example.com
    spec:
      group: example.com
      versions:
      - name: v1
        served: true
        storage: true
        schema:
          openAPIV3Schema:
            type: object
            properties:
              spec:
                type: object
                properties:
                  replicas:
                    type: integer
      scope: Namespaced
      names:
        plural: myresources
        singular: myresource
        kind: MyResource
    ---
    apiVersion: example.com/v1
    kind: MyResource
    metadata:
      name: my-instance
    spec:
      replicas: 2
```

# Example - RBAC Resources

Deploy ServiceAccount with RBAC:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesManifest
metadata:
  name: app-rbac
spec:
  namespace: my-app
  create_namespace: false
  manifest_yaml: |
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: my-app-sa
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: Role
    metadata:
      name: my-app-role
    rules:
    - apiGroups: [""]
      resources: ["configmaps", "secrets"]
      verbs: ["get", "list", "watch"]
    - apiGroups: [""]
      resources: ["pods"]
      verbs: ["get", "list"]
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: RoleBinding
    metadata:
      name: my-app-rolebinding
    subjects:
    - kind: ServiceAccount
      name: my-app-sa
    roleRef:
      kind: Role
      name: my-app-role
      apiGroup: rbac.authorization.k8s.io
```

# Example - Network Policy

Deploy a NetworkPolicy for security:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesManifest
metadata:
  name: network-security
spec:
  namespace: secure-app
  create_namespace: true
  manifest_yaml: |
    apiVersion: networking.k8s.io/v1
    kind: NetworkPolicy
    metadata:
      name: deny-all-ingress
    spec:
      podSelector: {}
      policyTypes:
      - Ingress
    ---
    apiVersion: networking.k8s.io/v1
    kind: NetworkPolicy
    metadata:
      name: allow-from-frontend
    spec:
      podSelector:
        matchLabels:
          app: backend
      policyTypes:
      - Ingress
      ingress:
      - from:
        - podSelector:
            matchLabels:
              app: frontend
        ports:
        - protocol: TCP
          port: 8080
```

# Example - With Target Cluster

Deploy to a specific Kubernetes cluster:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesManifest
metadata:
  name: targeted-deployment
spec:
  target_cluster:
    cluster_kind: GcpGkeCluster
    cluster_name: prod-cluster
  namespace: production
  create_namespace: false
  manifest_yaml: |
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: cluster-specific-config
    data:
      cluster: prod-cluster
      region: us-central1
```

# Example - Existing Namespace

Deploy to an existing namespace without creating it:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesManifest
metadata:
  name: use-existing-ns
spec:
  namespace: kube-system
  create_namespace: false
  manifest_yaml: |
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: custom-kube-config
    data:
      setting: value
```

