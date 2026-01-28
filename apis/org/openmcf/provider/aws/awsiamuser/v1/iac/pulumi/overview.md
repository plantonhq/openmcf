# Pulumi Module Architecture: AWS IAM User

## Overview

This Pulumi module provisions AWS IAM users with programmatic access keys through a declarative, protobuf-defined specification. IAM users are long-lived service accounts designed for CI/CD pipelines, third-party integrations, and legacy applications requiring AWS API access with permanent credentials.

The module follows OpenMCF's standard pattern: **input transformation → resource provisioning → secret encryption → output extraction**, with special emphasis on secure handling of access key secrets.

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
    ├── user.go          # IAM user resource implementation
    └── outputs.go       # Output key constants
```

### File Responsibilities

#### `main.go` (entrypoint)
- Unmarshals `AwsIamUserStackInput` from Pulumi stack configuration
- Delegates to `module.Resources()` for actual provisioning
- Minimal logic—purely a thin entrypoint wrapper

#### `module/main.go` (orchestrator)
**Key Function:** `Resources(ctx, stackInput)`

Responsibilities:
1. **Locals Initialization**: Transform stack input into typed locals struct
2. **Provider Configuration**: Handle two provider scenarios:
   - **Default**: Create provider with ambient AWS credentials (IAM role, environment variables)
   - **Explicit**: Create provider with credentials from `stackInput.ProviderConfig` (access key, secret key, session token)
3. **Resource Creation**: Invoke `iamUser()` with initialized locals and provider
4. **Output Export**: Export six IAM user outputs including encrypted secret access key

**Design Decision**: Provider handling supports both CI/CD environments (IRSA, instance profiles) and local development (explicit credentials).

#### `module/locals.go` (input transformer)
**Key Function:** `initializeLocals(ctx, stackInput)`

Transforms the protobuf `AwsIamUserStackInput` into a strongly-typed `Locals` struct and constructs AWS resource tags.

**Tag Construction**: Every IAM user receives five mandatory tags:
- `resource=true`: Marks this as a OpenMCF managed resource
- `organization=<org>`: Organization ID from metadata
- `environment=<env>`: Environment (dev/staging/prod) from metadata
- `resource-kind=AwsIamUser`: CloudResourceKind enum string
- `resource-id=<id>`: Unique resource identifier from metadata

**Why Tags Matter**: These tags enable:
- Cost allocation reporting by org/env/resource-kind
- Policy enforcement (e.g., "only users tagged env=prod can access production resources")
- Resource discovery and inventory management
- Drift detection (manually created users lack these tags)

#### `module/user.go` (resource implementation)
**Key Function:** `iamUser(ctx, locals, provider)`

This is the core implementation file containing all IAM user provisioning logic. It translates the protobuf `AwsIamUserSpec` into Pulumi's AWS IAM resources.

**Implementation Pattern**: The module follows a **four-phase provisioning** approach:

1. **Phase 1: Core User Creation**
   - Create the IAM user with username and tags
   - The user is the anchor resource for all subsequent attachments

2. **Phase 2: Managed Policy Attachments**
   - Iterate through `managed_policy_arns` list
   - Create `UserPolicyAttachment` resources for each ARN
   - Support both AWS-managed policies (e.g., `ReadOnlyAccess`) and customer-managed policies

3. **Phase 3: Inline Policy Creation**
   - Iterate through `inline_policies` map (policy name → policy JSON)
   - Convert each policy from protobuf Struct to JSON string
   - Create `UserPolicy` resources embedded directly in the user

4. **Phase 4: Access Key Creation (Conditional)**
   - Check `disable_access_keys` flag (default: false)
   - If false, create an `AccessKey` resource
   - Mark secret access key as **sensitive** (automatic Pulumi encryption)
   - Base64-encode secret for uniform format and safe transmission

**Struct to JSON Conversion**: The `structToJSONString()` helper safely converts protobuf Struct types to JSON strings, handling nil values gracefully.

**Secret Handling**: Pulumi's built-in secret management automatically encrypts the secret access key in state files and logs. The secret is base64-encoded for additional safety during transmission.

#### `module/outputs.go` (constants)
Defines string constants for output keys, preventing typos and providing a single source of truth.

**Exported Outputs**:
- `user_arn`: Amazon Resource Name (ARN) of the created IAM user
- `user_name`: Name of the IAM user in AWS
- `user_id`: Stable unique ID of the IAM user
- `access_key_id`: Access key ID for programmatic access (if keys were created)
- `secret_access_key`: Base64-encoded secret key (if keys were created, **sensitive**)
- `console_url`: AWS console sign-in URL

## Data Flow Diagram

```
┌──────────────────────────────────┐
│ AwsIamUserStackInput (protobuf)  │
│  ├─ target: AwsIamUser           │
│  │   ├─ metadata                 │
│  │   └─ spec: AwsIamUserSpec     │
│  │       ├─ user_name            │
│  │       ├─ managed_policy_arns  │
│  │       ├─ inline_policies      │
│  │       └─ disable_access_keys  │
│  └─ provider_config (optional)   │
└──────────────┬───────────────────┘
               │
               ▼
      ┌────────────────┐
      │initializeLocals│
      └────────┬────────┘
               │ Creates:
               │ - Locals.AwsIamUser
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
       │   iamUser()     │
       └────────┬────────┘
                │
                ├─ Create IAM User:
                │    iam.NewUser(ctx, userName, &iam.UserArgs{
                │      Name: userName,
                │      Tags: awsTags
                │    })
                │
                ├─ For each managed_policy_arn:
                │    iam.NewUserPolicyAttachment(ctx, attachName, &iam.UserPolicyAttachmentArgs{
                │      User:      user.Name,
                │      PolicyArn: policyArn
                │    })
                │
                ├─ For each inline_policy (name → JSON):
                │    Convert policy Struct → JSON string
                │    iam.NewUserPolicy(ctx, inlineName, &iam.UserPolicyArgs{
                │      User:   user.Name,
                │      Policy: policyJSON
                │    })
                │
                ├─ If !disable_access_keys:
                │    iam.NewAccessKey(ctx, akName, &iam.AccessKeyArgs{
                │      User: user.Name
                │    })
                │    Base64-encode secret for safe output
                │
                └─ Export outputs:
                     ctx.Export("user_arn", user.Arn)
                     ctx.Export("user_name", user.Name)
                     ctx.Export("user_id", user.UniqueId)
                     ctx.Export("console_url", consoleUrl)
                     ctx.Export("access_key_id", accessKey.ID()) [if created]
                     ctx.Export("secret_access_key", base64(accessKey.Secret)) [if created, sensitive]
                │
                ▼
       ┌─────────────────────────┐
       │  AWS IAM Service        │
       │   ├─ IAM User           │
       │   ├─ Managed Policies   │
       │   ├─ Inline Policies    │
       │   └─ Access Key (opt)   │
       └─────────────────────────┘
```

## Resource Relationships

```
AwsIamUserSpec
  │
  ├─ user_name ──────────────────┐
  │                               │
  ├─ managed_policy_arns ────────┼─ IAM User
  │    (repeated string)          │   ├─ Name (from spec.user_name)
  │                                │   ├─ Tags (org/env/resource metadata)
  ├─ inline_policies ────────────┼───┤   └─ UniqueId (AWS-generated)
  │    (map<string, Struct>)      │   │
  │                                │   │
  └─ disable_access_keys ────────┼───┤
       (bool, default: false)     │   │
                                  │   │
                                  │   ├─ Managed Policy Attachments
                                  │   │   (UserPolicyAttachment resources)
                                  │   │   ├─ AWS-managed policies
                                  │   │   └─ Customer-managed policies
                                  │   │
                                  │   ├─ Inline Policies
                                  │   │   (UserPolicy resources)
                                  │   │   ├─ Policy name
                                  │   │   └─ Policy JSON document
                                  │   │
                                  │   └─ Access Key (conditional)
                                  │       (AccessKey resource, if !disable_access_keys)
                                  │       ├─ Access Key ID
                                  │       └─ Secret Access Key (encrypted)
                                  │
                                  └─ Console URL (computed)
```

### Critical Relationships

**User → Policy Attachments**:
- Managed policy attachments reference policies by ARN
- Policies can be AWS-managed (e.g., `arn:aws:iam::aws:policy/ReadOnlyAccess`) or customer-managed
- Multiple users can share the same managed policy
- Changes to a managed policy affect all users using it

**User → Inline Policies**:
- Inline policies are embedded directly in the user
- Deleted when the user is deleted (tight coupling by design)
- Cannot be shared with other users
- Used for user-specific permissions that shouldn't apply elsewhere

**User → Access Key (1:1 Relationship)**:
- One access key per user (this module creates only one)
- AWS allows up to 2 access keys per user (for rotation), but this module creates 1 initially
- Key rotation requires creating a new key, updating applications, then deleting the old key
- Secret key is only available at creation time—after that, it's encrypted and unrecoverable

**Access Key → Secret Handling**:
- Secret access key is automatically marked as **sensitive** by Pulumi
- Encrypted in Pulumi state files with your chosen encryption provider (passphrase, KMS, etc.)
- Base64-encoded for safe transmission (prevents special character issues)
- Should be immediately stored in a secret manager (AWS Secrets Manager, HashiCorp Vault)

## Key Design Decisions

### 1. Single Access Key by Default

**Decision**: The module creates exactly one access key if `disable_access_keys` is false.

**Rationale**:
- **Simplicity**: Most use cases need one key initially
- **Rotation Pattern**: Best practice is to create a second key during rotation, not at initial provisioning
- **Security**: Fewer keys = smaller attack surface
- **Explicit Intent**: If you want multiple keys, create additional ones manually or via separate stacks

**Alternative Approach**: Some tools create zero keys by default, requiring explicit key creation. This module defaults to creating keys because the primary use case (service accounts) always needs programmatic access.

### 2. Base64-Encoded Secret Output

**Decision**: The secret access key is base64-encoded before exporting.

**Rationale**:
- **Safe Transmission**: Prevents issues with special characters in shell scripts or CI/CD systems
- **Uniform Format**: All secrets are encoded consistently
- **Decoding Clarity**: Explicit decode step reminds users to handle secrets carefully
- **No Security Reduction**: Base64 is encoding, not encryption (Pulumi encryption is still applied)

**Implementation**:
```go
secretAccessKey = accessKey.Secret.ApplyT(func(s string) *string {
    enc := base64.StdEncoding.EncodeToString([]byte(s))
    return &enc
}).(pulumi.StringPtrOutput)
```

**Usage**: Users must base64-decode the output to get the actual secret: `echo $SECRET | base64 -d`

### 3. Conditional Access Key Creation

**Decision**: Check `disable_access_keys` flag; if true, skip key creation entirely.

**Rationale**:
- **Identity-Only Users**: Some scenarios need IAM user identity without API keys (e.g., federated access)
- **Console-Only Users**: Legacy console users shouldn't have programmatic keys
- **Explicit Opt-Out**: Default behavior (create keys) matches most use cases; explicit flag required to disable

**Implementation**:
```go
var accessKey *iam.AccessKey
if !spec.DisableAccessKeys {
    accessKey, err = iam.NewAccessKey(ctx, akName, &iam.AccessKeyArgs{
        User: usr.Name,
    }, pulumi.Provider(provider))
}
```

**Output Behavior**: If no keys are created, `access_key_id` and `secret_access_key` outputs are nil.

### 4. Separate Attachment Resources

**Decision**: Managed policies are attached via separate `UserPolicyAttachment` resources rather than inline in the user definition.

**Rationale**:
- **AWS API Semantics**: IAM treats attachments as separate resources
- **Update Safety**: Changing attached policies doesn't require user replacement
- **Pulumi State Tracking**: Each attachment is independently tracked for drift detection
- **Error Isolation**: If one attachment fails, others can still succeed

**Naming Convention**: Attachments are named `{userName}-attach-{index}` for uniqueness and stability.

### 5. Console URL Generation

**Decision**: Generate a static AWS console sign-in URL as an output.

**Rationale**:
- **Convenience**: Users can quickly access the console without remembering the URL
- **Consistency**: All IAM user deployments have this reference
- **Minimal Value**: The URL is generic (`https://signin.aws.amazon.com/console`), not user-specific

**Note**: To create a console login profile (username + password), you'd need additional configuration not included in this module (console access is discouraged for service accounts).

## IAM-Specific Implementation Details

### Managed Policy Attachment Loop

The module creates one attachment resource per managed policy ARN.

**Implementation**:
```go
for idx, policyArn := range spec.ManagedPolicyArns {
    attachName := fmt.Sprintf("%s-attach-%d", userName, idx)
    _, err := iam.NewUserPolicyAttachment(ctx, attachName, &iam.UserPolicyAttachmentArgs{
        User:      usr.Name,
        PolicyArn: pulumi.String(policyArn),
    }, pulumi.Provider(provider))
    // error handling...
}
```

**Naming Convention**: Attachments are named `{userName}-attach-{index}` to ensure uniqueness and maintain stable resource names across updates.

**Error Handling**: If an attachment fails (e.g., policy ARN doesn't exist), the error propagates immediately—no silent skipping.

### Inline Policy Conversion and Creation

Inline policies require both name and JSON document conversion.

**Implementation**:
```go
for policyName, inlineStruct := range spec.InlinePolicies {
    inlinePolicyString, err := structToJSONString(inlineStruct)
    if err != nil {
        return nil, errors.Wrapf(err, "failed to marshal inline policy for %s", policyName)
    }

    inlineName := fmt.Sprintf("%s-inline-%s", userName, policyName)
    _, err = iam.NewUserPolicy(ctx, inlineName, &iam.UserPolicyArgs{
        User:   usr.Name,
        Policy: pulumi.String(inlinePolicyString),
    }, pulumi.Provider(provider))
    // error handling...
}
```

**Naming Convention**: Inline policies are named `{userName}-inline-{policyName}` to prevent collisions and maintain clarity.

**Lifecycle**: When the user is deleted, all inline policies are automatically deleted (tight coupling by design).

### Access Key Secret Handling

The access key secret is the most sensitive output and receives special treatment.

**Security Measures**:
1. **Pulumi Secret Marking**: The `AccessKey.Secret` field is automatically marked as sensitive by the Pulumi AWS provider
2. **State Encryption**: Pulumi encrypts sensitive values in state files using your configured encryption provider
3. **Log Masking**: Sensitive values are masked in Pulumi logs and CLI output
4. **Base64 Encoding**: Additional layer for safe transmission

**Implementation**:
```go
secretAccessKey = accessKey.Secret.ApplyT(func(s string) *string {
    enc := base64.StdEncoding.EncodeToString([]byte(s))
    return &enc
}).(pulumi.StringPtrOutput)
```

**Best Practice Flow**:
1. Deploy IAM user via Pulumi
2. Immediately retrieve secret with `pulumi stack output secret_access_key --show-secrets`
3. Base64-decode the secret
4. Store in AWS Secrets Manager or vault
5. Delete the secret from local terminal history
6. Configure application/CI/CD to retrieve from secret manager

### Output Exports

Six outputs are exported, with conditional presence based on access key creation.

**Implementation**:
```go
ctx.Export(OpUserArn, results.UserArn)
ctx.Export(OpUserName, results.UserName)
ctx.Export(OpUserId, results.UserId)
ctx.Export(OpConsoleUrl, results.ConsoleUrl)
ctx.Export(OpAccessKeyId, results.AccessKeyId)       // nil if disable_access_keys=true
ctx.Export(OpSecretAccessKey, results.SecretAccessKey) // nil if disable_access_keys=true
```

**Output Usage**:
- **user_arn**: Reference in IAM policies, CloudFormation, or other infrastructure
- **user_name**: Human-readable identifier for documentation or scripts
- **user_id**: Stable identifier that doesn't change even if username changes
- **access_key_id**: Configure in applications or CI/CD (not sensitive, safe to expose)
- **secret_access_key**: Retrieve once, decode, store securely, never log
- **console_url**: Convenience link for console access (generic URL)

## Error Handling Philosophy

The module follows a **fail-fast** approach:

1. **Validation at Protobuf Level**: The `spec.proto` validation rules catch configuration errors before Pulumi runs
2. **AWS API Errors Propagate**: If AWS rejects a configuration, the error propagates immediately (no silent fallbacks)
3. **Wrapped Errors**: All errors include context (`errors.Wrap()`) for easier debugging

**Rationale**: Infrastructure as code demands predictability. Silent defaults or error recovery can mask misconfiguration and create surprise behavior.

## Common Pitfalls and Gotchas

### Pitfall 1: Policy ARN Doesn't Exist
**Symptom**: AWS API error: `NoSuchEntity: The policy with the specified ARN does not exist`

**Cause**: Managed policy ARN is incorrect or the policy hasn't been created yet.

**Solution**: Verify policy ARN format and ensure customer-managed policies are created before referencing them. AWS-managed policies should use the format `arn:aws:iam::aws:policy/{policyName}`.

### Pitfall 2: Invalid Username Format
**Symptom**: Protobuf validation error or AWS API error about invalid username.

**Cause**: Username contains invalid characters or exceeds 64 characters.

**Solution**: Follow AWS naming rules: 1-64 characters, alphanumeric plus `+=,.@_-` only. The protobuf validation enforces this pattern.

### Pitfall 3: Losing Secret Access Key
**Symptom**: Unable to retrieve secret key after initial deployment.

**Cause**: Forgot to save the secret immediately; AWS doesn't allow retrieving it again.

**Solution**: Always run `pulumi stack output secret_access_key --show-secrets` immediately after deployment and store securely. If lost, delete the access key and create a new one.

### Pitfall 4: Inline Policy Size Limit
**Symptom**: AWS API error about policy size exceeding limit.

**Cause**: Inline policy exceeds 2,048 characters (per policy) or 10,240 total (all inline policies combined).

**Solution**: Move large policies to customer-managed policies and reference them via `managed_policy_arns`. Managed policies support up to 6,144 characters.

### Pitfall 5: Access Key in Version Control
**Symptom**: Security scanner alerts about exposed credentials in Git history.

**Cause**: Secret access key was printed to terminal or config file that got committed.

**Solution**: Never commit secrets to Git. Use `.gitignore` for any files containing secrets. If accidentally committed, rotate the key immediately (create new key, update applications, delete old key).

### Pitfall 6: Base64 Confusion
**Symptom**: Applications fail to authenticate with "Invalid credentials" error.

**Cause**: Using the base64-encoded secret directly without decoding.

**Solution**: Remember to base64-decode the secret before configuring applications: `echo $SECRET | base64 -d`

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

**Use Case**: Debugging policy JSON conversion, attachment failures, or AWS API errors.

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
pulumi stack output user_arn
pulumi stack output user_name
pulumi stack output access_key_id
pulumi stack output secret_access_key --show-secrets | base64 -d

# Verify in AWS Console or CLI
aws iam get-user --user-name <user_name>
aws iam list-attached-user-policies --user-name <user_name>
aws iam list-access-keys --user-name <user_name>
```

## Performance Considerations

### Resource Creation Time

IAM user creation is typically fast:
- **User creation**: 1-3 seconds
- **Managed policy attachments**: <1 second each
- **Inline policy creation**: <1 second each
- **Access key creation**: <1 second

**Total**: Most users provision in 5-10 seconds.

### Pulumi State Tracking

The module creates multiple Pulumi resources:
1. One `iam.User` resource
2. N `iam.UserPolicyAttachment` resources (where N = number of managed policies)
3. M `iam.UserPolicy` resources (where M = number of inline policies)
4. 0 or 1 `iam.AccessKey` resource (based on `disable_access_keys`)

**State Size Impact**: Minimal—IAM users are lightweight resources.

### Update Detection

Pulumi automatically detects drift in:
- User tags changes
- Managed policy attachments (added/removed externally)
- Inline policy modifications
- Access key existence (but not the secret value—it's write-only)

**Recommendation**: Run `pulumi refresh` periodically to detect out-of-band changes made via AWS console or other tools.

## Security Recommendations

### Immediate Actions After Deployment

1. **Retrieve and store secrets immediately**:
   ```bash
   SECRET=$(pulumi stack output secret_access_key --show-secrets)
   echo $SECRET | base64 -d | aws secretsmanager create-secret --name /app/user-creds --secret-string -
   ```

2. **Configure IAM password policy** (separate from this module):
   - Minimum password length
   - Require symbols, numbers, uppercase
   - Password expiration
   - Prevent password reuse

3. **Enable CloudTrail** (if not already enabled):
   - Track all IAM API calls
   - Monitor for unusual access patterns
   - Set up alerts for key usage from unexpected IPs

4. **Set up key rotation reminders**:
   - Calendar reminder every 90 days
   - Automated key age monitoring via CloudWatch
   - Scripted rotation process

### Ongoing Security Practices

- **Monitor with Access Advisor**: Check which permissions are actually used; remove unused policies
- **Review CloudTrail logs**: Identify anomalous API calls or access patterns
- **Tag consistently**: Enables IAM policy enforcement (e.g., "users tagged team=devops can only assume devops roles")
- **Least privilege**: Start with minimal permissions; add only when proven necessary
- **Avoid wildcard permissions**: Never use `"Action": "*"` or `"Resource": "*"` in policies
- **Consider alternatives**: Always ask if you can use IAM roles or federation instead

## Conclusion

This Pulumi module demonstrates OpenMCF's philosophy: **secure by default, explicit where necessary**.

The architecture is deliberately focused:
- Creates programmatic service accounts (the 80% use case)
- Handles access key secrets securely (automatic encryption)
- Provides clear, typed outputs
- Fails fast on misconfiguration

For teams familiar with Terraform or CloudFormation, the mapping should feel natural. For teams new to IAM users, the research documentation (`docs/README.md`) provides critical context about when to use IAM users versus roles or federation.

The result is a production-ready module that provisions secure service accounts with proper permission scoping, encrypted credential handling, and comprehensive audit trails—minimizing the security risks inherent in long-lived credentials.

