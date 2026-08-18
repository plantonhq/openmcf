package awseventbridgeapidestinationv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsEventBridgeApiDestinationSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsEventBridgeApiDestinationSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func apiKeyConnection() *AwsEventBridgeConnection {
	return &AwsEventBridgeConnection{
		Name: "partner-api",
		ApiKey: &AwsEventBridgeConnectionApiKeyAuth{
			Key:   "x-api-key",
			Value: "super-secret-key",
		},
	}
}

func fullSpec() *AwsEventBridgeApiDestinationSpec {
	return &AwsEventBridgeApiDestinationSpec{
		Region:     "us-east-1",
		Connection: apiKeyConnection(),
		Destination: &AwsEventBridgeDestination{
			Name:               "partner-webhook",
			InvocationEndpoint: "https://api.example.com/events",
			HttpMethod:         "POST",
		},
	}
}

var _ = ginkgo.Describe("AwsEventBridgeApiDestinationSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a connection + destination instance", func() {
			gomega.Expect(protovalidate.Validate(fullSpec())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a connection-only instance (shared trust anchor)", func() {
			spec := fullSpec()
			spec.Destination = nil
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a destination-only instance referencing an external connection", func() {
			spec := fullSpec()
			spec.Connection = nil
			spec.Destination.ConnectionArn = svr("arn:aws:events:us-east-1:123456789012:connection/shared/abc")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts basic auth", func() {
			spec := fullSpec()
			spec.Connection.ApiKey = nil
			spec.Connection.Basic = &AwsEventBridgeConnectionBasicAuth{
				Username: "svc-events",
				Password: "hunter2-but-long",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts oauth with token-request parameters", func() {
			spec := fullSpec()
			spec.Connection.ApiKey = nil
			spec.Connection.Oauth = &AwsEventBridgeConnectionOAuth{
				AuthorizationEndpoint: "https://auth.example.com/oauth2/token",
				HttpMethod:            "POST",
				ClientId:              "planton-events",
				ClientSecret:          "oauth-client-secret",
				OauthHttpParameters: &AwsEventBridgeConnectionHttpParameters{
					Body: []*AwsEventBridgeConnectionHttpParameter{
						{Key: "grant_type", Value: "client_credentials"},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects an empty spec", func() {
			spec := &AwsEventBridgeApiDestinationSpec{Region: "us-east-1"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a destination with both an owned connection and an external ARN", func() {
			spec := fullSpec()
			spec.Destination.ConnectionArn = svr("arn:aws:events:us-east-1:123456789012:connection/other/xyz")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a destination-only instance without a connection reference", func() {
			spec := fullSpec()
			spec.Connection = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a connection with two auth modes", func() {
			spec := fullSpec()
			spec.Connection.Basic = &AwsEventBridgeConnectionBasicAuth{Username: "u", Password: "p"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a connection with no auth mode", func() {
			spec := fullSpec()
			spec.Connection.ApiKey = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects oauth without token-request parameters", func() {
			spec := fullSpec()
			spec.Connection.ApiKey = nil
			spec.Connection.Oauth = &AwsEventBridgeConnectionOAuth{
				AuthorizationEndpoint: "https://auth.example.com/oauth2/token",
				HttpMethod:            "POST",
				ClientId:              "planton-events",
				ClientSecret:          "oauth-client-secret",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty http-parameter set", func() {
			spec := fullSpec()
			spec.Connection.InvocationHttpParameters = &AwsEventBridgeConnectionHttpParameters{}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a non-https invocation endpoint", func() {
			spec := fullSpec()
			spec.Destination.InvocationEndpoint = "http://api.example.com/events"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid connection name charset", func() {
			spec := fullSpec()
			spec.Connection.Name = "partner api"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown destination http method", func() {
			spec := fullSpec()
			spec.Destination.HttpMethod = "TRACE"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a zero rate limit", func() {
			spec := fullSpec()
			zero := int32(0)
			spec.Destination.InvocationRateLimitPerSecond = &zero
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
