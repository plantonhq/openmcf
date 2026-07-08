package module

import (
	"fmt"

	awskinesisfirehose "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awskinesisfirehose/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/kinesis"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The spec models the transformation pipeline as an ordered list of TYPED
// processor blocks (lambda, metadata_extraction, decompression,
// cloudwatch_log_processing, append_delimiter, record_deaggregation) so a
// manifest reads as intent instead of AWS's internal parameter vocabulary.
// The provider takes the raw {type, parameters[]} shape, so the typed arms
// are normalized here ONCE; each destination file adapts the normalized list
// to its destination-specific SDK types (which are structurally identical
// but nominally distinct).

// processorParameter is one AWS ProcessorParameterName/value pair.
type processorParameter struct {
	Name  string
	Value string
}

// normalizedProcessor is one provider-shaped processor: the AWS ProcessorType
// plus its parameter pairs, in the spec's declared order.
type normalizedProcessor struct {
	Type       string
	Parameters []processorParameter
}

// normalizeProcessors translates the typed processor arms into the provider's
// {type, parameters[]} shape. Returns nil when processing is absent or
// disabled. Optional numeric knobs are only sent when set (> 0) so AWS
// applies its own defaults otherwise.
func normalizeProcessors(processing *awskinesisfirehose.AwsKinesisFirehoseProcessing) []normalizedProcessor {
	if processing == nil || !processing.Enabled {
		return nil
	}

	normalized := make([]normalizedProcessor, 0, len(processing.Processors))
	for _, p := range processing.Processors {
		switch {
		case p.Lambda != nil:
			params := []processorParameter{
				{Name: "LambdaArn", Value: p.Lambda.LambdaArn.GetValue()},
			}
			if p.Lambda.BufferSizeInMbs > 0 {
				params = append(params, processorParameter{Name: "BufferSizeInMBs", Value: fmt.Sprintf("%d", p.Lambda.BufferSizeInMbs)})
			}
			if p.Lambda.BufferIntervalInSeconds > 0 {
				params = append(params, processorParameter{Name: "BufferIntervalInSeconds", Value: fmt.Sprintf("%d", p.Lambda.BufferIntervalInSeconds)})
			}
			if p.Lambda.NumberOfRetries > 0 {
				params = append(params, processorParameter{Name: "NumberOfRetries", Value: fmt.Sprintf("%d", p.Lambda.NumberOfRetries)})
			}
			normalized = append(normalized, normalizedProcessor{Type: "Lambda", Parameters: params})

		case p.MetadataExtraction != nil:
			// JsonParsingEngine is mandatory for MetadataExtraction; JQ-1.6 is
			// the only engine AWS supports today.
			engine := p.MetadataExtraction.JsonParsingEngine
			if engine == "" {
				engine = "JQ-1.6"
			}
			normalized = append(normalized, normalizedProcessor{
				Type: "MetadataExtraction",
				Parameters: []processorParameter{
					{Name: "MetadataExtractionQuery", Value: p.MetadataExtraction.Query},
					{Name: "JsonParsingEngine", Value: engine},
				},
			})

		case p.Decompression != nil:
			normalized = append(normalized, normalizedProcessor{
				Type: "Decompression",
				Parameters: []processorParameter{
					{Name: "CompressionFormat", Value: p.Decompression.CompressionFormat},
				},
			})

		case p.CloudwatchLogProcessing != nil:
			value := "false"
			if p.CloudwatchLogProcessing.DataMessageExtraction {
				value = "true"
			}
			normalized = append(normalized, normalizedProcessor{
				Type: "CloudWatchLogProcessing",
				Parameters: []processorParameter{
					{Name: "DataMessageExtraction", Value: value},
				},
			})

		case p.AppendDelimiter != nil:
			normalized = append(normalized, normalizedProcessor{
				Type: "AppendDelimiterToRecord",
				Parameters: []processorParameter{
					{Name: "Delimiter", Value: p.AppendDelimiter.Delimiter},
				},
			})

		case p.RecordDeaggregation != nil:
			params := []processorParameter{
				{Name: "SubRecordType", Value: p.RecordDeaggregation.SubRecordType},
			}
			if p.RecordDeaggregation.Delimiter != "" {
				params = append(params, processorParameter{Name: "Delimiter", Value: p.RecordDeaggregation.Delimiter})
			}
			normalized = append(normalized, normalizedProcessor{Type: "RecordDeAggregation", Parameters: params})
		}
	}

	return normalized
}

// buildExtendedS3Processing adapts the normalized pipeline to the Extended S3
// destination's SDK types.
func buildExtendedS3Processing(processing *awskinesisfirehose.AwsKinesisFirehoseProcessing) *kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationProcessingConfigurationArgs {
	normalized := normalizeProcessors(processing)
	if normalized == nil {
		return nil
	}
	processors := kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationProcessingConfigurationProcessorArray{}
	for _, p := range normalized {
		params := kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationProcessingConfigurationProcessorParameterArray{}
		for _, param := range p.Parameters {
			params = append(params, &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationProcessingConfigurationProcessorParameterArgs{
				ParameterName:  pulumi.String(param.Name),
				ParameterValue: pulumi.String(param.Value),
			})
		}
		processors = append(processors, &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationProcessingConfigurationProcessorArgs{
			Type:       pulumi.String(p.Type),
			Parameters: params,
		})
	}
	return &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationProcessingConfigurationArgs{
		Enabled:    pulumi.Bool(true),
		Processors: processors,
	}
}

// buildCloudwatchLogging constructs CloudWatch logging options for the
// Extended S3 destination. Returns nil if not configured or not enabled.
func buildCloudwatchLogging(logging *awskinesisfirehose.AwsKinesisFirehoseCloudwatchLogging) *kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationCloudwatchLoggingOptionsArgs {
	if logging == nil || !logging.Enabled {
		return nil
	}

	return &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationCloudwatchLoggingOptionsArgs{
		Enabled:       pulumi.Bool(true),
		LogGroupName:  pulumi.StringPtr(logging.LogGroupName),
		LogStreamName: pulumi.StringPtr(logging.LogStreamName),
	}
}
