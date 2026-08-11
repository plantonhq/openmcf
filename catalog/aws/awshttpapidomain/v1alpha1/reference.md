# AwsHttpApiDomain

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsHttpApiDomainSpec defines the desired configuration for an API Gateway
v2 custom domain name -- the production front door for HTTP APIs. It binds
an owned domain (e.g. "api.example.com") to an ACM certificate and maps one
or more APIs onto the domain through API mappings.

A custom domain is deliberately its own resource rather than a field on the
API: a domain outlives any one API, many APIs can be mapped onto one domain
under different path keys (e.g. "orders" and "billing"), and the domain's
certificate rotates on its own lifecycle.

DNS is composed, not embedded: the domain exports target_domain_name and
hosted_zone_id -- create an alias record with AwsRoute53DnsRecord pointing
your domain at those outputs. The certificate is composed the same way,
referencing an AwsCertManagerCert (or any ACM certificate ARN) that covers
the domain name.

Design notes:
- The domain_name itself is immutable (changing it replaces the resource).
- API Gateway v2 domains only support the REGIONAL endpoint type and the
  TLS_1_2 security policy, so neither is a spec field -- the modules set
  them. (Edge-optimized domains are a REST API v1 feature.)
- The certificate must be issued in the SAME region as the domain and must
  cover the domain name (exact match or wildcard).
- routing_mode (routing rules) is deliberately not modeled: header/path
  rule-based routing is a separate resource family; the default
  API-mapping-only mode covers the mapping surface this component models.

Credentials, region, and deployment workflow live outside this spec in
stack inputs.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsHttpApiDomain
metadata:
  name: test-api-domain
  org: test-org
  env: dev
  id: test-api-domain-dev
spec:
  region: us-west-2
  domainName: api.example.com
  certificateArn:
    value: arn:aws:acm:us-west-2:123456789012:certificate/abc-123
  # AWS-issued public certificate proving domain ownership -- required when
  # the TLS certificate is Private-CA-issued or mTLS uses an imported cert.
  ownershipVerificationCertificateArn:
    value: arn:aws:acm:us-west-2:123456789012:certificate/own-456
  # Rules evaluate in ascending priority; the first match wins. Rules
  # invoke ONLY REST-protocol APIs (ids passed literally), and AWS rejects
  # HTTP-API mappings alongside rule modes (api_mappings and rule-routed
  # domains are mutually exclusive on this kind), so a rule-routed domain
  # carries no apiMappings -- see the 01-single-api-domain preset for the
  # mapping-based shape.
  routingMode: ROUTING_RULE_ONLY
  routingRules:
    - priority: 10
      conditions:
        - basePaths:
            - orders
      apiId:
        value: e5f6a7b8
      stage: prod
      stripBasePath: true
    - priority: 20
      conditions:
        - header:
            name: x-tenant-id
            valueGlob: tenant-a-*
      apiId:
        value: a1b2c3d4
      stage: prod
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.domainName` | `string` | yes |  |  |
| `spec.certificateArn` | `string \| valueFrom` | yes |  | AwsCertManagerCert (`status.outputs.cert_arn`) |
| `spec.ipAddressType` | `string` |  |  |  |
| `spec.mutualTls` | `AwsHttpApiDomainMutualTls` |  |  |  |
| `spec.mutualTls.truststoreUri` | `string` | yes |  |  |
| `spec.mutualTls.truststoreVersion` | `string` |  |  |  |
| `spec.apiMappings` | `[]AwsHttpApiDomainApiMapping` |  |  |  |
| `spec.apiMappings[].apiId` | `string \| valueFrom` | yes |  | AwsHttpApiGateway (`status.outputs.api_id`) |
| `spec.apiMappings[].stage` | `string` | yes |  |  |
| `spec.apiMappings[].apiMappingKey` | `string` |  |  |  |
| `spec.ownershipVerificationCertificateArn` | `string \| valueFrom` |  |  | AwsCertManagerCert (`status.outputs.cert_arn`) |
| `spec.routingMode` | `string` |  |  |  |
| `spec.routingRules` | `[]AwsHttpApiDomainRoutingRule` |  |  |  |
| `spec.routingRules[].priority` | `int32` |  |  |  |
| `spec.routingRules[].conditions` | `[]AwsHttpApiDomainRoutingRuleCondition` | yes |  |  |
| `spec.routingRules[].conditions[].basePaths` | `[]string` |  |  |  |
| `spec.routingRules[].conditions[].header` | `AwsHttpApiDomainRoutingRuleHeaderMatch` |  |  |  |
| `spec.routingRules[].conditions[].header.name` | `string` | yes |  |  |
| `spec.routingRules[].conditions[].header.valueGlob` | `string` | yes |  |  |
| `spec.routingRules[].apiId` | `string \| valueFrom` | yes |  |  |
| `spec.routingRules[].stage` | `string` | yes |  |  |
| `spec.routingRules[].stripBasePath` | `bool` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the domain name will be created. Must match the
region of the certificate and of the APIs mapped onto the domain.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.domainName

`string` · required

The fully qualified custom domain name (e.g. "api.example.com").
Immutable after creation -- changing it replaces the domain. Must be
lowercase. A wildcard domain ("*.example.com") is allowed and matches
all first-level subdomains, but requires a certificate that covers the
wildcard.

- rule: {"string":{"minLen":"1","maxLen":"512","pattern":"^[a-z0-9*.-]+$"}}

### spec.certificateArn

`string | valueFrom` · required

The ACM certificate for TLS termination on this domain. The certificate
must be issued (or imported) in the same region as the domain and must
cover domain_name (exact or wildcard match). Accepts a direct certificate
ARN or a reference to an AwsCertManagerCert resource.

- references: AwsCertManagerCert (`status.outputs.cert_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCertManagerCert, name: <that resource's name>, fieldPath: status.outputs.cert_arn}} -- a bare string does not parse

### spec.ipAddressType

`string`

IP address type for the domain's endpoint.
- "ipv4": Resolve to IPv4 addresses only.
- "dualstack": Resolve to both IPv4 and IPv6.
When omitted, AWS applies its default (dualstack for new domains).

### spec.mutualTls

`AwsHttpApiDomainMutualTls`

Mutual TLS (mTLS) authentication for the domain. When configured, API
Gateway requires clients to present a certificate chaining to a CA in
the truststore before any request reaches an API. Common for B2B and
machine-to-machine APIs. When mTLS is enabled, also set
disable_execute_api_endpoint=true on the mapped APIs -- otherwise
callers can bypass mTLS via the default execute-api endpoint.

### spec.mutualTls.truststoreUri

`string` · required

S3 URI of the truststore -- a PEM bundle of the CA certificates that
client certificates must chain to (e.g.
"s3://my-bucket/truststore.pem"). The bucket must be in the same region
as the domain.

- rule: {"string":{"minLen":"1","pattern":"^s3://.+"}}

### spec.mutualTls.truststoreVersion

`string`

Optional S3 object version of the truststore. Pin a version so
truststore updates are an explicit, auditable change rather than a
silent side effect of overwriting the object.

### spec.apiMappings

`[]AwsHttpApiDomainApiMapping`

APIs mapped onto this domain. Each mapping binds one API's stage under
an optional path key: an empty api_mapping_key serves the API at the
domain root ("https://api.example.com/"), while a key like "orders"
serves it under "https://api.example.com/orders/". Multiple APIs
compose onto one domain by using distinct keys.

### spec.apiMappings[].apiId

`string | valueFrom` · required

The API to map onto the domain. Accepts a direct API ID or a reference
to an AwsHttpApiGateway resource.

- references: AwsHttpApiGateway (`status.outputs.api_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsHttpApiGateway, name: <that resource's name>, fieldPath: status.outputs.api_id}} -- a bare string does not parse

### spec.apiMappings[].stage

`string` · required

The stage of the API to serve. For HTTP APIs managed by
AwsHttpApiGateway this is the stage name exported in the API's outputs
(typically "$default").

- rule: {"string":{"minLen":"1"}}

### spec.apiMappings[].apiMappingKey

`string`

Path key under which the API is served (e.g. "orders" serves the API at
"https://<domain>/orders/..."). Leave empty to serve the API at the
domain root. Must not contain slashes -- nested keys are not supported
by API Gateway v2.

- rule: {"string":{"pattern":"^[^/]*$"}}

### spec.ownershipVerificationCertificateArn

`string | valueFrom`

- references: AwsCertManagerCert (`status.outputs.cert_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCertManagerCert, name: <that resource's name>, fieldPath: status.outputs.cert_arn}} -- a bare string does not parse

### spec.routingMode

`string`

### spec.routingRules

`[]AwsHttpApiDomainRoutingRule`

### spec.routingRules[].priority

`int32`

- rule: {"int32":{"lte":1000000,"gte":1}}

### spec.routingRules[].conditions

`[]AwsHttpApiDomainRoutingRuleCondition` · required

- rule: {"repeated":{"minItems":"1"}}
- rule: each routing rule condition must set exactly one of base_paths or header

### spec.routingRules[].conditions[].basePaths

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.routingRules[].conditions[].header

`AwsHttpApiDomainRoutingRuleHeaderMatch`

### spec.routingRules[].conditions[].header.name

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"40"}}

### spec.routingRules[].conditions[].header.valueGlob

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"128"}}

### spec.routingRules[].apiId

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.routingRules[].stage

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.routingRules[].stripBasePath

`bool`

## Validation Rules

- `ip_address_type_valid`: ip_address_type must be 'ipv4' or 'dualstack' when set
- `api_mapping_keys_unique`: api_mapping_key values must be unique across api_mappings (only one API can serve the domain root, and each path key can carry one API)
- `routing_mode_valid`: routing_mode must be 'API_MAPPING_ONLY', 'ROUTING_RULE_ONLY', or 'ROUTING_RULE_THEN_API_MAPPING' when set
- `routing_rules_require_rule_mode`: routing_rules require routing_mode 'ROUTING_RULE_ONLY' or 'ROUTING_RULE_THEN_API_MAPPING' -- under the default API_MAPPING_ONLY mode the rules would never be evaluated
- `api_mappings_conflict_with_rule_modes`: api_mappings cannot be combined with routing_mode 'ROUTING_RULE_ONLY' or 'ROUTING_RULE_THEN_API_MAPPING' -- AWS rejects HTTP-API mappings on rule-routed domains; route through routing_rules instead
- `rule_mode_requires_routing_rules`: routing_mode 'ROUTING_RULE_ONLY' or 'ROUTING_RULE_THEN_API_MAPPING' requires at least one routing_rules entry
- `routing_rule_priorities_unique`: routing_rules priority values must be unique -- AWS rejects two rules with the same priority on one domain

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsHttpApiDomain, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.domain_name` | `string` | The custom domain name (e.g. "api.example.com"). Exported as an output because it is the domain's join key -- downstream resources resolve references against outputs. |
| `status.outputs.domain_name_arn` | `string` | The ARN of the domain name resource. Useful for IAM policies and tag-based governance. |
| `status.outputs.target_domain_name` | `string` | The API Gateway-managed regional domain name to target from DNS (e.g. "d-abc123.execute-api.us-east-1.amazonaws.com"). Create an alias or CNAME record from the custom domain to this value. |
| `status.outputs.hosted_zone_id` | `string` | The Route 53 hosted zone ID of the API Gateway regional endpoint. Use as the alias target zone when creating a Route 53 alias record for the domain. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.certificateArn` | AwsCertManagerCert | `status.outputs.cert_arn` |
| `spec.apiMappings[].apiId` | AwsHttpApiGateway | `status.outputs.api_id` |
| `spec.ownershipVerificationCertificateArn` | AwsCertManagerCert | `status.outputs.cert_arn` |

## See Also

- [Overview](../README.md)
