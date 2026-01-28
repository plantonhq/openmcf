# GCP Compute Instance Module Architecture

## Overview

This document describes the architecture and design decisions for the GCP Compute Instance Pulumi module.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    GcpComputeInstanceStackInput                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ target: GcpComputeInstance                               │    │
│  │   ├── metadata (name, org, env, labels)                  │    │
│  │   └── spec                                               │    │
│  │         ├── projectId (StringValueOrRef)                 │    │
│  │         ├── zone                                         │    │
│  │         ├── machineType                                  │    │
│  │         ├── bootDisk (image, size, type)                 │    │
│  │         ├── networkInterfaces[]                          │    │
│  │         ├── serviceAccount                               │    │
│  │         ├── scheduling                                   │    │
│  │         └── ...                                          │    │
│  └─────────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ providerConfig: GcpProviderConfig                        │    │
│  │   └── gcpCredentialJson                                  │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Pulumi Module                              │
│  ┌──────────────────┐                                           │
│  │    main.go       │ ── Load stack input, setup provider       │
│  └────────┬─────────┘                                           │
│           │                                                      │
│           ▼                                                      │
│  ┌──────────────────┐                                           │
│  │   locals.go      │ ── Transform spec to locals, build labels │
│  └────────┬─────────┘                                           │
│           │                                                      │
│           ▼                                                      │
│  ┌──────────────────┐                                           │
│  │  instance.go     │ ── Create compute instance resource       │
│  └────────┬─────────┘                                           │
│           │                                                      │
│           ▼                                                      │
│  ┌──────────────────┐                                           │
│  │   outputs.go     │ ── Export stack outputs                   │
│  └──────────────────┘                                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    GCP Compute Engine                            │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ google_compute_instance                                  │    │
│  │   ├── name: {metadata.name}                              │    │
│  │   ├── project: {spec.projectId}                          │    │
│  │   ├── zone: {spec.zone}                                  │    │
│  │   ├── machine_type: {spec.machineType}                   │    │
│  │   ├── boot_disk                                          │    │
│  │   │     ├── image: {spec.bootDisk.image}                 │    │
│  │   │     ├── size: {spec.bootDisk.sizeGb}                 │    │
│  │   │     └── type: {spec.bootDisk.type}                   │    │
│  │   ├── network_interface[]                                │    │
│  │   │     ├── network: {networkInterfaces[].network}       │    │
│  │   │     ├── subnetwork: {networkInterfaces[].subnetwork} │    │
│  │   │     └── access_config[]                              │    │
│  │   ├── service_account                                    │    │
│  │   ├── scheduling                                         │    │
│  │   ├── labels                                             │    │
│  │   ├── tags                                               │    │
│  │   └── metadata                                           │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Stack Outputs                                 │
│  ├── instance_name                                              │
│  ├── instance_id                                                │
│  ├── self_link                                                  │
│  ├── internal_ip                                                │
│  ├── external_ip                                                │
│  ├── zone                                                       │
│  ├── machine_type                                               │
│  └── cpu_platform                                               │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow

1. **Stack Input Loading**: The `main.go` entrypoint loads the `GcpComputeInstanceStackInput` from the `STACK_INPUT` environment variable.

2. **Provider Setup**: GCP provider is configured using credentials from `providerConfig.gcpCredentialJson`.

3. **Locals Initialization**: The `locals.go` module transforms the spec into local variables and builds the label map.

4. **Instance Creation**: The `instance.go` module creates the `google_compute_instance` resource with all configurations.

5. **Output Export**: Stack outputs are exported for use by other modules or the CLI.

## Key Design Decisions

### StringValueOrRef for Cross-Resource References

Fields like `projectId`, `network`, and `subnetwork` use `StringValueOrRef` to support:
- **Literal values**: Direct string values for simple deployments
- **Value references**: References to other OpenMCF resources

This enables declarative resource composition without hardcoding values.

### Label Standardization

All instances get standard OpenMCF labels:
- `resource`: "true"
- `resource_kind`: "gcpcomputeinstance"
- `resource_name`: Instance name
- `resource_id`: Instance ID (if provided)
- `organization`: Organization (if provided)
- `environment`: Environment (if provided)

User labels are merged with these standard labels.

### Network Interface Flexibility

The module supports:
- Multiple network interfaces for multi-NIC VMs
- Optional external IPs via access configs
- Alias IP ranges for container workloads
- Both network and subnetwork specification

### Scheduling Options

Supports all GCP scheduling options:
- Standard VMs with live migration
- Preemptible VMs (legacy)
- Spot VMs (recommended for cost savings)
- Custom termination actions

## Error Handling

- All resource creation functions return errors that propagate to Pulumi
- Provider setup errors are wrapped with context
- Missing required fields are caught by proto validation before reaching the module

## Testing Strategy

1. **Unit Tests**: Proto validation tests in `spec_test.go`
2. **Integration Tests**: Deploy to test project with `iac/hack/manifest.yaml`
3. **Manual Testing**: Use `debug.sh` for step-through debugging
