//go:build e2e

// Package aws contains end-to-end tests that provision real AWS resources via
// Planton IaC modules and verify them through the AWS SDK. Credentials come from
// the ambient chain (local AWS SSO or GitHub Actions OIDC -- never a stored
// secret); see the aa_e2e harness package.
//
// Run with: go test -tags=e2e -timeout=30m -v ./e2e/aws/...
package aws

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	awse2e "github.com/plantonhq/planton/apis/dev/planton/provider/aws/aa_e2e"
	"github.com/plantonhq/planton/e2e/framework/discovery"
	"github.com/plantonhq/planton/e2e/framework/provider"
	"github.com/plantonhq/planton/e2e/framework/runner"
)

var (
	testHarness      *awse2e.Harness
	repoRoot         string
	runID            string
	pulumiBackendURL string
)

func TestMain(m *testing.M) {
	var err error
	repoRoot, err = filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve repo root: %v\n", err)
		os.Exit(1)
	}

	runID = uuid.New().String()[:8]

	backendDir, err := os.MkdirTemp("", "planton-e2e-aws-pulumi-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp backend dir: %v\n", err)
		os.Exit(1)
	}
	pulumiBackendURL = "file://" + backendDir
	defer os.RemoveAll(backendDir)

	if err := runner.PulumiLogin(pulumiBackendURL); err != nil {
		fmt.Fprintf(os.Stderr, "failed to login to pulumi backend: %v\n", err)
		os.Exit(1)
	}

	testHarness = awse2e.NewHarness()
	ctx := context.Background()
	if err := testHarness.Setup(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup AWS harness: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	if err := testHarness.Teardown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to teardown AWS harness: %v\n", err)
	}

	os.Exit(code)
}

// --- AWS S3 Bucket (walking skeleton) ---

func TestAwsS3Bucket_Pulumi(t *testing.T) { runAllScenariosForComponent(t, "awss3bucket", "pulumi") }
func TestAwsS3Bucket_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awss3bucket", "terraform")
}

// --- AWS ECR Repository (true leaf; repository + lifecycle policy in one fast lane) ---

func TestAwsEcrRepo_Pulumi(t *testing.T) { runAllScenariosForComponent(t, "awsecrrepo", "pulumi") }
func TestAwsEcrRepo_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsecrrepo", "terraform")
}

// --- AWS VPC (thin root of the networking graph) ---

func TestAwsVpc_Pulumi(t *testing.T)    { runAllScenariosForComponent(t, "awsvpc", "pulumi") }
func TestAwsVpc_Terraform(t *testing.T) { runAllScenariosForComponent(t, "awsvpc", "terraform") }

// --- AWS Subnet (first composed topology: deploys an AwsVpc prerequisite) ---

func TestAwsSubnet_Pulumi(t *testing.T)    { runAllScenariosForComponent(t, "awssubnet", "pulumi") }
func TestAwsSubnet_Terraform(t *testing.T) { runAllScenariosForComponent(t, "awssubnet", "terraform") }

// --- AWS NAT Gateway (deep composed topology: AwsVpc -> AwsSubnet -> AwsElasticIp) ---

func TestAwsNatGateway_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsnatgateway", "pulumi")
}
func TestAwsNatGateway_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsnatgateway", "terraform")
}

// --- AWS Internet Gateway (attaches to a gateway-free AwsVpc prerequisite) ---

func TestAwsInternetGateway_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsinternetgateway", "pulumi")
}
func TestAwsInternetGateway_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsinternetgateway", "terraform")
}

// --- AWS Egress-Only Internet Gateway (IPv6 outbound-only; attaches to an AwsVpc prerequisite) ---

func TestAwsEgressOnlyInternetGateway_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsegressonlyinternetgateway", "pulumi")
}
func TestAwsEgressOnlyInternetGateway_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsegressonlyinternetgateway", "terraform")
}

// --- AWS VPC Endpoint (private service access: gateway on the VPC's default route table + interface on the subnet pair) ---

func TestAwsVpcEndpoint_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsvpcendpoint", "pulumi")
}
func TestAwsVpcEndpoint_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsvpcendpoint", "terraform")
}

// --- AWS IAM Policy (leaf of the identity graph: a standalone managed policy) ---

func TestAwsIamPolicy_Pulumi(t *testing.T) { runAllScenariosForComponent(t, "awsiampolicy", "pulumi") }
func TestAwsIamPolicy_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsiampolicy", "terraform")
}

// --- AWS IAM Role (trust policy + attachments + inline policies) ---

func TestAwsIamRole_Pulumi(t *testing.T) { runAllScenariosForComponent(t, "awsiamrole", "pulumi") }
func TestAwsIamRole_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsiamrole", "terraform")
}

// --- AWS IAM Instance Profile (composed identity chain: deploys an AwsIamRole prerequisite) ---

func TestAwsIamInstanceProfile_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsiaminstanceprofile", "pulumi")
}
func TestAwsIamInstanceProfile_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsiaminstanceprofile", "terraform")
}

// --- AWS IAM User (long-lived identity with access key + force-destroy teardown) ---

func TestAwsIamUser_Pulumi(t *testing.T) { runAllScenariosForComponent(t, "awsiamuser", "pulumi") }
func TestAwsIamUser_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsiamuser", "terraform")
}

// --- AWS IAM OIDC Provider (keyless-federation trust anchor; synthetic issuer) ---

func TestAwsIamOidcProvider_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsiamoidcprovider", "pulumi")
}
func TestAwsIamOidcProvider_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsiamoidcprovider", "terraform")
}

// --- AWS LB Target Group (routing destination; deploys an AwsVpc prerequisite) ---

func TestAwsLbTargetGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awslbtargetgroup", "pulumi")
}
func TestAwsLbTargetGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awslbtargetgroup", "terraform")
}

// --- AWS ALB (two-AZ placement: deploys the VPC + two-subnet prerequisite pair) ---

func TestAwsAlb_Pulumi(t *testing.T)    { runAllScenariosForComponent(t, "awsalb", "pulumi") }
func TestAwsAlb_Terraform(t *testing.T) { runAllScenariosForComponent(t, "awsalb", "terraform") }

// --- AWS NLB (internal single-mapping smoke; deploys the VPC + subnet prerequisites) ---

func TestAwsNlb_Pulumi(t *testing.T)    { runAllScenariosForComponent(t, "awsnlb", "pulumi") }
func TestAwsNlb_Terraform(t *testing.T) { runAllScenariosForComponent(t, "awsnlb", "terraform") }

// --- AWS LB Listener (deepest chain: VPC -> subnets -> ALB -> target group -> listener) ---

func TestAwsLbListener_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awslblistener", "pulumi")
}
func TestAwsLbListener_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awslblistener", "terraform")
}

// --- AWS LB Listener Rule (per-service routing; rides the full family chain) ---

func TestAwsLbListenerRule_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awslblistenerrule", "pulumi")
}
func TestAwsLbListenerRule_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awslblistenerrule", "terraform")
}

func TestAwsLaunchTemplate_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awslaunchtemplate", "pulumi")
}
func TestAwsLaunchTemplate_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awslaunchtemplate", "terraform")
}

func TestAwsAutoScalingGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsautoscalinggroup", "pulumi")
}
func TestAwsAutoScalingGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsautoscalinggroup", "terraform")
}

// --- AWS RDS Cluster (Aurora Serverless v2 with one folded db.serverless instance; managed master password) ---

func TestAwsRdsCluster_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsrdscluster", "pulumi")
}
func TestAwsRdsCluster_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsrdscluster", "terraform")
}

// --- AWS RDS Instance (single-AZ postgres db.t4g.micro on gp3; managed master password) ---

func TestAwsRdsInstance_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsrdsinstance", "pulumi")
}
func TestAwsRdsInstance_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsrdsinstance", "terraform")
}

// --- AWS DocumentDB (one db.t4g.medium folded instance; managed master password) ---

func TestAwsDocumentDb_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsdocumentdb", "pulumi")
}
func TestAwsDocumentDb_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsdocumentdb", "terraform")
}

// --- AWS Neptune Cluster (Serverless 1-2 NCU with one folded db.serverless instance) ---

func TestAwsNeptuneCluster_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsneptunecluster", "pulumi")
}
func TestAwsNeptuneCluster_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsneptunecluster", "terraform")
}

// --- AWS Redshift Cluster (single ra3.large node; managed admin password) ---

func TestAwsRedshiftCluster_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsredshiftcluster", "pulumi")
}
func TestAwsRedshiftCluster_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsredshiftcluster", "terraform")
}

// --- AWS Redshift Serverless Namespace (account-level data plane; managed admin password) ---

func TestAwsRedshiftServerlessNamespace_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsredshiftserverlessnamespace", "pulumi")
}
func TestAwsRedshiftServerlessNamespace_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsredshiftserverlessnamespace", "terraform")
}

// --- AWS Redshift Serverless Workgroup (8-RPU compute plane on the namespace + three-AZ subnet trio) ---

func TestAwsRedshiftServerlessWorkgroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsredshiftserverlessworkgroup", "pulumi")
}
func TestAwsRedshiftServerlessWorkgroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsredshiftserverlessworkgroup", "terraform")
}

// --- AWS DynamoDB (true leaf; on-demand table with GSI, streams, PITR, folded satellites) ---

func TestAwsDynamodb_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsdynamodb", "pulumi")
}
func TestAwsDynamodb_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsdynamodb", "terraform")
}

// --- AWS Kinesis Data Stream (true leaf; ON_DEMAND full surface with folded resource policy) ---

func TestAwsKinesisStream_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awskinesisstream", "pulumi")
}
func TestAwsKinesisStream_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awskinesisstream", "terraform")
}

// --- AWS Kinesis Stream Consumer (enhanced fan-out; AwsKinesisStream prerequisite chain) ---

func TestAwsKinesisStreamConsumer_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awskinesisstreamconsumer", "pulumi")
}
func TestAwsKinesisStreamConsumer_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awskinesisstreamconsumer", "terraform")
}

// --- AWS Kinesis Firehose (Direct PUT -> extended_s3; S3 bucket + IAM role prerequisite chain) ---

func TestAwsKinesisFirehose_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awskinesisfirehose", "pulumi")
}
func TestAwsKinesisFirehose_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awskinesisfirehose", "terraform")
}

// --- AWS ElastiCache RBAC (account-level; no VPC prerequisite) ---

func TestAwsElasticacheUser_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awselasticacheuser", "pulumi")
}
func TestAwsElasticacheUser_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awselasticacheuser", "terraform")
}

func TestAwsElasticacheUserGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awselasticacheusergroup", "pulumi")
}
func TestAwsElasticacheUserGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awselasticacheusergroup", "terraform")
}

// --- AWS ElastiCache caches (AwsSubnet two-AZ prerequisite chain) ---

func TestAwsRedisElasticache_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsrediselasticache", "pulumi")
}
func TestAwsRedisElasticache_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsrediselasticache", "terraform")
}

func TestAwsMemcachedElasticache_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsmemcachedelasticache", "pulumi")
}
func TestAwsMemcachedElasticache_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsmemcachedelasticache", "terraform")
}

func TestAwsServerlessElasticache_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsserverlesselasticache", "pulumi")
}
func TestAwsServerlessElasticache_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsserverlesselasticache", "terraform")
}

// --- AWS MSK Cluster (two kafka.t3.small brokers, SASL/IAM; subnet + security-group prerequisite chain) ---

func TestAwsMskCluster_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsmskcluster", "pulumi")
}
func TestAwsMskCluster_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsmskcluster", "terraform")
}

// --- AWS MSK Serverless Cluster (capacity-managed Kafka, SASL/IAM only; subnet prerequisite chain) ---

func TestAwsMskServerlessCluster_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsmskserverlesscluster", "pulumi")
}
func TestAwsMskServerlessCluster_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsmskserverlesscluster", "terraform")
}

// --- AWS MWAA Environment (managed Airflow; S3 + IAM-role fixtures via the scenario's e2e-prerequisites annotation) ---

func TestAwsMwaaEnvironment_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsmwaaenvironment", "pulumi")
}
func TestAwsMwaaEnvironment_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsmwaaenvironment", "terraform")
}

// --- AWS OpenSearch Domain (true leaf; public single-node t3.small.search domain) ---

func TestAwsOpenSearchDomain_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsopensearchdomain", "pulumi")
}
func TestAwsOpenSearchDomain_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsopensearchdomain", "terraform")
}

// --- AWS EKS Cluster (the control plane; slowest single resource in the suite, ~25 min/engine) ---

func TestAwsEksCluster_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsekscluster", "pulumi")
}
func TestAwsEksCluster_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsekscluster", "terraform")
}

// --- AWS EKS Node Group (zero-capacity pool on a full control-plane prerequisite chain) ---

func TestAwsEksNodeGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awseksnodegroup", "pulumi")
}
func TestAwsEksNodeGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awseksnodegroup", "terraform")
}

// --- AWS EKS Add-on (kube-proxy adopted with OVERWRITE on the control-plane prerequisite chain) ---

func TestAwsEksAddon_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awseksaddon", "pulumi")
}
func TestAwsEksAddon_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awseksaddon", "terraform")
}

// --- AWS EKS Fargate Profile (namespace selector; no pods scheduled; full control-plane prerequisite chain) ---

func TestAwsEksFargateProfile_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awseksfargateprofile", "pulumi")
}
func TestAwsEksFargateProfile_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awseksfargateprofile", "terraform")
}

// --- AWS EKS Access Entry (STANDARD entry + view-policy association on an API-mode cluster) ---

func TestAwsEksAccessEntry_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awseksaccessentry", "pulumi")
}
func TestAwsEksAccessEntry_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awseksaccessentry", "terraform")
}

// --- AWS ECS Task Definition (revision registration only; no task launches; execution-role prerequisite) ---

func TestAwsEcsTaskDefinition_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsecstaskdefinition", "pulumi")
}
func TestAwsEcsTaskDefinition_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsecstaskdefinition", "terraform")
}

// --- AWS ECS Service (desiredCount 0 on the cluster + task-definition + subnet prerequisite chain) ---

func TestAwsEcsService_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsecsservice", "pulumi")
}
func TestAwsEcsService_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsecsservice", "terraform")
}

func TestAwsPlantonRunner_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsplantonrunner", "pulumi")
}

func TestAwsPlantonRunner_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsplantonrunner", "terraform")
}

// --- AWS ECS Cluster (Fargate leaf + the EC2-capacity chain via the scenario's e2e-prerequisites annotation) ---

func TestAwsEcsCluster_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsecscluster", "pulumi")
}
func TestAwsEcsCluster_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsecscluster", "terraform")
}

// --- AWS EC2 Instance (one t3.micro; subnet + security-group fixtures via the scenario's e2e-prerequisites annotation) ---

func TestAwsEc2Instance_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awsec2instance", "pulumi")
}
func TestAwsEc2Instance_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awsec2instance", "terraform")
}

// --- AWS Security Group (rules-rich group on the VPC prerequisite; prefix lists, self-reference, SG-to-SG refs) ---

func TestAwsSecurityGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awssecuritygroup", "pulumi")
}
func TestAwsSecurityGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awssecuritygroup", "terraform")
}

// --- AWS CloudFront (placeholder-origin distribution on the default certificate; true leaf, slow edge propagation) ---

func TestAwsCloudFront_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awscloudfront", "pulumi")
}
func TestAwsCloudFront_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awscloudfront", "terraform")
}

// --- AWS Certificate Manager (no-zone requested certificate resting in PENDING_VALIDATION; true leaf) ---

func TestAwsCertManagerCert_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awscertmanagercert", "pulumi")
}
func TestAwsCertManagerCert_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awscertmanagercert", "terraform")
}

// --- AWS KMS Key (symmetric key with one alias; true leaf) ---

func TestAwsKmsKey_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awskmskey", "pulumi")
}
func TestAwsKmsKey_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awskmskey", "terraform")
}

// --- AWS Lambda (zip-backed function; S3 bucket + object-set + execution-role chain) ---

func TestAwsLambda_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awslambda", "pulumi")
}
func TestAwsLambda_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awslambda", "terraform")
}

// --- AWS Lambda Event Source Mapping (SQS queue -> mapping -> function chain) ---

func TestAwsLambdaEventSourceMapping_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awslambdaeventsourcemapping", "pulumi")
}
func TestAwsLambdaEventSourceMapping_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awslambdaeventsourcemapping", "terraform")
}

// --- AWS SQS Queue (FIFO queue, full FIFO delivery surface) ---

func TestAwsSqsQueue_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awssqsqueue", "pulumi")
}
func TestAwsSqsQueue_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awssqsqueue", "terraform")
}

// --- AWS SNS Topic (standard topic with tracing + policy struct) ---

func TestAwsSnsTopic_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awssnstopic", "pulumi")
}
func TestAwsSnsTopic_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awssnstopic", "terraform")
}

// --- AWS SNS Subscription (topic registry prerequisite + composed SQS queue) ---

func TestAwsSnsSubscription_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awssnssubscription", "pulumi")
}
func TestAwsSnsSubscription_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awssnssubscription", "terraform")
}

// --- AWS EventBridge Bus (custom bus with log_config + composed DLQ) ---

func TestAwsEventBridgeBus_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awseventbridgebus", "pulumi")
}
func TestAwsEventBridgeBus_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awseventbridgebus", "terraform")
}

// --- AWS EventBridge Rule (event-pattern rule on custom bus + SQS target) ---

func TestAwsEventBridgeRule_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "awseventbridgerule", "pulumi")
}
func TestAwsEventBridgeRule_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "awseventbridgerule", "terraform")
}

// runAllScenariosForComponent discovers and runs all E2E scenarios for an AWS component.
func runAllScenariosForComponent(t *testing.T, component, engine string) {
	t.Helper()

	var moduleDir string
	switch engine {
	case "pulumi":
		moduleDir = filepath.Join(repoRoot, "apis", "dev", "planton", "provider", "aws", component, "v1", "iac", "pulumi")
	case "terraform":
		moduleDir = filepath.Join(repoRoot, "apis", "dev", "planton", "provider", "aws", component, "v1", "iac", "tf")
	default:
		t.Fatalf("unsupported engine: %s", engine)
	}

	if !fileExists(moduleDir) {
		t.Skipf("component %s %s module not found at %s", component, engine, moduleDir)
	}

	scenarios, err := discovery.DiscoverTestScenarios(repoRoot, "aws", component)
	if err != nil {
		t.Fatalf("failed to discover test scenarios for %s: %v", component, err)
	}

	if len(scenarios) == 0 {
		t.Skipf("no test scenarios found for %s", component)
	}

	t.Logf("Discovered %d scenarios for %s [%s]", len(scenarios), component, engine)

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.Name, func(t *testing.T) {
			runSingleScenario(t, component, moduleDir, engine, scenario)
		})
	}
}

func runSingleScenario(t *testing.T, component, moduleDir, engine string, scenario discovery.TestScenario) {
	t.Helper()

	tc := &provider.ComponentTestContext{
		Component:    component,
		Provider:     "aws",
		Engine:       engine,
		ModuleDir:    moduleDir,
		ManifestPath: scenario.ManifestPath,
		RepoRoot:     repoRoot,
		RunID:        runID,
		T:            t,
	}

	if engine == "pulumi" {
		stackName := runner.GenerateStackName(component+"-"+scenario.Name, runID)
		if len(stackName) > 50 {
			stackName = stackName[:50]
		}
		tc.StackName = stackName
		tc.BackendURL = pulumiBackendURL
	}

	ctx := context.Background()
	result := runner.RunComponentTest(ctx, tc, testHarness)

	for _, phase := range result.Phases {
		status := "PASS"
		if !phase.Passed {
			status = "FAIL"
		}
		t.Logf("  %s: %s (%s)", phase.Phase, status, phase.Duration)
		if phase.Error != nil {
			t.Logf("    Error: %v", phase.Error)
		}
	}

	if !result.Passed {
		t.Fatalf("scenario %s/%s [%s] failed (total: %s)", component, scenario.Name, engine, result.Duration)
	}

	t.Logf("scenario %s/%s [%s] passed (total: %s)", component, scenario.Name, engine, result.Duration)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
