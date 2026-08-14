# AwsRestApiDomain

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsRestApiDomain example (hack/dev manifest and refgen
# Example source): a REGIONAL custom domain with ACM certificate, one
# root mapping, and (commented in spirit) the PRIVATE access-association
# arm shown as a second mapping. Literal ARNs stand in for composed
# references so the offline `tofu plan` renders every coexisting arm.
# Mutual TLS and uploaded certificates are XOR with ACM and live in
# spec tests; PRIVATE access associations require endpoint type PRIVATE.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRestApiDomain
metadata:
  name: api-example-com
  id: api-example-com
  org: test-org
  env: dev
spec:
  region: us-west-2
  domainName: api.example.com
  certificateArn:
    value: arn:aws:acm:us-west-2:123456789012:certificate/abcd-1234
  endpointConfiguration:
    type: REGIONAL
    ipAddressType: ipv4
  # endpoint_access_mode pairs only with the SecurityPolicy_* family
  # (AWS rejects it with legacy policy names -- the CEL enforces it).
  endpointAccessMode: BASIC
  securityPolicy: SecurityPolicy_TLS13_1_2_2021_06
  routingMode: BASE_PATH_MAPPING_ONLY
  basePathMappings:
    - basePath: orders
      restApiId:
        value: abcdef1234
      stageName:
        value: prod
    - restApiId:
        value: abcdef1234
      stageName:
        value: prod
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.domainName` | `string` | yes |  |  |
| `spec.certificateArn` | `string \| valueFrom` |  |  | AwsCertManagerCert (`status.outputs.cert_arn`) |
| `spec.uploadedCertificate` | `AwsRestApiDomainUploadedCertificate` |  |  |  |
| `spec.uploadedCertificate.name` | `string` | yes |  |  |
| `spec.uploadedCertificate.body` | `string` | yes |  |  |
| `spec.uploadedCertificate.chain` | `string` |  |  |  |
| `spec.uploadedCertificate.privateKey` | `string` (sensitive) | yes |  |  |
| `spec.endpointConfiguration` | `AwsRestApiDomainEndpointConfiguration` |  |  |  |
| `spec.endpointConfiguration.type` | `string` |  |  |  |
| `spec.endpointConfiguration.ipAddressType` | `string` |  |  |  |
| `spec.endpointAccessMode` | `string` |  |  |  |
| `spec.securityPolicy` | `string` |  |  |  |
| `spec.mutualTls` | `AwsRestApiDomainMutualTls` |  |  |  |
| `spec.mutualTls.truststoreUri` | `string` | yes |  |  |
| `spec.mutualTls.truststoreVersion` | `string` |  |  |  |
| `spec.ownershipVerificationCertificateArn` | `string \| valueFrom` |  |  | AwsCertManagerCert (`status.outputs.cert_arn`) |
| `spec.policy` | `object` |  |  |  |
| `spec.routingMode` | `string` |  |  |  |
| `spec.basePathMappings` | `[]AwsRestApiDomainBasePathMapping` |  |  |  |
| `spec.basePathMappings[].basePath` | `string` |  |  |  |
| `spec.basePathMappings[].restApiId` | `string \| valueFrom` | yes |  | AwsRestApiGateway (`status.outputs.rest_api_id`) |
| `spec.basePathMappings[].stageName` | `string \| valueFrom` |  |  | AwsRestApiGateway (`status.outputs.stage_name`) |
| `spec.accessAssociations` | `[]AwsRestApiDomainAccessAssociation` |  |  |  |
| `spec.accessAssociations[].vpcEndpointId` | `string \| valueFrom` | yes |  | AwsVpcEndpoint (`status.outputs.vpc_endpoint_id`) |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.domainName

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"253"}}

### spec.certificateArn

`string | valueFrom`

- references: AwsCertManagerCert (`status.outputs.cert_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCertManagerCert, name: <that resource's name>, fieldPath: status.outputs.cert_arn}} -- a bare string does not parse

### spec.uploadedCertificate

`AwsRestApiDomainUploadedCertificate`

### spec.uploadedCertificate.name

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.uploadedCertificate.body

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.uploadedCertificate.chain

`string`

### spec.uploadedCertificate.privateKey

`string` · required · sensitive

- rule: {"required":true}

### spec.endpointConfiguration

`AwsRestApiDomainEndpointConfiguration`

- rule: PRIVATE endpoints require ip_address_type 'dualstack' (or omit it for the AWS default)

### spec.endpointConfiguration.type

`string`

- rule: {"string":{"in":["REGIONAL","EDGE","PRIVATE"]}}

### spec.endpointConfiguration.ipAddressType

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ipv4","dualstack"]}}

### spec.endpointAccessMode

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["BASIC","STRICT"]}}

### spec.securityPolicy

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["TLS_1_0","TLS_1_2","SecurityPolicy_TLS13_1_3_2025_09","SecurityPolicy_TLS13_1_3_FIPS_2025_09","SecurityPolicy_TLS13_1_2_PFS_PQ_2025_09","SecurityPolicy_TLS13_1_2_FIPS_PQ_2025_09","SecurityPolicy_TLS13_1_2_FIPS_PFS_PQ_2025_09","SecurityPolicy_TLS13_1_2_PQ_2025_09","SecurityPolicy_TLS13_1_2_2021_06","SecurityPolicy_TLS13_2025_EDGE","SecurityPolicy_TLS12_PFS_2025_EDGE","SecurityPolicy_TLS12_2018_EDGE"]}}

### spec.mutualTls

`AwsRestApiDomainMutualTls`

### spec.mutualTls.truststoreUri

`string` · required

- rule: {"string":{"minLen":"1","pattern":"^s3://"}}

### spec.mutualTls.truststoreVersion

`string`

### spec.ownershipVerificationCertificateArn

`string | valueFrom`

- references: AwsCertManagerCert (`status.outputs.cert_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCertManagerCert, name: <that resource's name>, fieldPath: status.outputs.cert_arn}} -- a bare string does not parse

### spec.policy

`object`

### spec.routingMode

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["BASE_PATH_MAPPING_ONLY","ROUTING_RULE_ONLY","ROUTING_RULE_THEN_BASE_PATH_MAPPING"]}}

### spec.basePathMappings

`[]AwsRestApiDomainBasePathMapping`

### spec.basePathMappings[].basePath

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"128","pattern":"^[^/]+$"}}

### spec.basePathMappings[].restApiId

`string | valueFrom` · required

- references: AwsRestApiGateway (`status.outputs.rest_api_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsRestApiGateway, name: <that resource's name>, fieldPath: status.outputs.rest_api_id}} -- a bare string does not parse

### spec.basePathMappings[].stageName

`string | valueFrom`

- references: AwsRestApiGateway (`status.outputs.stage_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsRestApiGateway, name: <that resource's name>, fieldPath: status.outputs.stage_name}} -- a bare string does not parse

### spec.accessAssociations

`[]AwsRestApiDomainAccessAssociation`

### spec.accessAssociations[].vpcEndpointId

`string | valueFrom` · required

- references: AwsVpcEndpoint (`status.outputs.vpc_endpoint_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsVpcEndpoint, name: <that resource's name>, fieldPath: status.outputs.vpc_endpoint_id}} -- a bare string does not parse

## Validation Rules

- `certificate_exactly_one_source`: set exactly one of certificate_arn (ACM) or uploaded_certificate
- `base_paths_unique`: base_path_mappings entries must have unique base_path values
- `access_associations_private_only`: access_associations apply only to PRIVATE endpoint types
- `access_association_endpoints_unique`: access_associations entries must reference unique VPC endpoints
- `policy_private_only`: policy applies only to PRIVATE endpoint types
- `access_mode_requires_modern_security_policy`: endpoint_access_mode requires a security_policy from the SecurityPolicy_* family

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsRestApiDomain, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.domain_name` | `string` |  |
| `status.outputs.domain_name_arn` | `string` |  |
| `status.outputs.domain_name_id` | `string` |  |
| `status.outputs.regional_domain_name` | `string` |  |
| `status.outputs.regional_zone_id` | `string` |  |
| `status.outputs.cloudfront_domain_name` | `string` |  |
| `status.outputs.cloudfront_zone_id` | `string` |  |
| `status.outputs.base_path_mapping_ids` | `map<string, string>` |  |
| `status.outputs.access_association_arns` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.certificateArn` | AwsCertManagerCert | `status.outputs.cert_arn` |
| `spec.ownershipVerificationCertificateArn` | AwsCertManagerCert | `status.outputs.cert_arn` |
| `spec.basePathMappings[].restApiId` | AwsRestApiGateway | `status.outputs.rest_api_id` |
| `spec.basePathMappings[].stageName` | AwsRestApiGateway | `status.outputs.stage_name` |
| `spec.accessAssociations[].vpcEndpointId` | AwsVpcEndpoint | `status.outputs.vpc_endpoint_id` |

## See Also

- [Overview](../README.md)
