package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// autoscalePool provisions the droplet autoscale pool and exports its
// outputs. The spec's scaling oneof (static XOR dynamic) is rendered into
// the SDK's single Config object, which the provider itself never
// validates for shape consistency. Create waits for the pool AND every
// member droplet to reach "active" (up to 15 minutes upstream). DESTROY
// DESTROYS THE MEMBERS: the API's only delete terminates every droplet
// the pool owns.
func autoscalePool(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.DropletAutoscale, error) {
	spec := locals.DigitalOceanDropletAutoscalePool.Spec

	// Exactly one scaling branch is set (spec oneof); only that branch's
	// leaves are rendered. Unset leaves stay nil and never reach the API,
	// matching the wire's zero-means-unset behavior.
	config := &digitalocean.DropletAutoscaleConfigArgs{}
	if static := spec.GetStatic(); static != nil {
		config.TargetNumberInstances = pulumi.IntPtr(int(static.TargetInstances))
	}
	if dynamic := spec.GetDynamic(); dynamic != nil {
		config.MinInstances = pulumi.IntPtr(int(dynamic.MinInstances))
		config.MaxInstances = pulumi.IntPtr(int(dynamic.MaxInstances))
		if dynamic.TargetCpuUtilization != nil {
			config.TargetCpuUtilization = pulumi.Float64Ptr(*dynamic.TargetCpuUtilization)
		}
		if dynamic.TargetMemoryUtilization != nil {
			config.TargetMemoryUtilization = pulumi.Float64Ptr(*dynamic.TargetMemoryUtilization)
		}
		if dynamic.CooldownMinutes != nil {
			config.CooldownMinutes = pulumi.IntPtr(int(*dynamic.CooldownMinutes))
		}
	}

	template := spec.DropletTemplate

	// SSH key references resolve to literal numeric key ids before the
	// module runs.
	var sshKeys pulumi.StringArray
	for _, key := range template.SshKeys {
		sshKeys = append(sshKeys, pulumi.String(key.GetValue()))
	}

	// spec.droplet_template.tags plus the standard Planton labels rendered
	// as DigitalOcean "key:value" tags -- the exact set the Terraform
	// module applies. Tags land on every member droplet.
	tagSet := map[string]bool{}
	var tagInputs pulumi.StringArray
	for _, t := range template.Tags {
		if !tagSet[t] {
			tagSet[t] = true
			tagInputs = append(tagInputs, pulumi.String(t))
		}
	}
	for k, v := range locals.DigitalOceanLabels {
		t := k + ":" + v
		if !tagSet[t] {
			tagSet[t] = true
			tagInputs = append(tagInputs, pulumi.String(t))
		}
	}

	templateArgs := &digitalocean.DropletAutoscaleDropletTemplateArgs{
		Size: pulumi.String(template.Size),
		// The spec's region enum value names ARE the provider's region
		// slugs.
		Region: pulumi.String(template.Region.String()),
		// Slug or numeric image id; DigitalOcean reads the image back as
		// a numeric id, and the provider itself persists the configured
		// value to avoid the drift.
		Image:            pulumi.String(template.Image),
		SshKeys:          sshKeys,
		Tags:             tagInputs,
		WithDropletAgent: pulumi.BoolPtr(template.WithDropletAgent),
		Ipv6:             pulumi.BoolPtr(template.Ipv6),
	}
	if template.Vpc.GetValue() != "" {
		templateArgs.VpcUuid = pulumi.String(template.Vpc.GetValue())
	}
	if template.ProjectId.GetValue() != "" {
		templateArgs.ProjectId = pulumi.String(template.ProjectId.GetValue())
	}
	if template.UserData != "" {
		templateArgs.UserData = pulumi.String(template.UserData)
	}
	// public_networking is deliberately never rendered: the provider
	// declares it but never copies it into any create/update request --
	// dead on write at the pinned version (and absent from the SDK).

	createdPool, err := digitalocean.NewDropletAutoscale(
		ctx,
		"autoscale-pool",
		&digitalocean.DropletAutoscaleArgs{
			Name:            pulumi.String(spec.PoolName),
			Config:          config,
			DropletTemplate: templateArgs,
		},
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean droplet autoscale pool")
	}

	ctx.Export(OpPoolId, createdPool.ID())
	ctx.Export(OpStatus, createdPool.Status)

	return createdPool, nil
}
