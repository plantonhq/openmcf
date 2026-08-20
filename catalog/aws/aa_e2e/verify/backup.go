package verify

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/pkg/errors"
)

// The AWS Backup family's verifiers. Identity notes that shape them:
//
//   - vaults (both types) and the three standard satellites are
//     name-identified -- one DescribeBackupVault answers for either
//     vault arm;
//   - plans are the family's one UUID-identified surface; selections
//     verify through the full-outputs path (their AWS-generated ids
//     live in the selection_ids output map);
//   - frameworks / report plans / restore testing plans are
//     name-identified but export ARNs (the chart-useful output); the
//     name parses out of the ARN's last segment unambiguously because
//     these names FORBID hyphens (the segment is "name" or
//     "name-<uuid>", and the first '-' can only start the uuid);
//   - the settings singleton's delete is a NO-OP on both arms --
//     absent means the settings object still answers with the
//     last-applied values, never "gone" (the token-vault precedent).

func backupClient(cfg aws.Config, region string) *backup.Client {
	return backup.NewFromConfig(cfg, func(o *backup.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

// nameFromBackupArn extracts the resource name from a backup ARN's
// last segment ("name" or "name-<uuid>"; the name itself cannot
// contain '-' for the kinds this is used on).
func nameFromBackupArn(arn string) string {
	segment := arn[strings.LastIndex(arn, ":")+1:]
	if i := strings.Index(segment, "-"); i > 0 {
		return segment[:i]
	}
	return segment
}

func isBackupNotFound(err error) bool {
	return strings.Contains(err.Error(), "ResourceNotFoundException")
}

// backupVaultVerifier verifies AwsBackupVault via DescribeBackupVault,
// keyed on vault_name -- the import identity for BOTH vault arms.
type backupVaultVerifier struct{}

func (*backupVaultVerifier) IDOutputKey() string { return "vault_name" }

func (*backupVaultVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := backupClient(cfg, region).DescribeBackupVault(ctx, &backup.DescribeBackupVaultInput{
		BackupVaultName: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "DescribeBackupVault(%s)", id)
	}
	if out.BackupVaultName == nil || *out.BackupVaultName != id {
		return errors.Errorf("backup vault (%s) answered with the wrong name %v", id, out.BackupVaultName)
	}
	return nil
}

func (*backupVaultVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := backupClient(cfg, region).DescribeBackupVault(ctx, &backup.DescribeBackupVaultInput{
		BackupVaultName: aws.String(id),
	})
	if err != nil {
		if isBackupNotFound(err) {
			return nil
		}
		return errors.Wrapf(err, "DescribeBackupVault(%s)", id)
	}
	return errors.Errorf("backup vault (%s) still exists after destroy", id)
}

// backupPlanVerifier verifies AwsBackupPlan via GetBackupPlan, keyed
// on plan_id (the AWS-generated UUID). The full-outputs path
// additionally proves every folded selection from the selection_ids
// map.
type backupPlanVerifier struct{}

func (*backupPlanVerifier) IDOutputKey() string { return "plan_id" }

func (*backupPlanVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := backupClient(cfg, region).GetBackupPlan(ctx, &backup.GetBackupPlanInput{
		BackupPlanId: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "GetBackupPlan(%s)", id)
	}
	if out.DeletionDate != nil {
		return errors.Errorf("backup plan (%s) reports a deletion date after deploy", id)
	}
	if out.BackupPlan == nil || len(out.BackupPlan.Rules) == 0 {
		return errors.Errorf("backup plan (%s) reports no rules after deploy", id)
	}
	return nil
}

func (*backupPlanVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := backupClient(cfg, region).GetBackupPlan(ctx, &backup.GetBackupPlanInput{
		BackupPlanId: aws.String(id),
	})
	if err != nil {
		if isBackupNotFound(err) {
			return nil
		}
		return errors.Wrapf(err, "GetBackupPlan(%s)", id)
	}
	// AWS keeps deleted plans describable with a DeletionDate -- the
	// soft-deleted echo, not an orphan.
	if out.DeletionDate != nil {
		return nil
	}
	return errors.Errorf("backup plan (%s) still exists after destroy", id)
}

func (v *backupPlanVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	planId, _ := outputs["plan_id"].(string)
	if planId == "" {
		return errors.New("backup plan outputs carry no plan_id")
	}
	if err := v.VerifyExists(ctx, cfg, planId, region); err != nil {
		return err
	}
	// Every folded selection must answer by its AWS-generated id.
	if selectionIds, ok := outputs["selection_ids"].(map[string]interface{}); ok {
		for name, raw := range selectionIds {
			selectionId, _ := raw.(string)
			if selectionId == "" {
				return errors.Errorf("selection %q carries an empty id in selection_ids", name)
			}
			if _, err := backupClient(cfg, region).GetBackupSelection(ctx, &backup.GetBackupSelectionInput{
				BackupPlanId: aws.String(planId),
				SelectionId:  aws.String(selectionId),
			}); err != nil {
				return errors.Wrapf(err, "GetBackupSelection(%s/%s)", planId, selectionId)
			}
		}
	}
	return nil
}

func (v *backupPlanVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	planId, _ := outputs["plan_id"].(string)
	if planId == "" {
		return errors.New("backup plan outputs carry no plan_id")
	}
	// Selections cannot outlive the plan -- the plan's absence is the
	// family's absence.
	return v.VerifyAbsent(ctx, cfg, planId, region)
}

// backupFrameworkVerifier verifies AwsBackupFramework via
// DescribeFramework, keyed on framework_arn (the name parses out of
// the ARN). The deployment status assertion is the honest bar: a
// framework without an ACTIVE Config recorder deploys FAILED while the
// provider treats the apply as complete.
type backupFrameworkVerifier struct{}

func (*backupFrameworkVerifier) IDOutputKey() string { return "framework_arn" }

func (*backupFrameworkVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	name := nameFromBackupArn(id)
	out, err := backupClient(cfg, region).DescribeFramework(ctx, &backup.DescribeFrameworkInput{
		FrameworkName: aws.String(name),
	})
	if err != nil {
		return errors.Wrapf(err, "DescribeFramework(%s)", name)
	}
	if out.DeploymentStatus == nil {
		return errors.Errorf("framework (%s) reports no deployment status after deploy", name)
	}
	// FAILED is a completed apply at the provider but never a healthy
	// framework -- reject it loudly (the Config-recorder gap class).
	if strings.Contains(*out.DeploymentStatus, "FAILED") {
		return errors.Errorf("framework (%s) deployment status is %s after deploy (an ACTIVE Config recorder recording the backup types is the usual missing prerequisite)", name, *out.DeploymentStatus)
	}
	if len(out.FrameworkControls) == 0 {
		return errors.Errorf("framework (%s) reports no controls after deploy", name)
	}
	return nil
}

func (*backupFrameworkVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	name := nameFromBackupArn(id)
	_, err := backupClient(cfg, region).DescribeFramework(ctx, &backup.DescribeFrameworkInput{
		FrameworkName: aws.String(name),
	})
	if err != nil {
		if isBackupNotFound(err) {
			return nil
		}
		return errors.Wrapf(err, "DescribeFramework(%s)", name)
	}
	return errors.Errorf("framework (%s) still exists after destroy", name)
}

// backupReportPlanVerifier verifies AwsBackupReportPlan via
// DescribeReportPlan, keyed on report_plan_arn (the name parses out of
// the ARN).
type backupReportPlanVerifier struct{}

func (*backupReportPlanVerifier) IDOutputKey() string { return "report_plan_arn" }

func (*backupReportPlanVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	name := nameFromBackupArn(id)
	out, err := backupClient(cfg, region).DescribeReportPlan(ctx, &backup.DescribeReportPlanInput{
		ReportPlanName: aws.String(name),
	})
	if err != nil {
		return errors.Wrapf(err, "DescribeReportPlan(%s)", name)
	}
	if out.ReportPlan == nil || out.ReportPlan.ReportDeliveryChannel == nil || out.ReportPlan.ReportDeliveryChannel.S3BucketName == nil {
		return errors.Errorf("report plan (%s) reports no delivery channel after deploy", name)
	}
	if out.ReportPlan.DeploymentStatus != nil && strings.Contains(*out.ReportPlan.DeploymentStatus, "FAILED") {
		return errors.Errorf("report plan (%s) deployment status is %s after deploy", name, *out.ReportPlan.DeploymentStatus)
	}
	return nil
}

func (*backupReportPlanVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	name := nameFromBackupArn(id)
	_, err := backupClient(cfg, region).DescribeReportPlan(ctx, &backup.DescribeReportPlanInput{
		ReportPlanName: aws.String(name),
	})
	if err != nil {
		if isBackupNotFound(err) {
			return nil
		}
		return errors.Wrapf(err, "DescribeReportPlan(%s)", name)
	}
	return errors.Errorf("report plan (%s) still exists after destroy", name)
}

// restoreTestingPlanVerifier verifies AwsBackupRestoreTestingPlan via
// GetRestoreTestingPlan, keyed on restore_testing_plan_arn (the name
// parses out of the ARN).
type restoreTestingPlanVerifier struct{}

func (*restoreTestingPlanVerifier) IDOutputKey() string { return "restore_testing_plan_arn" }

func (*restoreTestingPlanVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	name := nameFromBackupArn(id)
	out, err := backupClient(cfg, region).GetRestoreTestingPlan(ctx, &backup.GetRestoreTestingPlanInput{
		RestoreTestingPlanName: aws.String(name),
	})
	if err != nil {
		return errors.Wrapf(err, "GetRestoreTestingPlan(%s)", name)
	}
	if out.RestoreTestingPlan == nil || out.RestoreTestingPlan.RecoveryPointSelection == nil {
		return errors.Errorf("restore testing plan (%s) reports no recovery point selection after deploy", name)
	}
	if out.RestoreTestingPlan.ScheduleExpression == nil || *out.RestoreTestingPlan.ScheduleExpression == "" {
		return errors.Errorf("restore testing plan (%s) reports no schedule after deploy", name)
	}
	return nil
}

func (*restoreTestingPlanVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	name := nameFromBackupArn(id)
	_, err := backupClient(cfg, region).GetRestoreTestingPlan(ctx, &backup.GetRestoreTestingPlanInput{
		RestoreTestingPlanName: aws.String(name),
	})
	if err != nil {
		if isBackupNotFound(err) {
			return nil
		}
		return errors.Wrapf(err, "GetRestoreTestingPlan(%s)", name)
	}
	return errors.Errorf("restore testing plan (%s) still exists after destroy", name)
}

// backupSettingsVerifier verifies AwsBackupSettings via
// DescribeRegionSettings, keyed on region. DESTROY IS A NO-OP on both
// arms: absence means the settings object still answers with the
// last-applied values -- never "gone" (the token-vault precedent).
type backupSettingsVerifier struct{}

func (*backupSettingsVerifier) IDOutputKey() string { return "region" }

func (*backupSettingsVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := backupClient(cfg, region).DescribeRegionSettings(ctx, &backup.DescribeRegionSettingsInput{})
	if err != nil {
		return errors.Wrap(err, "DescribeRegionSettings")
	}
	// The scenario's whole point is asserting opt-ins -- verify the
	// applied posture, not mere readability (the settings object
	// always exists). EBS=true is the scenario's declared floor.
	if optedIn, ok := out.ResourceTypeOptInPreference["EBS"]; !ok || !optedIn {
		return errors.Errorf("backup region settings (%s) do not show the declared EBS opt-in after deploy", id)
	}
	return nil
}

func (*backupSettingsVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	// Provider deletes are no-ops on both arms: the last-applied
	// settings REMAIN after destroy. Asserting disappearance would fail
	// every honest run -- the contract is that the object still
	// answers.
	if _, err := backupClient(cfg, region).DescribeRegionSettings(ctx, &backup.DescribeRegionSettingsInput{}); err != nil {
		return errors.Wrap(err, "DescribeRegionSettings after destroy (the settings must persist -- delete is a no-op)")
	}
	return nil
}
