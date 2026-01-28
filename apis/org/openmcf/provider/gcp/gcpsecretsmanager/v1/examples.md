# GcpSecretsManager Examples

This document provides comprehensive examples for using the `GcpSecretsManager` API resource to manage secrets in Google Cloud Secret Manager.

## Create using CLI

Create a YAML file using one of the examples shown below. After the YAML is created, use the following command to apply it:

```shell
planton apply -f <yaml-path>
```

## Example 1: Basic Secret Creation with Literal Project ID

This example demonstrates creating a simple set of secrets in a GCP project using a literal project ID value. The secrets will be created with automatic replication across multiple regions for high availability.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSecretsManager
metadata:
  name: app-secrets
spec:
  projectId:
    value: my-gcp-project-123456
  secretNames:
    - api-key
    - database-password
    - oauth-client-secret
```

**What this creates:**
- Three empty secret containers in GCP Secret Manager
- Each secret has automatic multi-region replication
- Secrets are created with placeholder values that should be updated with actual values

**Setting secret values:**

After creation, update secret values using the `gcloud` CLI:

```shell
echo -n "actual-api-key-value" | gcloud secrets versions add api-key \
  --project=my-gcp-project-123456 \
  --data-file=-
```

## Example 2: Using Project Reference (Cross-Resource Dependency)

This example demonstrates using a reference to another resource's output instead of a literal project ID. This is useful when you want to create secrets in a project that is managed by another `GcpProject` resource, enabling cross-resource dependencies.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSecretsManager
metadata:
  name: app-secrets-with-ref
  env: prod
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: main-project
      fieldPath: status.outputs.project_id
  secretNames:
    - api-key
    - database-password
    - oauth-client-secret
```

**What this creates:**
- Same as Example 1, but the project ID is dynamically resolved from the `GcpProject` resource named `main-project`
- The `fieldPath` specifies which output field to use (`status.outputs.project_id`)
- This pattern enables declarative infrastructure composition where secrets are automatically created in the correct project

## Example 3: Environment-Specific Secrets

This example shows how to create environment-specific secrets by using the `env` metadata field. The environment name will be prepended to each secret name.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSecretsManager
metadata:
  name: prod-app-secrets
  env: prod
spec:
  projectId:
    value: my-gcp-project-123456
  secretNames:
    - database-url
    - redis-password
    - jwt-signing-key
    - stripe-api-key
```

**What this creates:**
- Secret names will be: `prod-database-url`, `prod-redis-password`, `prod-jwt-signing-key`, `prod-stripe-api-key`
- Allows separation between production and non-production secrets
- Prevents accidental access to production secrets from other environments

## Example 4: Multi-Service Application Secrets

This example demonstrates organizing secrets for a complex application with multiple services.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSecretsManager
metadata:
  name: ecommerce-platform-secrets
  org: acme-corp
  env: production
spec:
  projectId:
    value: acme-ecommerce-prod
  secretNames:
    # Database credentials
    - postgres-admin-password
    - postgres-app-password
    - redis-password
    # External service API keys
    - stripe-secret-key
    - sendgrid-api-key
    - twilio-auth-token
    # Application secrets
    - jwt-private-key
    - encryption-key
    - oauth-client-secret
    # Third-party integrations
    - github-webhook-secret
    - slack-webhook-url
```

**Best practices demonstrated:**
- Organized secret naming (grouped by purpose)
- Clear metadata (org, env) for tracking and access control
- Comprehensive coverage of different secret types

## Example 5: Development Environment Secrets

This example creates a minimal set of secrets for a development environment.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSecretsManager
metadata:
  name: dev-secrets
  env: dev
spec:
  projectId:
    value: my-project-dev
  secretNames:
    - dev-database-password
    - dev-api-key
    - dev-oauth-secret
```

**Development considerations:**
- Use a separate GCP project for dev secrets
- Secret values can be less complex for development
- Consider using the same secret structure as production for consistency

## Example 6: Microservices Secrets with Service Account Access

This example shows creating secrets for a microservices architecture where different services need access to different secrets.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSecretsManager
metadata:
  name: microservices-secrets
  env: staging
spec:
  projectId:
    value: microservices-staging-project
  secretNames:
    # Auth service secrets
    - auth-service-jwt-key
    - auth-service-db-password
    # Payment service secrets
    - payment-service-stripe-key
    - payment-service-db-password
    # Notification service secrets
    - notification-service-sendgrid-key
    - notification-service-twilio-token
    # Shared secrets
    - shared-redis-password
    - shared-rabbitmq-password
```

**IAM configuration** (applied separately using Terraform/Pulumi or gcloud):

```shell
# Grant auth service access to its secrets
gcloud secrets add-iam-policy-binding staging-auth-service-jwt-key \
  --project=microservices-staging-project \
  --member="serviceAccount:auth-service@microservices-staging-project.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"

# Grant payment service access to its secrets
gcloud secrets add-iam-policy-binding staging-payment-service-stripe-key \
  --project=microservices-staging-project \
  --member="serviceAccount:payment-service@microservices-staging-project.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

## Example 7: Secrets for CI/CD Pipeline

This example creates secrets specifically for CI/CD pipelines and deployment automation.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSecretsManager
metadata:
  name: cicd-secrets
spec:
  projectId:
    value: my-devops-project
  secretNames:
    # Container registry credentials
    - docker-registry-token
    - gcr-service-account-key
    # Deployment credentials
    - kubernetes-cluster-token
    - helm-repository-password
    # External service tokens
    - github-api-token
    - slack-notification-webhook
    - pagerduty-integration-key
    # Artifact signing
    - code-signing-certificate
    - artifact-signing-key
```

**CI/CD access pattern:**

In Cloud Build:

```yaml
availableSecrets:
  secretManager:
    - versionName: "projects/my-devops-project/secrets/github-api-token/versions/latest"
      env: "GITHUB_TOKEN"
    - versionName: "projects/my-devops-project/secrets/docker-registry-token/versions/latest"
      env: "DOCKER_TOKEN"
```

## Example 8: Integration with Kubernetes via External Secrets Operator

This example shows secrets that will be synced to a Kubernetes cluster using External Secrets Operator (ESO).

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSecretsManager
metadata:
  name: k8s-app-secrets
  env: prod
spec:
  projectId:
    value: my-gke-project
  secretNames:
    # Application secrets to be synced to K8s
    - app-database-url
    - app-redis-url
    - app-api-keys
    - app-tls-cert
    - app-tls-key
```

**External Secrets Operator configuration:**

```yaml
# Create after the GcpSecretsManager resource is deployed
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: app-secrets
  namespace: production
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: gcp-secret-store
  target:
    name: app-secrets
    creationPolicy: Owner
  data:
    - secretKey: database-url
      remoteRef:
        key: prod-app-database-url
        version: latest
    - secretKey: redis-url
      remoteRef:
        key: prod-app-redis-url
        version: latest
```

## Best Practices Summary

1. **Use environment prefixes**: Add `env` to metadata to namespace secrets by environment
2. **Organize by purpose**: Group related secrets with clear naming conventions
3. **Separate projects**: Use different GCP projects for production vs. non-production secrets
4. **Set values securely**: Never commit actual secret values to Git; use CLI or secure CI/CD variables
5. **Apply least privilege**: Grant IAM permissions per-secret, not project-wide
6. **Document access patterns**: Clearly document which services need access to which secrets
7. **Plan for rotation**: Design secret names and application code to support credential rotation

## Common Integration Patterns

### Pattern 1: Application Access (Go)

```go
import (
    "context"
    secretmanager "cloud.google.com/go/secretmanager/apiv1"
    "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func getSecret(projectID, secretID string) (string, error) {
    ctx := context.Background()
    client, err := secretmanager.NewClient(ctx)
    if err != nil {
        return "", err
    }
    defer client.Close()

    req := &secretmanagerpb.AccessSecretVersionRequest{
        Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID),
    }
    result, err := client.AccessSecretVersion(ctx, req)
    if err != nil {
        return "", err
    }
    return string(result.Payload.Data), nil
}
```

### Pattern 2: Cloud Run Environment Variable

In Cloud Run deployment:

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: my-service
spec:
  template:
    spec:
      containers:
        - image: gcr.io/my-project/my-app
          env:
            - name: DATABASE_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: prod-database-password
                  key: latest
```

### Pattern 3: Cloud Function Access

```python
from google.cloud import secretmanager

def access_secret(project_id, secret_id):
    client = secretmanager.SecretManagerServiceClient()
    name = f"projects/{project_id}/secrets/{secret_id}/versions/latest"
    response = client.access_secret_version(request={"name": name})
    return response.payload.data.decode('UTF-8')
```

## Notes

- Secrets are created with placeholder values initially
- Update secret values using `gcloud secrets versions add` command
- Secret names must be unique within a GCP project
- Labels are automatically added based on metadata (org, env, resource name)
- Automatic replication provides high availability across multiple regions
- Consider using Customer-Managed Encryption Keys (CMEK) for sensitive compliance requirements
