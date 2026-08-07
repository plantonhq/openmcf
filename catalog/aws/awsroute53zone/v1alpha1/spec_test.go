package awsroute53zonev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsRoute53ZoneSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRoute53ZoneSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

var _ = ginkgo.Describe("AwsRoute53ZoneSpec validations", func() {
	var spec *AwsRoute53ZoneSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: a public zone.
		spec = &AwsRoute53ZoneSpec{
			Region: "us-west-2",
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal public zone", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts a public zone with comment, force_destroy, and accelerated recovery", func() {
		spec.Comment = "production apex zone - platform team"
		spec.ForceDestroy = true
		spec.EnableAcceleratedRecovery = true
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts a public zone on a reusable delegation set", func() {
		spec.DelegationSetId = "N1PA6795SAMPLE"
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts a private zone with VPC associations", func() {
		spec.IsPrivate = true
		spec.VpcAssociations = []*AwsRoute53ZoneVpcAssociation{
			{VpcId: strRef("vpc-0123456789abcdef0")},
			{VpcId: strRef("vpc-0fedcba9876543210"), VpcRegion: "eu-west-1"},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts query logging on a public zone", func() {
		spec.QueryLogging = &AwsRoute53ZoneQueryLogging{
			CloudwatchLogGroupArn: strRef("arn:aws:logs:us-east-1:123456789012:log-group:/aws/route53/example.com"),
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts DNSSEC on a public zone", func() {
		spec.Dnssec = &AwsRoute53ZoneDnssec{
			KmsKeyArn:         strRef("arn:aws:kms:us-east-1:123456789012:key/abc"),
			KeySigningKeyName: "example_com_ksk",
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Private-zone contract
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a private zone without VPC associations", func() {
		spec.IsPrivate = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects VPC associations on a public zone", func() {
		spec.VpcAssociations = []*AwsRoute53ZoneVpcAssociation{
			{VpcId: strRef("vpc-0123456789abcdef0")},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a VPC association without a VPC", func() {
		spec.IsPrivate = true
		spec.VpcAssociations = []*AwsRoute53ZoneVpcAssociation{{}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Public-only features on private zones
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a delegation set on a private zone", func() {
		spec.IsPrivate = true
		spec.VpcAssociations = []*AwsRoute53ZoneVpcAssociation{
			{VpcId: strRef("vpc-0123456789abcdef0")},
		}
		spec.DelegationSetId = "N1PA6795SAMPLE"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects accelerated recovery on a private zone", func() {
		spec.IsPrivate = true
		spec.VpcAssociations = []*AwsRoute53ZoneVpcAssociation{
			{VpcId: strRef("vpc-0123456789abcdef0")},
		}
		spec.EnableAcceleratedRecovery = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects query logging on a private zone", func() {
		spec.IsPrivate = true
		spec.VpcAssociations = []*AwsRoute53ZoneVpcAssociation{
			{VpcId: strRef("vpc-0123456789abcdef0")},
		}
		spec.QueryLogging = &AwsRoute53ZoneQueryLogging{
			CloudwatchLogGroupArn: strRef("arn:aws:logs:us-east-1:123456789012:log-group:/aws/route53/x"),
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects DNSSEC on a private zone", func() {
		spec.IsPrivate = true
		spec.VpcAssociations = []*AwsRoute53ZoneVpcAssociation{
			{VpcId: strRef("vpc-0123456789abcdef0")},
		}
		spec.Dnssec = &AwsRoute53ZoneDnssec{
			KmsKeyArn: strRef("arn:aws:kms:us-east-1:123456789012:key/abc"),
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Nested requirements
	// -------------------------------------------------------------------------

	ginkgo.It("rejects query logging without a log group", func() {
		spec.QueryLogging = &AwsRoute53ZoneQueryLogging{}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects DNSSEC without a KMS key", func() {
		spec.Dnssec = &AwsRoute53ZoneDnssec{}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a KSK name with invalid characters", func() {
		spec.Dnssec = &AwsRoute53ZoneDnssec{
			KmsKeyArn:         strRef("arn:aws:kms:us-east-1:123456789012:key/abc"),
			KeySigningKeyName: "bad name!",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Scalars
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a missing region", func() {
		spec.Region = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a comment over 256 characters", func() {
		long := make([]byte, 257)
		for i := range long {
			long[i] = 'a'
		}
		spec.Comment = string(long)
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})
})
