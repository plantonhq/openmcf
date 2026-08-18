package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/servicediscovery"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// namespace creates the namespace (whichever type arm), its services
// and statically registered instances, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the three namespace types are three provider resources; exactly
//     one exists per the spec's type field, all exposing the same
//     downstream surface (id/arn);
//   - the HTTP namespace has NO update path at the provider - changing
//     its description replaces it; the DNS namespaces update
//     description in place;
//   - the private namespace's Vpc is never read back by the provider
//     (imports carry it as "{namespace_id}:{vpc_id}");
//   - a service binds its namespace through DnsConfig.NamespaceId when
//     it publishes DNS records, and through the top-level NamespaceId
//     otherwise (the provider's own documented split for its legacy
//     duplicated pointer);
//   - instance registration is an AWS upsert (create and update are
//     the same RegisterInstance call, keyed by instance_id); the
//     provider derives AWS_INSTANCE_IPV4 from AWS_EC2_INSTANCE_ID and
//     drops the derived echo on read;
//   - deregistering an already-gone instance errors at the provider
//     (no NotFound tolerance) - destroy instances through this module,
//     never out-of-band;
//   - a service's ForceDestroy deregisters EVERY instance in the
//     service first, including runtime-registered ones this manifest
//     never declared.
func namespace(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	var namespaceId pulumi.StringOutput
	var namespaceArn pulumi.StringOutput
	var hostedZoneId pulumi.StringOutput
	var httpName pulumi.StringOutput

	switch spec.Type {
	case "HTTP":
		args := &servicediscovery.HttpNamespaceArgs{
			Name: pulumi.String(locals.Target.Metadata.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if spec.Description != "" {
			args.Description = pulumi.String(spec.Description)
		}
		createdNamespace, err := servicediscovery.NewHttpNamespace(ctx, "namespace", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create http namespace")
		}
		namespaceId = createdNamespace.ID().ToStringOutput()
		namespaceArn = createdNamespace.Arn
		hostedZoneId = pulumi.String("").ToStringOutput()
		httpName = createdNamespace.HttpName

	case "PRIVATE_DNS":
		args := &servicediscovery.PrivateDnsNamespaceArgs{
			Name: pulumi.String(locals.Target.Metadata.Name),
			Vpc:  pulumi.String(spec.VpcId.GetValue()),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if spec.Description != "" {
			args.Description = pulumi.String(spec.Description)
		}
		createdNamespace, err := servicediscovery.NewPrivateDnsNamespace(ctx, "namespace", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create private dns namespace")
		}
		namespaceId = createdNamespace.ID().ToStringOutput()
		namespaceArn = createdNamespace.Arn
		hostedZoneId = createdNamespace.HostedZone
		httpName = pulumi.String("").ToStringOutput()

	default: // PUBLIC_DNS per the spec's vocabulary
		args := &servicediscovery.PublicDnsNamespaceArgs{
			Name: pulumi.String(locals.Target.Metadata.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if spec.Description != "" {
			args.Description = pulumi.String(spec.Description)
		}
		createdNamespace, err := servicediscovery.NewPublicDnsNamespace(ctx, "namespace", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create public dns namespace")
		}
		namespaceId = createdNamespace.ID().ToStringOutput()
		namespaceArn = createdNamespace.Arn
		hostedZoneId = createdNamespace.HostedZone
		httpName = pulumi.String("").ToStringOutput()
	}

	// Services, keyed by service name. Each registration's owning
	// service ID is also exported keyed "{service}//{instance_id}" -
	// the first half of the instance's composite import ID.
	serviceIds := pulumi.StringMap{}
	serviceArns := pulumi.StringMap{}
	instanceServiceIds := pulumi.StringMap{}
	for _, service := range spec.Services {
		args := &servicediscovery.ServiceArgs{
			Name: pulumi.String(service.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if service.Description != "" {
			args.Description = pulumi.String(service.Description)
		}

		if service.DnsConfig != nil {
			records := servicediscovery.ServiceDnsConfigDnsRecordArray{}
			for _, record := range service.DnsConfig.Records {
				records = append(records, &servicediscovery.ServiceDnsConfigDnsRecordArgs{
					Type: pulumi.String(record.Type),
					Ttl:  pulumi.Int(int(record.Ttl)),
				})
			}
			dnsConfig := &servicediscovery.ServiceDnsConfigArgs{
				NamespaceId: namespaceId,
				DnsRecords:  records,
			}
			if service.DnsConfig.RoutingPolicy != "" {
				dnsConfig.RoutingPolicy = pulumi.String(service.DnsConfig.RoutingPolicy)
			}
			args.DnsConfig = dnsConfig
		} else {
			// API-only services bind the namespace at the top level.
			args.NamespaceId = namespaceId
		}

		if service.HealthCheckConfig != nil {
			healthCheck := &servicediscovery.ServiceHealthCheckConfigArgs{}
			if service.HealthCheckConfig.Type != "" {
				healthCheck.Type = pulumi.String(service.HealthCheckConfig.Type)
			}
			if service.HealthCheckConfig.ResourcePath != "" {
				healthCheck.ResourcePath = pulumi.String(service.HealthCheckConfig.ResourcePath)
			}
			if service.HealthCheckConfig.FailureThreshold > 0 {
				healthCheck.FailureThreshold = pulumi.Int(int(service.HealthCheckConfig.FailureThreshold))
			}
			args.HealthCheckConfig = healthCheck
		}
		if service.HealthCheckCustomConfig != nil {
			args.HealthCheckCustomConfig = &servicediscovery.ServiceHealthCheckCustomConfigArgs{}
		}
		if service.ForceDestroy {
			args.ForceDestroy = pulumi.Bool(true)
		}

		createdService, err := servicediscovery.NewService(ctx,
			fmt.Sprintf("service-%s", service.Name), args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create service %s", service.Name)
		}
		serviceIds[service.Name] = createdService.ID().ToStringOutput()
		serviceArns[service.Name] = createdService.Arn

		// Statically registered instances - the typed fields compose
		// AWS's documented attribute keys; anything else rides
		// custom_attributes verbatim.
		for _, instance := range service.Instances {
			attributes := map[string]string{}
			for key, value := range instance.CustomAttributes {
				attributes[key] = value
			}
			if instance.Ip != "" {
				attributes["AWS_INSTANCE_IPV4"] = instance.Ip
			}
			if instance.Ipv6 != "" {
				attributes["AWS_INSTANCE_IPV6"] = instance.Ipv6
			}
			if instance.Port > 0 {
				attributes["AWS_INSTANCE_PORT"] = fmt.Sprintf("%d", instance.Port)
			}
			if instance.Cname != "" {
				attributes["AWS_INSTANCE_CNAME"] = instance.Cname
			}
			if instance.AliasDnsName.GetValue() != "" {
				attributes["AWS_ALIAS_DNS_NAME"] = instance.AliasDnsName.GetValue()
			}
			if instance.Ec2InstanceId.GetValue() != "" {
				attributes["AWS_EC2_INSTANCE_ID"] = instance.Ec2InstanceId.GetValue()
			}

			if _, err := servicediscovery.NewInstance(ctx,
				fmt.Sprintf("instance-%s-%s", service.Name, instance.InstanceId),
				&servicediscovery.InstanceArgs{
					InstanceId: pulumi.String(instance.InstanceId),
					ServiceId:  createdService.ID(),
					Attributes: pulumi.ToStringMap(attributes),
				}, pulumi.Provider(provider)); err != nil {
				return errors.Wrapf(err, "register instance %s/%s", service.Name, instance.InstanceId)
			}
			instanceServiceIds[fmt.Sprintf("%s//%s", service.Name, instance.InstanceId)] = createdService.ID().ToStringOutput()
		}
	}

	ctx.Export(OpNamespaceId, namespaceId)
	ctx.Export(OpNamespaceArn, namespaceArn)
	ctx.Export(OpHostedZoneId, hostedZoneId)
	ctx.Export(OpHttpName, httpName)
	ctx.Export(OpServiceIds, serviceIds)
	ctx.Export(OpServiceArns, serviceArns)
	ctx.Export(OpInstanceServiceIds, instanceServiceIds)
	return nil
}
