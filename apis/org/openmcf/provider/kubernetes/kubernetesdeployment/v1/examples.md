# Create using CLI

Create a yaml using the example shown below. After the yaml is created, use the below command to apply.

```shell
planton apply -f <yaml-path>
```

# Basic Example

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: todo-list-api
spec:
  version: main
  container:
    app:
      image:
        repo: nginx
        tag: latest
      ports:
        - appProtocol: http
          containerPort: 8080
          isIngressPort: true
          name: rest-api
          networkProtocol: TCP
          servicePort: 80
      resources:
        requests:
          cpu: 100m
          memory: 100Mi
        limits:
          cpu: 2000m
          memory: 2Gi
```

# Example w/ Environment Variables

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: todo-list-api
spec:
  version: main
  container:
    app:
      env:
        variables:
          DATABASE_NAME: todo
      image:
        repo: nginx
        tag: latest
      ports:
        - appProtocol: http
          containerPort: 8080
          isIngressPort: true
          name: rest-api
          networkProtocol: TCP
          servicePort: 80
      resources:
        requests:
          cpu: 100m
          memory: 100Mi
        limits:
          cpu: 2000m
          memory: 2Gi
```

# Example w/ Environment Secrets (Direct String Values)

This example shows how to provide secrets as direct string values. A Kubernetes Secret is automatically
created to store these values securely. This approach is suitable for development and testing.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: todo-list-api
spec:
  namespace:
    value: my-namespace
  version: main
  container:
    app:
      env:
        variables:
          DATABASE_NAME: todo
        secrets:
          DATABASE_PASSWORD:
            value: my-secret-password
          API_KEY:
            value: abc123
      image:
        repo: nginx
        tag: latest
      ports:
        - appProtocol: http
          containerPort: 8080
          isIngressPort: true
          name: rest-api
          networkProtocol: TCP
          servicePort: 80
      resources:
        requests:
          cpu: 100m
          memory: 100Mi
        limits:
          cpu: 2000m
          memory: 2Gi
```

# Example w/ Environment Secrets (Kubernetes Secret References)

This example shows how to reference existing Kubernetes Secrets. This is the recommended approach for
production deployments as it avoids storing sensitive values in configuration files.

**Prerequisites:** The referenced Kubernetes Secrets must exist in the cluster before deploying.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: todo-list-api
spec:
  namespace:
    value: my-namespace
  version: main
  container:
    app:
      env:
        variables:
          DATABASE_NAME: todo
        secrets:
          DATABASE_PASSWORD:
            secretRef:
              name: my-app-secrets
              key: db-password
          API_KEY:
            secretRef:
              name: external-api-credentials
              key: api-key
      image:
        repo: nginx
        tag: latest
      ports:
        - appProtocol: http
          containerPort: 8080
          isIngressPort: true
          name: rest-api
          networkProtocol: TCP
          servicePort: 80
      resources:
        requests:
          cpu: 100m
          memory: 100Mi
        limits:
          cpu: 2000m
          memory: 2Gi
```

# Example w/ Mixed Secret Types

This example demonstrates using both direct string values and Kubernetes Secret references together.
You can mix and match based on your needs - use direct values for development secrets and references
for production credentials.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: todo-list-api
spec:
  namespace:
    value: my-namespace
  version: main
  container:
    app:
      env:
        variables:
          DATABASE_NAME: todo
          LOG_LEVEL: info
        secrets:
          # Direct value - suitable for non-critical secrets in dev
          DEBUG_TOKEN:
            value: debug-only-token
          # External secret reference - recommended for production credentials
          DATABASE_PASSWORD:
            secretRef:
              name: postgres-credentials
              key: password
          # Another external reference
          AWS_ACCESS_KEY_ID:
            secretRef:
              name: aws-credentials
              key: access-key-id
          AWS_SECRET_ACCESS_KEY:
            secretRef:
              name: aws-credentials
              key: secret-access-key
      image:
        repo: nginx
        tag: latest
      ports:
        - appProtocol: http
          containerPort: 8080
          isIngressPort: true
          name: rest-api
          networkProtocol: TCP
          servicePort: 80
      resources:
        requests:
          cpu: 100m
          memory: 100Mi
        limits:
          cpu: 2000m
          memory: 2Gi
```

# Example w/ ConfigMaps and Volume Mounts

This example shows how to create ConfigMaps and mount them as files in the container.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: tekton-cloudevents-router
spec:
  namespace:
    value: tekton-pipelines
  version: main
  configMaps:
    router-config: |
      routes:
        - namespace_prefix: planton-dev-
          target: http://servicehub-tekton-webhooks.planton-dev.svc/webhook
        - namespace_prefix: planton-prod-
          target: http://servicehub-tekton-webhooks.planton-prod.svc/webhook
  container:
    app:
      image:
        repo: ghcr.io/plantoncloud/tekton-cloud-event-router
        tag: v0.1.0
      volumeMounts:
        - name: router-config
          mountPath: /etc/router/config.yaml
          configMap:
            name: router-config
            key: router-config
      ports:
        - name: http
          containerPort: 8080
          networkProtocol: TCP
          appProtocol: http
          servicePort: 80
      resources:
        requests:
          cpu: 50m
          memory: 64Mi
        limits:
          cpu: 200m
          memory: 256Mi
```

# Example w/ Command and Args Override

This example shows how to override the container's command and arguments.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: custom-entrypoint-app
spec:
  namespace:
    value: my-namespace
  version: main
  container:
    app:
      image:
        repo: busybox
        tag: latest
      command:
        - /bin/sh
        - -c
      args:
        - echo "Hello, World!" && sleep 3600
      ports:
        - name: http
          containerPort: 8080
          networkProtocol: TCP
          appProtocol: http
          servicePort: 80
      resources:
        requests:
          cpu: 50m
          memory: 64Mi
        limits:
          cpu: 100m
          memory: 128Mi
```

# Example w/ Multiple Volume Types

This example demonstrates mounting various volume types including ConfigMaps, Secrets, EmptyDir, and HostPath.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesDeployment
metadata:
  name: multi-volume-app
spec:
  namespace:
    value: my-namespace
  version: main
  configMaps:
    app-config: |
      database:
        host: postgres.default.svc
        port: 5432
      cache:
        host: redis.default.svc
        port: 6379
  container:
    app:
      image:
        repo: my-app
        tag: v1.0.0
      volumeMounts:
        # Mount ConfigMap as a single file
        - name: config-volume
          mountPath: /etc/app/config.yaml
          configMap:
            name: app-config
            key: app-config
        # Mount a Secret
        - name: tls-certs
          mountPath: /etc/tls
          readOnly: true
          secret:
            name: my-tls-secret
        # EmptyDir for temporary storage
        - name: cache
          mountPath: /tmp/cache
          emptyDir:
            medium: Memory
            sizeLimit: 256Mi
      ports:
        - name: http
          containerPort: 8080
          networkProtocol: TCP
          appProtocol: http
          servicePort: 80
      resources:
        requests:
          cpu: 100m
          memory: 128Mi
        limits:
          cpu: 500m
          memory: 512Mi
```
