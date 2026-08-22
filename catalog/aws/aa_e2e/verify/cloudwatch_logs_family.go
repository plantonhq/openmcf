package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	pkgerrors "github.com/pkg/errors"
)

func isLogsNotFound(err error) bool {
	var notFound *cloudwatchlogstypes.ResourceNotFoundException
	return pkgerrors.As(err, &notFound)
}

// --- AwsCloudwatchLogDelivery -----------------------------------------------

// logDeliveryVerifier verifies an AwsCloudwatchLogDelivery instance
// arm-for-arm from its outputs: the vended source by name, every owned
// destination by name, every delivery by its AWS-generated ID, and the
// cross-account destination by name. Either arm may be absent (the
// spec's at-least-one-arm contract) - empty outputs skip their checks.
type logDeliveryVerifier struct{}

func (*logDeliveryVerifier) IDOutputKey() string { return "source_name" }

func (v *logDeliveryVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	if id == "" {
		return nil
	}
	exists, err := logDeliverySourceExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchlogdelivery verify-exists failed for source %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscloudwatchlogdelivery source %q not found after deploy", id)
	}
	return nil
}

func (v *logDeliveryVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	if id == "" {
		return nil
	}
	exists, err := logDeliverySourceExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchlogdelivery verify-absent failed for source %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscloudwatchlogdelivery source %q still exists after destroy", id)
	}
	return nil
}

func (v *logDeliveryVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	client := cloudwatchlogs.NewFromConfig(cfg)

	if sourceName := stringOutput(outputs, "source_name"); sourceName != "" {
		if err := v.VerifyExists(ctx, cfg, sourceName, region); err != nil {
			return err
		}
	}
	for destinationName := range mapOutputValues(outputs, "destination_arns") {
		if _, err := client.GetDeliveryDestination(ctx, &cloudwatchlogs.GetDeliveryDestinationInput{
			Name: aws.String(destinationName),
		}); err != nil {
			return pkgerrors.Wrapf(err, "GetDeliveryDestination(%s)", destinationName)
		}
	}
	for deliveryName, deliveryId := range mapOutputValues(outputs, "delivery_ids") {
		if deliveryId == "" {
			return pkgerrors.Errorf("delivery %q carries an empty id in delivery_ids", deliveryName)
		}
		if _, err := client.GetDelivery(ctx, &cloudwatchlogs.GetDeliveryInput{Id: aws.String(deliveryId)}); err != nil {
			return pkgerrors.Wrapf(err, "GetDelivery(%s/%s)", deliveryName, deliveryId)
		}
	}
	if crossAccountName := stringOutput(outputs, "cross_account_destination_name"); crossAccountName != "" {
		exists, err := logDestinationExists(ctx, cfg, crossAccountName)
		if err != nil {
			return pkgerrors.Wrapf(err, "awscloudwatchlogdelivery verify-exists failed for cross-account destination %q", crossAccountName)
		}
		if !exists {
			return pkgerrors.Errorf("awscloudwatchlogdelivery cross-account destination %q not found after deploy", crossAccountName)
		}
	}
	return nil
}

func (v *logDeliveryVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	client := cloudwatchlogs.NewFromConfig(cfg)

	if sourceName := stringOutput(outputs, "source_name"); sourceName != "" {
		if err := v.VerifyAbsent(ctx, cfg, sourceName, region); err != nil {
			return err
		}
	}
	for destinationName := range mapOutputValues(outputs, "destination_arns") {
		_, err := client.GetDeliveryDestination(ctx, &cloudwatchlogs.GetDeliveryDestinationInput{
			Name: aws.String(destinationName),
		})
		if err == nil {
			return pkgerrors.Errorf("awscloudwatchlogdelivery destination %q still exists after destroy", destinationName)
		}
		if !isLogsNotFound(err) {
			return pkgerrors.Wrapf(err, "GetDeliveryDestination(%s)", destinationName)
		}
	}
	if crossAccountName := stringOutput(outputs, "cross_account_destination_name"); crossAccountName != "" {
		exists, err := logDestinationExists(ctx, cfg, crossAccountName)
		if err != nil {
			return pkgerrors.Wrapf(err, "awscloudwatchlogdelivery verify-absent failed for cross-account destination %q", crossAccountName)
		}
		if exists {
			return pkgerrors.Errorf("awscloudwatchlogdelivery cross-account destination %q still exists after destroy", crossAccountName)
		}
	}
	return nil
}

func logDeliverySourceExists(ctx context.Context, cfg aws.Config, sourceName string) (bool, error) {
	client := cloudwatchlogs.NewFromConfig(cfg)
	_, err := client.GetDeliverySource(ctx, &cloudwatchlogs.GetDeliverySourceInput{Name: aws.String(sourceName)})
	if err != nil {
		if isLogsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// logDestinationExists matches the legacy cross-account destination by
// exact name (DescribeDestinations matches by prefix, so the page is
// re-checked for equality).
func logDestinationExists(ctx context.Context, cfg aws.Config, destinationName string) (bool, error) {
	client := cloudwatchlogs.NewFromConfig(cfg)
	out, err := client.DescribeDestinations(ctx, &cloudwatchlogs.DescribeDestinationsInput{
		DestinationNamePrefix: aws.String(destinationName),
	})
	if err != nil {
		return false, err
	}
	for _, destination := range out.Destinations {
		if aws.ToString(destination.DestinationName) == destinationName {
			return true, nil
		}
	}
	return false, nil
}

// --- AwsCloudwatchLogAccountPolicy ------------------------------------------

// logAccountPolicyVerifier verifies an AwsCloudwatchLogAccountPolicy
// via DescribeAccountPolicies, keyed on the (policy_name, policy_type)
// pair from the outputs (the type scopes the describe; the name is
// exact-matched in the page).
type logAccountPolicyVerifier struct{}

func (*logAccountPolicyVerifier) IDOutputKey() string { return "policy_name" }

func (v *logAccountPolicyVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	return pkgerrors.New("awscloudwatchlogaccountpolicy verifies from outputs (needs policy_type)")
}

func (v *logAccountPolicyVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	return pkgerrors.New("awscloudwatchlogaccountpolicy verifies from outputs (needs policy_type)")
}

func (v *logAccountPolicyVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, _ string) error {
	policyName := stringOutput(outputs, "policy_name")
	policyType := stringOutput(outputs, "policy_type")
	exists, err := logAccountPolicyExists(ctx, cfg, policyName, policyType)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchlogaccountpolicy verify-exists failed for %q/%s", policyName, policyType)
	}
	if !exists {
		return pkgerrors.Errorf("awscloudwatchlogaccountpolicy %q (%s) not found after deploy", policyName, policyType)
	}
	return nil
}

func (v *logAccountPolicyVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, _ string) error {
	policyName := stringOutput(outputs, "policy_name")
	policyType := stringOutput(outputs, "policy_type")
	exists, err := logAccountPolicyExists(ctx, cfg, policyName, policyType)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchlogaccountpolicy verify-absent failed for %q/%s", policyName, policyType)
	}
	if exists {
		return pkgerrors.Errorf("awscloudwatchlogaccountpolicy %q (%s) still exists after destroy", policyName, policyType)
	}
	return nil
}

func logAccountPolicyExists(ctx context.Context, cfg aws.Config, policyName, policyType string) (bool, error) {
	client := cloudwatchlogs.NewFromConfig(cfg)
	out, err := client.DescribeAccountPolicies(ctx, &cloudwatchlogs.DescribeAccountPoliciesInput{
		PolicyType: cloudwatchlogstypes.PolicyType(policyType),
		PolicyName: aws.String(policyName),
	})
	if err != nil {
		if isLogsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	for _, policy := range out.AccountPolicies {
		if aws.ToString(policy.PolicyName) == policyName {
			return true, nil
		}
	}
	return false, nil
}

// --- AwsCloudwatchLogAnomalyDetector ----------------------------------------

// logAnomalyDetectorVerifier verifies an AwsCloudwatchLogAnomalyDetector
// via GetLogAnomalyDetector, keyed on the detector ARN. A deleted
// detector returns the typed ResourceNotFoundException.
type logAnomalyDetectorVerifier struct{}

func (*logAnomalyDetectorVerifier) IDOutputKey() string { return "anomaly_detector_arn" }

func (*logAnomalyDetectorVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := logAnomalyDetectorExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchloganomalydetector verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscloudwatchloganomalydetector %q not found after deploy", id)
	}
	return nil
}

func (*logAnomalyDetectorVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := logAnomalyDetectorExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchloganomalydetector verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscloudwatchloganomalydetector %q still exists after destroy", id)
	}
	return nil
}

func logAnomalyDetectorExists(ctx context.Context, cfg aws.Config, detectorArn string) (bool, error) {
	client := cloudwatchlogs.NewFromConfig(cfg)
	_, err := client.GetLogAnomalyDetector(ctx, &cloudwatchlogs.GetLogAnomalyDetectorInput{
		AnomalyDetectorArn: aws.String(detectorArn),
	})
	if err != nil {
		if isLogsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// --- AwsCloudwatchLogResourcePolicy -----------------------------------------

// logResourcePolicyVerifier verifies an AwsCloudwatchLogResourcePolicy
// via DescribeResourcePolicies, keyed on the policy_id output (the
// policy name for account scope; the target ARN would select the
// resource scope, which today's lanes do not exercise - the describe
// walks the account-scope page and exact-matches the name).
type logResourcePolicyVerifier struct{}

func (*logResourcePolicyVerifier) IDOutputKey() string { return "policy_id" }

func (*logResourcePolicyVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := logResourcePolicyExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchlogresourcepolicy verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscloudwatchlogresourcepolicy %q not found after deploy", id)
	}
	return nil
}

func (*logResourcePolicyVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := logResourcePolicyExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchlogresourcepolicy verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscloudwatchlogresourcepolicy %q still exists after destroy", id)
	}
	return nil
}

func logResourcePolicyExists(ctx context.Context, cfg aws.Config, policyName string) (bool, error) {
	client := cloudwatchlogs.NewFromConfig(cfg)
	input := &cloudwatchlogs.DescribeResourcePoliciesInput{}
	for {
		page, err := client.DescribeResourcePolicies(ctx, input)
		if err != nil {
			return false, err
		}
		for _, policy := range page.ResourcePolicies {
			if aws.ToString(policy.PolicyName) == policyName {
				return true, nil
			}
		}
		if page.NextToken == nil {
			return false, nil
		}
		input.NextToken = page.NextToken
	}
}
