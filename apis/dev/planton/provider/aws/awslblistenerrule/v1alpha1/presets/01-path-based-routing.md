# Path-Based Routing

This preset routes one URL prefix to one service: requests whose path matches
`/api/*` forward to the service's target group, while everything else falls
through to lower-priority rules and the listener's default action. It is the
building block of a path-multiplexed ALB, where several services share one
domain and one certificate.

## When to Use

- Several services behind one domain, split by URL prefix
  (`/api/*`, `/auth/*`, `/static/*`)
- Migrating a monolith one path at a time -- each extracted service takes
  over its prefix with a rule like this
- Any deployment where the service owns its route and the listener stays
  shared

## Key Configuration Choices

- **Explicit `priority: 10`** -- path rules shadow each other, so explicit
  numbers matter: a `/api/*` rule must outrank any catch-all `/*` rule.
  Leave gaps (10, 20, 30...) to slot rules in later without renumbering
- **Wildcard, not regex** -- `/api/*` is a wildcard pattern, matched without
  any load balancer attribute changes; use `regexValues` only when wildcards
  cannot express the match
- **Path excludes the query string** -- add a `queryString` condition to
  match parameters
- **Backend still sees `/api/...`** -- add a `url-rewrite` transform
  (`regex: ^/api/(.*)$`, `replace: /$1`) if the service expects to be
  mounted at `/`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<service-name>` | Name prefix for the rule resource | Your service's name |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<listener-resource-name>` | Name of the AwsLbListener to attach to | Your AwsLbListener manifest's `metadata.name` |
| `<target-group-resource-name>` | Name of the AwsLbTargetGroup receiving traffic | Your AwsLbTargetGroup manifest's `metadata.name` |

## Common Additions

- Add a `hostHeader` condition block to scope the path match to one domain
  (blocks AND together)
- Chain an `authenticate-oidc` action before the forward to gate the prefix
  behind SSO

## Related Presets

- **02-host-based-routing** -- split by domain instead of path
- **03-canary-weighted** -- shift this route's traffic gradually between two groups
