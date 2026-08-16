// Cost Explorer verifiers. CE is an account-global service (the region
// parameter is ignored) and identifies everything by ARN: the anomaly
// monitor and its folded subscriptions (exported in the
// subscription_arns map keyed by subscription name), and the cost
// category. The absent signals are the service's typed
// UnknownMonitorException / UnknownSubscriptionException /
// ResourceNotFoundException, never string matches.
package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	cetypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	pkgerrors "github.com/pkg/errors"
)

// costAnomalyMonitorVerifier verifies an AwsCostAnomalyMonitor (and its
// folded subscriptions through the outputs path), keyed on the monitor
// ARN.
type costAnomalyMonitorVerifier struct{}

func (*costAnomalyMonitorVerifier) IDOutputKey() string { return "monitor_arn" }

func (*costAnomalyMonitorVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := anomalyMonitorExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscostanomalymonitor verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscostanomalymonitor %q not found after deploy", id)
	}
	return nil
}

func (*costAnomalyMonitorVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := anomalyMonitorExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscostanomalymonitor verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscostanomalymonitor %q still exists after destroy", id)
	}
	return nil
}

// VerifyExistsFromOutputs additionally walks the subscription_arns
// map: each folded subscription must answer GetAnomalySubscriptions.
func (v *costAnomalyMonitorVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	monitorArn, _ := outputs["monitor_arn"].(string)
	if err := v.VerifyExists(ctx, cfg, monitorArn, region); err != nil {
		return err
	}
	client := costexplorer.NewFromConfig(cfg)
	for name, subscriptionArn := range stringMapOutput(outputs["subscription_arns"]) {
		out, err := client.GetAnomalySubscriptions(ctx, &costexplorer.GetAnomalySubscriptionsInput{
			SubscriptionArnList: []string{subscriptionArn},
		})
		if err != nil {
			return pkgerrors.Wrapf(err, "awscostanomalymonitor subscription %q (%s) lookup failed", name, subscriptionArn)
		}
		if len(out.AnomalySubscriptions) == 0 {
			return pkgerrors.Errorf("awscostanomalymonitor subscription %q (%s) not found after deploy", name, subscriptionArn)
		}
	}
	return nil
}

// VerifyAbsentFromOutputs asserts the monitor and every folded
// subscription are gone.
func (v *costAnomalyMonitorVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	monitorArn, _ := outputs["monitor_arn"].(string)
	if err := v.VerifyAbsent(ctx, cfg, monitorArn, region); err != nil {
		return err
	}
	client := costexplorer.NewFromConfig(cfg)
	for name, subscriptionArn := range stringMapOutput(outputs["subscription_arns"]) {
		out, err := client.GetAnomalySubscriptions(ctx, &costexplorer.GetAnomalySubscriptionsInput{
			SubscriptionArnList: []string{subscriptionArn},
		})
		if err != nil {
			var unknown *cetypes.UnknownSubscriptionException
			if errors.As(err, &unknown) {
				continue
			}
			return pkgerrors.Wrapf(err, "awscostanomalymonitor subscription %q verify-absent failed", name)
		}
		if len(out.AnomalySubscriptions) > 0 {
			return pkgerrors.Errorf("awscostanomalymonitor subscription %q (%s) still exists after destroy", name, subscriptionArn)
		}
	}
	return nil
}

func anomalyMonitorExists(ctx context.Context, cfg aws.Config, monitorArn string) (bool, error) {
	out, err := costexplorer.NewFromConfig(cfg).GetAnomalyMonitors(ctx, &costexplorer.GetAnomalyMonitorsInput{
		MonitorArnList: []string{monitorArn},
	})
	if err != nil {
		var unknown *cetypes.UnknownMonitorException
		if errors.As(err, &unknown) {
			return false, nil
		}
		return false, err
	}
	return len(out.AnomalyMonitors) > 0, nil
}

// costCategoryVerifier verifies an AwsCostCategory via
// DescribeCostCategoryDefinition, keyed on the category ARN.
type costCategoryVerifier struct{}

func (*costCategoryVerifier) IDOutputKey() string { return "category_arn" }

func (*costCategoryVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := costCategoryExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscostcategory verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscostcategory %q not found after deploy", id)
	}
	return nil
}

func (*costCategoryVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := costCategoryExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscostcategory verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscostcategory %q still exists after destroy", id)
	}
	return nil
}

func costCategoryExists(ctx context.Context, cfg aws.Config, categoryArn string) (bool, error) {
	out, err := costexplorer.NewFromConfig(cfg).DescribeCostCategoryDefinition(ctx, &costexplorer.DescribeCostCategoryDefinitionInput{
		CostCategoryArn: aws.String(categoryArn),
	})
	if err != nil {
		var notFound *cetypes.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	// A deleted category keeps answering Describe until month end with
	// EffectiveEnd set -- treat an end-dated definition as absent so
	// the destroy assertion reflects the user-visible state.
	if out.CostCategory != nil && aws.ToString(out.CostCategory.EffectiveEnd) != "" {
		return false, nil
	}
	return true, nil
}
