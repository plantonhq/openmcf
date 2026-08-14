package awsguarddutyv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsGuardDutySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsGuardDutySpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalDetector is the smallest valid instance: region only - a
// plain enabled detector.
func minimalDetector() *AwsGuardDutySpec {
	return &AwsGuardDutySpec{
		Region: "us-west-2",
	}
}

var _ = ginkgo.Describe("AwsGuardDutySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal detector", func() {
			gomega.Expect(protovalidate.Validate(minimalDetector())).To(gomega.BeNil())
		})

		ginkgo.It("accepts features with sub-toggles", func() {
			spec := minimalDetector()
			spec.Features = []*AwsGuardDutyFeature{
				{Name: "S3_DATA_EVENTS"},
				{
					Name: "RUNTIME_MONITORING",
					AdditionalConfiguration: []*AwsGuardDutyFeatureAdditionalConfiguration{
						{Name: "EC2_AGENT_MANAGEMENT"},
						{Name: "ECS_FARGATE_AGENT_MANAGEMENT", Enabled: proto.Bool(false)},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a filter with criteria and bounds", func() {
			spec := minimalDetector()
			spec.Filters = []*AwsGuardDutyFilter{{
				Name:   "archive-low-severity",
				Action: "ARCHIVE",
				Rank:   1,
				Criteria: []*AwsGuardDutyFilterCriterion{
					{Field: "severity", LessThan: "4"},
					{Field: "region", Equals: []string{"us-west-2"}},
				},
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts IP sets, threat intel sets, and a publishing destination", func() {
			spec := minimalDetector()
			spec.IpSets = []*AwsGuardDutyIpSet{{
				Name:     "trusted-office",
				Format:   "TXT",
				Location: "https://s3.amazonaws.com/my-bucket/trusted.txt",
				Activate: true,
			}}
			spec.ThreatIntelSets = []*AwsGuardDutyThreatIntelSet{{
				Name:     "known-bad",
				Format:   "TXT",
				Location: "https://s3.amazonaws.com/my-bucket/bad.txt",
				Activate: true,
			}}
			spec.PublishingDestination = &AwsGuardDutyPublishingDestination{
				BucketArn: svr("arn:aws:s3:::findings-archive"),
				KmsKeyArn: svr("arn:aws:kms:us-west-2:123456789012:key/abc"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the admin posture with organization and members", func() {
			spec := minimalDetector()
			spec.Organization = &AwsGuardDutyOrganization{
				AutoEnableOrganizationMembers: "NEW",
				Features: []*AwsGuardDutyOrganizationFeature{{
					Name:       "RUNTIME_MONITORING",
					AutoEnable: "ALL",
					AdditionalConfiguration: []*AwsGuardDutyOrganizationFeatureAdditionalConfiguration{{
						Name:       "EC2_AGENT_MANAGEMENT",
						AutoEnable: "NEW",
					}},
				}},
			}
			spec.Members = []*AwsGuardDutyMember{{
				AccountId: "111111111111",
				Email:     "security@example.com",
				Features: []*AwsGuardDutyMemberFeature{{
					Name: "EKS_AUDIT_LOGS",
				}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the member posture", func() {
			spec := minimalDetector()
			spec.AcceptInvitationFromAccountId = "222222222222"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalDetector()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate feature names", func() {
			spec := minimalDetector()
			spec.Features = []*AwsGuardDutyFeature{
				{Name: "S3_DATA_EVENTS"},
				{Name: "S3_DATA_EVENTS", Enabled: proto.Bool(false)},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown feature name", func() {
			spec := minimalDetector()
			spec.Features = []*AwsGuardDutyFeature{{Name: "DNS_LOGS"}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a filter without criteria", func() {
			spec := minimalDetector()
			spec.Filters = []*AwsGuardDutyFilter{{
				Name:   "empty",
				Action: "NOOP",
				Rank:   1,
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a criterion without any condition", func() {
			spec := minimalDetector()
			spec.Filters = []*AwsGuardDutyFilter{{
				Name:     "no-condition",
				Action:   "NOOP",
				Rank:     1,
				Criteria: []*AwsGuardDutyFilterCriterion{{Field: "severity"}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate criterion fields in one filter", func() {
			spec := minimalDetector()
			spec.Filters = []*AwsGuardDutyFilter{{
				Name:   "dup-fields",
				Action: "NOOP",
				Rank:   1,
				Criteria: []*AwsGuardDutyFilterCriterion{
					{Field: "severity", GreaterThan: "4"},
					{Field: "severity", LessThan: "8"},
				},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects more than five match patterns", func() {
			spec := minimalDetector()
			spec.Filters = []*AwsGuardDutyFilter{{
				Name:   "too-many-matches",
				Action: "NOOP",
				Rank:   1,
				Criteria: []*AwsGuardDutyFilterCriterion{{
					Field:   "type",
					Matches: []string{"a", "b", "c", "d", "e", "f"},
				}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a filter rank above 100", func() {
			spec := minimalDetector()
			spec.Filters = []*AwsGuardDutyFilter{{
				Name:     "rank-overflow",
				Action:   "NOOP",
				Rank:     101,
				Criteria: []*AwsGuardDutyFilterCriterion{{Field: "severity", Equals: []string{"8"}}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown ip set format", func() {
			spec := minimalDetector()
			spec.IpSets = []*AwsGuardDutyIpSet{{
				Name:     "bad-format",
				Format:   "CSV",
				Location: "https://s3.amazonaws.com/my-bucket/list.csv",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a publishing destination without its KMS key", func() {
			spec := minimalDetector()
			spec.PublishingDestination = &AwsGuardDutyPublishingDestination{
				BucketArn: svr("arn:aws:s3:::findings-archive"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects the org feature vocabulary violation AI_ANALYST", func() {
			spec := minimalDetector()
			spec.Organization = &AwsGuardDutyOrganization{
				AutoEnableOrganizationMembers: "NEW",
				Features: []*AwsGuardDutyOrganizationFeature{{
					Name:       "AI_ANALYST",
					AutoEnable: "ALL",
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an organization block without the auto-enable posture", func() {
			spec := minimalDetector()
			spec.Organization = &AwsGuardDutyOrganization{}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate member accounts", func() {
			spec := minimalDetector()
			spec.Members = []*AwsGuardDutyMember{
				{AccountId: "111111111111", Email: "a@example.com"},
				{AccountId: "111111111111", Email: "b@example.com"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a member with a malformed account id", func() {
			spec := minimalDetector()
			spec.Members = []*AwsGuardDutyMember{{AccountId: "1234", Email: "a@example.com"}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects mixing the member and admin postures", func() {
			spec := minimalDetector()
			spec.AcceptInvitationFromAccountId = "222222222222"
			spec.Organization = &AwsGuardDutyOrganization{AutoEnableOrganizationMembers: "NEW"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
