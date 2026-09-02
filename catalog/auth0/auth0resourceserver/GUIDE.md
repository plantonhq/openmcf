# Auth0 Resource Server Guide

## Security
## Platform Security Posture

The certifications below are Auth0's own published claims about their hosted platform (verify current status on Auth0's compliance page). They describe the vendor's service — never this Planton component, and never your deployment: configuring this resource does not make your application certified, authorized, or compliant with any framework.

Auth0's published certifications and security standards:

- SOC 2 Type II (annual audit)
- ISO 27001, ISO 27018 (privacy controls)
- HIPAA BAA available on enterprise plans
- PCI DSS Level 1 Service Provider
- FedRAMP Authorized (moderate baseline)
- CSA STAR Level 2
- GDPR compliant with Data Processing Agreement

## Data Protection

- **Data residency**: US, EU, AU regions
- **Encryption in transit**: TLS 1.2+
- **Encryption at rest**: AES-256
- **Penetration testing**: Annual third-party assessments

## Resource Server-Specific Security Notes

### Signing Algorithms

Resource servers support two token signing algorithms with different security characteristics:

- **RS256 (recommended)**: Asymmetric signing using RSA public/private key pairs. The API validates tokens using Auth0's public JWKS endpoint. No shared secret required. Supports key rotation without coordination.
- **HS256**: Symmetric signing using a shared secret. The API must store the signing secret to validate tokens. Key rotation requires coordinated updates between Auth0 and all consuming APIs.

Use RS256 unless you have a specific requirement for HS256. The signing secret for HS256 resource servers is sensitive and must be treated as a credential.

### Scope Design

Scopes define the permissions available for an API. Follow least-privilege principles:

- Define granular scopes (e.g., `read:orders`, `write:orders`) rather than coarse ones (`admin`).
- Scope names are arbitrary strings but conventionally use `action:resource` format.
- Scopes are requested at authentication time and granted based on client grants or user consent.

### Token Validation

APIs must validate access tokens on every request. Validation must include:

- Signature verification against Auth0's JWKS endpoint (RS256) or shared secret (HS256)
- Issuer (`iss`) claim matches your Auth0 domain
- Audience (`aud`) claim matches the resource server identifier
- Expiration (`exp`) claim is in the future

## Permissions
## Management API Scopes

Auth0 resource server resources require the following Management API scopes for CRUD operations. These scopes must be granted to the M2M application used for infrastructure automation.

| Operation | Scope | Description |
|-----------|-------|-------------|
| Read | `read:resource_servers` | List and retrieve API definitions and their scopes |
| Create | `create:resource_servers` | Create new APIs with identifier, scopes, and signing config |
| Update | `update:resource_servers` | Modify API settings, scopes, and token configuration |
| Delete | `delete:resource_servers` | Remove APIs from the tenant |

## Additional Scopes

Managing client access to resource servers requires client grant scopes:

| Operation | Scope | Description |
|-----------|-------|-------------|
| Read grants | `read:client_grants` | List which clients can access this API |
| Create grants | `create:client_grants` | Grant a client access to this API with specific scopes |
| Update grants | `update:client_grants` | Modify the scopes a client has for this API |
| Delete grants | `delete:client_grants` | Revoke a client's access to this API |

## Minimum Required Scopes

For basic lifecycle management (create, read, update, delete), the minimum required scopes are:

```
read:resource_servers create:resource_servers update:resource_servers delete:resource_servers
```

If the automation also manages which clients can access the API, add:

```
read:client_grants create:client_grants update:client_grants delete:client_grants
```

## Compliance
## Regulatory Frameworks

The table below records Auth0's own published compliance posture for their hosted platform, as Auth0 states it (verify current status in Auth0's trust documentation). These are vendor facts about the service this component configures — not properties of the component or of your deployment, and nothing here transfers to your application without your own assessment.

Auth0's published framework posture:

| Framework | Status | Notes |
|-----------|--------|-------|
| SOC 2 Type II | Certified | Annual audit cycle |
| ISO 27001:2022 | Certified | Information security management |
| ISO 27018:2019 | Certified | PII protection in public clouds |
| HIPAA | BAA Available | Enterprise plans only |
| PCI DSS Level 1 | Certified | Service provider certification |
| FedRAMP Moderate | Authorized | US government workloads |
| CSA STAR Level 2 | Certified | Cloud security assurance |
| GDPR | Compliant | Data Processing Agreement available |
| CCPA | Compliant | California consumer privacy |
| Privacy Shield | Historical | Deprecated framework, replaced by SCCs |

## Resource Server-Specific Compliance Notes

### Access Token Content

Access tokens issued for resource servers may contain user identifiers and granted scopes. If your API's access tokens include custom claims with PII (via Actions or Rules), ensure token handling complies with GDPR data minimization principles. Only include claims that the API needs.

### Scope-Based Access Control

Resource server scopes provide the foundation for authorization decisions. For compliance-sensitive workloads, document the mapping between scopes and business-level access rights. Auditors may request evidence that scope assignments follow least-privilege principles.

### Audit Trail

All resource server CRUD operations are logged in Auth0 tenant logs. Token issuance events for each resource server are also logged, providing an audit trail of API access. Log retention depends on plan tier (2 days free, 30 days enterprise).

### API Identifier Stability

The resource server identifier (audience) is immutable after creation. Changing an API's audience requires creating a new resource server. This is relevant for compliance documentation that references specific API identifiers.

## Cost
## Pricing Model

Auth0 pricing is based on Monthly Active Users (MAUs), not on the number of resources created. Resource servers (APIs) are free API objects with no per-resource cost.

## Free Tier

The Auth0 Free plan includes:

- 25,000 MAUs
- 1 tenant
- Unlimited resource servers (APIs)
- Unlimited scopes per resource server

## Cost Impact

Creating, updating, or deleting Auth0 resource server resources has no direct billing impact. There is no charge per API definition or per scope/permission defined.

The only cost driver is the number of monthly active users authenticating through your tenant. API definitions and their associated scopes are purely configuration objects.

## Token Volume Considerations

While resource servers are free to define, the access tokens issued for these APIs consume tenant resources. High-volume API usage patterns should consider:

| Factor | Impact |
|--------|--------|
| Token issuance rate | Contributes to tenant rate limits |
| Token lifetime | Shorter tokens increase issuance frequency |
| Number of scopes per token | Increases token size but not cost |
| M2M token grants | Counted against plan M2M token quota |

## Rate Limits

The Auth0 Management API enforces rate limits on resource server operations. Creating or updating resource servers counts against the Management API rate limit (varies by plan tier). This affects automation speed but not cost.
