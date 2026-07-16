package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/efs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func policies(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, fs *efs.FileSystem) error {
	spec := locals.AwsElasticFileSystem.Spec

	// Backup policy — the resource has no true delete (removal PUTs status
	// DISABLED), so it is only materialized when backups are enabled; absent
	// means AWS's default (disabled) and a toggle-off flows through the same
	// resource lifecycle.
	if spec.BackupEnabled {
		_, err := efs.NewBackupPolicy(ctx, "backup-policy", &efs.BackupPolicyArgs{
			FileSystemId: fs.ID(),
			BackupPolicy: &efs.BackupPolicyBackupPolicyArgs{
				Status: pulumi.String("ENABLED"),
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "failed to create efs backup policy")
		}
	}

	// File system resource policy — the spec models it as a Struct, so it is
	// serialized to the JSON document the provider expects.
	if spec.Policy != nil {
		policyJSON, err := json.Marshal(spec.Policy.AsMap())
		if err != nil {
			return errors.Wrap(err, "failed to serialize efs file system policy to JSON")
		}

		policyArgs := &efs.FileSystemPolicyArgs{
			FileSystemId: fs.ID(),
			Policy:       pulumi.String(string(policyJSON)),
		}

		// Only set when the user deliberately opts out of AWS's lockout
		// safety check (a policy that denies the deploying principal future
		// PutFileSystemPolicy calls is otherwise rejected).
		if spec.BypassPolicyLockoutSafetyCheck {
			policyArgs.BypassPolicyLockoutSafetyCheck = pulumi.Bool(true)
		}

		_, err = efs.NewFileSystemPolicy(ctx, "fs-policy", policyArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "failed to create efs file system policy")
		}
	}

	return nil
}
