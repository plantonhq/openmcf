# Pulumi Module Architecture: AWS IAM Role

## Overview

This Pulumi module provisions AWS IAM roles through a declarative, protobuf-defined specification. IAM roles are the foundation of AWS security, enabling secure delegation of permissions to services, applications, and users without long-lived credentials.

The module follows OpenMCF's standard pattern: **input transformation → resource provisioning → output extraction**, with special attention to IAM's two-policy model (trust policy + permissions policies).

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
    ├── iam_role.go      # IAM role resource implementation
    └── outputs.go       # Output key constants
```

### File Responsibilities

#### `main.go` (entrypoint)
- Unmarshals `AwsIamRoleStackInput` from Pulumi stack configuration
- Delegates to `module.Resources()` for actual provisioning
- Minimal logic—purely a thin entrypoint wrapper

#### `module/main.go` (orchestrator)
**Key Function:** `Resources(ctx, stackInput)`

Responsibilities:
1. **Locals Initialization**: Transform stack input into typed locals struct
2. **Provider Configuration**: Handle two provider scenarios:
   - **Default**: Create provider with ambient AWS credentials (IAM role, environment variables)
   - **Explicit**: Create provider with credentials from `stackInput.ProviderConfig` (access key, secret key, session token)
3. **Resource Creation**: Invoke `iamRole()` with initialized locals and provider
4. **Output Export**: Export IAM role ARN and name using Pulumi's export mechanism

**Design Decision**: Provider handling supports both CI/CD environments (IRSA, instance profiles) and local development (explicit credentials).

#### `module/locals.go` (input transformer)
**Key Function:** `initializeLocals(ctx, stackInput)`

Transforms the protobuf `AwsIamRoleStackInput` into a strongly-typed `Locals` struct and constructs AWS resource tags.

**Tag Construction**: Every IAM role receives five mandatory tags:
- `resource=true`: Marks this as a OpenMCF managed resource
- `organization=<org>`: Organization ID from metadata
- `environment=<env>`: Environment (dev/staging/prod) from metadata
- `resource-kind=AwsIamRole`: CloudResourceKind enum string
- `resource-id=<id>`: Unique resource identifier from metadata

**Why Tags Matter**: These tags enable:
- Cost allocation reporting by org/env/resource-kind
- Policy enforcement (e.g., "only roles with env=prod can assume sensitive policies")
- Resource discovery and inventory management
- Drift detection (manually created roles lack these tags)

#### `module/iam_role.go` (resource implementation)
**Key Function:** `iamRole(ctx, locals, provider)`

This is the core implementation file containing all IAM role provisioning logic. It translates the protobuf `AwsIamRoleSpec` into Pulumi's AWS IAM resources.

**Implementation Pattern**: The module follows a **three-phase provisioning** approach:

1. **Phase 1: Core Role Creation**
   - Convert trust policy from protobuf Struct to JSON string
   - Create the IAM role with name, assume role policy, description, path, and tags
   - Validate trust policy structure

2. **Phase 2: Managed Policy Attachments**
   - Iterate through `managed_policy_arns` list
   - Create `RolePolicyAttachment` resources for each ARN
   - Support both AWS-managed policies (e.g., `arn:aws:iam::aws:policy/...`) and customer-managed policies

3. **Phase 3: Inline Policy Creation**
   - Iterate through `inline_policies` map (policy name → policy JSON)
   - Convert each policy from protobuf Struct to JSON string
   - Create `RolePolicy` resources embedded directly in the role

**Struct to JSON Conversion**: The `structToJSONString()` helper safely converts protobuf Struct types to JSON strings, handling nil values gracefully.

#### `module/outputs.go` (constants)
Defines string constants for output keys, preventing typos and providing a single source of truth.

**Exported Outputs**:
- `role_arn`: Amazon Resource Name (ARN) of the created IAM role
- `role_name`: Name of the IAM role in AWS

## Data Flow Diagram

```
┌──────────────────────────────────┐
│ AwsIamRoleStackInput (protobuf)  │
│  ├─ target: AwsIamRole           │
│  │   ├─ metadata                 │
│  │   └─ spec: AwsIamRoleSpec     │
│  │       ├─ description          │
│  │       ├─ path                 │
│  │       ├─ trust_policy         │
│  │       ├─ managed_policy_arns  │
│  │       └─ inline_policies      │
│  └─ provider_config (optional)   │
└──────────────┬───────────────────┘
               │
               ▼
      ┌────────────────┐
      │initializeLocals│
      └────────┬────────┘
               │ Creates:
               │ - Locals.AwsIamRole
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
       │   iamRole()     │
       └────────┬────────┘
                │
                ├─ Convert trust_policy Struct → JSON string
                │
                ├─ Create IAM Role:
                │    iam.NewRole(ctx, roleName, &iam.RoleArgs{
                │      Name:             roleName,
                │      AssumeRolePolicy: trustPolicyJSON,
                │      Description:      description,
                │      Path:             path,
                │      Tags:             awsTags
                │    })
                │
                ├─ For each managed_policy_arn:
                │    iam.NewRolePolicyAttachment(ctx, attachName, &iam.RolePolicyAttachmentArgs{
                │      Role:      role.Name,
                │      PolicyArn: policyArn
                │    })
                │
                ├─ For each inline_policy (name → JSON):
                │    Convert policy Struct → JSON string
                │    iam.NewRolePolicy(ctx, inlineName, &iam.RolePolicyArgs{
                │      Role:   role.Name,
                │      Policy: policyJSON
                │    })
                │
                └─ Export outputs:
                     ctx.Export("role_arn", role.Arn)
                     ctx.Export("role_name", role.Name)
                │
                ▼
       ┌─────────────────────────┐
       │  AWS IAM Service        │
       │   ├─ IAM Role           │
       │   ├─ Trust Policy       │
       │   ├─ Managed Policies   │
       │   └─ Inline Policies    │
       └─────────────────────────┘
```

## Resource Relationships

```
AwsIamRoleSpec
  │
  ├─ trust_policy ────────────────┐
  │    (google.protobuf.Struct)   │
  │    → Converted to JSON        │
  │                                │
  ├─ managed_policy_arns ────────┼─ IAM Role
  │    (repeated string)          │   ├─ Name (from metadata)
  │                                │   ├─ AssumeRolePolicy (trust policy JSON)
  ├─ inline_policies ────────────┼───┤   ├─ Description
  │    (map<string, Struct>)      │   │   ├─ Path (default: "/")
  │                                │   │   └─ Tags (org/env/resource metadata)
  ├─ description ────────────────┘   │
  │                                   │
  └─ path ───────────────────────────┤
                                      │
                                      ├─ Managed Policy Attachments
                                      │   (RolePolicyAttachment resources)
                                      │   ├─ AWS-managed policies
                                      │   └─ Customer-managed policies
                                      │
                                      └─ Inline Policies
                                          (RolePolicy resources)
                                          ├─ Policy name
                                          └─ Policy JSON document
```

### Critical Relationships

**Trust Policy → AssumeRole Permission**:
- The trust policy controls **who** can assume the role
- Must be valid JSON with Version, Statement, Effect, Principal, Action
- Separate from permissions policies—validates identity before granting access

**Managed Policies → Reusability**:
- Attached by ARN reference
- Can be AWS-managed (e.g., `AWSLambdaBasicExecutionRole`) or customer-managed
- Changes to the policy affect all roles using it
- Support version rollback (IAM maintains last 5 versions)

**Inline Policies → Tight Coupling**:
- Embedded directly in the role definition
- Deleted when the role is deleted
- Cannot be shared with other roles
- Used for role-specific permissions that shouldn't apply elsewhere

**Path → Organizational Hierarchy**:
- Defaults to "/" (root)
- Can use paths like "/service-roles/" or "/application/frontend/" for grouping
- Affects role ARN: `arn:aws:iam::123456789012:role/service-roles/my-role`

## Key Design Decisions

### 1. Protobuf Struct for Policy Documents

**Decision**: Trust policies and inline policies are defined as `google.protobuf.Struct` rather than strings.

**Rationale**:
- **Type Safety**: Policy documents are structured JSON, not arbitrary strings
- **Validation**: Can validate policy structure at protobuf level (e.g., required fields)
- **Flexibility**: Supports nested policy structures without escaping issues
- **Consistency**: Matches AWS API expectations (JSON documents)

**Trade-off**: Requires conversion from Struct to JSON string at runtime, but this is a small overhead compared to API call latency.

### 2. Separate Attachment Resources

**Decision**: Managed policies are attached via separate `RolePolicyAttachment` resources rather than inline in the role definition.

**Rationale**:
- **AWS API Semantics**: IAM treats attachments as separate resources
- **Update Safety**: Changing attached policies doesn't require role replacement
- **Pulumi State Tracking**: Each attachment is independently tracked for drift detection
- **Error Isolation**: If one attachment fails, others can still succeed

**Implementation**: Each policy ARN gets a unique attachment resource named `{roleName}-attach-{index}`.

### 3. Inline Policy Name as Map Key

**Decision**: Inline policies use a map with policy name as key, not a repeated field with name+document.

**Rationale**:
- **Uniqueness Guarantee**: Map keys are unique by definition
- **Simpler API**: No need for separate name field within each policy
- **Clearer Intent**: Policy name is naturally the identifier

**Trade-off**: Protobuf maps don't preserve order, but policy order doesn't matter for IAM.

### 4. Conservative Path Default

**Decision**: Path defaults to "/" if omitted.

**Rationale**:
- **AWS Default Alignment**: AWS console and CLI default to "/"
- **Simplicity**: Most users don't need organizational paths
- **Override Available**: Advanced users can specify custom paths for grouping

### 5. Tag Injection at Creation Time

**Decision**: Tags are constructed from metadata and attached at role creation (not separate API call).

**Rationale**:
- **Atomicity**: Tags exist from the moment the role is created
- **Consistency**: Every OpenMCF resource has identical tagging structure
- **Cost Allocation**: AWS billing sees tags immediately

## IAM-Specific Implementation Details

### Trust Policy Conversion

The trust policy must be converted from protobuf Struct to a JSON string that AWS IAM accepts.

**Implementation**:
```go
trustPolicyString, err := structToJSONString(spec.TrustPolicy)
if err != nil {
    return errors.Wrap(err, "failed to marshal trust policy JSON")
}
```

**Validation Note**: The protobuf spec requires `trust_policy` to be non-nil. However, the helper handles nil gracefully (returns `"{}"`) to prevent panics if validation is bypassed.

**Common Trust Policy Patterns**:
- **Lambda**: `{"Principal": {"Service": "lambda.amazonaws.com"}}`
- **ECS Tasks**: `{"Principal": {"Service": "ecs-tasks.amazonaws.com"}}`
- **Cross-Account**: `{"Principal": {"AWS": "arn:aws:iam::123456789012:root"}}`
- **Federated Users**: `{"Principal": {"Federated": "arn:aws:iam::123456789012:saml-provider/..."}}`

### Managed Policy Attachment Loop

The module creates one attachment resource per managed policy ARN.

**Implementation**:
```go
for idx, policyArn := range spec.ManagedPolicyArns {
    attachName := fmt.Sprintf("%s-attach-%d", roleName, idx)
    _, err := iam.NewRolePolicyAttachment(ctx, attachName, &iam.RolePolicyAttachmentArgs{
        Role:      iamRole.Name,
        PolicyArn: pulumi.String(policyArn),
    }, pulumi.Provider(provider))
    // error handling...
}
```

**Naming Convention**: Attachments are named `{roleName}-attach-{index}` to ensure uniqueness and maintain stable resource names across updates.

**Error Handling**: If an attachment fails (e.g., policy ARN doesn't exist), the error propagates immediately—no silent skipping.

### Inline Policy Conversion and Creation

Inline policies require both name and JSON document conversion.

**Implementation**:
```go
for policyName, inlineStruct := range spec.InlinePolicies {
    inlinePolicyString, err := structToJSONString(inlineStruct)
    if err != nil {
        return errors.Wrapf(err, "failed to marshal inline policy for %s", policyName)
    }

    inlineName := fmt.Sprintf("%s-inline-%s", roleName, policyName)
    _, err = iam.NewRolePolicy(ctx, inlineName, &iam.RolePolicyArgs{
        Role:   iamRole.Name,
        Policy: pulumi.String(inlinePolicyString),
    }, pulumi.Provider(provider))
    // error handling...
}
```

**Naming Convention**: Inline policies are named `{roleName}-inline-{policyName}` to prevent collisions and maintain clarity.

**Lifecycle**: When the role is deleted, all inline policies are automatically deleted (tight coupling by design).

### Output Exports

Only two outputs are exported:

**Implementation**:
```go
ctx.Export(OpRoleArn, iamRole.Arn)
ctx.Export(OpRoleName, iamRole.Name)
```

**Output Usage**:
- **role_arn**: Used in Lambda function definitions, ECS task definitions, EC2 instance profiles, and other services requiring an execution role
- **role_name**: Used for IAM policy references, CloudFormation stack parameters, and human-readable identification

## Error Handling Philosophy

The module follows a **fail-fast** approach:

1. **Validation at Protobuf Level**: The `spec.proto` validation rules catch configuration errors before Pulumi runs
2. **AWS API Errors Propagate**: If AWS rejects a configuration, the error propagates immediately (no silent fallbacks)
3. **Wrapped Errors**: All errors include context (`errors.Wrap()`) for easier debugging

**Rationale**: Infrastructure as code demands predictability. Silent defaults or error recovery can mask misconfiguration and create surprise behavior.

## Common Pitfalls and Gotchas

### Pitfall 1: Invalid Trust Policy JSON
**Symptom**: AWS API error: `MalformedPolicyDocument: The policy contains the following error: Invalid JSON`

**Cause**: Trust policy Struct doesn't convert to valid JSON (missing required fields, syntax errors).

**Solution**: Validate trust policy structure before deployment. Ensure it has `Version`, `Statement` with `Effect`, `Principal`, and `Action`.

### Pitfall 2: Policy ARN Doesn't Exist
**Symptom**: AWS API error: `NoSuchEntity: The policy with the specified ARN does not exist`

**Cause**: Managed policy ARN is incorrect or the policy hasn't been created yet.

**Solution**: Verify policy ARN format and ensure customer-managed policies are created before referencing them. AWS-managed policies should use the format `arn:aws:iam::aws:policy/{policyName}`.

### Pitfall 3: Circular Dependency with Policy
**Symptom**: Pulumi error about resource dependency cycle.

**Cause**: Trying to create a role and a policy that references that role in the same stack, with circular references.

**Solution**: Use managed policy ARNs for existing policies, or structure the stack to create the policy first, then the role.

### Pitfall 4: Overly Permissive Trust Policy
**Symptom**: Security audit flags role as high-risk.

**Cause**: Trust policy uses wildcard principals (`"Principal": "*"`) or missing conditions.

**Solution**: Always specify explicit principals (service names, account ARNs). Add conditions like `aws:SourceAccount` or `sts:ExternalId` for additional security.

### Pitfall 5: Inline Policy Size Limit
**Symptom**: AWS API error about policy size exceeding limit.

**Cause**: Inline policy exceeds 2,048 characters (per policy) or 10,240 total (all inline policies combined).

**Solution**: Move large policies to customer-managed policies and reference them via `managed_policy_arns`. Managed policies support up to 6,144 characters.

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

**Use Case**: Debugging trust policy JSON conversion, attachment failures, or AWS API errors.

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
pulumi stack output role_arn
pulumi stack output role_name

# Verify in AWS Console or CLI
aws iam get-role --role-name <role_name>
aws iam list-attached-role-policies --role-name <role_name>
```

## Performance Considerations

### Resource Creation Time

IAM role creation is typically fast:
- **Role creation**: 1-3 seconds
- **Managed policy attachments**: <1 second each
- **Inline policy creation**: <1 second each

**Total**: Most roles provision in 5-10 seconds.

### Pulumi State Tracking

The module creates multiple Pulumi resources:
1. One `iam.Role` resource
2. N `iam.RolePolicyAttachment` resources (where N = number of managed policies)
3. M `iam.RolePolicy` resources (where M = number of inline policies)

**State Size Impact**: Minimal—IAM roles are lightweight resources.

### Update Detection

Pulumi automatically detects drift in:
- Role trust policy changes
- Managed policy attachments (added/removed externally)
- Inline policy modifications
- Tag changes

**Recommendation**: Run `pulumi refresh` periodically to detect out-of-band changes made via AWS console or other tools.

## Conclusion

This Pulumi module demonstrates OpenMCF's philosophy: **expose the full power of IAM roles while providing guardrails through validation**.

The architecture is deliberately explicit:
- No hidden defaults or magic mappings
- Clear separation of concerns (locals, provider, resource creation, outputs)
- Validation at the protobuf layer prevents misconfiguration
- Resource creation follows AWS IAM API semantics closely

For teams familiar with Terraform or CloudFormation, the mapping should feel natural. For teams new to IAM roles, the explicit field names and comprehensive research documentation (`docs/README.md`) provide the context needed to make informed security decisions.

The result is a production-ready module that can provision simple roles (single managed policy, basic trust policy) or complex roles (multiple managed policies, multiple inline policies, custom paths, detailed trust conditions) with equal ease.

