package module

import (
	"strconv"
	"strings"

	awsrdsclusterv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsrdscluster/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsRdsCluster *awsrdsclusterv1.AwsRdsCluster

	// ClusterIdentifier is metadata.name -- create-only in AWS, and the
	// basis both engines share so a manifest deploys identically on either.
	ClusterIdentifier string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsrdsclusterv1.AwsRdsClusterStackInput) *Locals {
	locals := &Locals{}
	locals.AwsRdsCluster = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.ClusterIdentifier = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsRdsCluster.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}

// engineFamily derives the cluster parameter-group family from the pinned
// engine + engine_version (inline parameters require a pinned version,
// CEL-enforced, so this never sees an empty version):
//
//	aurora-postgresql 16.4                -> aurora-postgresql16
//	postgres          16.4                -> postgres16
//	aurora-mysql      8.0.mysql_aurora... -> aurora-mysql8.0
//	mysql             8.0.39              -> mysql8.0
//
// PostgreSQL families are keyed by major version; MySQL families by
// major.minor -- AWS's own family naming, not a convention of ours.
func engineFamily(engine, engineVersion string) string {
	parts := strings.Split(engineVersion, ".")
	switch engine {
	case "aurora-postgresql", "postgres":
		return engine + parts[0]
	default:
		if len(parts) >= 2 {
			return engine + parts[0] + "." + parts[1]
		}
		return engine + parts[0]
	}
}
