# AwsRestApiGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsRestApiGateway example (hack/dev manifest and refgen
# Example source): a typed-route REST API exercising every coexisting
# arm -- MOCK + HTTP integrations, models, a request validator, a TOKEN
# authorizer, a gateway response, a resource policy, documentation, and
# a stage with method settings. Literal ARNs stand in for composed
# references so the offline `tofu plan` renders every arm. The OpenAPI
# definition source is the XOR alternative (exactly-one with routes)
# and is proven by spec tests, not this manifest.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRestApiGateway
metadata:
  name: orders-api
  id: orders-api
  org: test-org
  env: dev
spec:
  region: us-west-2
  description: Orders REST API
  apiKeySource: HEADER
  binaryMediaTypes:
    - application/octet-stream
  minimumCompressionSize: 1024
  endpointConfiguration:
    type: REGIONAL
    ipAddressType: ipv4
  endpointAccessMode: BASIC
  securityPolicy: TLS_1_2
  policy:
    Version: "2012-10-17"
    Statement:
      - Effect: Allow
        Principal: "*"
        Action: execute-api:Invoke
        Resource: "*"
  models:
    - name: Order
      contentType: application/json
      description: An order payload
      schema: '{"type":"object","properties":{"id":{"type":"string"}}}'
  requestValidators:
    - name: body-and-params
      validateRequestBody: true
      validateRequestParameters: true
  authorizers:
    - name: orders-token
      type: TOKEN
      lambdaInvokeUri:
        value: arn:aws:apigateway:us-west-2:lambda:path/2015-03-31/functions/arn:aws:lambda:us-west-2:123456789012:function:orders-auth/invocations
      identitySource: method.request.header.Authorization
      resultTtlSeconds: 300
  gatewayResponses:
    - responseType: UNAUTHORIZED
      statusCode: "401"
      responseTemplates:
        application/json: '{"message":"unauthorized"}'
  documentation:
    parts:
      - properties: '{"description":"Orders API"}'
        location:
          type: API
    publishedVersion:
      version: "1"
      description: First published snapshot
  stage:
    name: prod
    description: Production stage
    xrayTracingEnabled: true
    methodSettings:
      - methodPath: "*/*"
        metricsEnabled: true
        loggingLevel: ERROR
        throttlingBurstLimit: 100
        throttlingRateLimit: 50
  routes:
    - path: /health
      method: GET
      operationName: getHealth
      integration:
        type: MOCK
        requestTemplates:
          application/json: '{"statusCode": 200}'
      responses:
        - statusCode: "200"
          integrationResponseTemplates:
            application/json: '{"ok": true}'
    - path: /orders/{id}
      method: GET
      authorization: CUSTOM
      authorizerName: orders-token
      apiKeyRequired: true
      requestValidatorName: body-and-params
      requestParameters:
        method.request.path.id: true
      requestModels:
        application/json: Order
      integration:
        type: HTTP
        httpMethod: GET
        uri:
          value: https://orders.internal.example.com/orders/{id}
        connectionType: INTERNET
        passthroughBehavior: WHEN_NO_MATCH
        timeoutMilliseconds: 29000
      responses:
        - statusCode: "200"
          responseModels:
            application/json: Order
          integrationResponseTemplates:
            application/json: "$input.json('$')"
        - statusCode: "404"
          selectionPattern: ".*not found.*"
          integrationResponseTemplates:
            application/json: '{"error":"not found"}'
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.apiKeySource` | `string` |  |  |  |
| `spec.binaryMediaTypes` | `[]string` |  |  |  |
| `spec.minimumCompressionSize` | `int32` |  |  |  |
| `spec.disableExecuteApiEndpoint` | `bool` |  |  |  |
| `spec.endpointConfiguration` | `AwsRestApiGatewayEndpointConfiguration` |  |  |  |
| `spec.endpointConfiguration.type` | `string` |  |  |  |
| `spec.endpointConfiguration.ipAddressType` | `string` |  |  |  |
| `spec.endpointConfiguration.vpcEndpointIds` | `[]string \| valueFrom` |  |  | AwsVpcEndpoint (`status.outputs.vpc_endpoint_id`) |
| `spec.endpointAccessMode` | `string` |  |  |  |
| `spec.securityPolicy` | `string` |  |  |  |
| `spec.policy` | `object` |  |  |  |
| `spec.routes` | `[]AwsRestApiGatewayRoute` |  |  |  |
| `spec.routes[].path` | `string` |  |  |  |
| `spec.routes[].method` | `string` |  |  |  |
| `spec.routes[].authorization` | `string` |  |  |  |
| `spec.routes[].authorizerName` | `string` |  |  |  |
| `spec.routes[].authorizationScopes` | `[]string` |  |  |  |
| `spec.routes[].apiKeyRequired` | `bool` |  |  |  |
| `spec.routes[].operationName` | `string` |  |  |  |
| `spec.routes[].requestParameters` | `map<string, bool>` |  |  |  |
| `spec.routes[].requestModels` | `map<string, string>` |  |  |  |
| `spec.routes[].requestValidatorName` | `string` |  |  |  |
| `spec.routes[].integration` | `AwsRestApiGatewayIntegration` | yes |  |  |
| `spec.routes[].integration.type` | `string` |  |  |  |
| `spec.routes[].integration.uri` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.invoke_arn`) |
| `spec.routes[].integration.httpMethod` | `string` |  |  |  |
| `spec.routes[].integration.credentialsArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.routes[].integration.connectionType` | `string` |  |  |  |
| `spec.routes[].integration.vpcLinkId` | `string \| valueFrom` |  |  | AwsRestApiVpcLink (`status.outputs.vpc_link_id`) |
| `spec.routes[].integration.passthroughBehavior` | `string` |  |  |  |
| `spec.routes[].integration.contentHandling` | `string` |  |  |  |
| `spec.routes[].integration.cacheKeyParameters` | `[]string` |  |  |  |
| `spec.routes[].integration.cacheNamespace` | `string` |  |  |  |
| `spec.routes[].integration.requestParameters` | `map<string, string>` |  |  |  |
| `spec.routes[].integration.requestTemplates` | `map<string, string>` |  |  |  |
| `spec.routes[].integration.timeoutMilliseconds` | `int32` |  |  |  |
| `spec.routes[].integration.responseTransferMode` | `string` |  |  |  |
| `spec.routes[].integration.tlsInsecureSkipVerification` | `bool` |  |  |  |
| `spec.routes[].responses` | `[]AwsRestApiGatewayRouteResponse` |  |  |  |
| `spec.routes[].responses[].statusCode` | `string` |  |  |  |
| `spec.routes[].responses[].responseModels` | `map<string, string>` |  |  |  |
| `spec.routes[].responses[].responseParameters` | `map<string, bool>` |  |  |  |
| `spec.routes[].responses[].selectionPattern` | `string` |  |  |  |
| `spec.routes[].responses[].integrationResponseParameters` | `map<string, string>` |  |  |  |
| `spec.routes[].responses[].integrationResponseTemplates` | `map<string, string>` |  |  |  |
| `spec.routes[].responses[].contentHandling` | `string` |  |  |  |
| `spec.openapi` | `AwsRestApiGatewayOpenApiDefinition` |  |  |  |
| `spec.openapi.body` | `string` | yes |  |  |
| `spec.openapi.failOnWarnings` | `bool` |  |  |  |
| `spec.openapi.parameters` | `map<string, string>` |  |  |  |
| `spec.openapi.mode` | `string` |  |  |  |
| `spec.models` | `[]AwsRestApiGatewayModel` |  |  |  |
| `spec.models[].name` | `string` | yes |  |  |
| `spec.models[].contentType` | `string` | yes |  |  |
| `spec.models[].description` | `string` |  |  |  |
| `spec.models[].schema` | `string` |  |  |  |
| `spec.requestValidators` | `[]AwsRestApiGatewayRequestValidator` |  |  |  |
| `spec.requestValidators[].name` | `string` | yes |  |  |
| `spec.requestValidators[].validateRequestBody` | `bool` |  |  |  |
| `spec.requestValidators[].validateRequestParameters` | `bool` |  |  |  |
| `spec.authorizers` | `[]AwsRestApiGatewayAuthorizer` |  |  |  |
| `spec.authorizers[].name` | `string` | yes |  |  |
| `spec.authorizers[].type` | `string` |  |  |  |
| `spec.authorizers[].lambdaInvokeUri` | `string \| valueFrom` |  |  | AwsLambda (`status.outputs.invoke_arn`) |
| `spec.authorizers[].credentialsArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.authorizers[].providerArns` | `[]string \| valueFrom` |  |  | AwsCognitoUserPool (`status.outputs.user_pool_arn`) |
| `spec.authorizers[].identitySource` | `string` |  |  |  |
| `spec.authorizers[].identityValidationExpression` | `string` |  |  |  |
| `spec.authorizers[].resultTtlSeconds` | `int32` |  |  |  |
| `spec.gatewayResponses` | `[]AwsRestApiGatewayGatewayResponse` |  |  |  |
| `spec.gatewayResponses[].responseType` | `string` |  |  |  |
| `spec.gatewayResponses[].statusCode` | `string` |  |  |  |
| `spec.gatewayResponses[].responseParameters` | `map<string, string>` |  |  |  |
| `spec.gatewayResponses[].responseTemplates` | `map<string, string>` |  |  |  |
| `spec.stage` | `AwsRestApiGatewayStage` |  |  |  |
| `spec.stage.name` | `string` |  |  |  |
| `spec.stage.description` | `string` |  |  |  |
| `spec.stage.stageVariables` | `map<string, string>` |  |  |  |
| `spec.stage.xrayTracingEnabled` | `bool` |  |  |  |
| `spec.stage.cacheCluster` | `AwsRestApiGatewayCacheCluster` |  |  |  |
| `spec.stage.cacheCluster.enabled` | `bool` |  |  |  |
| `spec.stage.cacheCluster.size` | `string` |  |  |  |
| `spec.stage.clientCertificate` | `AwsRestApiGatewayClientCertificate` |  |  |  |
| `spec.stage.clientCertificate.generate` | `bool` |  |  |  |
| `spec.stage.clientCertificate.existingCertificateId` | `string` |  |  |  |
| `spec.stage.clientCertificate.description` | `string` |  |  |  |
| `spec.stage.accessLog` | `AwsRestApiGatewayAccessLog` |  |  |  |
| `spec.stage.accessLog.destinationArn` | `string \| valueFrom` | yes |  | AwsCloudwatchLogGroup (`status.outputs.log_group_arn`) |
| `spec.stage.accessLog.format` | `string` | yes |  |  |
| `spec.stage.documentationVersion` | `string` |  |  |  |
| `spec.stage.methodSettings` | `[]AwsRestApiGatewayMethodSettings` |  |  |  |
| `spec.stage.methodSettings[].methodPath` | `string` | yes |  |  |
| `spec.stage.methodSettings[].metricsEnabled` | `bool` |  |  |  |
| `spec.stage.methodSettings[].loggingLevel` | `string` |  |  |  |
| `spec.stage.methodSettings[].dataTraceEnabled` | `bool` |  |  |  |
| `spec.stage.methodSettings[].throttlingBurstLimit` | `int32` |  |  |  |
| `spec.stage.methodSettings[].throttlingRateLimit` | `double` |  |  |  |
| `spec.stage.methodSettings[].cachingEnabled` | `bool` |  |  |  |
| `spec.stage.methodSettings[].cacheTtlInSeconds` | `int32` |  |  |  |
| `spec.stage.methodSettings[].cacheDataEncrypted` | `bool` |  |  |  |
| `spec.stage.methodSettings[].requireAuthorizationForCacheControl` | `bool` |  |  |  |
| `spec.stage.methodSettings[].unauthorizedCacheControlHeaderStrategy` | `string` |  |  |  |
| `spec.documentation` | `AwsRestApiGatewayDocumentation` |  |  |  |
| `spec.documentation.parts` | `[]AwsRestApiGatewayDocumentationPart` |  |  |  |
| `spec.documentation.parts[].location` | `AwsRestApiGatewayDocumentationLocation` | yes |  |  |
| `spec.documentation.parts[].location.type` | `string` |  |  |  |
| `spec.documentation.parts[].location.path` | `string` |  |  |  |
| `spec.documentation.parts[].location.method` | `string` |  |  |  |
| `spec.documentation.parts[].location.name` | `string` |  |  |  |
| `spec.documentation.parts[].location.statusCode` | `string` |  |  |  |
| `spec.documentation.parts[].properties` | `string` | yes |  |  |
| `spec.documentation.publishedVersion` | `AwsRestApiGatewayDocumentationVersion` |  |  |  |
| `spec.documentation.publishedVersion.version` | `string` | yes |  |  |
| `spec.documentation.publishedVersion.description` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.description

`string`

- rule: {"string":{"maxLen":"1024"}}

### spec.apiKeySource

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["HEADER","AUTHORIZER"]}}

### spec.binaryMediaTypes

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.minimumCompressionSize

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":10485760,"gte":-1}}

### spec.disableExecuteApiEndpoint

`bool`

### spec.endpointConfiguration

`AwsRestApiGatewayEndpointConfiguration`

- rule: PRIVATE endpoints require ip_address_type 'dualstack' (or omit it for the AWS default)
- rule: vpc_endpoint_ids apply only to PRIVATE endpoint types

### spec.endpointConfiguration.type

`string`

- rule: {"string":{"in":["REGIONAL","EDGE","PRIVATE"]}}

### spec.endpointConfiguration.ipAddressType

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ipv4","dualstack"]}}

### spec.endpointConfiguration.vpcEndpointIds

`[]string | valueFrom`

- references: AwsVpcEndpoint (`status.outputs.vpc_endpoint_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsVpcEndpoint, name: <that resource's name>, fieldPath: status.outputs.vpc_endpoint_id}} -- a bare string does not parse

### spec.endpointAccessMode

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["BASIC","STRICT"]}}

### spec.securityPolicy

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["TLS_1_0","TLS_1_2","SecurityPolicy_TLS13_1_3_2025_09","SecurityPolicy_TLS13_1_3_FIPS_2025_09","SecurityPolicy_TLS13_1_2_PFS_PQ_2025_09","SecurityPolicy_TLS13_1_2_FIPS_PQ_2025_09","SecurityPolicy_TLS13_1_2_FIPS_PFS_PQ_2025_09","SecurityPolicy_TLS13_1_2_PQ_2025_09","SecurityPolicy_TLS13_1_2_2021_06","SecurityPolicy_TLS13_2025_EDGE","SecurityPolicy_TLS12_PFS_2025_EDGE","SecurityPolicy_TLS12_2018_EDGE"]}}

### spec.policy

`object`

### spec.routes

`[]AwsRestApiGatewayRoute`

- rule: responses entries must have unique status_code values

### spec.routes[].path

`string`

- rule: {"string":{"pattern":"^/"}}

### spec.routes[].method

`string`

- rule: {"string":{"in":["ANY","DELETE","GET","HEAD","OPTIONS","PATCH","POST","PUT"]}}

### spec.routes[].authorization

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["NONE","AWS_IAM","CUSTOM","COGNITO_USER_POOLS"]}}

### spec.routes[].authorizerName

`string`

### spec.routes[].authorizationScopes

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.routes[].apiKeyRequired

`bool`

### spec.routes[].operationName

`string`

- rule: {"string":{"maxLen":"64"}}

### spec.routes[].requestParameters

`map<string, bool>`

### spec.routes[].requestModels

`map<string, string>`

### spec.routes[].requestValidatorName

`string`

### spec.routes[].integration

`AwsRestApiGatewayIntegration` · required

- rule: {"required":true}
- rule: MOCK integrations take no uri; every other integration type requires one
- rule: AWS_PROXY (Lambda proxy) integrations use http_method 'POST' (or omit it and the modules send POST)
- rule: http_method is required for HTTP, HTTP_PROXY, and AWS integrations (AWS_PROXY defaults to POST; MOCK takes none)
- rule: integrations with connection_type 'VPC_LINK' must set vpc_link_id (the link to route through), and vpc_link_id must not be set otherwise
- rule: integrations with connection_type 'VPC_LINK' must use integration type 'HTTP' or 'HTTP_PROXY' - REST VPC links front NLB-based HTTP backends
- rule: timeout_milliseconds above 300000 requires response_transfer_mode 'STREAM' (BUFFERED responses cap at 300000)

### spec.routes[].integration.type

`string`

- rule: {"string":{"in":["MOCK","HTTP","HTTP_PROXY","AWS","AWS_PROXY"]}}

### spec.routes[].integration.uri

`string | valueFrom`

- references: AwsLambda (`status.outputs.invoke_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.invoke_arn}} -- a bare string does not parse

### spec.routes[].integration.httpMethod

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ANY","DELETE","GET","HEAD","OPTIONS","PATCH","POST","PUT"]}}

### spec.routes[].integration.credentialsArn

`string | valueFrom`

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.routes[].integration.connectionType

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["INTERNET","VPC_LINK"]}}

### spec.routes[].integration.vpcLinkId

`string | valueFrom`

- references: AwsRestApiVpcLink (`status.outputs.vpc_link_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsRestApiVpcLink, name: <that resource's name>, fieldPath: status.outputs.vpc_link_id}} -- a bare string does not parse

### spec.routes[].integration.passthroughBehavior

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["WHEN_NO_MATCH","WHEN_NO_TEMPLATES","NEVER"]}}

### spec.routes[].integration.contentHandling

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["CONVERT_TO_BINARY","CONVERT_TO_TEXT"]}}

### spec.routes[].integration.cacheKeyParameters

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.routes[].integration.cacheNamespace

`string`

### spec.routes[].integration.requestParameters

`map<string, string>`

### spec.routes[].integration.requestTemplates

`map<string, string>`

### spec.routes[].integration.timeoutMilliseconds

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":900000,"gte":50}}

### spec.routes[].integration.responseTransferMode

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["BUFFERED","STREAM"]}}

### spec.routes[].integration.tlsInsecureSkipVerification

`bool`

### spec.routes[].responses

`[]AwsRestApiGatewayRouteResponse`

### spec.routes[].responses[].statusCode

`string`

- rule: {"string":{"pattern":"^[1-5][0-9][0-9]$"}}

### spec.routes[].responses[].responseModels

`map<string, string>`

### spec.routes[].responses[].responseParameters

`map<string, bool>`

### spec.routes[].responses[].selectionPattern

`string`

### spec.routes[].responses[].integrationResponseParameters

`map<string, string>`

### spec.routes[].responses[].integrationResponseTemplates

`map<string, string>`

### spec.routes[].responses[].contentHandling

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["CONVERT_TO_BINARY","CONVERT_TO_TEXT"]}}

### spec.openapi

`AwsRestApiGatewayOpenApiDefinition`

### spec.openapi.body

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.openapi.failOnWarnings

`bool`

### spec.openapi.parameters

`map<string, string>`

### spec.openapi.mode

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["overwrite","merge"]}}

### spec.models

`[]AwsRestApiGatewayModel`

### spec.models[].name

`string` · required

- rule: {"string":{"minLen":"1","pattern":"^[a-zA-Z0-9]+$"}}

### spec.models[].contentType

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.models[].description

`string`

- rule: {"string":{"maxLen":"1024"}}

### spec.models[].schema

`string`

### spec.requestValidators

`[]AwsRestApiGatewayRequestValidator`

### spec.requestValidators[].name

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.requestValidators[].validateRequestBody

`bool`

### spec.requestValidators[].validateRequestParameters

`bool`

### spec.authorizers

`[]AwsRestApiGatewayAuthorizer`

- rule: TOKEN and REQUEST authorizers require lambda_invoke_uri (and take no provider_arns); COGNITO_USER_POOLS authorizers require provider_arns (and take no lambda_invoke_uri)
- rule: identity_validation_expression applies only to TOKEN authorizers

### spec.authorizers[].name

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"128"}}

### spec.authorizers[].type

`string`

- rule: {"string":{"in":["TOKEN","REQUEST","COGNITO_USER_POOLS"]}}

### spec.authorizers[].lambdaInvokeUri

`string | valueFrom`

- references: AwsLambda (`status.outputs.invoke_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.invoke_arn}} -- a bare string does not parse

### spec.authorizers[].credentialsArn

`string | valueFrom`

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.authorizers[].providerArns

`[]string | valueFrom`

- references: AwsCognitoUserPool (`status.outputs.user_pool_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCognitoUserPool, name: <that resource's name>, fieldPath: status.outputs.user_pool_arn}} -- a bare string does not parse

### spec.authorizers[].identitySource

`string`

### spec.authorizers[].identityValidationExpression

`string`

### spec.authorizers[].resultTtlSeconds

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":3600,"gte":0}}

### spec.gatewayResponses

`[]AwsRestApiGatewayGatewayResponse`

### spec.gatewayResponses[].responseType

`string`

- rule: {"string":{"in":["DEFAULT_4XX","DEFAULT_5XX","RESOURCE_NOT_FOUND","UNAUTHORIZED","INVALID_API_KEY","ACCESS_DENIED","AUTHORIZER_FAILURE","AUTHORIZER_CONFIGURATION_ERROR","INVALID_SIGNATURE","EXPIRED_TOKEN","MISSING_AUTHENTICATION_TOKEN","INTEGRATION_FAILURE","INTEGRATION_TIMEOUT","API_CONFIGURATION_ERROR","UNSUPPORTED_MEDIA_TYPE","BAD_REQUEST_PARAMETERS","BAD_REQUEST_BODY","REQUEST_TOO_LARGE","THROTTLED","QUOTA_EXCEEDED","WAF_FILTERED"]}}

### spec.gatewayResponses[].statusCode

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^[1-5][0-9][0-9]$"}}

### spec.gatewayResponses[].responseParameters

`map<string, string>`

### spec.gatewayResponses[].responseTemplates

`map<string, string>`

### spec.stage

`AwsRestApiGatewayStage`

- rule: method_settings entries must have unique method_path values

### spec.stage.name

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"128","pattern":"^[a-zA-Z0-9_-]+$"}}

### spec.stage.description

`string`

- rule: {"string":{"maxLen":"1024"}}

### spec.stage.stageVariables

`map<string, string>`

### spec.stage.xrayTracingEnabled

`bool`

### spec.stage.cacheCluster

`AwsRestApiGatewayCacheCluster`

### spec.stage.cacheCluster.enabled

`bool`

### spec.stage.cacheCluster.size

`string`

- rule: {"string":{"in":["0.5","1.6","6.1","13.5","28.4","58.2","118","237"]}}

### spec.stage.clientCertificate

`AwsRestApiGatewayClientCertificate`

- rule: set exactly one of generate or existing_certificate_id
- rule: description applies only to a generated certificate

### spec.stage.clientCertificate.generate

`bool`

### spec.stage.clientCertificate.existingCertificateId

`string`

### spec.stage.clientCertificate.description

`string`

- rule: {"string":{"maxLen":"1024"}}

### spec.stage.accessLog

`AwsRestApiGatewayAccessLog`

### spec.stage.accessLog.destinationArn

`string | valueFrom` · required

- references: AwsCloudwatchLogGroup (`status.outputs.log_group_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCloudwatchLogGroup, name: <that resource's name>, fieldPath: status.outputs.log_group_arn}} -- a bare string does not parse

### spec.stage.accessLog.format

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.stage.documentationVersion

`string`

### spec.stage.methodSettings

`[]AwsRestApiGatewayMethodSettings`

### spec.stage.methodSettings[].methodPath

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.stage.methodSettings[].metricsEnabled

`bool` · optional (explicit presence)

### spec.stage.methodSettings[].loggingLevel

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["OFF","ERROR","INFO"]}}

### spec.stage.methodSettings[].dataTraceEnabled

`bool` · optional (explicit presence)

### spec.stage.methodSettings[].throttlingBurstLimit

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":0}}

### spec.stage.methodSettings[].throttlingRateLimit

`double` · optional (explicit presence)

- rule: {"double":{"gte":0}}

### spec.stage.methodSettings[].cachingEnabled

`bool` · optional (explicit presence)

### spec.stage.methodSettings[].cacheTtlInSeconds

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":0}}

### spec.stage.methodSettings[].cacheDataEncrypted

`bool` · optional (explicit presence)

### spec.stage.methodSettings[].requireAuthorizationForCacheControl

`bool` · optional (explicit presence)

### spec.stage.methodSettings[].unauthorizedCacheControlHeaderStrategy

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["FAIL_WITH_403","SUCCEED_WITH_RESPONSE_HEADER","SUCCEED_WITHOUT_RESPONSE_HEADER"]}}

### spec.documentation

`AwsRestApiGatewayDocumentation`

### spec.documentation.parts

`[]AwsRestApiGatewayDocumentationPart`

### spec.documentation.parts[].location

`AwsRestApiGatewayDocumentationLocation` · required

- rule: {"required":true}

### spec.documentation.parts[].location.type

`string`

- rule: {"string":{"in":["API","AUTHORIZER","MODEL","RESOURCE","METHOD","PATH_PARAMETER","QUERY_PARAMETER","REQUEST_HEADER","REQUEST_BODY","RESPONSE","RESPONSE_HEADER","RESPONSE_BODY"]}}

### spec.documentation.parts[].location.path

`string`

### spec.documentation.parts[].location.method

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ANY","DELETE","GET","HEAD","OPTIONS","PATCH","POST","PUT"]}}

### spec.documentation.parts[].location.name

`string`

### spec.documentation.parts[].location.statusCode

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^([1-5][0-9][0-9]|\\*)$"}}

### spec.documentation.parts[].properties

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.documentation.publishedVersion

`AwsRestApiGatewayDocumentationVersion`

### spec.documentation.publishedVersion.version

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.documentation.publishedVersion.description

`string`

- rule: {"string":{"maxLen":"1024"}}

## Validation Rules

- `definition_exactly_one_source`: define the API with exactly one of routes (typed) or openapi (imported document)
- `route_path_method_unique`: route path+method pairs must be unique - two routes with the same path and method would conflict in API Gateway
- `route_path_depth_max_five`: route paths support at most five segments - use a greedy-proxy tail ({proxy+}) for deeper hierarchies
- `model_names_unique`: model names must be unique - routes reference models by name
- `request_validator_names_unique`: request validator names must be unique - routes reference validators by name
- `authorizer_names_unique`: authorizer names must be unique - routes reference authorizers by name
- `gateway_response_types_unique`: gateway_responses entries must have unique response_type values
- `authorized_route_requires_authorizer_name`: routes with authorization 'CUSTOM' or 'COGNITO_USER_POOLS' must specify an authorizer_name
- `route_authorizer_name_must_exist`: route authorizer_name must match a defined authorizer name
- `route_authorizer_type_matches`: a route's authorization must match its authorizer: 'CUSTOM' routes reference TOKEN or REQUEST authorizers, 'COGNITO_USER_POOLS' routes reference COGNITO_USER_POOLS authorizers
- `authorization_scopes_cognito_only`: authorization_scopes apply only to routes with authorization 'COGNITO_USER_POOLS'
- `route_validator_name_must_exist`: route request_validator_name must match a defined request validator name
- `route_request_models_must_exist`: route request_models values must reference a defined model name or the AWS built-ins 'Empty'/'Error'
- `route_response_models_must_exist`: route response model values must reference a defined model name or the AWS built-ins 'Empty'/'Error'
- `stage_documentation_version_requires_documentation`: stage.documentation_version requires documentation.published_version with the same version - the stage can only serve a version this spec publishes

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsRestApiGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.rest_api_id` | `string` |  |
| `status.outputs.rest_api_arn` | `string` |  |
| `status.outputs.execution_arn` | `string` |  |
| `status.outputs.root_resource_id` | `string` |  |
| `status.outputs.stage_name` | `string` |  |
| `status.outputs.stage_arn` | `string` |  |
| `status.outputs.invoke_url` | `string` |  |
| `status.outputs.deployment_id` | `string` |  |
| `status.outputs.client_certificate_id` | `string` |  |
| `status.outputs.client_certificate_pem` | `string` |  |
| `status.outputs.resource_ids` | `map<string, string>` |  |
| `status.outputs.authorizer_ids` | `map<string, string>` |  |
| `status.outputs.model_ids` | `map<string, string>` |  |
| `status.outputs.request_validator_ids` | `map<string, string>` |  |
| `status.outputs.documentation_part_ids` | `map<string, string>` |  |
| `status.outputs.route_resource_ids` | `map<string, string>` |  |
| `status.outputs.route_methods` | `map<string, string>` |  |
| `status.outputs.response_resource_ids` | `map<string, string>` |  |
| `status.outputs.response_methods` | `map<string, string>` |  |
| `status.outputs.response_status_codes` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.endpointConfiguration.vpcEndpointIds` | AwsVpcEndpoint | `status.outputs.vpc_endpoint_id` |
| `spec.routes[].integration.uri` | AwsLambda | `status.outputs.invoke_arn` |
| `spec.routes[].integration.credentialsArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.routes[].integration.vpcLinkId` | AwsRestApiVpcLink | `status.outputs.vpc_link_id` |
| `spec.authorizers[].lambdaInvokeUri` | AwsLambda | `status.outputs.invoke_arn` |
| `spec.authorizers[].credentialsArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.authorizers[].providerArns` | AwsCognitoUserPool | `status.outputs.user_pool_arn` |
| `spec.stage.accessLog.destinationArn` | AwsCloudwatchLogGroup | `status.outputs.log_group_arn` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsRestApiDomain | `spec.basePathMappings[].restApiId` | `status.outputs.rest_api_id` |
| AwsRestApiDomain | `spec.basePathMappings[].stageName` | `status.outputs.stage_name` |
| AwsRestApiUsagePlan | `spec.apiStages[].restApiId` | `status.outputs.rest_api_id` |
| AwsRestApiUsagePlan | `spec.apiStages[].stageName` | `status.outputs.stage_name` |

## See Also

- [Overview](../README.md)
