package awsbudgetv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBudgetSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBudgetSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalBudget is the smallest valid instance: a monthly cost budget
// with a fixed limit.
func minimalBudget() *AwsBudgetSpec {
	return &AwsBudgetSpec{
		Region:     "us-east-1",
		BudgetName: "Monthly AWS Spend",
		BudgetType: "COST",
		TimeUnit:   "MONTHLY",
		Limit: &AwsBudgetLimit{
			Amount: "1000",
			Unit:   "USD",
		},
	}
}

func validAction() *AwsBudgetAction {
	return &AwsBudgetAction{
		Name:             "freeze-dev",
		ActionType:       "APPLY_IAM_POLICY",
		ApprovalModel:    "AUTOMATIC",
		NotificationType: "ACTUAL",
		ExecutionRoleArn: svr("arn:aws:iam::123456789012:role/budget-actions"),
		ActionThreshold: &AwsBudgetActionThreshold{
			ActionThresholdType:  "PERCENTAGE",
			ActionThresholdValue: 100,
		},
		Subscribers: []*AwsBudgetActionSubscriber{{
			Address:          svr("finops@example.com"),
			SubscriptionType: "EMAIL",
		}},
		IamActionDefinition: &AwsBudgetIamActionDefinition{
			PolicyArn: svr("arn:aws:iam::aws:policy/AWSDenyAll"),
			Groups:    []*foreignkeyv1.StringValueOrRef{svr("developers")},
		},
	}
}

var _ = ginkgo.Describe("AwsBudgetSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal fixed-limit budget", func() {
			gomega.Expect(protovalidate.Validate(minimalBudget())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a planned-limits budget", func() {
			spec := minimalBudget()
			spec.Limit = nil
			spec.PlannedLimits = []*AwsBudgetPlannedLimit{
				{StartTime: "2026-01-01_00:00", Amount: "1000", Unit: "USD"},
				{StartTime: "2026-02-01_00:00", Amount: "1200", Unit: "USD"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an auto-adjusting budget", func() {
			spec := minimalBudget()
			spec.Limit = nil
			spec.AutoAdjust = &AwsBudgetAutoAdjust{
				AutoAdjustType:         "HISTORICAL",
				BudgetAdjustmentPeriod: 6,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the modern metric + filter_expression pair", func() {
			spec := minimalBudget()
			spec.Metric = "UnblendedCost"
			spec.FilterExpression = &AwsBudgetFilterExpression{
				Dimension: &AwsBudgetFilterDimension{
					Key:    "SERVICE",
					Values: []string{"Amazon Elastic Compute Cloud - Compute"},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a nested filter expression", func() {
			spec := minimalBudget()
			spec.Metric = "UnblendedCost"
			spec.FilterExpression = &AwsBudgetFilterExpression{
				And: []*AwsBudgetFilterExpressionNode{
					{Dimension: &AwsBudgetFilterDimension{Key: "REGION", Values: []string{"us-east-1"}}},
					{Not: &AwsBudgetFilterExpressionLeaf{
						Tag: &AwsBudgetFilterTag{Key: "environment", Values: []string{"sandbox"}},
					}},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts legacy cost_filters with cost_types", func() {
			spec := minimalBudget()
			spec.CostFilters = []*AwsBudgetCostFilter{{Name: "Service", Values: []string{"Amazon Relational Database Service"}}}
			include := false
			spec.CostTypes = &AwsBudgetCostTypes{IncludeCredit: &include}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a notification with subscribers", func() {
			spec := minimalBudget()
			spec.Notifications = []*AwsBudgetNotification{{
				ComparisonOperator:       "GREATER_THAN",
				NotificationType:         "FORECASTED",
				Threshold:                80,
				ThresholdType:            "PERCENTAGE",
				SubscriberEmailAddresses: []string{"finops@example.com"},
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an IAM-arm action", func() {
			spec := minimalBudget()
			spec.Actions = []*AwsBudgetAction{validAction()}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an SSM-arm action", func() {
			spec := minimalBudget()
			action := validAction()
			action.ActionType = "RUN_SSM_DOCUMENTS"
			action.IamActionDefinition = nil
			action.SsmActionDefinition = &AwsBudgetSsmActionDefinition{
				ActionSubType: "STOP_EC2_INSTANCES",
				Region:        "us-west-2",
				InstanceIds:   []*foreignkeyv1.StringValueOrRef{svr("i-0abc123def4567890")},
			}
			spec.Actions = []*AwsBudgetAction{action}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalBudget()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a budget name containing a colon", func() {
			spec := minimalBudget()
			spec.BudgetName = "team:budget"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown budget type", func() {
			spec := minimalBudget()
			spec.BudgetType = "SPEND"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects no funding shape", func() {
			spec := minimalBudget()
			spec.Limit = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects two funding shapes together", func() {
			spec := minimalBudget()
			spec.AutoAdjust = &AwsBudgetAutoAdjust{AutoAdjustType: "FORECAST"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects both filter generations together", func() {
			spec := minimalBudget()
			spec.Metric = "UnblendedCost"
			spec.FilterExpression = &AwsBudgetFilterExpression{
				Dimension: &AwsBudgetFilterDimension{Key: "SERVICE", Values: []string{"AmazonEC2"}},
			}
			spec.CostFilters = []*AwsBudgetCostFilter{{Name: "Service", Values: []string{"AmazonEC2"}}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a metric without its filter expression", func() {
			spec := minimalBudget()
			spec.Metric = "UnblendedCost"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a filter expression without its metric", func() {
			spec := minimalBudget()
			spec.FilterExpression = &AwsBudgetFilterExpression{
				Dimension: &AwsBudgetFilterDimension{Key: "SERVICE", Values: []string{"AmazonEC2"}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects the metric alongside cost_types", func() {
			spec := minimalBudget()
			spec.Metric = "UnblendedCost"
			spec.FilterExpression = &AwsBudgetFilterExpression{
				Dimension: &AwsBudgetFilterDimension{Key: "SERVICE", Values: []string{"AmazonEC2"}},
			}
			include := true
			spec.CostTypes = &AwsBudgetCostTypes{IncludeTax: &include}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects the ABSENT match option (budgets rejects it at plan time)", func() {
			spec := minimalBudget()
			spec.Metric = "UnblendedCost"
			spec.FilterExpression = &AwsBudgetFilterExpression{
				Tag: &AwsBudgetFilterTag{Key: "environment", MatchOptions: []string{"ABSENT"}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed time period", func() {
			spec := minimalBudget()
			spec.TimePeriodStart = "2026-01-01 00:00"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an out-of-range adjustment period", func() {
			spec := minimalBudget()
			spec.Limit = nil
			spec.AutoAdjust = &AwsBudgetAutoAdjust{
				AutoAdjustType:         "HISTORICAL",
				BudgetAdjustmentPeriod: 61,
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a subscriber-less notification", func() {
			spec := minimalBudget()
			spec.Notifications = []*AwsBudgetNotification{{
				ComparisonOperator: "GREATER_THAN",
				NotificationType:   "ACTUAL",
				Threshold:          100,
				ThresholdType:      "PERCENTAGE",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate action names", func() {
			spec := minimalBudget()
			spec.Actions = []*AwsBudgetAction{validAction(), validAction()}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a definition arm that mismatches the action type", func() {
			spec := minimalBudget()
			action := validAction()
			action.ActionType = "APPLY_SCP_POLICY"
			spec.Actions = []*AwsBudgetAction{action}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an IAM action naming no principals", func() {
			spec := minimalBudget()
			action := validAction()
			action.IamActionDefinition.Groups = nil
			spec.Actions = []*AwsBudgetAction{action}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an action without subscribers", func() {
			spec := minimalBudget()
			action := validAction()
			action.Subscribers = nil
			spec.Actions = []*AwsBudgetAction{action}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed account id", func() {
			spec := minimalBudget()
			spec.AccountId = "12345"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
