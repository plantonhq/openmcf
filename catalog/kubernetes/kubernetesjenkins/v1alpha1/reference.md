# KubernetesJenkins

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**KubernetesJenkinsSpec** deploys Jenkins — the extendable CI/CD
automation server — from the official jenkinsci Helm chart
(https://github.com/jenkinsci/helm-charts). The module generates a
random admin password into a per-deploy Kubernetes Secret (never
authored in the manifest) and exports the admin username plus the
Secret name/key as outputs, so consumers reference credentials instead
of pasting them. The spec covers container resource sizing, a
`helm_values` escape hatch for chart options beyond the typed surface,
and an ingress toggle with the external hostname.

## Example

```yaml
# Full-surface offline-proof manifest: exercises the namespace reference
# with namespace creation, explicit container resource sizing, the
# helm_values escape hatch, and ingress with an external hostname — the
# spec's whole typed surface. The admin password is module-generated into
# a Secret and never appears here. Placeholder values; never applied to a
# real cluster.
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesJenkins
metadata:
  name: hack-jenkins
spec:
  namespace:
    value: hack-jenkins
  createNamespace: true
  containerResources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      cpu: "1"
      memory: 2Gi
  helmValues:
    controller.componentName: jenkins-controller
    persistence.size: 10Gi
  ingress:
    enabled: true
    hostname: jenkins.example.com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.containerResources` | `ContainerResources` |  |  |  |
| `spec.containerResources.limits` | `CpuMemory` |  |  |  |
| `spec.containerResources.limits.cpu` | `string` |  |  |  |
| `spec.containerResources.limits.memory` | `string` |  |  |  |
| `spec.containerResources.requests` | `CpuMemory` |  |  |  |
| `spec.containerResources.requests.cpu` | `string` |  |  |  |
| `spec.containerResources.requests.memory` | `string` |  |  |  |
| `spec.helmValues` | `map<string, string>` |  |  |  |
| `spec.ingress` | `KubernetesJenkinsIngress` |  |  |  |
| `spec.ingress.enabled` | `bool` |  |  |  |
| `spec.ingress.hostname` | `string` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Kubernetes Namespace

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.createNamespace

`bool`

flag to indicate if the namespace should be created

### spec.containerResources

`ContainerResources`

The CPU and memory resources allocated to the Jenkins container.

### spec.containerResources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.containerResources.limits.cpu

`string`

### spec.containerResources.limits.memory

`string`

### spec.containerResources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.containerResources.requests.cpu

`string`

### spec.containerResources.requests.memory

`string`

### spec.helmValues

`map<string, string>`

A map of key-value pairs that provide additional customization options for the Helm chart used to deploy Jenkins.
These values allow for further refinement of the deployment, such as customizing resource limits, setting environment variables,
or specifying version tags. For detailed information on the available options, refer to the Helm chart documentation at:
https://github.com/jenkinsci/helm-charts/blob/main/charts/jenkins/values.yaml

### spec.ingress

`KubernetesJenkinsIngress`

The ingress configuration for the Jenkins deployment.

- rule: hostname is required when ingress is enabled

### spec.ingress.enabled

`bool`

Flag to enable or disable ingress.

### spec.ingress.hostname

`string`

The full hostname for external access (e.g., "jenkins.example.com").
Required when enabled is true.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesJenkins, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | kubernetes namespace in which jenkins-kubernetes is created. |
| `status.outputs.service` | `string` | kubernetes service name for jenkins-kubernetes. ex: main-jenkins-kubernetes in the above example, "main" is the name of the jenkins-kubernetes |
| `status.outputs.port_forward_command` | `string` | command to setup port-forwarding to open jenkins-kubernetes from developers laptop. this might come handy when jenkins-kubernetes ingress is disabled for security reasons. this is rendered by combining jenkins_kubernetes_kubernetes_service and kubernetes_namespace ex: kubectl port-forward svc/jenkins_kubernetes_kubernetes_service -n kubernetes_namespace 8080:8080 running the command from this attribute makes it possible to access jenkins-kubernetes using http://localhost:8080 |
| `status.outputs.kube_endpoint` | `string` | kubernetes endpoint to connect to jenkins-kubernetes from the web browser. ex: main-jenkins-kubernetes.namespace.svc.cluster.local:8080 |
| `status.outputs.external_hostname` | `string` | public endpoint to open jenkins-kubernetes from clients outside kubernetes. ex: https://jnk-planton-pcs-dev-main.data.dev.planton.live:8080 |
| `status.outputs.internal_hostname` | `string` | internal postgres-kubernetes hostname. ex: https://jnk-planton-pcs-dev-main.data-internal.dev.planton.live:8080 |
| `status.outputs.username` | `string` | jenkins username |
| `status.outputs.password_secret` | `KubernetesSecretKey` | kubernetes secret key for the password. |
| `status.outputs.password_secret.name` | `string` | The name of the Kubernetes Secret. |
| `status.outputs.password_secret.key` | `string` | The key within the Kubernetes Secret. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |

## See Also

- [Overview](../README.md)
