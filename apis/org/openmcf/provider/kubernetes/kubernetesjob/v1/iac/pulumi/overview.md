# KubernetesJob Pulumi Module - Architecture Overview

## Design Principles

This module follows OpenMCF's standard patterns for Kubernetes workload deployments:

1. **Declarative Configuration**: All configuration comes from the protobuf-defined spec
2. **Idempotent Operations**: Resources are created/updated based on desired state
3. **Separation of Concerns**: Each resource type has its own file
4. **Consistent Labeling**: All resources share common labels for tracking

## Resource Creation Flow

```
┌─────────────────────────────────────────────────────────────┐
│                     Stack Input                              │
│  (KubernetesJob + ProviderConfig + DockerConfigJson)        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Initialize Locals                          │
│  - Extract target, namespace, labels                         │
│  - Compute resource names (secrets, etc.)                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Create Kubernetes Provider                      │
│  - Configure from provider_config (kubeconfig, context)      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│            Create Namespace (if createNamespace)             │
│  - Uses locals.Namespace                                     │
│  - Applies standard labels                                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Create Secrets                             │
│  - Environment secrets with direct values                    │
│  - Image pull secret (if docker credentials provided)        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Create ConfigMaps                           │
│  - From spec.config_maps map                                 │
│  - One ConfigMap per entry                                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Create Job                                │
│  - Container with image, resources, env vars                 │
│  - Volume mounts from spec                                   │
│  - Job settings (parallelism, completions, etc.)             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Export Outputs                             │
│  - namespace                                                 │
│  - job_name                                                  │
└─────────────────────────────────────────────────────────────┘
```

## Key Design Decisions

### 1. Environment Variable Handling

Environment variables support two modes:
- **Direct values**: Plain text values set directly
- **Secret references**: References to existing Kubernetes Secrets

For direct secret values, we create an internal Secret (`{name}-env-secrets`) to avoid exposing sensitive data in the Job spec.

### 2. Volume Mount Architecture

Volume mounts are processed by `buildVolumeMountsAndVolumes()` which:
1. Creates volume mount entries for the container
2. Creates corresponding volume definitions for the pod

Supported volume types:
- ConfigMaps (inline or external)
- Secrets
- HostPath
- EmptyDir
- PersistentVolumeClaim (PVC)

### 3. Resource Naming

To avoid conflicts when multiple jobs share a namespace:
- Secrets are named `{job-name}-env-secrets`
- Image pull secrets are named `{job-name}-image-pull-secret`
- ConfigMaps use the key from `spec.configMaps`

### 4. Label Consistency

All resources share these labels:
- `openmcf.org/resource: "true"`
- `openmcf.org/resource-name: {name}`
- `openmcf.org/resource-kind: KubernetesJob`
- `openmcf.org/resource-id: {id}` (if set)
- `openmcf.org/organization: {org}` (if set)
- `openmcf.org/environment: {env}` (if set)

## Job-Specific Configuration

### Parallelism and Completions

- `parallelism`: Max pods running simultaneously
- `completions`: Required successful completions

Common patterns:
- **Single pod job**: parallelism=1, completions=1 (default)
- **Parallel batch**: parallelism=N, completions=M (where M >= N)
- **Work queue**: parallelism=N, completions unset (pods run until queue empty)

### Completion Mode

- **NonIndexed** (default): Pods are interchangeable
- **Indexed**: Each pod gets `JOB_COMPLETION_INDEX` (0 to completions-1)

### Failure Handling

- `backoffLimit`: Retries before job failure (default: 6)
- `activeDeadlineSeconds`: Max job duration
- `restartPolicy`: "Never" (new pod per retry) or "OnFailure" (restart in-place)

### Cleanup

- `ttlSecondsAfterFinished`: Auto-delete after completion
- `suspend`: Pause job creation without deleting

## Dependencies

```
main.go
   └── module.Resources()
         ├── initializeLocals() → Locals struct
         ├── pulumikubernetesprovider.Get() → Provider
         ├── namespace() → Namespace (optional)
         ├── secret() → Secret (optional)
         ├── createImagePullSecret() → Secret (optional)
         ├── configMaps() → ConfigMaps (optional)
         ├── job() → Job
         └── exportOutputs()
```

## Error Handling

All resource creation functions:
1. Return errors wrapped with context
2. Use `errors.Wrap()` for stack traces
3. Allow partial cleanup via Pulumi's state management

## Testing

The module can be tested using the hack manifest:

```bash
cd iac/pulumi
export STACK_INPUT=$(cat ../hack/manifest.yaml | base64)
pulumi preview
```
