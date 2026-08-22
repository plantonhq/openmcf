package awseventbridgepipev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsEventBridgePipeSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsEventBridgePipeSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func minimalPipe() *AwsEventBridgePipeSpec {
	return &AwsEventBridgePipeSpec{
		Region:  "us-east-1",
		Source:  svr("arn:aws:sqs:us-east-1:123456789012:orders"),
		Target:  svr("arn:aws:sqs:us-east-1:123456789012:orders-processed"),
		RoleArn: svr("arn:aws:iam::123456789012:role/pipe-exec"),
	}
}

var _ = ginkgo.Describe("AwsEventBridgePipeSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal SQS-to-SQS pipe", func() {
			gomega.Expect(protovalidate.Validate(minimalPipe())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a stopped pipe", func() {
			spec := minimalPipe()
			spec.DesiredState = "STOPPED"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts SQS tuning with a filter", func() {
			spec := minimalPipe()
			spec.SourceParameters = &AwsEventBridgePipeSourceParameters{
				FilterCriteria: &AwsEventBridgePipeFilterCriteria{
					Filters: []*AwsEventBridgePipeFilter{{Pattern: `{"body":{"type":["order"]}}`}},
				},
				Sqs: &AwsEventBridgePipeSqsSourceParameters{
					BatchSize:                      proto.Int32(10),
					MaximumBatchingWindowInSeconds: proto.Int32(0),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a Kinesis source at a timestamp", func() {
			spec := minimalPipe()
			spec.Source = svr("arn:aws:kinesis:us-east-1:123456789012:stream/clicks")
			spec.SourceParameters = &AwsEventBridgePipeSourceParameters{
				Kinesis: &AwsEventBridgePipeKinesisSourceParameters{
					StartingPosition:          "AT_TIMESTAMP",
					StartingPositionTimestamp: "2026-08-01T00:00:00Z",
					MaximumRecordAgeInSeconds: proto.Int32(-1),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an enrichment with parameters", func() {
			spec := minimalPipe()
			spec.Enrichment = svr("arn:aws:lambda:us-east-1:123456789012:function:enrich")
			spec.EnrichmentParameters = &AwsEventBridgePipeEnrichmentParameters{
				InputTemplate: `{"id": <$.messageId>}`,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an ECS task target at full depth", func() {
			spec := minimalPipe()
			spec.Target = svr("arn:aws:ecs:us-east-1:123456789012:cluster/apps")
			spec.TargetParameters = &AwsEventBridgePipeTargetParameters{
				EcsTask: &AwsEventBridgePipeEcsTaskTargetParameters{
					TaskDefinitionArn: svr("arn:aws:ecs:us-east-1:123456789012:task-definition/worker:3"),
					LaunchType:        "FARGATE",
					NetworkConfiguration: &AwsEventBridgePipeEcsNetworkConfiguration{
						Subnets: []*foreignkeyv1.StringValueOrRef{svr("subnet-0123456789abcdef0")},
					},
					Overrides: &AwsEventBridgePipeEcsTaskOverrides{
						ContainerOverrides: []*AwsEventBridgePipeEcsContainerOverride{{
							Name:    "worker",
							Command: []string{"process", "--once"},
							Environment: []*AwsEventBridgePipeEcsEnvironmentVariable{
								{Name: "MODE", Value: "batch"},
							},
						}},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts MSK credentials with exactly one mechanism", func() {
			spec := minimalPipe()
			spec.Source = svr("arn:aws:kafka:us-east-1:123456789012:cluster/orders/uuid")
			spec.SourceParameters = &AwsEventBridgePipeSourceParameters{
				Msk: &AwsEventBridgePipeMskSourceParameters{
					TopicName: "orders",
					Credentials: &AwsEventBridgePipeMskCredentials{
						SaslScram_512Auth: "arn:aws:secretsmanager:us-east-1:123456789012:secret:kafka-abc",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a Redshift target with secret auth", func() {
			spec := minimalPipe()
			spec.TargetParameters = &AwsEventBridgePipeTargetParameters{
				RedshiftData: &AwsEventBridgePipeRedshiftDataTargetParameters{
					Database:         "analytics",
					Sqls:             []string{"INSERT INTO events SELECT 1"},
					SecretManagerArn: "arn:aws:secretsmanager:us-east-1:123456789012:secret:rs-abc",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts TRACE logging to CloudWatch", func() {
			spec := minimalPipe()
			spec.LogConfiguration = &AwsEventBridgePipeLogConfiguration{
				Level:                "TRACE",
				IncludeExecutionData: true,
				CloudwatchLogs: &AwsEventBridgePipeCloudWatchLogsLogDestination{
					LogGroupArn: svr("arn:aws:logs:us-east-1:123456789012:log-group:/pipes/orders"),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a pipe without a source", func() {
			spec := minimalPipe()
			spec.Source = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a pipe without a role", func() {
			spec := minimalPipe()
			spec.RoleArn = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects two source family blocks", func() {
			spec := minimalPipe()
			spec.SourceParameters = &AwsEventBridgePipeSourceParameters{
				Sqs:     &AwsEventBridgePipeSqsSourceParameters{},
				Kinesis: &AwsEventBridgePipeKinesisSourceParameters{StartingPosition: "LATEST"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects two target family blocks", func() {
			spec := minimalPipe()
			spec.TargetParameters = &AwsEventBridgePipeTargetParameters{
				Sqs:     &AwsEventBridgePipeSqsTargetParameters{MessageGroupId: "g"},
				Kinesis: &AwsEventBridgePipeKinesisTargetParameters{PartitionKey: "pk"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects enrichment parameters without an enrichment", func() {
			spec := minimalPipe()
			spec.EnrichmentParameters = &AwsEventBridgePipeEnrichmentParameters{InputTemplate: "{}"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects AT_TIMESTAMP without a timestamp", func() {
			spec := minimalPipe()
			spec.SourceParameters = &AwsEventBridgePipeSourceParameters{
				Kinesis: &AwsEventBridgePipeKinesisSourceParameters{StartingPosition: "AT_TIMESTAMP"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a timestamp on LATEST", func() {
			spec := minimalPipe()
			spec.SourceParameters = &AwsEventBridgePipeSourceParameters{
				Kinesis: &AwsEventBridgePipeKinesisSourceParameters{
					StartingPosition:          "LATEST",
					StartingPositionTimestamp: "2026-08-01T00:00:00Z",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects more than five filters", func() {
			filters := make([]*AwsEventBridgePipeFilter, 6)
			for i := range filters {
				filters[i] = &AwsEventBridgePipeFilter{Pattern: `{"a":[1]}`}
			}
			spec := minimalPipe()
			spec.SourceParameters = &AwsEventBridgePipeSourceParameters{
				FilterCriteria: &AwsEventBridgePipeFilterCriteria{Filters: filters},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects MSK credentials with both mechanisms", func() {
			spec := minimalPipe()
			spec.SourceParameters = &AwsEventBridgePipeSourceParameters{
				Msk: &AwsEventBridgePipeMskSourceParameters{
					TopicName: "orders",
					Credentials: &AwsEventBridgePipeMskCredentials{
						ClientCertificateTlsAuth: "arn:aws:secretsmanager:us-east-1:123456789012:secret:a",
						SaslScram_512Auth:        "arn:aws:secretsmanager:us-east-1:123456789012:secret:b",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a non-Secrets-Manager credential reference", func() {
			spec := minimalPipe()
			spec.SourceParameters = &AwsEventBridgePipeSourceParameters{
				Rabbitmq: &AwsEventBridgePipeRabbitMqSourceParameters{
					QueueName:            "orders",
					BasicAuthCredentials: "my-plaintext-password",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects Redshift with both auth paths", func() {
			spec := minimalPipe()
			spec.TargetParameters = &AwsEventBridgePipeTargetParameters{
				RedshiftData: &AwsEventBridgePipeRedshiftDataTargetParameters{
					Database:         "analytics",
					Sqls:             []string{"SELECT 1"},
					DbUser:           "admin",
					SecretManagerArn: "arn:aws:secretsmanager:us-east-1:123456789012:secret:rs-abc",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a logging level with no destination", func() {
			spec := minimalPipe()
			spec.LogConfiguration = &AwsEventBridgePipeLogConfiguration{Level: "ERROR"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown desired state", func() {
			spec := minimalPipe()
			spec.DesiredState = "PAUSED"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a batch job target without a job name", func() {
			spec := minimalPipe()
			spec.TargetParameters = &AwsEventBridgePipeTargetParameters{
				BatchJob: &AwsEventBridgePipeBatchJobTargetParameters{JobDefinition: "worker:3"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
