# AzureApplicationInsightsStandardWebTest

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureApplicationInsightsStandardWebTestSpec** defines an Application
Insights Standard Web Test: a synthetic availability monitor that issues
an HTTP request to a URL from one or more Azure regions on a schedule and
records whether the endpoint responded correctly. It is how you prove an
endpoint is reachable and healthy from the outside -- and, paired with a
metric alert on the web test's availability, how you get paged when it is
not.

The test binds to an Application Insights component, which stores its
results and is where its availability metric lives. Run it from multiple
geo-locations to distinguish a real outage from a single-region network
blip.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureApplicationInsightsStandardWebTest
metadata:
  name: test-web-test
spec:
  resourceGroup:
    value: test-rg
  name: homepage-health
  region: eastus
  applicationInsightsId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Insights/components/test-ai
  frequency: 300
  timeout: 30
  request:
    url: https://example.com/health
    httpVerb: GET
  validationRules:
    expectedStatusCode: 200
    sslCheckEnabled: true
    sslCertRemainingLifetime: 30
    content:
      contentMatch: healthy
  geoLocations:
    - us-east-1-azr
    - us-west-2-azr
  description: Homepage availability from two US regions
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.applicationInsightsId` | `string \| valueFrom` | yes |  | AzureApplicationInsights (`status.outputs.application_insights_id`) |
| `spec.frequency` | `int32` |  |  |  |
| `spec.timeout` | `int32` |  |  |  |
| `spec.enabled` | `bool` |  |  |  |
| `spec.retryEnabled` | `bool` |  |  |  |
| `spec.request` | `AzureApplicationInsightsStandardWebTestRequest` | yes |  |  |
| `spec.request.url` | `string` | yes |  |  |
| `spec.request.httpVerb` | `string` |  |  |  |
| `spec.request.body` | `string` |  |  |  |
| `spec.request.followRedirectsEnabled` | `bool` |  |  |  |
| `spec.request.parseDependentRequestsEnabled` | `bool` |  |  |  |
| `spec.request.headers` | `[]AzureApplicationInsightsStandardWebTestHeader` |  |  |  |
| `spec.request.headers[].name` | `string` | yes |  |  |
| `spec.request.headers[].value` | `string` | yes |  |  |
| `spec.validationRules` | `AzureApplicationInsightsStandardWebTestValidationRules` |  |  |  |
| `spec.validationRules.expectedStatusCode` | `int32` |  |  |  |
| `spec.validationRules.sslCertRemainingLifetime` | `int32` |  |  |  |
| `spec.validationRules.sslCheckEnabled` | `bool` |  |  |  |
| `spec.validationRules.content` | `AzureApplicationInsightsStandardWebTestContentValidation` |  |  |  |
| `spec.validationRules.content.contentMatch` | `string` | yes |  |  |
| `spec.validationRules.content.ignoreCase` | `bool` |  |  |  |
| `spec.validationRules.content.passIfTextFound` | `bool` |  |  |  |
| `spec.geoLocations` | `[]string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the web test is created in. Can be a literal
resource-group name or a reference to an AzureResourceGroup's name
output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the web test, unique within the resource group. Fixed at
creation.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.region

`string` · required

The region the web test resource lives in -- typically the same region
as its Application Insights component. (This is distinct from
geo_locations, which are where the test RUNS FROM.) Fixed at creation.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.applicationInsightsId

`string | valueFrom` · required

The Application Insights component that stores the test's results and
hosts its availability metric. Defaults to referencing an
AzureApplicationInsights component's application_insights_id output.
Fixed at creation.

- references: AzureApplicationInsights (`status.outputs.application_insights_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureApplicationInsights, name: <that resource's name>, fieldPath: status.outputs.application_insights_id}} -- a bare string does not parse

### spec.frequency

`int32` · optional (explicit presence)

How often the test runs, in seconds. One of 300 (5 min, the default),
600 (10 min), or 900 (15 min). Unspecified applies 300.

- rule: frequency must be 300, 600, or 900 seconds

### spec.timeout

`int32` · optional (explicit presence)

How long, in seconds, to wait for a response before the run is a
failure. Unspecified applies Azure's default (30). Must be positive.

- rule: {"int32":{"gte":0}}

### spec.enabled

`bool` · optional (explicit presence)

Whether the test is active. Unspecified leaves it enabled (Azure's
default). Set false to keep the definition but stop running it.

### spec.retryEnabled

`bool` · optional (explicit presence)

Whether a failed run is retried before it counts as a failure --
reduces false alarms from transient blips. Unspecified applies Azure's
default.

### spec.request

`AzureApplicationInsightsStandardWebTestRequest` · required

The HTTP request the test issues. Required -- the test needs a URL to
hit.

- rule: {"required":true}

### spec.request.url

`string` · required

The URL the test requests. Required.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.request.httpVerb

`string`

The HTTP method. Unspecified applies Azure's default (GET). One of GET,
POST, PUT, PATCH, DELETE, HEAD, OPTIONS.

- rule: http_verb must be one of GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS

### spec.request.body

`string`

A request body (for POST/PUT/PATCH). Leave empty for methods without a
body.

### spec.request.followRedirectsEnabled

`bool` · optional (explicit presence)

Whether the test follows HTTP redirects. Unspecified applies Azure's
default (true).

### spec.request.parseDependentRequestsEnabled

`bool` · optional (explicit presence)

Whether the test also loads the page's dependent requests (images,
scripts) and counts their success. Unspecified applies Azure's default
(true).

### spec.request.headers

`[]AzureApplicationInsightsStandardWebTestHeader`

Extra request headers to send.

### spec.request.headers[].name

`string` · required

The header name.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.request.headers[].value

`string` · required

The header value.

- rule: {"required":true}

### spec.validationRules

`AzureApplicationInsightsStandardWebTestValidationRules`

Optional response validation beyond a 200 -- assert a status code, an
SSL certificate lifetime floor, or that the body contains (or lacks)
specific text.

- rule: ssl_cert_remaining_lifetime requires ssl_check_enabled to be true (the lifetime assertion is silently dropped when the SSL check is off)

### spec.validationRules.expectedStatusCode

`int32` · optional (explicit presence)

The HTTP status code that counts as success. Unspecified applies
Azure's default (200).

### spec.validationRules.sslCertRemainingLifetime

`int32` · optional (explicit presence)

Fail the run if the endpoint's SSL certificate has fewer than this many
days remaining (1-365). Requires ssl_check_enabled to be true -- the
lifetime assertion rides the SSL check, and Azure silently drops it
when the check is off. Leave unset to skip the lifetime assertion.

- rule: ssl_cert_remaining_lifetime must be between 1 and 365 days

### spec.validationRules.sslCheckEnabled

`bool` · optional (explicit presence)

Whether to validate the SSL certificate at all. Unspecified applies
Azure's default (false). Requires the request url to be https.

### spec.validationRules.content

`AzureApplicationInsightsStandardWebTestContentValidation`

Assert the response body contains (or does not contain) specific text.

### spec.validationRules.content.contentMatch

`string` · required

The text to look for in the response body. Required.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.validationRules.content.ignoreCase

`bool` · optional (explicit presence)

Whether the match is case-insensitive. Unspecified applies Azure's
default (false, case-sensitive).

### spec.validationRules.content.passIfTextFound

`bool` · optional (explicit presence)

Whether finding the text is a PASS (true) or a FAILURE (false). Set
false to fail when an error string appears. Unspecified applies Azure's
default (false).

### spec.geoLocations

`[]string` · required

The Azure regions the test runs FROM (e.g. "us-east-1-azr",
"us-west-2-azr"). At least one required; use several so a single
region's network trouble does not read as an outage. These are Azure
web-test location ids, not standard region names.

- rule: {"repeated":{"minItems":"1"}}

### spec.description

`string`

A human-readable description of what the test checks.

### spec.tags

`map<string, string>`

Free-form tags applied to the web test, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag with
the same key wins. Updatable in place.

## Validation Rules

- `web_test_ssl_checks_require_https`: ssl_check_enabled and ssl_cert_remaining_lifetime require the request url to be https

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureApplicationInsightsStandardWebTest, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.web_test_id` | `string` | The Azure Resource Manager ID of the web test. This is the seam an AzureMonitorMetricAlert references (web_test_id) to alert on the test's availability. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Insights/webTests/{name} |
| `status.outputs.web_test_name` | `string` | The name of the web test resource. |
| `status.outputs.synthetic_monitor_id` | `string` | The synthetic monitor id Azure assigns the test -- the id availability metrics and the portal use to identify the monitor. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.applicationInsightsId` | AzureApplicationInsights | `status.outputs.application_insights_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureMonitorMetricAlert | `spec.webTestAvailabilityCriteria.webTestId` | `status.outputs.web_test_id` |

## See Also

- [Overview](../README.md)
