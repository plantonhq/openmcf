// Package verify checks that AWS resources created by an E2E scenario exist after
// DEPLOY and are gone after DESTROY. Each component family has its own verifier
// because AWS verification is service-specific (HeadBucket for S3,
// DescribeSubnets for a subnet, ...) -- unlike the single Management-API path a
// SaaS provider uses. All verifiers run against the same ambient credential
// chain the deploy used, so a verification failure reflects real cloud state.
package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"
)

// Verifier checks a single component's AWS resource for existence/absence.
type Verifier interface {
	// IDOutputKey is the stack-output key carrying the identifier used to verify
	// the resource (e.g. "bucket_id").
	IDOutputKey() string
	// VerifyExists returns an error unless the resource exists.
	VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error
	// VerifyAbsent returns an error unless the resource is gone.
	VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error
}

// OutputsVerifier inspects the full stack output map when a single string id is
// insufficient (e.g. AwsS3ObjectSet verifies HeadObject per key in object_etags).
type OutputsVerifier interface {
	Verifier
	VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error
	VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error
}

// verifiers maps a component name to its verifier. New AWS components register
// here as they are forged; today it carries the S3 walking-skeleton only.
var verifiers = map[string]Verifier{
	"awscertmanagercert":             &acmCertificateVerifier{},
	"awscloudfront":                  &cloudFrontDistributionVerifier{},
	"awssqsqueue":                    &sqsQueueVerifier{},
	"awssnstopic":                    &snsTopicVerifier{},
	"awssnssubscription":             &snsSubscriptionVerifier{},
	"awseventbridgebus":              &eventBridgeBusVerifier{},
	"awseventbridgerule":             &eventBridgeRuleVerifier{},
	"awss3bucket":                    &s3Verifier{},
	"awss3objectset":                 &s3ObjectSetVerifier{},
	"awssubnet":                      &subnetVerifier{},
	"awsvpc":                         &vpcVerifier{},
	"awsinternetgateway":             &internetGatewayVerifier{},
	"awsegressonlyinternetgateway":   &egressOnlyInternetGatewayVerifier{},
	"awsvpcendpoint":                 &vpcEndpointVerifier{},
	"awsnatgateway":                  &natGatewayVerifier{},
	"awselasticip":                   &elasticIpVerifier{},
	"awsiampolicy":                   &iamPolicyVerifier{},
	"awsiamrole":                     &iamRoleVerifier{},
	"awsiaminstanceprofile":          &iamInstanceProfileVerifier{},
	"awsiamuser":                     &iamUserVerifier{},
	"awsiamoidcprovider":             &iamOidcProviderVerifier{},
	"awsalb":                         &loadBalancerVerifier{component: "awsalb"},
	"awsnlb":                         &loadBalancerVerifier{component: "awsnlb"},
	"awslbtargetgroup":               &targetGroupVerifier{},
	"awslblistener":                  &listenerVerifier{},
	"awslblistenerrule":              &listenerRuleVerifier{},
	"awslaunchtemplate":              &launchTemplateVerifier{},
	"awsautoscalinggroup":            &autoScalingGroupVerifier{},
	"awsecscluster":                  &ecsClusterVerifier{},
	"awsecstaskdefinition":           &ecsTaskDefinitionVerifier{},
	"awsecsservice":                  &ecsServiceVerifier{},
	"awsrdscluster":                  &rdsClusterVerifier{},
	"awsrdsinstance":                 &rdsInstanceVerifier{},
	"awsdocumentdb":                  &docdbClusterVerifier{},
	"awsneptunecluster":              &neptuneClusterVerifier{},
	"awsredshiftcluster":             &redshiftClusterVerifier{},
	"awsredshiftserverlessnamespace": &redshiftServerlessNamespaceVerifier{},
	"awsredshiftserverlessworkgroup": &redshiftServerlessWorkgroupVerifier{},
	"awsekscluster":                  &eksClusterVerifier{},
	"awseksnodegroup":                &eksNodeGroupVerifier{},
	"awseksaddon":                    &eksAddonVerifier{},
	"awseksfargateprofile":           &eksFargateProfileVerifier{},
	"awseksaccessentry":              &eksAccessEntryVerifier{},
	"awsdynamodb":                    &dynamodbTableVerifier{},
	"awskinesisstream":               &kinesisStreamVerifier{},
	"awskinesisstreamconsumer":       &kinesisStreamConsumerVerifier{},
	"awskinesisfirehose":             &kinesisFirehoseVerifier{},
	"awslambda":                      &lambdaFunctionVerifier{},
	"awslambdaeventsourcemapping":    &lambdaEventSourceMappingVerifier{},
	"awskmskey":                      &kmsKeyVerifier{},
	"awselasticacheuser":             &elasticacheUserVerifier{},
	"awselasticacheusergroup":        &elasticacheUserGroupVerifier{},
	"awsrediselasticache":            &elasticacheReplicationGroupVerifier{},
	"awsmemcachedelasticache":        &elasticacheClusterVerifier{},
	"awsserverlesselasticache":       &elasticacheServerlessCacheVerifier{},
	"awssecuritygroup":               &securityGroupVerifier{},
	"awsmskcluster":                  &mskClusterVerifier{},
	"awsmskserverlesscluster":        &mskServerlessClusterVerifier{},
	"awsmwaaenvironment":             &mwaaEnvironmentVerifier{},
	"awsopensearchdomain":            &opensearchDomainVerifier{},
	"awsec2instance":                 &ec2InstanceVerifier{},
	"awsecrrepo":                     &ecrRepoVerifier{},
	"awsroute53zone":                 &route53ZoneVerifier{},
	"awsroute53dnsrecord":            &route53DnsRecordVerifier{},
	"awsroute53healthcheck":          &route53HealthCheckVerifier{},
	"awscloudwatchloggroup":          &cloudwatchLogGroupVerifier{},
	"awscloudwatchalarm":             &cloudwatchAlarmVerifier{},
	"awscloudwatchcompositealarm":    &cloudwatchCompositeAlarmVerifier{},
	"awsstepfunction":                &stepFunctionVerifier{},
	"awshttpapigateway":              &httpApiGatewayVerifier{},
	"awshttpapivpclink":              &httpApiVpcLinkVerifier{},
	"awshttpapidomain":               &httpApiDomainVerifier{},
	"awscognitouserpool":             &cognitoUserPoolVerifier{},
	"awscognitouserpoolclient":       &cognitoUserPoolClientVerifier{},
	"awscognitoidentityprovider":     &cognitoIdentityProviderVerifier{},
	"awscognitoresourceserver":       &cognitoResourceServerVerifier{},
	"awselasticfilesystem":           &efsFileSystemVerifier{},
	"awsefsaccesspoint":              &efsAccessPointVerifier{},
	"awswafwebacl":                   &wafWebAclVerifier{},
	"awswafipset":                    &wafIpSetVerifier{},
	"awswafregexpatternset":          &wafRegexPatternSetVerifier{},
	"awsbatchcomputeenvironment":     &batchComputeEnvironmentVerifier{},
	"awsbatchjobqueue":               &batchJobQueueVerifier{},
	"awsbatchschedulingpolicy":       &batchSchedulingPolicyVerifier{},
	"awsbatchjobdefinition":          &batchJobDefinitionVerifier{},
}

// GetVerifier returns the verifier for a component, or an error if none is registered.
func GetVerifier(component string) (Verifier, error) {
	v, ok := verifiers[component]
	if !ok {
		return nil, errors.Errorf("no AWS verifier registered for component %q", component)
	}
	return v, nil
}
