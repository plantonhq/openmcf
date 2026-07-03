package module

import (
	"github.com/plantonhq/planton/internal/valuefrom"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// alb provisions the Application Load Balancer. The ALB carries no routing
// configuration by design: listeners, rules, and target groups are separate
// resources that attach to it by ARN, so this module owns only what is truly
// load-balancer-wide -- placement, security groups, and the HTTP behavior
// attributes. Changing "internal" replaces the load balancer; the attributes
// update in place.
func alb(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*lb.LoadBalancer, error) {
	spec := locals.AwsAlb.Spec

	// AWS limits load balancer names to 32 characters; truncate
	// deterministically so the same manifest always yields the same name.
	albName := truncateName(locals.AwsAlb.Metadata.Name, 32)

	args := &lb.LoadBalancerArgs{
		Name:                     pulumi.String(albName),
		LoadBalancerType:         pulumi.String("application"),
		Subnets:                  pulumi.ToStringArray(valuefrom.ToStringArray(spec.Subnets)),
		SecurityGroups:           pulumi.ToStringArray(valuefrom.ToStringArray(spec.SecurityGroups)),
		Internal:                 pulumi.Bool(spec.Internal),
		EnableDeletionProtection: pulumi.Bool(spec.DeleteProtectionEnabled),
		Tags:                     pulumi.ToStringMap(locals.AwsTags),
	}

	// Only explicitly set attributes are sent, so AWS keeps its own defaults
	// for the rest -- the module never bakes in opinions the spec does not
	// express.
	if spec.IpAddressType != "" {
		args.IpAddressType = pulumi.StringPtr(spec.IpAddressType)
	}
	if spec.IdleTimeoutSeconds > 0 {
		args.IdleTimeout = pulumi.IntPtr(int(spec.IdleTimeoutSeconds))
	}
	if spec.ClientKeepAliveSeconds > 0 {
		args.ClientKeepAlive = pulumi.IntPtr(int(spec.ClientKeepAliveSeconds))
	}

	// http2_enabled is optional bool: nil keeps the AWS default (true); an
	// explicit false downgrades clients to HTTP/1.1.
	if spec.Http2Enabled != nil {
		args.EnableHttp2 = pulumi.BoolPtr(spec.GetHttp2Enabled())
	}
	if spec.WafFailOpenEnabled {
		args.EnableWafFailOpen = pulumi.BoolPtr(true)
	}
	if spec.ZonalShiftEnabled {
		args.EnableZonalShift = pulumi.BoolPtr(true)
	}
	if spec.DropInvalidHeaderFields {
		args.DropInvalidHeaderFields = pulumi.BoolPtr(true)
	}
	if spec.PreserveHostHeader {
		args.PreserveHostHeader = pulumi.BoolPtr(true)
	}
	if spec.XffClientPortEnabled {
		args.EnableXffClientPort = pulumi.BoolPtr(true)
	}
	if spec.XffHeaderProcessingMode != "" {
		args.XffHeaderProcessingMode = pulumi.StringPtr(spec.XffHeaderProcessingMode)
	}
	if spec.DesyncMitigationMode != "" {
		args.DesyncMitigationMode = pulumi.StringPtr(spec.DesyncMitigationMode)
	}
	if spec.TlsVersionAndCipherSuiteHeadersEnabled {
		args.EnableTlsVersionAndCipherSuiteHeaders = pulumi.BoolPtr(true)
	}

	// The three S3 log streams share one shape; "enabled" is implied by the
	// block's presence in the spec (a bucket with logging off is meaningless).
	if spec.AccessLogs != nil {
		accessLogs := &lb.LoadBalancerAccessLogsArgs{
			Bucket:  pulumi.String(spec.AccessLogs.Bucket.GetValue()),
			Enabled: pulumi.BoolPtr(true),
		}
		if spec.AccessLogs.Prefix != "" {
			accessLogs.Prefix = pulumi.StringPtr(spec.AccessLogs.Prefix)
		}
		args.AccessLogs = accessLogs
	}
	if spec.ConnectionLogs != nil {
		connectionLogs := &lb.LoadBalancerConnectionLogsArgs{
			Bucket:  pulumi.String(spec.ConnectionLogs.Bucket.GetValue()),
			Enabled: pulumi.BoolPtr(true),
		}
		if spec.ConnectionLogs.Prefix != "" {
			connectionLogs.Prefix = pulumi.StringPtr(spec.ConnectionLogs.Prefix)
		}
		args.ConnectionLogs = connectionLogs
	}
	if spec.HealthCheckLogs != nil {
		healthCheckLogs := &lb.LoadBalancerHealthCheckLogsArgs{
			Bucket:  pulumi.String(spec.HealthCheckLogs.Bucket.GetValue()),
			Enabled: pulumi.BoolPtr(true),
		}
		if spec.HealthCheckLogs.Prefix != "" {
			healthCheckLogs.Prefix = pulumi.StringPtr(spec.HealthCheckLogs.Prefix)
		}
		args.HealthCheckLogs = healthCheckLogs
	}

	createdLoadBalancer, err := lb.NewLoadBalancer(ctx, albName, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create AWS ALB")
	}

	ctx.Export(OpAlbArn, createdLoadBalancer.Arn)
	ctx.Export(OpAlbName, createdLoadBalancer.Name)
	ctx.Export(OpAlbDnsName, createdLoadBalancer.DnsName)
	ctx.Export(OpAlbHostedZoneId, createdLoadBalancer.ZoneId)
	// The CloudWatch LoadBalancer dimension -- request-count autoscaling
	// policies compose it with a target group's arn_suffix.
	ctx.Export(OpAlbArnSuffix, createdLoadBalancer.ArnSuffix)

	return createdLoadBalancer, nil
}

// truncateName enforces AWS's 32-character load balancer name limit.
func truncateName(name string, maxLen int) string {
	if len(name) <= maxLen {
		return name
	}
	return name[:maxLen]
}
