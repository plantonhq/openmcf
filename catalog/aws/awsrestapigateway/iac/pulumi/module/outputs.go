package module

import (
	awsrestapigatewayv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsrestapigateway/v1alpha1"
)

const (
	OpRestApiId            = "rest_api_id"
	OpRestApiArn           = "rest_api_arn"
	OpExecutionArn         = "execution_arn"
	OpRootResourceId       = "root_resource_id"
	OpStageName            = "stage_name"
	OpStageArn             = "stage_arn"
	OpInvokeUrl            = "invoke_url"
	OpDeploymentId         = "deployment_id"
	OpClientCertificateId  = "client_certificate_id"
	OpClientCertificatePem = "client_certificate_pem"
	OpResourceIds          = "resource_ids"
	OpAuthorizerIds        = "authorizer_ids"
	OpModelIds             = "model_ids"
	OpRequestValidatorIds  = "request_validator_ids"
	OpDocumentationPartIds = "documentation_part_ids"
	OpRouteResourceIds     = "route_resource_ids"
	OpRouteMethods         = "route_methods"
	OpResponseResourceIds  = "response_resource_ids"
	OpResponseMethods      = "response_methods"
	OpResponseStatusCodes  = "response_status_codes"
)

// awsRestApiGatewayStageDefaults stands in when the spec omits the
// stage entirely (name resolution already happened in locals).
var awsRestApiGatewayStageDefaults = awsrestapigatewayv1alpha1.AwsRestApiGatewayStage{}
