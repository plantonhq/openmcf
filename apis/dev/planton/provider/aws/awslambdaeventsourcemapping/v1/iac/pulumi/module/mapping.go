package module

import (
	"strings"

	"github.com/pkg/errors"
	awslambdaeventsourcemappingv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awslambdaeventsourcemapping/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// mapping provisions the event source mapping. The event source (ARN or
// self-managed Kafka bootstrap servers) is create-time immutable; batching,
// filters, failure handling, and the target function edit in place.
func mapping(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*lambda.EventSourceMapping, error) {
	spec := locals.AwsLambdaEventSourceMapping.Spec

	args := &lambda.EventSourceMappingArgs{
		FunctionName: pulumi.String(spec.FunctionArn.GetValue()),
		Enabled:      pulumi.Bool(!spec.Disabled),
		Tags:         pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.SelfManagedKafka != nil {
		args.SelfManagedEventSource = &lambda.EventSourceMappingSelfManagedEventSourceArgs{
			Endpoints: pulumi.StringMap{
				"KAFKA_BOOTSTRAP_SERVERS": pulumi.String(strings.Join(spec.SelfManagedKafka.BootstrapServers, ",")),
			},
		}
	} else if spec.EventSourceArn.GetValue() != "" {
		args.EventSourceArn = pulumi.String(spec.EventSourceArn.GetValue())
	}

	if spec.BatchSize != 0 {
		args.BatchSize = pulumi.Int(int(spec.BatchSize))
	}
	if spec.MaximumBatchingWindowSeconds != 0 {
		args.MaximumBatchingWindowInSeconds = pulumi.Int(int(spec.MaximumBatchingWindowSeconds))
	}

	if len(spec.Filters) > 0 {
		filters := lambda.EventSourceMappingFilterCriteriaFilterArray{}
		for _, f := range spec.Filters {
			filterArgs := &lambda.EventSourceMappingFilterCriteriaFilterArgs{}
			if f.Pattern != "" {
				filterArgs.Pattern = pulumi.String(f.Pattern)
			}
			filters = append(filters, filterArgs)
		}
		args.FilterCriteria = &lambda.EventSourceMappingFilterCriteriaArgs{
			Filters: filters,
		}
	}

	if spec.KmsKeyArn.GetValue() != "" {
		args.KmsKeyArn = pulumi.String(spec.KmsKeyArn.GetValue())
	}

	if len(spec.FunctionResponseTypes) > 0 {
		args.FunctionResponseTypes = pulumi.ToStringArray(spec.FunctionResponseTypes)
	}

	if spec.ScalingMaxConcurrency != 0 {
		args.ScalingConfig = &lambda.EventSourceMappingScalingConfigArgs{
			MaximumConcurrency: pulumi.Int(int(spec.ScalingMaxConcurrency)),
		}
	}

	if len(spec.Metrics) > 0 {
		args.MetricsConfig = &lambda.EventSourceMappingMetricsConfigArgs{
			Metrics: pulumi.ToStringArray(spec.Metrics),
		}
	}

	if spec.StartingPosition != "" {
		args.StartingPosition = pulumi.String(spec.StartingPosition)
	}
	if spec.StartingPositionTimestamp != "" {
		args.StartingPositionTimestamp = pulumi.String(spec.StartingPositionTimestamp)
	}
	if spec.ParallelizationFactor != 0 {
		args.ParallelizationFactor = pulumi.Int(int(spec.ParallelizationFactor))
	}
	if spec.MaximumRecordAgeSeconds != 0 {
		args.MaximumRecordAgeInSeconds = pulumi.Int(int(spec.MaximumRecordAgeSeconds))
	}
	if spec.MaximumRetryAttempts != nil {
		args.MaximumRetryAttempts = pulumi.Int(int(*spec.MaximumRetryAttempts))
	}
	if spec.BisectBatchOnFunctionError {
		args.BisectBatchOnFunctionError = pulumi.Bool(true)
	}
	if spec.TumblingWindowSeconds != 0 {
		args.TumblingWindowInSeconds = pulumi.Int(int(spec.TumblingWindowSeconds))
	}

	if spec.OnFailureDestinationArn.GetValue() != "" {
		args.DestinationConfig = &lambda.EventSourceMappingDestinationConfigArgs{
			OnFailure: &lambda.EventSourceMappingDestinationConfigOnFailureArgs{
				DestinationArn: pulumi.String(spec.OnFailureDestinationArn.GetValue()),
			},
		}
	}

	if len(spec.Topics) > 0 {
		args.Topics = pulumi.ToStringArray(spec.Topics)
	}

	if spec.MqQueue != "" {
		args.Queues = pulumi.String(spec.MqQueue)
	}

	if spec.DocumentDb != nil {
		docArgs := &lambda.EventSourceMappingDocumentDbEventSourceConfigArgs{
			DatabaseName: pulumi.String(spec.DocumentDb.DatabaseName),
		}
		if spec.DocumentDb.CollectionName != "" {
			docArgs.CollectionName = pulumi.String(spec.DocumentDb.CollectionName)
		}
		if spec.DocumentDb.FullDocument != "" {
			docArgs.FullDocument = pulumi.String(spec.DocumentDb.FullDocument)
		}
		args.DocumentDbEventSourceConfig = docArgs
	}

	if len(spec.SourceAccessConfigurations) > 0 {
		args.SourceAccessConfigurations = sourceAccessConfigurations(spec.SourceAccessConfigurations)
	}

	if spec.ProvisionedPollers != nil {
		pollerArgs := &lambda.EventSourceMappingProvisionedPollerConfigArgs{}
		if spec.ProvisionedPollers.MinimumPollers != 0 {
			pollerArgs.MinimumPollers = pulumi.Int(int(spec.ProvisionedPollers.MinimumPollers))
		}
		if spec.ProvisionedPollers.MaximumPollers != 0 {
			pollerArgs.MaximumPollers = pulumi.Int(int(spec.ProvisionedPollers.MaximumPollers))
		}
		args.ProvisionedPollerConfig = pollerArgs
	}

	// Kafka source config: consumer group + schema registry. AWS models the
	// same shape once per source family (MSK vs self-managed), so both arms
	// mirror the Terraform module's paired dynamic blocks.
	if spec.KafkaConsumerGroupId != "" || spec.SchemaRegistry != nil {
		if spec.SelfManagedKafka != nil {
			cfg := &lambda.EventSourceMappingSelfManagedKafkaEventSourceConfigArgs{}
			if spec.KafkaConsumerGroupId != "" {
				cfg.ConsumerGroupId = pulumi.String(spec.KafkaConsumerGroupId)
			}
			if spec.SchemaRegistry != nil {
				cfg.SchemaRegistryConfig = &lambda.EventSourceMappingSelfManagedKafkaEventSourceConfigSchemaRegistryConfigArgs{
					SchemaRegistryUri:       pulumi.String(spec.SchemaRegistry.Uri),
					EventRecordFormat:       pulumi.String(spec.SchemaRegistry.EventRecordFormat),
					SchemaValidationConfigs: selfManagedSchemaValidationConfigs(spec.SchemaRegistry.ValidationAttributes),
					AccessConfigs:           selfManagedSchemaAccessConfigs(spec.SchemaRegistry.AccessConfigurations),
				}
			}
			args.SelfManagedKafkaEventSourceConfig = cfg
		} else if len(spec.Topics) > 0 {
			cfg := &lambda.EventSourceMappingAmazonManagedKafkaEventSourceConfigArgs{}
			if spec.KafkaConsumerGroupId != "" {
				cfg.ConsumerGroupId = pulumi.String(spec.KafkaConsumerGroupId)
			}
			if spec.SchemaRegistry != nil {
				cfg.SchemaRegistryConfig = &lambda.EventSourceMappingAmazonManagedKafkaEventSourceConfigSchemaRegistryConfigArgs{
					SchemaRegistryUri:       pulumi.String(spec.SchemaRegistry.Uri),
					EventRecordFormat:       pulumi.String(spec.SchemaRegistry.EventRecordFormat),
					SchemaValidationConfigs: amazonManagedSchemaValidationConfigs(spec.SchemaRegistry.ValidationAttributes),
					AccessConfigs:           amazonManagedSchemaAccessConfigs(spec.SchemaRegistry.AccessConfigurations),
				}
			}
			args.AmazonManagedKafkaEventSourceConfig = cfg
		}
	}

	createdMapping, err := lambda.NewEventSourceMapping(ctx, locals.MappingName, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create event source mapping")
	}

	return createdMapping, nil
}

func sourceAccessConfigurations(configs []*awslambdaeventsourcemappingv1.AwsLambdaEventSourceMappingSourceAccess) lambda.EventSourceMappingSourceAccessConfigurationArray {
	result := lambda.EventSourceMappingSourceAccessConfigurationArray{}
	for _, c := range configs {
		result = append(result, &lambda.EventSourceMappingSourceAccessConfigurationArgs{
			Type: pulumi.String(c.Type),
			Uri:  pulumi.String(c.Uri),
		})
	}
	return result
}

// The schema-registry validation/access builders below exist once per Kafka
// source family because the provider generates a distinct (but identically
// shaped) type tree for each -- MSK and self-managed configs cannot share
// argument structs.

func amazonManagedSchemaValidationConfigs(attributes []string) lambda.EventSourceMappingAmazonManagedKafkaEventSourceConfigSchemaRegistryConfigSchemaValidationConfigArray {
	result := lambda.EventSourceMappingAmazonManagedKafkaEventSourceConfigSchemaRegistryConfigSchemaValidationConfigArray{}
	for _, attribute := range attributes {
		result = append(result, &lambda.EventSourceMappingAmazonManagedKafkaEventSourceConfigSchemaRegistryConfigSchemaValidationConfigArgs{
			Attribute: pulumi.String(attribute),
		})
	}
	return result
}

func amazonManagedSchemaAccessConfigs(configs []*awslambdaeventsourcemappingv1.AwsLambdaEventSourceMappingSourceAccess) lambda.EventSourceMappingAmazonManagedKafkaEventSourceConfigSchemaRegistryConfigAccessConfigArray {
	result := lambda.EventSourceMappingAmazonManagedKafkaEventSourceConfigSchemaRegistryConfigAccessConfigArray{}
	for _, c := range configs {
		result = append(result, &lambda.EventSourceMappingAmazonManagedKafkaEventSourceConfigSchemaRegistryConfigAccessConfigArgs{
			Type: pulumi.String(c.Type),
			Uri:  pulumi.String(c.Uri),
		})
	}
	return result
}

func selfManagedSchemaValidationConfigs(attributes []string) lambda.EventSourceMappingSelfManagedKafkaEventSourceConfigSchemaRegistryConfigSchemaValidationConfigArray {
	result := lambda.EventSourceMappingSelfManagedKafkaEventSourceConfigSchemaRegistryConfigSchemaValidationConfigArray{}
	for _, attribute := range attributes {
		result = append(result, &lambda.EventSourceMappingSelfManagedKafkaEventSourceConfigSchemaRegistryConfigSchemaValidationConfigArgs{
			Attribute: pulumi.String(attribute),
		})
	}
	return result
}

func selfManagedSchemaAccessConfigs(configs []*awslambdaeventsourcemappingv1.AwsLambdaEventSourceMappingSourceAccess) lambda.EventSourceMappingSelfManagedKafkaEventSourceConfigSchemaRegistryConfigAccessConfigArray {
	result := lambda.EventSourceMappingSelfManagedKafkaEventSourceConfigSchemaRegistryConfigAccessConfigArray{}
	for _, c := range configs {
		result = append(result, &lambda.EventSourceMappingSelfManagedKafkaEventSourceConfigSchemaRegistryConfigAccessConfigArgs{
			Type: pulumi.String(c.Type),
			Uri:  pulumi.String(c.Uri),
		})
	}
	return result
}
