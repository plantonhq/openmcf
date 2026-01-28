# Kubernetes Tekton Operator - Architecture Overview

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           Pulumi Module                                  │
│                                                                          │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────────┐  │
│  │     main.go     │───▶│   module/       │───▶│ Kubernetes Cluster  │  │
│  │   (entrypoint)  │    │   main.go       │    │                     │  │
│  └─────────────────┘    │   locals.go     │    │  ┌───────────────┐  │  │
│                         │   outputs.go    │    │  │ tekton-       │  │  │
│                         │   vars.go       │    │  │ operator ns   │  │  │
│                         │   tekton_       │    │  │               │  │  │
│                         │   operator.go   │    │  │  ┌─────────┐  │  │  │
│                         └─────────────────┘    │  │  │Operator │  │  │  │
│                                                │  │  │Pod      │  │  │  │
│                                                │  │  └─────────┘  │  │  │
│                                                │  └───────────────┘  │  │
│                                                │                     │  │
│                                                │  ┌───────────────┐  │  │
│                                                │  │ tekton-       │  │  │
│                                                │  │ pipelines ns  │  │  │
│                                                │  │               │  │  │
│                                                │  │  Components   │  │  │
│                                                │  │  (managed by  │  │  │
│                                                │  │   operator)   │  │  │
│                                                │  └───────────────┘  │  │
│                                                └─────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

## Resource Flow

```
KubernetesTektonOperatorStackInput (proto)
         │
         ▼
    ┌─────────────────────────────────────────────────────┐
    │                 initializeLocals()                   │
    │                                                      │
    │  • Extract component settings                        │
    │  • Set namespace values                              │
    │  • Prepare Kubernetes labels                         │
    │  • Export initial outputs                            │
    └─────────────────────────────────────────────────────┘
         │
         ▼
    ┌─────────────────────────────────────────────────────┐
    │              tektonOperator()                        │
    │                                                      │
    │  1. Install Tekton Operator manifests               │
    │     └─▶ yamlv2.NewConfigFile()                      │
    │         (uses yaml/v2 for CRD ordering)             │
    │                                                      │
    │  2. Create TektonConfig CRD                         │
    │     └─▶ yamlv2.NewConfigGroup()                     │
    │         • profile: all/basic/lite                   │
    │         • targetNamespace: tekton-pipelines         │
    └─────────────────────────────────────────────────────┘
         │
         ▼
    ┌─────────────────────────────────────────────────────┐
    │          Tekton Operator (in-cluster)                │
    │                                                      │
    │  Reads TektonConfig and installs:                   │
    │  • Tekton Pipelines (if profile includes)           │
    │  • Tekton Triggers (if profile includes)            │
    │  • Tekton Dashboard (if profile = all)              │
    └─────────────────────────────────────────────────────┘
```

## Component Selection Logic

```
User Input (components)          Profile Selected        Components Installed
─────────────────────            ────────────────        ────────────────────
pipelines: true                  
triggers: true        ───────▶   profile: "all"   ───▶  Pipelines + Triggers
dashboard: true                                          + Dashboard

pipelines: true                  
triggers: true        ───────▶   profile: "basic" ───▶  Pipelines + Triggers
dashboard: false                 

pipelines: true                  
triggers: false       ───────▶   profile: "lite"  ───▶  Pipelines only
dashboard: false                 
```

## Key Design Decisions

### 1. YAML Manifests via yaml/v2

**Decision**: Use official YAML manifests with Pulumi's `yaml/v2` module.

**Rationale**:
- Tekton Operator is primarily distributed via YAML manifests
- No official Helm chart from Tekton project
- Manifests are stable and well-tested
- Simpler dependency management
- **yaml/v2** provides better CRD ordering and await behavior, ensuring
  CRDs are registered before TektonConfig is created (prevents 
  "no matches for kind TektonConfig" errors)

### 2. Operator-Managed Components

**Decision**: Let the Tekton Operator manage component installation.

**Rationale**:
- Operator handles version compatibility
- Automatic reconciliation of component state
- Simpler upgrade path
- TektonConfig provides unified configuration

### 3. Profile-Based Installation

**Decision**: Map component booleans to TektonConfig profiles.

**Rationale**:
- Aligns with Tekton Operator's design
- Reduces complexity in IaC module
- Ensures correct component dependencies

## Namespace Strategy

```
┌────────────────────────┐     ┌────────────────────────┐
│   tekton-operator      │     │   tekton-pipelines     │
│   (Operator Namespace) │     │   (Components Namespace)│
│                        │     │                        │
│   ┌──────────────────┐ │     │   ┌──────────────────┐ │
│   │ tekton-operator  │ │     │   │ tekton-pipelines │ │
│   │ controller       │─┼─────┼──▶│ controller       │ │
│   └──────────────────┘ │     │   └──────────────────┘ │
│                        │     │                        │
│   Watches:             │     │   ┌──────────────────┐ │
│   • TektonConfig       │     │   │ tekton-triggers  │ │
│   • TektonPipeline     │     │   │ controller       │ │
│   • TektonTrigger      │     │   └──────────────────┘ │
│   • TektonDashboard    │     │                        │
│                        │     │   ┌──────────────────┐ │
│                        │     │   │ tekton-dashboard │ │
│                        │     │   └──────────────────┘ │
└────────────────────────┘     └────────────────────────┘
```

## Error Handling

```
main.go
   │
   ├─▶ LoadStackInput() ─────┬─▶ Error: "failed to load stack-input"
   │                         │
   ▼                         │
module.Resources()           │
   │                         │
   ├─▶ initializeLocals() ───┤
   │                         │
   ├─▶ GetKubernetesProvider()┬─▶ Error: "setup kubernetes provider"
   │                         │
   ├─▶ tektonOperator() ─────┬─▶ Error: "deploy tekton operator"
   │       │                 │
   │       ├─▶ yamlv2.NewConfigFile()──▶ Error: "install tekton operator manifests"
   │       │                 │
   │       └─▶ yamlv2.NewConfigGroup()─▶ Error: "create tekton config"
   │                         │
   └─▶ Success               │
```

## Upgrade Path

```
Current State                  Upgrade Steps
─────────────                  ─────────────
v0.77.0                        1. Update spec.operator_version in manifest
                                  operator_version: "v0.78.0"
                               2. Run: planton apply -f manifest.yaml
                                       │
                                       ▼
                               Pulumi detects changes to release URL
                                       │
                                       ▼
                               Applies new operator manifests
                                       │
                                       ▼
v0.78.0                        Operator reconciles TektonConfig
                               and upgrades components
```

## Security Considerations

### RBAC Requirements

The Tekton Operator requires cluster-admin privileges:

```
ClusterRole: tekton-operator
  Rules:
  - apiGroups: ["*"]
    resources: ["*"]
    verbs: ["*"]
```

### Network Security

Recommended NetworkPolicies:

```
Operator Namespace:
- Egress to Kubernetes API
- Egress to Tekton registries

Components Namespace:
- Ingress for webhooks (Triggers)
- Ingress for Dashboard (if exposed)
- Egress for pipeline tasks
```

## References

- [Tekton Operator Design](https://github.com/tektoncd/operator/blob/main/docs/TektonConfig.md)
- [TektonConfig CRD](https://tekton.dev/docs/operator/)
- [Pulumi Kubernetes YAML](https://www.pulumi.com/registry/packages/kubernetes/api-docs/yaml/configfile/)
- [Pulumi Kubernetes YAML v2](https://www.pulumi.com/blog/kubernetes-yaml-v2/) - Better CRD ordering and await behavior
