package module

import (
	awseventbridgepipev1alpha1 "github.com/plantonhq/planton/catalog/aws/awseventbridgepipe/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/pipes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildSourceParameters renders the source-family tuning (the spec's
// CEL guarantees at most one family block) plus the event filter.
func buildSourceParameters(source *awseventbridgepipev1alpha1.AwsEventBridgePipeSourceParameters) *pipes.PipeSourceParametersArgs {
	args := &pipes.PipeSourceParametersArgs{}

	if source.FilterCriteria != nil {
		filters := pipes.PipeSourceParametersFilterCriteriaFilterArray{}
		for _, filter := range source.FilterCriteria.Filters {
			filters = append(filters, &pipes.PipeSourceParametersFilterCriteriaFilterArgs{
				Pattern: pulumi.String(filter.Pattern),
			})
		}
		args.FilterCriteria = &pipes.PipeSourceParametersFilterCriteriaArgs{
			Filters: filters,
		}
	}

	if source.Sqs != nil {
		sqs := &pipes.PipeSourceParametersSqsQueueParametersArgs{}
		if source.Sqs.BatchSize != nil {
			sqs.BatchSize = pulumi.Int(int(*source.Sqs.BatchSize))
		}
		if source.Sqs.MaximumBatchingWindowInSeconds != nil {
			sqs.MaximumBatchingWindowInSeconds = pulumi.Int(int(*source.Sqs.MaximumBatchingWindowInSeconds))
		}
		args.SqsQueueParameters = sqs
	}

	if source.Kinesis != nil {
		kinesis := &pipes.PipeSourceParametersKinesisStreamParametersArgs{
			StartingPosition: pulumi.String(source.Kinesis.StartingPosition),
		}
		if source.Kinesis.StartingPositionTimestamp != "" {
			kinesis.StartingPositionTimestamp = pulumi.String(source.Kinesis.StartingPositionTimestamp)
		}
		if source.Kinesis.BatchSize != nil {
			kinesis.BatchSize = pulumi.Int(int(*source.Kinesis.BatchSize))
		}
		if source.Kinesis.MaximumBatchingWindowInSeconds != nil {
			kinesis.MaximumBatchingWindowInSeconds = pulumi.Int(int(*source.Kinesis.MaximumBatchingWindowInSeconds))
		}
		if source.Kinesis.MaximumRecordAgeInSeconds != nil {
			kinesis.MaximumRecordAgeInSeconds = pulumi.Int(int(*source.Kinesis.MaximumRecordAgeInSeconds))
		}
		if source.Kinesis.MaximumRetryAttempts != nil {
			kinesis.MaximumRetryAttempts = pulumi.Int(int(*source.Kinesis.MaximumRetryAttempts))
		}
		if source.Kinesis.OnPartialBatchItemFailure != "" {
			kinesis.OnPartialBatchItemFailure = pulumi.String(source.Kinesis.OnPartialBatchItemFailure)
		}
		if source.Kinesis.ParallelizationFactor != nil {
			kinesis.ParallelizationFactor = pulumi.Int(int(*source.Kinesis.ParallelizationFactor))
		}
		if source.Kinesis.DeadLetterQueueArn.GetValue() != "" {
			kinesis.DeadLetterConfig = &pipes.PipeSourceParametersKinesisStreamParametersDeadLetterConfigArgs{
				Arn: pulumi.String(source.Kinesis.DeadLetterQueueArn.GetValue()),
			}
		}
		args.KinesisStreamParameters = kinesis
	}

	if source.Dynamodb != nil {
		dynamodb := &pipes.PipeSourceParametersDynamodbStreamParametersArgs{
			StartingPosition: pulumi.String(source.Dynamodb.StartingPosition),
		}
		if source.Dynamodb.BatchSize != nil {
			dynamodb.BatchSize = pulumi.Int(int(*source.Dynamodb.BatchSize))
		}
		if source.Dynamodb.MaximumBatchingWindowInSeconds != nil {
			dynamodb.MaximumBatchingWindowInSeconds = pulumi.Int(int(*source.Dynamodb.MaximumBatchingWindowInSeconds))
		}
		if source.Dynamodb.MaximumRecordAgeInSeconds != nil {
			dynamodb.MaximumRecordAgeInSeconds = pulumi.Int(int(*source.Dynamodb.MaximumRecordAgeInSeconds))
		}
		if source.Dynamodb.MaximumRetryAttempts != nil {
			dynamodb.MaximumRetryAttempts = pulumi.Int(int(*source.Dynamodb.MaximumRetryAttempts))
		}
		if source.Dynamodb.OnPartialBatchItemFailure != "" {
			dynamodb.OnPartialBatchItemFailure = pulumi.String(source.Dynamodb.OnPartialBatchItemFailure)
		}
		if source.Dynamodb.ParallelizationFactor != nil {
			dynamodb.ParallelizationFactor = pulumi.Int(int(*source.Dynamodb.ParallelizationFactor))
		}
		if source.Dynamodb.DeadLetterQueueArn.GetValue() != "" {
			dynamodb.DeadLetterConfig = &pipes.PipeSourceParametersDynamodbStreamParametersDeadLetterConfigArgs{
				Arn: pulumi.String(source.Dynamodb.DeadLetterQueueArn.GetValue()),
			}
		}
		args.DynamodbStreamParameters = dynamodb
	}

	if source.Msk != nil {
		msk := &pipes.PipeSourceParametersManagedStreamingKafkaParametersArgs{
			TopicName: pulumi.String(source.Msk.TopicName),
		}
		if source.Msk.ConsumerGroupId != "" {
			msk.ConsumerGroupId = pulumi.String(source.Msk.ConsumerGroupId)
		}
		if source.Msk.StartingPosition != "" {
			msk.StartingPosition = pulumi.String(source.Msk.StartingPosition)
		}
		if source.Msk.BatchSize != nil {
			msk.BatchSize = pulumi.Int(int(*source.Msk.BatchSize))
		}
		if source.Msk.MaximumBatchingWindowInSeconds != nil {
			msk.MaximumBatchingWindowInSeconds = pulumi.Int(int(*source.Msk.MaximumBatchingWindowInSeconds))
		}
		if source.Msk.Credentials != nil {
			credentials := &pipes.PipeSourceParametersManagedStreamingKafkaParametersCredentialsArgs{}
			if source.Msk.Credentials.ClientCertificateTlsAuth != "" {
				credentials.ClientCertificateTlsAuth = pulumi.String(source.Msk.Credentials.ClientCertificateTlsAuth)
			}
			if source.Msk.Credentials.SaslScram_512Auth != "" {
				credentials.SaslScram512Auth = pulumi.String(source.Msk.Credentials.SaslScram_512Auth)
			}
			msk.Credentials = credentials
		}
		args.ManagedStreamingKafkaParameters = msk
	}

	if source.SelfManagedKafka != nil {
		kafka := &pipes.PipeSourceParametersSelfManagedKafkaParametersArgs{
			TopicName: pulumi.String(source.SelfManagedKafka.TopicName),
		}
		if len(source.SelfManagedKafka.AdditionalBootstrapServers) > 0 {
			kafka.AdditionalBootstrapServers = pulumi.ToStringArray(source.SelfManagedKafka.AdditionalBootstrapServers)
		}
		if source.SelfManagedKafka.ConsumerGroupId != "" {
			kafka.ConsumerGroupId = pulumi.String(source.SelfManagedKafka.ConsumerGroupId)
		}
		if source.SelfManagedKafka.StartingPosition != "" {
			kafka.StartingPosition = pulumi.String(source.SelfManagedKafka.StartingPosition)
		}
		if source.SelfManagedKafka.BatchSize != nil {
			kafka.BatchSize = pulumi.Int(int(*source.SelfManagedKafka.BatchSize))
		}
		if source.SelfManagedKafka.MaximumBatchingWindowInSeconds != nil {
			kafka.MaximumBatchingWindowInSeconds = pulumi.Int(int(*source.SelfManagedKafka.MaximumBatchingWindowInSeconds))
		}
		if source.SelfManagedKafka.ServerRootCaCertificate != "" {
			kafka.ServerRootCaCertificate = pulumi.String(source.SelfManagedKafka.ServerRootCaCertificate)
		}
		if source.SelfManagedKafka.Credentials != nil {
			credentials := &pipes.PipeSourceParametersSelfManagedKafkaParametersCredentialsArgs{}
			if source.SelfManagedKafka.Credentials.BasicAuth != "" {
				credentials.BasicAuth = pulumi.String(source.SelfManagedKafka.Credentials.BasicAuth)
			}
			if source.SelfManagedKafka.Credentials.ClientCertificateTlsAuth != "" {
				credentials.ClientCertificateTlsAuth = pulumi.String(source.SelfManagedKafka.Credentials.ClientCertificateTlsAuth)
			}
			if source.SelfManagedKafka.Credentials.SaslScram_256Auth != "" {
				credentials.SaslScram256Auth = pulumi.String(source.SelfManagedKafka.Credentials.SaslScram_256Auth)
			}
			if source.SelfManagedKafka.Credentials.SaslScram_512Auth != "" {
				credentials.SaslScram512Auth = pulumi.String(source.SelfManagedKafka.Credentials.SaslScram_512Auth)
			}
			kafka.Credentials = credentials
		}
		if source.SelfManagedKafka.Vpc != nil {
			subnets := pulumi.StringArray{}
			for _, subnet := range source.SelfManagedKafka.Vpc.Subnets {
				subnets = append(subnets, pulumi.String(subnet.GetValue()))
			}
			securityGroups := pulumi.StringArray{}
			for _, securityGroup := range source.SelfManagedKafka.Vpc.SecurityGroups {
				securityGroups = append(securityGroups, pulumi.String(securityGroup.GetValue()))
			}
			vpc := &pipes.PipeSourceParametersSelfManagedKafkaParametersVpcArgs{
				Subnets: subnets,
			}
			if len(securityGroups) > 0 {
				vpc.SecurityGroups = securityGroups
			}
			kafka.Vpc = vpc
		}
		args.SelfManagedKafkaParameters = kafka
	}

	if source.Activemq != nil {
		activemq := &pipes.PipeSourceParametersActivemqBrokerParametersArgs{
			QueueName: pulumi.String(source.Activemq.QueueName),
			Credentials: &pipes.PipeSourceParametersActivemqBrokerParametersCredentialsArgs{
				BasicAuth: pulumi.String(source.Activemq.BasicAuthCredentials),
			},
		}
		if source.Activemq.BatchSize != nil {
			activemq.BatchSize = pulumi.Int(int(*source.Activemq.BatchSize))
		}
		if source.Activemq.MaximumBatchingWindowInSeconds != nil {
			activemq.MaximumBatchingWindowInSeconds = pulumi.Int(int(*source.Activemq.MaximumBatchingWindowInSeconds))
		}
		args.ActivemqBrokerParameters = activemq
	}

	if source.Rabbitmq != nil {
		rabbitmq := &pipes.PipeSourceParametersRabbitmqBrokerParametersArgs{
			QueueName: pulumi.String(source.Rabbitmq.QueueName),
			Credentials: &pipes.PipeSourceParametersRabbitmqBrokerParametersCredentialsArgs{
				BasicAuth: pulumi.String(source.Rabbitmq.BasicAuthCredentials),
			},
		}
		if source.Rabbitmq.VirtualHost != "" {
			rabbitmq.VirtualHost = pulumi.String(source.Rabbitmq.VirtualHost)
		}
		if source.Rabbitmq.BatchSize != nil {
			rabbitmq.BatchSize = pulumi.Int(int(*source.Rabbitmq.BatchSize))
		}
		if source.Rabbitmq.MaximumBatchingWindowInSeconds != nil {
			rabbitmq.MaximumBatchingWindowInSeconds = pulumi.Int(int(*source.Rabbitmq.MaximumBatchingWindowInSeconds))
		}
		args.RabbitmqBrokerParameters = rabbitmq
	}

	return args
}
