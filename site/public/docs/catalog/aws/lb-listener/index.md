---
title: "LB Listener"
description: "LB Listener deployment documentation"
icon: "package"
order: 100
componentName: "awslblistener"
---

# AWS LB Listener

Deploys an ELBv2 listener: the port/protocol entry point on a load balancer,
with TLS certificates (default plus SNI), an optional authentication step,
and the default action taken when no listener rule matches. One kind serves
both families -- ALB listeners take the full action set; NLB listeners
forward only.

## What Gets Created

When you deploy an AwsLbListener resource, Planton provisions:

- **Listener** — an `aws_lb_listener` / `lb.Listener` on the referenced load
  balancer, with the specified port, protocol, TLS configuration, and
  default-action chain
- **SNI certificate attachments** — one listener-certificate attachment per
  entry in `additionalCertificateArns`, for serving multiple domains from one
  listener
- **Listener attributes** — TCP idle timeout (NLB) and HTTP header
  injection/override attributes (ALB), when configured

Per-service routing is **not** created here — attach `AwsLbListenerRule`
resources to this listener's `listener_arn` output.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **A load balancer**: an `AwsAlb` (default reference) or `AwsNlb`.
- **An ACM certificate** (e.g. `AwsCertManagerCert`) for `HTTPS` and `TLS` listeners.
- **A target group** (`AwsLbTargetGroup`) for forward actions.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLbListener
metadata:
  name: https
spec:
  region: us-west-2
  loadBalancerArn:
    valueFrom:
      kind: AwsAlb
      name: main-alb
      fieldPath: status.outputs.load_balancer_arn
  port: 443
  protocol: HTTPS
  certificateArn:
    valueFrom:
      kind: AwsCertManagerCert
      name: main-cert
      fieldPath: status.outputs.cert_arn
  defaultActions:
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
planton apply -f listener.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the load balancer's region. | Required; non-empty |
| `loadBalancerArn` | `string \| valueFrom` | The load balancer this listener attaches to. Defaults to referencing an `AwsAlb`'s `load_balancer_arn` output; use explicit `valueFrom` for `AwsNlb`. Immutable. | Required |
| `port` | `int` | Port the listener accepts traffic on. | Required; 1–65535 |
| `protocol` | `string` | Listener protocol; decides the family and allowed actions. | Required. One of: `HTTP`, `HTTPS` (ALB), `TCP`, `UDP`, `TCP_UDP`, `TLS` (NLB) |
| `defaultActions` | `object[]` | Action chain for requests no rule matches. Auth actions run first; the chain ends in exactly one of `forward`, `redirect`, or `fixed-response`. NLB protocols: exactly one `forward`. | Required; min 1 item |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `certificateArn` | `string \| valueFrom` | — | Default server certificate. Required for `HTTPS`/`TLS`; not valid otherwise. References `AwsCertManagerCert.cert_arn` by default. |
| `additionalCertificateArns` | `(string \| valueFrom)[]` | `[]` | SNI certificates for serving several domains from one listener. |
| `sslPolicy` | `string` | AWS default | TLS security policy. Recommended: `ELBSecurityPolicy-TLS13-1-2-2021-06`. `HTTPS`/`TLS` only. |
| `alpnPolicy` | `string` | — | NLB `TLS` only: `HTTP1Only`, `HTTP2Only`, `HTTP2Optional`, `HTTP2Preferred`, `None`. |
| `mutualAuthentication` | `object` | off | ALB `HTTPS` only: mTLS `mode` (`off`, `passthrough`, `verify`), `trustStoreArn` (required for `verify`), expiry/CA-advertisement options. |
| `tcpIdleTimeoutSeconds` | `int` | `350` | NLB `TCP` only: idle timeout, 60–6000. |
| `httpHeaders` | `object` | — | ALB only: request-header injection (TLS/mTLS details toward targets) and response-header overrides (CORS, HSTS, CSP, X-Frame-Options). |
| `defaultActions[].type` | `string` | — | `forward`, `redirect`, `fixed-response`, `authenticate-cognito`, `authenticate-oidc`, or `jwt-validation`. Exactly one matching config block must be set. |
| `defaultActions[].forward.targetGroups` | `object[]` | — | 1–5 weighted destinations; `arn` references `AwsLbTargetGroup.target_group_arn` by default, `weight` 0–999 (default 1). |
| `defaultActions[].forward.stickiness` | `object` | disabled | Pins clients to the group that served their first request in a weighted forward. |
| `defaultActions[].redirect` | `object` | — | `statusCode` (`HTTP_301`/`HTTP_302`) plus optional protocol/port/host/path/query overrides; untouched parts pass through. |
| `defaultActions[].fixedResponse` | `object` | — | `contentType`, `statusCode` (2xx/4xx/5xx, default `503`), `messageBody` (max 1024 chars). |
| `defaultActions[].authenticateCognito` | `object` | — | `userPoolArn` (references `AwsCognitoUserPool`), `userPoolClientId`, `userPoolDomain`, scope/session options. |
| `defaultActions[].authenticateOidc` | `object` | — | Issuer, authorization/token/user-info endpoints, `clientId`, `clientSecret` (handled as a secret end to end), scope/session options. |
| `defaultActions[].jwtValidation` | `object` | — | `issuer`, `jwksEndpoint`, up to 10 additional required claims. |

## Examples

### HTTP-to-HTTPS redirect listener

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLbListener
metadata:
  name: http-redirect
spec:
  region: us-west-2
  loadBalancerArn:
    valueFrom:
      kind: AwsAlb
      name: main-alb
      fieldPath: status.outputs.load_balancer_arn
  port: 80
  protocol: HTTP
  defaultActions:
    - type: redirect
      redirect:
        statusCode: HTTP_301
        protocol: HTTPS
        port: "443"
```

### Multi-domain HTTPS with SNI and a 404 default

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLbListener
metadata:
  name: https-multi-domain
spec:
  region: us-west-2
  loadBalancerArn:
    valueFrom:
      kind: AwsAlb
      name: main-alb
      fieldPath: status.outputs.load_balancer_arn
  port: 443
  protocol: HTTPS
  sslPolicy: ELBSecurityPolicy-TLS13-1-2-2021-06
  certificateArn:
    valueFrom:
      kind: AwsCertManagerCert
      name: example-com
      fieldPath: status.outputs.cert_arn
  additionalCertificateArns:
    - valueFrom:
        kind: AwsCertManagerCert
        name: example-org
        fieldPath: status.outputs.cert_arn
  defaultActions:
    - type: fixed-response
      fixedResponse:
        contentType: text/plain
        statusCode: "404"
        messageBody: Not found
```

### NLB TLS listener

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLbListener
metadata:
  name: nlb-tls
spec:
  region: us-west-2
  loadBalancerArn:
    valueFrom:
      kind: AwsNlb
      name: edge-nlb
      fieldPath: status.outputs.load_balancer_arn
  port: 443
  protocol: TLS
  certificateArn:
    valueFrom:
      kind: AwsCertManagerCert
      name: edge-cert
      fieldPath: status.outputs.cert_arn
  alpnPolicy: HTTP2Preferred
  defaultActions:
    - type: forward
      forward:
        targetGroups:
          - arn:
              valueFrom:
                kind: AwsLbTargetGroup
                name: backend
                fieldPath: status.outputs.target_group_arn
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `listener_arn` | ARN of the listener — what `AwsLbListenerRule` resources attach through |

## Related Components

- [AwsAlb](/docs/catalog/aws/alb) — the Application Load Balancer this listener attaches to by default
- [AwsNlb](/docs/catalog/aws/nlb) — the Network Load Balancer alternative (forward-only actions)
- [AwsLbTargetGroup](/docs/catalog/aws/lb-target-group) — the destination of forward actions
- [AwsLbListenerRule](/docs/catalog/aws/lb-listener-rule) — per-service routing attached to this listener
- [AwsCertManagerCert](/docs/catalog/aws/certificate) — provides the default and SNI certificates
- [AwsCognitoUserPool](/docs/catalog/aws/cognito-user-pool) — the user pool behind `authenticate-cognito` actions
