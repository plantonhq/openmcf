package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// healthcheck creates the standalone health check. The unused config block is
// never sent (the spec's CEL guarantees http_config only on HTTP/HTTPS and
// tcp_config only on TCP), because both blocks are Computed upstream and
// sending the wrong one reads back as drift.
func healthcheck(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareHealthcheck.Spec

	args := &cloudflare.HealthcheckArgs{
		ZoneId:  pulumi.String(spec.ZoneId.GetValue()),
		Name:    pulumi.String(spec.Name),
		Address: pulumi.String(spec.Address),
	}

	if spec.Type != nil {
		args.Type = pulumi.StringPtr(spec.GetType())
	}
	if len(spec.CheckRegions) > 0 {
		args.CheckRegions = pulumi.ToStringArray(spec.CheckRegions)
	}
	if spec.ConsecutiveFails != nil {
		args.ConsecutiveFails = pulumi.IntPtr(int(spec.GetConsecutiveFails()))
	}
	if spec.ConsecutiveSuccesses != nil {
		args.ConsecutiveSuccesses = pulumi.IntPtr(int(spec.GetConsecutiveSuccesses()))
	}
	if spec.Interval != nil {
		args.Interval = pulumi.IntPtr(int(spec.GetInterval()))
	}
	if spec.Retries != nil {
		args.Retries = pulumi.IntPtr(int(spec.GetRetries()))
	}
	if spec.Timeout != nil {
		args.Timeout = pulumi.IntPtr(int(spec.GetTimeout()))
	}
	if spec.Suspended != nil {
		args.Suspended = pulumi.BoolPtr(spec.GetSuspended())
	}

	if spec.HttpConfig != nil {
		httpArgs := &cloudflare.HealthcheckHttpConfigArgs{}
		if spec.HttpConfig.Method != nil {
			httpArgs.Method = pulumi.StringPtr(spec.HttpConfig.GetMethod())
		}
		if spec.HttpConfig.Path != nil {
			httpArgs.Path = pulumi.StringPtr(spec.HttpConfig.GetPath())
		}
		if spec.HttpConfig.Port != nil {
			httpArgs.Port = pulumi.IntPtr(int(spec.HttpConfig.GetPort()))
		}
		if len(spec.HttpConfig.ExpectedCodes) > 0 {
			httpArgs.ExpectedCodes = pulumi.ToStringArray(spec.HttpConfig.ExpectedCodes)
		}
		if spec.HttpConfig.ExpectedBody != "" {
			httpArgs.ExpectedBody = pulumi.StringPtr(spec.HttpConfig.ExpectedBody)
		}
		if spec.HttpConfig.FollowRedirects != nil {
			httpArgs.FollowRedirects = pulumi.BoolPtr(spec.HttpConfig.GetFollowRedirects())
		}
		if spec.HttpConfig.AllowInsecure != nil {
			httpArgs.AllowInsecure = pulumi.BoolPtr(spec.HttpConfig.GetAllowInsecure())
		}
		if len(spec.HttpConfig.Headers) > 0 {
			// The spec wraps each header's values in a message (one header, many
			// values); the provider wants a plain map of string lists -- unwrap.
			headers := pulumi.StringArrayMap{}
			for name, wrapper := range spec.HttpConfig.Headers {
				headers[name] = pulumi.ToStringArray(wrapper.Values)
			}
			httpArgs.Header = headers
		}
		args.HttpConfig = httpArgs
	}

	if spec.TcpConfig != nil {
		tcpArgs := &cloudflare.HealthcheckTcpConfigArgs{}
		if spec.TcpConfig.Method != nil {
			tcpArgs.Method = pulumi.StringPtr(spec.TcpConfig.GetMethod())
		}
		if spec.TcpConfig.Port != nil {
			tcpArgs.Port = pulumi.IntPtr(int(spec.TcpConfig.GetPort()))
		}
		args.TcpConfig = tcpArgs
	}

	createdHealthcheck, err := cloudflare.NewHealthcheck(
		ctx,
		"healthcheck",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create healthcheck")
	}

	ctx.Export(OpHealthcheckId, createdHealthcheck.ID())
	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))

	return nil
}
