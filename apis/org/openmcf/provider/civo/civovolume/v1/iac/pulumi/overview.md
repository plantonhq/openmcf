# CivoVolume Pulumi Module Architecture

## High-Level Overview

The CivoVolume Pulumi module provisions Civo block storage volumes. It translates a declarative Protobuf specification (`CivoVolumeSpec`) into Civo infrastructure resources, handling volume creation and output management. The module gracefully handles provider limitations by logging informational messages when requested features (filesystem formatting, snapshots, tags) aren't supported by the Civo API.

```
┌─────────────────────────────────────────────────────────────────┐
│                    CivoVolume Manifest (YAML)                   │
│                                                                  │
│  apiVersion: civo.openmcf.org/v1                        │
│  kind: CivoVolume                                               │
│  metadata: {name: prod-db-data}                                 │
│  spec:                                                          │
│    volumeName: prod-db-data                                     │
│    region: LON1                                                 │
│    sizeGib: 100                                                 │
│    filesystemType: XFS                                          │
│    tags: [env:prod, backup:daily]                              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ (Protobuf deserialization)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                  CivoVolumeStackInput (Proto)                   │
│                                                                  │
│  CivoVolume: {metadata, spec}                                   │
│  ProviderConfig: {civo_token}                                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ (Pulumi module entry)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Pulumi Resources                           │
│                                                                  │
│  1. Initialize Locals (metadata, labels, spec)                  │
│  2. Create Civo Provider (with API token)                       │
│  3. Create Volume (with validation)                             │
│  4. Handle Limitations (log info/warnings)                      │
│  5. Export Outputs (volume ID)                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ (Civo API calls)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Civo Infrastructure                         │
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ Volume (Block Storage)                                    │ │
│  │                                                           │ │
│  │ Name: prod-db-data                                        │ │
│  │ Region: LON1                                              │ │
│  │ Size: 100 GiB                                             │ │
│  │ Status: unattached, unformatted                           │ │
│  └───────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ (Stack outputs)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                CivoVolumeStackOutputs (Proto)                   │
│                                                                  │
│  volume_id: "uuid-1234-5678"                                    │
│  attached_instance_id: "" (empty until attached)                │
│  device_path: "" (available after attachment)                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ (Next steps: attachment & formatting)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Post-Deployment Steps                         │
│                                                                  │
│  1. Attach volume to instance (manual or Kubernetes CSI)        │
│  2. Format volume (mkfs.ext4 or mkfs.xfs)                       │
│  3. Mount volume to filesystem                                  │
│  4. Configure /etc/fstab for persistence                        │
└─────────────────────────────────────────────────────────────────┘
```

## Module Components

### 1. Entry Point (`main.go`)

The Pulumi stack initialization file. Responsibilities:
- Parse command-line arguments or environment variables
- Deserialize `CivoVolumeStackInput` from JSON/YAML
- Call `module.Resources()` with Pulumi context and stack input
- Handle top-level errors

**Flow**:
```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Parse stack input (from env var or file)
        stackInput := parseStackInput()
        
        // Call module
        return module.Resources(ctx, stackInput)
    })
}
```

### 2. Module Entry (`module/main.go`)

The core module logic. Responsibilities:
- Initialize locals (metadata, labels, spec)
- Create Civo provider
- Orchestrate volume creation
- Return errors

**Flow**:
```go
func Resources(ctx *pulumi.Context, stackInput *CivoVolumeStackInput) error {
    // 1. Prepare locals
    locals := initializeLocals(ctx, stackInput)
    
    // 2. Setup Civo provider
    civoProvider, err := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
    if err != nil {
        return errors.Wrap(err, "failed to setup civo provider")
    }
    
    // 3. Create volume
    _, err = volume(ctx, locals, civoProvider)
    return err
}
```

### 3. Locals Initialization (`module/locals.go`)

Prepares local variables for use throughout the module. Responsibilities:
- Extract metadata (name, labels, description)
- Store CivoVolume spec
- Prepare Civo labels for resource tracking

**Data Structure**:
```go
type Locals struct {
    CivoProviderConfig *civoprovider.CivoProviderConfig
    CivoVolume         *civovolumev1.CivoVolume
    CivoLabels         map[string]string
}
```

**Civo Labels Applied**:
- `planton.org/resource`: "true"
- `planton.org/resource-kind`: "CivoVolume"
- `planton.org/resource-id`: volume ID
- `planton.org/resource-name`: volume name
- `planton.org/organization`: org name
- `planton.org/environment`: env name

**Why Locals?**
- **Reusability**: Shared across resource creation functions
- **Clarity**: Centralized data access
- **Immutability**: Prepared once, used many times

### 4. Volume Creation (`module/volume.go`)

Core resource provisioning logic. Responsibilities:
- Create Civo Volume with validated parameters
- Handle filesystem_type (informational - log message)
- Handle snapshot_id (warn if not supported)
- Handle tags (informational - not applied to Civo resource)
- Export outputs

**Volume Creation**:
```go
volumeArgs := &civo.VolumeArgs{
    Name:   pulumi.String(locals.CivoVolume.Spec.VolumeName),
    Region: pulumi.String(locals.CivoVolume.Spec.Region.String()),
    SizeGb: pulumi.Int(int(locals.CivoVolume.Spec.SizeGib)),
}

createdVolume, err := civo.NewVolume(
    ctx,
    "volume",
    volumeArgs,
    pulumi.Provider(civoProvider),
)
```

**Key Points**:
- `Name`: Validated in proto (lowercase, alphanumeric + hyphens)
- `Region`: Civo region enum mapped to string
- `SizeGb`: 1-16,000 GiB range (validated in proto)

**Filesystem Type Handling**:
```go
if spec.FilesystemType != NONE {
    ctx.Log.Info(fmt.Sprintf(
        "Filesystem type '%s' requested. Note: Civo provider doesn't "+
        "expose filesystem formatting. Volume created unformatted. "+
        "Use cloud-init or configuration management to format after attachment.",
        filesystemName,
    ), nil)
}
```

**Rationale**: Civo provider doesn't support filesystem formatting during creation. The module logs a reminder for post-deployment formatting.

**Snapshot Handling**:
```go
if spec.SnapshotId != "" {
    ctx.Log.Warn(fmt.Sprintf(
        "Snapshot ID '%s' specified. Note: Civo Volume snapshots are not "+
        "supported on public Civo cloud. This parameter is reserved for "+
        "CivoStack deployments. Volume will be created empty.",
        spec.SnapshotId,
    ), nil)
}
```

**Rationale**: Snapshot functionality only exists in CivoStack (private cloud), not public Civo.

**Tags Handling**:
```go
if len(spec.Tags) > 0 {
    ctx.Log.Info(fmt.Sprintf(
        "Tags specified: %v. Note: Civo Volume provider doesn't support tags. "+
        "Tags recorded in metadata but not applied to Civo resource. "+
        "Use Civo labels (applied automatically) for resource organization.",
        spec.Tags,
    ), nil)
}
```

**Rationale**: Civo Volume provider doesn't support tags. Tags are for OpenMCF metadata only.

### 5. Outputs (`module/outputs.go`)

Defines output constant names for `ctx.Export()`. Responsibilities:
- Standardize output key names
- Ensure consistency across modules
- Document what gets exported

**Constants**:
```go
const (
    OpVolumeId = "volume_id"
)
```

**Export Example**:
```go
ctx.Export(OpVolumeId, createdVolume.ID())
// Note: attached_instance_id and device_path are only available
// after attachment (handled separately or by Kubernetes CSI driver)
```

## Resource Dependency Graph

```
┌──────────────────────┐
│  Civo Provider       │
│  (API token)         │
└──────────┬───────────┘
           │
           │ (provides)
           │
           ▼
┌──────────────────────┐
│   Volume             │
│   (block storage)    │
│                      │
│  - Name              │
│  - Region            │
│  - Size (GiB)        │
│  - Status: active    │
└──────────────────────┘
```

**Dependency Chain**:
1. **Provider**: Must exist first (initialized in `Resources()`)
2. **Volume**: Depends on provider (authenticated API calls)

Volumes are created independently. Attachment is a separate operation (handled manually, via Civo CLI, or automatically by Kubernetes CSI driver).

## Limitations Handling

The module gracefully handles three categories of limitations:

### 1. Filesystem Formatting (Informational)

**Requested Feature**: `filesystemType: EXT4` or `XFS`  
**Provider Support**: ❌ Not exposed  
**Module Behavior**: Log info message, create unformatted volume  
**User Action**: Format manually after attachment using `mkfs.ext4` or `mkfs.xfs`

### 2. Snapshots (Warning)

**Requested Feature**: `snapshotId: "snapshot-12345"`  
**Provider Support**: ❌ Not available on public Civo (only CivoStack)  
**Module Behavior**: Log warning, create empty volume  
**User Action**: Implement application-level backups to object storage

### 3. Tags (Informational)

**Requested Feature**: `tags: ["env:prod", "backup:daily"]`  
**Provider Support**: ❌ Not exposed by Civo Volume provider  
**Module Behavior**: Log info message, tags recorded in metadata only  
**User Action**: Use Civo labels (automatically applied by OpenMCF) for resource organization

## State Management

### Pulumi State

Pulumi tracks resource state in a backend. State includes:
- Resource ID (volume ID)
- Resource attributes (name, region, size)
- Outputs (exported values)

**State Operations**:
- `pulumi up`: Creates/updates volume, writes state
- `pulumi refresh`: Syncs state with actual infrastructure
- `pulumi destroy`: Deletes volume, removes from state

### State Persistence

```
Pulumi State Backend
    │
    └─► Resource: civo:index/volume:Volume
            ├─► ID: "volume-uuid-1234"
            ├─► Name: "prod-db-data"
            ├─► Region: "LON1"
            ├─► SizeGb: 100
            └─► Status: "active"
```

## Performance Characteristics

### Resource Creation Times

| Operation | Typical Duration |
|-----------|------------------|
| Provider initialization | < 1 second |
| Volume creation | 5-15 seconds |
| Output export | < 1 second |
| **Total** | **~10-20 seconds** |

### Scaling Considerations

- **Single Volume**: ~15 seconds
- **10 Volumes (parallel)**: ~20-30 seconds (Pulumi parallelizes by default)
- **100 Volumes**: Limited by Civo API rate limits

**Optimization**: Pulumi automatically parallelizes independent resources (e.g., multiple volumes).

## Security Considerations

### Secrets Management

1. **Civo API Token**: Stored as Pulumi secret
   ```bash
   pulumi config set civo:token $CIVO_TOKEN --secret
   ```

2. **State Encryption**: Pulumi backend encrypts state at rest

### Least Privilege

- Use separate Civo API tokens per environment (dev, staging, prod)
- Rotate tokens periodically

### Audit Trail

- Pulumi logs all operations (create, update, delete)
- Civo provides audit logs for API calls

## Post-Deployment Workflow

After Pulumi creates a volume:

1. **Attach Volume**:
   - Manual: Civo Console or CLI
   - Kubernetes: CSI driver handles automatically

2. **Format Volume** (first-time setup):
   ```bash
   ssh root@<instance-ip>
   mkfs.ext4 /dev/vdb  # or mkfs.xfs
   ```

3. **Mount Volume**:
   ```bash
   mkdir -p /data
   mount /dev/vdb /data
   echo "/dev/vdb /data ext4 defaults,nofail 0 2" >> /etc/fstab
   ```

4. **Use Volume**: Write data, mount in applications

## Error Handling

### Common Error Scenarios

| Error | Cause | Recovery |
|-------|-------|----------|
| "volume name already exists" | Duplicate name in region | Choose different name |
| "401 Unauthorized" | Invalid Civo token | Verify token, update config |
| "region not found" | Invalid region enum | Use valid region (LON1, NYC1, etc.) |
| "size exceeds limit" | Size > 16,000 GiB | Reduce size or use multiple volumes |
| "state locked" | Concurrent `pulumi up` | Wait for other operation to finish |

## Related Resources

- **Civo Provider Docs**: [pulumi.com/registry/packages/civo](https://www.pulumi.com/registry/packages/civo/)
- **Pulumi Go SDK**: [pulumi.com/docs/reference/pkg/go](https://www.pulumi.com/docs/reference/pkg/go/)
- **Civo API Reference**: [civo.com/api/volumes](https://www.civo.com/api/volumes)

## Conclusion

The CivoVolume Pulumi module provides a production-ready implementation for provisioning block storage on Civo. It handles the 80% use case (name, region, size) while gracefully logging informational messages for features not yet supported by the Civo provider. The module's clear handling of provider limitations ensures users understand what's automated and what requires manual post-deployment steps.

For most teams, this module eliminates the need to write custom Pulumi code or manage raw API calls—simply declare what you want in YAML, and the module handles the rest, providing clear guidance for any manual steps required after provisioning.

