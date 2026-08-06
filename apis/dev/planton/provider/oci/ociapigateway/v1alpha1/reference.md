# OciApiGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciApiGatewaySpec defines the specification for an OCI API Gateway bundled
with a single API deployment.

The gateway provides the managed network endpoint (public or private) while
the deployment defines the API specification: routes, backends, authentication,
CORS, and rate limiting. Bundling them reflects the most common pattern where
one gateway serves one API.

Key behaviors:
  - endpoint_type and subnet_id are immutable after creation (ForceNew on gateway)
  - deployment.path_prefix is immutable after creation (ForceNew on deployment)
  - The deployment is always created in the same compartment as the gateway
  - JWT authentication validates tokens before forwarding to backends
  - Routes are evaluated in order; first match wins

Excluded from v1:
  - Gateway: ca_bundles, ip_mode, IPv4/IPv6 address config, response_cache_details (external Redis), locks
  - Authentication: Custom (Functions-based), Token (newer JWT variant), validation_failure_policy (OAuth2 redirect)
  - Deployment policies: dynamic_authentication, mutual_tls, usage_plans
  - Route policies: body_validation, header/query transformations/validations, response_cache
  - Backend types: dynamic_routing, oauth2_logout
  - Tags: defined_tags, system_tags, freeform_tags (auto from labels)

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.endpointType` | `enum` |  |  |  |
| `spec.subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.certificateId` | `string` |  |  |  |
| `spec.networkSecurityGroupIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.deployment` | `Deployment` | yes |  |  |
| `spec.deployment.pathPrefix` | `string` | yes |  |  |
| `spec.deployment.displayName` | `string` |  |  |  |
| `spec.deployment.loggingPolicies` | `LoggingPolicies` |  |  |  |
| `spec.deployment.loggingPolicies.accessLog` | `AccessLog` |  |  |  |
| `spec.deployment.loggingPolicies.accessLog.isEnabled` | `bool` |  |  |  |
| `spec.deployment.loggingPolicies.executionLog` | `ExecutionLog` |  |  |  |
| `spec.deployment.loggingPolicies.executionLog.isEnabled` | `bool` |  |  |  |
| `spec.deployment.loggingPolicies.executionLog.logLevel` | `enum` |  |  |  |
| `spec.deployment.requestPolicies` | `RequestPolicies` |  |  |  |
| `spec.deployment.requestPolicies.authentication` | `Authentication` |  |  |  |
| `spec.deployment.requestPolicies.authentication.issuers` | `[]string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.audiences` | `[]string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.tokenHeader` | `string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.tokenQueryParam` | `string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.tokenAuthScheme` | `string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.maxClockSkewInSeconds` | `float` |  |  |  |
| `spec.deployment.requestPolicies.authentication.isAnonymousAccessAllowed` | `bool` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys` | `PublicKeys` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.type` | `enum` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.uri` | `string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.isSslVerifyDisabled` | `bool` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.maxCacheDurationInHours` | `int32` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.keys` | `[]StaticKey` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.keys[].kid` | `string` | yes |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.keys[].format` | `enum` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.keys[].key` | `string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.keys[].kty` | `string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.keys[].alg` | `string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.keys[].e` | `string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.keys[].n` | `string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.publicKeys.keys[].use` | `string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.verifyClaims` | `[]VerifyClaim` |  |  |  |
| `spec.deployment.requestPolicies.authentication.verifyClaims[].key` | `string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.verifyClaims[].values` | `[]string` |  |  |  |
| `spec.deployment.requestPolicies.authentication.verifyClaims[].isRequired` | `bool` |  |  |  |
| `spec.deployment.requestPolicies.cors` | `CorsPolicy` |  |  |  |
| `spec.deployment.requestPolicies.cors.allowedOrigins` | `[]string` | yes |  |  |
| `spec.deployment.requestPolicies.cors.allowedMethods` | `[]string` |  |  |  |
| `spec.deployment.requestPolicies.cors.allowedHeaders` | `[]string` |  |  |  |
| `spec.deployment.requestPolicies.cors.exposedHeaders` | `[]string` |  |  |  |
| `spec.deployment.requestPolicies.cors.isAllowCredentialsEnabled` | `bool` |  |  |  |
| `spec.deployment.requestPolicies.cors.maxAgeInSeconds` | `int32` |  |  |  |
| `spec.deployment.requestPolicies.rateLimiting` | `RateLimiting` |  |  |  |
| `spec.deployment.requestPolicies.rateLimiting.rateInRequestsPerSecond` | `int32` |  |  |  |
| `spec.deployment.requestPolicies.rateLimiting.rateKey` | `enum` |  |  |  |
| `spec.deployment.routes` | `[]Route` | yes |  |  |
| `spec.deployment.routes[].path` | `string` | yes |  |  |
| `spec.deployment.routes[].methods` | `[]string` |  |  |  |
| `spec.deployment.routes[].backend` | `Backend` | yes |  |  |
| `spec.deployment.routes[].backend.type` | `enum` |  |  |  |
| `spec.deployment.routes[].backend.url` | `string` |  |  |  |
| `spec.deployment.routes[].backend.functionId` | `string` |  |  |  |
| `spec.deployment.routes[].backend.status` | `int32` |  |  |  |
| `spec.deployment.routes[].backend.body` | `string` |  |  |  |
| `spec.deployment.routes[].backend.connectTimeoutInSeconds` | `float` |  |  |  |
| `spec.deployment.routes[].backend.readTimeoutInSeconds` | `float` |  |  |  |
| `spec.deployment.routes[].backend.sendTimeoutInSeconds` | `float` |  |  |  |
| `spec.deployment.routes[].backend.isSslVerifyDisabled` | `bool` |  |  |  |
| `spec.deployment.routes[].backend.headers` | `[]BackendHeader` |  |  |  |
| `spec.deployment.routes[].backend.headers[].name` | `string` | yes |  |  |
| `spec.deployment.routes[].backend.headers[].value` | `string` |  |  |  |
| `spec.deployment.routes[].authorization` | `RouteAuthorization` |  |  |  |
| `spec.deployment.routes[].authorization.type` | `enum` |  |  |  |
| `spec.deployment.routes[].authorization.allowedScope` | `[]string` |  |  |  |
| `spec.deployment.routes[].loggingPolicies` | `LoggingPolicies` |  |  |  |
| `spec.deployment.routes[].loggingPolicies.accessLog` | `AccessLog` |  |  |  |
| `spec.deployment.routes[].loggingPolicies.accessLog.isEnabled` | `bool` |  |  |  |
| `spec.deployment.routes[].loggingPolicies.executionLog` | `ExecutionLog` |  |  |  |
| `spec.deployment.routes[].loggingPolicies.executionLog.isEnabled` | `bool` |  |  |  |
| `spec.deployment.routes[].loggingPolicies.executionLog.logLevel` | `enum` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the gateway and deployment will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.endpointType

`enum`

Whether the gateway is internet-facing (public) or VCN-internal (private).
Immutable after creation.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `endpoint_type_unspecified`
- `endpoint_type_public`
- `endpoint_type_private`

### spec.subnetId

`string | valueFrom` · required

OCID of the subnet where the gateway will be deployed.
For public gateways, this must be a public subnet. Immutable after creation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.displayName

`string`

Display name for the gateway. When omitted, the metadata name is used.

### spec.certificateId

`string`

OCID of the OCI Certificates service certificate for TLS termination
on the gateway. Only meaningful for public gateways serving HTTPS.

### spec.networkSecurityGroupIds

`[]string | valueFrom`

OCIDs of network security groups applied to the gateway.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.deployment

`Deployment` · required

The API deployment that defines routes, backends, and policies.

- rule: {"required":true}

### spec.deployment.pathPrefix

`string` · required

URL path prefix for all routes in this deployment (e.g., "/api/v1").
Must start with "/" and must be unique per gateway. Immutable after creation.

- rule: {"string":{"minLen":"1","prefix":"/"}}

### spec.deployment.displayName

`string`

Display name for the deployment. When omitted, defaults to
"{gateway_display_name}-deployment".

### spec.deployment.loggingPolicies

`LoggingPolicies`

Logging configuration for the deployment.

### spec.deployment.loggingPolicies.accessLog

`AccessLog`

Access log captures every request/response pair.

### spec.deployment.loggingPolicies.accessLog.isEnabled

`bool`

Whether access logging is enabled.

### spec.deployment.loggingPolicies.executionLog

`ExecutionLog`

Execution log captures gateway processing details.

### spec.deployment.loggingPolicies.executionLog.isEnabled

`bool`

Whether execution logging is enabled.

### spec.deployment.loggingPolicies.executionLog.logLevel

`enum`

Log verbosity level. Only meaningful when is_enabled is true.

Allowed values (use exactly as shown):

- `log_level_unspecified`
- `info`
- `warn`
- `error`

### spec.deployment.requestPolicies

`RequestPolicies`

Deployment-level request policies applied before route matching.

### spec.deployment.requestPolicies.authentication

`Authentication`

JWT-based authentication configuration.

- rule: public_keys is required for JWT authentication

### spec.deployment.requestPolicies.authentication.issuers

`[]string`

Token issuers (iss claim). Tokens with issuers not in this list are rejected.

### spec.deployment.requestPolicies.authentication.audiences

`[]string`

Allowed audiences (aud claim). Tokens not targeting any of these audiences are rejected.

### spec.deployment.requestPolicies.authentication.tokenHeader

`string`

HTTP header containing the token. Defaults to "Authorization" on the OCI side.

### spec.deployment.requestPolicies.authentication.tokenQueryParam

`string`

Query parameter containing the token. Mutually exclusive with token_header
in practice, but OCI allows both (header takes precedence).

### spec.deployment.requestPolicies.authentication.tokenAuthScheme

`string`

Authentication scheme prefix in the header value (e.g., "Bearer").
Defaults to "Bearer" on the OCI side.

### spec.deployment.requestPolicies.authentication.maxClockSkewInSeconds

`float` · optional (explicit presence)

Maximum allowed clock skew in seconds when verifying token exp/nbf/iat claims.

### spec.deployment.requestPolicies.authentication.isAnonymousAccessAllowed

`bool` · optional (explicit presence)

When true, unauthenticated requests are allowed through (routes can still
enforce authorization individually). Useful for APIs with mixed public/private routes.

### spec.deployment.requestPolicies.authentication.publicKeys

`PublicKeys`

Public key configuration for verifying token signatures.

- rule: uri is required when type is remote_jwks
- rule: keys must be non-empty when type is static_keys

### spec.deployment.requestPolicies.authentication.publicKeys.type

`enum`

How public keys are obtained.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `public_key_unspecified`
- `remote_jwks`
- `static_keys`

### spec.deployment.requestPolicies.authentication.publicKeys.uri

`string`

URI of the JWKS endpoint (e.g., "https://idp.example.com/.well-known/jwks.json").
Required when type is remote_jwks.

### spec.deployment.requestPolicies.authentication.publicKeys.isSslVerifyDisabled

`bool` · optional (explicit presence)

Whether to skip TLS certificate verification when fetching JWKS.
Only applicable when type is remote_jwks.

### spec.deployment.requestPolicies.authentication.publicKeys.maxCacheDurationInHours

`int32` · optional (explicit presence)

How long (in hours) to cache the JWKS response. Only applicable when
type is remote_jwks. OCI defaults to 1 hour.

### spec.deployment.requestPolicies.authentication.publicKeys.keys

`[]StaticKey`

Static keys for token verification. Required when type is static_keys.

- rule: key is required when format is pem
- rule: kty is required when format is json_web_key

### spec.deployment.requestPolicies.authentication.publicKeys.keys[].kid

`string` · required

Key identifier (kid). Used to match incoming tokens to the correct key.

- rule: {"string":{"minLen":"1"}}

### spec.deployment.requestPolicies.authentication.publicKeys.keys[].format

`enum`

Key encoding format.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `key_format_unspecified`
- `pem`
- `json_web_key`

### spec.deployment.requestPolicies.authentication.publicKeys.keys[].key

`string`

PEM-encoded public key or certificate. Required when format is pem.

### spec.deployment.requestPolicies.authentication.publicKeys.keys[].kty

`string`

JWK key type (e.g., "RSA", "EC"). Required when format is json_web_key.

### spec.deployment.requestPolicies.authentication.publicKeys.keys[].alg

`string`

JWK algorithm (e.g., "RS256"). Used when format is json_web_key.

### spec.deployment.requestPolicies.authentication.publicKeys.keys[].e

`string`

RSA public exponent (Base64url-encoded). Used for RSA keys in JWK format.

### spec.deployment.requestPolicies.authentication.publicKeys.keys[].n

`string`

RSA modulus (Base64url-encoded). Used for RSA keys in JWK format.

### spec.deployment.requestPolicies.authentication.publicKeys.keys[].use

`string`

JWK public key use (e.g., "sig"). Used when format is json_web_key.

### spec.deployment.requestPolicies.authentication.verifyClaims

`[]VerifyClaim`

Additional claims to verify on validated tokens.

### spec.deployment.requestPolicies.authentication.verifyClaims[].key

`string`

Claim name (e.g., "scope", "groups", "email").

### spec.deployment.requestPolicies.authentication.verifyClaims[].values

`[]string`

Expected claim values. If non-empty, the token's claim must match
at least one of these values.

### spec.deployment.requestPolicies.authentication.verifyClaims[].isRequired

`bool` · optional (explicit presence)

Whether the claim must be present in the token.

### spec.deployment.requestPolicies.cors

`CorsPolicy`

Cross-Origin Resource Sharing policy for browser-based clients.

### spec.deployment.requestPolicies.cors.allowedOrigins

`[]string` · required

Allowed origins. At least one origin is required (e.g., ["*"] or
["https://app.example.com"]).

- rule: {"repeated":{"minItems":"1"}}

### spec.deployment.requestPolicies.cors.allowedMethods

`[]string`

Allowed HTTP methods (e.g., ["GET", "POST", "PUT", "DELETE"]).

### spec.deployment.requestPolicies.cors.allowedHeaders

`[]string`

Allowed request headers.

### spec.deployment.requestPolicies.cors.exposedHeaders

`[]string`

Response headers exposed to the browser.

### spec.deployment.requestPolicies.cors.isAllowCredentialsEnabled

`bool` · optional (explicit presence)

Whether the browser may send credentials (cookies, authorization headers).

### spec.deployment.requestPolicies.cors.maxAgeInSeconds

`int32` · optional (explicit presence)

How long (in seconds) the browser may cache preflight responses.

### spec.deployment.requestPolicies.rateLimiting

`RateLimiting`

Rate limiting policy.

### spec.deployment.requestPolicies.rateLimiting.rateInRequestsPerSecond

`int32`

Maximum number of requests per second.

- rule: {"int32":{"gt":0}}

### spec.deployment.requestPolicies.rateLimiting.rateKey

`enum`

How requests are grouped for rate limiting.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `rate_key_unspecified`
- `client_ip`
- `total`

### spec.deployment.routes

`[]Route` · required

API routes defining path-to-backend mappings. At least one route
is required. Routes are evaluated in order; first match wins.

- rule: {"repeated":{"minItems":"1"}}

### spec.deployment.routes[].path

`string` · required

URL path pattern (e.g., "/users/{userId}", "/health").
Supports exact, prefix, and wildcard matching per OCI path semantics.

- rule: {"string":{"minLen":"1"}}

### spec.deployment.routes[].methods

`[]string`

HTTP methods this route handles (e.g., ["GET", "POST"]).
When empty, all methods are accepted.

### spec.deployment.routes[].backend

`Backend` · required

Backend that processes matching requests.

- rule: {"required":true}
- rule: url is required when type is http
- rule: function_id is required when type is oracle_functions

### spec.deployment.routes[].backend.type

`enum`

Backend type.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `backend_unspecified`
- `http`
- `oracle_functions`
- `stock_response`

### spec.deployment.routes[].backend.url

`string`

Target URL for HTTP backends (e.g., "https://backend.example.com:8080").
Required when type is http.

### spec.deployment.routes[].backend.functionId

`string`

OCID of the OCI function to invoke. Required when type is oracle_functions.

### spec.deployment.routes[].backend.status

`int32`

HTTP status code for stock responses (e.g., 200, 404).
Only applicable when type is stock_response.

### spec.deployment.routes[].backend.body

`string`

Response body for stock responses.
Only applicable when type is stock_response.

### spec.deployment.routes[].backend.connectTimeoutInSeconds

`float` · optional (explicit presence)

Connection timeout in seconds. Applicable to http and oracle_functions backends.

### spec.deployment.routes[].backend.readTimeoutInSeconds

`float` · optional (explicit presence)

Read timeout in seconds. Applicable to http and oracle_functions backends.

### spec.deployment.routes[].backend.sendTimeoutInSeconds

`float` · optional (explicit presence)

Send timeout in seconds. Applicable to http and oracle_functions backends.

### spec.deployment.routes[].backend.isSslVerifyDisabled

`bool` · optional (explicit presence)

Whether to skip TLS certificate verification for the backend.
Applicable to http backends.

### spec.deployment.routes[].backend.headers

`[]BackendHeader`

Custom headers to add to requests forwarded to the backend.

### spec.deployment.routes[].backend.headers[].name

`string` · required

Header name.

- rule: {"string":{"minLen":"1"}}

### spec.deployment.routes[].backend.headers[].value

`string`

Header value.

### spec.deployment.routes[].authorization

`RouteAuthorization`

Per-route authorization policy. When authentication is configured at
the deployment level, this controls which routes require valid tokens.

- rule: allowed_scope must be non-empty when type is any_of

### spec.deployment.routes[].authorization.type

`enum`

Authorization enforcement level.

Allowed values (use exactly as shown):

- `authorization_unspecified`
- `anonymous`
- `any_of`
- `authentication_only`

### spec.deployment.routes[].authorization.allowedScope

`[]string`

OAuth2 scopes required when type is any_of. The token must contain
at least one of these scopes.

### spec.deployment.routes[].loggingPolicies

`LoggingPolicies`

Per-route logging policy override.

### spec.deployment.routes[].loggingPolicies.accessLog

`AccessLog`

Access log captures every request/response pair.

### spec.deployment.routes[].loggingPolicies.accessLog.isEnabled

`bool`

Whether access logging is enabled.

### spec.deployment.routes[].loggingPolicies.executionLog

`ExecutionLog`

Execution log captures gateway processing details.

### spec.deployment.routes[].loggingPolicies.executionLog.isEnabled

`bool`

Whether execution logging is enabled.

### spec.deployment.routes[].loggingPolicies.executionLog.logLevel

`enum`

Log verbosity level. Only meaningful when is_enabled is true.

Allowed values (use exactly as shown):

- `log_level_unspecified`
- `info`
- `warn`
- `error`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciApiGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.gateway_id` | `string` | OCID of the API gateway. |
| `status.outputs.hostname` | `string` | Hostname assigned to the gateway by OCI (e.g., "abc123.apigateway.us-ashburn-1.oci.customer-oci.com"). Use this for DNS CNAME records or client configuration. |
| `status.outputs.deployment_endpoint` | `string` | Full endpoint URL of the deployment (gateway hostname + path prefix). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.networkSecurityGroupIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
