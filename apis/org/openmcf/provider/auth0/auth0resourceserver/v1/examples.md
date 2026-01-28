# Auth0 Resource Server Examples

This document provides practical examples for common Auth0 Resource Server configurations.

## Basic Examples

### Simple API

A minimal API configuration with default settings:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: simple-api
spec:
  identifier: https://api.mycompany.com/
```

### API with Custom Name

Specify a friendly display name:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: backend-api
spec:
  identifier: https://api.mycompany.com/backend
  name: Backend Services API
```

## Token Configuration Examples

### Short-Lived Tokens

For high-security scenarios with short token lifetimes:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: secure-api
spec:
  identifier: https://api.secure.com/
  name: Secure API
  signing_alg: RS256
  token_lifetime: 900         # 15 minutes
  token_lifetime_for_web: 300 # 5 minutes for implicit flows
```

### Refresh Token Support

Enable refresh tokens for mobile or desktop applications:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: mobile-api
spec:
  identifier: https://api.mobile.com/
  name: Mobile App API
  signing_alg: RS256
  token_lifetime: 3600
  allow_offline_access: true
```

### HS256 Signing

Use symmetric signing (not recommended for most cases):

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: hs256-api
spec:
  identifier: https://api.internal.com/
  name: Internal API
  signing_alg: HS256
  token_lifetime: 86400
```

## Scopes/Permissions Examples

### CRUD Scopes

Standard create, read, update, delete permissions:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: crud-api
spec:
  identifier: https://api.example.com/resources
  name: Resource Management API
  scopes:
    - name: create:resources
      description: Create new resources
    - name: read:resources
      description: Read resource data
    - name: update:resources
      description: Update existing resources
    - name: delete:resources
      description: Delete resources
```

### Domain-Specific Scopes

Scopes for an e-commerce API:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: ecommerce-api
spec:
  identifier: https://api.shop.com/
  name: E-Commerce API
  scopes:
    - name: read:products
      description: View product catalog
    - name: write:products
      description: Manage product listings
    - name: read:orders
      description: View order history
    - name: write:orders
      description: Create and update orders
    - name: process:payments
      description: Process payment transactions
    - name: admin:inventory
      description: Full inventory management access
```

### API with Admin Scopes

Include administrative permissions:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: admin-api
spec:
  identifier: https://api.admin.com/
  name: Admin API
  scopes:
    - name: read:users
      description: View user accounts
    - name: write:users
      description: Create and update users
    - name: delete:users
      description: Delete user accounts
    - name: read:logs
      description: View audit logs
    - name: admin:system
      description: Full system administration
```

## RBAC Examples

### Basic RBAC Setup

Enable RBAC with permission claims in tokens:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: rbac-api
spec:
  identifier: https://api.rbac.com/
  name: RBAC-Enabled API
  enforce_policies: true
  token_dialect: access_token_authz
  scopes:
    - name: read:data
      description: Read data
    - name: write:data
      description: Write data
```

### RBAC with RFC 9068 Tokens

Use IETF-compliant token format with RBAC:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: compliant-api
spec:
  identifier: https://api.compliant.com/
  name: IETF-Compliant API
  signing_alg: RS256
  enforce_policies: true
  token_dialect: rfc9068_profile_authz
  scopes:
    - name: read:records
      description: Read records
    - name: write:records
      description: Write records
```

## Multi-Tenant Examples

### Per-Environment APIs

Development environment:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: dev-api
  env: development
spec:
  identifier: https://api.dev.example.com/
  name: Development API
  token_lifetime: 86400  # Longer tokens for dev
  allow_offline_access: true
```

Production environment:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: prod-api
  env: production
spec:
  identifier: https://api.example.com/
  name: Production API
  signing_alg: RS256
  token_lifetime: 3600  # Shorter tokens for prod
  enforce_policies: true
  token_dialect: access_token_authz
```

## Complete Production Example

A fully-configured production API with all recommended settings:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: production-api
  org: acme-corp
  env: production
  labels:
    team: platform
    cost-center: engineering
spec:
  identifier: https://api.acme.com/v1
  name: ACME Production API v1
  
  # Token Configuration
  signing_alg: RS256
  token_lifetime: 3600           # 1 hour
  token_lifetime_for_web: 1800   # 30 minutes for web
  
  # Access Control
  skip_consent_for_verifiable_first_party_clients: true
  allow_offline_access: true
  enforce_policies: true
  token_dialect: access_token_authz
  
  # Scopes
  scopes:
    # User scopes
    - name: read:profile
      description: Read user profile information
    - name: write:profile
      description: Update user profile
    
    # Data scopes
    - name: read:data
      description: Read application data
    - name: write:data
      description: Create and update data
    - name: delete:data
      description: Delete data
    
    # Admin scopes
    - name: admin:users
      description: Full user administration
    - name: admin:settings
      description: Manage application settings
```

## Integration with Auth0Client

When creating APIs, you often need to grant access to applications. Use Auth0Client's `api_grants` to authorize access:

```yaml
# First, create the Resource Server
apiVersion: auth0.openmcf.org/v1
kind: Auth0ResourceServer
metadata:
  name: my-api
spec:
  identifier: https://api.example.com/
  scopes:
    - name: read:data
      description: Read data
---
# Then, create a client that can access it
apiVersion: auth0.openmcf.org/v1
kind: Auth0Client
metadata:
  name: backend-service
spec:
  application_type: non_interactive
  grant_types:
    - client_credentials
  api_grants:
    - audience: https://api.example.com/
      scopes:
        - read:data
```
