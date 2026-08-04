# AwsCloudFront

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsCloudFrontSpec defines an Amazon CloudFront distribution: the
global CDN front door that terminates TLS at the edge, caches
responses close to viewers, and routes requests to one or more
origins (S3 buckets, load balancers, API endpoints, or anything
HTTP-addressable).

The model mirrors CloudFront's own composition:

  - origins declare WHERE content comes from, each with an origin
    type arm (S3, custom/HTTP, or a provisioned VPC origin);
  - origin_groups compose two origins into a primary/failover pair;
  - default_cache_behavior (plus path-matched
    ordered_cache_behaviors) declares HOW requests are matched,
    cached, and forwarded -- each behavior targets one origin or
    origin group by its user-chosen ID;
  - viewer_certificate + aliases put the distribution on your own
    domain.

Validation proves every behavior target and origin-group member
resolves to a declared origin, so a dangling reference is caught at
manifest time instead of at deploy time.

Distributions have no name in AWS -- identity is the generated ID
(E2ABC...). metadata.name drives the Name identity tag, and
consumers compose through the domain_name / hosted_zone_id outputs
(Route53 alias records) and distribution_arn (WAF, invalidations).

Deploys are slow by nature: CloudFront pushes configuration to every
edge location, so create/update typically takes 5-15 minutes (AWS
caps at far more). wait_for_deployment controls whether the
deployment blocks on that propagation.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCloudFront
metadata:
  name: awscloudfront-demo
spec:
  region: us-east-1
  comment: Demo distribution on the default CloudFront certificate
  origins:
    - originId: s3-assets
      domainName: demo-assets.s3.us-east-1.amazonaws.com
      s3Origin:
        createOriginAccessControl: true
  defaultCacheBehavior:
    targetOriginId: s3-assets
    viewerProtocolPolicy: redirect-to-https
    compress: true
    # Managed-CachingOptimized: the AWS managed cache policy for static
    # content (respects origin cache headers, compresses the cache key).
    cachePolicyId: 658327ea-f89d-4fab-a63d-7e88639e58f6
  defaultRootObject: index.html
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.enabled` | `bool` |  | `true` |  |
| `spec.aliases` | `[]string` |  |  |  |
| `spec.comment` | `string` |  |  |  |
| `spec.defaultRootObject` | `string` |  |  |  |
| `spec.httpVersion` | `string` |  |  |  |
| `spec.isIpv6Enabled` | `bool` |  |  |  |
| `spec.priceClass` | `string` |  |  |  |
| `spec.webAclArn` | `string \| valueFrom` |  |  | AwsWafWebAcl (`status.outputs.web_acl_arn`) |
| `spec.origins` | `[]AwsCloudFrontOrigin` | yes |  |  |
| `spec.origins[].originId` | `string` | yes |  |  |
| `spec.origins[].domainName` | `string` | yes |  |  |
| `spec.origins[].originPath` | `string` |  |  |  |
| `spec.origins[].connectionAttempts` | `int32` |  |  |  |
| `spec.origins[].connectionTimeoutSeconds` | `int32` |  |  |  |
| `spec.origins[].customHeaders` | `[]AwsCloudFrontOriginCustomHeader` |  |  |  |
| `spec.origins[].customHeaders[].name` | `string` | yes |  |  |
| `spec.origins[].customHeaders[].value` | `string` | yes |  |  |
| `spec.origins[].originShield` | `AwsCloudFrontOriginShield` |  |  |  |
| `spec.origins[].originShield.originShieldRegion` | `string` | yes |  |  |
| `spec.origins[].s3Origin` | `AwsCloudFrontS3Origin` |  |  |  |
| `spec.origins[].s3Origin.createOriginAccessControl` | `bool` |  |  |  |
| `spec.origins[].s3Origin.originAccessControlId` | `string` |  |  |  |
| `spec.origins[].s3Origin.originAccessIdentity` | `string` |  |  |  |
| `spec.origins[].customOrigin` | `AwsCloudFrontCustomOrigin` |  |  |  |
| `spec.origins[].customOrigin.protocolPolicy` | `string` | yes |  |  |
| `spec.origins[].customOrigin.httpPort` | `int32` |  |  |  |
| `spec.origins[].customOrigin.httpsPort` | `int32` |  |  |  |
| `spec.origins[].customOrigin.sslProtocols` | `[]string` |  |  |  |
| `spec.origins[].customOrigin.keepaliveTimeoutSeconds` | `int32` |  |  |  |
| `spec.origins[].customOrigin.readTimeoutSeconds` | `int32` |  |  |  |
| `spec.origins[].vpcOrigin` | `AwsCloudFrontVpcOrigin` |  |  |  |
| `spec.origins[].vpcOrigin.vpcOriginId` | `string` | yes |  |  |
| `spec.origins[].vpcOrigin.keepaliveTimeoutSeconds` | `int32` |  |  |  |
| `spec.origins[].vpcOrigin.readTimeoutSeconds` | `int32` |  |  |  |
| `spec.originGroups` | `[]AwsCloudFrontOriginGroup` |  |  |  |
| `spec.originGroups[].originGroupId` | `string` | yes |  |  |
| `spec.originGroups[].memberOriginIds` | `[]string` | yes |  |  |
| `spec.originGroups[].failoverStatusCodes` | `[]int32` | yes |  |  |
| `spec.defaultCacheBehavior` | `AwsCloudFrontCacheBehavior` | yes |  |  |
| `spec.defaultCacheBehavior.targetOriginId` | `string` | yes |  |  |
| `spec.defaultCacheBehavior.viewerProtocolPolicy` | `string` | yes |  |  |
| `spec.defaultCacheBehavior.allowedMethods` | `[]string` |  |  |  |
| `spec.defaultCacheBehavior.cachedMethods` | `[]string` |  |  |  |
| `spec.defaultCacheBehavior.compress` | `bool` |  |  |  |
| `spec.defaultCacheBehavior.cachePolicyId` | `string` |  |  |  |
| `spec.defaultCacheBehavior.originRequestPolicyId` | `string` |  |  |  |
| `spec.defaultCacheBehavior.responseHeadersPolicyId` | `string` |  |  |  |
| `spec.defaultCacheBehavior.forwardedValues` | `AwsCloudFrontForwardedValues` |  |  |  |
| `spec.defaultCacheBehavior.forwardedValues.queryString` | `bool` |  |  |  |
| `spec.defaultCacheBehavior.forwardedValues.queryStringCacheKeys` | `[]string` |  |  |  |
| `spec.defaultCacheBehavior.forwardedValues.headers` | `[]string` |  |  |  |
| `spec.defaultCacheBehavior.forwardedValues.cookiesForward` | `string` | yes |  |  |
| `spec.defaultCacheBehavior.forwardedValues.whitelistedCookieNames` | `[]string` |  |  |  |
| `spec.defaultCacheBehavior.minTtlSeconds` | `int64` |  |  |  |
| `spec.defaultCacheBehavior.defaultTtlSeconds` | `int64` |  |  |  |
| `spec.defaultCacheBehavior.maxTtlSeconds` | `int64` |  |  |  |
| `spec.defaultCacheBehavior.functionAssociations` | `[]AwsCloudFrontFunctionAssociation` |  |  |  |
| `spec.defaultCacheBehavior.functionAssociations[].eventType` | `string` | yes |  |  |
| `spec.defaultCacheBehavior.functionAssociations[].functionArn` | `string` | yes |  |  |
| `spec.defaultCacheBehavior.lambdaFunctionAssociations` | `[]AwsCloudFrontLambdaFunctionAssociation` |  |  |  |
| `spec.defaultCacheBehavior.lambdaFunctionAssociations[].eventType` | `string` | yes |  |  |
| `spec.defaultCacheBehavior.lambdaFunctionAssociations[].lambdaArn` | `string` | yes |  |  |
| `spec.defaultCacheBehavior.lambdaFunctionAssociations[].includeBody` | `bool` |  |  |  |
| `spec.defaultCacheBehavior.trustedKeyGroupIds` | `[]string` |  |  |  |
| `spec.defaultCacheBehavior.trustedSigners` | `[]string` |  |  |  |
| `spec.defaultCacheBehavior.fieldLevelEncryptionId` | `string` |  |  |  |
| `spec.defaultCacheBehavior.realtimeLogConfigArn` | `string` |  |  |  |
| `spec.defaultCacheBehavior.smoothStreaming` | `bool` |  |  |  |
| `spec.defaultCacheBehavior.grpcEnabled` | `bool` |  |  |  |
| `spec.orderedCacheBehaviors` | `[]AwsCloudFrontOrderedCacheBehavior` |  |  |  |
| `spec.orderedCacheBehaviors[].pathPattern` | `string` | yes |  |  |
| `spec.orderedCacheBehaviors[].behavior` | `AwsCloudFrontCacheBehavior` | yes |  |  |
| `spec.orderedCacheBehaviors[].behavior.targetOriginId` | `string` | yes |  |  |
| `spec.orderedCacheBehaviors[].behavior.viewerProtocolPolicy` | `string` | yes |  |  |
| `spec.orderedCacheBehaviors[].behavior.allowedMethods` | `[]string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.cachedMethods` | `[]string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.compress` | `bool` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.cachePolicyId` | `string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.originRequestPolicyId` | `string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.responseHeadersPolicyId` | `string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.forwardedValues` | `AwsCloudFrontForwardedValues` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.forwardedValues.queryString` | `bool` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.forwardedValues.queryStringCacheKeys` | `[]string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.forwardedValues.headers` | `[]string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.forwardedValues.cookiesForward` | `string` | yes |  |  |
| `spec.orderedCacheBehaviors[].behavior.forwardedValues.whitelistedCookieNames` | `[]string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.minTtlSeconds` | `int64` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.defaultTtlSeconds` | `int64` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.maxTtlSeconds` | `int64` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.functionAssociations` | `[]AwsCloudFrontFunctionAssociation` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.functionAssociations[].eventType` | `string` | yes |  |  |
| `spec.orderedCacheBehaviors[].behavior.functionAssociations[].functionArn` | `string` | yes |  |  |
| `spec.orderedCacheBehaviors[].behavior.lambdaFunctionAssociations` | `[]AwsCloudFrontLambdaFunctionAssociation` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.lambdaFunctionAssociations[].eventType` | `string` | yes |  |  |
| `spec.orderedCacheBehaviors[].behavior.lambdaFunctionAssociations[].lambdaArn` | `string` | yes |  |  |
| `spec.orderedCacheBehaviors[].behavior.lambdaFunctionAssociations[].includeBody` | `bool` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.trustedKeyGroupIds` | `[]string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.trustedSigners` | `[]string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.fieldLevelEncryptionId` | `string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.realtimeLogConfigArn` | `string` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.smoothStreaming` | `bool` |  |  |  |
| `spec.orderedCacheBehaviors[].behavior.grpcEnabled` | `bool` |  |  |  |
| `spec.viewerCertificate` | `AwsCloudFrontViewerCertificate` |  |  |  |
| `spec.viewerCertificate.acmCertificateArn` | `string \| valueFrom` |  |  | AwsCertManagerCert (`status.outputs.cert_arn`) |
| `spec.viewerCertificate.iamCertificateId` | `string` |  |  |  |
| `spec.viewerCertificate.sslSupportMethod` | `string` |  |  |  |
| `spec.viewerCertificate.minimumProtocolVersion` | `string` |  |  |  |
| `spec.customErrorResponses` | `[]AwsCloudFrontCustomErrorResponse` |  |  |  |
| `spec.customErrorResponses[].errorCode` | `int32` | yes |  |  |
| `spec.customErrorResponses[].responseCode` | `int32` |  |  |  |
| `spec.customErrorResponses[].responsePagePath` | `string` |  |  |  |
| `spec.customErrorResponses[].errorCachingMinTtlSeconds` | `int64` |  |  |  |
| `spec.geoRestriction` | `AwsCloudFrontGeoRestriction` |  |  |  |
| `spec.geoRestriction.restrictionType` | `string` | yes |  |  |
| `spec.geoRestriction.locations` | `[]string` | yes |  |  |
| `spec.logging` | `AwsCloudFrontLogging` |  |  |  |
| `spec.logging.bucket` | `string` | yes |  |  |
| `spec.logging.prefix` | `string` |  |  |  |
| `spec.logging.includeCookies` | `bool` |  |  |  |
| `spec.waitForDeployment` | `bool` |  | `true` |  |
| `spec.retainOnDelete` | `bool` |  |  |  |
| `spec.enableAdditionalMetrics` | `bool` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region used for the provider connection. CloudFront itself
is a GLOBAL service -- the distribution is the same everywhere and
its supporting resources (certificates, CLOUDFRONT-scope WAF ACLs)
must live in us-east-1 -- but the deployment still runs through a
regional endpoint. Example: "us-east-1".

- rule: {"string":{"minLen":"1"}}

### spec.enabled

`bool` · optional (explicit presence)

Whether the distribution accepts and serves viewer requests. True
(the default) serves traffic; false keeps the distribution
deployed-but-dark (useful for staging a configuration before
cutover, and required by AWS before a distribution can be
deleted -- the modules handle that disable-on-destroy dance).

- default: `true`

### spec.aliases

`[]string`

Alternate domain names (CNAMEs) the distribution answers for,
e.g. "cdn.example.com". Requires viewer_certificate with your own
certificate covering every alias -- CloudFront rejects an alias
the certificate does not cover.

- rule: {"repeated":{"unique":true,"items":{"string":{"pattern":"^[A-Za-z0-9\\*\\-\\.]+\\.[A-Za-z]{2,}$"}}}}

### spec.comment

`string`

Free-form comment shown in the AWS Console distribution list.
Up to 128 characters.

- rule: {"string":{"maxLen":"128"}}

### spec.defaultRootObject

`string`

The object CloudFront serves when a viewer requests the root URL
("/"), e.g. "index.html". Applies to the root only -- subdirectory
index documents need a CloudFront Function or S3 website origin.

- rule: {"string":{"pattern":"^[A-Za-z0-9\\-\\.\\_\\/]*$"}}

### spec.httpVersion

`string`

The maximum HTTP version viewers can use: "http2" (the default),
"http1.1", "http2and3", or "http3". HTTP/3 (QUIC) improves
performance on lossy networks at no extra cost -- "http2and3" is
the safe way to adopt it. Empty keeps http2.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["http1.1","http2","http2and3","http3"]}}

### spec.isIpv6Enabled

`bool`

Whether CloudFront answers IPv6 viewer requests. False (the AWS
default) is IPv4-only; enabling costs nothing and serves
dual-stack (create the AAAA alias record alongside the A record).

### spec.priceClass

`string`

Which edge locations serve the distribution -- the cost/latency
dial: "PriceClass_All" (every edge location -- the default),
"PriceClass_200" (excludes South America + Australia/New Zealand),
or "PriceClass_100" (North America + Europe only, the cheapest).
Viewers outside the selected class are served from the nearest
included edge -- functional everywhere, just slower far away.
Empty keeps PriceClass_All.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["PriceClass_100","PriceClass_200","PriceClass_All"]}}

### spec.webAclArn

`string | valueFrom`

The AWS WAF Web ACL protecting the distribution, by ARN.
CloudFront-scope ACLs (scope CLOUDFRONT, which must live in
us-east-1) are addressed by ARN — never by the bare web ACL ID.
Can reference an AwsWafWebAcl resource.

- references: AwsWafWebAcl (`status.outputs.web_acl_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsWafWebAcl, name: <that resource's name>, fieldPath: status.outputs.web_acl_arn}} -- a bare string does not parse

### spec.origins

`[]AwsCloudFrontOrigin` · required

The content sources the distribution can route to. Each origin's
origin_id is your stable handle -- cache behaviors and origin
groups target it by that ID.

- rule: {"repeated":{"minItems":"1"}}
- rule: set at most one of s3_origin, custom_origin, or vpc_origin -- an origin has exactly one type

### spec.origins[].originId

`string` · required

Your stable handle for this origin -- cache behaviors and origin
groups target it (e.g. "s3-assets", "api-backend"). Must be
unique within the distribution.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128","pattern":"^[A-Za-z0-9\\-\\_\\.]+$"}}

### spec.origins[].domainName

`string` · required

The DNS name CloudFront connects to. For S3 REST origins use the
REGIONAL bucket endpoint ("bucket.s3.us-west-2.amazonaws.com" --
the bucket_regional_domain_name output of AwsS3Bucket); for
load balancers the LB DNS name; for S3 static websites the
website endpoint (as a custom_origin -- website endpoints speak
plain HTTP only).

- rule: {"required":true,"string":{"pattern":"^[A-Za-z0-9\\-\\.]+\\.[A-Za-z0-9]{2,}$"}}

### spec.origins[].originPath

`string`

Optional path CloudFront appends to every origin request, e.g.
"/production" to serve from a bucket sub-prefix. Must start with
"/" and not end with one.

- rule: {"string":{"pattern":"^$|^/[A-Za-z0-9\\-\\_\\/\\.]*[A-Za-z0-9\\-\\_\\.]$"}}

### spec.origins[].connectionAttempts

`int32`

How many times CloudFront attempts to connect to the origin per
request, 1-3. 0 keeps the AWS default (3).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":3,"gte":1}}

### spec.origins[].connectionTimeoutSeconds

`int32`

Seconds CloudFront waits when establishing each connection
attempt, 1-10. 0 keeps the AWS default (10).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":10,"gte":1}}

### spec.origins[].customHeaders

`[]AwsCloudFrontOriginCustomHeader`

Headers CloudFront adds to every request it sends to this origin.
The classic use is a shared-secret header (e.g. "X-Origin-Verify")
the origin checks so it only serves CloudFront traffic.

### spec.origins[].customHeaders[].name

`string` · required

The header name, e.g. "X-Origin-Verify".

- rule: {"required":true}

### spec.origins[].customHeaders[].value

`string` · required

The header value. Treat shared-secret values like configuration,
not state secrets -- rotate them by updating the distribution.

- rule: {"required":true}

### spec.origins[].originShield

`AwsCloudFrontOriginShield`

Origin Shield: an extra regional caching layer in front of the
origin that collapses requests from all edge locations, cutting
origin load dramatically for origin-heavy workloads. Billed per
request; choose the region closest to the origin.

### spec.origins[].originShield.originShieldRegion

`string` · required

The AWS region for the Origin Shield cache. Choose the region
with the lowest latency to the origin (usually the origin's own
region). Example: "us-west-2".

- rule: {"required":true}

### spec.origins[].s3Origin

`AwsCloudFrontS3Origin`

S3 REST origin arm: access control for a private bucket.

- rule: set at most one of create_origin_access_control, origin_access_control_id, or origin_access_identity

### spec.origins[].s3Origin.createOriginAccessControl

`bool`

Create and attach an Origin Access Control for this origin -- the
modern way to serve a private bucket (SigV4 signing, works with
SSE-KMS and all regions). The distribution's ARN must be allowed
in the bucket policy ("AWS:SourceArn" condition on the
"cloudfront.amazonaws.com" principal).

### spec.origins[].s3Origin.originAccessControlId

`string`

Attach an EXISTING Origin Access Control by ID instead of
creating one -- for sharing one OAC across distributions.

### spec.origins[].s3Origin.originAccessIdentity

`string`

Attach an existing legacy Origin Access Identity, as the full
"origin-access-identity/cloudfront/<ID>" path. OAI is the
predecessor of OAC (no SSE-KMS support, not recommended for new
configurations) -- accepted here for buckets already wired to
one, never created.

- rule: {"string":{"pattern":"^$|^origin-access-identity/cloudfront/[A-Z0-9]+$"}}

### spec.origins[].customOrigin

`AwsCloudFrontCustomOrigin`

Custom (HTTP) origin arm: ports, protocol, and timeouts for load
balancers, API endpoints, or any HTTP server.

### spec.origins[].customOrigin.protocolPolicy

`string` · required

How CloudFront connects to the origin: "https-only" (recommended),
"http-only" (required for S3 website endpoints, which do not speak
HTTPS), or "match-viewer" (mirrors the viewer's protocol).

- rule: {"required":true,"string":{"in":["http-only","https-only","match-viewer"]}}

### spec.origins[].customOrigin.httpPort

`int32`

The origin's HTTP port. 0 keeps the default (80).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":65535,"gte":1}}

### spec.origins[].customOrigin.httpsPort

`int32`

The origin's HTTPS port. 0 keeps the default (443).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":65535,"gte":1}}

### spec.origins[].customOrigin.sslProtocols

`[]string`

TLS versions CloudFront may use to the origin. Empty keeps
["TLSv1.2"] -- the safe modern floor. Only add older protocols for
legacy origins that cannot do better.

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["SSLv3","TLSv1","TLSv1.1","TLSv1.2"]}}}}

### spec.origins[].customOrigin.keepaliveTimeoutSeconds

`int32`

Seconds an idle connection to the origin is kept open, 1-60 (up
to 120 by AWS quota increase). Raising it helps request-heavy
workloads reuse connections. 0 keeps the AWS default (5).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":120,"gte":1}}

### spec.origins[].customOrigin.readTimeoutSeconds

`int32`

Seconds CloudFront waits for an origin response, 1-60 (up to 180
by AWS quota increase). Also the cap on how long a slow API may
take before viewers see a 504. 0 keeps the AWS default (30).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":180,"gte":1}}

### spec.origins[].vpcOrigin

`AwsCloudFrontVpcOrigin`

VPC origin arm: route to a provisioned CloudFront VPC origin that
reaches a private ALB/NLB/instance inside your VPC without any
public exposure.

### spec.origins[].vpcOrigin.vpcOriginId

`string` · required

The VPC origin ID (the CloudFront vpc_origin resource's ID, not a
VPC ID).

- rule: {"required":true}

### spec.origins[].vpcOrigin.keepaliveTimeoutSeconds

`int32`

Seconds an idle connection is kept open, 1-120. 0 keeps the AWS
default (5).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":120,"gte":1}}

### spec.origins[].vpcOrigin.readTimeoutSeconds

`int32`

Seconds CloudFront waits for a response, 1-180. 0 keeps the AWS
default (30).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":180,"gte":1}}

### spec.originGroups

`[]AwsCloudFrontOriginGroup`

Primary/failover origin pairs. A behavior that targets a group's
ID gets automatic failover: CloudFront retries the second member
when the first returns one of the configured status codes.

### spec.originGroups[].originGroupId

`string` · required

Your stable handle for this group -- unique across origins and
groups.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128","pattern":"^[A-Za-z0-9\\-\\_\\.]+$"}}

### spec.originGroups[].memberOriginIds

`[]string` · required

Exactly two member origin IDs, primary first. Members must be
declared origins (groups cannot nest).

- rule: {"repeated":{"minItems":"2","maxItems":"2","unique":true}}

### spec.originGroups[].failoverStatusCodes

`[]int32` · required

The origin HTTP status codes that trigger failover to the second
member, from 400, 403, 404, 416, 500, 502, 503, 504.

- rule: {"repeated":{"minItems":"1","unique":true,"items":{"int32":{"in":[400,403,404,416,500,502,503,504]}}}}

### spec.defaultCacheBehavior

`AwsCloudFrontCacheBehavior` · required

How requests that match no ordered behavior are cached and
forwarded -- every distribution has exactly one default behavior.

- rule: {"required":true}
- rule: set exactly one of cache_policy_id (modern, recommended) or forwarded_values (legacy) -- CloudFront requires a caching configuration
- rule: min_ttl_seconds/default_ttl_seconds/max_ttl_seconds only apply with forwarded_values -- with cache_policy_id the policy owns the TTLs

### spec.defaultCacheBehavior.targetOriginId

`string` · required

The origin_id or origin_group_id this behavior routes to.

- rule: {"required":true}

### spec.defaultCacheBehavior.viewerProtocolPolicy

`string` · required

How viewers reach this content: "redirect-to-https" (the
sensible default for websites), "https-only" (APIs; no redirect
round-trip), or "allow-all" (legacy plain-HTTP viewers).

- rule: {"required":true,"string":{"in":["allow-all","https-only","redirect-to-https"]}}

### spec.defaultCacheBehavior.allowedMethods

`[]string`

HTTP methods CloudFront forwards to the origin. Empty keeps
["GET", "HEAD"] -- static content. Use the full seven
(GET/HEAD/OPTIONS/PUT/POST/PATCH/DELETE) for APIs; only GET/HEAD
(+OPTIONS) responses are ever cached.

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["GET","HEAD","OPTIONS","PUT","POST","PATCH","DELETE"]}}}}

### spec.defaultCacheBehavior.cachedMethods

`[]string`

Methods whose responses CloudFront caches -- ["GET", "HEAD"] (the
default when empty) or ["GET", "HEAD", "OPTIONS"] (cache CORS
preflights).

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["GET","HEAD","OPTIONS"]}}}}

### spec.defaultCacheBehavior.compress

`bool`

Automatically gzip/brotli-compress responses for viewers that
accept it -- almost always what you want for text content (the
cache policy must enable the compression formats too when using
the modern generation).

### spec.defaultCacheBehavior.cachePolicyId

`string`

MODERN generation: the cache policy controlling the cache key and
TTLs, by ID -- a managed policy ID (see the message comment) or
your own. Mutually exclusive with forwarded_values and the
per-behavior TTLs.

### spec.defaultCacheBehavior.originRequestPolicyId

`string`

MODERN generation: the origin-request policy controlling which
headers/cookies/query strings are forwarded (WITHOUT joining the
cache key), by ID.

### spec.defaultCacheBehavior.responseHeadersPolicyId

`string`

The response-headers policy adding/removing headers on responses
(security headers, CORS, server-timing), by ID. Works with both
generations.

### spec.defaultCacheBehavior.forwardedValues

`AwsCloudFrontForwardedValues`

LEGACY generation: inline forwarding rules; forwarded values join
the cache key. Mutually exclusive with cache_policy_id.

- rule: whitelisted_cookie_names requires cookies_forward: whitelist

### spec.defaultCacheBehavior.forwardedValues.queryString

`bool`

Forward query strings to the origin (and cache on them).

### spec.defaultCacheBehavior.forwardedValues.queryStringCacheKeys

`[]string`

When query_string is true, cache only on these parameters (empty
caches on all of them).

- rule: {"repeated":{"unique":true}}

### spec.defaultCacheBehavior.forwardedValues.headers

`[]string`

Header names forwarded to the origin and joined into the cache
key. "*" forwards all headers and effectively disables caching.

- rule: {"repeated":{"unique":true}}

### spec.defaultCacheBehavior.forwardedValues.cookiesForward

`string` · required

Cookie forwarding: "none" (the safe default for cacheable
content), "whitelist" (only whitelisted_cookie_names), or "all"
(destroys cache efficiency -- every cookie combination is its own
cache entry).

- rule: {"required":true,"string":{"in":["none","whitelist","all"]}}

### spec.defaultCacheBehavior.forwardedValues.whitelistedCookieNames

`[]string`

Cookie names forwarded when cookies_forward is "whitelist".

- rule: {"repeated":{"unique":true}}

### spec.defaultCacheBehavior.minTtlSeconds

`int64`

LEGACY generation TTL floor in seconds (with forwarded_values
only). Default 0.

- rule: {"int64":{"gte":"0"}}

### spec.defaultCacheBehavior.defaultTtlSeconds

`int64`

LEGACY generation default TTL in seconds, used when the origin
sends no caching headers (with forwarded_values only). 0 keeps
the AWS default (86400 -- one day).

- rule: {"int64":{"gte":"0"}}

### spec.defaultCacheBehavior.maxTtlSeconds

`int64`

LEGACY generation TTL ceiling in seconds (with forwarded_values
only). 0 keeps the AWS default (31536000 -- one year).

- rule: {"int64":{"gte":"0"}}

### spec.defaultCacheBehavior.functionAssociations

`[]AwsCloudFrontFunctionAssociation`

CloudFront Functions (sub-millisecond JavaScript at every edge)
attached to this behavior -- viewer-request/viewer-response only,
at most one per event type. The lightweight choice for URL
rewrites, redirects, and header manipulation.

- rule: {"repeated":{"maxItems":"2"}}

### spec.defaultCacheBehavior.functionAssociations[].eventType

`string` · required

The event: "viewer-request" (before the cache lookup) or
"viewer-response" (before returning to the viewer).

- rule: {"required":true,"string":{"in":["viewer-request","viewer-response"]}}

### spec.defaultCacheBehavior.functionAssociations[].functionArn

`string` · required

The CloudFront Function ARN.

- rule: {"required":true,"string":{"pattern":"^arn:aws[a-zA-Z-]*:cloudfront::[0-9]{12}:function/.+$"}}

### spec.defaultCacheBehavior.lambdaFunctionAssociations

`[]AwsCloudFrontLambdaFunctionAssociation`

Lambda@Edge functions (full Lambda at regional edges) attached to
this behavior -- all four event types, at most one per type. For
logic that needs the network, the body, or more than 1ms.

- rule: {"repeated":{"maxItems":"4"}}

### spec.defaultCacheBehavior.lambdaFunctionAssociations[].eventType

`string` · required

The event: "viewer-request"/"viewer-response" (every request; 5s
limit) or "origin-request"/"origin-response" (cache misses only;
30s limit -- the usual choice for origin rewrites).

- rule: {"required":true,"string":{"in":["viewer-request","viewer-response","origin-request","origin-response"]}}

### spec.defaultCacheBehavior.lambdaFunctionAssociations[].lambdaArn

`string` · required

The Lambda function VERSION ARN (Lambda@Edge requires a numbered
version, never $LATEST or an alias). The function must live in
us-east-1.

- rule: {"required":true,"string":{"pattern":"^arn:aws[a-zA-Z-]*:lambda:us-east-1:[0-9]{12}:function:[^:]+:[0-9]+$"}}

### spec.defaultCacheBehavior.lambdaFunctionAssociations[].includeBody

`bool`

Expose the request body to the function (viewer-request and
origin-request events).

### spec.defaultCacheBehavior.trustedKeyGroupIds

`[]string`

Restrict this behavior's content to signed URLs/cookies from
these key groups (the modern private-content mechanism).

- rule: {"repeated":{"unique":true}}

### spec.defaultCacheBehavior.trustedSigners

`[]string`

LEGACY private content: AWS account numbers whose CloudFront key
pairs may sign URLs. Prefer trusted_key_group_ids.

- rule: {"repeated":{"unique":true}}

### spec.defaultCacheBehavior.fieldLevelEncryptionId

`string`

The field-level encryption configuration ID applied to this
behavior (encrypts specific POST fields at the edge with your
public key).

### spec.defaultCacheBehavior.realtimeLogConfigArn

`string`

The real-time log configuration ARN streaming this behavior's
requests to Kinesis within seconds.

### spec.defaultCacheBehavior.smoothStreaming

`bool`

Serve Microsoft Smooth Streaming media from this behavior.

### spec.defaultCacheBehavior.grpcEnabled

`bool`

Allow gRPC viewer traffic over HTTP/2 (requires POST in
allowed_methods and http_version http2 or above).

### spec.orderedCacheBehaviors

`[]AwsCloudFrontOrderedCacheBehavior`

Path-matched behaviors evaluated IN ORDER before the default --
first match wins. The classic use: route "/api/*" to a load
balancer with caching disabled while everything else serves from
S3.

### spec.orderedCacheBehaviors[].pathPattern

`string` · required

The path pattern this behavior matches, e.g. "/api/*",
"/images/*.jpg", "*.css". First matching behavior wins.

- rule: {"required":true}

### spec.orderedCacheBehaviors[].behavior

`AwsCloudFrontCacheBehavior` · required

The behavior applied to matching requests.

- rule: {"required":true}
- rule: set exactly one of cache_policy_id (modern, recommended) or forwarded_values (legacy) -- CloudFront requires a caching configuration
- rule: min_ttl_seconds/default_ttl_seconds/max_ttl_seconds only apply with forwarded_values -- with cache_policy_id the policy owns the TTLs

### spec.orderedCacheBehaviors[].behavior.targetOriginId

`string` · required

The origin_id or origin_group_id this behavior routes to.

- rule: {"required":true}

### spec.orderedCacheBehaviors[].behavior.viewerProtocolPolicy

`string` · required

How viewers reach this content: "redirect-to-https" (the
sensible default for websites), "https-only" (APIs; no redirect
round-trip), or "allow-all" (legacy plain-HTTP viewers).

- rule: {"required":true,"string":{"in":["allow-all","https-only","redirect-to-https"]}}

### spec.orderedCacheBehaviors[].behavior.allowedMethods

`[]string`

HTTP methods CloudFront forwards to the origin. Empty keeps
["GET", "HEAD"] -- static content. Use the full seven
(GET/HEAD/OPTIONS/PUT/POST/PATCH/DELETE) for APIs; only GET/HEAD
(+OPTIONS) responses are ever cached.

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["GET","HEAD","OPTIONS","PUT","POST","PATCH","DELETE"]}}}}

### spec.orderedCacheBehaviors[].behavior.cachedMethods

`[]string`

Methods whose responses CloudFront caches -- ["GET", "HEAD"] (the
default when empty) or ["GET", "HEAD", "OPTIONS"] (cache CORS
preflights).

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["GET","HEAD","OPTIONS"]}}}}

### spec.orderedCacheBehaviors[].behavior.compress

`bool`

Automatically gzip/brotli-compress responses for viewers that
accept it -- almost always what you want for text content (the
cache policy must enable the compression formats too when using
the modern generation).

### spec.orderedCacheBehaviors[].behavior.cachePolicyId

`string`

MODERN generation: the cache policy controlling the cache key and
TTLs, by ID -- a managed policy ID (see the message comment) or
your own. Mutually exclusive with forwarded_values and the
per-behavior TTLs.

### spec.orderedCacheBehaviors[].behavior.originRequestPolicyId

`string`

MODERN generation: the origin-request policy controlling which
headers/cookies/query strings are forwarded (WITHOUT joining the
cache key), by ID.

### spec.orderedCacheBehaviors[].behavior.responseHeadersPolicyId

`string`

The response-headers policy adding/removing headers on responses
(security headers, CORS, server-timing), by ID. Works with both
generations.

### spec.orderedCacheBehaviors[].behavior.forwardedValues

`AwsCloudFrontForwardedValues`

LEGACY generation: inline forwarding rules; forwarded values join
the cache key. Mutually exclusive with cache_policy_id.

- rule: whitelisted_cookie_names requires cookies_forward: whitelist

### spec.orderedCacheBehaviors[].behavior.forwardedValues.queryString

`bool`

Forward query strings to the origin (and cache on them).

### spec.orderedCacheBehaviors[].behavior.forwardedValues.queryStringCacheKeys

`[]string`

When query_string is true, cache only on these parameters (empty
caches on all of them).

- rule: {"repeated":{"unique":true}}

### spec.orderedCacheBehaviors[].behavior.forwardedValues.headers

`[]string`

Header names forwarded to the origin and joined into the cache
key. "*" forwards all headers and effectively disables caching.

- rule: {"repeated":{"unique":true}}

### spec.orderedCacheBehaviors[].behavior.forwardedValues.cookiesForward

`string` · required

Cookie forwarding: "none" (the safe default for cacheable
content), "whitelist" (only whitelisted_cookie_names), or "all"
(destroys cache efficiency -- every cookie combination is its own
cache entry).

- rule: {"required":true,"string":{"in":["none","whitelist","all"]}}

### spec.orderedCacheBehaviors[].behavior.forwardedValues.whitelistedCookieNames

`[]string`

Cookie names forwarded when cookies_forward is "whitelist".

- rule: {"repeated":{"unique":true}}

### spec.orderedCacheBehaviors[].behavior.minTtlSeconds

`int64`

LEGACY generation TTL floor in seconds (with forwarded_values
only). Default 0.

- rule: {"int64":{"gte":"0"}}

### spec.orderedCacheBehaviors[].behavior.defaultTtlSeconds

`int64`

LEGACY generation default TTL in seconds, used when the origin
sends no caching headers (with forwarded_values only). 0 keeps
the AWS default (86400 -- one day).

- rule: {"int64":{"gte":"0"}}

### spec.orderedCacheBehaviors[].behavior.maxTtlSeconds

`int64`

LEGACY generation TTL ceiling in seconds (with forwarded_values
only). 0 keeps the AWS default (31536000 -- one year).

- rule: {"int64":{"gte":"0"}}

### spec.orderedCacheBehaviors[].behavior.functionAssociations

`[]AwsCloudFrontFunctionAssociation`

CloudFront Functions (sub-millisecond JavaScript at every edge)
attached to this behavior -- viewer-request/viewer-response only,
at most one per event type. The lightweight choice for URL
rewrites, redirects, and header manipulation.

- rule: {"repeated":{"maxItems":"2"}}

### spec.orderedCacheBehaviors[].behavior.functionAssociations[].eventType

`string` · required

The event: "viewer-request" (before the cache lookup) or
"viewer-response" (before returning to the viewer).

- rule: {"required":true,"string":{"in":["viewer-request","viewer-response"]}}

### spec.orderedCacheBehaviors[].behavior.functionAssociations[].functionArn

`string` · required

The CloudFront Function ARN.

- rule: {"required":true,"string":{"pattern":"^arn:aws[a-zA-Z-]*:cloudfront::[0-9]{12}:function/.+$"}}

### spec.orderedCacheBehaviors[].behavior.lambdaFunctionAssociations

`[]AwsCloudFrontLambdaFunctionAssociation`

Lambda@Edge functions (full Lambda at regional edges) attached to
this behavior -- all four event types, at most one per type. For
logic that needs the network, the body, or more than 1ms.

- rule: {"repeated":{"maxItems":"4"}}

### spec.orderedCacheBehaviors[].behavior.lambdaFunctionAssociations[].eventType

`string` · required

The event: "viewer-request"/"viewer-response" (every request; 5s
limit) or "origin-request"/"origin-response" (cache misses only;
30s limit -- the usual choice for origin rewrites).

- rule: {"required":true,"string":{"in":["viewer-request","viewer-response","origin-request","origin-response"]}}

### spec.orderedCacheBehaviors[].behavior.lambdaFunctionAssociations[].lambdaArn

`string` · required

The Lambda function VERSION ARN (Lambda@Edge requires a numbered
version, never $LATEST or an alias). The function must live in
us-east-1.

- rule: {"required":true,"string":{"pattern":"^arn:aws[a-zA-Z-]*:lambda:us-east-1:[0-9]{12}:function:[^:]+:[0-9]+$"}}

### spec.orderedCacheBehaviors[].behavior.lambdaFunctionAssociations[].includeBody

`bool`

Expose the request body to the function (viewer-request and
origin-request events).

### spec.orderedCacheBehaviors[].behavior.trustedKeyGroupIds

`[]string`

Restrict this behavior's content to signed URLs/cookies from
these key groups (the modern private-content mechanism).

- rule: {"repeated":{"unique":true}}

### spec.orderedCacheBehaviors[].behavior.trustedSigners

`[]string`

LEGACY private content: AWS account numbers whose CloudFront key
pairs may sign URLs. Prefer trusted_key_group_ids.

- rule: {"repeated":{"unique":true}}

### spec.orderedCacheBehaviors[].behavior.fieldLevelEncryptionId

`string`

The field-level encryption configuration ID applied to this
behavior (encrypts specific POST fields at the edge with your
public key).

### spec.orderedCacheBehaviors[].behavior.realtimeLogConfigArn

`string`

The real-time log configuration ARN streaming this behavior's
requests to Kinesis within seconds.

### spec.orderedCacheBehaviors[].behavior.smoothStreaming

`bool`

Serve Microsoft Smooth Streaming media from this behavior.

### spec.orderedCacheBehaviors[].behavior.grpcEnabled

`bool`

Allow gRPC viewer traffic over HTTP/2 (requires POST in
allowed_methods and http_version http2 or above).

### spec.viewerCertificate

`AwsCloudFrontViewerCertificate`

The certificate presented to viewers. Absent means the default
*.cloudfront.net certificate (no aliases possible). Set the ACM
arm (recommended) or the IAM arm to serve your own domains --
the ACM certificate MUST live in us-east-1.

- rule: set at most one of acm_certificate_arn or iam_certificate_id

### spec.viewerCertificate.acmCertificateArn

`string | valueFrom`

The ACM certificate ARN -- MUST be in us-east-1 regardless of
where anything else lives. Can reference an AwsCertManagerCert
resource. The recommended arm.

- references: AwsCertManagerCert (`status.outputs.cert_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCertManagerCert, name: <that resource's name>, fieldPath: status.outputs.cert_arn}} -- a bare string does not parse

### spec.viewerCertificate.iamCertificateId

`string`

A certificate uploaded to IAM, by ID -- the legacy arm for
regions/paths where ACM is unavailable.

### spec.viewerCertificate.sslSupportMethod

`string`

How CloudFront serves HTTPS for custom certificates: "sni-only"
(the default -- free, supported by every modern client) or "vip"
(dedicated IPs at significant monthly cost, for ancient non-SNI
clients). Empty keeps sni-only.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["sni-only","vip"]}}

### spec.viewerCertificate.minimumProtocolVersion

`string`

The minimum TLS version viewers must speak, e.g. "TLSv1.2_2021"
(the recommended modern floor, and the module default for custom
certificates), "TLSv1.2_2019", "TLSv1.1_2016", "TLSv1_2016",
"TLSv1". Empty keeps TLSv1.2_2021.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["TLSv1","TLSv1_2016","TLSv1.1_2016","TLSv1.2_2018","TLSv1.2_2019","TLSv1.2_2021"]}}

### spec.customErrorResponses

`[]AwsCloudFrontCustomErrorResponse`

Replace origin error responses with custom pages and control how
long errors are cached -- e.g. serve "/errors/404.html" with a
404, or map S3's 403-for-missing-object to a 404.

- rule: response_page_path requires response_code

### spec.customErrorResponses[].errorCode

`int32` · required

The origin HTTP status code to intercept.

- rule: {"required":true,"int32":{"in":[400,403,404,405,414,416,500,501,502,503,504]}}

### spec.customErrorResponses[].responseCode

`int32`

The status code returned to the viewer (e.g. map S3's
403-for-missing-object to 404). 0 passes the original code
through.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":599,"gte":200}}

### spec.customErrorResponses[].responsePagePath

`string`

The page served for this error, as a distribution path (e.g.
"/errors/404.html"). The path must be servable by a behavior.

- rule: {"string":{"pattern":"^$|^/.*$"}}

### spec.customErrorResponses[].errorCachingMinTtlSeconds

`int64`

Seconds CloudFront caches this error response. 0 keeps the AWS
default (300); errors caching too long can mask origin recovery.

- rule: {"int64":{"gte":"0"}}

### spec.geoRestriction

`AwsCloudFrontGeoRestriction`

Allow or deny viewers by country (ISO 3166-1 alpha-2 codes).
Absent means no geographic restriction.

### spec.geoRestriction.restrictionType

`string` · required

"whitelist" (serve ONLY the listed countries) or "blacklist"
(serve everyone EXCEPT the listed countries).

- rule: {"required":true,"string":{"in":["whitelist","blacklist"]}}

### spec.geoRestriction.locations

`[]string` · required

ISO 3166-1 alpha-2 country codes, e.g. "US", "DE", "IN".

- rule: {"repeated":{"minItems":"1","unique":true,"items":{"string":{"pattern":"^[A-Z]{2}$"}}}}

### spec.logging

`AwsCloudFrontLogging`

Standard access logs delivered to an S3 bucket (best-effort
delivery, typically within minutes). Absent disables logging.

### spec.logging.bucket

`string` · required

The S3 bucket receiving logs, as a bucket DOMAIN NAME
("my-logs.s3.amazonaws.com"). The bucket must have ACLs enabled
(object ownership "Bucket owner preferred") -- CloudFront writes
logs via the awslogsdelivery canonical user.

- rule: {"required":true,"string":{"pattern":"^[A-Za-z0-9\\-\\.]+\\.s3(\\.[a-z0-9-]+)?\\.amazonaws\\.com$"}}

### spec.logging.prefix

`string`

Key prefix for log objects, e.g. "cdn/".

### spec.logging.includeCookies

`bool`

Include cookies in the logs.

### spec.waitForDeployment

`bool` · optional (explicit presence)

Whether the deployment blocks until CloudFront reports the
distribution Deployed at every edge location (typically 5-15
minutes). True (the default) means downstream resources see a
live distribution; false returns as soon as the configuration is
accepted -- the distribution keeps serving the previous
configuration until propagation completes.

- default: `true`

### spec.retainOnDelete

`bool`

Disable the distribution instead of deleting it on destroy. The
distribution stops serving but stays in the account for manual
deletion later -- an escape hatch for teardown ordering problems;
leave false for the normal full delete.

### spec.enableAdditionalMetrics

`bool`

Turn on CloudWatch additional metrics (cache hit rate, origin
latency, error rates by status code) for the distribution. Billed
per the CloudWatch custom-metric rate; indispensable for
production cache tuning.

## Validation Rules

- `aliases_require_certificate`: serving aliases requires viewer_certificate with the ACM or IAM arm set -- the default *.cloudfront.net certificate cannot cover custom domains
- `origin_ids_unique`: origin_id values must be unique across origins and origin_groups -- cache behaviors target them by ID
- `default_behavior_target_resolves`: default_cache_behavior.target_origin_id must match the origin_id of a declared origin or the origin_group_id of a declared origin group
- `ordered_behavior_targets_resolve`: every ordered_cache_behaviors[].behavior.target_origin_id must match the origin_id of a declared origin or the origin_group_id of a declared origin group
- `origin_group_members_resolve`: every origin_groups[].member_origin_ids entry must match the origin_id of a declared origin

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsCloudFront, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.distribution_id` | `string` | The distribution ID (e.g. "E2ABCDEF123456") -- what invalidation requests and monitoring subscriptions key on. |
| `status.outputs.distribution_arn` | `string` | The distribution ARN -- what WAF associations and resource policies reference. |
| `status.outputs.domain_name` | `string` | The CloudFront domain name (e.g. "d123abc.cloudfront.net") -- the target for Route53 alias records and CNAMEs pointing custom domains at the distribution. |
| `status.outputs.hosted_zone_id` | `string` | The Route53 hosted zone ID for CloudFront alias records -- always "Z2FDTNDATAQYW2" (CloudFront's global zone), exported so alias records can compose without hardcoding it. |
| `status.outputs.status` | `string` | The distribution status at the end of the deployment: "Deployed" (propagated to every edge location) or "InProgress" (still propagating -- the resting state when wait_for_deployment is false). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.webAclArn` | AwsWafWebAcl | `status.outputs.web_acl_arn` |
| `spec.viewerCertificate.acmCertificateArn` | AwsCertManagerCert | `status.outputs.cert_arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
