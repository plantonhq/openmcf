Here are a few examples for the `GcpSecretsManager` API resource, modeled in a similar way to how you created examples for the `MicroserviceKubernetes` API.

### Example 1: Basic Google Cloud Secrets Manager Setup

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSecretsManager
metadata:
  name: prod-secrets
spec:
  projectId: my-gcp-project
  secret_names:
    - database-password
    - api-key
```

### Example 2: Google Cloud Secrets Manager with Multiple Secrets

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSecretsManager
metadata:
  name: dev-secrets
spec:
  projectId: dev-gcp-project
  secret_names:
    - jwt-secret
    - database-url
    - oauth-token
```

### Example 3: Google Cloud Secrets Manager with Empty Spec

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSecretsManager
metadata:
  name: minimal-secrets
spec:
  projectId: my-gcp-project
```
