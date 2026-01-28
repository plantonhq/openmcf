# Auth0Connection Examples

This document provides complete, copy-paste ready examples for configuring Auth0 connections using OpenMCF.

## Database Connections

### Basic Auth0 Database

Minimal database connection with default settings:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: user-database
  org: my-organization
  env: development
spec:
  strategy: auth0
  display_name: Sign Up with Email
  enabled_clients:
    - value: "your-app-client-id"
```

### Production Database with Security

Database connection with enhanced security settings:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: production-database
  org: my-organization
  env: production
  labels:
    team: identity
    security-level: high
spec:
  strategy: auth0
  display_name: Create Account
  enabled_clients:
    - value: "web-app-client-id"
    - value: "mobile-app-client-id"
  database_options:
    password_policy: excellent
    requires_username: false
    disable_signup: false
    brute_force_protection: true
    password_history_size: 10
    password_no_personal_info: true
    password_dictionary: true
    mfa_enabled: true
```

## Social Connections

### Google OAuth

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: google-social
  org: my-organization
  env: production
spec:
  strategy: google-oauth2
  display_name: Continue with Google
  show_as_button: true
  enabled_clients:
    - value: "web-app-client-id"
  social_options:
    client_id: "123456789-abcdef.apps.googleusercontent.com"
    client_secret: "GOCSPX-your-secret-here"
    scopes:
      - openid
      - profile
      - email
```

### GitHub OAuth

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: github-login
  org: my-organization
  env: production
spec:
  strategy: github
  display_name: Sign in with GitHub
  show_as_button: true
  enabled_clients:
    - value: "developer-portal-client-id"
  social_options:
    client_id: "Iv1.your-github-app-id"
    client_secret: "your-github-client-secret"
    scopes:
      - read:user
      - user:email
      - read:org
```

### Facebook Login

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: facebook-social
  org: my-organization
  env: production
spec:
  strategy: facebook
  display_name: Continue with Facebook
  show_as_button: true
  enabled_clients:
    - value: "consumer-app-client-id"
  social_options:
    client_id: "your-facebook-app-id"
    client_secret: "your-facebook-app-secret"
    scopes:
      - email
      - public_profile
```

### Microsoft Account

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: microsoft-personal
  org: my-organization
  env: production
spec:
  strategy: microsoft-account
  display_name: Sign in with Microsoft
  show_as_button: true
  enabled_clients:
    - value: "web-app-client-id"
  social_options:
    client_id: "your-azure-app-id"
    client_secret: "your-azure-client-secret"
    scopes:
      - openid
      - profile
      - email
```

## Enterprise Connections

### SAML with Okta

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: okta-enterprise-sso
  org: my-organization
  env: production
spec:
  strategy: samlp
  display_name: Company SSO
  is_domain_connection: true
  realms:
    - company.com
    - company.io
  enabled_clients:
    - value: "internal-app-client-id"
  saml_options:
    sign_in_endpoint: "https://company.okta.com/app/app-id/sso/saml"
    signing_cert: |
      -----BEGIN CERTIFICATE-----
      MIICmTCCAYGgAwIBAgIJAKc...
      -----END CERTIFICATE-----
    sign_out_endpoint: "https://company.okta.com/app/app-id/slo/saml"
    entity_id: "http://www.okta.com/exk123abc"
    protocol_binding: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
    sign_request: true
    signature_algorithm: rsa-sha256
    digest_algorithm: sha256
    attribute_mappings:
      email: "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"
      given_name: "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname"
      family_name: "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname"
```

### OIDC with Keycloak

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: keycloak-oidc
  org: my-organization
  env: production
spec:
  strategy: oidc
  display_name: Login with Keycloak
  is_domain_connection: true
  enabled_clients:
    - value: "internal-app-client-id"
  oidc_options:
    issuer: "https://keycloak.company.com/realms/main"
    client_id: "auth0-integration"
    client_secret: "keycloak-client-secret"
    scopes:
      - openid
      - profile
      - email
      - groups
    type: front_channel
    attribute_mappings:
      email: email
      name: name
      given_name: given_name
      family_name: family_name
```

### Azure AD / Entra ID

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: azure-ad-enterprise
  org: my-organization
  env: production
spec:
  strategy: waad
  display_name: Microsoft Work Account
  is_domain_connection: true
  realms:
    - contoso.com
  enabled_clients:
    - value: "corporate-app-client-id"
  azure_ad_options:
    client_id: "your-azure-app-id-guid"
    client_secret: "your-azure-client-secret"
    domain: "contoso.onmicrosoft.com"
    tenant_id: "your-tenant-guid"
    use_common_endpoint: false
    max_groups_to_retrieve: 100
    should_trust_email_verified: true
    api_enable_users: false
```

### Azure AD Multi-Tenant

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: azure-ad-multi-tenant
  org: my-organization
  env: production
spec:
  strategy: waad
  display_name: Sign in with Work Account
  enabled_clients:
    - value: "saas-app-client-id"
  azure_ad_options:
    client_id: "your-multi-tenant-app-id"
    client_secret: "your-azure-client-secret"
    domain: "common"
    use_common_endpoint: true
    should_trust_email_verified: true
```

## Using Foreign Key References

Auth0Connection supports foreign key references to link with Auth0Client components managed by OpenMCF. This enables declarative dependencies between connections and applications.

### Reference Auth0Client by Name

Instead of hardcoding client IDs, reference managed Auth0Client resources:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: google-social
  org: my-organization
  env: production
spec:
  strategy: google-oauth2
  display_name: Continue with Google
  enabled_clients:
    # Reference Auth0Client components - client_id resolved automatically
    - value_from:
        kind: Auth0Client
        name: web-app
    - value_from:
        kind: Auth0Client
        name: mobile-app
  social_options:
    client_id: "123456789.apps.googleusercontent.com"
    client_secret: "GOCSPX-your-secret"
    scopes:
      - openid
      - profile
      - email
```

### Mix Direct Values and References

You can combine direct client IDs with references:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: enterprise-sso
  org: my-organization
  env: production
spec:
  strategy: samlp
  display_name: Company SSO
  is_domain_connection: true
  enabled_clients:
    # Managed Auth0Client
    - value_from:
        kind: Auth0Client
        name: internal-portal
    # External/legacy client ID
    - value: "legacy-app-client-id-abc123"
  saml_options:
    sign_in_endpoint: "https://idp.company.com/sso"
    signing_cert: |
      -----BEGIN CERTIFICATE-----
      MIICmTCCAYGgAwIBAgIJAKc...
      -----END CERTIFICATE-----
```

### Complete Example: Auth0Client + Auth0Connection

First, create the Auth0Client:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Client
metadata:
  name: my-web-app
  org: my-organization
  env: production
spec:
  application_type: spa
  name: My Web Application
  callbacks:
    - "https://app.example.com/callback"
  allowed_logout_urls:
    - "https://app.example.com"
  web_origins:
    - "https://app.example.com"
```

Then, enable the connection for that client:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: user-database
  org: my-organization
  env: production
spec:
  strategy: auth0
  display_name: Sign Up with Email
  enabled_clients:
    - value_from:
        kind: Auth0Client
        name: my-web-app
  database_options:
    password_policy: good
    brute_force_protection: true
```

**Benefits of Foreign Key References:**
- No hardcoded client IDs
- Automatic dependency resolution
- Changes to Auth0Client automatically propagate
- Clear infrastructure dependencies

---

## Multi-Environment Configuration

### Development

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: user-database
  org: my-organization
  env: development
spec:
  strategy: auth0
  display_name: Development Login
  database_options:
    password_policy: low
    brute_force_protection: false
    disable_signup: false
```

### Staging

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: user-database
  org: my-organization
  env: staging
spec:
  strategy: auth0
  display_name: Staging Login
  database_options:
    password_policy: good
    brute_force_protection: true
    mfa_enabled: false
```

### Production

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0Connection
metadata:
  name: user-database
  org: my-organization
  env: production
spec:
  strategy: auth0
  display_name: Sign Up
  database_options:
    password_policy: excellent
    brute_force_protection: true
    password_history_size: 10
    password_no_personal_info: true
    password_dictionary: true
    mfa_enabled: true
```

## Deployment Commands

### Deploy a connection

```bash
openmcf apply --manifest auth0-connection.yaml \
  --auth0-provider-config auth0-creds.yaml
```

### Preview changes

```bash
openmcf plan --manifest auth0-connection.yaml \
  --auth0-provider-config auth0-creds.yaml
```

### Destroy a connection

```bash
openmcf destroy --manifest auth0-connection.yaml \
  --auth0-provider-config auth0-creds.yaml
```

## Auth0 Credential File

Create `auth0-creds.yaml`:

```yaml
domain: your-tenant.auth0.com
clientId: your-m2m-client-id
clientSecret: your-m2m-client-secret
```

Ensure the Machine-to-Machine application has the following permissions:
- `create:connections`
- `read:connections`
- `update:connections`
- `delete:connections`

