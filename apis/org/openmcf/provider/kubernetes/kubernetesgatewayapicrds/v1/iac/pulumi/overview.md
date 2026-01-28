# KubernetesGatewayApiCrds Pulumi Module Overview

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     Pulumi Module                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────────┐ │
│  │   main.go    │────▶│  module/     │────▶│ Kubernetes API   │ │
│  │  (entrypoint)│     │  main.go     │     │                  │ │
│  └──────────────┘     └──────────────┘     └──────────────────┘ │
│                              │                                   │
│                              ▼                                   │
│                       ┌──────────────┐                          │
│                       │  locals.go   │                          │
│                       │  (computed)  │                          │
│                       └──────────────┘                          │
│                              │                                   │
│                              ▼                                   │
│                       ┌──────────────┐                          │
│                       │  vars.go     │                          │
│                       │  (constants) │                          │
│                       └──────────────┘                          │
│                              │                                   │
│                              ▼                                   │
│                       ┌──────────────────────────────────────┐  │
│                       │  Gateway API CRD Manifest (YAML)      │  │
│                       │  from github.com/kubernetes-sigs/     │  │
│                       │  gateway-api/releases                  │  │
│                       └──────────────────────────────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow

1. **Input**: `KubernetesGatewayApiCrdsStackInput` loaded from environment
2. **Locals**: Version and channel computed from spec
3. **Manifest URL**: Constructed from version and channel
4. **Apply**: YAML manifest applied to cluster via Kubernetes provider
5. **Output**: Version, channel, and CRD list exported

## Key Design Decisions

### Remote Manifest Fetching

The module fetches CRD manifests directly from GitHub releases:
- Ensures official, unmodified CRDs
- No need to bundle manifests in the module
- Easy version updates

### No Namespace Required

Gateway API CRDs are cluster-scoped:
- No namespace creation needed
- Simpler deployment
- Single point of CRD management

### Channel Selection

Two channels support different use cases:
- **Standard**: Production-ready, stable APIs
- **Experimental**: Cutting-edge features, may change

### Idempotent Application

Using Pulumi's YAML ConfigFile:
- CRDs are applied idempotently
- Updates handled automatically
- Proper cleanup on destroy

## Resources Created

The module creates cluster-scoped CRDs only:

### Standard Channel
- `gatewayclasses.gateway.networking.k8s.io`
- `gateways.gateway.networking.k8s.io`
- `httproutes.gateway.networking.k8s.io`
- `referencegrants.gateway.networking.k8s.io`

### Experimental Channel (includes standard)
- All standard CRDs, plus:
- `tcproutes.gateway.networking.k8s.io`
- `udproutes.gateway.networking.k8s.io`
- `tlsroutes.gateway.networking.k8s.io`
- `grpcroutes.gateway.networking.k8s.io`

## Error Handling

The module handles:
- Network failures (GitHub unreachable)
- Invalid version (404 from releases)
- Permission denied (missing cluster-admin)
- CRD conflicts (existing CRDs)

## Upgrade Path

To upgrade Gateway API version:
1. Update `version` in manifest
2. Run `pulumi up`
3. CRDs are updated in-place

**Note**: Downgrading CRDs is not recommended and may cause issues with existing Gateway API resources.
