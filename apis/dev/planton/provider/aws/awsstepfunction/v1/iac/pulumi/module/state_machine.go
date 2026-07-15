package module

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sfn"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func stateMachine(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	// The ASL definition arrives as a protobuf Struct; the AWS API takes it
	// as a JSON document. ASL key casing (StartAt, States, ...) survives the
	// round trip because Struct preserves map keys verbatim.
	definitionMap := spec.Definition.AsMap()
	definitionJSON, err := json.Marshal(definitionMap)
	if err != nil {
		return errors.Wrap(err, "failed to serialize state machine definition to JSON")
	}

	// STANDARD when unset. Changing the type replaces the state machine
	// (AWS ForceNew), which the spec documents.
	smType := "STANDARD"
	if spec.Type != "" {
		smType = spec.Type
	}

	args := &sfn.StateMachineArgs{
		// The cloud name is metadata.name -- the same basis the Terraform
		// module uses, so both engines create the same physical identity.
		Name:       pulumi.StringPtr(locals.Target.Metadata.Name),
		Definition: pulumi.String(string(definitionJSON)),
		RoleArn:    pulumi.String(spec.RoleArn.GetValue()),
		Type:       pulumi.StringPtr(smType),
		// Publish an immutable version on create and on every configuration
		// change. The latest version's ARN is exported as a stack output so
		// consumers can pin executions to a snapshot instead of the mutable
		// state machine.
		Publish: pulumi.BoolPtr(spec.Publish),
		Tags:    pulumi.ToStringMap(locals.AwsTags),
	}

	// X-Ray tracing is a single toggle; the role must be able to put trace
	// segments (xray:PutTraceSegments / PutTelemetryRecords).
	if spec.TracingEnabled {
		args.TracingConfiguration = &sfn.StateMachineTracingConfigurationArgs{
			Enabled: pulumi.BoolPtr(true),
		}
	}

	// Execution-history logging. Only rendered for a real level -- AWS treats
	// level OFF and an absent block identically.
	if spec.Logging != nil && spec.Logging.Level != "" && spec.Logging.Level != "OFF" {
		logArgs := &sfn.StateMachineLoggingConfigurationArgs{
			Level:                pulumi.StringPtr(spec.Logging.Level),
			IncludeExecutionData: pulumi.BoolPtr(spec.Logging.IncludeExecutionData),
		}

		if spec.Logging.LogDestination.GetValue() != "" {
			logDest := spec.Logging.LogDestination.GetValue()
			// AWS requires the CloudWatch log group ARN to end with ":*".
			// Referenced log-group ARNs arrive without the suffix, so append
			// it -- users should never have to know about this quirk.
			if !strings.HasSuffix(logDest, ":*") {
				logDest = logDest + ":*"
			}
			logArgs.LogDestination = pulumi.StringPtr(logDest)
		}

		args.LoggingConfiguration = logArgs
	}

	// Customer-managed KMS encryption. AWS's other arm (AWS_OWNED_KEY) is the
	// no-block default, so the spec models exactly one honest shape: a block
	// with a key means customer-managed.
	if spec.Encryption != nil && spec.Encryption.KmsKeyId.GetValue() != "" {
		encArgs := &sfn.StateMachineEncryptionConfigurationArgs{
			Type:     pulumi.StringPtr("CUSTOMER_MANAGED_KMS_KEY"),
			KmsKeyId: pulumi.StringPtr(spec.Encryption.KmsKeyId.GetValue()),
		}
		// A zero reuse period means "let AWS default" (300s); AWS accepts
		// 60-900.
		if spec.Encryption.KmsDataKeyReusePeriodSeconds != 0 {
			encArgs.KmsDataKeyReusePeriodSeconds = pulumi.IntPtr(int(spec.Encryption.KmsDataKeyReusePeriodSeconds))
		}
		args.EncryptionConfiguration = encArgs
	}

	sm, err := sfn.NewStateMachine(ctx, locals.Target.Metadata.Name, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create Step Functions state machine")
	}

	// Export outputs matching AwsStepFunctionStackOutputs.
	ctx.Export(OpStateMachineArn, sm.Arn)
	ctx.Export(OpStateMachineName, sm.Name)
	ctx.Export(OpStateMachineVersionArn, sm.StateMachineVersionArn)
	ctx.Export(OpRevisionId, sm.RevisionId)
	ctx.Export(OpStatus, sm.Status)
	ctx.Export(OpCreationDate, sm.CreationDate)

	return nil
}
