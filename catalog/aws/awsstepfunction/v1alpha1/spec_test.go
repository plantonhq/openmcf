package awsstepfunctionv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsStepFunctionSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsStepFunctionSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// helper to create a minimal valid ASL definition as a protobuf Struct.
func minimalDefinition() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]interface{}{
		"StartAt": "Hello",
		"States": map[string]interface{}{
			"Hello": map[string]interface{}{
				"Type": "Pass",
				"End":  true,
			},
		},
	})
	return s
}

// helper to create a Lambda task definition.
func lambdaTaskDefinition() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]interface{}{
		"StartAt": "ProcessOrder",
		"States": map[string]interface{}{
			"ProcessOrder": map[string]interface{}{
				"Type":     "Task",
				"Resource": "arn:aws:lambda:us-east-1:123456789012:function:process-order",
				"End":      true,
			},
		},
	})
	return s
}

var _ = ginkgo.Describe("AwsStepFunctionSpec validations", func() {
	var spec *AwsStepFunctionSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: definition + role_arn.
		spec = &AwsStepFunctionSpec{
			Region:     "us-west-2",
			Definition: minimalDefinition(),
			RoleArn:    strRef("arn:aws:iam::123456789012:role/StepFunctionsExecRole"),
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal spec (definition + role_arn only)", func() {
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts explicit STANDARD type", func() {
		spec.Type = "STANDARD"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts EXPRESS type", func() {
		spec.Type = "EXPRESS"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with publish enabled", func() {
		spec.Publish = true
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with tracing enabled", func() {
		spec.TracingEnabled = proto.Bool(true)
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with tracing explicitly disabled (the revert send)", func() {
		spec.TracingEnabled = proto.Bool(false)
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts aliases when publish is true", func() {
		spec.Publish = true
		spec.Aliases = []*AwsStepFunctionAlias{
			{Name: "prod", Description: "production traffic"},
			{Name: "canary"},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("rejects aliases without publish", func() {
		spec.Publish = false
		spec.Aliases = []*AwsStepFunctionAlias{{Name: "prod"}}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
		gomega.Expect(err.Error()).To(gomega.ContainSubstring("aliases require publish"))
	})

	ginkgo.It("rejects duplicate alias names", func() {
		spec.Publish = true
		spec.Aliases = []*AwsStepFunctionAlias{{Name: "prod"}, {Name: "prod"}}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
		gomega.Expect(err.Error()).To(gomega.ContainSubstring("must be unique"))
	})

	ginkgo.It("rejects an alias name with illegal characters", func() {
		spec.Publish = true
		spec.Aliases = []*AwsStepFunctionAlias{{Name: "prod alias"}}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a literal role_arn that is not an ARN", func() {
		spec.RoleArn = strRef("not-an-arn")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
		gomega.Expect(err.Error()).To(gomega.ContainSubstring("must be an IAM role ARN"))
	})

	ginkgo.It("accepts a spec with logging (level ALL + destination)", func() {
		spec.Logging = &AwsStepFunctionLoggingConfig{
			Level:                "ALL",
			IncludeExecutionData: true,
			LogDestination:       strRef("arn:aws:logs:us-east-1:123456789012:log-group:/aws/stepfunctions/my-workflow:*"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with logging level OFF (no destination required)", func() {
		spec.Logging = &AwsStepFunctionLoggingConfig{
			Level: "OFF",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with logging level ERROR + destination", func() {
		spec.Logging = &AwsStepFunctionLoggingConfig{
			Level:          "ERROR",
			LogDestination: strRef("arn:aws:logs:us-east-1:123456789012:log-group:/aws/stepfunctions/my-workflow:*"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with logging level FATAL + destination", func() {
		spec.Logging = &AwsStepFunctionLoggingConfig{
			Level:          "FATAL",
			LogDestination: strRef("arn:aws:logs:us-east-1:123456789012:log-group:/aws/stepfunctions/my-workflow:*"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec with customer-managed KMS encryption", func() {
		spec.Encryption = &AwsStepFunctionEncryptionConfig{
			KmsKeyId:                     strRef("arn:aws:kms:us-east-1:123456789012:key/mrk-abc123"),
			KmsDataKeyReusePeriodSeconds: 300,
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts encryption with kms_key_id only (reuse period defaults to AWS default)", func() {
		spec.Encryption = &AwsStepFunctionEncryptionConfig{
			KmsKeyId: strRef("arn:aws:kms:us-east-1:123456789012:key/mrk-abc123"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a full production-ready spec", func() {
		spec.Type = "STANDARD"
		spec.Definition = lambdaTaskDefinition()
		spec.Publish = true
		spec.TracingEnabled = proto.Bool(true)
		spec.Logging = &AwsStepFunctionLoggingConfig{
			Level:                "ALL",
			IncludeExecutionData: true,
			LogDestination:       strRef("arn:aws:logs:us-east-1:123456789012:log-group:/aws/stepfunctions/prod-workflow:*"),
		}
		spec.Encryption = &AwsStepFunctionEncryptionConfig{
			KmsKeyId:                     strRef("arn:aws:kms:us-east-1:123456789012:key/mrk-abc123"),
			KmsDataKeyReusePeriodSeconds: 600,
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Required field validations
	// -------------------------------------------------------------------------

	ginkgo.It("fails when definition is missing", func() {
		spec.Definition = nil
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when role_arn is missing", func() {
		spec.RoleArn = nil
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: type_valid_values
	// -------------------------------------------------------------------------

	ginkgo.It("fails when type has an invalid value", func() {
		spec.Type = "INVALID"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when type is lowercase", func() {
		spec.Type = "standard"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: logging_level_valid_values
	// -------------------------------------------------------------------------

	ginkgo.It("fails when logging level has an invalid value", func() {
		spec.Logging = &AwsStepFunctionLoggingConfig{
			Level:          "DEBUG",
			LogDestination: strRef("arn:aws:logs:us-east-1:123456789012:log-group:test:*"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when logging level is lowercase", func() {
		spec.Logging = &AwsStepFunctionLoggingConfig{
			Level:          "all",
			LogDestination: strRef("arn:aws:logs:us-east-1:123456789012:log-group:test:*"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: logging_destination_required_when_enabled
	// -------------------------------------------------------------------------

	ginkgo.It("fails when logging level is ALL but no log_destination", func() {
		spec.Logging = &AwsStepFunctionLoggingConfig{
			Level: "ALL",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when logging level is ERROR but no log_destination", func() {
		spec.Logging = &AwsStepFunctionLoggingConfig{
			Level: "ERROR",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when logging level is FATAL but no log_destination", func() {
		spec.Logging = &AwsStepFunctionLoggingConfig{
			Level: "FATAL",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: kms_data_key_reuse_range (on AwsStepFunctionEncryptionConfig)
	// -------------------------------------------------------------------------

	ginkgo.It("fails when kms_data_key_reuse_period_seconds is below 60", func() {
		spec.Encryption = &AwsStepFunctionEncryptionConfig{
			KmsKeyId:                     strRef("arn:aws:kms:us-east-1:123456789012:key/mrk-abc123"),
			KmsDataKeyReusePeriodSeconds: 30,
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when kms_data_key_reuse_period_seconds exceeds 900", func() {
		spec.Encryption = &AwsStepFunctionEncryptionConfig{
			KmsKeyId:                     strRef("arn:aws:kms:us-east-1:123456789012:key/mrk-abc123"),
			KmsDataKeyReusePeriodSeconds: 901,
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts kms_data_key_reuse_period_seconds at boundary 60", func() {
		spec.Encryption = &AwsStepFunctionEncryptionConfig{
			KmsKeyId:                     strRef("arn:aws:kms:us-east-1:123456789012:key/mrk-abc123"),
			KmsDataKeyReusePeriodSeconds: 60,
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts kms_data_key_reuse_period_seconds at boundary 900", func() {
		spec.Encryption = &AwsStepFunctionEncryptionConfig{
			KmsKeyId:                     strRef("arn:aws:kms:us-east-1:123456789012:key/mrk-abc123"),
			KmsDataKeyReusePeriodSeconds: 900,
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Encryption: kms_key_id is required in the encryption block
	// -------------------------------------------------------------------------

	ginkgo.It("fails when encryption block is present but kms_key_id is missing", func() {
		spec.Encryption = &AwsStepFunctionEncryptionConfig{
			KmsDataKeyReusePeriodSeconds: 300,
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})
})
