# KubernetesGhaRunnerScaleSet Examples

## Repository-Level Runners with Docker Support

Runners connected to a specific repository, with Docker-in-Docker for building images:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSet
metadata:
  name: repo-runners
spec:
  namespace:
    value: gha-runners
  createNamespace: true
  
  github:
    configUrl: https://github.com/myorg/myrepo
    patToken:
      token: ghp_xxxxxxxxxxxx
  
  containerMode:
    type: DIND
  
  scaling:
    minRunners: 1    # Keep one warm for fast startup
    maxRunners: 10
  
  runner:
    resources:
      requests:
        cpu: "1"
        memory: "2Gi"
      limits:
        cpu: "4"
        memory: "8Gi"
```

## Organization-Level Runners with GitHub App

Runners available to all repositories in an organization:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSet
metadata:
  name: org-runners
spec:
  namespace:
    value: gha-runners
  createNamespace: true
  
  github:
    configUrl: https://github.com/myorg
    githubApp:
      appId: "123456"
      installationId: "654321"
      # Private key must be base64 encoded
      # Generate with: cat private-key.pem | base64
      privateKeyBase64: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVB...
  
  runnerGroup: "production"
  runnerScaleSetName: "org-prod-runners"
  
  containerMode:
    type: DIND
  
  scaling:
    minRunners: 2
    maxRunners: 20
```

## Kubernetes Mode for Container Workflows

Each workflow step runs in its own container:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSet
metadata:
  name: k8s-runners
spec:
  namespace:
    value: gha-runners
  createNamespace: true
  
  github:
    configUrl: https://github.com/myorg
    patToken:
      token: ghp_xxxxxxxxxxxx
  
  containerMode:
    type: KUBERNETES
    workVolumeClaim:
      storageClass: fast-ssd
      size: "50Gi"
      accessModes:
        - ReadWriteOnce
  
  scaling:
    minRunners: 0
    maxRunners: 15
```

Workflow example:

```yaml
jobs:
  build:
    runs-on: [self-hosted, k8s-runners]
    container: node:20-alpine
    steps:
      - uses: actions/checkout@v4
      - run: npm ci
      - run: npm test
```

## Runners with Persistent Caches

Cache npm, gradle, and Docker layers for faster builds:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSet
metadata:
  name: cached-runners
spec:
  namespace:
    value: gha-runners
  createNamespace: true
  
  github:
    configUrl: https://github.com/myorg/myrepo
    patToken:
      token: ghp_xxxxxxxxxxxx
  
  containerMode:
    type: DIND
  
  scaling:
    minRunners: 0
    maxRunners: 5
  
  runner:
    resources:
      requests:
        cpu: "2"
        memory: "4Gi"
      limits:
        cpu: "4"
        memory: "8Gi"
  
  persistentVolumes:
    - name: npm-cache
      size: "20Gi"
      mountPath: /home/runner/.npm
      storageClass: standard
    
    - name: gradle-cache
      size: "30Gi"
      mountPath: /home/runner/.gradle
    
    - name: maven-cache
      size: "30Gi"
      mountPath: /home/runner/.m2
    
    - name: docker-cache
      size: "100Gi"
      mountPath: /var/lib/docker
      storageClass: fast-ssd
```

## Custom Runner Image

Use a custom runner image with pre-installed tools:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSet
metadata:
  name: custom-runners
spec:
  namespace:
    value: gha-runners
  createNamespace: true
  
  github:
    configUrl: https://github.com/myorg
    patToken:
      token: ghp_xxxxxxxxxxxx
  
  containerMode:
    type: DIND
  
  runner:
    image:
      repository: ghcr.io/myorg/custom-runner
      tag: v1.2.3
      pullPolicy: IfNotPresent
    
    env:
      - name: RUNNER_DEBUG
        value: "1"
      - name: CUSTOM_TOOL_PATH
        value: /opt/tools
    
    resources:
      requests:
        cpu: "500m"
        memory: "1Gi"
      limits:
        cpu: "2"
        memory: "4Gi"
  
  imagePullSecrets:
    - ghcr-secret
  
  scaling:
    minRunners: 0
    maxRunners: 10
```

## Enterprise Runners with Controller Reference

For clusters where automatic controller discovery doesn't work:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSet
metadata:
  name: enterprise-runners
spec:
  namespace:
    value: gha-runners
  createNamespace: true
  
  github:
    configUrl: https://github.com/enterprises/myenterprise
    existingSecretName: enterprise-github-credentials
  
  containerMode:
    type: DIND
  
  controllerServiceAccount:
    namespace: arc-system
    name: arc-gha-rs-controller
  
  scaling:
    minRunners: 5    # Always keep some warm
    maxRunners: 50
  
  runnerGroup: "enterprise-default"
  
  labels:
    team: platform
    tier: production
  
  annotations:
    description: "Enterprise-wide GitHub Actions runners"
```

## Minimal Configuration

Smallest possible configuration for testing:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGhaRunnerScaleSet
metadata:
  name: test-runners
spec:
  namespace:
    value: default
  
  github:
    configUrl: https://github.com/myuser/myrepo
    patToken:
      token: ghp_xxxxxxxxxxxx
  
  containerMode:
    type: DEFAULT
  
  scaling:
    maxRunners: 2
```

