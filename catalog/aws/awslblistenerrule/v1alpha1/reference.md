# AwsLbListenerRule

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsLbListenerRuleSpec defines an ALB listener rule: a condition-action pair
that routes matching requests, evaluated in priority order before the
listener's default action.

The rule is the unit of per-service routing. A shared listener (say HTTPS
443 on the environment's ALB) stays untouched while each service deploys
its own rule -- "host api.example.com forwards to the api target group",
"path /admin/* requires OIDC login first" -- and removes it when the
service goes away. That independent lifecycle, many rules per listener, is
why the rule is a first-class kind rather than a listener detail.

Rules are an Application Load Balancer concept: Network Load Balancer
listeners route purely by port/protocol and take no rules.

The listener is create-only: moving a rule replaces it. Priority,
conditions, actions, and transforms all update in place.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsLbListenerRule
metadata:
  name: api-route-demo
spec:
  region: us-west-2
  listenerArn:
    value: arn:aws:elasticloadbalancing:us-west-2:123456789012:listener/app/demo/50dc6c495c0c9188/f2f7dc8efc522ab2
  priority: 10
  # Host AND path conditions plus a url-rewrite transform deliberately
  # exercise the multi-criterion condition shape and the transform block --
  # the parts of the variable contract that would break silently if the
  # nested objects were mistyped.
  conditions:
    - hostHeader:
        values:
          - api.example.com
    - pathPattern:
        values:
          - /v1/*
  actions:
    - type: forward
      forward:
        targetGroups:
          - arn:
              value: arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/api/943f017f100becff
  transforms:
    - type: url-rewrite
      urlRewrite:
        regex: "^/v1/(.*)$"
        replace: "/v2/$1"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.listenerArn` | `string \| valueFrom` | yes |  | AwsLbListener (`status.outputs.listener_arn`) |
| `spec.priority` | `int32` |  |  |  |
| `spec.conditions` | `[]AwsLbListenerRuleCondition` | yes |  |  |
| `spec.conditions[].hostHeader` | `AwsLbListenerRuleHostHeaderCondition` |  |  |  |
| `spec.conditions[].hostHeader.values` | `[]string` |  |  |  |
| `spec.conditions[].hostHeader.regexValues` | `[]string` |  |  |  |
| `spec.conditions[].pathPattern` | `AwsLbListenerRulePathPatternCondition` |  |  |  |
| `spec.conditions[].pathPattern.values` | `[]string` |  |  |  |
| `spec.conditions[].pathPattern.regexValues` | `[]string` |  |  |  |
| `spec.conditions[].httpHeader` | `AwsLbListenerRuleHttpHeaderCondition` |  |  |  |
| `spec.conditions[].httpHeader.httpHeaderName` | `string` | yes |  |  |
| `spec.conditions[].httpHeader.values` | `[]string` |  |  |  |
| `spec.conditions[].httpHeader.regexValues` | `[]string` |  |  |  |
| `spec.conditions[].httpRequestMethod` | `AwsLbListenerRuleHttpRequestMethodCondition` |  |  |  |
| `spec.conditions[].httpRequestMethod.values` | `[]string` | yes |  |  |
| `spec.conditions[].queryString` | `AwsLbListenerRuleQueryStringCondition` |  |  |  |
| `spec.conditions[].queryString.pairs` | `[]AwsLbListenerRuleQueryStringPair` | yes |  |  |
| `spec.conditions[].queryString.pairs[].key` | `string` |  |  |  |
| `spec.conditions[].queryString.pairs[].value` | `string` | yes |  |  |
| `spec.conditions[].sourceIp` | `AwsLbListenerRuleSourceIpCondition` |  |  |  |
| `spec.conditions[].sourceIp.values` | `[]string` | yes |  |  |
| `spec.actions` | `[]AwsLbListenerRuleAction` | yes |  |  |
| `spec.actions[].type` | `string` | yes |  |  |
| `spec.actions[].order` | `int32` |  |  |  |
| `spec.actions[].forward` | `AwsLbListenerRuleActionForward` |  |  |  |
| `spec.actions[].forward.targetGroups` | `[]AwsLbListenerRuleActionForwardTargetGroup` | yes |  |  |
| `spec.actions[].forward.targetGroups[].arn` | `string \| valueFrom` | yes |  | AwsLbTargetGroup (`status.outputs.target_group_arn`) |
| `spec.actions[].forward.targetGroups[].weight` | `int32` |  |  |  |
| `spec.actions[].forward.stickiness` | `AwsLbListenerRuleActionForwardStickiness` |  |  |  |
| `spec.actions[].forward.stickiness.enabled` | `bool` |  |  |  |
| `spec.actions[].forward.stickiness.durationSeconds` | `int32` |  |  |  |
| `spec.actions[].redirect` | `AwsLbListenerRuleActionRedirect` |  |  |  |
| `spec.actions[].redirect.statusCode` | `string` | yes |  |  |
| `spec.actions[].redirect.protocol` | `string` |  |  |  |
| `spec.actions[].redirect.port` | `string` |  |  |  |
| `spec.actions[].redirect.host` | `string` |  |  |  |
| `spec.actions[].redirect.path` | `string` |  |  |  |
| `spec.actions[].redirect.query` | `string` |  |  |  |
| `spec.actions[].fixedResponse` | `AwsLbListenerRuleActionFixedResponse` |  |  |  |
| `spec.actions[].fixedResponse.contentType` | `string` | yes |  |  |
| `spec.actions[].fixedResponse.statusCode` | `string` |  |  |  |
| `spec.actions[].fixedResponse.messageBody` | `string` |  |  |  |
| `spec.actions[].authenticateCognito` | `AwsLbListenerRuleActionAuthenticateCognito` |  |  |  |
| `spec.actions[].authenticateCognito.userPoolArn` | `string \| valueFrom` | yes |  | AwsCognitoUserPool (`status.outputs.user_pool_arn`) |
| `spec.actions[].authenticateCognito.userPoolClientId` | `string \| valueFrom` | yes |  | AwsCognitoUserPoolClient (`status.outputs.client_id`) |
| `spec.actions[].authenticateCognito.userPoolDomain` | `string \| valueFrom` | yes |  | AwsCognitoUserPool (`status.outputs.user_pool_domain`) |
| `spec.actions[].authenticateCognito.authenticationRequestExtraParams` | `map<string, string>` |  |  |  |
| `spec.actions[].authenticateCognito.onUnauthenticatedRequest` | `string` |  |  |  |
| `spec.actions[].authenticateCognito.scope` | `string` |  |  |  |
| `spec.actions[].authenticateCognito.sessionCookieName` | `string` |  |  |  |
| `spec.actions[].authenticateCognito.sessionTimeoutSeconds` | `int32` |  |  |  |
| `spec.actions[].authenticateOidc` | `AwsLbListenerRuleActionAuthenticateOidc` |  |  |  |
| `spec.actions[].authenticateOidc.issuer` | `string` | yes |  |  |
| `spec.actions[].authenticateOidc.authorizationEndpoint` | `string` | yes |  |  |
| `spec.actions[].authenticateOidc.tokenEndpoint` | `string` | yes |  |  |
| `spec.actions[].authenticateOidc.userInfoEndpoint` | `string` | yes |  |  |
| `spec.actions[].authenticateOidc.clientId` | `string` | yes |  |  |
| `spec.actions[].authenticateOidc.clientSecret` | `string` (sensitive) | yes |  |  |
| `spec.actions[].authenticateOidc.authenticationRequestExtraParams` | `map<string, string>` |  |  |  |
| `spec.actions[].authenticateOidc.onUnauthenticatedRequest` | `string` |  |  |  |
| `spec.actions[].authenticateOidc.scope` | `string` |  |  |  |
| `spec.actions[].authenticateOidc.sessionCookieName` | `string` |  |  |  |
| `spec.actions[].authenticateOidc.sessionTimeoutSeconds` | `int32` |  |  |  |
| `spec.actions[].jwtValidation` | `AwsLbListenerRuleActionJwtValidation` |  |  |  |
| `spec.actions[].jwtValidation.issuer` | `string` | yes |  |  |
| `spec.actions[].jwtValidation.jwksEndpoint` | `string` | yes |  |  |
| `spec.actions[].jwtValidation.additionalClaims` | `[]AwsLbListenerRuleActionJwtClaim` |  |  |  |
| `spec.actions[].jwtValidation.additionalClaims[].name` | `string` | yes |  |  |
| `spec.actions[].jwtValidation.additionalClaims[].format` | `string` | yes |  |  |
| `spec.actions[].jwtValidation.additionalClaims[].values` | `[]string` | yes |  |  |
| `spec.transforms` | `[]AwsLbListenerRuleTransform` |  |  |  |
| `spec.transforms[].type` | `string` | yes |  |  |
| `spec.transforms[].hostHeaderRewrite` | `AwsLbListenerRuleRewrite` |  |  |  |
| `spec.transforms[].hostHeaderRewrite.regex` | `string` | yes |  |  |
| `spec.transforms[].hostHeaderRewrite.replace` | `string` |  |  |  |
| `spec.transforms[].urlRewrite` | `AwsLbListenerRuleRewrite` |  |  |  |
| `spec.transforms[].urlRewrite.regex` | `string` | yes |  |  |
| `spec.transforms[].urlRewrite.replace` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the rule is created. Must match the listener's
region. Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.listenerArn

`string | valueFrom` · required

The listener this rule attaches to. Immutable: changing the listener
replaces the rule.

- references: AwsLbListener (`status.outputs.listener_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLbListener, name: <that resource's name>, fieldPath: status.outputs.listener_arn}} -- a bare string does not parse

### spec.priority

`int32`

Evaluation priority among the listener's rules, 1-50000; lower values
are evaluated first and must be unique per listener. When omitted, AWS
assigns the next free priority after the current highest -- fine for
append-only routing, but set explicit priorities when rules shadow each
other (a "/api/*" rule must outrank a catch-all "/*" rule).

### spec.conditions

`[]AwsLbListenerRuleCondition` · required

What a request must match for the rule to fire. Required, 1-5 condition
blocks; a request must satisfy ALL of them (conditions AND together;
the values inside one condition OR together).

- rule: {"required":true,"repeated":{"minItems":"1","maxItems":"5"}}
- rule: exactly one of host_header, path_pattern, http_header, http_request_method, query_string, or source_ip must be set

### spec.conditions[].hostHeader

`AwsLbListenerRuleHostHeaderCondition`

Match on the Host header, e.g. "api.example.com" or "*.example.com"
(wildcards: "*" spans any characters, "?" one character).

- rule: at least one of values or regex_values must be non-empty

### spec.conditions[].hostHeader.values

`[]string`

Wildcard hostname patterns, each up to 128 characters; the condition
matches when ANY pattern matches.

- rule: {"repeated":{"items":{"string":{"maxLen":"128"}}}}

### spec.conditions[].hostHeader.regexValues

`[]string`

Regular-expression hostname patterns, each up to 128 characters. Regex
matching must be enabled on the load balancer's attributes; prefer
wildcards when they express the intent.

- rule: {"repeated":{"items":{"string":{"maxLen":"128"}}}}

### spec.conditions[].pathPattern

`AwsLbListenerRulePathPatternCondition`

Match on the request path, e.g. "/admin/*".

- rule: at least one of values or regex_values must be non-empty

### spec.conditions[].pathPattern.values

`[]string`

Wildcard path patterns, each up to 128 characters (e.g. "/api/*"); the
condition matches when ANY pattern matches.

- rule: {"repeated":{"items":{"string":{"maxLen":"128"}}}}

### spec.conditions[].pathPattern.regexValues

`[]string`

Regular-expression path patterns, each up to 128 characters. Regex
matching must be enabled on the load balancer's attributes; prefer
wildcards when they express the intent.

- rule: {"repeated":{"items":{"string":{"maxLen":"128"}}}}

### spec.conditions[].httpHeader

`AwsLbListenerRuleHttpHeaderCondition`

Match on an arbitrary request header.

- rule: at least one of values or regex_values must be non-empty

### spec.conditions[].httpHeader.httpHeaderName

`string` · required

The header name to inspect, 1-40 characters (RFC 7230 token characters).

- rule: {"required":true,"string":{"maxLen":"40"}}

### spec.conditions[].httpHeader.values

`[]string`

Wildcard patterns for the header value, each up to 128 characters; the
condition matches when ANY pattern matches.

- rule: {"repeated":{"items":{"string":{"maxLen":"128"}}}}

### spec.conditions[].httpHeader.regexValues

`[]string`

Regular-expression patterns for the header value, each up to 128
characters. Regex matching must be enabled on the load balancer's
attributes.

- rule: {"repeated":{"items":{"string":{"maxLen":"128"}}}}

### spec.conditions[].httpRequestMethod

`AwsLbListenerRuleHttpRequestMethodCondition`

Match on the HTTP method, e.g. "GET", "POST".

### spec.conditions[].httpRequestMethod.values

`[]string` · required

Methods to match, e.g. "GET", "POST", or a custom method name (up to 40
characters, letters/hyphens/underscores). Matching is case-sensitive and
exact, per AWS.

- rule: {"required":true,"repeated":{"minItems":"1","items":{"string":{"pattern":"^[A-Za-z-_]{1,40}$"}}}}

### spec.conditions[].queryString

`AwsLbListenerRuleQueryStringCondition`

Match on query-string key/value pairs.

### spec.conditions[].queryString.pairs

`[]AwsLbListenerRuleQueryStringPair` · required

The pairs to match.

- rule: {"required":true,"repeated":{"minItems":"1"}}

### spec.conditions[].queryString.pairs[].key

`string`

The parameter name (wildcards allowed). When omitted, the value is
matched against every parameter.

### spec.conditions[].queryString.pairs[].value

`string` · required

The value pattern (wildcards allowed).

- rule: {"required":true}

### spec.conditions[].sourceIp

`AwsLbListenerRuleSourceIpCondition`

Match on the client source IP (CIDR blocks). Uses the address that
connected to the ALB, not X-Forwarded-For entries.

### spec.conditions[].sourceIp.values

`[]string` · required

CIDR blocks to match (IPv4 or IPv6), e.g. "10.0.0.0/8",
"203.0.113.0/24". The condition matches when the client address falls in
ANY block.

- rule: {"required":true,"repeated":{"minItems":"1"}}

### spec.actions

`[]AwsLbListenerRuleAction` · required

What happens to a matching request. Required, at least one action. With
several actions (an authenticate action chained before its forward),
order follows the list; every chain ends in exactly one of forward,
redirect, or fixed-response.

- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: type must be one of: forward, redirect, fixed-response, authenticate-cognito, authenticate-oidc, jwt-validation
- rule: order must be between 1 and 50000 when set
- rule: forward configuration must be set when (and only when) type is 'forward'
- rule: redirect configuration must be set when (and only when) type is 'redirect'
- rule: fixed_response configuration must be set when (and only when) type is 'fixed-response'
- rule: authenticate_cognito configuration must be set when (and only when) type is 'authenticate-cognito'
- rule: authenticate_oidc configuration must be set when (and only when) type is 'authenticate-oidc'
- rule: jwt_validation configuration must be set when (and only when) type is 'jwt-validation'

### spec.actions[].type

`string` · required

The action type. Valid values: "forward", "redirect", "fixed-response",
"authenticate-cognito", "authenticate-oidc", "jwt-validation".

- rule: {"required":true}

### spec.actions[].order

`int32`

Explicit evaluation order within the chain, 1-50000 (lower runs first).
When omitted, list position decides -- which is almost always enough.

### spec.actions[].forward

`AwsLbListenerRuleActionForward`

Forward the request to one or more weighted target groups.

### spec.actions[].forward.targetGroups

`[]AwsLbListenerRuleActionForwardTargetGroup` · required

The destination target groups, 1-5, each with an optional weight.

- rule: {"required":true,"repeated":{"minItems":"1","maxItems":"5"}}
- rule: weight must be between 0 and 999

### spec.actions[].forward.targetGroups[].arn

`string | valueFrom` · required

The target group receiving the traffic.

- references: AwsLbTargetGroup (`status.outputs.target_group_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLbTargetGroup, name: <that resource's name>, fieldPath: status.outputs.target_group_arn}} -- a bare string does not parse

### spec.actions[].forward.targetGroups[].weight

`int32`

Relative traffic share, 0-999 (proportional to the sum across groups;
0 drains the group). AWS default: 1, so equal weights need no
configuration.

### spec.actions[].forward.stickiness

`AwsLbListenerRuleActionForwardStickiness`

Target-group stickiness for multi-group forwards: pins a client to the
group (not the individual target) that served its first request, so a
session does not flap between blue and green mid-canary.

- rule: duration_seconds is required when stickiness is enabled
- rule: duration_seconds must be between 1 and 604800 when set

### spec.actions[].forward.stickiness.enabled

`bool`

Whether group-level stickiness is active.

### spec.actions[].forward.stickiness.durationSeconds

`int32`

Seconds the group association lasts, 1-604800 (7 days). Required when
enabled.

### spec.actions[].redirect

`AwsLbListenerRuleActionRedirect`

Reply with an HTTP redirect.

- rule: status_code must be 'HTTP_301' or 'HTTP_302'
- rule: protocol must be 'HTTP', 'HTTPS', or '#{protocol}' when set

### spec.actions[].redirect.statusCode

`string` · required

The redirect status code. Required. Valid values: "HTTP_301" (permanent
-- browsers cache it) or "HTTP_302" (temporary).

- rule: {"required":true}

### spec.actions[].redirect.protocol

`string`

Target protocol: "HTTPS", "HTTP", or "#{protocol}" (keep the original).

### spec.actions[].redirect.port

`string`

Target port: a port number as a string, or "#{port}" (keep the
original).

### spec.actions[].redirect.host

`string`

Target hostname, 1-128 characters; supports placeholders (e.g.
"#{host}").

### spec.actions[].redirect.path

`string`

Target absolute path, starting with "/"; supports placeholders (e.g.
"/#{path}" or "/new/#{path}").

### spec.actions[].redirect.query

`string`

Target query string, without the leading "?"; supports placeholders.

### spec.actions[].fixedResponse

`AwsLbListenerRuleActionFixedResponse`

Reply with a canned HTTP response, never touching a target.

- rule: content_type must be one of: text/plain, text/css, text/html, application/javascript, application/json
- rule: status_code must be a 2xx, 4xx, or 5xx code when set

### spec.actions[].fixedResponse.contentType

`string` · required

The Content-Type of the response. Required. Valid values: "text/plain",
"text/css", "text/html", "application/javascript", "application/json".

- rule: {"required":true}

### spec.actions[].fixedResponse.statusCode

`string`

The HTTP status code, as a string: 2xx, 4xx, or 5xx (e.g. "404").
AWS default: "503".

### spec.actions[].fixedResponse.messageBody

`string`

The response body, up to 1024 characters.

- rule: {"string":{"maxLen":"1024"}}

### spec.actions[].authenticateCognito

`AwsLbListenerRuleActionAuthenticateCognito`

Authenticate the client against an Amazon Cognito user pool before the
terminal action runs. HTTPS listeners only.

- rule: on_unauthenticated_request must be 'deny', 'allow', or 'authenticate' when set

### spec.actions[].authenticateCognito.userPoolArn

`string | valueFrom` · required

The Cognito user pool performing authentication.

- references: AwsCognitoUserPool (`status.outputs.user_pool_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCognitoUserPool, name: <that resource's name>, fieldPath: status.outputs.user_pool_arn}} -- a bare string does not parse

### spec.actions[].authenticateCognito.userPoolClientId

`string | valueFrom` · required

The app client ID within the user pool. Accepts a direct client ID or a
reference to an AwsCognitoUserPoolClient resource.

- references: AwsCognitoUserPoolClient (`status.outputs.client_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCognitoUserPoolClient, name: <that resource's name>, fieldPath: status.outputs.client_id}} -- a bare string does not parse

### spec.actions[].authenticateCognito.userPoolDomain

`string | valueFrom` · required

The user pool's hosted-UI domain: the prefix of a Cognito domain (e.g.
"my-app" for my-app.auth.us-west-2.amazoncognito.com) or a full custom
domain. Accepts a direct domain string or a reference to the pool's
user_pool_domain output.

- references: AwsCognitoUserPool (`status.outputs.user_pool_domain`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCognitoUserPool, name: <that resource's name>, fieldPath: status.outputs.user_pool_domain}} -- a bare string does not parse

### spec.actions[].authenticateCognito.authenticationRequestExtraParams

`map<string, string>`

Extra query parameters appended to the authorization-endpoint request.

### spec.actions[].authenticateCognito.onUnauthenticatedRequest

`string`

What happens to an unauthenticated request: "authenticate" (AWS default
-- redirect to login), "deny" (401), or "allow" (pass through; the app
sees the request without identity headers).

### spec.actions[].authenticateCognito.scope

`string`

Space-separated OAuth scopes to request. AWS default: "openid".

### spec.actions[].authenticateCognito.sessionCookieName

`string`

The session cookie name. AWS default: "AWSELBAuthSessionCookie".

### spec.actions[].authenticateCognito.sessionTimeoutSeconds

`int32`

Seconds the authentication session stays valid. AWS default: 604800
(7 days).

### spec.actions[].authenticateOidc

`AwsLbListenerRuleActionAuthenticateOidc`

Authenticate the client against any OpenID Connect provider before the
terminal action runs. HTTPS listeners only.

- rule: on_unauthenticated_request must be 'deny', 'allow', or 'authenticate' when set

### spec.actions[].authenticateOidc.issuer

`string` · required

The OIDC issuer identifier (e.g. "https://accounts.google.com").

- rule: {"required":true}

### spec.actions[].authenticateOidc.authorizationEndpoint

`string` · required

The provider's authorization endpoint URL.

- rule: {"required":true}

### spec.actions[].authenticateOidc.tokenEndpoint

`string` · required

The provider's token endpoint URL.

- rule: {"required":true}

### spec.actions[].authenticateOidc.userInfoEndpoint

`string` · required

The provider's UserInfo endpoint URL.

- rule: {"required":true}

### spec.actions[].authenticateOidc.clientId

`string` · required

The OAuth client ID registered with the provider.

- rule: {"required":true}

### spec.actions[].authenticateOidc.clientSecret

`string` · required · sensitive

The OAuth client secret. Handled as a secret end to end -- never exposed
in plans, logs, or the console.

- rule: {"required":true}

### spec.actions[].authenticateOidc.authenticationRequestExtraParams

`map<string, string>`

Extra query parameters appended to the authorization-endpoint request.

### spec.actions[].authenticateOidc.onUnauthenticatedRequest

`string`

What happens to an unauthenticated request: "authenticate" (AWS default
-- redirect to login), "deny" (401), or "allow" (pass through).

### spec.actions[].authenticateOidc.scope

`string`

Space-separated OAuth scopes to request. AWS default: "openid".

### spec.actions[].authenticateOidc.sessionCookieName

`string`

The session cookie name. AWS default: "AWSELBAuthSessionCookie".

### spec.actions[].authenticateOidc.sessionTimeoutSeconds

`int32`

Seconds the authentication session stays valid. AWS default: 604800
(7 days).

### spec.actions[].jwtValidation

`AwsLbListenerRuleActionJwtValidation`

Validate a JWT bearer token on the request before the terminal action
runs -- API-style auth with no session cookies or login redirects.
HTTPS listeners only.

### spec.actions[].jwtValidation.issuer

`string` · required

The expected token issuer (the "iss" claim), 1-256 characters.

- rule: {"required":true,"string":{"maxLen":"256"}}

### spec.actions[].jwtValidation.jwksEndpoint

`string` · required

The JWKS endpoint URL serving the signing keys tokens are verified
against, 1-256 characters.

- rule: {"required":true,"string":{"maxLen":"256"}}

### spec.actions[].jwtValidation.additionalClaims

`[]AwsLbListenerRuleActionJwtClaim`

Extra claims the token must carry, beyond issuer and signature. Up to 10.

- rule: {"repeated":{"maxItems":"10"}}
- rule: format must be 'single-string', 'string-array', or 'space-separated-values'

### spec.actions[].jwtValidation.additionalClaims[].name

`string` · required

The claim name in the token payload (e.g. "aud", "scope").

- rule: {"required":true}

### spec.actions[].jwtValidation.additionalClaims[].format

`string` · required

How the claim value is encoded in the token. Valid values:
"single-string", "string-array", "space-separated-values" (how OAuth
"scope" claims are typically encoded).

- rule: {"required":true}

### spec.actions[].jwtValidation.additionalClaims[].values

`[]string` · required

Accepted values for the claim, 1-10 entries of up to 256 characters.
The claim satisfies validation when any listed value matches.

- rule: {"required":true,"repeated":{"minItems":"1","maxItems":"10","items":{"string":{"maxLen":"256"}}}}

### spec.transforms

`[]AwsLbListenerRuleTransform`

Request rewrites applied before the action runs: rewrite the host
header, the URL, or both (at most one of each).

- rule: {"repeated":{"maxItems":"2"}}
- rule: type must be 'host-header-rewrite' or 'url-rewrite'
- rule: host_header_rewrite must be set when (and only when) type is 'host-header-rewrite'
- rule: url_rewrite must be set when (and only when) type is 'url-rewrite'

### spec.transforms[].type

`string` · required

The transform type. Required. Valid values: "host-header-rewrite",
"url-rewrite".

- rule: {"required":true}

### spec.transforms[].hostHeaderRewrite

`AwsLbListenerRuleRewrite`

Rewrite configuration for "host-header-rewrite".

### spec.transforms[].hostHeaderRewrite.regex

`string` · required

The regular expression matched against the component being rewritten,
1-1024 characters.

- rule: {"required":true,"string":{"maxLen":"1024"}}

### spec.transforms[].hostHeaderRewrite.replace

`string`

The replacement, up to 1024 characters; may reference capture groups
from the regex.

- rule: {"string":{"maxLen":"1024"}}

### spec.transforms[].urlRewrite

`AwsLbListenerRuleRewrite`

Rewrite configuration for "url-rewrite".

### spec.transforms[].urlRewrite.regex

`string` · required

The regular expression matched against the component being rewritten,
1-1024 characters.

- rule: {"required":true,"string":{"maxLen":"1024"}}

### spec.transforms[].urlRewrite.replace

`string`

The replacement, up to 1024 characters; may reference capture groups
from the regex.

- rule: {"string":{"maxLen":"1024"}}

## Validation Rules

- `priority_range`: priority must be between 1 and 50000 when set

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsLbListenerRule, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.rule_arn` | `string` | The ARN of the rule (e.g. "arn:aws:elasticloadbalancing:us-west-2: 123456789012:listener-rule/app/api/50dc.../f2f7.../9683b2d02a6cabee"). The handle audit tooling and imports reference. |
| `status.outputs.priority` | `string` | The priority AWS assigned to the rule -- meaningful when the spec left priority unset and AWS picked the next free slot. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.listenerArn` | AwsLbListener | `status.outputs.listener_arn` |
| `spec.actions[].forward.targetGroups[].arn` | AwsLbTargetGroup | `status.outputs.target_group_arn` |
| `spec.actions[].authenticateCognito.userPoolArn` | AwsCognitoUserPool | `status.outputs.user_pool_arn` |
| `spec.actions[].authenticateCognito.userPoolClientId` | AwsCognitoUserPoolClient | `status.outputs.client_id` |
| `spec.actions[].authenticateCognito.userPoolDomain` | AwsCognitoUserPool | `status.outputs.user_pool_domain` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsEcsService | `spec.loadBalancers[].advancedConfiguration.productionListenerRule` | `status.outputs.rule_arn` |
| AwsEcsService | `spec.loadBalancers[].advancedConfiguration.testListenerRule` | `status.outputs.rule_arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
