---
title: "HTTPS Forward"
description: "This preset creates the workhorse ALB listener: HTTPS on 443, terminating TLS with an ACM certificate and forwarding everything to one target group. It is the entry point most architectures start..."
type: "preset"
rank: "01"
presetSlug: "01-https-forward"
componentSlug: "lb-listener"
componentTitle: "LB Listener"
provider: "aws"
icon: "package"
order: 1
---

# HTTPS Forward

This preset creates the workhorse ALB listener: HTTPS on 443, terminating TLS
with an ACM certificate and forwarding everything to one target group. It is
the entry point most architectures start from -- add `AwsLbListenerRule`
resources against this listener's `listener_arn` output as more services
share the port.

## When to Use

- A single service (or a service mesh entry point) behind an ALB
- The starting HTTPS listener that rules will later split by host or path
- Pairing with the **02-http-redirect-to-https** preset for the standard
  80/443 duo

## Key Configuration Choices

- **Everything by reference** -- the load balancer, certificate, and target
  group all resolve from other components' outputs at deploy time, so the
  routing graph stays explicit
- **`ELBSecurityPolicy-TLS13-1-2-2021-06`** -- negotiates TLS 1.3 with a 1.2
  fallback; the sensible modern default (the AWS default still permits 1.0/1.1)
- **Single forward action** -- add weighted `targetGroups` entries later for
  blue/green shifts; the listener itself never needs replacing
- **Certificate rotation is in-place** -- pointing `certificateArn` at a new
  certificate updates the listener without downtime or rule loss

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<alb-resource-name>` | Name of the AwsAlb resource to attach to | Your AwsAlb manifest's `metadata.name` |
| `<certificate-resource-name>` | Name of the AwsCertManagerCert for the domain | Your AwsCertManagerCert manifest's `metadata.name` |
| `<target-group-resource-name>` | Name of the AwsLbTargetGroup receiving traffic | Your AwsLbTargetGroup manifest's `metadata.name` |

## Common Additions

- Add `additionalCertificateArns` to serve more domains via SNI from the same
  listener
- Change the default action to a `fixed-response` 404 once all real routes
  live in listener rules
- Add `httpHeaders.response` to set HSTS and other security headers at the
  edge

## Related Presets

- **02-http-redirect-to-https** -- the port-80 companion redirect
- **03-oidc-protected** -- the same listener with a login step in front
