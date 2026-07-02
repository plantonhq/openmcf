# AWS LB Listener Rule: The Unit of Per-Service Routing

## What a Listener Rule Is

An ALB listener rule is a condition-action pair attached to a listener:
requests matching the conditions get the rule's actions; everything else
falls through to lower-priority rules and finally the listener's default
action. Rules are how one Application Load Balancer -- one set of IPs, one
certificate, one DNS name -- serves an entire environment of services.

Rules are an ALB-only concept. Network Load Balancer listeners route purely
by port and protocol; there is no request to inspect at Layer 4, so
`AwsLbListenerRule` resources only attach to HTTP/HTTPS listeners.

## Why a First-Class Component

The rule's lifecycle belongs to the service, not the listener. When a
service deploys, it brings a target group and a rule ("host
`api.example.com` forwards to my group"); when it is decommissioned, both
leave. The listener -- shared, certificate-bearing, slow-changing -- stays
untouched throughout. Folding rules into the listener definition would force
every service deployment to edit one shared manifest, serializing exactly
the operations that should be independent. Many rules per listener, each
with its own lifecycle, is the reason the rule is a kind of its own.

A practical division of labor follows:

- The **listener** carries TLS material and a `fixed-response` 404 default.
- Each **service** carries its own rule and target group.
- Cross-cutting concerns (an OIDC gate on `/admin/*`, a maintenance-mode
  fixed response) are rules too, deployed and removed like any other.

## How Matching Works

A rule takes 1-5 condition blocks, and a request must satisfy **all** of
them (conditions AND together). Within one condition, multiple values **OR**
together. So "host is `api.example.com` AND path starts with `/v2/`" is two
condition blocks; "path is `/v1/*` or `/v2/*`" is one block with two values.
Exactly one criterion is set per block -- host header, path pattern, HTTP
header, HTTP method, query string, or source IP -- mirroring the AWS API
shape.

Details that matter in practice:

- **Wildcards vs regex**: host, path, and header conditions accept wildcard
  patterns (`*` spans any characters, `?` one character) everywhere; regex
  patterns additionally require regex matching to be enabled on the load
  balancer's attributes. Prefer wildcards when they express the intent.
- **Path conditions ignore the query string**; match query parameters with a
  `queryString` condition instead.
- **`sourceIp` uses the connecting address**, not `X-Forwarded-For` -- behind
  CloudFront or another proxy, the connecting address is the proxy's.
- **Method matching is exact and case-sensitive**, per AWS.

## Priorities

Priorities order evaluation: 1-50000, unique per listener, lower first. When
the spec omits `priority`, AWS assigns the next free slot after the current
highest -- fine for disjoint routes (different hosts), but shadowing rules
need explicit numbers: a `/api/*` rule must outrank a catch-all `/*` rule or
it will never fire. Leaving gaps between assigned priorities (10, 20, 30...)
keeps room to slot rules in later without renumbering.

The `priority` stack output reports the value AWS actually assigned, which
is how an auto-assigned rule's position becomes visible to tooling.

## Actions and Transforms

The action chain is the same discriminated-union shape as the listener's
default actions -- `type` plus exactly one matching configuration block --
with authentication steps (`authenticate-cognito`, `authenticate-oidc`,
`jwt-validation`) running before a terminal `forward`, `redirect`, or
`fixed-response`. Everything the listener's actions support works here,
scoped to just the matching requests. Weighted forwards make per-route
canaries: two target groups, weights 95/5, optional group stickiness so a
session does not flap between versions mid-canary.

Transforms rewrite the request before the action runs: at most one
`host-header-rewrite` and one `url-rewrite` per rule, each a regex
find-and-replace with capture-group support. The classic use is path
stripping for services that expect to be mounted at `/` -- match `/svc/*`,
rewrite `^/svc/(.*)$` to `/$1`, forward.

## Mutability

The listener is the rule's only create-only field: moving a rule to a
different listener replaces it. Priority, conditions, actions, and
transforms all update in place -- re-prioritizing routing or tightening a
match never re-creates the rule.

## Dual-Engine Implementation

`AwsLbListenerRule` ships both a Terraform/OpenTofu module and a Pulumi (Go)
module at behavioral parity. Both express conditions and action chains
identically, apply the same transform semantics, and export the same outputs
(`rule_arn`, `priority`). Whichever engine a team standardizes on, the rule
behaves identically.
