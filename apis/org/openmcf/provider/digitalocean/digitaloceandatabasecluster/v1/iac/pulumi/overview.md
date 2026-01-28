# DigitalOcean Database Cluster - Pulumi Module Architecture

## Module Architecture

The module follows OpenMCF's standard pattern for infrastructure provisioning:

1. **Entry Point** (`main.go`) - Parses stack input, invokes module
2. **Locals Initialization** (`module/locals.go`) - Extracts metadata, labels, credentials
3. **Provider Setup** - Configures DigitalOcean provider from credentials
4. **Resource Creation** (`module/database_cluster.go`) - Provisions database cluster
5. **Output Export** - Exports connection details and cluster metadata

---

## Data Flow

### 1. Input Processing

```protobuf
message DigitalOceanDatabaseClusterStackInput {
  ProviderConfig provider_config = 1;
  DigitalOceanDatabaseCluster target = 2;
}
```

### 2. Resource Creation

```go
clusterArgs := &digitalocean.DatabaseClusterArgs{
    Name:     pulumi.String(spec.ClusterName),
    Engine:   pulumi.String(engineSlug),
    Version:  pulumi.String(spec.EngineVersion),
    Size:     pulumi.String(spec.SizeSlug),
    Region:   pulumi.String(regionSlug),
    NodeCount: pulumi.Int(int(spec.NodeCount)),
}

// Optional fields
if spec.Vpc != nil && spec.Vpc.GetValue() != "" {
    clusterArgs.PrivateNetworkUuid = pulumi.StringPtr(spec.Vpc.GetValue())
}

if spec.StorageGib > 0 {
    clusterArgs.StorageSizeMib = pulumi.IntPtr(int(spec.StorageGib) * 1024)
}
```

### 3. Output Export

```go
ctx.Export(OpClusterId, cluster.ID())
ctx.Export(OpConnectionUri, pulumi.ToSecret(cluster.Uri))
ctx.Export(OpPassword, pulumi.ToSecret(cluster.Password))
```

Sensitive outputs are wrapped in `pulumi.ToSecret()` for encryption.

---

## Engine Slug Mapping

DigitalOcean uses non-standard engine slugs:

| Protobuf Enum | DigitalOcean API Slug | Notes |
|---------------|----------------------|-------|
| `postgres` | `pg` | Abbreviated |
| `mysql` | `mysql` | Direct match |
| `redis` | `redis` | Redis 7+ = Valkey |
| `mongodb` | `mongodb` | Direct match |

---

## Critical Design Patterns

### 1. Optional VPC Handling

```go
var vpcUuid *string
if spec.Vpc != nil {
    uuid := spec.Vpc.GetValue()
    if uuid != "" {
        vpcUuid = &uuid
    }
}

if vpcUuid != nil {
    clusterArgs.PrivateNetworkUuid = pulumi.StringPtr(*vpcUuid)
}
```

**Why**: VPC is optional. If not specified, DigitalOcean creates cluster outside VPC (deprecated pattern).

### 2. Storage Size Conversion

DigitalOcean API expects storage in **MiB**, protobuf spec uses **GiB**:

```go
storageMib := int(spec.StorageGib) * 1024
clusterArgs.StorageSizeMib = pulumi.IntPtr(storageMib)
```

### 3. Sensitive Output Handling

```go
ctx.Export(OpPassword, pulumi.ToSecret(cluster.Password))
ctx.Export(OpConnectionUri, pulumi.ToSecret(cluster.Uri))
```

**Why**: Prevents credentials from appearing in plain text in logs or console output.

---

## Limitations and Workarounds

### What the Module Handles

✅ Cluster provisioning (name, engine, version, size, nodes)  
✅ VPC integration  
✅ Custom storage sizing  
✅ Public/private connectivity configuration

### What Requires Separate Resources

❌ Firewall rules (`DigitalOceanDatabaseFirewall`)  
❌ Connection pools (`DigitalOceanDatabaseConnectionPool`)  
❌ Additional users (`DigitalOceanDatabaseUser`)  
❌ Additional databases (`DigitalOceanDatabaseDb`)  
❌ Read replicas (`DigitalOceanDatabaseReplica`)

**Rationale**: These have independent lifecycles and should be managed separately.

---

## Production Best Practices

### High Availability

```go
if spec.NodeCount < 2 {
    log.Warning("Single-node cluster is not HA; suitable for dev/test only")
}
```

**Recommendation**: Enforce `node_count ≥ 2` for production via policy-as-code.

### PostgreSQL Connection Pooling

**Critical**: PostgreSQL connection limits are severely restricted. The module exports connection details but **does not create PgBouncer pools**.

**Solution**: Deploy separate connection pool resource:

```go
pool := digitalocean.NewDatabaseConnectionPool(ctx, "pool", &digitalocean.DatabaseConnectionPoolArgs{
    ClusterId: cluster.ID(),
    Name:      pulumi.String("app-pool"),
    Mode:      pulumi.String("transaction"),
    Size:      pulumi.Int(25),
})
```

### Secret Management

Credentials are exported as sensitive outputs. For production:

1. **Pulumi Secrets**: Already encrypted in state
2. **Vault Integration**: Store in HashiCorp Vault post-deployment
3. **K8s Secrets**: Create Kubernetes secret from outputs for app consumption

---

## Testing and Validation

### Pre-Deployment Validation

```go
// Validate node_count is in allowed range (1-3)
if spec.NodeCount < 1 || spec.NodeCount > 3 {
    return errors.New("node_count must be between 1 and 3")
}
```

### Post-Deployment Verification

```bash
# Verify cluster is accessible
psql "$(pulumi stack output --show-secrets connection_uri)" -c "SELECT version();"

# Check cluster status
doctl databases get $(pulumi stack output cluster_id)
```

---

## Further Reading

- **DigitalOcean API Docs**: [Databases API](https://docs.digitalocean.com/reference/api/api-reference/#tag/Databases)
- **Pulumi DigitalOcean Provider**: [pulumi-digitalocean](https://www.pulumi.com/registry/packages/digitalocean/)
- **Component Docs**: [../../README.md](../../README.md)
