package module

import (
	"github.com/pkg/errors"
	awsredshiftserverlessworkgroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsredshiftserverlessworkgroup/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources provisions the Redshift Serverless workgroup -- the compute
// plane of the serverless warehouse. The workgroup composes onto its
// neighbors instead of embedding them: the namespace it serves, the
// subnets it places compute in, and the security groups on its endpoint
// all attach by reference, and warehouse ingress rules live on the
// referenced AwsSecurityGroup nodes -- this module never creates or
// mutates a resource that deserves to be its own node.
func Resources(ctx *pulumi.Context, stackInput *awsredshiftserverlessworkgroupv1alpha1.AwsRedshiftServerlessWorkgroupStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder,
	// which resolves the right credential mechanism (static keys, keyless
	// web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsRedshiftServerlessWorkgroup.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdWorkgroup, err := workgroup(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create Redshift Serverless workgroup")
	}

	ctx.Export(OpWorkgroupName, createdWorkgroup.WorkgroupName)
	ctx.Export(OpWorkgroupId, createdWorkgroup.WorkgroupId)
	ctx.Export(OpArn, createdWorkgroup.Arn)
	ctx.Export(OpPort, createdWorkgroup.Port)

	// The connection hostname lives on the workgroup's endpoint list
	// (exactly one endpoint once the workgroup is available). Index and
	// Elem both resolve to zero values when the endpoint is not yet
	// known, so the export shape is stable without an ApplyT applier.
	ctx.Export(OpEndpointAddress, createdWorkgroup.Endpoints.Index(pulumi.Int(0)).Address().Elem())

	return nil
}
