# Pulumi Implementation Architecture Overview

## Purpose

This document provides an architectural overview of the Pulumi implementation for Cloudflare Zero Trust Access Application, explaining design decisions, resource dependencies, and implementation patterns.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│ OpenMCF Manifest (YAML)                            │
│ CloudflareZeroTrustAccessApplication                        │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ Pulumi Program (main.go)                                    │
│ - Loads StackInput                                          │
│ - Initializes Cloudflare Provider                           │
│ - Calls module.Resources()                                  │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ Module (module/main.go)                                     │
│ - Initializes Locals                                        │
│ - Creates Cloudflare Provider with API token                │
│ - Calls application() to create resources                   │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ application() - module/application.go                       │
│                                                             │
│ 1. Create Access Application                                │
│    ├─ Name: spec.application_name                           │
│    ├─ Domain: spec.hostname                                 │
│    ├─ Type: "self_hosted"                                   │
│    └─ Session Duration: spec.session_duration_minutes       │
│                                                             │
│ 2. Lookup Zone (to get account_id)                          │
│    └─ Query by spec.zone_id                                 │
│                                                             │
│ 3. Create Access Policy                                     │
│    ├─ Account ID: from zone lookup                          │
│    ├─ Application ID: from step 1                           │
│    ├─ Decision: ALLOW or DENY                               │
│    ├─ Include Rules:                                        │
│    │  ├─ Email blocks (spec.allowed_emails)                 │
│    │  └─ Group blocks (spec.allowed_google_groups)          │
│    └─ Require Rules:                                        │
│       └─ MFA block (if spec.require_mfa)                    │
│                                                             │
│ 4. Export Outputs                                           │
│    ├─ application_id                                        │
│    ├─ public_hostname                                       │
│    └─ policy_id                                             │
└─────────────────────────────────────────────────────────────┘
```

## Resource Creation Flow

### 1. Cloudflare Access Application

**Purpose**: Creates the protected application resource in Cloudflare Zero Trust.

**Implementation** (`module/application.go` lines 19-41):

```go
createdAccessApplication, err := cloudflare.NewAccessApplication(
    ctx,
    "access_application",
    &cloudflare.AccessApplicationArgs{
        Name:   pulumi.String(spec.ApplicationName),
        ZoneId: pulumi.String(spec.ZoneId),
        Domain: pulumi.String(spec.Hostname),
        Type:   pulumi.StringPtr("self_hosted"),
        SessionDuration: pulumi.StringPtr(fmt.Sprintf("%dm", spec.SessionDurationMinutes)),
    },
    pulumi.Provider(cloudflareProvider),
)
```

**Key Design Decisions**:
- **Type**: Hardcoded to `"self_hosted"` because this is the only relevant type for protecting custom applications. Cloudflare also supports SaaS app types, but those are not in scope for this 80/20 implementation.
- **Session Duration**: Converted from minutes (int32) to Cloudflare's duration string format (e.g., "480m" for 8 hours).
- **Provider**: Uses scoped Cloudflare provider with API token for authentication.

### 2. Zone Lookup

**Purpose**: Retrieve the Cloudflare account ID associated with the DNS zone. This is required to create the Access Policy.

**Implementation** (`module/application.go` lines 44-49):

```go
zone := cloudflare.LookupZoneOutput(ctx, cloudflare.LookupZoneOutputArgs{
    ZoneId: pulumi.String(spec.ZoneId),
}, pulumi.Provider(cloudflareProvider))

accountId := zone.Account().Id()
```

**Why This is Necessary**:
- Cloudflare's Access Policy resource requires `account_id` (not zone_id) to bind the policy
- The zone lookup provides the account ID associated with the zone
- This is a read-only operation; it doesn't create or modify any resources

### 3. Access Policy

**Purpose**: Define who can access the application and under what conditions.

**Implementation** (`module/application.go` lines 51-100):

#### Include Rules (Who Can Access)

**Email-Based Access**:
```go
if len(spec.AllowedEmails) > 0 {
    for _, e := range spec.AllowedEmails {
        includeBlocks = append(includeBlocks, &cloudflare.AccessPolicyIncludeArgs{
            Email: &cloudflare.AccessPolicyIncludeEmailArgs{
                Email: pulumi.String(e),
            },
        })
    }
}
```

Each email (or email domain) is added as a **separate Include block**. This is a requirement of the Cloudflare Pulumi provider v6 API.

**Group-Based Access**:
```go
if len(spec.AllowedGoogleGroups) > 0 {
    for _, g := range spec.AllowedGoogleGroups {
        includeBlocks = append(includeBlocks, &cloudflare.AccessPolicyIncludeArgs{
            Group: &cloudflare.AccessPolicyIncludeGroupArgs{
                Id: pulumi.String(g),
            },
        })
    }
}
```

Google Workspace groups are referenced by their group email address (e.g., `engineering@example.com`).

#### Require Rules (Additional Requirements)

**MFA Enforcement**:
```go
var requireBlocks cloudflare.AccessPolicyRequireArray
if spec.RequireMfa {
    requireBlocks = append(requireBlocks, &cloudflare.AccessPolicyRequireArgs{
        AuthMethod: &cloudflare.AccessPolicyRequireAuthMethodArgs{
            AuthMethod: pulumi.String("mfa"),
        },
    })
}
```

If `require_mfa: true`, adds a Require block enforcing that users must authenticate with MFA.

**Critical Fix Applied**: The initial implementation had an empty `AccessPolicyRequireArgs{}` block, which didn't enforce MFA. The fix adds the `AuthMethod` field with value `"mfa"` to properly enforce multi-factor authentication.

#### Policy Decision

```go
decision := "allow"
if spec.PolicyType == cloudflarezerotrustaccessapplicationv1.CloudflareZeroTrustPolicyType_BLOCK {
    decision = "deny"
}
```

Maps the protobuf enum `CloudflareZeroTrustPolicyType` to Cloudflare's policy decision string:
- `ALLOW` (enum value 0) → `"allow"`
- `BLOCK` (enum value 1) → `"deny"`

### 4. Outputs

**Implementation** (`module/application.go` lines 103-105):

```go
ctx.Export(OpApplicationId, createdAccessApplication.ID())
ctx.Export(OpPublicHostname, pulumi.String(spec.Hostname))
ctx.Export(OpPolicyId, createdAccessPolicy.ID())
```

Exports match the schema defined in `stack_outputs.proto`.

## Design Patterns

### 80/20 Principle

This implementation focuses on the **20% of Access configuration that covers 80% of use cases**:

**Included (Essential)**:
- Email domain access (e.g., `@example.com`)
- Specific email access (e.g., `user@example.com`)
- Google Workspace group access
- MFA enforcement
- Session duration configuration
- Allow/Block policy types

**Excluded (Advanced, Rare)**:
- Device posture checks (requires Cloudflare WARP client)
- Service tokens (different workflow for machine-to-machine access)
- Custom OIDC claims or SAML attributes
- Geolocation restrictions
- IP-based rules
- App launcher visibility settings

These advanced features can be added via direct Pulumi code customization if needed, but they're intentionally excluded from the protobuf API to keep it simple.

### Provider Scoping

The Cloudflare provider is created with explicit API token configuration:

```go
cloudflareProvider, err := cloudflare.NewProvider(ctx, "cloudflare_provider", &cloudflare.ProviderArgs{
    ApiToken: pulumi.String(locals.CloudflareProviderConfig.ApiToken),
})
```

This ensures:
- Credentials are scoped per resource
- Multiple Cloudflare accounts can be managed in the same Pulumi program
- API tokens are never hardcoded (passed via stack input)

### Error Handling

Errors are wrapped with context using `github.com/pkg/errors`:

```go
if err != nil {
    return nil, errors.Wrap(err, "failed to create access application")
}
```

This provides clear error messages during debugging and deployment.

### Resource Dependencies

Pulumi automatically manages dependencies via output values:

1. **Access Application** is created first
2. **Zone Lookup** can happen concurrently with Access Application creation
3. **Access Policy** depends on both Access Application ID and Account ID from Zone Lookup
4. Pulumi builds a dependency graph and executes resources in correct order

Explicit dependency is added for clarity:

```go
pulumi.DependsOn([]pulumi.Resource{createdAccessApplication})
```

## Cloudflare API Version Considerations

### Pulumi Provider v6

The implementation uses **Pulumi Cloudflare provider v6**, which introduced breaking changes from v5:

**Changes from v5 → v6**:
- Resource renamed: `cloudflare_access_application` → `cloudflare.AccessApplication`
- Include/Require rules use nested structures instead of flat arrays
- Account ID is now required for `AccessPolicy` (previously optional)

**Why v6**: It's the latest stable version with full Zero Trust Access support.

### Terraform Provider Equivalent

The Terraform implementation (in `iac/tf/`) uses **Cloudflare Terraform provider v4.x**, which has similar semantics but slightly different syntax (HCL vs. Go).

## Comparison to Terraform Implementation

| Aspect | Pulumi (Go) | Terraform (HCL) |
|--------|-------------|-----------------|
| **Language** | Go (imperative) | HCL (declarative) |
| **Type Safety** | Strong typing (compile-time checks) | Weak typing (runtime checks) |
| **Logic** | Full programming language (loops, conditionals) | Limited (for_each, dynamic blocks) |
| **State Management** | Pulumi state backend | Terraform state backend |
| **Provider Version** | v6 | v4 |

Both implementations produce **identical infrastructure**—the choice is a matter of team preference and existing tooling.

## Testing Strategy

### Unit Tests

Currently, unit tests are in `spec_test.go` at the v1 level, validating protobuf validation rules:
- Required fields (application_name, zone_id, hostname)
- Field constraints (string min length, foreign key references)

### Integration Tests

Integration testing requires:
1. Valid Cloudflare API token
2. Existing Cloudflare DNS zone
3. Configured identity provider in Cloudflare Zero Trust dashboard

Integration tests are run manually using `debug.sh` with test credentials.

### Future: E2E Tests

Potential future improvements:
- Automated E2E tests using ephemeral Cloudflare zones
- Test deployment and teardown in CI/CD pipeline
- Validation that Access Application correctly intercepts requests

## Security Considerations

### API Token Handling

**Best Practice**: API tokens are passed via environment variables or Pulumi secrets, never hardcoded.

```bash
export CLOUDFLARE_API_TOKEN="your-token"
pulumi config set cloudflare:apiToken $CLOUDFLARE_API_TOKEN --secret
```

### MFA Enforcement

**Critical Security Feature**: The MFA enforcement is implemented via Access Policy Require rules. If `require_mfa: true`, the policy **blocks access** unless the user authenticates with MFA.

**Important**: MFA enforcement relies on the identity provider supporting MFA prompting. Ensure your IdP (Google Workspace, Okta, Azure AD) is properly configured.

### Session Duration

**Security vs. UX Trade-off**: Shorter sessions (1-4 hours) are more secure but require frequent re-authentication. Longer sessions (12-24 hours) improve UX but increase risk if a device is compromised.

**Recommendation**: Use short sessions (2-4 hours) for production/sensitive applications, longer sessions (8-12 hours) for internal tools.

## Performance Considerations

### Resource Creation Time

Typical deployment times:
- **Access Application**: ~2-5 seconds
- **Zone Lookup**: <1 second (read-only)
- **Access Policy**: ~2-5 seconds
- **Total**: ~5-10 seconds

### Pulumi State Size

Each Access Application deployment adds ~2-3 KB to Pulumi state (minimal overhead).

### Cloudflare Edge Latency

Once deployed:
- Authentication checks happen at Cloudflare's edge (low latency, <50ms globally)
- No additional latency for end users compared to direct access
- Session tokens are cached in browser (subsequent requests don't re-authenticate)

## Monitoring and Observability

### Pulumi Logs

```bash
pulumi logs --follow
```

Shows real-time deployment progress and resource creation events.

### Cloudflare Access Logs

After deployment, view authentication logs in Cloudflare Zero Trust dashboard:
- **Zero Trust** → **Logs** → **Access**
- Shows every login attempt, user email, IP, country, MFA status

### Outputs for Monitoring

The exported outputs can be used for monitoring dashboards:
```bash
pulumi stack output application_id  # Use in Cloudflare API queries
pulumi stack output policy_id       # Track policy changes
```

## Future Enhancements

Potential improvements to the Pulumi implementation:

1. **Support for Multiple Policies**: Allow multiple policies per application (e.g., different rules for different paths)
2. **Device Posture Integration**: Add support for WARP client device posture checks
3. **Service Tokens**: Support for machine-to-machine authentication
4. **Custom IdP Claims**: Support for advanced OIDC/SAML claims-based rules
5. **App Launcher Configuration**: Allow customization of app launcher visibility

These enhancements would expand beyond the 80/20 principle but could be added as optional advanced fields.

## References

- **Component Documentation**: [../../docs/README.md](../../docs/README.md)
- **User Guide**: [../../README.md](../../README.md)
- **Examples**: [../../examples.md](../../examples.md)
- **Pulumi Cloudflare Provider Docs**: [pulumi.com/registry/packages/cloudflare](https://www.pulumi.com/registry/packages/cloudflare/)
- **Cloudflare Zero Trust Docs**: [developers.cloudflare.com/cloudflare-one/applications](https://developers.cloudflare.com/cloudflare-one/applications)
- **Cloudflare Access API Reference**: [developers.cloudflare.com/api/operations/access-applications-list-access-applications](https://developers.cloudflare.com/api/operations/access-applications-list-access-applications)

## Conclusion

This Pulumi implementation provides a production-ready, type-safe, and maintainable Infrastructure-as-Code solution for Cloudflare Zero Trust Access Applications. It follows the 80/20 principle, focusing on essential features while maintaining extensibility for advanced use cases.

The architecture prioritizes:
- **Simplicity**: Minimal configuration surface area
- **Security**: MFA enforcement, scoped credentials, short sessions
- **Reliability**: Proper error handling, dependency management, idempotent deployments
- **Observability**: Clear outputs, logging, and monitoring integration

