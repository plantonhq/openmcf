# Cloudflare Zero Trust Access Application

Secure internal applications with identity-based authentication at the edge using OpenMCF's unified API.

## Overview

Cloudflare Zero Trust Access (formerly Cloudflare Access) replaces traditional VPNs with application-layer access control enforced at Cloudflare's global edge network. Instead of granting network-wide access based on connection, Access evaluates **every request** against identity-based policies—checking user email, group membership, MFA status, device posture, and more—before proxying traffic to your origin.

This component provides a clean, protobuf-defined API for provisioning Zero Trust Access Applications, following the **80/20 principle**: exposing the essential configuration fields that 80% of teams need while keeping the API simple and predictable.

## Key Features

- **Identity-Based Access**: Authenticate users via email domain or Google Workspace groups
- **MFA Enforcement**: Require multi-factor authentication for sensitive applications
- **Session Management**: Configure session duration from minutes to hours
- **Policy Types**: Allow or block access based on flexible rules
- **No VPN Required**: Users access applications via browser; Cloudflare handles authentication
- **Edge Enforcement**: Authentication happens at Cloudflare's edge, close to users
- **Audit Logging**: Every authentication attempt is logged for compliance and security monitoring

## Prerequisites

1. **Cloudflare Account**: Active Cloudflare account with Zero Trust Access enabled (Free plan supports up to 50 users)
2. **DNS Zone**: A Cloudflare DNS zone for the domain you want to protect
3. **Identity Provider**: Configured identity provider (Google Workspace, Okta, Azure AD, GitHub, etc.)
4. **API Token**: Cloudflare API token with Zero Trust permissions
5. **OpenMCF CLI**: Install from [openmcf.org](https://openmcf.org)

## Quick Start

### Minimal Configuration

Protect an internal application with email domain access:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareZeroTrustAccessApplication
metadata:
  name: internal-dashboard
spec:
  application_name: "Internal Dashboard"
  zone_id: "your-cloudflare-zone-id"
  hostname: "dash.example.com"
  policy_type: ALLOW
  allowed_emails:
    - "@example.com"
  session_duration_minutes: 480  # 8 hours
  require_mfa: true
```

Deploy:

```bash
planton apply -f access-app.yaml
```

### With Google Workspace Groups

Restrict access to specific teams using Google Workspace groups:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareZeroTrustAccessApplication
metadata:
  name: staging-admin
spec:
  application_name: "Staging Admin Console"
  zone_id: "your-cloudflare-zone-id"
  hostname: "staging-admin.example.com"
  policy_type: ALLOW
  allowed_google_groups:
    - "engineering@example.com"
    - "devops@example.com"
  session_duration_minutes: 720  # 12 hours
  require_mfa: true
```

### Production Console with Strict Security

High-security application with short session duration:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareZeroTrustAccessApplication
metadata:
  name: prod-console
spec:
  application_name: "Production Console"
  zone_id: "your-cloudflare-zone-id"
  hostname: "prod.example.com"
  policy_type: ALLOW
  allowed_google_groups:
    - "prod-admins@example.com"
  session_duration_minutes: 120  # 2 hours
  require_mfa: true
```

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `application_name` | string | Display name of the Access Application (shown in logs and dashboard) |
| `zone_id` | string | Cloudflare DNS zone ID (reference from CloudflareDnsZone resource) |
| `hostname` | string | Fully qualified domain name to protect (e.g., `app.example.com`) |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `policy_type` | enum | `ALLOW` | Access policy type: `ALLOW` or `BLOCK` |
| `allowed_emails` | list(string) | `[]` | Email addresses allowed access (e.g., `user@example.com`, `@example.com` for domain) |
| `session_duration_minutes` | int32 | `1440` | Session duration in minutes (default: 24 hours) |
| `require_mfa` | bool | `false` | Whether multi-factor authentication is required |
| `allowed_google_groups` | list(string) | `[]` | Google Workspace group email addresses allowed access |

### Policy Types

- **ALLOW**: Grant access to users matching the policy rules (most common)
- **BLOCK**: Explicitly deny access to users matching the rules (for exceptions or blacklisting)

### Session Duration Guidelines

- **1-4 hours (60-240 minutes)**: High-security applications (production consoles, financial tools)
- **8 hours (480 minutes)**: Standard workday access (internal dashboards, staging environments)
- **12-24 hours (720-1440 minutes)**: Low-security internal tools (wikis, documentation sites)

## Outputs

After deployment, the following outputs are available:

- `application_id`: The unique ID of the Cloudflare Access Application
- `public_hostname`: The hostname being protected (echoes input)
- `policy_id`: The ID of the Cloudflare Access policy

Access outputs:

```bash
planton output application_id
planton output public_hostname
```

## Identity Provider Integration

### Google Workspace

Configure Google Workspace as your identity provider in the Cloudflare Zero Trust dashboard:

1. Navigate to **Settings** → **Authentication** → **Login Methods**
2. Add **Google** as an identity provider
3. Configure OAuth app and consent screen
4. Use group email addresses in `allowed_google_groups` field

### Other Identity Providers

Cloudflare Access supports:
- **Okta**: Use Okta groups with SAML/OIDC
- **Azure AD**: Use Azure AD groups
- **GitHub**: Use GitHub organization membership
- **Generic SAML 2.0 / OIDC**: Integrate any standards-compliant IdP

## Multi-Environment Pattern

### Separate Applications Per Environment

Create distinct Access Applications for each environment:

**Development**:
```yaml
metadata:
  name: my-app-dev
spec:
  application_name: "My App - Development"
  hostname: "dev.example.com"
  allowed_emails:
    - "@example.com"
  session_duration_minutes: 720
  require_mfa: false
```

**Production**:
```yaml
metadata:
  name: my-app-prod
spec:
  application_name: "My App - Production"
  hostname: "prod.example.com"
  allowed_google_groups:
    - "prod-admins@example.com"
  session_duration_minutes: 120
  require_mfa: true
```

## Common Use Cases

### 1. Company-Wide Internal Tool

Internal wiki or documentation site accessible to all employees:

```yaml
spec:
  application_name: "Company Wiki"
  hostname: "wiki.example.com"
  allowed_emails:
    - "@example.com"
  session_duration_minutes: 480
  require_mfa: false
```

### 2. Team-Specific Dashboard

Dashboard accessible only to engineering team:

```yaml
spec:
  application_name: "Engineering Dashboard"
  hostname: "eng-dash.example.com"
  allowed_google_groups:
    - "engineering@example.com"
  session_duration_minutes: 480
  require_mfa: true
```

### 3. High-Security Admin Console

Production admin console with strict security:

```yaml
spec:
  application_name: "Production Admin"
  hostname: "admin.example.com"
  allowed_google_groups:
    - "prod-admins@example.com"
  session_duration_minutes: 120
  require_mfa: true
```

## Best Practices

1. **Always Require MFA for Production**: Set `require_mfa: true` for any application accessing production data or systems
2. **Use Group-Based Access**: Leverage Google Workspace or IdP groups instead of hard-coding individual emails
3. **Shorten Sessions for Sensitive Apps**: Use 1-4 hour sessions for high-security applications
4. **Descriptive Application Names**: Use clear names like "Production Console" or "Staging Dashboard" for easy identification in logs
5. **Separate Environments**: Create distinct Access Applications for dev/staging/prod rather than reusing the same application
6. **Version Control Configs**: Store Access Application manifests in git alongside application code
7. **Plan for Emergency Access**: Configure fallback authentication methods in case your primary IdP goes down

## What This Component Does NOT Include

Following the 80/20 principle, these fields are **intentionally excluded** to keep the API simple:

- ❌ **Device Posture Checks**: Requires Cloudflare WARP client deployment (advanced feature)
- ❌ **Service Tokens**: Different workflow for machine-to-machine access
- ❌ **Custom OIDC Claims**: Rare use case; use raw Terraform/Pulumi if needed
- ❌ **Geolocation Restrictions**: Advanced security feature
- ❌ **App Launcher Visibility**: Cosmetic setting; can be added if requested

These features can be added via direct Terraform/Pulumi customization if needed.

## Cloudflare Tunnel Integration

Cloudflare Tunnel (formerly Argo Tunnel) allows you to expose internal applications (running on private IPs) securely through Cloudflare's edge without opening inbound firewall ports.

**Workflow**:
1. Deploy Cloudflare Tunnel to establish outbound connection from your network to Cloudflare
2. Map tunnel to a public hostname (e.g., `internal-app.example.com`)
3. Create Zero Trust Access Application for that hostname (this component)
4. Users access the public hostname; Cloudflare enforces authentication and proxies traffic through the tunnel

**Note**: Cloudflare Tunnel is a separate resource. This component secures the public-facing hostname with identity-based policies.

## Monitoring and Logging

Cloudflare Access provides comprehensive logging:

1. **Zero Trust Dashboard**: View authentication attempts, active sessions, user activity
2. **Access Audit Logs**: See every login (success/failure), user email, IP address, country, timestamp, MFA status
3. **Cloudflare Logpush**: Stream logs to Splunk, Datadog, S3, or your SIEM for long-term retention and alerting

**Recommended Alerts**:
- Repeated failed login attempts (potential brute force)
- Logins from unexpected countries (potential account compromise)
- Disabled users attempting access (misconfiguration or malicious activity)

## Troubleshooting

### "Access Denied" After Deployment

**Cause**: User's email or group membership doesn't match policy rules.

**Solution**: Verify:
1. User's email domain matches `allowed_emails`
2. User is member of Google Workspace group specified in `allowed_google_groups`
3. Identity provider is correctly configured in Cloudflare Zero Trust dashboard

### Session Expired Too Quickly

**Cause**: `session_duration_minutes` is set too low.

**Solution**: Increase session duration (e.g., from 60 to 480 minutes).

### MFA Not Being Enforced

**Cause**: `require_mfa: false` in spec, or IdP doesn't support MFA prompting.

**Solution**: 
1. Set `require_mfa: true` in manifest
2. Verify IdP supports MFA (Google Workspace, Okta, Azure AD all support this)
3. Re-apply manifest: `planton apply -f access-app.yaml`

### DNS Not Resolving

**Cause**: Hostname isn't properly proxied through Cloudflare.

**Solution**: Ensure DNS record for hostname is:
1. Created in Cloudflare DNS
2. Proxy status is enabled (orange cloud icon)
3. Points to origin server or Cloudflare Tunnel

## Examples

For detailed usage examples, see [examples.md](examples.md).

## Architecture Details

For in-depth architectural guidance, comparison of deployment methods (VPN vs. Zero Trust, Terraform vs. Pulumi), and production best practices, see [docs/README.md](docs/README.md).

## Terraform and Pulumi

This component supports both Pulumi (default) and Terraform:

- **Pulumi**: `iac/pulumi/` - Go-based implementation
- **Terraform**: `iac/tf/` - HCL-based implementation

Both produce identical infrastructure. Choose based on your team's preference.

## Support

- **Documentation**: [docs/README.md](docs/README.md)
- **Cloudflare Zero Trust Docs**: [developers.cloudflare.com/cloudflare-one/applications](https://developers.cloudflare.com/cloudflare-one/applications)
- **OpenMCF**: [openmcf.org](https://openmcf.org)

## License

This component is part of OpenMCF and follows the same license.

