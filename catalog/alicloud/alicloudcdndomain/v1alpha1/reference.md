# AliCloudCdnDomain

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudCdnDomainSpec defines the configuration for an Alibaba Cloud CDN
accelerated domain managed by the CDN service.

A CDN domain maps a user-facing domain name (e.g., "cdn.example.com") to one
or more origin servers. Alibaba Cloud's globally distributed edge nodes cache
and serve content from the origins, reducing latency for end users. After
creating the CDN domain, point a CNAME record at your DNS provider to the
`cname` value returned in the stack outputs.

This component creates a single CDN domain with origin sources and optional
HTTPS certificate configuration. Advanced CDN function configs (caching rules,
referer whitelists, HTTP headers) are outside the scope of this component and
should be managed via raw Terraform or Pulumi if needed.

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudCdnDomain
metadata:
  name: alicloudcdndomain-demo
spec:
  region: cn-hangzhou
  domainName: cdn.example.com
  cdnType: web
  sources:
    - type: ipaddr
      content: "203.0.113.10"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.domainName` | `string` | yes |  |  |
| `spec.cdnType` | `string` | yes |  |  |
| `spec.scope` | `string` |  |  |  |
| `spec.sources` | `[]AliCloudCdnDomainSource` | yes |  |  |
| `spec.sources[].type` | `string` | yes |  |  |
| `spec.sources[].content` | `string` | yes |  |  |
| `spec.sources[].port` | `int32` |  |  |  |
| `spec.sources[].priority` | `int32` |  |  |  |
| `spec.sources[].weight` | `int32` |  |  |  |
| `spec.certificateConfig` | `AliCloudCdnDomainCertificateConfig` |  |  |  |
| `spec.certificateConfig.certName` | `string` |  |  |  |
| `spec.certificateConfig.certType` | `string` |  |  |  |
| `spec.certificateConfig.certId` | `string` |  |  |  |
| `spec.certificateConfig.certRegion` | `string` |  |  |  |
| `spec.certificateConfig.serverCertificate` | `string` |  |  |  |
| `spec.certificateConfig.privateKey` | `string` (sensitive) |  |  |  |
| `spec.certificateConfig.serverCertificateStatus` | `string` |  |  |  |
| `spec.checkUrl` | `string` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region for provider initialization.
CDN is a global service, but the provider requires a region for API calls.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.domainName

`string` · required

The accelerated domain name (e.g., "cdn.example.com").
Must be a valid domain name between 1 and 63 characters.
Cannot be changed after creation (ForceNew).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63"}}

### spec.cdnType

`string` · required

The type of content the CDN accelerates.
Cannot be changed after creation (ForceNew).
  - "web"      : images, small files, and web pages
  - "download" : large file downloads (software, games)
  - "video"    : on-demand and live video streaming

- rule: cdn_type must be one of: web, download, video
- rule: {"required":true}

### spec.scope

`string`

Geographic scope for CDN acceleration.
  - "domestic" : mainland China only
  - "overseas" : outside mainland China
  - "global"   : worldwide
If omitted, the Alibaba Cloud API defaults to "domestic".

- rule: scope must be one of: domestic, overseas, global

### spec.sources

`[]AliCloudCdnDomainSource` · required

Origin server sources for the CDN domain. At least one source is required.
Multiple sources enable failover and load distribution via priority and weight.

- rule: {"repeated":{"minItems":"1"}}

### spec.sources[].type

`string` · required

The type of origin source.
  - "ipaddr"  : origin is an IP address
  - "domain"  : origin is a domain name
  - "oss"     : origin is an Alibaba Cloud OSS bucket
  - "common"  : origin is a common source (L2 CDN)

- rule: source type must be one of: ipaddr, domain, oss, common
- rule: {"required":true}

### spec.sources[].content

`string` · required

The address of the origin server.
For "ipaddr": an IP address (e.g., "1.2.3.4").
For "domain": a domain name (e.g., "origin.example.com").
For "oss": the OSS bucket domain (e.g., "my-bucket.oss-cn-hangzhou.aliyuncs.com").
Each source's content must be unique within the sources list.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.sources[].port

`int32`

The port number for origin requests.
Typically 80 for HTTP or 443 for HTTPS.
If omitted, defaults to 80.

### spec.sources[].priority

`int32`

Priority of this origin source. Lower values indicate higher priority.
Valid range: 0-100. If omitted, defaults to 20.
Use 20 for primary sources and 30 for standby sources.

- rule: priority must be between 0 and 100

### spec.sources[].weight

`int32`

Weight of this origin source for load balancing across multiple origins.
Valid range: 0-100. If omitted, defaults to 10.
Only effective when multiple sources share the same priority.

- rule: weight must be between 0 and 100

### spec.certificateConfig

`AliCloudCdnDomainCertificateConfig`

HTTPS certificate configuration.
If omitted, the CDN domain serves content over HTTP only.

### spec.certificateConfig.certName

`string`

Display name of the certificate.

### spec.certificateConfig.certType

`string`

The type of certificate.
  - "upload" : upload your own certificate and private key
  - "cas"    : use a certificate from Alibaba Cloud Certificate Management Service
  - "free"   : use a free DV certificate issued by Alibaba Cloud

- rule: cert_type must be one of: upload, cas, free

### spec.certificateConfig.certId

`string`

Certificate ID from Alibaba Cloud Certificate Management Service.
Required when cert_type is "cas".

### spec.certificateConfig.certRegion

`string`

The region of the CAS certificate.
Only relevant when cert_type is "cas".
Supports "cn-hangzhou" (domestic) and "ap-southeast-1" (international).
Defaults to "cn-hangzhou" in the provider.

### spec.certificateConfig.serverCertificate

`string`

PEM-encoded certificate content.
Required when cert_type is "upload".

### spec.certificateConfig.privateKey

`string` · sensitive

PEM-encoded private key content.
Required when cert_type is "upload".

### spec.certificateConfig.serverCertificateStatus

`string`

Whether HTTPS is enabled on the CDN domain.
  - "on"  : HTTPS enabled (default when certificate_config is provided)
  - "off" : HTTPS disabled

- rule: server_certificate_status must be one of: on, off

### spec.checkUrl

`string`

URL used by the CDN service to verify that the origin is reachable.
Alibaba Cloud sends a request to this URL during domain creation.
If omitted, no health check is performed during creation.

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for access control and cost attribution.
If omitted, the domain is placed in the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the CDN domain.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudCdnDomain, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.domain_name` | `string` | The accelerated domain name as registered in CDN. |
| `status.outputs.cname` | `string` | The CNAME assigned by Alibaba Cloud CDN. Create a CNAME record at your DNS provider pointing domain_name to this value for CDN acceleration to take effect. |
| `status.outputs.status` | `string` | The current status of the CDN domain. Typical values: "online", "offline", "configuring", "checking", "check_failed". |

## See Also

- [Overview](../README.md)
