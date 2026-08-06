package awscognitouserpoolclientv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsCognitoUserPoolClientSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCognitoUserPoolClientSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

var _ = ginkgo.Describe("AwsCognitoUserPoolClientSpec validations", func() {
	var spec *AwsCognitoUserPoolClientSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: region + pool reference.
		spec = &AwsCognitoUserPoolClientSpec{
			Region:     "us-west-2",
			UserPoolId: strRef("us-east-1_Ab1Cd2EfG"),
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal spec (region + pool)", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a web-app OAuth client", func() {
		spec.AllowedOauthFlowsUserPoolClient = true
		spec.AllowedOauthFlows = []string{"code"}
		spec.AllowedOauthScopes = []string{"openid", "email", "profile"}
		spec.CallbackUrls = []string{"https://app.example.com/callback"}
		spec.LogoutUrls = []string{"https://app.example.com/"}
		spec.DefaultRedirectUri = "https://app.example.com/callback"
		spec.ExplicitAuthFlows = []string{"ALLOW_USER_SRP_AUTH", "ALLOW_REFRESH_TOKEN_AUTH"}
		spec.PreventUserExistenceErrors = "ENABLED"
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a machine-to-machine client_credentials client", func() {
		spec.GenerateSecret = true
		spec.AllowedOauthFlowsUserPoolClient = true
		spec.AllowedOauthFlows = []string{"client_credentials"}
		spec.AllowedOauthScopes = []string{"https://api.example.com/read"}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts federated providers by reference and literal", func() {
		spec.SupportedIdentityProviders = []*foreignkeyv1.StringValueOrRef{
			strRef("COGNITO"),
			strRef("Google"),
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts explicit token lifetimes with units", func() {
		spec.AccessTokenValidity = proto.Int32(30)
		spec.IdTokenValidity = proto.Int32(30)
		spec.RefreshTokenValidity = proto.Int32(90)
		spec.TokenValidityUnits = &AwsCognitoUserPoolClientTokenValidityUnits{
			AccessToken:  "minutes",
			IdToken:      "minutes",
			RefreshToken: "days",
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts default-unit lifetimes (hours/hours/days)", func() {
		spec.AccessTokenValidity = proto.Int32(1)
		spec.IdTokenValidity = proto.Int32(24)
		spec.RefreshTokenValidity = proto.Int32(3650)
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts refresh-token rotation with a grace period", func() {
		spec.RefreshTokenRotation = &AwsCognitoUserPoolClientRefreshTokenRotation{
			Feature:                 "ENABLED",
			RetryGracePeriodSeconds: 30,
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts Pinpoint analytics by application ARN", func() {
		spec.AnalyticsConfiguration = &AwsCognitoUserPoolClientAnalyticsConfig{
			ApplicationArn: "arn:aws:mobiletargeting:us-west-2:123456789012:apps/abc123",
			UserDataShared: true,
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts Pinpoint analytics by application ID with role wiring", func() {
		spec.AnalyticsConfiguration = &AwsCognitoUserPoolClientAnalyticsConfig{
			ApplicationId: "abc123",
			ExternalId:    "cognito-analytics",
			RoleArn:       strRef("arn:aws:iam::123456789012:role/pinpoint-publish"),
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Required fields
	// -------------------------------------------------------------------------

	ginkgo.It("fails when region is empty", func() {
		spec.Region = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when user_pool_id is missing", func() {
		spec.UserPoolId = nil
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: OAuth contract
	// -------------------------------------------------------------------------

	ginkgo.It("fails on an invalid OAuth flow", func() {
		spec.AllowedOauthFlowsUserPoolClient = true
		spec.AllowedOauthFlows = []string{"password"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when client_credentials is mixed with code", func() {
		spec.AllowedOauthFlowsUserPoolClient = true
		spec.AllowedOauthFlows = []string{"client_credentials", "code"}
		spec.CallbackUrls = []string{"https://app.example.com/callback"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when OAuth flows are set without enabling OAuth", func() {
		spec.AllowedOauthFlows = []string{"code"}
		spec.CallbackUrls = []string{"https://app.example.com/callback"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when the code grant has no callback URLs", func() {
		spec.AllowedOauthFlowsUserPoolClient = true
		spec.AllowedOauthFlows = []string{"code"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when default_redirect_uri is not a callback URL", func() {
		spec.AllowedOauthFlowsUserPoolClient = true
		spec.AllowedOauthFlows = []string{"code"}
		spec.CallbackUrls = []string{"https://app.example.com/callback"}
		spec.DefaultRedirectUri = "https://other.example.com/callback"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: auth flows
	// -------------------------------------------------------------------------

	ginkgo.It("fails on an invalid explicit auth flow", func() {
		spec.ExplicitAuthFlows = []string{"ALLOW_MAGIC_LINK_AUTH"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when legacy and ALLOW_* auth flows are mixed", func() {
		spec.ExplicitAuthFlows = []string{"USER_PASSWORD_AUTH", "ALLOW_REFRESH_TOKEN_AUTH"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when auth_session_validity is out of range", func() {
		spec.AuthSessionValidity = proto.Int32(20)
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: token lifetimes
	// -------------------------------------------------------------------------

	ginkgo.It("fails when the access token lifetime exceeds 24 hours", func() {
		spec.AccessTokenValidity = proto.Int32(25)
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when the access token lifetime is under 5 minutes", func() {
		spec.AccessTokenValidity = proto.Int32(2)
		spec.TokenValidityUnits = &AwsCognitoUserPoolClientTokenValidityUnits{AccessToken: "minutes"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when the ID token lifetime exceeds 24 hours in minutes", func() {
		spec.IdTokenValidity = proto.Int32(1500)
		spec.TokenValidityUnits = &AwsCognitoUserPoolClientTokenValidityUnits{IdToken: "minutes"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when the refresh token lifetime is under 60 minutes", func() {
		spec.RefreshTokenValidity = proto.Int32(30)
		spec.TokenValidityUnits = &AwsCognitoUserPoolClientTokenValidityUnits{RefreshToken: "minutes"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when the refresh token lifetime exceeds 10 years", func() {
		spec.RefreshTokenValidity = proto.Int32(3651)
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails on an invalid token validity unit", func() {
		spec.TokenValidityUnits = &AwsCognitoUserPoolClientTokenValidityUnits{AccessToken: "weeks"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: rotation, posture, analytics
	// -------------------------------------------------------------------------

	ginkgo.It("fails on an invalid refresh_token_rotation feature", func() {
		spec.RefreshTokenRotation = &AwsCognitoUserPoolClientRefreshTokenRotation{Feature: "ON"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when the rotation grace period exceeds 60 seconds", func() {
		spec.RefreshTokenRotation = &AwsCognitoUserPoolClientRefreshTokenRotation{
			Feature:                 "ENABLED",
			RetryGracePeriodSeconds: 61,
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when rotation is ENABLED alongside ALLOW_REFRESH_TOKEN_AUTH", func() {
		spec.RefreshTokenRotation = &AwsCognitoUserPoolClientRefreshTokenRotation{Feature: "ENABLED"}
		spec.ExplicitAuthFlows = []string{"ALLOW_USER_SRP_AUTH", "ALLOW_REFRESH_TOKEN_AUTH"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts ALLOW_REFRESH_TOKEN_AUTH when rotation is DISABLED", func() {
		spec.RefreshTokenRotation = &AwsCognitoUserPoolClientRefreshTokenRotation{Feature: "DISABLED"}
		spec.ExplicitAuthFlows = []string{"ALLOW_USER_SRP_AUTH", "ALLOW_REFRESH_TOKEN_AUTH"}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("fails on an invalid prevent_user_existence_errors", func() {
		spec.PreventUserExistenceErrors = "DISABLED"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when analytics sets both application ARN and ID", func() {
		spec.AnalyticsConfiguration = &AwsCognitoUserPoolClientAnalyticsConfig{
			ApplicationArn: "arn:aws:mobiletargeting:us-west-2:123456789012:apps/abc123",
			ApplicationId:  "abc123",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when analytics application ID lacks role wiring", func() {
		spec.AnalyticsConfiguration = &AwsCognitoUserPoolClientAnalyticsConfig{
			ApplicationId: "abc123",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})
})
