# Cloudflare D1 Database - Pulumi Implementation Overview

## Architecture

The Pulumi module for CloudflareD1Database follows OpenMCF's standard architecture for IaC modules:

```
Stack Input (Protobuf) → Module Entry Point → Provider Setup → Resource Creation → Outputs
```

### Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│ Stack Input (CloudflareD1DatabaseStackInput)                            │
│ ┌─────────────────────────┐  ┌──────────────────────────────────────┐  │
│ │ ProviderConfig          │  │ CloudflareD1Database                 │  │
│ │ - credential_id         │  │ - metadata (name, labels)            │  │
│ │                         │  │ - spec (account_id, database_name,  │  │
│ └─────────────────────────┘  │         region, read_replication)    │  │
│                               └──────────────────────────────────────┘  │
└──────────────────────┬──────────────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────────────┐
│ Module Entry Point (module.Resources)                                   │
│ 1. Initialize locals (copy stack input)                                 │
│ 2. Setup Cloudflare provider (with credential)                          │
│ 3. Create D1 database resource                                          │
└──────────────────────┬──────────────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────────────┐
│ database() Function                                                      │
│ 1. Build D1DatabaseArgs from spec                                       │
│ 2. Map region enum → string (if specified)                              │
│ 3. Add read replication config (if specified)                           │
│ 4. Call cloudflare.NewD1Database()                                      │
│ 5. Export outputs (database-id, database-name, connection-string)       │
└──────────────────────┬──────────────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────────────┐
│ Cloudflare API                                                           │
│ POST /accounts/{account_id}/d1/database                                 │
│ - Creates D1 database                                                    │
│ - Returns database ID                                                    │
└─────────────────────────────────────────────────────────────────────────┘
```

## Design Decisions

### 1. Why Pulumi (Go)?

**Decision**: Use Pulumi with Go SDK instead of Terraform.

**Rationale**:
- **Language Consistency**: OpenMCF's backend is Go-based. Using Pulumi Go keeps the entire stack in one language.
- **Type Safety**: Go's strong typing catches errors at compile time. Protobuf structs map naturally to Go structs.
- **Expressiveness**: Go's conditionals and functions make complex logic (like region mapping) cleaner than HCL.
- **Future-Proofing**: Easier to add custom orchestration (e.g., pairing `wrangler d1 migrations apply` with resource provisioning) in Go than via Terraform external data sources.

**Trade-offs**:
- Smaller community than Terraform (fewer public examples).
- Requires Go runtime (more setup overhead than Terraform's single binary).

**Verdict**: For OpenMCF's use case, the benefits outweigh the trade-offs.

---

### 2. Module Structure: Separate Files

**Decision**: Split module logic into separate files (`main.go`, `locals.go`, `database.go`, `outputs.go`).

**Rationale**:
- **Readability**: Each file has a single responsibility. `database.go` contains only D1 resource creation logic.
- **Maintainability**: Easier to update or debug when logic is isolated.
- **Consistency**: Mirrors OpenMCF's standard pattern for all IaC modules.

**Trade-offs**:
- Slightly more files to navigate.
- Requires jumping between files to follow flow.

**Verdict**: Standard pattern across OpenMCF. Consistency is more valuable than file count.

---

### 3. Region Mapping: Enum to String

**Decision**: Convert CloudflareD1Region enum (proto) to string (Pulumi API) via explicit `mapRegionToString` function.

**Rationale**:
- **Type Safety**: Proto enums are strongly typed. Pulumi expects strings. Explicit mapping catches typos.
- **Clarity**: A `switch` statement documents all valid region values in one place.
- **Flexibility**: Easy to add new regions or rename values if Cloudflare changes API.

**Trade-offs**:
- Adds a helper function (slightly more code).
- Must be updated if Cloudflare adds new regions.

**Verdict**: Explicit is better than implicit. The clarity is worth the extra function.

---

### 4. Optional Fields: Nil Checks

**Decision**: Use `if != nil` checks for optional fields (region, read_replication) instead of always passing them.

**Rationale**:
- **API Defaults**: If `region` is unspecified, Cloudflare selects a default. Passing an empty string might override this behavior.
- **Cleaner Diffs**: Omitting optional fields produces cleaner Pulumi diffs (only shows what's explicitly configured).
- **Protobuf Semantics**: Proto3 uses "optional" keyword. Respecting that in Go code maintains consistency.

**Trade-offs**:
- Slightly more verbose code (if checks instead of always passing).
- Must remember to check for `nil` or `unspecified` enum values.

**Verdict**: Respects API semantics. Cleaner than passing empty/default values.

---

### 5. Connection String: Empty Export

**Decision**: Export `connection-string` as an empty string.

**Rationale**:
- **Protobuf Schema Requirement**: The `CloudflareD1DatabaseStackOutputs` protobuf message defines a `connection_string` field.
- **API Reality**: Cloudflare's D1 API does not expose a connection string (D1 uses Worker bindings, not connection strings).
- **Provider Limitation**: Pulumi's Cloudflare provider (v6.4.1) does not return a connection string output.

**Trade-offs**:
- Confusing for users who expect a connection string.
- Must document that D1 doesn't use connection strings.

**Verdict**: Protobuf schema must be satisfied. Document clearly that D1 uses bindings instead.

---

### 6. Error Handling: Descriptive Wrapping

**Decision**: Use `errors.Wrap` from `github.com/pkg/errors` for all error returns.

**Rationale**:
- **Context**: `errors.Wrap` adds descriptive context to each error (e.g., "failed to create cloudflare d1 database").
- **Debugging**: Stack traces and wrapped messages help debug failures in CI/CD.
- **Consistency**: Standard pattern across OpenMCF modules.

**Trade-offs**:
- Slightly more verbose than `fmt.Errorf`.
- Requires external dependency.

**Verdict**: Better error messages are worth the dependency.

---

### 7. No Unit Tests in Module

**Decision**: Do not include unit tests in `iac/pulumi/module/`.

**Rationale**:
- **Pulumi Testing Challenges**: Mocking Pulumi's context and providers is complex. Most Pulumi modules are integration-tested, not unit-tested.
- **Component-Level Tests**: The `v1/spec_test.go` file validates the protobuf spec and validation rules. This is more valuable than mocking Pulumi internals.
- **Integration Testing**: Real deployments (via `planton apply`) test the full stack.

**Trade-offs**:
- No automated tests for module logic (e.g., region mapping).
- Relies on integration tests to catch bugs.

**Verdict**: Integration tests are sufficient. Unit-testing Pulumi modules has diminishing returns.

---

## Key Implementation Patterns

### Pattern 1: Locals Initialization

**Purpose**: Centralize stack input data for easy access throughout the module.

**Implementation**:

```go
type Locals struct {
	CloudflareProviderConfig *pulumicloudflareprovider.CloudflareCredential
	CloudflareD1Database     *cloudflared1databasev1.CloudflareD1Database
}

func initializeLocals(ctx *pulumi.Context, stackInput *cloudflared1databasev1.CloudflareD1DatabaseStackInput) *Locals {
	return &Locals{
		CloudflareProviderConfig: stackInput.ProviderConfig,
		CloudflareD1Database:     stackInput.Target,
	}
}
```

**Why**: Avoids passing `stackInput` to every function. Cleaner function signatures.

---

### Pattern 2: Provider Abstraction

**Purpose**: Reuse provider setup logic across all Cloudflare components.

**Implementation**:

```go
cloudflareProvider, err := pulumicloudflareprovider.Get(ctx, stackInput.ProviderConfig)
```

**Why**: `pulumicloudflareprovider.Get` is a shared helper that:
- Loads credentials from OpenMCF's credential store
- Configures the Cloudflare provider with API token
- Returns a `*cloudflare.Provider` ready to use

This keeps module code focused on resource creation, not credential plumbing.

---

### Pattern 3: Conditional Field Addition

**Purpose**: Only pass optional fields to Pulumi if they're explicitly set.

**Implementation**:

```go
// Only add region if specified
if locals.CloudflareD1Database.Spec.Region != cloudflared1databasev1.CloudflareD1Region_cloudflare_d1_region_unspecified {
	regionStr := mapRegionToString(locals.CloudflareD1Database.Spec.Region)
	if regionStr != "" {
		d1Args.PrimaryLocationHint = pulumi.String(regionStr)
	}
}

// Only add read replication if specified
if locals.CloudflareD1Database.Spec.ReadReplication != nil {
	d1Args.ReadReplication = &cloudflare.D1DatabaseReadReplicationArgs{
		Mode: pulumi.String(locals.CloudflareD1Database.Spec.ReadReplication.Mode),
	}
}
```

**Why**: Respects API defaults. Produces cleaner Pulumi diffs.

---

## Comparison to Terraform

| Aspect | Pulumi (Go) | Terraform (HCL) |
|--------|-------------|-----------------|
| **Language** | Go (strongly typed) | HCL (declarative DSL) |
| **Conditionals** | Native Go `if` statements | `count`, `for_each`, `dynamic` blocks |
| **Region Mapping** | `switch` statement in `mapRegionToString()` | `locals` with `lookup()` or `terraform.tfvars` |
| **Error Handling** | `errors.Wrap` with stack traces | `terraform validate`, runtime errors |
| **Testing** | Integration tests via `planton apply` | `terraform plan`, integration tests |
| **State Management** | Pulumi state (file, S3, Pulumi Cloud) | Terraform state (file, S3, Terraform Cloud) |
| **Provider Version** | `pulumi-cloudflare/sdk/v6` (bridged from Terraform) | `cloudflare/cloudflare ~> 4.0` |

**Key Insight**: Both produce identical infrastructure. The choice is about team preference and existing tooling.

---

## When to Use Pulumi vs. Terraform

### Use Pulumi (Go) When:
- ✅ Your team prefers writing infrastructure in a real programming language
- ✅ You need complex logic (dynamic resource creation, advanced conditionals)
- ✅ You want type safety and compile-time checks
- ✅ You're already using Go for backend services

### Use Terraform (HCL) When:
- ✅ Your team is already familiar with HCL
- ✅ You want the industry-standard IaC tool
- ✅ You prefer declarative syntax over imperative code
- ✅ You want the largest community and ecosystem

### Both Work Equally Well For:
- ✅ Standard D1 provisioning (account_id, database_name, region, read_replication)
- ✅ Multi-environment deployments (dev, preview, prod)
- ✅ CI/CD integration

---

## Schema Management: The Orchestration Gap

### What Pulumi Does

Pulumi provisions the D1 database **resource** (the "container"). It creates the database in your Cloudflare account and exports the database ID.

### What Pulumi Does NOT Do

Pulumi does **not** create tables or manage schema. This is by design.

**Reason**: Cloudflare does not expose a "migrations" API endpoint. The only way to apply schema is via the Wrangler CLI, which implements client-side migration logic (reading `.sql` files, tracking applied migrations, executing transactionally).

### The Hybrid Workflow

Production deployments require orchestrating IaC (Pulumi) and Wrangler CLI:

```bash
# Step 1: Provision database resource (Pulumi)
planton apply -f database.yaml

# Step 2: Apply schema migrations (Wrangler)
npx wrangler d1 migrations apply my-app-db --remote

# Step 3: Deploy Worker (Wrangler)
npx wrangler deploy
```

**Why This Matters**: You cannot provision a "complete" D1 database (with tables) using pure IaC. The Orchestration Gap is unavoidable.

---

## Future Enhancements

### 1. Wrangler Integration (Potential)

**Idea**: Use Pulumi's `local.Command` resource to shell out to `wrangler d1 migrations apply` as part of the Pulumi stack.

**Pros**: Single command (`pulumi up`) provisions database and applies schema.

**Cons**:
- Non-idempotent (Wrangler modifies state outside Pulumi's tracking).
- Brittle (relies on Wrangler binary being available).
- Debugging complexity (errors buried in command output).

**Verdict**: Not recommended. Keep IaC and schema management separate for clarity.

---

### 2. Worker Binding Support

**Idea**: Add Worker binding configuration to the CloudflareD1DatabaseSpec.

**Pros**: Single resource provisions database and binds it to a Worker.

**Cons**:
- Tight coupling (database lifecycle tied to Worker lifecycle).
- Multi-Worker scenarios (one database, many Workers) become complex.

**Verdict**: Bindings are better managed in the Worker's configuration (`wrangler.toml` or separate Terraform `cloudflare_workers_script` resource).

---

## Summary

The Pulumi module for CloudflareD1Database is a thin, clean wrapper around Cloudflare's D1 API. It:
- ✅ Translates protobuf specs to Pulumi resources
- ✅ Handles optional fields (region, read_replication) correctly
- ✅ Maps enums to strings explicitly
- ✅ Exports structured outputs for use by Workers

It does **not**:
- ❌ Manage schema (use Wrangler CLI)
- ❌ Configure Worker bindings (use `wrangler.toml`)

For detailed architectural guidance, see [../../docs/README.md](../../docs/README.md).

