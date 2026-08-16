package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/pkg/errors"
)

// The SSM family's verifiers. Identity notes that shape them:
//
//   - parameters are name-identified (the explicit spec field --
//     hierarchical paths metadata.name cannot carry); GetParameter
//     answers without decryption so SecureString lanes never pull
//     plaintext into logs;
//   - documents are name-identified; DeleteDocument is asynchronous,
//     so VerifyAbsent tolerates a still-deleting read the same way a
//     not-found reads (the scenario name is run-scoped for the same
//     reason);
//   - associations are UUID-identified (the document name is NOT the
//     identity at AWS);
//   - maintenance windows are "mw-..."-identified; the folded target
//     and task registrations verify through the full-outputs path
//     (their AWS-generated ids live in the target_ids/task_ids output
//     maps);
//   - patch baselines are "pb-..."-identified; the folded patch groups
//     and the default designation are proof-lane evidence (their
//     identities compose spec entries with the baseline id -- see the
//     import map), not output-addressable surfaces.
func ssmClient(cfg aws.Config, region string) *ssm.Client {
	return ssm.NewFromConfig(cfg, func(o *ssm.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

// ssmParameterVerifier verifies AwsSsmParameter via GetParameter,
// keyed on parameter_name.
type ssmParameterVerifier struct{}

func (*ssmParameterVerifier) IDOutputKey() string { return "parameter_name" }

func (*ssmParameterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	// WithDecryption stays false: existence is the assertion, and
	// SecureString plaintext must never reach lane logs.
	_, err := ssmClient(cfg, region).GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "expected SSM parameter %s to exist", id)
	}
	return nil
}

func (*ssmParameterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := ssmClient(cfg, region).GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(id),
	})
	if err == nil {
		return errors.Errorf("expected SSM parameter %s to be gone, but it still exists", id)
	}
	var notFound *ssmtypes.ParameterNotFound
	if errors.As(err, &notFound) {
		return nil
	}
	return errors.Wrapf(err, "unexpected error verifying SSM parameter %s absence", id)
}

// ssmDocumentVerifier verifies AwsSsmDocument via DescribeDocument,
// keyed on document_name.
type ssmDocumentVerifier struct{}

func (*ssmDocumentVerifier) IDOutputKey() string { return "document_name" }

func (*ssmDocumentVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := ssmClient(cfg, region).DescribeDocument(ctx, &ssm.DescribeDocumentInput{
		Name: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "expected SSM document %s to exist", id)
	}
	// A document that reached Failed never carries runnable content --
	// existence alone is not the honest bar.
	if out.Document != nil && out.Document.Status == ssmtypes.DocumentStatusFailed {
		return errors.Errorf("SSM document %s exists but is in status Failed", id)
	}
	return nil
}

func (*ssmDocumentVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := ssmClient(cfg, region).DescribeDocument(ctx, &ssm.DescribeDocumentInput{
		Name: aws.String(id),
	})
	if err != nil {
		var invalid *ssmtypes.InvalidDocument
		if errors.As(err, &invalid) {
			return nil
		}
		return errors.Wrapf(err, "unexpected error verifying SSM document %s absence", id)
	}
	// DeleteDocument is asynchronous -- a still-deleting document reads
	// as Deleting and counts as gone (the config-rule DELETING
	// precedent).
	if out.Document != nil && out.Document.Status == ssmtypes.DocumentStatusDeleting {
		return nil
	}
	return errors.Errorf("expected SSM document %s to be gone, but it still exists", id)
}

// ssmAssociationVerifier verifies AwsSsmAssociation via
// DescribeAssociation, keyed on association_id (the AWS-generated
// UUID -- the document name is NOT the identity).
type ssmAssociationVerifier struct{}

func (*ssmAssociationVerifier) IDOutputKey() string { return "association_id" }

func (*ssmAssociationVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := ssmClient(cfg, region).DescribeAssociation(ctx, &ssm.DescribeAssociationInput{
		AssociationId: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "expected SSM association %s to exist", id)
	}
	return nil
}

func (*ssmAssociationVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := ssmClient(cfg, region).DescribeAssociation(ctx, &ssm.DescribeAssociationInput{
		AssociationId: aws.String(id),
	})
	if err == nil {
		return errors.Errorf("expected SSM association %s to be gone, but it still exists", id)
	}
	var notFound *ssmtypes.AssociationDoesNotExist
	if errors.As(err, &notFound) {
		return nil
	}
	return errors.Wrapf(err, "unexpected error verifying SSM association %s absence", id)
}

// ssmMaintenanceWindowVerifier verifies AwsSsmMaintenanceWindow via
// GetMaintenanceWindow, keyed on window_id; the folded target and task
// registrations verify through the full-outputs path (their ids live
// in the target_ids/task_ids output maps).
type ssmMaintenanceWindowVerifier struct{}

func (*ssmMaintenanceWindowVerifier) IDOutputKey() string { return "window_id" }

func (*ssmMaintenanceWindowVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := ssmClient(cfg, region).GetMaintenanceWindow(ctx, &ssm.GetMaintenanceWindowInput{
		WindowId: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "expected SSM maintenance window %s to exist", id)
	}
	return nil
}

func (*ssmMaintenanceWindowVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := ssmClient(cfg, region).GetMaintenanceWindow(ctx, &ssm.GetMaintenanceWindowInput{
		WindowId: aws.String(id),
	})
	if err == nil {
		return errors.Errorf("expected SSM maintenance window %s to be gone, but it still exists", id)
	}
	var notFound *ssmtypes.DoesNotExistException
	if errors.As(err, &notFound) {
		return nil
	}
	return errors.Wrapf(err, "unexpected error verifying SSM maintenance window %s absence", id)
}

func (v *ssmMaintenanceWindowVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	windowId, _ := outputs["window_id"].(string)
	if windowId == "" {
		return errors.Errorf("maintenance window outputs carry no window_id")
	}
	if err := v.VerifyExists(ctx, cfg, windowId, region); err != nil {
		return err
	}
	client := ssmClient(cfg, region)
	// Every folded target registration must answer by its AWS-generated
	// id (DescribeMaintenanceWindowTargets filtered to the one id must
	// return exactly it).
	if targetIds, ok := outputs["target_ids"].(map[string]interface{}); ok {
		for name, raw := range targetIds {
			targetId, _ := raw.(string)
			if targetId == "" {
				return errors.Errorf("target %q carries an empty id in target_ids", name)
			}
			out, err := client.DescribeMaintenanceWindowTargets(ctx, &ssm.DescribeMaintenanceWindowTargetsInput{
				WindowId: aws.String(windowId),
				Filters: []ssmtypes.MaintenanceWindowFilter{{
					Key:    aws.String("WindowTargetId"),
					Values: []string{targetId},
				}},
			})
			if err != nil {
				return errors.Wrapf(err, "DescribeMaintenanceWindowTargets(%s/%s)", windowId, targetId)
			}
			if len(out.Targets) != 1 {
				return errors.Errorf("expected window target %s (%s) to exist, found %d matches", name, targetId, len(out.Targets))
			}
		}
	}
	// Every folded task must answer by its AWS-generated id.
	if taskIds, ok := outputs["task_ids"].(map[string]interface{}); ok {
		for name, raw := range taskIds {
			taskId, _ := raw.(string)
			if taskId == "" {
				return errors.Errorf("task %q carries an empty id in task_ids", name)
			}
			if _, err := client.GetMaintenanceWindowTask(ctx, &ssm.GetMaintenanceWindowTaskInput{
				WindowId:     aws.String(windowId),
				WindowTaskId: aws.String(taskId),
			}); err != nil {
				return errors.Wrapf(err, "GetMaintenanceWindowTask(%s/%s)", windowId, taskId)
			}
		}
	}
	return nil
}

func (v *ssmMaintenanceWindowVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	windowId, _ := outputs["window_id"].(string)
	if windowId == "" {
		return errors.Errorf("maintenance window outputs carry no window_id")
	}
	// The window's absence implies its registrations' absence -- AWS
	// deletes targets and tasks with the window.
	return v.VerifyAbsent(ctx, cfg, windowId, region)
}

// ssmPatchBaselineVerifier verifies AwsSsmPatchBaseline via
// GetPatchBaseline, keyed on baseline_id. The folded patch groups and
// the default designation compose SPEC entries with the baseline id
// (no output map), so their live evidence rides the proof lane's
// watchers rather than this verifier.
type ssmPatchBaselineVerifier struct{}

func (*ssmPatchBaselineVerifier) IDOutputKey() string { return "baseline_id" }

func (*ssmPatchBaselineVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := ssmClient(cfg, region).GetPatchBaseline(ctx, &ssm.GetPatchBaselineInput{
		BaselineId: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "expected SSM patch baseline %s to exist", id)
	}
	return nil
}

func (*ssmPatchBaselineVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := ssmClient(cfg, region).GetPatchBaseline(ctx, &ssm.GetPatchBaselineInput{
		BaselineId: aws.String(id),
	})
	if err == nil {
		return errors.Errorf("expected SSM patch baseline %s to be gone, but it still exists", id)
	}
	var notFound *ssmtypes.DoesNotExistException
	if errors.As(err, &notFound) {
		return nil
	}
	return errors.Wrapf(err, "unexpected error verifying SSM patch baseline %s absence", id)
}
