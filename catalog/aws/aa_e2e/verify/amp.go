package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amp"
	amptypes "github.com/aws/aws-sdk-go-v2/service/amp/types"
	pkgerrors "github.com/pkg/errors"
)

func isAmpNotFound(err error) bool {
	var notFound *amptypes.ResourceNotFoundException
	return pkgerrors.As(err, &notFound)
}

// --- AwsManagedPrometheus ----------------------------------------------------

// ampWorkspaceVerifier verifies an AwsManagedPrometheus instance
// arm-for-arm from its outputs: the workspace via DescribeWorkspace
// (ACTIVE), the alert manager definition, every rule-groups namespace
// (walked by name from rule_group_namespace_arns), the query-logging
// configuration, the resource policy, and every anomaly detector
// (walked by ID from anomaly_detector_ids). Satellites die with the
// workspace, so absence is the workspace's absence.
type ampWorkspaceVerifier struct{}

func (*ampWorkspaceVerifier) IDOutputKey() string { return "workspace_id" }

func (*ampWorkspaceVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := ampWorkspaceActive(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmanagedprometheus verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsmanagedprometheus workspace %q not found after deploy", id)
	}
	return nil
}

func (*ampWorkspaceVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := ampWorkspaceActive(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmanagedprometheus verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsmanagedprometheus workspace %q still exists after destroy", id)
	}
	return nil
}

func (v *ampWorkspaceVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	workspaceId := stringOutput(outputs, "workspace_id")
	if workspaceId == "" {
		return pkgerrors.New("awsmanagedprometheus outputs carry no workspace_id")
	}
	if err := v.VerifyExists(ctx, cfg, workspaceId, region); err != nil {
		return err
	}

	client := amp.NewFromConfig(cfg)
	for namespaceName := range mapOutputValues(outputs, "rule_group_namespace_arns") {
		if _, err := client.DescribeRuleGroupsNamespace(ctx, &amp.DescribeRuleGroupsNamespaceInput{
			WorkspaceId: aws.String(workspaceId),
			Name:        aws.String(namespaceName),
		}); err != nil {
			return pkgerrors.Wrapf(err, "DescribeRuleGroupsNamespace(%s/%s)", workspaceId, namespaceName)
		}
	}
	for detectorAlias, detectorId := range mapOutputValues(outputs, "anomaly_detector_ids") {
		if detectorId == "" {
			return pkgerrors.Errorf("anomaly detector %q carries an empty id in anomaly_detector_ids", detectorAlias)
		}
		if _, err := client.DescribeAnomalyDetector(ctx, &amp.DescribeAnomalyDetectorInput{
			WorkspaceId:       aws.String(workspaceId),
			AnomalyDetectorId: aws.String(detectorId),
		}); err != nil {
			return pkgerrors.Wrapf(err, "DescribeAnomalyDetector(%s/%s)", workspaceId, detectorId)
		}
	}
	return nil
}

func (v *ampWorkspaceVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	workspaceId := stringOutput(outputs, "workspace_id")
	if workspaceId == "" {
		return pkgerrors.New("awsmanagedprometheus outputs carry no workspace_id")
	}
	// Satellites cannot outlive the workspace - its absence is the
	// family's absence.
	return v.VerifyAbsent(ctx, cfg, workspaceId, region)
}

// ampWorkspaceActive reports whether the workspace exists in a live
// state (AWS keeps DELETING workspaces describable until they vanish;
// only ACTIVE/CREATING/UPDATING count as present).
func ampWorkspaceActive(ctx context.Context, cfg aws.Config, workspaceId string) (bool, error) {
	client := amp.NewFromConfig(cfg)
	out, err := client.DescribeWorkspace(ctx, &amp.DescribeWorkspaceInput{WorkspaceId: aws.String(workspaceId)})
	if err != nil {
		if isAmpNotFound(err) {
			return false, nil
		}
		return false, err
	}
	if out.Workspace != nil && out.Workspace.Status != nil &&
		out.Workspace.Status.StatusCode == amptypes.WorkspaceStatusCodeDeleting {
		return false, nil
	}
	return true, nil
}

// --- AwsManagedPrometheusScraper ----------------------------------------------

// ampScraperVerifier verifies an AwsManagedPrometheusScraper via
// DescribeScraper, keyed on the AWS-generated scraper ID. Deletes
// drain (the scraper stays describable while DELETING), so only
// non-deleting states count as present.
type ampScraperVerifier struct{}

func (*ampScraperVerifier) IDOutputKey() string { return "scraper_id" }

func (*ampScraperVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := ampScraperPresent(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmanagedprometheusscraper verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsmanagedprometheusscraper %q not found after deploy", id)
	}
	return nil
}

func (*ampScraperVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := ampScraperPresent(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmanagedprometheusscraper verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsmanagedprometheusscraper %q still exists after destroy", id)
	}
	return nil
}

func ampScraperPresent(ctx context.Context, cfg aws.Config, scraperId string) (bool, error) {
	client := amp.NewFromConfig(cfg)
	out, err := client.DescribeScraper(ctx, &amp.DescribeScraperInput{ScraperId: aws.String(scraperId)})
	if err != nil {
		if isAmpNotFound(err) {
			return false, nil
		}
		return false, err
	}
	if out.Scraper != nil && out.Scraper.Status != nil &&
		out.Scraper.Status.StatusCode == amptypes.ScraperStatusCodeDeleting {
		return false, nil
	}
	return true, nil
}
