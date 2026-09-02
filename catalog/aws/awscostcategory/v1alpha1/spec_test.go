package awscostcategoryv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsCostCategorySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCostCategorySpec Validation Suite")
}

// minimalCategory is the smallest valid instance: one REGULAR rule
// assigning a value to one service's spend (SERVICE_CODE - the only
// by-service dimension cost category rules accept).
func minimalCategory() *AwsCostCategorySpec {
	return &AwsCostCategorySpec{
		Region:       "us-east-1",
		CategoryName: "Cost Center",
		Rules: []*AwsCostCategoryRule{{
			Value: "compute",
			Rule: &AwsCostCategoryExpression{
				Dimension: &AwsCostCategoryExpressionDimension{
					Key:    "SERVICE_CODE",
					Values: []string{"AmazonEC2"},
				},
			},
		}},
	}
}

var _ = ginkgo.Describe("AwsCostCategorySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal single-rule category", func() {
			gomega.Expect(protovalidate.Validate(minimalCategory())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an inherited-value rule", func() {
			spec := minimalCategory()
			spec.Rules = []*AwsCostCategoryRule{{
				Type: "INHERITED_VALUE",
				InheritedValue: &AwsCostCategoryInheritedValue{
					DimensionName: "TAG",
					DimensionKey:  "team",
				},
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a nested expression with tag ABSENT matching", func() {
			spec := minimalCategory()
			spec.Rules[0].Rule = &AwsCostCategoryExpression{
				Or: []*AwsCostCategoryExpressionNode{
					{Tag: &AwsCostCategoryExpressionTag{Key: "team", MatchOptions: []string{"ABSENT"}}},
					{Dimension: &AwsCostCategoryExpressionDimension{Key: "REGION", Values: []string{"us-east-1"}}},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts default_value, effective_start, and split-charge rules", func() {
			spec := minimalCategory()
			spec.DefaultValue = "shared"
			spec.EffectiveStart = "2026-01-01T00:00:00Z"
			spec.SplitChargeRules = []*AwsCostCategorySplitChargeRule{{
				Source:  "shared",
				Targets: []string{"compute", "data"},
				Method:  "PROPORTIONAL",
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a FIXED split with allocation percentages", func() {
			spec := minimalCategory()
			spec.SplitChargeRules = []*AwsCostCategorySplitChargeRule{{
				Source:  "shared",
				Targets: []string{"compute", "data"},
				Method:  "FIXED",
				Parameters: []*AwsCostCategorySplitChargeParameter{{
					Type:   "ALLOCATION_PERCENTAGES",
					Values: []string{"60", "40"},
				}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalCategory()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a category name over 50 characters", func() {
			spec := minimalCategory()
			spec.CategoryName = "an-extremely-long-cost-category-name-over-the-cap-x"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty rules list", func() {
			spec := minimalCategory()
			spec.Rules = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rule with both shapes", func() {
			spec := minimalCategory()
			spec.Rules[0].InheritedValue = &AwsCostCategoryInheritedValue{DimensionName: "TAG"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rule with neither shape", func() {
			spec := minimalCategory()
			spec.Rules[0].Rule = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an INHERITED_VALUE type on an expression rule", func() {
			spec := minimalCategory()
			spec.Rules[0].Type = "INHERITED_VALUE"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a REGULAR rule without its value", func() {
			spec := minimalCategory()
			spec.Rules[0].Value = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a mid-month effective_start", func() {
			spec := minimalCategory()
			spec.EffectiveStart = "2026-01-15T00:00:00Z"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a FIXED split without allocation percentages", func() {
			spec := minimalCategory()
			spec.SplitChargeRules = []*AwsCostCategorySplitChargeRule{{
				Source:  "shared",
				Targets: []string{"compute"},
				Method:  "FIXED",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown split method", func() {
			spec := minimalCategory()
			spec.SplitChargeRules = []*AwsCostCategorySplitChargeRule{{
				Source:  "shared",
				Targets: []string{"compute"},
				Method:  "WEIGHTED",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects the SERVICE display-name dimension (rules take SERVICE_CODE - the server names the allowed set)", func() {
			spec := minimalCategory()
			spec.Rules[0].Rule.Dimension.Key = "SERVICE"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
