package module

import (
	"github.com/pkg/errors"
	awsroute53healthcheckv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsroute53healthcheck/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources creates one Route 53 health check. The spec's CEL rules guarantee
// the arguments already match the chosen check_type (endpoint probing,
// calculated aggregation, CloudWatch mirroring, or recovery-control
// mirroring), so this module maps fields 1:1 without re-validating.
func Resources(ctx *pulumi.Context, stackInput *awsroute53healthcheckv1alpha1.AwsRoute53HealthCheckStackInput) error {
	locals := initializeLocals(ctx, stackInput)
	spec := locals.AwsRoute53HealthCheck.Spec

	// Build the AWS provider from the stack input via the shared builder, which
	// resolves the right credential mechanism (static keys, keyless web
	// identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	// The monitoring model (ForceNew).
	args := &route53.HealthCheckArgs{
		Type: pulumi.String(spec.CheckType),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	// --- Endpoint checks (HTTP / HTTPS / *_STR_MATCH / TCP) ------------------
	// fqdn alone: Route 53 resolves and probes it (and sends it as the Host
	// header). ip_address alone or with fqdn: the probe goes to the IP.
	// Unset optionals stay nil so the provider applies its own defaults
	// (port 80/443, resource path "/", all checker regions).
	if spec.Fqdn != "" {
		args.Fqdn = pulumi.String(spec.Fqdn)
	}
	if spec.IpAddress != "" {
		args.IpAddress = pulumi.String(spec.IpAddress)
	}
	if spec.Port != 0 {
		args.Port = pulumi.Int(int(spec.Port))
	}
	if spec.ResourcePath != "" {
		args.ResourcePath = pulumi.String(spec.ResourcePath)
	}
	if spec.SearchString != "" {
		args.SearchString = pulumi.String(spec.SearchString)
	}
	// request_interval and measure_latency are create-time (ForceNew).
	if spec.RequestInterval != nil {
		args.RequestInterval = pulumi.Int(int(spec.GetRequestInterval()))
	}
	if spec.FailureThreshold != nil {
		args.FailureThreshold = pulumi.Int(int(spec.GetFailureThreshold()))
	}
	if spec.MeasureLatency {
		args.MeasureLatency = pulumi.Bool(true)
	}
	if spec.EnableSni != nil {
		args.EnableSni = pulumi.Bool(spec.GetEnableSni())
	}
	if len(spec.Regions) > 0 {
		args.Regions = pulumi.ToStringArray(spec.Regions)
	}

	// --- State shaping (any type) ---------------------------------------------
	if spec.InvertHealthcheck {
		args.InvertHealthcheck = pulumi.Bool(true)
	}
	if spec.Disabled {
		args.Disabled = pulumi.Bool(true)
	}

	// --- CALCULATED: aggregate child checks -----------------------------------
	// A zero threshold means "all children must be healthy" (the AWS default
	// when the parameter is omitted).
	if len(spec.ChildHealthChecks) > 0 {
		childIds := make([]string, 0, len(spec.ChildHealthChecks))
		for _, child := range spec.ChildHealthChecks {
			childIds = append(childIds, child.GetValue())
		}
		args.ChildHealthchecks = pulumi.ToStringArray(childIds)
	}
	if spec.ChildHealthThreshold != 0 {
		args.ChildHealthThreshold = pulumi.Int(int(spec.ChildHealthThreshold))
	}

	// --- CLOUDWATCH_METRIC: mirror an alarm (the private-resource pattern) ----
	if spec.CloudwatchAlarmName != "" {
		args.CloudwatchAlarmName = pulumi.String(spec.CloudwatchAlarmName)
	}
	if spec.CloudwatchAlarmRegion != "" {
		args.CloudwatchAlarmRegion = pulumi.String(spec.CloudwatchAlarmRegion)
	}
	if spec.InsufficientDataHealthStatus != "" {
		args.InsufficientDataHealthStatus = pulumi.String(spec.InsufficientDataHealthStatus)
	}

	// --- RECOVERY_CONTROL: mirror an ARC routing control (ForceNew) -----------
	if spec.RoutingControlArn != "" {
		args.RoutingControlArn = pulumi.String(spec.RoutingControlArn)
	}

	createdHealthCheck, err := route53.NewHealthCheck(ctx,
		locals.AwsRoute53HealthCheck.Metadata.Name,
		args,
		pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create health check")
	}

	ctx.Export(OpHealthCheckId, createdHealthCheck.ID())
	ctx.Export(OpHealthCheckArn, createdHealthCheck.Arn)

	return nil
}
