# Pulumi Module Architecture: AWS KMS Key

## Overview

This Pulumi module provisions AWS KMS (Key Management Service) customer-managed encryption keys through a declarative, protobuf-defined specification. KMS keys provide cryptographic operations for data encryption, digital signatures, and key generation with automatic rotation, fine-grained access control, and comprehensive audit logging.

The module follows OpenMCF's standard pattern: **input transformation → resource provisioning → output extraction**, focusing on the 80/20 configuration that covers most production KMS use cases while maintaining simplicity.

## Module Structure

```
iac/pulumi/
├── main.go              # Pulumi entrypoint
├── Pulumi.yaml          # Pulumi project metadata
├── Makefile             # Build/test helpers
├── debug.sh             # Delve debugging wrapper
└── module/
    ├── main.go          # Orchestration logic (provider setup, output export)
    ├── locals.go        # Input transformation and tag construction
    ├── kms.go           # KMS key and alias resource implementation
    └── outputs.go       # Output key constants
```

### File Responsibilities

#### `main.go` (entrypoint)
- Unmarshals `AwsKmsKeyStackInput` from Pulumi stack configuration
- Delegates to `module.Resources()` for actual provisioning
- Minimal logic—purely a thin entrypoint wrapper

#### `module/main.go` (orchestrator)
**Key Function:** `Resources(ctx, stackInput)`

Responsibilities:
1. **Locals Initialization**: Transform stack input into typed locals struct
2. **Provider Configuration**: Handle two provider scenarios:
   - **Default**: Create provider with ambient AWS credentials (IAM role, environment variables)
   - **Explicit**: Create provider with credentials from `stackInput.ProviderConfig` (access key, secret key, session token)
3. **Resource Creation**: Invoke `kmsKey()` with initialized locals and provider
4. **Output Export**: Export four KMS key outputs (key_id, key_arn, alias_name, rotation_enabled)

**Design Decision**: Provider handling supports both CI/CD environments (IRSA, instance profiles) and local development (explicit credentials).

#### `module/locals.go` (input transformer)
**Key Function:** `initializeLocals(ctx, stackInput)`

Transforms the protobuf `AwsKmsKeyStackInput` into a strongly-typed `Locals` struct and constructs AWS resource tags.

**Tag Construction**: Every KMS key receives five mandatory tags:
- `resource=true`: Marks this as a OpenMCF managed resource
- `organization=<org>`: Organization ID from metadata
- `environment=<env>`: Environment (dev/staging/prod) from metadata
- `resource-kind=AwsKmsKey`: CloudResourceKind enum string
- `resource-id=<id>`: Unique resource identifier from metadata

**Why Tags Matter**: These tags enable:
- Cost allocation reporting by org/env/resource-kind ($1/month per key adds up)
- Policy enforcement (e.g., "only keys tagged env=prod can encrypt production data")
- Resource discovery and inventory management
- Drift detection (manually created keys lack these tags)

#### `module/kms.go` (resource implementation)
**Key Function:** `kmsKey(ctx, locals, provider)`

This is the core implementation file containing all KMS key provisioning logic. It translates the protobuf `AwsKmsKeySpec` into Pulumi's AWS KMS resources.

**Implementation Pattern**: The module follows a **two-phase provisioning** approach:

1. **Phase 1: Key Creation**
   - Map protobuf enum `key_spec` to AWS key spec strings
   - Create the KMS key with description, rotation, deletion window, and tags
   - Handle key type variations (symmetric, RSA, ECC)

2. **Phase 2: Alias Creation (Conditional)**
   - If `alias_name` is provided, create a KMS alias resource
   - Link alias to the key ID
   - If no alias provided, export empty string

**Key Type Mapping**: Explicit enum-to-string conversion ensures type safety:
- `0 (symmetric)` → `"SYMMETRIC_DEFAULT"`
- `1 (rsa_2048)` → `"RSA_2048"`
- `2 (rsa_4096)` → `"RSA_4096"`
- `3 (ecc_nist_p256)` → `"ECC_NIST_P256"`

#### `module/outputs.go` (constants)
Defines string constants for output keys, preventing typos and providing a single source of truth.

**Exported Outputs**:
- `key_id`: Unique identifier for the KMS key
- `key_arn`: Amazon Resource Name (ARN) for policies and cross-account access
- `alias_name`: Human-readable alias (if provided)
- `rotation_enabled`: Whether automatic rotation is enabled (boolean)

## Data Flow Diagram

```
┌──────────────────────────────────┐
│ AwsKmsKeyStackInput (protobuf)   │
│  ├─ target: AwsKmsKey            │
│  │   ├─ metadata                 │
│  │   └─ spec: AwsKmsKeySpec      │
│  │       ├─ key_spec             │
│  │       ├─ description          │
│  │       ├─ disable_key_rotation │
│  │       ├─ deletion_window_days │
│  │       └─ alias_name           │
│  └─ provider_config (optional)   │
└──────────────┬───────────────────┘
               │
               ▼
      ┌────────────────┐
      │initializeLocals│
      └────────┬────────┘
               │ Creates:
               │ - Locals.AwsKmsKey
               │ - Locals.AwsTags (map[string]string)
               ▼
     ┌──────────────────────┐
     │ AWS Provider Setup   │
     │  ├─ Ambient creds OR │
     │  └─ Explicit creds   │
     └──────────┬───────────┘
                │
                ▼
       ┌─────────────────┐
       │    kmsKey()     │
       └────────┬────────┘
                │
                ├─ Map key_spec enum to AWS string:
                │    0 → "SYMMETRIC_DEFAULT"
                │    1 → "RSA_2048"
                │    2 → "RSA_4096"
                │    3 → "ECC_NIST_P256"
                │
                ├─ Create KMS Key:
                │    kms.NewKey(ctx, name, &kms.KeyArgs{
                │      Description:           description,
                │      DeletionWindowInDays:  deletionWindowDays,
                │      EnableKeyRotation:     !disable_key_rotation,
                │      CustomerMasterKeySpec: keySpec,
                │      Tags:                  awsTags
                │    })
                │
                ├─ If alias_name provided:
                │    kms.NewAlias(ctx, name+"-alias", &kms.AliasArgs{
                │      Name:        alias_name,
                │      TargetKeyId: key.KeyId
                │    })
                │
                └─ Export outputs:
                     ctx.Export("key_id", key.KeyId)
                     ctx.Export("key_arn", key.Arn)
                     ctx.Export("alias_name", aliasNameOutput)
                     ctx.Export("rotation_enabled", key.EnableKeyRotation)
                │
                ▼
       ┌─────────────────────────┐
       │  AWS KMS Service        │
       │   ├─ KMS Key            │
       │   └─ Alias (optional)   │
       └─────────────────────────┘
```

## Resource Relationships

```
AwsKmsKeySpec
  │
  ├─ key_spec ───────────────────┐
  │    (enum: symmetric, RSA, ECC)│
  │                                │
  ├─ description ─────────────────┼─ KMS Key
  │                                │   ├─ Key ID (UUID)
  ├─ disable_key_rotation ────────┼───┤   ├─ Key ARN
  │    (bool, default: false)     │   │   ├─ Description
  │                                │   │   ├─ Key Spec (SYMMETRIC_DEFAULT, RSA_2048, etc.)
  ├─ deletion_window_days ────────┼───┤   ├─ Deletion Window (7-30 days)
  │    (7-30, recommended: 30)    │   │   ├─ Rotation Enabled (true/false)
  │                                │   │   └─ Tags (org/env/resource metadata)
  └─ alias_name ──────────────────┼───┤
       (optional, pattern: alias/*)│   │
                                   │   └─ Alias (conditional)
                                   │       ├─ Alias Name (alias/*)
                                   │       └─ Target Key ID
                                   │
                                   └─ Outputs:
                                       ├─ key_id
                                       ├─ key_arn
                                       ├─ alias_name
                                       └─ rotation_enabled
```

### Critical Relationships

**Key → Alias (0:1 Relationship)**:
- A key can have zero or one alias (in this module)
- AWS supports multiple aliases per key, but this module creates only one for simplicity
- Aliases provide human-readable names: `alias/prod/app-data` instead of UUID key IDs
- Aliases can be updated to point to different keys (for key rotation)

**Key Type → Rotation Compatibility**:
- **Symmetric keys**: Support automatic annual rotation
- **Asymmetric keys (RSA, ECC)**: Do NOT support automatic rotation
- The module doesn't enforce this—AWS will reject rotation for asymmetric keys at API level
- Best practice: Set `disable_key_rotation: true` for RSA/ECC keys

**Deletion Window → Data Recovery**:
- Keys enter "pending deletion" state for 7-30 days
- During this window, the key cannot be used but can be cancelled
- After the window, the key is permanently deleted and data becomes unrecoverable
- Recommended: Always use 30 days (maximum) for production keys

## Key Design Decisions

### 1. Explicit Enum Mapping Over Reflection

**Decision**: All protobuf enums are explicitly mapped to AWS string constants.

**Implementation**:
```go
var keySpec pulumi.StringPtrInput
switch spec.KeySpec {
case 0:
    keySpec = pulumi.StringPtr("SYMMETRIC_DEFAULT")
case 1:
    keySpec = pulumi.StringPtr("RSA_2048")
// ... other cases
}
```

**Rationale**:
- **Type Safety**: Compiler catches missing enum cases
- **API Versioning**: If AWS adds new key types, code fails to compile, forcing explicit handling
- **Clarity**: No hidden mappings—the transformation is visible
- **Documentation**: Switch statement serves as documentation of supported key types

### 2. Rotation Defaults to Enabled

**Decision**: By default, keys have automatic rotation enabled (`disable_key_rotation: false`).

**Implementation**:
```go
EnableKeyRotation: pulumi.BoolPtr(!spec.DisableKeyRotation)
```

**Rationale**:
- **Security Best Practice**: AWS recommends annual rotation for symmetric keys
- **Compliance**: Many security frameworks require key rotation
- **Explicit Opt-Out**: Users must explicitly set `disable_key_rotation: true` to disable
- **Protobuf Default**: Boolean defaults to false, so rotation is enabled unless explicitly disabled

**Trade-off**: Asymmetric keys don't support rotation, so users must remember to disable it for RSA/ECC keys. Alternative would be to check key type and conditionally enable rotation, but explicit configuration is clearer.

### 3. Optional Alias Creation

**Decision**: Aliases are optional; created only if `alias_name` is provided.

**Implementation**:
```go
if spec.AliasName != "" {
    _, err := kms.NewAlias(ctx, locals.AwsKmsKey.Metadata.Name+"-alias", &kms.AliasArgs{
        Name:        pulumi.String(spec.AliasName),
        TargetKeyId: createdKey.KeyId,
    }, pulumi.Provider(provider))
}
```

**Rationale**:
- **Flexibility**: Some use cases reference keys directly by ARN
- **Multi-Alias Support**: Users might create additional aliases separately
- **Validation**: Protobuf pattern enforces `alias/` prefix if alias is provided
- **Output Safety**: Empty string exported if no alias (not nil, preventing output errors)

**Best Practice**: Always create aliases for production keys. Distributing key IDs makes key rotation impossible.

### 4. Deletion Window Validation at Protobuf Level

**Decision**: Validate deletion window (7-30 days) in protobuf validation, not Go code.

**Protobuf Validation**:
```protobuf
int32 deletion_window_days = 4 [
  (buf.validate.field).int32.gte = 7,
  (buf.validate.field).int32.lte = 30
];
```

**Rationale**:
- **Early Validation**: Errors caught before Pulumi runs
- **Consistent Validation**: Same rules apply in Terraform and other tools
- **Clear Errors**: Protobuf validation messages are descriptive
- **No Redundancy**: Don't duplicate validation in multiple places

### 5. Tag Injection at Resource Creation

**Decision**: Tags are constructed from metadata and attached at key creation (not separate API call).

**Rationale**:
- **Atomicity**: Tags exist from the moment the key is created
- **Consistency**: Every OpenMCF resource has identical tagging structure
- **Cost Allocation**: AWS billing sees tags immediately, no gap in cost reporting
- **No Tag Drift**: Tags can't be missing or inconsistent with metadata

## KMS-Specific Implementation Details

### Key Type Enum Mapping

The module maps protobuf enums to AWS-specific key spec strings.

**Enum Definitions** (spec.proto):
```protobuf
enum AwsKmsKeyType {
  symmetric = 0;
  rsa_2048 = 1;
  rsa_4096 = 2;
  ecc_nist_p256 = 3;
}
```

**AWS Mappings**:
- `symmetric (0)` → `"SYMMETRIC_DEFAULT"` - AES-256-GCM, fastest, supports envelope encryption
- `rsa_2048 (1)` → `"RSA_2048"` - 2048-bit RSA for encryption and signing
- `rsa_4096 (2)` → `"RSA_4096"` - 4096-bit RSA for higher security
- `ecc_nist_p256 (3)` → `"ECC_NIST_P256"` - Elliptic curve for signing (faster than RSA)

**Default Handling**: If enum value is unrecognized (shouldn't happen with validation), defaults to `"SYMMETRIC_DEFAULT"`.

### Rotation Logic Inversion

The protobuf field is `disable_key_rotation` (boolean), but AWS API expects `enable_key_rotation`.

**Implementation**:
```go
EnableKeyRotation: pulumi.BoolPtr(!spec.DisableKeyRotation)
```

**Why Inversion?**
- Protobuf field name makes the default behavior clear: "disable_key_rotation: false" = rotation enabled
- Users explicitly set `true` to disable rotation, making it an opt-in to less secure behavior
- AWS API uses positive phrasing (`enable`), so we invert the boolean

**Protobuf Default**: Boolean fields default to `false`, so rotation is enabled by default (secure-by-default design).

### Conditional Alias Output

If no alias is provided, the module exports an empty string rather than nil.

**Implementation**:
```go
var aliasNameOutput pulumi.StringOutput
if spec.AliasName != "" {
    aliasNameOutput = pulumi.String(spec.AliasName).ToStringOutput()
} else {
    aliasNameOutput = pulumi.String("").ToStringOutput()
}
```

**Rationale**:
- **Consistent Output Type**: Always a string output, never nil
- **Downstream Safety**: Consumers can safely reference the output without nil checks
- **Empty String Convention**: Standard pattern in OpenMCF for optional fields

### Deletion Window Range

AWS enforces 7-30 day deletion window. The protobuf validation reflects this constraint.

**Validation**:
```protobuf
deletion_window_days = 4 [
  (buf.validate.field).int32.gte = 7,
  (buf.validate.field).int32.lte = 30,
  (org.openmcf.shared.options.recommended_default) = "30"
];
```

**Recommended Default**: 30 days (maximum) provides maximum protection against accidental deletion.

**Common Mistake**: Setting 7 days for production keys to "clean up faster." This is dangerous—7 days is insufficient for detecting and reversing accidental deletions.

## Error Handling Philosophy

The module follows a **fail-fast** approach:

1. **Validation at Protobuf Level**: The `spec.proto` validation rules catch configuration errors before Pulumi runs
2. **AWS API Errors Propagate**: If AWS rejects a configuration, the error propagates immediately (no silent fallbacks)
3. **Wrapped Errors**: All errors include context (`errors.Wrap()`) for easier debugging

**Rationale**: Infrastructure as code demands predictability. Silent defaults or error recovery can mask misconfiguration and create surprise behavior.

## Common Pitfalls and Gotchas

### Pitfall 1: Rotation Enabled for Asymmetric Keys
**Symptom**: AWS API error: `InvalidParameterException: Automatic key rotation is not supported for asymmetric keys`

**Cause**: Default rotation setting is enabled, but RSA/ECC keys don't support rotation.

**Solution**: Explicitly set `disable_key_rotation: true` for asymmetric keys (rsa_2048, rsa_4096, ecc_nist_p256).

### Pitfall 2: Alias Without "alias/" Prefix
**Symptom**: Protobuf validation error about alias pattern.

**Cause**: Alias doesn't start with required `alias/` prefix.

**Solution**: All aliases must follow pattern `alias/[A-Za-z0-9/_-]{1,250}`, e.g., `alias/prod/app-data`.

### Pitfall 3: Short Deletion Window in Production
**Symptom**: Key scheduled for deletion in 7 days, insufficient time to detect and cancel.

**Cause**: Set `deletion_window_days: 7` for "faster cleanup."

**Solution**: Always use 30 days (maximum) for production keys. Only use 7 days for dev/test environments.

### Pitfall 4: Changing Key Type After Creation
**Symptom**: Pulumi wants to replace (delete + recreate) the key.

**Cause**: Changed `key_spec` from symmetric to RSA (or vice versa).

**Solution**: Key type is immutable. Changing it requires creating a new key and migrating data. Plan carefully:
1. Create new key with desired type
2. Re-encrypt data using new key
3. Update applications to use new key
4. Schedule old key for deletion after migration is complete

### Pitfall 5: Forgetting to Create Alias
**Symptom**: Application references key by key ID (UUID), making rotation impossible.

**Cause**: Omitted `alias_name` field.

**Solution**: Always create aliases for keys referenced by applications. Aliases can be updated to point to different keys, enabling zero-downtime key rotation.

## Testing and Debugging

### Debugging with Delve

The `debug.sh` script enables step-through debugging:

1. Uncomment the binary option in `Pulumi.yaml`:
   ```yaml
   runtime:
     options:
       binary: ./debug.sh
   ```
2. Run Pulumi CLI commands normally
3. The debug script launches Delve, allowing breakpoints in any module file

**Use Case**: Debugging enum mapping, investigating AWS API errors, or verifying tag construction.

### Manual Testing with Sample Manifest

Use the sample manifest in `iac/hack/manifest.yaml` for local testing:

```bash
cd iac/pulumi
pulumi stack init dev
pulumi config set aws:region us-east-1
# Set AWS credentials via environment variables or AWS_PROFILE
pulumi up
```

**Verification**:
```bash
pulumi stack output key_id
pulumi stack output key_arn
pulumi stack output alias_name
pulumi stack output rotation_enabled

# Verify in AWS Console or CLI
aws kms describe-key --key-id alias/demo/app-data
aws kms get-key-rotation-status --key-id alias/demo/app-data
```

## Performance Considerations

### Resource Creation Time

KMS key creation is very fast:
- **Key creation**: 1-2 seconds
- **Alias creation**: <1 second

**Total**: Most keys provision in 2-5 seconds.

### Pulumi State Tracking

The module creates two Pulumi resources:
1. One `kms.Key` resource
2. Zero or one `kms.Alias` resource (conditional)

**State Size Impact**: Minimal—KMS keys are lightweight resources.

### Cost Impact

- **Monthly Cost**: $1 per customer-managed key
- **API Costs**: $0.03 per 10,000 requests (encrypt, decrypt, generate data key)
- **Free Tier**: 20,000 requests/month free

**Cost Optimization**: Reuse keys across multiple resources (one key per environment, not per resource).

## Security Recommendations

### Immediate Actions After Deployment

1. **Configure Key Policy** (not managed by this module):
   - Add least-privilege policies for key administrators vs users
   - Require encryption context for decrypt operations
   - Enable CloudTrail logging for all KMS API calls

2. **Set Up Monitoring**:
   ```bash
   # CloudWatch alarm for unauthorized access attempts
   aws cloudwatch put-metric-alarm \
     --alarm-name kms-unauthorized-access \
     --metric-name UnauthorizedApiCalls \
     --namespace AWS/KMS
   ```

3. **Document Key Purpose**:
   - What data does this key protect?
   - Which applications/services use it?
   - Who has admin access?

### Ongoing Security Practices

- **Monitor CloudTrail**: Review KMS API calls weekly for anomalies
- **Rotate Asymmetric Keys Manually**: Since automatic rotation isn't supported
- **Test Disaster Recovery**: Verify encrypted data can be recovered
- **Audit Key Policies**: Quarterly review of who has access
- **Key Inventory**: Maintain documentation of all keys and their purposes
- **Deletion Review**: Before scheduling deletion, verify no resources still use the key

## Conclusion

This Pulumi module demonstrates OpenMCF's philosophy: **secure by default, explicit where necessary**.

The architecture is deliberately focused:
- Enables automatic rotation (unless explicitly disabled)
- Supports all AWS key types with type-safe enum mapping
- Creates optional aliases for human-readable key references
- Enforces maximum deletion window via validation

For teams familiar with Terraform or CloudFormation, the mapping should feel natural. For teams new to KMS, the research documentation (`docs/README.md`) provides critical context about when to use customer-managed keys versus AWS-managed keys.

The result is a production-ready module that provisions secure encryption keys with proper rotation, auditing, and deletion protection—while keeping the API surface small and the defaults secure.

