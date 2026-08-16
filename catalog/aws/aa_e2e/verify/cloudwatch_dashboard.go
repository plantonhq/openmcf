package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	pkgerrors "github.com/pkg/errors"
)

func isDashboardNotFound(err error) bool {
	var notFound *cloudwatchtypes.DashboardNotFoundError
	return pkgerrors.As(err, &notFound)
}

// cloudwatchDashboardVerifier verifies an AwsCloudwatchDashboard via
// GetDashboard, keyed on the dashboard name (the resource's identity
// and import ID). A deleted dashboard returns the typed
// DashboardNotFoundError, which is the "absent" signal.
type cloudwatchDashboardVerifier struct{}

func (*cloudwatchDashboardVerifier) IDOutputKey() string { return "dashboard_name" }

func (*cloudwatchDashboardVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := cloudwatchDashboardExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchdashboard verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscloudwatchdashboard %q not found after deploy", id)
	}
	return nil
}

func (*cloudwatchDashboardVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := cloudwatchDashboardExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchdashboard verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscloudwatchdashboard %q still exists after destroy", id)
	}
	return nil
}

func cloudwatchDashboardExists(ctx context.Context, cfg aws.Config, dashboardName string) (bool, error) {
	client := cloudwatch.NewFromConfig(cfg)
	_, err := client.GetDashboard(ctx, &cloudwatch.GetDashboardInput{DashboardName: aws.String(dashboardName)})
	if err != nil {
		if isDashboardNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
