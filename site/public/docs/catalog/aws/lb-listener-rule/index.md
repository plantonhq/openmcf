---
title: "LB Listener Rule"
description: "LB Listener Rule deployment documentation"
icon: "package"
order: 100
componentName: "awslblistenerrule"
---

# AWS LB Listener Rule

Deploys an ALB listener rule: a condition-action pair evaluated in priority
order before the listener's default action. Rules are how one Application
Load Balancer serves many services -- each service brings its own rule
("host `api.example.com` forwards to my target group") and removes it when
it goes away, while the shared listener stays untouched.

## What Gets Created

When you deploy an AwsLbListenerRule resource, Planton provisions:

- **Listener rule** — an `aws_lb_listener_rule` / `lb.ListenerRule` on the
  referenced listener, with the specified priority, condition blocks, action
  chain, and optional host/URL rewrite transforms

Rules apply to Application Load Balancers only — NLB listeners route purely
by port/protocol and take no rules.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **An ALB listener** (`AwsLbListener` with `HTTP` or `HTTPS` protocol) to attach to.
- **A target group** (`AwsLbTargetGroup`) for forward actions.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsLbListenerRule
metadata:
  name: api-route
spec:
  region: us-west-2
  listenerArn:
    valueFrom:
      kind: AwsLbListener
      name: https
      fieldPath: status.outputs.listener_arn
  priority: 10
  conditions:
    - pathPattern:
        values:
          - /api/*
  actions:
    - type: forward
      forward:
        targetGroups:
          - arn:
              valueFrom:
                kind: AwsLbTargetGroup
                name: api
                fieldPath: status.outputs.target_group_arn
```

```shell
planton apply -f rule.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the listener's region. | Required; non-empty |
| `listenerArn` | `string \| valueFrom` | The listener this rule attaches to. References `AwsLbListener.listener_arn` by default. Immutable. | Required |
| `conditions` | `object[]` | What a request must match. Blocks AND together; values within one block OR together. Exactly one criterion per block. | Required; 1–5 items |
| `actions` | `object[]` | Action chain for matching requests. Auth actions run first; the chain ends in exactly one of `forward`, `redirect`, or `fixed-response`. | Required; min 1 item |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `priority` | `int` | next free slot | Evaluation order, 1–50000, unique per listener, lower first. Set explicitly when rules shadow each other. |
| `conditions[].hostHeader` | `object` | — | Host header match: wildcard `values` and/or `regexValues` (regex requires the load balancer attribute). |
| `conditions[].pathPattern` | `object` | — | Path match (query string excluded): wildcard `values` and/or `regexValues`. |
| `conditions[].httpHeader` | `object` | — | Arbitrary header match: `httpHeaderName` plus `values`/`regexValues`. |
| `conditions[].httpRequestMethod` | `object` | — | Method match (`values`), case-sensitive and exact. |
| `conditions[].queryString` | `object` | — | Query-string `pairs` (key optional, value required; wildcards allowed). |
| `conditions[].sourceIp` | `object` | — | Client CIDR match. Uses the connecting address, not `X-Forwarded-For`. |
| `actions[].type` | `string` | — | `forward`, `redirect`, `fixed-response`, `authenticate-cognito`, `authenticate-oidc`, or `jwt-validation`. Exactly one matching config block must be set. |
| `actions[].forward.targetGroups` | `object[]` | — | 1–5 weighted destinations; `arn` references `AwsLbTargetGroup.target_group_arn` by default, `weight` 0–999 (default 1). |
| `actions[].forward.stickiness` | `object` | disabled | Pins clients to the group that served their first request in a weighted forward. |
| `actions[].redirect` | `object` | — | `statusCode` (`HTTP_301`/`HTTP_302`) plus optional protocol/port/host/path/query overrides; untouched parts pass through. |
| `actions[].fixedResponse` | `object` | — | `contentType`, `statusCode` (2xx/4xx/5xx, default `503`), `messageBody` (max 1024 chars). |
| `actions[].authenticateCognito` | `object` | — | `userPoolArn` (references `AwsCognitoUserPool`), `userPoolClientId`, `userPoolDomain`, scope/session options. HTTPS listeners only. |
| `actions[].authenticateOidc` | `object` | — | Issuer, endpoints, `clientId`, `clientSecret` (handled as a secret end to end), scope/session options. HTTPS listeners only. |
| `actions[].jwtValidation` | `object` | — | `issuer`, `jwksEndpoint`, up to 10 additional required claims. HTTPS listeners only. |
| `transforms` | `object[]` | `[]` | Request rewrites before the action runs: at most one `host-header-rewrite` and one `url-rewrite`, each a regex find-and-replace with capture groups. |

## Examples

### OIDC gate on the admin path

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsLbListenerRule
metadata:
  name: admin-auth
spec:
  region: us-west-2
  listenerArn:
    valueFrom:
      kind: AwsLbListener
      name: https
      fieldPath: status.outputs.listener_arn
  priority: 5
  conditions:
    - pathPattern:
        values:
          - /admin/*
  actions:
    - type: authenticate-oidc
      authenticateOidc:
        issuer: https://idp.example.com
        authorizationEndpoint: https://idp.example.com/authorize
        tokenEndpoint: https://idp.example.com/oauth/token
        userInfoEndpoint: https://idp.example.com/userinfo
        clientId: admin-portal
        clientSecret: replace-with-client-secret
    - type: forward
      forward:
        targetGroups:
          - arn:
              valueFrom:
                kind: AwsLbTargetGroup
                name: admin
                fieldPath: status.outputs.target_group_arn
```

### Path stripping with a URL rewrite

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsLbListenerRule
metadata:
  name: svc-mount
spec:
  region: us-west-2
  listenerArn:
    valueFrom:
      kind: AwsLbListener
      name: https
      fieldPath: status.outputs.listener_arn
  priority: 30
  conditions:
    - pathPattern:
        values:
          - /svc/*
  transforms:
    - type: url-rewrite
      urlRewrite:
        regex: ^/svc/(.*)$
        replace: /$1
  actions:
    - type: forward
      forward:
        targetGroups:
          - arn:
              valueFrom:
                kind: AwsLbTargetGroup
                name: svc
                fieldPath: status.outputs.target_group_arn
```

### Maintenance mode for one host

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsLbListenerRule
metadata:
  name: maintenance
spec:
  region: us-west-2
  listenerArn:
    valueFrom:
      kind: AwsLbListener
      name: https
      fieldPath: status.outputs.listener_arn
  priority: 1
  conditions:
    - hostHeader:
        values:
          - app.example.com
  actions:
    - type: fixed-response
      fixedResponse:
        contentType: text/html
        statusCode: "503"
        messageBody: <h1>Down for maintenance</h1>
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `rule_arn` | ARN of the rule — the handle audit tooling and imports reference |
| `priority` | The priority AWS assigned — meaningful when the spec left `priority` unset |

## Related Components

- [AwsLbListener](/docs/catalog/aws/lb-listener) — the listener this rule attaches to
- [AwsLbTargetGroup](/docs/catalog/aws/lb-target-group) — the destination of forward actions
- [AwsAlb](/docs/catalog/aws/alb) — the Application Load Balancer carrying the listener
- [AwsCognitoUserPool](/docs/catalog/aws/cognito-user-pool) — the user pool behind `authenticate-cognito` actions
