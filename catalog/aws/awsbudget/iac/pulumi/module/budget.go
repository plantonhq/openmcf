package module

import (
	"fmt"

	"github.com/pkg/errors"
	awsbudgetv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbudget/v1alpha1"
	"github.com/plantonhq/planton/internal/valuefrom"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/budgets"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// budget provisions the Budgets budget, its notifications, and its
// folded actions, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the budget's name is spec.budget_name (an explicit field -
//     budget names legally carry spaces metadata.name cannot) and
//     changing it replaces the budget;
//   - exactly ONE funding shape renders (spec-enforced): the fixed
//     limit, the planned per-period limits, or auto-adjustment;
//   - the two filter GENERATIONS are mutually exclusive
//     (spec-enforced): legacy cost_filters/cost_types vs the modern
//     metric + filter_expression pair. The provider's `metrics` is a
//     single-element list flattened to one string by the SDK - the
//     spec's singular `metric` maps 1:1;
//   - each action is one aws_budgets_budget_action keyed by its spec
//     name (the for_each key and the outputs-map key); its definition
//     arm matches action_type (spec-enforced);
//   - both resources are taggable; account_id (member-account budgets
//     managed from a payer account) is create-only on both.
func budget(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	spec := locals.Spec

	args := &budgets.BudgetArgs{
		Name:       pulumi.StringPtr(spec.BudgetName),
		BudgetType: pulumi.String(spec.BudgetType),
		TimeUnit:   pulumi.String(spec.TimeUnit),
		Tags:       pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.AccountId != "" {
		args.AccountId = pulumi.StringPtr(spec.AccountId)
	}
	if spec.BillingViewArn != "" {
		args.BillingViewArn = pulumi.StringPtr(spec.BillingViewArn)
	}
	if spec.Limit != nil {
		args.LimitAmount = pulumi.StringPtr(spec.Limit.Amount)
		args.LimitUnit = pulumi.StringPtr(spec.Limit.Unit)
	}
	if len(spec.PlannedLimits) > 0 {
		plannedLimits := budgets.BudgetPlannedLimitArray{}
		for _, p := range spec.PlannedLimits {
			plannedLimits = append(plannedLimits, budgets.BudgetPlannedLimitArgs{
				StartTime: pulumi.String(p.StartTime),
				Amount:    pulumi.String(p.Amount),
				Unit:      pulumi.String(p.Unit),
			})
		}
		args.PlannedLimits = plannedLimits
	}
	if spec.AutoAdjust != nil {
		autoAdjust := &budgets.BudgetAutoAdjustDataArgs{
			AutoAdjustType: pulumi.String(spec.AutoAdjust.AutoAdjustType),
		}
		if spec.AutoAdjust.BudgetAdjustmentPeriod != 0 {
			autoAdjust.HistoricalOptions = &budgets.BudgetAutoAdjustDataHistoricalOptionsArgs{
				BudgetAdjustmentPeriod: pulumi.Int(int(spec.AutoAdjust.BudgetAdjustmentPeriod)),
			}
		}
		args.AutoAdjustData = autoAdjust
	}
	if spec.TimePeriodStart != "" {
		args.TimePeriodStart = pulumi.StringPtr(spec.TimePeriodStart)
	}
	if spec.TimePeriodEnd != "" {
		args.TimePeriodEnd = pulumi.StringPtr(spec.TimePeriodEnd)
	}
	if spec.CostTypes != nil {
		args.CostTypes = buildCostTypes(spec.CostTypes)
	}
	if len(spec.CostFilters) > 0 {
		costFilters := budgets.BudgetCostFilterArray{}
		for _, f := range spec.CostFilters {
			costFilters = append(costFilters, budgets.BudgetCostFilterArgs{
				Name:   pulumi.String(f.Name),
				Values: pulumi.ToStringArray(f.Values),
			})
		}
		args.CostFilters = costFilters
	}
	if spec.Metric != "" {
		args.Metrics = pulumi.StringPtr(spec.Metric)
	}
	if spec.FilterExpression != nil {
		args.FilterExpression = buildFilterExpression(spec.FilterExpression)
	}
	if len(spec.Notifications) > 0 {
		notifications := budgets.BudgetNotificationArray{}
		for _, n := range spec.Notifications {
			notification := budgets.BudgetNotificationArgs{
				ComparisonOperator: pulumi.String(n.ComparisonOperator),
				NotificationType:   pulumi.String(n.NotificationType),
				Threshold:          pulumi.Float64(n.Threshold),
				ThresholdType:      pulumi.String(n.ThresholdType),
			}
			if len(n.SubscriberEmailAddresses) > 0 {
				notification.SubscriberEmailAddresses = pulumi.ToStringArray(n.SubscriberEmailAddresses)
			}
			if snsArns := valuefrom.ToStringArray(n.SubscriberSnsTopicArns); len(snsArns) > 0 {
				notification.SubscriberSnsTopicArns = pulumi.ToStringArray(snsArns)
			}
			notifications = append(notifications, notification)
		}
		args.Notifications = notifications
	}

	createdBudget, err := budgets.NewBudget(ctx, "budget", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create budget")
	}

	actionIds := pulumi.StringMap{}
	for _, action := range spec.Actions {
		createdAction, err := budgetAction(ctx, locals, provider, createdBudget, action)
		if err != nil {
			return errors.Wrapf(err, "action %s", action.Name)
		}
		actionIds[action.Name] = createdAction.ActionId
	}

	ctx.Export(OpBudgetName, createdBudget.Name)
	ctx.Export(OpBudgetArn, createdBudget.Arn)
	ctx.Export(OpAccountId, createdBudget.AccountId)
	ctx.Export(OpActionIds, actionIds)

	return nil
}

// buildCostTypes renders the cost-component toggles. Presence-typed
// fields render only when set - the provider then applies AWS's
// defaults (include_* true; use_blended/use_amortized false) for the
// rest.
func buildCostTypes(costTypes *awsbudgetv1alpha1.AwsBudgetCostTypes) *budgets.BudgetCostTypesArgs {
	out := &budgets.BudgetCostTypesArgs{}
	if costTypes.IncludeCredit != nil {
		out.IncludeCredit = pulumi.BoolPtr(costTypes.GetIncludeCredit())
	}
	if costTypes.IncludeDiscount != nil {
		out.IncludeDiscount = pulumi.BoolPtr(costTypes.GetIncludeDiscount())
	}
	if costTypes.IncludeOtherSubscription != nil {
		out.IncludeOtherSubscription = pulumi.BoolPtr(costTypes.GetIncludeOtherSubscription())
	}
	if costTypes.IncludeRecurring != nil {
		out.IncludeRecurring = pulumi.BoolPtr(costTypes.GetIncludeRecurring())
	}
	if costTypes.IncludeRefund != nil {
		out.IncludeRefund = pulumi.BoolPtr(costTypes.GetIncludeRefund())
	}
	if costTypes.IncludeSubscription != nil {
		out.IncludeSubscription = pulumi.BoolPtr(costTypes.GetIncludeSubscription())
	}
	if costTypes.IncludeSupport != nil {
		out.IncludeSupport = pulumi.BoolPtr(costTypes.GetIncludeSupport())
	}
	if costTypes.IncludeTax != nil {
		out.IncludeTax = pulumi.BoolPtr(costTypes.GetIncludeTax())
	}
	if costTypes.IncludeUpfront != nil {
		out.IncludeUpfront = pulumi.BoolPtr(costTypes.GetIncludeUpfront())
	}
	if costTypes.UseAmortized != nil {
		out.UseAmortized = pulumi.BoolPtr(costTypes.GetUseAmortized())
	}
	if costTypes.UseBlended != nil {
		out.UseBlended = pulumi.BoolPtr(costTypes.GetUseBlended())
	}
	return out
}

// budgetAction provisions one folded budget action. The definition arm
// was already validated against action_type by the spec's CEL rules;
// valueFrom references (execution role, policy ARNs, principals,
// instance IDs) were resolved before the module ran.
func budgetAction(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource,
	createdBudget *budgets.Budget, action *awsbudgetv1alpha1.AwsBudgetAction) (*budgets.BudgetAction, error) {

	definition := budgets.BudgetActionDefinitionArgs{}
	if iamArm := action.IamActionDefinition; iamArm != nil {
		iamDefinition := &budgets.BudgetActionDefinitionIamActionDefinitionArgs{
			PolicyArn: pulumi.String(iamArm.PolicyArn.GetValue()),
		}
		if groups := valuefrom.ToStringArray(iamArm.Groups); len(groups) > 0 {
			iamDefinition.Groups = pulumi.ToStringArray(groups)
		}
		if roles := valuefrom.ToStringArray(iamArm.Roles); len(roles) > 0 {
			iamDefinition.Roles = pulumi.ToStringArray(roles)
		}
		if users := valuefrom.ToStringArray(iamArm.Users); len(users) > 0 {
			iamDefinition.Users = pulumi.ToStringArray(users)
		}
		definition.IamActionDefinition = iamDefinition
	}
	if scpArm := action.ScpActionDefinition; scpArm != nil {
		definition.ScpActionDefinition = &budgets.BudgetActionDefinitionScpActionDefinitionArgs{
			PolicyId:  pulumi.String(scpArm.PolicyId.GetValue()),
			TargetIds: pulumi.ToStringArray(scpArm.TargetIds),
		}
	}
	if ssmArm := action.SsmActionDefinition; ssmArm != nil {
		definition.SsmActionDefinition = &budgets.BudgetActionDefinitionSsmActionDefinitionArgs{
			ActionSubType: pulumi.String(ssmArm.ActionSubType),
			Region:        pulumi.String(ssmArm.Region),
			InstanceIds:   pulumi.ToStringArray(valuefrom.ToStringArray(ssmArm.InstanceIds)),
		}
	}

	subscribers := budgets.BudgetActionSubscriberArray{}
	for _, s := range action.Subscribers {
		subscribers = append(subscribers, budgets.BudgetActionSubscriberArgs{
			Address:          pulumi.String(s.Address.GetValue()),
			SubscriptionType: pulumi.String(s.SubscriptionType),
		})
	}

	args := &budgets.BudgetActionArgs{
		BudgetName:    createdBudget.Name,
		ActionType:    pulumi.String(action.ActionType),
		ApprovalModel: pulumi.String(action.ApprovalModel),
		ActionThreshold: budgets.BudgetActionActionThresholdArgs{
			ActionThresholdType:  pulumi.String(action.ActionThreshold.ActionThresholdType),
			ActionThresholdValue: pulumi.Float64(action.ActionThreshold.ActionThresholdValue),
		},
		Definition:       definition,
		ExecutionRoleArn: pulumi.String(action.ExecutionRoleArn.GetValue()),
		NotificationType: pulumi.String(action.NotificationType),
		Subscribers:      subscribers,
		Tags:             pulumi.ToStringMap(locals.AwsTags),
	}
	// The action inherits the budget's account scope: a member-account
	// budget's actions must live in the same account.
	if locals.Spec.AccountId != "" {
		args.AccountId = pulumi.StringPtr(locals.Spec.AccountId)
	}

	createdAction, err := budgets.NewBudgetAction(ctx, fmt.Sprintf("action-%s", action.Name), args,
		pulumi.Provider(provider), pulumi.Parent(createdBudget))
	if err != nil {
		return nil, errors.Wrap(err, "create budget action")
	}
	return createdAction, nil
}
