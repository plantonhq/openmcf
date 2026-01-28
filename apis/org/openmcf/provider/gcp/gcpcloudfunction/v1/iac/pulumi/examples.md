Here are a few examples for the `GcpCloudFunction` API resource, showcasing different configurations for deploying Google Cloud Functions (Gen 2). These examples demonstrate basic configurations, environment variables, and use of secrets.

---

### Basic Example

Deploy a simple HTTP-triggered Cloud Function with minimal configuration:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudFunction
metadata:
  name: my-cloud-function
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  buildConfig:
    runtime: python311
    entryPoint: hello_http
    source:
      bucket: my-code-bucket
      object: functions/my-function-v1.0.0.zip
```

---

### Example with Project Reference (ValueFrom)

Deploy a function using a reference to a GcpProject resource for dynamic project ID resolution:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudFunction
metadata:
  name: my-cloud-function-with-ref
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: main-project
      fieldPath: status.outputs.project_id
  region: us-central1
  buildConfig:
    runtime: python311
    entryPoint: hello_http
    source:
      bucket: my-code-bucket
      object: functions/my-function-v1.0.0.zip
```

---

### Example with Environment Variables

Deploy a function with custom environment variables for configuration:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudFunction
metadata:
  name: my-cloud-function-env-vars
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  buildConfig:
    runtime: nodejs20
    entryPoint: handleRequest
    source:
      bucket: my-code-bucket
      object: functions/api-v1.0.0.zip
  serviceConfig:
    environmentVariables:
      DATABASE_NAME: my_database
      LOG_LEVEL: info
```

---

### Example with Secrets from GCP Secret Manager

Deploy a function with secrets injected from Google Secret Manager:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudFunction
metadata:
  name: my-cloud-function-secrets
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  buildConfig:
    runtime: go122
    entryPoint: ProcessRequest
    source:
      bucket: my-code-bucket
      object: functions/processor-v1.0.0.zip
  serviceConfig:
    secretEnvironmentVariables:
      DATABASE_PASSWORD: database-password
      API_KEY: external-api-key
```

---

### Production Example with VPC and Scaling

Deploy a production-ready function with VPC connectivity, scaling, and authentication:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudFunction
metadata:
  name: production-function
  env: prod
  org: my-org
spec:
  projectId:
    value: my-gcp-project-prod
  region: us-east1
  buildConfig:
    runtime: nodejs22
    entryPoint: handleRequest
    source:
      bucket: my-prod-code-bucket
      object: functions/api-gateway-v3.4.1.zip
  serviceConfig:
    serviceAccountEmail: api-gateway@my-gcp-project-prod.iam.gserviceaccount.com
    availableMemoryMb: 1024
    timeoutSeconds: 60
    maxInstanceRequestConcurrency: 80
    vpcConnector: projects/my-gcp-project-prod/locations/us-east1/connectors/vpc-connector
    ingressSettings: 0  # ALLOW_ALL
    scaling:
      minInstanceCount: 2
      maxInstanceCount: 100
    allowUnauthenticated: true
```
