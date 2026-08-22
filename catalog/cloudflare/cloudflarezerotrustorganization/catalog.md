# Cloudflare Zero Trust Organization

The Zero Trust organization: the Access login experience (team domain, login page design, session defaults, MFA policy) plus the Access service-key rotation cadence. A singleton upsert -- applying mutates the organization the account already carries, and destroy abandons the configuration.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Organization configuration** -- one `cloudflare_zero_trust_organization` (upsert of the account/zone singleton)
- **Key-rotation cadence** -- one `cloudflare_zero_trust_access_key_configuration` (only when `keyRotationIntervalDays` is set)

## Prerequisites

- **Zero Trust enabled on the account** (the team-name onboarding step -- it creates the organization this resource configures)
- **A Cloudflare API token** with Account → Access: Organizations, Identity Providers, and Groups → Edit

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustOrganization
metadata:
  name: acme-zero-trust-org
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  authDomain: acme
  name: Acme Zero Trust
  sessionDuration: 24h
```

```shell
planton apply -f organization.yaml
```

## Configuration Reference

### Scope (exactly one)

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The account-wide organization (common case). | 32-hex; exactly one scope. |
| `zoneId` | string/ref | A zone-scoped organization. | Exactly one scope; not importable. |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `authDomain` | string | Team domain without `.cloudflareaccess.com`. | Changing it breaks every Access login -- see GUIDE. |
| `name` | string | Display name on the login page. | |
| `sessionDuration` | string | Access session lifetime. | Go duration (e.g. `24h`, `2h45m`). |
| `warpAuthSessionDuration` | string | WARP session lifetime. | Minutes/hours only. |
| `userSeatExpirationInactiveTime` | string | Seat release after inactivity. | Min `730h`. |
| `denyUnmatchedRequests` | bool | Deny requests matching no application. | + `denyUnmatchedRequestsExemptedZoneNames`. |
| `customPages` | object | Custom block/denied page UIDs. | |
| `loginDesign` | object | Login page colors, logo, header/footer. | |
| `mfaConfig` | object | Allowed authenticators + MFA sessions. | totp, biometrics, security_key, ssh_piv_key. |
| `mfaSshPivKeyRequirements` | object | Hardware PIV key constraints. | pin/touch policies, key types/sizes CEL-walled. |
| `allowAuthenticateViaWarp` | bool | WARP session counts as authentication. | |
| `autoRedirectToIdentity` | bool | Skip the IdP picker with one provider. | |
| `mfaRequiredForAllApps` | bool | Global MFA by default. | Needs MFA config + session duration. |
| `isUiReadOnly` | bool | Lock the dashboard read-only. | + `uiReadOnlyToggleReason`. |
| `keyRotationIntervalDays` | int | Access service-key rotation cadence. | 21-365; account scope only. |

## Destroy Semantics

Destroy is a NO-OP for both folded surfaces: the live configuration stays exactly as last applied. Nothing is reverted, nothing is deleted.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `auth_domain` | string | The team domain -- what Access applications and WARP enrollment reference |
| `account_id` | string | The account applied to (the singleton's identity) |

## Related Components

- [Cloudflare Zero Trust Access Identity Provider](/docs/catalog/cloudflare/cloudflarezerotrustaccessidentityprovider) -- how users sign in
- [Cloudflare Zero Trust Access Application](/docs/catalog/cloudflare/cloudflarezerotrustaccessapplication) -- what sits behind the login
- [Cloudflare Zero Trust Access Policy](/docs/catalog/cloudflare/cloudflarezerotrustaccesspolicy) -- who gets in
