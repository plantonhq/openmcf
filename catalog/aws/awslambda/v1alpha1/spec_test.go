package awslambdav1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsLambdaSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsLambdaSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func intPtr(v int32) *int32 { return &v }

var _ = ginkgo.Describe("AwsLambdaSpec validations", func() {
	var spec *AwsLambdaSpec

	ginkgo.BeforeEach(func() {
		spec = &AwsLambdaSpec{
			Region:  "us-west-2",
			RoleArn: literal("arn:aws:iam::123456789012:role/order-api-exec"),
			S3: &AwsLambdaS3Code{
				Bucket: literal("artifacts-bucket"),
				Key:    "functions/order-api.zip",
			},
			Runtime: "python3.13",
			Handler: "app.handler",
		}
	})

	// -----------------------------------------------------------------
	// Happy paths
	// -----------------------------------------------------------------

	ginkgo.It("accepts a minimal valid zip function", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a minimal valid container-image function", func() {
		spec.S3 = nil
		spec.Runtime = ""
		spec.Handler = ""
		spec.ImageUri = "123456789012.dkr.ecr.us-west-2.amazonaws.com/order-api:v42"
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a full-featured zip function", func() {
		spec.Description = "order intake API"
		spec.SourceCodeHash = "q0Wep3z1sTImroLBJKcRuBIfqkKing9UVePUFmSpjKQ="
		spec.SourceKmsKeyArn = literal("arn:aws:kms:us-west-2:123456789012:key/1234abcd")
		spec.Architecture = "arm64"
		spec.MemorySizeMb = 1024
		spec.TimeoutSeconds = 30
		spec.EphemeralStorageMb = 2048
		spec.Environment = map[string]string{"STAGE": "prod"}
		spec.KmsKeyArn = literal("arn:aws:kms:us-west-2:123456789012:key/5678efgh")
		spec.SubnetIds = []*foreignkeyv1.StringValueOrRef{literal("subnet-aaa"), literal("subnet-bbb")}
		spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{literal("sg-ccc")}
		spec.Ipv6AllowedForDualStack = true
		spec.DeadLetterTargetArn = literal("arn:aws:sqs:us-west-2:123456789012:order-api-dlq")
		spec.TracingMode = "Active"
		spec.FileSystemConfig = &AwsLambdaFileSystemConfig{
			AccessPointArn: literal("arn:aws:elasticfilesystem:us-west-2:123456789012:access-point/fsap-1"),
			LocalMountPath: "/mnt/models",
		}
		spec.LayerArns = []*foreignkeyv1.StringValueOrRef{literal("arn:aws:lambda:us-west-2:123456789012:layer:shared:3")}
		spec.Publish = true
		spec.ReservedConcurrentExecutions = intPtr(50)
		spec.SnapStart = true
		spec.LoggingConfig = &AwsLambdaLoggingConfig{
			LogFormat:           "JSON",
			ApplicationLogLevel: "INFO",
			SystemLogLevel:      "WARN",
			LogGroup:            literal("/platform/order-api"),
		}
		spec.CodeSigningConfigArn = "arn:aws:lambda:us-west-2:123456789012:code-signing-config:csc-1"
		// Provisioned concurrency and canary weights are mutually exclusive
		// per alias in AWS, so the fixture carries one alias of each shape.
		spec.Aliases = []*AwsLambdaAlias{
			{
				Name:                            "live",
				FunctionVersion:                 "1",
				ProvisionedConcurrentExecutions: intPtr(5),
			},
			{
				Name:                            "canary",
				FunctionVersion:                 "1",
				RoutingAdditionalVersionWeights: map[string]float64{"2": 0.1},
			},
		}
		spec.FunctionUrl = &AwsLambdaFunctionUrl{
			AuthorizationType: "AWS_IAM",
			InvokeMode:        "RESPONSE_STREAM",
			Cors: &AwsLambdaFunctionUrlCors{
				AllowOrigins:  []string{"https://app.example.com"},
				AllowMethods:  []string{"GET", "POST"},
				MaxAgeSeconds: 3600,
			},
		}
		spec.InvokePermissions = []*AwsLambdaInvokePermission{{
			StatementId: "allow-uploads-bucket",
			Principal:   "s3.amazonaws.com",
			SourceArn:   "arn:aws:s3:::uploads-bucket",
		}}
		spec.AsyncInvokeConfig = &AwsLambdaAsyncInvokeConfig{
			MaximumRetryAttempts:    intPtr(1),
			MaximumEventAgeSeconds:  3600,
			OnFailureDestinationArn: literal("arn:aws:sqs:us-west-2:123456789012:order-api-failures"),
		}
		spec.RecursiveLoop = "Terminate"
		spec.RuntimeManagement = &AwsLambdaRuntimeManagement{UpdateRuntimeOn: "Auto"}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Code source (exactly one)
	// -----------------------------------------------------------------

	ginkgo.It("rejects a function with no code source", func() {
		spec.S3 = nil
		spec.Runtime = ""
		spec.Handler = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a function with both code sources", func() {
		spec.ImageUri = "123456789012.dkr.ecr.us-west-2.amazonaws.com/order-api:v42"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a zip function missing runtime", func() {
		spec.Runtime = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a zip function missing handler", func() {
		spec.Handler = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an image function that sets runtime/handler", func() {
		spec.S3 = nil
		spec.ImageUri = "123456789012.dkr.ecr.us-west-2.amazonaws.com/order-api:v42"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects image_config on a zip function", func() {
		spec.ImageConfig = &AwsLambdaImageConfig{Command: []string{"app.handler"}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects source_kms_key_arn on an image function", func() {
		spec.S3 = nil
		spec.Runtime = ""
		spec.Handler = ""
		spec.ImageUri = "123456789012.dkr.ecr.us-west-2.amazonaws.com/order-api:v42"
		spec.SourceKmsKeyArn = literal("arn:aws:kms:us-west-2:123456789012:key/1234abcd")
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects code signing on an image function", func() {
		spec.S3 = nil
		spec.Runtime = ""
		spec.Handler = ""
		spec.ImageUri = "123456789012.dkr.ecr.us-west-2.amazonaws.com/order-api:v42"
		spec.CodeSigningConfigArn = "arn:aws:lambda:us-west-2:123456789012:code-signing-config:csc-1"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Execution environment bounds
	// -----------------------------------------------------------------

	ginkgo.It("rejects memory below the AWS floor", func() {
		spec.MemorySizeMb = 64
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects timeout above the AWS ceiling", func() {
		spec.TimeoutSeconds = 901
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects ephemeral storage below the AWS floor", func() {
		spec.EphemeralStorageMb = 256
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an unknown architecture", func() {
		spec.Architecture = "riscv64"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts reserved concurrency of zero (the kill switch)", func() {
		spec.ReservedConcurrentExecutions = intPtr(0)
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects negative reserved concurrency", func() {
		spec.ReservedConcurrentExecutions = intPtr(-1)
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// VPC attachment coupling
	// -----------------------------------------------------------------

	ginkgo.It("rejects subnets without security groups", func() {
		spec.SubnetIds = []*foreignkeyv1.StringValueOrRef{literal("subnet-aaa")}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects security groups without subnets", func() {
		spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{literal("sg-ccc")}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects dual-stack IPv6 egress without a VPC attachment", func() {
		spec.Ipv6AllowedForDualStack = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an EFS mount without a VPC attachment", func() {
		spec.FileSystemConfig = &AwsLambdaFileSystemConfig{
			AccessPointArn: literal("arn:aws:elasticfilesystem:us-west-2:123456789012:access-point/fsap-1"),
			LocalMountPath: "/mnt/models",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an EFS mount path outside /mnt", func() {
		spec.SubnetIds = []*foreignkeyv1.StringValueOrRef{literal("subnet-aaa")}
		spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{literal("sg-ccc")}
		spec.FileSystemConfig = &AwsLambdaFileSystemConfig{
			AccessPointArn: literal("arn:aws:elasticfilesystem:us-west-2:123456789012:access-point/fsap-1"),
			LocalMountPath: "/data/models",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Versioning couplings
	// -----------------------------------------------------------------

	ginkgo.It("rejects SnapStart without publish", func() {
		spec.SnapStart = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects aliases without publish", func() {
		spec.Aliases = []*AwsLambdaAlias{{Name: "live", FunctionVersion: "1"}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects duplicate alias names", func() {
		spec.Publish = true
		spec.Aliases = []*AwsLambdaAlias{
			{Name: "live", FunctionVersion: "1"},
			{Name: "live", FunctionVersion: "2"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects provisioned concurrency on a weighted (canary) alias", func() {
		spec.Publish = true
		spec.Aliases = []*AwsLambdaAlias{{
			Name:                            "live",
			FunctionVersion:                 "1",
			RoutingAdditionalVersionWeights: map[string]float64{"2": 0.1},
			ProvisionedConcurrentExecutions: intPtr(5),
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects more than one additional routing version on an alias", func() {
		spec.Publish = true
		spec.Aliases = []*AwsLambdaAlias{{
			Name:                            "live",
			FunctionVersion:                 "1",
			RoutingAdditionalVersionWeights: map[string]float64{"2": 0.1, "3": 0.1},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a routing weight above 1.0", func() {
		spec.Publish = true
		spec.Aliases = []*AwsLambdaAlias{{
			Name:                            "live",
			FunctionVersion:                 "1",
			RoutingAdditionalVersionWeights: map[string]float64{"2": 1.5},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects provisioned concurrency on a $LATEST alias", func() {
		spec.Publish = true
		spec.Aliases = []*AwsLambdaAlias{{
			Name:                            "dev",
			FunctionVersion:                 "$LATEST",
			ProvisionedConcurrentExecutions: intPtr(2),
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Logging couplings
	// -----------------------------------------------------------------

	ginkgo.It("rejects log-level filtering without JSON format", func() {
		spec.LoggingConfig = &AwsLambdaLoggingConfig{ApplicationLogLevel: "INFO"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an unknown log level", func() {
		spec.LoggingConfig = &AwsLambdaLoggingConfig{LogFormat: "JSON", SystemLogLevel: "TRACE"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Satellites
	// -----------------------------------------------------------------

	ginkgo.It("rejects duplicate permission statement ids", func() {
		spec.InvokePermissions = []*AwsLambdaInvokePermission{
			{StatementId: "grant", Principal: "s3.amazonaws.com"},
			{StatementId: "grant", Principal: "sns.amazonaws.com"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects function_url_auth_type on a non-URL permission action", func() {
		spec.InvokePermissions = []*AwsLambdaInvokePermission{{
			StatementId:         "grant",
			Principal:           "*",
			FunctionUrlAuthType: "NONE",
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a Manual runtime pin without a version ARN", func() {
		spec.RuntimeManagement = &AwsLambdaRuntimeManagement{UpdateRuntimeOn: "Manual"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a runtime version ARN without the Manual mode", func() {
		spec.RuntimeManagement = &AwsLambdaRuntimeManagement{
			UpdateRuntimeOn:   "Auto",
			RuntimeVersionArn: "arn:aws:lambda:us-west-2::runtime:abc",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an async retry count above the AWS ceiling", func() {
		spec.AsyncInvokeConfig = &AwsLambdaAsyncInvokeConfig{MaximumRetryAttempts: intPtr(3)}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an unknown recursive_loop value", func() {
		spec.RecursiveLoop = "Detect"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a function URL with an unknown authorization type", func() {
		spec.FunctionUrl = &AwsLambdaFunctionUrl{AuthorizationType: "COGNITO"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects more than five layers", func() {
		for i := 0; i < 6; i++ {
			spec.LayerArns = append(spec.LayerArns, literal("arn:aws:lambda:us-west-2:123456789012:layer:l:1"))
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Execution platform: Managed Instances, durability, tenancy
	// -----------------------------------------------------------------

	ginkgo.It("accepts a Managed Instances function with environment sizing", func() {
		spec.ManagedInstances = &AwsLambdaManagedInstances{
			CapacityProviderArn:          "arn:aws:lambda:us-west-2:123456789012:capacity-provider:steady",
			MemoryGibPerVcpu:             4.0,
			MaxConcurrencyPerEnvironment: 8,
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects managed_instances without a capacity provider ARN", func() {
		spec.ManagedInstances = &AwsLambdaManagedInstances{MemoryGibPerVcpu: 4.0}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts memory above the standard ceiling for Managed Instances sizes", func() {
		spec.MemorySizeMb = 32768
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects memory above the platform maximum", func() {
		spec.MemorySizeMb = 32769
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts a durable function with retention", func() {
		spec.DurableConfig = &AwsLambdaDurableConfig{
			ExecutionTimeoutSeconds: 86400,
			RetentionPeriodDays:     30,
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects a durable config without an execution timeout", func() {
		spec.DurableConfig = &AwsLambdaDurableConfig{RetentionPeriodDays: 30}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a durable execution timeout beyond 366 days", func() {
		spec.DurableConfig = &AwsLambdaDurableConfig{ExecutionTimeoutSeconds: 31622401}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a durable retention period beyond 90 days", func() {
		spec.DurableConfig = &AwsLambdaDurableConfig{
			ExecutionTimeoutSeconds: 3600,
			RetentionPeriodDays:     91,
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts PER_TENANT isolation and rejects unknown modes", func() {
		spec.TenantIsolationMode = "PER_TENANT"
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		spec.TenantIsolationMode = "SHARED"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// publish_to, code_sha256, and per-qualifier scaling configs
	// -----------------------------------------------------------------

	ginkgo.It("accepts publish_to LATEST_PUBLISHED and rejects other values", func() {
		spec.PublishTo = "LATEST_PUBLISHED"
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		spec.PublishTo = "LATEST"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts a deployed-package digest", func() {
		spec.CodeSha256 = "q0Wep3z1sTImroLBJKcRuBIfqkKing9UVePUFmSpjKQ="
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts scaling configs for a version and the published head", func() {
		spec.ScalingConfigs = []*AwsLambdaScalingConfig{
			{Qualifier: "3", MinExecutionEnvironments: intPtr(1)},
			{Qualifier: "$LATEST.PUBLISHED", MaxExecutionEnvironments: intPtr(20)},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects a scaling config qualified by alias name", func() {
		spec.ScalingConfigs = []*AwsLambdaScalingConfig{
			{Qualifier: "live", MinExecutionEnvironments: intPtr(1)},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an empty scaling config", func() {
		spec.ScalingConfigs = []*AwsLambdaScalingConfig{{Qualifier: "3"}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects scaling bounds out of order", func() {
		spec.ScalingConfigs = []*AwsLambdaScalingConfig{{
			Qualifier:                "3",
			MinExecutionEnvironments: intPtr(10),
			MaxExecutionEnvironments: intPtr(5),
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects duplicate scaling config qualifiers", func() {
		spec.ScalingConfigs = []*AwsLambdaScalingConfig{
			{Qualifier: "3", MinExecutionEnvironments: intPtr(1)},
			{Qualifier: "3", MaxExecutionEnvironments: intPtr(5)},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Per-qualifier satellite targeting
	// -----------------------------------------------------------------

	ginkgo.It("accepts a function URL qualified by a declared alias", func() {
		spec.Publish = true
		spec.Aliases = []*AwsLambdaAlias{{Name: "live", FunctionVersion: "1"}}
		spec.FunctionUrl = &AwsLambdaFunctionUrl{
			AuthorizationType: "AWS_IAM",
			Qualifier:         "live",
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects a function URL qualifier that names no alias", func() {
		spec.Publish = true
		spec.Aliases = []*AwsLambdaAlias{{Name: "live", FunctionVersion: "1"}}
		spec.FunctionUrl = &AwsLambdaFunctionUrl{
			AuthorizationType: "AWS_IAM",
			Qualifier:         "staging",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts qualified permissions, async config, and runtime management", func() {
		spec.InvokePermissions = []*AwsLambdaInvokePermission{{
			StatementId:           "allow-alexa",
			Principal:             "alexa-appkit.amazon.com",
			Qualifier:             "live",
			EventSourceToken:      "amzn1.ask.skill.12345678-1234-1234-1234-123456789012",
			InvokedViaFunctionUrl: false,
		}}
		spec.AsyncInvokeConfig = &AwsLambdaAsyncInvokeConfig{
			MaximumRetryAttempts: intPtr(1),
			Qualifier:            "live",
		}
		spec.RuntimeManagement = &AwsLambdaRuntimeManagement{
			UpdateRuntimeOn: "FunctionUpdate",
			Qualifier:       "live",
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})
})
