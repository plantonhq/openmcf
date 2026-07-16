package awslambdaeventsourcemappingv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsLambdaEventSourceMappingSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsLambdaEventSourceMappingSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func intPtr(v int32) *int32 { return &v }

var _ = ginkgo.Describe("AwsLambdaEventSourceMappingSpec validations", func() {
	var spec *AwsLambdaEventSourceMappingSpec

	ginkgo.BeforeEach(func() {
		spec = &AwsLambdaEventSourceMappingSpec{
			Region:         "us-west-2",
			FunctionArn:    literal("arn:aws:lambda:us-west-2:123456789012:function:order-worker"),
			EventSourceArn: literal("arn:aws:sqs:us-west-2:123456789012:orders"),
		}
	})

	// -----------------------------------------------------------------
	// Happy paths
	// -----------------------------------------------------------------

	ginkgo.It("accepts a minimal SQS mapping", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a tuned SQS mapping", func() {
		spec.BatchSize = 100
		spec.MaximumBatchingWindowSeconds = 30
		spec.ScalingMaxConcurrency = 50
		spec.FunctionResponseTypes = []string{"ReportBatchItemFailures"}
		spec.Metrics = []string{"EventCount"}
		spec.Filters = []*AwsLambdaEventSourceMappingFilter{
			{Pattern: `{"body":{"type":["order.created"]}}`},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a full-featured Kinesis mapping", func() {
		spec.EventSourceArn = literal("arn:aws:kinesis:us-west-2:123456789012:stream/clicks")
		spec.StartingPosition = "TRIM_HORIZON"
		spec.ParallelizationFactor = 4
		spec.MaximumRecordAgeSeconds = 3600
		spec.MaximumRetryAttempts = intPtr(5)
		spec.BisectBatchOnFunctionError = true
		spec.TumblingWindowSeconds = 60
		spec.OnFailureDestinationArn = literal("arn:aws:sqs:us-west-2:123456789012:clicks-failures")
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts an MSK mapping", func() {
		spec.EventSourceArn = literal("arn:aws:kafka:us-west-2:123456789012:cluster/events/abc-123")
		spec.StartingPosition = "LATEST"
		spec.Topics = []string{"orders"}
		spec.KafkaConsumerGroupId = "order-workers"
		spec.ProvisionedPollers = &AwsLambdaEventSourceMappingProvisionedPollers{MinimumPollers: 2, MaximumPollers: 10}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a self-managed Kafka mapping with SASL auth", func() {
		spec.EventSourceArn = nil
		spec.SelfManagedKafka = &AwsLambdaEventSourceMappingSelfManagedKafka{
			BootstrapServers: []string{"kafka1.internal:9092", "kafka2.internal:9092"},
		}
		spec.StartingPosition = "TRIM_HORIZON"
		spec.Topics = []string{"orders"}
		spec.SourceAccessConfigurations = []*AwsLambdaEventSourceMappingSourceAccess{
			{Type: "SASL_SCRAM_512_AUTH", Uri: "arn:aws:secretsmanager:us-west-2:123456789012:secret:kafka-sasl"},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a DocumentDB mapping", func() {
		spec.EventSourceArn = literal("arn:aws:rds:us-west-2:123456789012:cluster:orders-docdb")
		spec.StartingPosition = "LATEST"
		spec.DocumentDb = &AwsLambdaEventSourceMappingDocumentDb{
			DatabaseName: "orders",
			FullDocument: "UpdateLookup",
		}
		spec.SourceAccessConfigurations = []*AwsLambdaEventSourceMappingSourceAccess{
			{Type: "BASIC_AUTH", Uri: "arn:aws:secretsmanager:us-west-2:123456789012:secret:docdb-creds"},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Event source (exactly one)
	// -----------------------------------------------------------------

	ginkgo.It("rejects a mapping with no event source", func() {
		spec.EventSourceArn = nil
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a mapping with both event sources", func() {
		spec.SelfManagedKafka = &AwsLambdaEventSourceMappingSelfManagedKafka{
			BootstrapServers: []string{"kafka1.internal:9092"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects self-managed Kafka without bootstrap servers", func() {
		spec.EventSourceArn = nil
		spec.SelfManagedKafka = &AwsLambdaEventSourceMappingSelfManagedKafka{}
		spec.StartingPosition = "TRIM_HORIZON"
		spec.Topics = []string{"orders"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects self-managed Kafka without a starting position", func() {
		spec.EventSourceArn = nil
		spec.SelfManagedKafka = &AwsLambdaEventSourceMappingSelfManagedKafka{
			BootstrapServers: []string{"kafka1.internal:9092"},
		}
		spec.Topics = []string{"orders"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Stream tuning couplings
	// -----------------------------------------------------------------

	ginkgo.It("rejects a timestamp without AT_TIMESTAMP", func() {
		spec.StartingPosition = "LATEST"
		spec.StartingPositionTimestamp = "2026-07-04T00:00:00Z"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects AT_TIMESTAMP without a timestamp", func() {
		spec.StartingPosition = "AT_TIMESTAMP"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an out-of-range record age", func() {
		spec.MaximumRecordAgeSeconds = 30
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts the -1 never-age-out sentinel", func() {
		spec.MaximumRecordAgeSeconds = -1
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects retry attempts above the AWS ceiling", func() {
		spec.MaximumRetryAttempts = intPtr(10001)
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a parallelization factor above 10", func() {
		spec.ParallelizationFactor = 11
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Consumption control
	// -----------------------------------------------------------------

	ginkgo.It("rejects a scaling concurrency below 2", func() {
		spec.ScalingMaxConcurrency = 1
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an unknown function response type", func() {
		spec.FunctionResponseTypes = []string{"RetryWholeBatch"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects more than ten filters", func() {
		for i := 0; i < 11; i++ {
			spec.Filters = append(spec.Filters, &AwsLambdaEventSourceMappingFilter{Pattern: `{"a":[1]}`})
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -----------------------------------------------------------------
	// Kafka couplings
	// -----------------------------------------------------------------

	ginkgo.It("rejects provisioned pollers on a non-Kafka source", func() {
		spec.ProvisionedPollers = &AwsLambdaEventSourceMappingProvisionedPollers{MaximumPollers: 10}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a schema registry on a non-Kafka source", func() {
		spec.SchemaRegistry = &AwsLambdaEventSourceMappingSchemaRegistry{
			Uri:               "https://registry.example.com",
			EventRecordFormat: "JSON",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects inverted poller bounds", func() {
		spec.EventSourceArn = literal("arn:aws:kafka:us-west-2:123456789012:cluster/events/abc-123")
		spec.StartingPosition = "LATEST"
		spec.Topics = []string{"orders"}
		spec.ProvisionedPollers = &AwsLambdaEventSourceMappingProvisionedPollers{MinimumPollers: 20, MaximumPollers: 5}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects mixing mq_queue with Kafka topics", func() {
		spec.Topics = []string{"orders"}
		spec.MqQueue = "orders-queue"
		spec.StartingPosition = "LATEST"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an unknown source access type", func() {
		spec.SourceAccessConfigurations = []*AwsLambdaEventSourceMappingSourceAccess{
			{Type: "OAUTH", Uri: "arn:aws:secretsmanager:us-west-2:123456789012:secret:x"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})
})
