package awshttpapigatewayv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsHttpApiGatewaySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsHttpApiGatewaySpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// helper to create a minimal valid spec with a single $default route.
func minimalValidSpec() *AwsHttpApiGatewaySpec {
	return &AwsHttpApiGatewaySpec{
		Region: "us-west-2",
		Routes: []*AwsHttpApiGatewayRoute{
			{
				RouteKey: "$default",
				Integration: &AwsHttpApiGatewayIntegration{
					IntegrationType: "AWS_PROXY",
					IntegrationUri:  strRef("arn:aws:lambda:us-east-1:123456789012:function:my-func"),
				},
			},
		},
	}
}

var _ = ginkgo.Describe("AwsHttpApiGatewaySpec validations", func() {
	var spec *AwsHttpApiGatewaySpec

	ginkgo.BeforeEach(func() {
		spec = minimalValidSpec()
	})

	// -------------------------------------------------------------------------
	// Happy paths
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal spec with single $default route", func() {
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts multiple routes to different Lambda functions", func() {
		spec.Routes = []*AwsHttpApiGatewayRoute{
			{
				RouteKey: "GET /users",
				Integration: &AwsHttpApiGatewayIntegration{
					IntegrationType: "AWS_PROXY",
					IntegrationUri:  strRef("arn:aws:lambda:us-east-1:123456789012:function:users"),
				},
			},
			{
				RouteKey: "POST /orders",
				Integration: &AwsHttpApiGatewayIntegration{
					IntegrationType: "AWS_PROXY",
					IntegrationUri:  strRef("arn:aws:lambda:us-east-1:123456789012:function:orders"),
				},
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with HTTP_PROXY integration", func() {
		spec.Routes[0].Integration.IntegrationType = "HTTP_PROXY"
		spec.Routes[0].Integration.IntegrationUri = strRef("https://api.example.com")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with CORS configuration", func() {
		spec.CorsConfiguration = &AwsHttpApiGatewayCorsConfig{
			AllowOrigins:     []string{"https://example.com"},
			AllowMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			MaxAgeSeconds:    3600,
			AllowCredentials: true,
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with stage configuration", func() {
		spec.Stage = &AwsHttpApiGatewayStageConfig{
			Name:       "prod",
			AutoDeploy: proto.Bool(true),
			StageVariables: map[string]string{
				"env": "production",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a stage with auto_deploy explicitly disabled", func() {
		spec.Stage = &AwsHttpApiGatewayStageConfig{
			AutoDeploy: proto.Bool(false),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with access logging", func() {
		spec.Stage = &AwsHttpApiGatewayStageConfig{
			AccessLog: &AwsHttpApiGatewayAccessLogConfig{
				DestinationArn: strRef("arn:aws:logs:us-east-1:123456789012:log-group:/aws/apigateway/my-api"),
				Format:         `{"requestId":"$context.requestId","ip":"$context.identity.sourceIp"}`,
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with throttling", func() {
		spec.Stage = &AwsHttpApiGatewayStageConfig{
			DefaultThrottle: &AwsHttpApiGatewayThrottleConfig{
				BurstLimit: 100,
				RateLimit:  50.0,
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with JWT authorizer", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:           "cognito",
				AuthorizerType: "JWT",
				JwtConfiguration: &AwsHttpApiGatewayJwtConfig{
					Issuer:    strRef("https://cognito-idp.us-east-1.amazonaws.com/us-east-1_abc123"),
					Audiences: []*foreignkeyv1.StringValueOrRef{strRef("my-app-client-id")},
				},
				IdentitySources: []string{"$request.header.Authorization"},
			},
		}
		spec.Routes[0].AuthorizationType = "JWT"
		spec.Routes[0].AuthorizerName = "cognito"
		spec.Routes[0].AuthorizationScopes = []string{"read:users"}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with REQUEST authorizer", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:                           "custom-lambda",
				AuthorizerType:                 "REQUEST",
				AuthorizerUri:                  strRef("arn:aws:lambda:us-east-1:123456789012:function:auth"),
				AuthorizerCredentialsArn:       strRef("arn:aws:iam::123456789012:role/api-auth-role"),
				IdentitySources:                []string{"$request.header.Authorization"},
				ResultTtlSeconds:               300,
				EnableSimpleResponses:          true,
				AuthorizerPayloadFormatVersion: "2.0",
			},
		}
		spec.Routes[0].AuthorizationType = "CUSTOM"
		spec.Routes[0].AuthorizerName = "custom-lambda"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with NONE authorization (explicit)", func() {
		spec.Routes[0].AuthorizationType = "NONE"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with AWS_IAM authorization", func() {
		spec.Routes[0].AuthorizationType = "AWS_IAM"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts integration with explicit payload format version 1.0", func() {
		spec.Routes[0].Integration.PayloadFormatVersion = "1.0"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts integration with explicit payload format version 2.0", func() {
		spec.Routes[0].Integration.PayloadFormatVersion = "2.0"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts integration with timeout", func() {
		spec.Routes[0].Integration.TimeoutMilliseconds = 5000
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a fully-configured production spec", func() {
		spec.Description = "Production API for order management"
		spec.CorsConfiguration = &AwsHttpApiGatewayCorsConfig{
			AllowOrigins:     []string{"https://app.example.com"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"X-Request-Id"},
			MaxAgeSeconds:    7200,
			AllowCredentials: true,
		}
		spec.Stage = &AwsHttpApiGatewayStageConfig{
			AutoDeploy: proto.Bool(true),
			AccessLog: &AwsHttpApiGatewayAccessLogConfig{
				DestinationArn: strRef("arn:aws:logs:us-east-1:123456789012:log-group:/aws/apigateway/orders"),
				Format:         `{"requestId":"$context.requestId"}`,
			},
			DefaultThrottle: &AwsHttpApiGatewayThrottleConfig{
				BurstLimit: 500,
				RateLimit:  100.0,
			},
			DetailedMetricsEnabled: true,
			RouteSettings: []*AwsHttpApiGatewayRouteSettings{
				{
					RouteKey:               "GET /orders",
					ThrottlingBurstLimit:   50,
					ThrottlingRateLimit:    25.0,
					DetailedMetricsEnabled: true,
				},
			},
		}
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:           "cognito",
				AuthorizerType: "JWT",
				JwtConfiguration: &AwsHttpApiGatewayJwtConfig{
					Issuer:    strRef("https://cognito-idp.us-east-1.amazonaws.com/us-east-1_abc123"),
					Audiences: []*foreignkeyv1.StringValueOrRef{strRef("orders-client")},
				},
				IdentitySources: []string{"$request.header.Authorization"},
			},
		}
		spec.Routes = []*AwsHttpApiGatewayRoute{
			{
				RouteKey: "GET /orders",
				Integration: &AwsHttpApiGatewayIntegration{
					IntegrationType:      "AWS_PROXY",
					IntegrationUri:       strRef("arn:aws:lambda:us-east-1:123456789012:function:get-orders"),
					PayloadFormatVersion: "2.0",
				},
				AuthorizationType:   "JWT",
				AuthorizerName:      "cognito",
				AuthorizationScopes: []string{"orders:read"},
			},
			{
				RouteKey: "POST /orders",
				Integration: &AwsHttpApiGatewayIntegration{
					IntegrationType:      "AWS_PROXY",
					IntegrationUri:       strRef("arn:aws:lambda:us-east-1:123456789012:function:create-order"),
					PayloadFormatVersion: "2.0",
				},
				AuthorizationType:   "JWT",
				AuthorizerName:      "cognito",
				AuthorizationScopes: []string{"orders:write"},
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Route validation failures
	// -------------------------------------------------------------------------

	ginkgo.It("fails when no routes are provided", func() {
		spec.Routes = nil
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when routes is empty", func() {
		spec.Routes = []*AwsHttpApiGatewayRoute{}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when route_key is empty", func() {
		spec.Routes[0].RouteKey = ""
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when route integration is nil", func() {
		spec.Routes[0].Integration = nil
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when integration_type is empty", func() {
		spec.Routes[0].Integration.IntegrationType = ""
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when integration_uri is nil", func() {
		spec.Routes[0].Integration.IntegrationUri = nil
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL validation failures: integration_type
	// -------------------------------------------------------------------------

	ginkgo.It("fails when integration_type is invalid", func() {
		spec.Routes[0].Integration.IntegrationType = "MOCK"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL validation failures: authorization
	// -------------------------------------------------------------------------

	ginkgo.It("fails when authorization_type is invalid", func() {
		spec.Routes[0].AuthorizationType = "CUSTOM"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when JWT authorization has no authorizer_name", func() {
		spec.Routes[0].AuthorizationType = "JWT"
		spec.Routes[0].AuthorizerName = ""
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when authorizer_name references non-existent authorizer", func() {
		spec.Routes[0].AuthorizationType = "JWT"
		spec.Routes[0].AuthorizerName = "does-not-exist"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL validation failures: authorizer
	// -------------------------------------------------------------------------

	ginkgo.It("fails when authorizer_type is invalid", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:           "bad",
				AuthorizerType: "INVALID",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when JWT authorizer has no jwt_configuration", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:           "jwt-missing-config",
				AuthorizerType: "JWT",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when JWT authorizer has empty issuer", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:           "jwt-empty-issuer",
				AuthorizerType: "JWT",
				JwtConfiguration: &AwsHttpApiGatewayJwtConfig{
					Issuer: nil,
				},
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when REQUEST authorizer has no authorizer_uri", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:           "lambda-missing-uri",
				AuthorizerType: "REQUEST",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL validation failures: payload format version
	// -------------------------------------------------------------------------

	ginkgo.It("fails when payload_format_version is invalid", func() {
		spec.Routes[0].Integration.PayloadFormatVersion = "3.0"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when authorizer_payload_format_version is invalid", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:                           "bad-payload-ver",
				AuthorizerType:                 "REQUEST",
				AuthorizerUri:                  strRef("arn:aws:lambda:us-east-1:123456789012:function:auth"),
				AuthorizerPayloadFormatVersion: "3.0",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL validation failures: range validations
	// -------------------------------------------------------------------------

	ginkgo.It("fails when integration timeout is below minimum", func() {
		spec.Routes[0].Integration.TimeoutMilliseconds = 10
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when integration timeout exceeds maximum", func() {
		spec.Routes[0].Integration.TimeoutMilliseconds = 31000
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when authorizer TTL exceeds 3600", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:             "bad-ttl",
				AuthorizerType:   "REQUEST",
				AuthorizerUri:    strRef("arn:aws:lambda:us-east-1:123456789012:function:auth"),
				ResultTtlSeconds: 3601,
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Field-level validation failures
	// -------------------------------------------------------------------------

	ginkgo.It("fails when CORS max_age_seconds is negative", func() {
		spec.CorsConfiguration = &AwsHttpApiGatewayCorsConfig{
			MaxAgeSeconds: -1,
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when CORS max_age_seconds exceeds 86400", func() {
		spec.CorsConfiguration = &AwsHttpApiGatewayCorsConfig{
			MaxAgeSeconds: 86401,
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when authorizer name is empty", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:           "",
				AuthorizerType: "JWT",
				JwtConfiguration: &AwsHttpApiGatewayJwtConfig{
					Issuer: strRef("https://example.com"),
				},
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when authorizer name exceeds 128 characters", func() {
		longName := ""
		for i := 0; i < 130; i++ {
			longName += "a"
		}
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:           longName,
				AuthorizerType: "JWT",
				JwtConfiguration: &AwsHttpApiGatewayJwtConfig{
					Issuer: strRef("https://example.com"),
				},
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when access log format is empty", func() {
		spec.Stage = &AwsHttpApiGatewayStageConfig{
			AccessLog: &AwsHttpApiGatewayAccessLogConfig{
				DestinationArn: strRef("arn:aws:logs:us-east-1:123456789012:log-group:test"),
				Format:         "",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when access log destination_arn is missing", func() {
		spec.Stage = &AwsHttpApiGatewayStageConfig{
			AccessLog: &AwsHttpApiGatewayAccessLogConfig{
				Format: `{"requestId":"$context.requestId"}`,
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when description exceeds 1024 characters", func() {
		longDesc := ""
		for i := 0; i < 1025; i++ {
			longDesc += "a"
		}
		spec.Description = longDesc
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// AWS service integrations (integration_subtype)
	// -------------------------------------------------------------------------

	ginkgo.It("accepts an SQS service integration (subtype + credentials + request_parameters, no URI)", func() {
		spec.Routes[0].Integration = &AwsHttpApiGatewayIntegration{
			IntegrationType:    "AWS_PROXY",
			IntegrationSubtype: "SQS-SendMessage",
			CredentialsArn:     strRef("arn:aws:iam::123456789012:role/apigw-sqs-role"),
			RequestParameters: map[string]string{
				"QueueUrl":    "https://sqs.us-east-1.amazonaws.com/123456789012/orders",
				"MessageBody": "$request.body",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a Step Functions service integration", func() {
		spec.Routes[0].Integration = &AwsHttpApiGatewayIntegration{
			IntegrationType:    "AWS_PROXY",
			IntegrationSubtype: "StepFunctions-StartExecution",
			CredentialsArn:     strRef("arn:aws:iam::123456789012:role/apigw-sfn-role"),
			RequestParameters: map[string]string{
				"StateMachineArn": "arn:aws:states:us-east-1:123456789012:stateMachine:orders",
				"Input":           "$request.body",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("fails when a service integration also sets integration_uri", func() {
		spec.Routes[0].Integration = &AwsHttpApiGatewayIntegration{
			IntegrationType:    "AWS_PROXY",
			IntegrationSubtype: "SQS-SendMessage",
			IntegrationUri:     strRef("arn:aws:lambda:us-east-1:123456789012:function:x"),
			CredentialsArn:     strRef("arn:aws:iam::123456789012:role/apigw-sqs-role"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a service integration omits credentials_arn", func() {
		spec.Routes[0].Integration = &AwsHttpApiGatewayIntegration{
			IntegrationType:    "AWS_PROXY",
			IntegrationSubtype: "SQS-SendMessage",
			RequestParameters:  map[string]string{"QueueUrl": "https://sqs.us-east-1.amazonaws.com/123456789012/q"},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a service integration uses HTTP_PROXY", func() {
		spec.Routes[0].Integration = &AwsHttpApiGatewayIntegration{
			IntegrationType:    "HTTP_PROXY",
			IntegrationSubtype: "SQS-SendMessage",
			CredentialsArn:     strRef("arn:aws:iam::123456789012:role/apigw-sqs-role"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a proxy integration omits integration_uri", func() {
		spec.Routes[0].Integration = &AwsHttpApiGatewayIntegration{
			IntegrationType: "AWS_PROXY",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Private integrations (VPC link)
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a private HTTP_PROXY integration through a VPC link", func() {
		spec.Routes[0].Integration = &AwsHttpApiGatewayIntegration{
			IntegrationType:       "HTTP_PROXY",
			IntegrationUri:        strRef("arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/internal/50dc6c495c0c9188/f2f7dc8efc522ab2"),
			ConnectionType:        "VPC_LINK",
			ConnectionId:          strRef("vpclink-abc123"),
			TlsServerNameToVerify: "api.internal.example.com",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("fails when VPC_LINK is set without connection_id", func() {
		spec.Routes[0].Integration = &AwsHttpApiGatewayIntegration{
			IntegrationType: "HTTP_PROXY",
			IntegrationUri:  strRef("https://internal.example.com"),
			ConnectionType:  "VPC_LINK",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when connection_id is set without VPC_LINK", func() {
		spec.Routes[0].Integration = &AwsHttpApiGatewayIntegration{
			IntegrationType: "HTTP_PROXY",
			IntegrationUri:  strRef("https://internal.example.com"),
			ConnectionId:    strRef("vpclink-abc123"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a VPC_LINK integration uses AWS_PROXY", func() {
		spec.Routes[0].Integration = &AwsHttpApiGatewayIntegration{
			IntegrationType: "AWS_PROXY",
			IntegrationUri:  strRef("arn:aws:lambda:us-east-1:123456789012:function:x"),
			ConnectionType:  "VPC_LINK",
			ConnectionId:    strRef("vpclink-abc123"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when connection_type is invalid", func() {
		spec.Routes[0].Integration.ConnectionType = "PRIVATE"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Response parameters
	// -------------------------------------------------------------------------

	ginkgo.It("accepts response parameter mappings", func() {
		spec.Routes[0].Integration.ResponseParameters = []*AwsHttpApiGatewayResponseParameters{
			{
				StatusCode: "500",
				Mappings: map[string]string{
					"overwrite:statuscode":    "503",
					"append:header.x-request": "$context.requestId",
				},
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("fails when response parameter status_code is not a status code", func() {
		spec.Routes[0].Integration.ResponseParameters = []*AwsHttpApiGatewayResponseParameters{
			{StatusCode: "5xx", Mappings: map[string]string{"overwrite:statuscode": "503"}},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when response parameter mappings are empty", func() {
		spec.Routes[0].Integration.ResponseParameters = []*AwsHttpApiGatewayResponseParameters{
			{StatusCode: "500"},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// API-level fields
	// -------------------------------------------------------------------------

	ginkgo.It("accepts ip_address_type dualstack and an api_version", func() {
		spec.IpAddressType = "dualstack"
		spec.ApiVersion = "2026-07-07"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("fails when ip_address_type is invalid", func() {
		spec.IpAddressType = "ipv6"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when route keys are duplicated", func() {
		spec.Routes = append(spec.Routes, &AwsHttpApiGatewayRoute{
			RouteKey: "$default",
			Integration: &AwsHttpApiGatewayIntegration{
				IntegrationType: "HTTP_PROXY",
				IntegrationUri:  strRef("https://example.com"),
			},
		})
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts an operation_name on a route", func() {
		spec.Routes[0].OperationName = "handleDefault"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("fails when operation_name exceeds 64 characters", func() {
		longName := ""
		for i := 0; i < 65; i++ {
			longName += "a"
		}
		spec.Routes[0].OperationName = longName
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Authorizer cross-checks
	// -------------------------------------------------------------------------

	ginkgo.It("fails when authorizer names are duplicated", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:           "auth",
				AuthorizerType: "JWT",
				JwtConfiguration: &AwsHttpApiGatewayJwtConfig{
					Issuer:    strRef("https://example.com"),
					Audiences: []*foreignkeyv1.StringValueOrRef{strRef("a")},
				},
			},
			{
				Name:           "auth",
				AuthorizerType: "REQUEST",
				AuthorizerUri:  strRef("arn:aws:lambda:us-east-1:123456789012:function:auth"),
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a JWT authorizer has no audiences", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:           "jwt-no-aud",
				AuthorizerType: "JWT",
				JwtConfiguration: &AwsHttpApiGatewayJwtConfig{
					Issuer: strRef("https://example.com"),
				},
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a JWT route references a REQUEST authorizer", func() {
		spec.Authorizers = []*AwsHttpApiGatewayAuthorizer{
			{
				Name:           "lambda-auth",
				AuthorizerType: "REQUEST",
				AuthorizerUri:  strRef("arn:aws:lambda:us-east-1:123456789012:function:auth"),
			},
		}
		spec.Routes[0].AuthorizationType = "JWT"
		spec.Routes[0].AuthorizerName = "lambda-auth"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a CUSTOM route has no authorizer_name", func() {
		spec.Routes[0].AuthorizationType = "CUSTOM"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Stage route settings
	// -------------------------------------------------------------------------

	ginkgo.It("fails when route_settings target an unknown route key", func() {
		spec.Stage = &AwsHttpApiGatewayStageConfig{
			RouteSettings: []*AwsHttpApiGatewayRouteSettings{
				{RouteKey: "GET /nonexistent"},
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a route setting has an empty route_key", func() {
		spec.Stage = &AwsHttpApiGatewayStageConfig{
			RouteSettings: []*AwsHttpApiGatewayRouteSettings{
				{RouteKey: ""},
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})
})
