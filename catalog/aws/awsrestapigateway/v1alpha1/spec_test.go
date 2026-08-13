package awsrestapigatewayv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsRestApiGatewaySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRestApiGatewaySpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func i32(v int32) *int32 { return &v }

func boolPtr(v bool) *bool { return &v }

const lambdaInvokeArn = "arn:aws:apigateway:us-west-2:lambda:path/2015-03-31/functions/arn:aws:lambda:us-west-2:123456789012:function:orders/invocations"

// mockRoute is the smallest valid route: a MOCK health check.
func mockRoute() *AwsRestApiGatewayRoute {
	return &AwsRestApiGatewayRoute{
		Path:   "/health",
		Method: "GET",
		Integration: &AwsRestApiGatewayIntegration{
			Type:             "MOCK",
			RequestTemplates: map[string]string{"application/json": `{"statusCode": 200}`},
		},
		Responses: []*AwsRestApiGatewayRouteResponse{
			{
				StatusCode:                   "200",
				IntegrationResponseTemplates: map[string]string{"application/json": `{"ok": true}`},
			},
		},
	}
}

// minimalGateway is the smallest valid API: one MOCK route.
func minimalGateway() *AwsRestApiGatewaySpec {
	return &AwsRestApiGatewaySpec{
		Region: "us-west-2",
		Routes: []*AwsRestApiGatewayRoute{mockRoute()},
	}
}

// lambdaRoute exercises the AWS_PROXY arm.
func lambdaRoute() *AwsRestApiGatewayRoute {
	return &AwsRestApiGatewayRoute{
		Path:   "/orders/{id}",
		Method: "GET",
		Integration: &AwsRestApiGatewayIntegration{
			Type: "AWS_PROXY",
			Uri:  svr(lambdaInvokeArn),
		},
	}
}

var _ = ginkgo.Describe("AwsRestApiGatewaySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with a minimal MOCK route", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalGateway())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an OpenAPI definition instead of routes", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsRestApiGatewaySpec{
					Region: "us-west-2",
					Openapi: &AwsRestApiGatewayOpenApiDefinition{
						Body: `{"openapi":"3.0.1","info":{"title":"orders","version":"1"},"paths":{}}`,
						Mode: "overwrite",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full surface configured", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGateway()
				spec.Description = "orders REST API"
				spec.ApiKeySource = "HEADER"
				spec.BinaryMediaTypes = []string{"application/octet-stream"}
				spec.MinimumCompressionSize = i32(1024)
				spec.DisableExecuteApiEndpoint = true
				spec.EndpointConfiguration = &AwsRestApiGatewayEndpointConfiguration{
					Type:          "REGIONAL",
					IpAddressType: "dualstack",
				}
				spec.EndpointAccessMode = "STRICT"
				spec.SecurityPolicy = "SecurityPolicy_TLS13_1_2_2021_06"
				spec.Models = []*AwsRestApiGatewayModel{
					{
						Name:        "OrderInput",
						ContentType: "application/json",
						Schema:      `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object"}`,
					},
				}
				spec.RequestValidators = []*AwsRestApiGatewayRequestValidator{
					{Name: "body-only", ValidateRequestBody: true},
				}
				spec.Authorizers = []*AwsRestApiGatewayAuthorizer{
					{
						Name:                         "jwt-lambda",
						Type:                         "TOKEN",
						LambdaInvokeUri:              svr(lambdaInvokeArn),
						ResultTtlSeconds:             i32(0),
						IdentityValidationExpression: "^Bearer .+$",
					},
					{
						Name:         "cognito",
						Type:         "COGNITO_USER_POOLS",
						ProviderArns: []*foreignkeyv1.StringValueOrRef{svr("arn:aws:cognito-idp:us-west-2:123456789012:userpool/us-west-2_abc")},
					},
				}
				spec.GatewayResponses = []*AwsRestApiGatewayGatewayResponse{
					{
						ResponseType:       "MISSING_AUTHENTICATION_TOKEN",
						StatusCode:         "404",
						ResponseParameters: map[string]string{"gatewayresponse.header.x-request-id": "'$context.requestId'"},
					},
				}
				orders := lambdaRoute()
				orders.Authorization = "CUSTOM"
				orders.AuthorizerName = "jwt-lambda"
				orders.RequestValidatorName = "body-only"
				orders.RequestModels = map[string]string{"application/json": "OrderInput"}
				orders.ApiKeyRequired = true
				me := &AwsRestApiGatewayRoute{
					Path:                "/me",
					Method:              "GET",
					Authorization:       "COGNITO_USER_POOLS",
					AuthorizerName:      "cognito",
					AuthorizationScopes: []string{"orders/read"},
					Integration: &AwsRestApiGatewayIntegration{
						Type:       "HTTP_PROXY",
						Uri:        svr("https://backend.example.com/me"),
						HttpMethod: "GET",
					},
				}
				private := &AwsRestApiGatewayRoute{
					Path:   "/internal",
					Method: "ANY",
					Integration: &AwsRestApiGatewayIntegration{
						Type:           "HTTP_PROXY",
						Uri:            svr("https://internal.example.com"),
						HttpMethod:     "ANY",
						ConnectionType: "VPC_LINK",
						VpcLinkId:      svr("vl-abc123"),
					},
				}
				spec.Routes = append(spec.Routes, orders, me, private)
				spec.Stage = &AwsRestApiGatewayStage{
					Name:                 "prod",
					XrayTracingEnabled:   true,
					StageVariables:       map[string]string{"backendHost": "backend.example.com"},
					CacheCluster:         &AwsRestApiGatewayCacheCluster{Enabled: true, Size: "0.5"},
					ClientCertificate:    &AwsRestApiGatewayClientCertificate{Generate: true, Description: "backend trust"},
					DocumentationVersion: "1.0.0",
					AccessLog: &AwsRestApiGatewayAccessLog{
						DestinationArn: svr("arn:aws:logs:us-west-2:123456789012:log-group:/apigw/orders"),
						Format:         `{"requestId":"$context.requestId"}`,
					},
					MethodSettings: []*AwsRestApiGatewayMethodSettings{
						{
							MethodPath:     "*/*",
							MetricsEnabled: boolPtr(true),
							LoggingLevel:   "ERROR",
						},
						{
							MethodPath:           "orders/{id}/GET",
							CachingEnabled:       boolPtr(true),
							CacheTtlInSeconds:    i32(60),
							ThrottlingBurstLimit: i32(100),
							ThrottlingRateLimit:  func() *float64 { v := 50.0; return &v }(),
						},
					},
				}
				spec.Documentation = &AwsRestApiGatewayDocumentation{
					Parts: []*AwsRestApiGatewayDocumentationPart{
						{
							Location:   &AwsRestApiGatewayDocumentationLocation{Type: "API"},
							Properties: `{"description":"The orders API"}`,
						},
					},
					PublishedVersion: &AwsRestApiGatewayDocumentationVersion{Version: "1.0.0"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a spec with neither routes nor openapi", func() {
			spec := &AwsRestApiGatewaySpec{Region: "us-west-2"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of routes"))
		})

		ginkgo.It("rejects a spec with both routes and openapi", func() {
			spec := minimalGateway()
			spec.Openapi = &AwsRestApiGatewayOpenApiDefinition{Body: "{}"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of routes"))
		})

		ginkgo.It("rejects duplicate path+method pairs", func() {
			spec := minimalGateway()
			spec.Routes = append(spec.Routes, mockRoute())
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("path+method pairs must be unique"))
		})

		ginkgo.It("rejects a path not starting with /", func() {
			spec := minimalGateway()
			spec.Routes[0].Path = "health"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a MOCK integration carrying a uri", func() {
			spec := minimalGateway()
			spec.Routes[0].Integration.Uri = svr("https://example.com")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("MOCK integrations take no uri"))
		})

		ginkgo.It("rejects an HTTP integration without a uri", func() {
			spec := minimalGateway()
			spec.Routes[0].Integration = &AwsRestApiGatewayIntegration{
				Type:       "HTTP_PROXY",
				HttpMethod: "GET",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("MOCK integrations take no uri"))
		})

		ginkgo.It("rejects an HTTP integration without a backend method", func() {
			spec := minimalGateway()
			spec.Routes[0].Integration = &AwsRestApiGatewayIntegration{
				Type: "HTTP_PROXY",
				Uri:  svr("https://example.com"),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("http_method is required"))
		})

		ginkgo.It("rejects an AWS_PROXY integration with a non-POST method", func() {
			spec := minimalGateway()
			route := lambdaRoute()
			route.Integration.HttpMethod = "GET"
			spec.Routes = []*AwsRestApiGatewayRoute{route}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("use http_method 'POST'"))
		})

		ginkgo.It("rejects VPC_LINK without a vpc_link_id", func() {
			spec := minimalGateway()
			spec.Routes[0].Integration = &AwsRestApiGatewayIntegration{
				Type:           "HTTP_PROXY",
				Uri:            svr("https://internal.example.com"),
				HttpMethod:     "GET",
				ConnectionType: "VPC_LINK",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("must set vpc_link_id"))
		})

		ginkgo.It("rejects a vpc_link_id without VPC_LINK", func() {
			spec := minimalGateway()
			spec.Routes[0].Integration = &AwsRestApiGatewayIntegration{
				Type:       "HTTP_PROXY",
				Uri:        svr("https://example.com"),
				HttpMethod: "GET",
				VpcLinkId:  svr("vl-abc"),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("must set vpc_link_id"))
		})

		ginkgo.It("rejects a VPC_LINK integration that is not HTTP", func() {
			spec := minimalGateway()
			route := lambdaRoute()
			route.Integration.ConnectionType = "VPC_LINK"
			route.Integration.VpcLinkId = svr("vl-abc")
			spec.Routes = []*AwsRestApiGatewayRoute{route}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("REST VPC links front NLB"))
		})

		ginkgo.It("rejects a BUFFERED timeout above 300 seconds", func() {
			spec := minimalGateway()
			route := lambdaRoute()
			route.Integration.TimeoutMilliseconds = 400000
			spec.Routes = []*AwsRestApiGatewayRoute{route}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("requires response_transfer_mode 'STREAM'"))
		})

		ginkgo.It("accepts a STREAM timeout up to 900 seconds", func() {
			spec := minimalGateway()
			route := lambdaRoute()
			route.Integration.TimeoutMilliseconds = 900000
			route.Integration.ResponseTransferMode = "STREAM"
			spec.Routes = []*AwsRestApiGatewayRoute{route}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects a CUSTOM route without an authorizer name", func() {
			spec := minimalGateway()
			spec.Routes[0].Authorization = "CUSTOM"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("must specify an authorizer_name"))
		})

		ginkgo.It("rejects an authorizer reference that does not resolve", func() {
			spec := minimalGateway()
			spec.Routes[0].Authorization = "CUSTOM"
			spec.Routes[0].AuthorizerName = "ghost"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("must match a defined authorizer name"))
		})

		ginkgo.It("rejects a COGNITO route bound to a Lambda authorizer", func() {
			spec := minimalGateway()
			spec.Authorizers = []*AwsRestApiGatewayAuthorizer{
				{Name: "lambda-auth", Type: "TOKEN", LambdaInvokeUri: svr(lambdaInvokeArn)},
			}
			spec.Routes[0].Authorization = "COGNITO_USER_POOLS"
			spec.Routes[0].AuthorizerName = "lambda-auth"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("must match its authorizer"))
		})

		ginkgo.It("rejects authorization scopes on a non-Cognito route", func() {
			spec := minimalGateway()
			spec.Routes[0].AuthorizationScopes = []string{"orders/read"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("authorization_scopes apply only"))
		})

		ginkgo.It("rejects a TOKEN authorizer without a Lambda", func() {
			spec := minimalGateway()
			spec.Authorizers = []*AwsRestApiGatewayAuthorizer{
				{Name: "bad", Type: "TOKEN"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("COGNITO_USER_POOLS authorizers require provider_arns"))
		})

		ginkgo.It("rejects a COGNITO authorizer without provider pools", func() {
			spec := minimalGateway()
			spec.Authorizers = []*AwsRestApiGatewayAuthorizer{
				{Name: "bad", Type: "COGNITO_USER_POOLS"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("COGNITO_USER_POOLS authorizers require provider_arns"))
		})

		ginkgo.It("rejects an identity validation expression on a REQUEST authorizer", func() {
			spec := minimalGateway()
			spec.Authorizers = []*AwsRestApiGatewayAuthorizer{
				{
					Name:                         "req",
					Type:                         "REQUEST",
					LambdaInvokeUri:              svr(lambdaInvokeArn),
					IdentityValidationExpression: "^Bearer .+$",
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("applies only to TOKEN authorizers"))
		})

		ginkgo.It("rejects a request model referencing an undefined model", func() {
			spec := minimalGateway()
			spec.Routes[0].RequestModels = map[string]string{"application/json": "Ghost"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("must reference a defined model name"))
		})

		ginkgo.It("accepts the AWS built-in models without defining them", func() {
			spec := minimalGateway()
			spec.Routes[0].RequestModels = map[string]string{"application/json": "Empty"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects a validator reference that does not resolve", func() {
			spec := minimalGateway()
			spec.Routes[0].RequestValidatorName = "ghost"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("must match a defined request validator"))
		})

		ginkgo.It("rejects duplicate response status codes on one route", func() {
			spec := minimalGateway()
			spec.Routes[0].Responses = append(spec.Routes[0].Responses, &AwsRestApiGatewayRouteResponse{StatusCode: "200"})
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("unique status_code values"))
		})

		ginkgo.It("rejects duplicate gateway response types", func() {
			spec := minimalGateway()
			spec.GatewayResponses = []*AwsRestApiGatewayGatewayResponse{
				{ResponseType: "DEFAULT_4XX"},
				{ResponseType: "DEFAULT_4XX"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("unique response_type values"))
		})

		ginkgo.It("rejects a PRIVATE endpoint forced to ipv4", func() {
			spec := minimalGateway()
			spec.EndpointConfiguration = &AwsRestApiGatewayEndpointConfiguration{
				Type:          "PRIVATE",
				IpAddressType: "ipv4",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("PRIVATE endpoints require ip_address_type"))
		})

		ginkgo.It("rejects VPC endpoints on a non-PRIVATE API", func() {
			spec := minimalGateway()
			spec.EndpointConfiguration = &AwsRestApiGatewayEndpointConfiguration{
				Type:           "REGIONAL",
				VpcEndpointIds: []*foreignkeyv1.StringValueOrRef{svr("vpce-abc")},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("vpc_endpoint_ids apply only"))
		})

		ginkgo.It("rejects a client certificate with both arms", func() {
			spec := minimalGateway()
			spec.Stage = &AwsRestApiGatewayStage{
				ClientCertificate: &AwsRestApiGatewayClientCertificate{
					Generate:              true,
					ExistingCertificateId: "cert-abc",
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of generate or existing_certificate_id"))
		})

		ginkgo.It("rejects a description on a referenced certificate", func() {
			spec := minimalGateway()
			spec.Stage = &AwsRestApiGatewayStage{
				ClientCertificate: &AwsRestApiGatewayClientCertificate{
					ExistingCertificateId: "cert-abc",
					Description:           "not ours to describe",
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("description applies only to a generated certificate"))
		})

		ginkgo.It("rejects duplicate method_settings paths", func() {
			spec := minimalGateway()
			spec.Stage = &AwsRestApiGatewayStage{
				MethodSettings: []*AwsRestApiGatewayMethodSettings{
					{MethodPath: "*/*"},
					{MethodPath: "*/*"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("unique method_path values"))
		})

		ginkgo.It("rejects a stage documentation version this spec does not publish", func() {
			spec := minimalGateway()
			spec.Stage = &AwsRestApiGatewayStage{DocumentationVersion: "1.0.0"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("documentation.published_version"))
		})

		ginkgo.It("rejects an invalid cache size", func() {
			spec := minimalGateway()
			spec.Stage = &AwsRestApiGatewayStage{
				CacheCluster: &AwsRestApiGatewayCacheCluster{Enabled: true, Size: "2.0"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an authorizer TTL above one hour", func() {
			spec := minimalGateway()
			spec.Authorizers = []*AwsRestApiGatewayAuthorizer{
				{
					Name:             "jwt",
					Type:             "TOKEN",
					LambdaInvokeUri:  svr(lambdaInvokeArn),
					ResultTtlSeconds: i32(3601),
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
