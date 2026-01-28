# Civo Database Pulumi Module - Architecture Overview

## Purpose

This document provides an architectural overview of the Pulumi module for Civo Database provisioning within the OpenMCF ecosystem.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     OpenMCF CLI                          │
│                  (Orchestration Layer)                           │
└─────────────────────────────┬───────────────────────────────────┘
                              │
                              │ CivoDatabaseStackInput (protobuf)
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                   Pulumi Entrypoint                              │
│                   (iac/pulumi/main.go)                           │
│                                                                  │
│  • Unmarshals stack input                                       │
│  • Calls module.Resources()                                     │
│  • Handles errors                                               │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                 Module Entry Point                               │
│               (iac/pulumi/module/main.go)                        │
│                                                                  │
│  • Initializes locals from stack input                          │
│  • Configures Civo provider                                     │
│  • Orchestrates resource creation                               │
└─────────┬───────────────────┬───────────────────────────────────┘
          │                   │
          │                   └──────────────────────┐
          │                                          │
┌─────────▼──────────────────┐         ┌───────────▼────────────┐
│   Provider Configuration   │         │  Resource Provisioning │
│ (pulumicivoprovider.Get()) │         │  (module/database.go)  │
│                            │         │                        │
│ • Civo API token           │         │ • Database instance    │
│ • Region configuration     │         │ • Network attachment   │
└────────────────────────────┘         │ • Firewall rules       │
                                       │ • Storage config       │
                                       │ • HA replicas          │
                                       │ • Tag management       │
                                       └───────────┬────────────┘
                                                   │
                                       ┌───────────▼────────────┐
                                       │   Civo API             │
                                       │   (pulumi-civo)        │
                                       │                        │
                                       │ • CreateDatabase       │
                                       │ • AttachNetwork        │
                                       │ • AttachFirewall       │
                                       └───────────┬────────────┘
                                                   │
                                       ┌───────────▼────────────┐
                                       │   Outputs              │
                                       │                        │
                                       │ • database_id          │
                                       │ • host (DNS endpoint)  │
                                       │ • port                 │
                                       │ • username             │
                                       │ • password (sensitive) │
                                       └────────────────────────┘
```

## Component Responsibilities

### 1. Pulumi Entrypoint (`main.go`)

**Purpose**: Thin wrapper that serves as the Pulumi program entry point.

**Responsibilities**:
- Parse command-line flags
- Unmarshal stack input from JSON
- Invoke `module.Resources()`
- Handle top-level errors

**Key Code**:
```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        stackInput := &civodatabasev1.CivoDatabaseStackInput{}
        // ... unmarshal logic ...
        return module.Resources(ctx, stackInput)
    })
}
```

**Why separate**: Decouples Pulumi runtime concerns from business logic, making the module testable.

---

### 2. Module Entry (`module/main.go`)

**Purpose**: Primary orchestration point for resource provisioning.

**Responsibilities**:
- Initialize `Locals` struct with stack input
- Configure Civo provider from credentials
- Call `database()` to provision resources
- Return errors to entrypoint

**Key Code**:
```go
func Resources(ctx *pulumi.Context, stackInput *civodatabasev1.CivoDatabaseStackInput) error {
    locals := initializeLocals(ctx, stackInput)
    civoProvider, err := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
    if err != nil {
        return errors.Wrap(err, "failed to setup Civo provider")
    }
    if _, err := database(ctx, locals, civoProvider); err != nil {
        return errors.Wrap(err, "failed to create Civo database")
    }
    return nil
}
```

**Why separate**: Centralizes orchestration logic and provider setup, making the module easier to extend.

---

### 3. Locals Initialization (`module/locals.go`)

**Purpose**: Consolidate frequently used values and derive computed values.

**Responsibilities**:
- Store references to provider config and spec
- Derive Civo-compatible tags from spec
- Prepare labels for potential future use

**Key Code**:
```go
type Locals struct {
    CivoProviderConfig *civoprovider.CivoProviderConfig
    CivoDatabase       *civodatabasev1.CivoDatabase
    CivoTags           pulumi.StringArray
    CivoLabels         map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *civodatabasev1.CivoDatabaseStackInput) *Locals {
    // ... initialization logic ...
}
```

**Why separate**: Avoids recomputing values and keeps business logic DRY.

---

### 4. Database Provisioning (`module/database.go`)

**Purpose**: Core resource provisioning logic.

**Responsibilities**:
- Translate proto enum values to Civo-compatible strings
- Map spec fields to `civo.DatabaseArgs`
- Handle optional fields (firewall, storage, tags)
- Calculate total nodes (primary + replicas)
- Create `civo.Database` resource
- Export outputs

**Key Code**:
```go
func database(ctx *pulumi.Context, locals *Locals, civoProvider *civo.Provider) (*civo.Database, error) {
    // Translate engine
    var engineSlug string
    switch locals.CivoDatabase.Spec.Engine {
    case civodatabasev1.CivoDatabaseEngine_mysql:
        engineSlug = "mysql"
    case civodatabasev1.CivoDatabaseEngine_postgres:
        engineSlug = "postgres"
    default:
        return nil, errors.Errorf("unsupported database engine: %v", locals.CivoDatabase.Spec.Engine)
    }

    // Build args
    databaseArgs := &civo.DatabaseArgs{
        Name:    pulumi.String(locals.CivoDatabase.Spec.DbInstanceName),
        Engine:  pulumi.String(engineSlug),
        Version: pulumi.String(locals.CivoDatabase.Spec.EngineVersion),
        Region:  pulumi.String(locals.CivoDatabase.Spec.Region.String()),
        Size:    pulumi.String(locals.CivoDatabase.Spec.SizeSlug),
        Nodes:   pulumi.Int(int(locals.CivoDatabase.Spec.Replicas) + 1),
    }

    // ... attach network, firewall, storage, tags ...

    // Provision
    createdDatabase, err := civo.NewDatabase(ctx, "database", databaseArgs, pulumi.Provider(civoProvider))
    
    // Export outputs
    ctx.Export(OpDatabaseId, createdDatabase.ID())
    ctx.Export(OpHost, createdDatabase.DnsEndpoint)
    // ... more exports ...

    return createdDatabase, nil
}
```

**Why separate**: Encapsulates all database-specific logic in one file, making it easy to locate and modify.

---

### 5. Output Constants (`module/outputs.go`)

**Purpose**: Define standardized output keys for cross-resource wiring.

**Responsibilities**:
- Declare constant strings for output keys
- Ensure consistency across OpenMCF

**Key Code**:
```go
const (
    OpDatabaseId        = "database_id"
    OpHost              = "host"
    OpPort              = "port"
    OpUsername          = "username"
    OpPasswordSecretRef = "password"
)
```

**Why separate**: Centralizes output naming convention, preventing typos and inconsistencies.

---

## Data Flow

### Input → Processing → Output

```
Stack Input (JSON/Protobuf)
│
├─ provider_config
│  └─ civo_token: "abc123..."
│
└─ target
   ├─ metadata
   │  ├─ name: "prod-db"
   │  └─ env: "production"
   │
   └─ spec
      ├─ db_instance_name: "production-db"
      ├─ engine: postgres (enum value 2)
      ├─ engine_version: "16"
      ├─ region: lon1 (enum value)
      ├─ size_slug: "g3.db.large"
      ├─ replicas: 2
      ├─ network_id: { value: "net-12345678" }
      ├─ firewall_ids: [{ value: "fw-87654321" }]
      ├─ storage_gib: 200
      └─ tags: ["production", "backend"]

      ↓ (module/main.go: initializeLocals)

Locals Struct
│
├─ CivoProviderConfig: { token: "abc123..." }
├─ CivoDatabase: { spec: { ... } }
├─ CivoTags: ["production", "backend"]
└─ CivoLabels: { "planton:resource": "true", ... }

      ↓ (module/database.go: translate & map)

civo.DatabaseArgs
│
├─ Name: "production-db"
├─ Engine: "postgres"  ← translated from enum
├─ Version: "16"
├─ Region: "lon1"  ← converted from enum to string
├─ Size: "g3.db.large"
├─ Nodes: 3  ← computed (replicas + 1)
├─ NetworkId: "net-12345678"
├─ FirewallId: "fw-87654321"  ← first element
├─ SizeGb: 200
└─ Tags: ["production", "backend"]

      ↓ (pulumi-civo: civo.NewDatabase)

Civo API Call
│
└─ POST /v2/databases
   {
     "name": "production-db",
     "engine": "postgres",
     "version": "16",
     ...
   }

      ↓ (Civo API response)

civo.Database Resource
│
├─ id: "db-abc123..."
├─ dns_endpoint: "db-abc123.civo.com"
├─ endpoint: "10.0.0.50"
├─ port: 5432
├─ username: "civo"
├─ password: "SecurePass123!"
└─ status: "Active"

      ↓ (module/database.go: ctx.Export)

Pulumi Outputs
│
├─ database_id: "db-abc123..."
├─ host: "db-abc123.civo.com"
├─ port: 5432
├─ username: "civo"  (sensitive)
└─ password: "SecurePass123!"  (sensitive)
```

---

## Key Design Decisions

### 1. Why Separate `database.go` from `main.go`?

**Rationale**: Separation of concerns. `main.go` handles orchestration, while `database.go` focuses purely on database resource logic.

**Benefit**: Easy to extend. If we need to add additional resources (e.g., monitoring), we create `monitoring.go` and call it from `main.go`.

---

### 2. Why Use `Locals` Struct?

**Rationale**: Avoids passing multiple parameters to every function and centralizes derived values.

**Benefit**: Reduces function signatures and makes code DRY. Tags and labels are computed once and reused.

---

### 3. Why Translate Enums to Strings?

**Rationale**: Protobuf enums are integers. Civo API expects string literals (`"mysql"`, `"postgres"`).

**Benefit**: Type safety in the protobuf spec, while maintaining API compatibility.

**Example**:
```go
// Proto: CivoDatabaseEngine_postgres = 2
// Civo API expects: "postgres"
switch locals.CivoDatabase.Spec.Engine {
case civodatabasev1.CivoDatabaseEngine_postgres:
    engineSlug = "postgres"
}
```

---

### 4. Why Compute `Nodes` as `Replicas + 1`?

**Rationale**: Civo API expects **total node count** (primary + replicas), but the protobuf spec uses `replicas` (just the replica count).

**Benefit**: Matches intuitive user expectation. Users specify "I want 2 replicas" (meaning 3 total nodes), not "I want 3 nodes."

**Example**:
```go
// User specifies: replicas: 2
// Civo expects: nodes: 3
databaseArgs.Nodes = pulumi.Int(int(locals.CivoDatabase.Spec.Replicas) + 1)
```

---

### 5. Why Export `dns_endpoint` as `host`?

**Rationale**: The DNS endpoint is the **recommended connection string** for HA configurations. It automatically updates during failover.

**Benefit**: Applications using the `host` output get automatic failover for free.

**Alternative**: We also export the static `endpoint`, but DNS endpoint is the primary output.

---

### 6. Why Mark Password as Sensitive?

**Rationale**: Pulumi automatically handles sensitive values by encrypting them in state and masking them in logs.

**Benefit**: Prevents accidental credential leakage.

**Implementation**: Pulumi's `pulumi-civo` provider marks the password as sensitive automatically. We just export it:
```go
ctx.Export(OpPasswordSecretRef, createdDatabase.Password)
```

---

## Error Handling Strategy

### Explicit Validation

The module validates inputs early and returns descriptive errors:

```go
if locals.CivoDatabase.Spec.Engine == civodatabasev1.CivoDatabaseEngine_civo_database_engine_unspecified {
    return nil, errors.New("engine must be specified")
}
```

### Error Wrapping

Errors are wrapped with context to aid debugging:

```go
if err != nil {
    return errors.Wrap(err, "failed to create Civo database")
}
```

### Fail-Fast Philosophy

The module stops at the first error. Pulumi handles rollback automatically.

---

## Testing Strategy

### Unit Tests

Test individual functions in isolation:

```go
func TestEngineTranslation(t *testing.T) {
    // Test that CivoDatabaseEngine_postgres translates to "postgres"
}
```

### Integration Tests

Deploy actual resources to Civo and verify outputs:

```bash
export CIVO_TOKEN="test-token"
pulumi stack init test
pulumi up
pulumi stack output host  # Verify output exists
pulumi destroy
```

### Validation Tests

Verify that invalid inputs are rejected:

```go
func TestInvalidEngine(t *testing.T) {
    // Test that unspecified engine returns error
}
```

---

## Extension Points

### Adding New Fields

To add a new field from `CivoDatabaseSpec`:

1. **Update `spec.proto`** (if field doesn't exist)
2. **Regenerate Go stubs**: `make protos`
3. **Map field in `database.go`**:
   ```go
   if locals.CivoDatabase.Spec.NewField != "" {
       databaseArgs.NewField = pulumi.String(locals.CivoDatabase.Spec.NewField)
   }
   ```
4. **Test the change**
5. **Update documentation**

### Adding New Resources

To add a new related resource (e.g., database user):

1. **Create `module/user.go`**:
   ```go
   func user(ctx *pulumi.Context, locals *Locals, db *civo.Database) (*civo.DatabaseUser, error) {
       // ...
   }
   ```
2. **Call from `module/main.go`**:
   ```go
   db, err := database(ctx, locals, civoProvider)
   if err != nil { return err }
   
   user, err := user(ctx, locals, db)
   if err != nil { return err }
   ```

---

## Performance Considerations

### Sequential Provisioning

Resources are created sequentially. This is intentional:
- Database must exist before we can create users
- Network must exist before we can attach firewall
- Pulumi handles parallelization where safe

### State Management

Pulumi stores state remotely (e.g., Pulumi Cloud, S3). This ensures:
- No local state files
- Concurrent safety
- Audit trail

---

## Security Considerations

### Credential Handling

1. **API Token**: Provided via environment variable (`CIVO_TOKEN`), never hardcoded
2. **Database Password**: Marked as sensitive, encrypted in state
3. **Network Isolation**: Database always deployed in private network

### Least Privilege

Firewall rules should follow least-privilege principle:
```go
// Only allow traffic from Kubernetes node CIDR
databaseArgs.FirewallId = pulumi.String(firewallId)
```

---

## Debugging Tips

### Enable Verbose Logging

```bash
pulumi up --logtostderr --logflow -v=9
```

### Inspect Pulumi State

```bash
pulumi stack export > state.json
cat state.json | jq '.deployment.resources[] | select(.type == "civo:index:Database")'
```

### Check Civo API Calls

```bash
# Enable Civo SDK debug logging
export CIVO_DEBUG=1
pulumi up
```

---

## Additional Resources

- **Pulumi Civo Provider Docs**: [pulumi.com/docs/clouds/civo](https://www.pulumi.com/docs/clouds/civo/)
- **Civo API Reference**: [civo.com/api/databases](https://www.civo.com/api/databases)
- **Parent README**: [`../../README.md`](../../README.md)
- **Module README**: [`README.md`](README.md)
- **Examples**: [`examples.md`](examples.md)

