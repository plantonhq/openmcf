package awsroute53healthcheckv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsRoute53HealthCheckSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRoute53HealthCheckSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

var _ = ginkgo.Describe("AwsRoute53HealthCheckSpec validations", func() {
	var spec *AwsRoute53HealthCheckSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: an HTTPS endpoint check.
		spec = &AwsRoute53HealthCheckSpec{
			Region:    "us-west-2",
			CheckType: "HTTPS",
			Fqdn:      "app.example.com",
		}
	})

	// -------------------------------------------------------------------------
	// Happy path per check type
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal HTTPS check", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts a fully tuned HTTP check", func() {
		spec.CheckType = "HTTP"
		spec.IpAddress = "192.0.2.1"
		spec.Port = 8080
		spec.ResourcePath = "/healthz"
		spec.RequestInterval = proto.Int32(10)
		spec.FailureThreshold = proto.Int32(2)
		spec.MeasureLatency = true
		spec.Regions = []string{"us-east-1", "us-west-2", "eu-west-1"}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts a string-match check", func() {
		spec.CheckType = "HTTPS_STR_MATCH"
		spec.ResourcePath = "/status"
		spec.SearchString = "\"healthy\":true"
		spec.EnableSni = proto.Bool(true)
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts a TCP check with a port", func() {
		spec.CheckType = "TCP"
		spec.Fqdn = ""
		spec.IpAddress = "192.0.2.1"
		spec.Port = 5432
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts a calculated check over children", func() {
		spec = &AwsRoute53HealthCheckSpec{
			Region:    "us-west-2",
			CheckType: "CALCULATED",
			ChildHealthChecks: []*foreignkeyv1.StringValueOrRef{
				strRef("abcdef11-2222-3333-4444-555555fedcba"),
				strRef("fedcba55-4444-3333-2222-11fedcbaabcd"),
			},
			ChildHealthThreshold: 1,
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts a CloudWatch metric check", func() {
		spec = &AwsRoute53HealthCheckSpec{
			Region:                       "us-west-2",
			CheckType:                    "CLOUDWATCH_METRIC",
			CloudwatchAlarmName:          "api-5xx-rate",
			CloudwatchAlarmRegion:        "us-west-2",
			InsufficientDataHealthStatus: "LastKnownStatus",
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts a recovery-control check", func() {
		spec = &AwsRoute53HealthCheckSpec{
			Region:            "us-west-2",
			CheckType:         "RECOVERY_CONTROL",
			RoutingControlArn: "arn:aws:route53-recovery-control::123456789012:controlpanel/abc/routingcontrol/def",
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts state shaping (invert + disabled) on any type", func() {
		spec.InvertHealthcheck = true
		spec.Disabled = true
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// check_type
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a missing check type", func() {
		spec.CheckType = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects an unknown check type", func() {
		spec.CheckType = "PING"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Endpoint contracts
	// -------------------------------------------------------------------------

	ginkgo.It("rejects an endpoint check without a target", func() {
		spec.Fqdn = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a TCP check without a port", func() {
		spec.CheckType = "TCP"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a resource path on a TCP check", func() {
		spec.CheckType = "TCP"
		spec.Port = 443
		spec.ResourcePath = "/healthz"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a string-match check without a search string", func() {
		spec.CheckType = "HTTP_STR_MATCH"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a search string on a plain HTTP check", func() {
		spec.CheckType = "HTTP"
		spec.SearchString = "ok"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects an invalid request interval", func() {
		spec.RequestInterval = proto.Int32(20)
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a failure threshold above 10", func() {
		spec.FailureThreshold = proto.Int32(11)
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects fewer than 3 checker regions", func() {
		spec.Regions = []string{"us-east-1", "us-west-2"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects an invalid checker region", func() {
		spec.Regions = []string{"us-east-1", "us-west-2", "eu-central-1"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Non-endpoint types must not carry endpoint fields
	// -------------------------------------------------------------------------

	ginkgo.It("rejects an fqdn on a calculated check", func() {
		spec.CheckType = "CALCULATED"
		spec.ChildHealthChecks = []*foreignkeyv1.StringValueOrRef{strRef("abc")}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects probe tuning on a CloudWatch check", func() {
		spec = &AwsRoute53HealthCheckSpec{
			Region:                "us-west-2",
			CheckType:             "CLOUDWATCH_METRIC",
			CloudwatchAlarmName:   "api-5xx-rate",
			CloudwatchAlarmRegion: "us-west-2",
			MeasureLatency:        true,
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// CALCULATED contracts
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a calculated check without children", func() {
		spec = &AwsRoute53HealthCheckSpec{
			Region:    "us-west-2",
			CheckType: "CALCULATED",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects children on a non-calculated check", func() {
		spec.ChildHealthChecks = []*foreignkeyv1.StringValueOrRef{strRef("abc")}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// CLOUDWATCH_METRIC contracts
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a CloudWatch check missing the alarm region", func() {
		spec = &AwsRoute53HealthCheckSpec{
			Region:              "us-west-2",
			CheckType:           "CLOUDWATCH_METRIC",
			CloudwatchAlarmName: "api-5xx-rate",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects an invalid insufficient-data status", func() {
		spec = &AwsRoute53HealthCheckSpec{
			Region:                       "us-west-2",
			CheckType:                    "CLOUDWATCH_METRIC",
			CloudwatchAlarmName:          "api-5xx-rate",
			CloudwatchAlarmRegion:        "us-west-2",
			InsufficientDataHealthStatus: "Unknown",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// RECOVERY_CONTROL contracts
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a recovery-control check without the routing control ARN", func() {
		spec = &AwsRoute53HealthCheckSpec{
			Region:    "us-west-2",
			CheckType: "RECOVERY_CONTROL",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a routing control ARN on an endpoint check", func() {
		spec.RoutingControlArn = "arn:aws:route53-recovery-control::123456789012:controlpanel/abc/routingcontrol/def"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// region
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a missing region", func() {
		spec.Region = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})
})
