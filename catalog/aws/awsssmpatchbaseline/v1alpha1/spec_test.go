package awsssmpatchbaselinev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsSsmPatchBaselineSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSsmPatchBaselineSpec Validation Suite")
}

func int32Ptr(i int32) *int32 { return &i }

// minimalBaseline is the smallest valid instance: an empty baseline
// for the provider-default WINDOWS operating system.
func minimalBaseline() *AwsSsmPatchBaselineSpec {
	return &AwsSsmPatchBaselineSpec{Region: "us-west-2"}
}

func securityCriticalFilter() []*AwsSsmPatchBaselinePatchFilter {
	return []*AwsSsmPatchBaselinePatchFilter{{
		Key:    "CLASSIFICATION",
		Values: []string{"SecurityUpdates"},
	}}
}

var _ = ginkgo.Describe("AwsSsmPatchBaselineSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal baseline", func() {
			gomega.Expect(protovalidate.Validate(minimalBaseline())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a days-based approval rule including zero days", func() {
			spec := minimalBaseline()
			spec.ApprovalRules = []*AwsSsmPatchBaselineApprovalRule{{
				PatchFilters:     securityCriticalFilter(),
				ApproveAfterDays: int32Ptr(0),
				ComplianceLevel:  "CRITICAL",
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a date-based approval rule", func() {
			spec := minimalBaseline()
			spec.ApprovalRules = []*AwsSsmPatchBaselineApprovalRule{{
				PatchFilters:     securityCriticalFilter(),
				ApproveUntilDate: "2026-12-31",
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a rule with neither approval arm (the Debian/Ubuntu shape)", func() {
			spec := minimalBaseline()
			spec.OperatingSystem = "UBUNTU"
			spec.ApprovalRules = []*AwsSsmPatchBaselineApprovalRule{{
				PatchFilters: []*AwsSsmPatchBaselinePatchFilter{{
					Key:    "PRIORITY",
					Values: []string{"Required"},
				}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the full surface: filters, patch lists, sources, groups, designation", func() {
			spec := minimalBaseline()
			spec.OperatingSystem = "AMAZON_LINUX_2023"
			spec.Description = "AL2023 security baseline"
			spec.GlobalFilters = []*AwsSsmPatchBaselinePatchFilter{{
				Key:    "PRODUCT",
				Values: []string{"AmazonLinux2023"},
			}}
			spec.ApprovedPatches = []string{"kernel-6.1.0"}
			spec.ApprovedPatchesComplianceLevel = "HIGH"
			spec.ApprovedPatchesEnableNonSecurity = true
			spec.RejectedPatches = []string{"nginx-1.25.0"}
			spec.RejectedPatchesAction = "BLOCK"
			spec.Sources = []*AwsSsmPatchBaselineSource{{
				Name:          "internal-repo",
				Configuration: "[internal]\nbaseurl=https://repo.example.com/al2023",
				Products:      []string{"AmazonLinux2023.3"},
			}}
			spec.PatchGroups = []string{"prod-servers", "staging-servers"}
			spec.SetAsDefaultBaseline = true
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalBaseline()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown operating system", func() {
			spec := minimalBaseline()
			spec.OperatingSystem = "FREEBSD"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rule carrying both approval arms", func() {
			spec := minimalBaseline()
			spec.ApprovalRules = []*AwsSsmPatchBaselineApprovalRule{{
				PatchFilters:     securityCriticalFilter(),
				ApproveAfterDays: int32Ptr(7),
				ApproveUntilDate: "2026-12-31",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rule without patch filters", func() {
			spec := minimalBaseline()
			spec.ApprovalRules = []*AwsSsmPatchBaselineApprovalRule{{ApproveAfterDays: int32Ptr(7)}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects approve_after_days above 360", func() {
			spec := minimalBaseline()
			spec.ApprovalRules = []*AwsSsmPatchBaselineApprovalRule{{
				PatchFilters:     securityCriticalFilter(),
				ApproveAfterDays: int32Ptr(361),
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed approve_until_date", func() {
			spec := minimalBaseline()
			spec.ApprovalRules = []*AwsSsmPatchBaselineApprovalRule{{
				PatchFilters:     securityCriticalFilter(),
				ApproveUntilDate: "2026-13-01",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown patch filter key", func() {
			spec := minimalBaseline()
			spec.GlobalFilters = []*AwsSsmPatchBaselinePatchFilter{{
				Key:    "KERNEL",
				Values: []string{"*"},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects more than 4 global filters", func() {
			spec := minimalBaseline()
			for i := 0; i < 5; i++ {
				spec.GlobalFilters = append(spec.GlobalFilters, &AwsSsmPatchBaselinePatchFilter{
					Key:    "PRODUCT",
					Values: []string{"*"},
				})
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate patch groups", func() {
			spec := minimalBaseline()
			spec.PatchGroups = []string{"prod-servers", "prod-servers"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a source without products", func() {
			spec := minimalBaseline()
			spec.Sources = []*AwsSsmPatchBaselineSource{{
				Name:          "internal-repo",
				Configuration: "[internal]",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a source name with illegal characters", func() {
			spec := minimalBaseline()
			spec.Sources = []*AwsSsmPatchBaselineSource{{
				Name:          "internal repo",
				Configuration: "[internal]",
				Products:      []string{"*"},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown rejected_patches_action", func() {
			spec := minimalBaseline()
			spec.RejectedPatchesAction = "IGNORE"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
