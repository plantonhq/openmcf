package module

import (
	"fmt"

	"github.com/pkg/errors"
	kubernetespostgresv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetespostgres/v1"
	zalandov1 "github.com/plantonhq/openmcf/pkg/kubernetes/kubernetestypes/zalandooperator/kubernetes/acid/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildRestoreConfig generates Zalando operator's spec:standby configuration and STANDBY_* environment
// variables for cross-cluster disaster recovery using the Standby-then-Promote pattern.
//
// When restore.enabled=true:
//   - Returns a populated standby block with s3_wal_path
//   - Returns STANDBY_* environment variables for R2 access
//   - Database bootstraps as read-only standby from R2 backups
//
// When restore.enabled=false or restore=nil:
//   - Returns nil for both standby and env vars
//   - Database runs as normal read-write primary
//   - If previously in standby mode, triggers promotion
//
// Parameters:
//   - restoreConfig: User's restore configuration from API spec
//   - operatorBucketName: Fallback bucket from operator-level config (optional)
//
// Returns:
//   - standby block for Zalando manifest (or nil)
//   - STANDBY_* environment variables (or nil)
//   - error if configuration is invalid
func buildRestoreConfig(
	restoreConfig *kubernetespostgresv1.KubernetesPostgresRestoreConfig,
	operatorBucketName string,
) (*zalandov1.PostgresqlSpecStandbyArgs, []pulumi.MapInput, error) {

	// If restore is not configured or disabled, return nil (normal primary mode)
	if restoreConfig == nil || !restoreConfig.Enabled {
		return nil, nil, nil
	}

	// Validate s3_path is provided
	if restoreConfig.S3Path == "" {
		return nil, nil, errors.New("restore.s3_path is required when restore.enabled=true")
	}

	// Determine bucket name (per-database overrides operator-level)
	var bucketName string
	if restoreConfig.BucketName != nil && *restoreConfig.BucketName != "" {
		// Use per-database bucket
		bucketName = *restoreConfig.BucketName
	} else if operatorBucketName != "" {
		// Fallback to operator-level bucket
		bucketName = operatorBucketName
	} else {
		return nil, nil, errors.New("restore.bucket_name is required when restore.enabled=true (not found in database or operator config)")
	}

	// Construct full S3 path for Zalando's spec:standby.s3_wal_path
	// Format: s3://bucket-name/path/to/backups
	fullS3Path := fmt.Sprintf("s3://%s/%s", bucketName, restoreConfig.S3Path)

	// Create Zalando standby block
	standbyBlock := &zalandov1.PostgresqlSpecStandbyArgs{
		S3_wal_path: pulumi.String(fullS3Path),
	}

	// Build STANDBY_* environment variables for R2 access
	var envVars []pulumi.MapInput

	if restoreConfig.R2Config != nil {
		// Construct R2 endpoint URL
		r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com",
			restoreConfig.R2Config.CloudflareAccountId)

		// STANDBY_* env vars are used by Spilo/Patroni during standby bootstrap
		// These are distinct from WALG_* (ongoing backups) and CLONE_* (clone operations)
		envVars = []pulumi.MapInput{
			pulumi.Map{
				"name":  pulumi.String("STANDBY_AWS_ENDPOINT"),
				"value": pulumi.String(r2Endpoint),
			},
			pulumi.Map{
				"name":  pulumi.String("STANDBY_AWS_FORCE_PATH_STYLE"),
				"value": pulumi.String("true"),
			},
			pulumi.Map{
				"name":  pulumi.String("STANDBY_AWS_ACCESS_KEY_ID"),
				"value": pulumi.String(restoreConfig.R2Config.AccessKeyId),
			},
			pulumi.Map{
				"name":  pulumi.String("STANDBY_AWS_SECRET_ACCESS_KEY"),
				"value": pulumi.String(restoreConfig.R2Config.SecretAccessKey),
			},
		}
	}

	return standbyBlock, envVars, nil
}
