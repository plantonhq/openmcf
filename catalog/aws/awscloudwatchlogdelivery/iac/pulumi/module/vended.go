package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// vended creates the vended-log pipeline: the delivery source, owned
// destinations (with their policies), and the deliveries joining them.
//
// Lifecycle facts the render below depends on:
//   - the source is per (resource, log_type): name, log_type, and
//     resource_arn all replace on change;
//   - a delivery's source and destination both replace on change; only
//     the wire-format settings update in place;
//   - AWS allows ONE delivery per (source, destination-type) pair;
//   - for CloudFront sources AWS prepends
//     "AWSLogs/{account-id}/CloudFront/" to the S3 suffix path; the
//     provider strips that prefix on reads, so configure only your own
//     path segment.
func vended(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec.Vended

	destinationArns := pulumi.StringMap{}
	deliveryIds := pulumi.StringMap{}
	deliveryArns := pulumi.StringMap{}

	if spec == nil {
		ctx.Export(OpSourceName, pulumi.String(""))
		ctx.Export(OpSourceArn, pulumi.String(""))
		ctx.Export(OpSourceService, pulumi.String(""))
		ctx.Export(OpDestinationArns, destinationArns)
		ctx.Export(OpDeliveryIds, deliveryIds)
		ctx.Export(OpDeliveryArns, deliveryArns)
		return nil
	}

	var createdSource *cloudwatch.LogDeliverySource
	if source := spec.Source; source != nil {
		var err error
		createdSource, err = cloudwatch.NewLogDeliverySource(ctx, "source", &cloudwatch.LogDeliverySourceArgs{
			Name:        pulumi.String(source.Name),
			LogType:     pulumi.String(source.LogType),
			ResourceArn: pulumi.String(source.ResourceArn.GetValue()),
			Tags:        pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create delivery source")
		}
		ctx.Export(OpSourceName, createdSource.Name)
		ctx.Export(OpSourceArn, createdSource.Arn)
		ctx.Export(OpSourceService, createdSource.Service)
	} else {
		ctx.Export(OpSourceName, pulumi.String(""))
		ctx.Export(OpSourceArn, pulumi.String(""))
		ctx.Export(OpSourceService, pulumi.String(""))
	}

	createdDestinations := map[string]*cloudwatch.LogDeliveryDestination{}
	for _, destination := range spec.Destinations {
		args := &cloudwatch.LogDeliveryDestinationArgs{
			Name: pulumi.String(destination.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}

		// XRAY destinations carry no configuration block
		// (spec-guaranteed); every other type requires one.
		if destination.DestinationResourceArn.GetValue() != "" {
			args.DeliveryDestinationConfiguration = &cloudwatch.LogDeliveryDestinationDeliveryDestinationConfigurationArgs{
				DestinationResourceArn: pulumi.String(destination.DestinationResourceArn.GetValue()),
			}
		}
		if destination.DeliveryDestinationType != "" {
			args.DeliveryDestinationType = pulumi.String(destination.DeliveryDestinationType)
		}
		if destination.OutputFormat != "" {
			args.OutputFormat = pulumi.String(destination.OutputFormat)
		}

		createdDestination, err := cloudwatch.NewLogDeliveryDestination(ctx, "destination-"+destination.Name, args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create delivery destination %s", destination.Name)
		}
		createdDestinations[destination.Name] = createdDestination
		destinationArns[destination.Name] = createdDestination.Arn

		if destination.Policy != nil {
			policyJson, err := json.Marshal(destination.Policy.AsMap())
			if err != nil {
				return errors.Wrapf(err, "marshal policy for destination %s", destination.Name)
			}
			if _, err := cloudwatch.NewLogDeliveryDestinationPolicy(ctx, "destination-policy-"+destination.Name, &cloudwatch.LogDeliveryDestinationPolicyArgs{
				DeliveryDestinationName:   createdDestination.Name,
				DeliveryDestinationPolicy: pulumi.String(string(policyJson)),
			}, pulumi.Provider(provider)); err != nil {
				return errors.Wrapf(err, "create policy for destination %s", destination.Name)
			}
		}
	}

	for _, delivery := range spec.Deliveries {
		var destinationArn pulumi.StringInput
		options := []pulumi.ResourceOption{pulumi.Provider(provider)}
		if delivery.DestinationName != "" {
			ownedDestination := createdDestinations[delivery.DestinationName]
			destinationArn = ownedDestination.Arn
			options = append(options, pulumi.DependsOn([]pulumi.Resource{ownedDestination}))
		} else {
			destinationArn = pulumi.String(delivery.DestinationArn.GetValue())
		}

		args := &cloudwatch.LogDeliveryArgs{
			DeliverySourceName:     createdSource.Name,
			DeliveryDestinationArn: destinationArn,
			Tags:                   pulumi.ToStringMap(locals.AwsTags),
		}

		if len(delivery.RecordFields) > 0 {
			args.RecordFields = pulumi.ToStringArray(delivery.RecordFields)
		}
		if delivery.FieldDelimiter != "" {
			args.FieldDelimiter = pulumi.String(delivery.FieldDelimiter)
		}
		if s3Configuration := delivery.S3Configuration; s3Configuration != nil {
			args.S3DeliveryConfigurations = cloudwatch.LogDeliveryS3DeliveryConfigurationArray{
				&cloudwatch.LogDeliveryS3DeliveryConfigurationArgs{
					EnableHiveCompatiblePath: pulumi.Bool(s3Configuration.EnableHiveCompatiblePath),
					SuffixPath:               pulumi.String(s3Configuration.SuffixPath),
				},
			}
		}

		createdDelivery, err := cloudwatch.NewLogDelivery(ctx, "delivery-"+delivery.Name, args, options...)
		if err != nil {
			return errors.Wrapf(err, "create delivery %s", delivery.Name)
		}
		deliveryIds[delivery.Name] = createdDelivery.ID().ToStringOutput()
		deliveryArns[delivery.Name] = createdDelivery.Arn
	}

	ctx.Export(OpDestinationArns, destinationArns)
	ctx.Export(OpDeliveryIds, deliveryIds)
	ctx.Export(OpDeliveryArns, deliveryArns)
	return nil
}
