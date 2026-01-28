# Auth0EventStream Examples

This document provides comprehensive examples for configuring Auth0 Event Streams.

## Table of Contents

- [EventBridge Examples](#eventbridge-examples)
- [Webhook Examples](#webhook-examples)
- [Common Patterns](#common-patterns)

---

## EventBridge Examples

### Security Monitoring with EventBridge

Stream security-relevant events to AWS EventBridge for SIEM integration:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: security-events
  org: acme-corp
  env: production
  labels:
    team: security
    purpose: siem-integration
spec:
  destination_type: eventbridge
  subscriptions:
    - authentication.success
    - authentication.failure
    - user.blocked
    - user.unblocked
    - api.authorization.failure
  eventbridge_configuration:
    aws_account_id: "123456789012"
    aws_region: us-east-1
```

### User Lifecycle Events for Analytics

Track user lifecycle events for analytics and reporting:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: user-analytics
  org: acme-corp
  env: production
  labels:
    team: analytics
    purpose: user-tracking
spec:
  destination_type: eventbridge
  subscriptions:
    - user.created
    - user.updated
    - user.deleted
  eventbridge_configuration:
    aws_account_id: "123456789012"
    aws_region: us-west-2
```

### Multi-Region EventBridge Setup

For organizations with regional requirements:

```yaml
# US Region
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: events-us
  org: global-corp
  env: production
spec:
  destination_type: eventbridge
  subscriptions:
    - user.created
    - authentication.success
  eventbridge_configuration:
    aws_account_id: "111111111111"
    aws_region: us-east-1
---
# EU Region
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: events-eu
  org: global-corp
  env: production
spec:
  destination_type: eventbridge
  subscriptions:
    - user.created
    - authentication.success
  eventbridge_configuration:
    aws_account_id: "222222222222"
    aws_region: eu-west-1
```

---

## Webhook Examples

### Bearer Token Authentication

Simple webhook with bearer token authentication:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: user-events-webhook
  org: startup-inc
  env: production
spec:
  destination_type: webhook
  subscriptions:
    - user.created
    - user.updated
    - user.deleted
  webhook_configuration:
    webhook_endpoint: "https://api.startup.io/webhooks/auth0"
    webhook_authorization:
      method: bearer
      token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Basic Authentication

Webhook with HTTP Basic authentication:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: siem-webhook
  org: enterprise-corp
  env: production
spec:
  destination_type: webhook
  subscriptions:
    - authentication.success
    - authentication.failure
  webhook_configuration:
    webhook_endpoint: "https://siem.enterprise.com/api/v1/events"
    webhook_authorization:
      method: basic
      username: auth0-integration
      password: super-secret-password-123
```

### Slack/Teams Integration via Webhook

Send events to a workflow automation platform:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: new-user-alerts
  org: saas-company
  env: production
spec:
  destination_type: webhook
  subscriptions:
    - user.created
  webhook_configuration:
    webhook_endpoint: "https://hooks.zapier.com/hooks/catch/123456/abcdef/"
    webhook_authorization:
      method: bearer
      token: "zapier-webhook-token"
```

### Pipedream Integration

For development and testing:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: dev-events
  org: dev-team
  env: development
spec:
  destination_type: webhook
  subscriptions:
    - user.created
    - user.updated
    - authentication.success
    - authentication.failure
  webhook_configuration:
    webhook_endpoint: "https://eof28wtn4v4506o.m.pipedream.net"
    webhook_authorization:
      method: bearer
      token: "test-token-123"
```

---

## Common Patterns

### Comprehensive Audit Trail

Capture all significant events for compliance:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: audit-trail
  org: fintech-corp
  env: production
  labels:
    compliance: sox
    retention: 7years
spec:
  destination_type: eventbridge
  subscriptions:
    # User lifecycle
    - user.created
    - user.updated
    - user.deleted
    - user.blocked
    - user.unblocked
    # Authentication
    - authentication.success
    - authentication.failure
    # Authorization
    - api.authorization.success
    - api.authorization.failure
    # Management operations
    - management.client.created
    - management.client.updated
    - management.client.deleted
    - management.connection.created
    - management.connection.updated
    - management.connection.deleted
  eventbridge_configuration:
    aws_account_id: "999888777666"
    aws_region: us-east-1
```

### Login Analytics

Track authentication patterns:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: login-analytics
  org: analytics-team
  env: production
spec:
  destination_type: webhook
  subscriptions:
    - authentication.success
  webhook_configuration:
    webhook_endpoint: "https://analytics.company.com/auth0/logins"
    webhook_authorization:
      method: bearer
      token: "analytics-api-key"
```

### Failed Login Alerts

Monitor for potential security issues:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: failed-login-alerts
  org: security-team
  env: production
spec:
  destination_type: webhook
  subscriptions:
    - authentication.failure
    - user.blocked
  webhook_configuration:
    webhook_endpoint: "https://alerts.pagerduty.com/integration/auth0"
    webhook_authorization:
      method: bearer
      token: "pagerduty-integration-key"
```

### Development/Staging Environment

Minimal event stream for non-production:

```yaml
apiVersion: auth0.openmcf.org/v1
kind: Auth0EventStream
metadata:
  name: dev-events
  org: platform-team
  env: staging
spec:
  destination_type: webhook
  subscriptions:
    - user.created
  webhook_configuration:
    webhook_endpoint: "https://webhook.site/unique-id"
    webhook_authorization:
      method: bearer
      token: "dev-token"
```

---

## Deployment Commands

### Deploy with CLI

```bash
# Deploy EventBridge event stream
openmcf apply --manifest eventbridge-events.yaml \
  --auth0-provider-config auth0-creds.yaml

# Deploy webhook event stream
openmcf apply --manifest webhook-events.yaml \
  --auth0-provider-config auth0-creds.yaml
```

### Check Status

```bash
# Preview changes before applying
openmcf plan --manifest eventstream.yaml \
  --auth0-provider-config auth0-creds.yaml

# Destroy an event stream
openmcf destroy --manifest eventstream.yaml \
  --auth0-provider-config auth0-creds.yaml
```

### Auth0 Credentials File Format

```yaml
# auth0-creds.yaml
domain: your-tenant.auth0.com
clientId: your-m2m-client-id
clientSecret: your-m2m-client-secret
```

---

## Notes

### EventBridge Considerations

- **Immutable Configuration**: EventBridge settings (AWS account ID, region) cannot be changed after creation. You must delete and recreate the stream.
- **Partner Event Source**: After creation, you must associate the partner event source with an EventBridge event bus in AWS.
- **IAM Permissions**: Ensure your AWS IAM configuration allows Auth0 to publish events.

### Webhook Considerations

- **HTTPS Required**: Webhook endpoints must use HTTPS.
- **Timeout**: Auth0 expects a 2xx response within 10 seconds.
- **Retries**: Auth0 will retry failed deliveries with exponential backoff.
- **Secret Management**: Store tokens/passwords securely and rotate regularly.

