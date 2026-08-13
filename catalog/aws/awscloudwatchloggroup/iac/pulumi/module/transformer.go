package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	awscloudwatchloggroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awscloudwatchloggroup/v1alpha1"
)

// transformer creates the group's aws_cloudwatch_log_transformer twin: one
// PutTransformer document carrying the ordered processor pipeline. Each spec
// processors entry carries exactly one processor (CEL-enforced); the builder
// maps that arm onto the SDK's per-type members. Optional strings are sent
// only when set (nested defaults are not materialized inside repeated
// messages); booleans inside entries are always sent — PutTransformer
// replaces the whole pipeline, so explicit false and omitted false are the
// same write.
func transformer(
	ctx *pulumi.Context,
	spec *awscloudwatchloggroupv1alpha1.AwsCloudwatchLogGroupTransformer,
	createdLogGroup *cloudwatch.LogGroup,
	provider *aws.Provider,
) error {
	configs := make(cloudwatch.LogTransformerTransformerConfigArray, 0, len(spec.Processors))
	for _, processor := range spec.Processors {
		configs = append(configs, buildProcessorConfig(processor))
	}

	if _, err := cloudwatch.NewLogTransformer(
		ctx,
		"transformer",
		&cloudwatch.LogTransformerArgs{
			LogGroupArn:        createdLogGroup.Arn,
			TransformerConfigs: configs,
		},
		pulumi.Provider(provider),
		pulumi.DependsOn([]pulumi.Resource{createdLogGroup}),
	); err != nil {
		return errors.Wrap(err, "failed to create transformer")
	}
	return nil
}

// buildProcessorConfig maps one spec processor entry (exactly one arm set)
// onto the SDK's transformer_config member. The bridge types single-use
// members as pointers and multi-use members as arrays; our one-processor-
// per-entry model always contributes a single element to array members.
func buildProcessorConfig(
	processor *awscloudwatchloggroupv1alpha1.AwsCloudwatchLogGroupTransformerProcessor,
) cloudwatch.LogTransformerTransformerConfigArgs {
	config := cloudwatch.LogTransformerTransformerConfigArgs{}

	if arm := processor.AddKeys; arm != nil {
		entries := make(cloudwatch.LogTransformerTransformerConfigAddKeysEntryArray, 0, len(arm.Entries))
		for _, entry := range arm.Entries {
			entries = append(entries, cloudwatch.LogTransformerTransformerConfigAddKeysEntryArgs{
				Key:               pulumi.String(entry.Key),
				Value:             pulumi.String(entry.Value),
				OverwriteIfExists: pulumi.BoolPtr(entry.OverwriteIfExists),
			})
		}
		config.AddKeys = &cloudwatch.LogTransformerTransformerConfigAddKeysArgs{Entries: entries}
	}

	if arm := processor.CopyValue; arm != nil {
		entries := make(cloudwatch.LogTransformerTransformerConfigCopyValueEntryArray, 0, len(arm.Entries))
		for _, entry := range arm.Entries {
			entries = append(entries, cloudwatch.LogTransformerTransformerConfigCopyValueEntryArgs{
				Source:            pulumi.String(entry.Source),
				Target:            pulumi.String(entry.Target),
				OverwriteIfExists: pulumi.BoolPtr(entry.OverwriteIfExists),
			})
		}
		config.CopyValue = &cloudwatch.LogTransformerTransformerConfigCopyValueArgs{Entries: entries}
	}

	if arm := processor.Csv; arm != nil {
		csvArgs := cloudwatch.LogTransformerTransformerConfigCsvArgs{}
		if len(arm.Columns) > 0 {
			csvArgs.Columns = pulumi.ToStringArray(arm.Columns)
		}
		if arm.Delimiter != "" {
			csvArgs.Delimiter = pulumi.StringPtr(arm.Delimiter)
		}
		if arm.QuoteCharacter != "" {
			csvArgs.QuoteCharacter = pulumi.StringPtr(arm.QuoteCharacter)
		}
		if arm.Source != "" {
			csvArgs.Source = pulumi.StringPtr(arm.Source)
		}
		config.Csvs = cloudwatch.LogTransformerTransformerConfigCsvArray{csvArgs}
	}

	if arm := processor.DateTimeConverter; arm != nil {
		converterArgs := cloudwatch.LogTransformerTransformerConfigDateTimeConverterArgs{
			Source:        pulumi.String(arm.Source),
			Target:        pulumi.String(arm.Target),
			MatchPatterns: pulumi.ToStringArray(arm.MatchPatterns),
		}
		if arm.Locale != "" {
			converterArgs.Locale = pulumi.StringPtr(arm.Locale)
		}
		if arm.SourceTimezone != "" {
			converterArgs.SourceTimezone = pulumi.StringPtr(arm.SourceTimezone)
		}
		if arm.TargetFormat != "" {
			converterArgs.TargetFormat = pulumi.StringPtr(arm.TargetFormat)
		}
		if arm.TargetTimezone != "" {
			converterArgs.TargetTimezone = pulumi.StringPtr(arm.TargetTimezone)
		}
		config.DateTimeConverters = cloudwatch.LogTransformerTransformerConfigDateTimeConverterArray{converterArgs}
	}

	if arm := processor.DeleteKeys; arm != nil {
		config.DeleteKeys = cloudwatch.LogTransformerTransformerConfigDeleteKeyArray{
			cloudwatch.LogTransformerTransformerConfigDeleteKeyArgs{
				WithKeys: pulumi.ToStringArray(arm.WithKeys),
			},
		}
	}

	if arm := processor.Grok; arm != nil {
		grokArgs := &cloudwatch.LogTransformerTransformerConfigGrokArgs{
			Match: pulumi.String(arm.Match),
		}
		if arm.Source != "" {
			grokArgs.Source = pulumi.StringPtr(arm.Source)
		}
		config.Grok = grokArgs
	}

	if arm := processor.ListToMap; arm != nil {
		listToMapArgs := cloudwatch.LogTransformerTransformerConfigListToMapArgs{
			Source:  pulumi.String(arm.Source),
			Key:     pulumi.String(arm.Key),
			Flatten: pulumi.BoolPtr(arm.Flatten),
		}
		if arm.ValueKey != "" {
			listToMapArgs.ValueKey = pulumi.StringPtr(arm.ValueKey)
		}
		if arm.Target != "" {
			listToMapArgs.Target = pulumi.StringPtr(arm.Target)
		}
		if arm.FlattenedElement != "" {
			listToMapArgs.FlattenedElement = pulumi.StringPtr(arm.FlattenedElement)
		}
		config.ListToMaps = cloudwatch.LogTransformerTransformerConfigListToMapArray{listToMapArgs}
	}

	if arm := processor.LowerCaseString; arm != nil {
		config.LowerCaseStrings = cloudwatch.LogTransformerTransformerConfigLowerCaseStringArray{
			cloudwatch.LogTransformerTransformerConfigLowerCaseStringArgs{
				WithKeys: pulumi.ToStringArray(arm.WithKeys),
			},
		}
	}

	if arm := processor.MoveKeys; arm != nil {
		entries := make(cloudwatch.LogTransformerTransformerConfigMoveKeyEntryArray, 0, len(arm.Entries))
		for _, entry := range arm.Entries {
			entries = append(entries, cloudwatch.LogTransformerTransformerConfigMoveKeyEntryArgs{
				Source:            pulumi.String(entry.Source),
				Target:            pulumi.String(entry.Target),
				OverwriteIfExists: pulumi.BoolPtr(entry.OverwriteIfExists),
			})
		}
		config.MoveKeys = cloudwatch.LogTransformerTransformerConfigMoveKeyArray{
			cloudwatch.LogTransformerTransformerConfigMoveKeyArgs{Entries: entries},
		}
	}

	if arm := processor.ParseCloudfront; arm != nil {
		parseArgs := &cloudwatch.LogTransformerTransformerConfigParseCloudfrontArgs{}
		if arm.Source != "" {
			parseArgs.Source = pulumi.StringPtr(arm.Source)
		}
		config.ParseCloudfront = parseArgs
	}

	if arm := processor.ParseJson; arm != nil {
		parseArgs := cloudwatch.LogTransformerTransformerConfigParseJsonArgs{}
		if arm.Source != "" {
			parseArgs.Source = pulumi.StringPtr(arm.Source)
		}
		if arm.Destination != "" {
			parseArgs.Destination = pulumi.StringPtr(arm.Destination)
		}
		config.ParseJsons = cloudwatch.LogTransformerTransformerConfigParseJsonArray{parseArgs}
	}

	if arm := processor.ParseKeyValue; arm != nil {
		parseArgs := cloudwatch.LogTransformerTransformerConfigParseKeyValueArgs{
			OverwriteIfExists: pulumi.BoolPtr(arm.OverwriteIfExists),
		}
		if arm.Source != "" {
			parseArgs.Source = pulumi.StringPtr(arm.Source)
		}
		if arm.Destination != "" {
			parseArgs.Destination = pulumi.StringPtr(arm.Destination)
		}
		if arm.FieldDelimiter != "" {
			parseArgs.FieldDelimiter = pulumi.StringPtr(arm.FieldDelimiter)
		}
		if arm.KeyValueDelimiter != "" {
			parseArgs.KeyValueDelimiter = pulumi.StringPtr(arm.KeyValueDelimiter)
		}
		if arm.KeyPrefix != "" {
			parseArgs.KeyPrefix = pulumi.StringPtr(arm.KeyPrefix)
		}
		if arm.NonMatchValue != "" {
			parseArgs.NonMatchValue = pulumi.StringPtr(arm.NonMatchValue)
		}
		config.ParseKeyValues = cloudwatch.LogTransformerTransformerConfigParseKeyValueArray{parseArgs}
	}

	if arm := processor.ParsePostgres; arm != nil {
		parseArgs := &cloudwatch.LogTransformerTransformerConfigParsePostgresArgs{}
		if arm.Source != "" {
			parseArgs.Source = pulumi.StringPtr(arm.Source)
		}
		config.ParsePostgres = parseArgs
	}

	if arm := processor.ParseRoute53; arm != nil {
		parseArgs := &cloudwatch.LogTransformerTransformerConfigParseRoute53Args{}
		if arm.Source != "" {
			parseArgs.Source = pulumi.StringPtr(arm.Source)
		}
		config.ParseRoute53 = parseArgs
	}

	if arm := processor.ParseToOcsf; arm != nil {
		parseArgs := &cloudwatch.LogTransformerTransformerConfigParseToOcsfArgs{
			EventSource: pulumi.String(arm.EventSource),
			OcsfVersion: pulumi.String(arm.OcsfVersion),
		}
		if arm.Source != "" {
			parseArgs.Source = pulumi.StringPtr(arm.Source)
		}
		config.ParseToOcsf = parseArgs
	}

	if arm := processor.ParseVpc; arm != nil {
		parseArgs := &cloudwatch.LogTransformerTransformerConfigParseVpcArgs{}
		if arm.Source != "" {
			parseArgs.Source = pulumi.StringPtr(arm.Source)
		}
		config.ParseVpc = parseArgs
	}

	if arm := processor.ParseWaf; arm != nil {
		parseArgs := &cloudwatch.LogTransformerTransformerConfigParseWafArgs{}
		if arm.Source != "" {
			parseArgs.Source = pulumi.StringPtr(arm.Source)
		}
		config.ParseWaf = parseArgs
	}

	if arm := processor.RenameKeys; arm != nil {
		entries := make(cloudwatch.LogTransformerTransformerConfigRenameKeyEntryArray, 0, len(arm.Entries))
		for _, entry := range arm.Entries {
			entries = append(entries, cloudwatch.LogTransformerTransformerConfigRenameKeyEntryArgs{
				Key:               pulumi.String(entry.Key),
				RenameTo:          pulumi.String(entry.RenameTo),
				OverwriteIfExists: pulumi.BoolPtr(entry.OverwriteIfExists),
			})
		}
		config.RenameKeys = cloudwatch.LogTransformerTransformerConfigRenameKeyArray{
			cloudwatch.LogTransformerTransformerConfigRenameKeyArgs{Entries: entries},
		}
	}

	if arm := processor.SplitString; arm != nil {
		entries := make(cloudwatch.LogTransformerTransformerConfigSplitStringEntryArray, 0, len(arm.Entries))
		for _, entry := range arm.Entries {
			entries = append(entries, cloudwatch.LogTransformerTransformerConfigSplitStringEntryArgs{
				Source:    pulumi.String(entry.Source),
				Delimiter: pulumi.String(entry.Delimiter),
			})
		}
		config.SplitStrings = cloudwatch.LogTransformerTransformerConfigSplitStringArray{
			cloudwatch.LogTransformerTransformerConfigSplitStringArgs{Entries: entries},
		}
	}

	if arm := processor.SubstituteString; arm != nil {
		entries := make(cloudwatch.LogTransformerTransformerConfigSubstituteStringEntryArray, 0, len(arm.Entries))
		for _, entry := range arm.Entries {
			entries = append(entries, cloudwatch.LogTransformerTransformerConfigSubstituteStringEntryArgs{
				Source: pulumi.String(entry.Source),
				From:   pulumi.String(entry.From),
				To:     pulumi.String(entry.To),
			})
		}
		config.SubstituteStrings = cloudwatch.LogTransformerTransformerConfigSubstituteStringArray{
			cloudwatch.LogTransformerTransformerConfigSubstituteStringArgs{Entries: entries},
		}
	}

	if arm := processor.TrimString; arm != nil {
		config.TrimStrings = cloudwatch.LogTransformerTransformerConfigTrimStringArray{
			cloudwatch.LogTransformerTransformerConfigTrimStringArgs{
				WithKeys: pulumi.ToStringArray(arm.WithKeys),
			},
		}
	}

	if arm := processor.TypeConverter; arm != nil {
		entries := make(cloudwatch.LogTransformerTransformerConfigTypeConverterEntryArray, 0, len(arm.Entries))
		for _, entry := range arm.Entries {
			entries = append(entries, cloudwatch.LogTransformerTransformerConfigTypeConverterEntryArgs{
				Key:  pulumi.String(entry.Key),
				Type: pulumi.String(entry.Type),
			})
		}
		config.TypeConverters = cloudwatch.LogTransformerTransformerConfigTypeConverterArray{
			cloudwatch.LogTransformerTransformerConfigTypeConverterArgs{Entries: entries},
		}
	}

	if arm := processor.UpperCaseString; arm != nil {
		config.UpperCaseStrings = cloudwatch.LogTransformerTransformerConfigUpperCaseStringArray{
			cloudwatch.LogTransformerTransformerConfigUpperCaseStringArgs{
				WithKeys: pulumi.ToStringArray(arm.WithKeys),
			},
		}
	}

	return config
}
