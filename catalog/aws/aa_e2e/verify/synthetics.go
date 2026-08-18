package verify

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/synthetics"
	syntheticstypes "github.com/aws/aws-sdk-go-v2/service/synthetics/types"
	pkgerrors "github.com/pkg/errors"
)

func isSyntheticsNotFound(err error) bool {
	var notFound *syntheticstypes.ResourceNotFoundException
	return pkgerrors.As(err, &notFound)
}

// syntheticsVerifier verifies an AwsCloudwatchSynthetics instance: the
// canary via GetCanary (keyed on canary_name) and every owned group via
// GetGroup (walked from the group_ids output map), so groups-only
// instances verify their groups and canary instances verify both arms.
// Deletion is synchronous for both resource types; the typed
// ResourceNotFoundException is the "absent" signal.
type syntheticsVerifier struct{}

func (*syntheticsVerifier) IDOutputKey() string { return "canary_name" }

func (v *syntheticsVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	if id == "" {
		// Groups-only instance: the outputs path covers the groups.
		return nil
	}
	exists, err := syntheticsCanaryExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchsynthetics verify-exists failed for canary %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscloudwatchsynthetics canary %q not found after deploy", id)
	}
	return nil
}

func (v *syntheticsVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	if id == "" {
		return nil
	}
	exists, err := syntheticsCanaryExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchsynthetics verify-absent failed for canary %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscloudwatchsynthetics canary %q still exists after destroy", id)
	}
	return nil
}

func (v *syntheticsVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	if err := v.VerifyExists(ctx, cfg, stringOutput(outputs, "canary_name"), region); err != nil {
		return err
	}
	for _, groupName := range mapOutputKeys(outputs, "group_ids") {
		exists, err := syntheticsGroupExists(ctx, cfg, groupName)
		if err != nil {
			return pkgerrors.Wrapf(err, "awscloudwatchsynthetics verify-exists failed for group %q", groupName)
		}
		if !exists {
			return pkgerrors.Errorf("awscloudwatchsynthetics group %q not found after deploy", groupName)
		}
	}
	return nil
}

func (v *syntheticsVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	if err := v.VerifyAbsent(ctx, cfg, stringOutput(outputs, "canary_name"), region); err != nil {
		return err
	}
	for _, groupName := range mapOutputKeys(outputs, "group_ids") {
		exists, err := syntheticsGroupExists(ctx, cfg, groupName)
		if err != nil {
			return pkgerrors.Wrapf(err, "awscloudwatchsynthetics verify-absent failed for group %q", groupName)
		}
		if exists {
			return pkgerrors.Errorf("awscloudwatchsynthetics group %q still exists after destroy", groupName)
		}
	}
	return nil
}

func syntheticsCanaryExists(ctx context.Context, cfg aws.Config, canaryName string) (bool, error) {
	client := synthetics.NewFromConfig(cfg)
	_, err := client.GetCanary(ctx, &synthetics.GetCanaryInput{Name: aws.String(canaryName)})
	if err != nil {
		if isSyntheticsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func syntheticsGroupExists(ctx context.Context, cfg aws.Config, groupName string) (bool, error) {
	client := synthetics.NewFromConfig(cfg)
	_, err := client.GetGroup(ctx, &synthetics.GetGroupInput{GroupIdentifier: aws.String(groupName)})
	if err != nil {
		if isSyntheticsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// stringOutput reads a string-valued stack output, tolerating absence.
func stringOutput(outputs map[string]interface{}, key string) string {
	if value, ok := outputs[key].(string); ok {
		return value
	}
	return ""
}

// mapOutputKeys reads the keys of a map-valued stack output, tolerating
// absence.
func mapOutputKeys(outputs map[string]interface{}, key string) []string {
	raw, ok := outputs[key].(map[string]interface{})
	if !ok {
		return nil
	}
	keys := make([]string, 0, len(raw))
	for k := range raw {
		keys = append(keys, k)
	}
	return keys
}

// mapOutputValues reads the string values of a map-valued stack output
// keyed by entry name, tolerating absence.
func mapOutputValues(outputs map[string]interface{}, key string) map[string]string {
	raw, ok := outputs[key].(map[string]interface{})
	if !ok {
		return nil
	}
	values := make(map[string]string, len(raw))
	for k, v := range raw {
		values[k] = fmt.Sprintf("%v", v)
	}
	return values
}
