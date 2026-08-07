package generators

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/plantonhq/planton/pkg/crkreflect"
	"google.golang.org/protobuf/proto"
)

// migratedKinds is the allowlist of cloud-resource kinds whose committed
// variables.tf is owned by the generator (ProtoToVariablesTF) and guarded
// against drift. A kind is added here only after its module has been regenerated
// and validated (tofu validate against a null-pruned tfvars), so the guard can
// never be red for an unmigrated module. Remaining providers/kinds are migrated
// in tracked batches, each appended here.
//
// Scope note: only providers whose module conventions match the generator (it
// flattens wrapper types like StringValueOrRef to primitives and emits the
// canonical metadata block) belong here. AWS modules follow these conventions.
// Providers that intentionally diverge (e.g. OCI modules expose the wrapper
// object) are out of scope until their modules are migrated to the generator.
var migratedKinds = []string{
	// aws-ecs-environment chart kinds (the set that surfaced the schema bug).
	"AwsRoute53Zone",
	"AwsEcsCluster",
	"AwsIamRole",
	"AwsEcrRepo",
	"AwsSecurityGroup",
	"AwsAlb",
	"AwsCertManagerCert",
	"AwsEcsService",
	"AwsEcsTaskDefinition",
	// AWS networking primitives already on the modern schema, brought under the
	// guard so they cannot regress.
	"AwsVpc",
	"AwsSubnet",
	"AwsInternetGateway",
	"AwsNatGateway",
	"AwsElasticIp",
	// EC2 fleet compute kinds, generator-owned from the day they were forged.
	"AwsLaunchTemplate",
	"AwsAutoScalingGroup",
	// EKS control plane + managed node group, migrated off the legacy
	// hand-written contract together.
	"AwsEksCluster",
	"AwsEksNodeGroup",
	// EKS satellites, generator-owned from the day they were forged.
	"AwsEksAddon",
	"AwsEksFargateProfile",
	"AwsEksAccessEntry",
	// Networking fast-follow, generator-owned from the day it was forged.
	"AwsVpcEndpoint",
	// RDS pair, migrated off the legacy hand-written contracts together.
	"AwsRdsCluster",
	"AwsRdsInstance",
	// ElastiCache family: RBAC kinds + three cache kinds, migrated off the
	// legacy type = any contracts together.
	"AwsElasticacheUser",
	"AwsElasticacheUserGroup",
	"AwsRedisElasticache",
	"AwsMemcachedElasticache",
	"AwsServerlessElasticache",
	// Aurora-shaped siblings, migrated off the legacy hand-written
	// contracts together.
	"AwsDocumentDb",
	"AwsNeptuneCluster",
	// Redshift, migrated off the legacy hand-written contract.
	"AwsRedshiftCluster",
	// Redshift Serverless pair, generator-owned from the day it was forged.
	"AwsRedshiftServerlessNamespace",
	"AwsRedshiftServerlessWorkgroup",
	// DynamoDB, migrated off the legacy hand-written contract.
	"AwsDynamodb",
	// Streaming + search depth pass: MSK migrated off its legacy hand-written
	// contract, OpenSearch off its legacy type = any contract.
	"AwsMskCluster",
	"AwsOpenSearchDomain",
	// EC2 instance, migrated off its legacy hand-written contract.
	"AwsEc2Instance",
	// MWAA, migrated off its legacy hand-written contract.
	"AwsMwaaEnvironment",
	// MSK Serverless, generator-owned from the day it was forged.
	"AwsMskServerlessCluster",
	// Lambda + KMS depth pass: both migrated off legacy hand-written
	// contracts; the event source mapping generator-owned from the day
	// it was forged.
	"AwsLambda",
	"AwsKmsKey",
	"AwsLambdaEventSourceMapping",
	// Edge pair: CloudFront migrated off its legacy hand-written contract
	// (ACM was already enrolled and regenerated with its rebuilt spec).
	"AwsCloudFront",
	// Messaging + eventing depth pass: SQS, SNS topic/subscription, and
	// EventBridge bus/rule, generator-owned from the day they were forged.
	"AwsSqsQueue",
	"AwsSnsTopic",
	"AwsSnsSubscription",
	"AwsEventBridgeBus",
	"AwsEventBridgeRule",
	// Object storage + streaming depth pass: S3 migrated off its legacy
	// hand-written contract; the Kinesis family off its legacy type = any
	// contracts.
	"AwsS3Bucket",
	"AwsKinesisStream",
	"AwsKinesisStreamConsumer",
	"AwsKinesisFirehose",
	// DNS: the record migrated off its hand-written contract (the zone was
	// already enrolled); the health check generator-owned from the day it
	// was forged.
	"AwsRoute53DnsRecord",
	"AwsRoute53HealthCheck",
	// Observability: the CloudWatch pair migrated off their legacy type = any
	// contracts; the composite alarm generator-owned from the day it was
	// forged.
	"AwsCloudwatchLogGroup",
	"AwsCloudwatchAlarm",
	"AwsCloudwatchCompositeAlarm",
	// Serverless front door: Step Functions migrated off its legacy type = any
	// contract, the HTTP API off its hand-written typed contract; the VPC link
	// and custom domain generator-owned from the day they were forged.
	"AwsStepFunction",
	"AwsHttpApiGateway",
	"AwsHttpApiVpcLink",
	"AwsHttpApiDomain",
	// Cognito family: the user pool migrated off its region-only typed
	// skeleton, the identity provider off its legacy type = any contract; the
	// app client and resource server generator-owned from the day they were
	// forged.
	"AwsCognitoUserPool",
	"AwsCognitoIdentityProvider",
	"AwsCognitoUserPoolClient",
	"AwsCognitoResourceServer",
	// EFS family: the file system migrated off its legacy hand-written flat
	// contract; the access point generator-owned from the day it was forged.
	"AwsElasticFileSystem",
	"AwsEfsAccessPoint",
	// WAF family: the web ACL migrated off its legacy type = any contract;
	// the IP set and regex pattern set generator-owned from the day they
	// were forged.
	"AwsWafWebAcl",
	"AwsWafIpSet",
	"AwsWafRegexPatternSet",
	// Batch family: the compute environment migrated off its legacy
	// hand-written contract; the job queue, scheduling policy, and job
	// definition generator-owned from the day they were forged.
	"AwsBatchComputeEnvironment",
	"AwsBatchJobQueue",
	"AwsBatchSchedulingPolicy",
	"AwsBatchJobDefinition",
	// SES family, generator-owned from the day it was forged.
	"AwsSesConfigurationSet",
	"AwsSesEmailIdentity",
	// App Runner family: the service migrated off its legacy hand-written
	// contract; the three companion kinds generator-owned from the day they
	// were forged.
	"AwsAppRunnerService",
	"AwsAppRunnerAutoScalingConfiguration",
	"AwsAppRunnerVpcConnector",
	"AwsAppRunnerObservabilityConfiguration",
	// Transit Gateway family: the gateway migrated off its hand-written
	// contract; the VPC attachment and route table generator-owned from the
	// day they were forged.
	"AwsTransitGateway",
	"AwsTransitGatewayVpcAttachment",
	"AwsTransitGatewayRouteTable",
	// Analytics pair, migrated off their hand-written contracts together.
	"AwsAthenaWorkgroup",
	"AwsGlueCatalogDatabase",
	// CI/CD pair, migrated off their hand-written contracts together.
	"AwsCodeBuildProject",
	"AwsCodePipeline",
	// MemoryDB family: the cluster migrated off its hand-written contract;
	// the user and ACL generator-owned from the day they were forged.
	"AwsMemorydbCluster",
	"AwsMemorydbUser",
	"AwsMemorydbAcl",
	// Client VPN, migrated off its legacy hand-written contract.
	"AwsClientVpn",
	// Global Accelerator, migrated off its hand-written contract.
	"AwsGlobalAccelerator",
	// FSx standalone trio, migrated off their hand-written contracts
	// together; the data repository association generator-owned from the
	// day it was forged.
	"AwsFsxLustreFileSystem",
	"AwsFsxOpenzfsFileSystem",
	"AwsFsxWindowsFileSystem",
	"AwsFsxDataRepositoryAssociation",
	// SageMaker Domain, migrated off its legacy hand-written contract.
	"AwsSagemakerDomain",
	// FSx ONTAP trio, migrated off their hand-written contracts together.
	"AwsFsxOntapFileSystem",
	"AwsFsxOntapStorageVirtualMachine",
	"AwsFsxOntapVolume",
	// S3 object set, migrated off its legacy hand-written contract.
	"AwsS3ObjectSet",
}

// TestVariablesTFDrift asserts that every migrated module's committed
// variables.tf is byte-identical to the generator output. This makes the
// generator the single source of truth: a hand-edit or a legacy schema can never
// silently ship. Run with PLANTON_REGEN_VARIABLES=1 to (re)write the files from
// the generator instead of comparing.
func TestVariablesTFDrift(t *testing.T) {
	root := repoRoot(t)
	regenerate := os.Getenv("PLANTON_REGEN_VARIABLES") == "1"

	for _, kindName := range migratedKinds {
		kindName := kindName
		t.Run(kindName, func(t *testing.T) {
			kind := crkreflect.KindFromString(kindName)
			msg, err := crkreflect.NewInstance(kind)
			if err != nil {
				t.Fatalf("NewInstance(%s): %v", kindName, err)
			}

			want, err := ProtoToVariablesTF(msg)
			if err != nil {
				t.Fatalf("ProtoToVariablesTF(%s): %v", kindName, err)
			}
			want = strings.TrimRight(want, "\n") + "\n"

			path := moduleVariablesPath(root, msg)

			if regenerate {
				if err := os.WriteFile(path, []byte(want), 0o644); err != nil {
					t.Fatalf("write %s: %v", path, err)
				}
				t.Logf("regenerated %s", path)
				return
			}

			gotBytes, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s (did you run PLANTON_REGEN_VARIABLES=1?): %v", path, err)
			}
			if strings.TrimRight(string(gotBytes), "\n")+"\n" != want {
				t.Errorf("variables.tf for %s is out of sync with the generator.\n"+
					"Run: PLANTON_REGEN_VARIABLES=1 go test ./pkg/iac/tofu/generators/ -run TestVariablesTFDrift\n"+
					"path: %s", kindName, path)
			}
		})
	}
}

// moduleVariablesPath derives a kind's module variables.tf path from its proto
// descriptor's source file: catalog/<p>/<kind>/<version>/api.proto ->
// <repo>/catalog/<p>/<kind>/iac/tf/variables.tf (the module lives at the
// component root, one level above the versioned contract).
func moduleVariablesPath(root string, msg proto.Message) string {
	protoPath := msg.ProtoReflect().Descriptor().ParentFile().Path()
	componentDir := filepath.Dir(filepath.Dir(protoPath))
	return filepath.Join(root, componentDir, "iac", "tf", "variables.tf")
}

// repoRoot walks up from this test file to the directory containing go.mod.
// Under the Bazel sandbox the repo checkout (and its committed variables.tf
// files) is not present, so the drift guard cannot run there -- it is
// enforced by the plain `go test` lane, and skips explicitly under Bazel
// instead of failing on the unreachable go.mod.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			if os.Getenv("TEST_WORKSPACE") != "" {
				t.Skip("skipping drift guard under the Bazel sandbox: the repo checkout (go.mod + committed variables.tf) is not available; the guard runs via `go test`")
			}
			t.Fatal("could not locate repo root (go.mod)")
		}
		dir = parent
	}
}
