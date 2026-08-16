package awsssmassociationv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsSsmAssociationSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSsmAssociationSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalAssociation is the smallest valid instance: an AWS-managed
// document bound with no targets.
func minimalAssociation() *AwsSsmAssociationSpec {
	return &AwsSsmAssociationSpec{
		Region:       "us-west-2",
		DocumentName: svr("AWS-RunShellScript"),
	}
}

var _ = ginkgo.Describe("AwsSsmAssociationSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal AWS-managed-document association", func() {
			gomega.Expect(protovalidate.Validate(minimalAssociation())).To(gomega.BeNil())
		})

		ginkgo.It("accepts the full scheduling and compliance surface", func() {
			spec := minimalAssociation()
			spec.AssociationName = "nightly-patch-scan"
			spec.DocumentVersion = "$DEFAULT"
			spec.Parameters = map[string]string{"Operation": "Scan"}
			spec.Targets = []*AwsSsmAssociationTarget{{
				Key:    "tag:env",
				Values: []string{"prod"},
			}}
			spec.ScheduleExpression = "cron(0 2 ? * SUN *)"
			spec.ApplyOnlyAtCronInterval = true
			spec.ComplianceSeverity = "HIGH"
			spec.SyncCompliance = "AUTO"
			spec.MaxConcurrency = "10%"
			spec.MaxErrors = "0"
			spec.CalendarNames = []string{"change-freeze-calendar"}
			spec.OutputLocation = &AwsSsmAssociationOutputLocation{
				S3BucketName: svr("command-output-bucket"),
				S3KeyPrefix:  "ssm/runs",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a concrete document version pin", func() {
			spec := minimalAssociation()
			spec.DocumentVersion = "3"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing document reference", func() {
			spec := minimalAssociation()
			spec.DocumentName = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an association name below 3 characters", func() {
			spec := minimalAssociation()
			spec.AssociationName = "ab"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed document version", func() {
			spec := minimalAssociation()
			spec.DocumentVersion = "latest"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects more than 5 targets", func() {
			spec := minimalAssociation()
			for i := 0; i < 6; i++ {
				spec.Targets = append(spec.Targets, &AwsSsmAssociationTarget{
					Key:    "InstanceIds",
					Values: []string{"i-0abc"},
				})
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a target without values", func() {
			spec := minimalAssociation()
			spec.Targets = []*AwsSsmAssociationTarget{{Key: "tag:env"}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects max_concurrency of zero", func() {
			spec := minimalAssociation()
			spec.MaxConcurrency = "0"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed max_errors", func() {
			spec := minimalAssociation()
			spec.MaxErrors = "-1"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown compliance severity", func() {
			spec := minimalAssociation()
			spec.ComplianceSeverity = "SEVERE"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an output location without a bucket", func() {
			spec := minimalAssociation()
			spec.OutputLocation = &AwsSsmAssociationOutputLocation{S3KeyPrefix: "ssm/runs"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a negative wait_for_success timeout", func() {
			spec := minimalAssociation()
			spec.WaitForSuccessTimeoutSeconds = -30
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
