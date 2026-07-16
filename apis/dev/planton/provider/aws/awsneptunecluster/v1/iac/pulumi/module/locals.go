package module

import (
	"strconv"
	"strings"

	awsneptuneclusterv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsneptunecluster/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsNeptuneCluster *awsneptuneclusterv1.AwsNeptuneCluster

	// ClusterIdentifier is metadata.name -- create-only in AWS, and the
	// basis both engines share so a manifest deploys identically on either.
	ClusterIdentifier string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsneptuneclusterv1.AwsNeptuneClusterStackInput) *Locals {
	locals := &Locals{}
	locals.AwsNeptuneCluster = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.ClusterIdentifier = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsNeptuneCluster.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}

// engineFamily derives the cluster parameter-group family from the pinned
// engine_version (inline parameters require a pinned version,
// CEL-enforced, so this never sees an empty version): "1.4.5.1" ->
// neptune1.4. Neptune families are keyed by major.minor -- AWS's own
// family naming, not a convention of ours.
func engineFamily(engineVersion string) string {
	parts := strings.Split(engineVersion, ".")
	if len(parts) >= 2 {
		return "neptune" + parts[0] + "." + parts[1]
	}
	return "neptune" + parts[0]
}
