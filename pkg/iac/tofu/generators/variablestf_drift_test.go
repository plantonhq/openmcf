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
	// ElastiCache family (Session 013): RBAC kinds + three cache kinds, migrated
	// off the legacy type = any contracts together.
	"AwsElasticacheUser",
	"AwsElasticacheUserGroup",
	"AwsRedisElasticache",
	"AwsMemcachedElasticache",
	"AwsServerlessElasticache",
	// Aurora-shaped siblings (Session 014), migrated off the legacy
	// hand-written contracts together.
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
// descriptor's source file: dev/planton/provider/<p>/<kind>/v1/api.proto ->
// <repo>/apis/dev/planton/provider/<p>/<kind>/v1/iac/tf/variables.tf.
func moduleVariablesPath(root string, msg proto.Message) string {
	protoPath := msg.ProtoReflect().Descriptor().ParentFile().Path()
	dir := filepath.Dir(protoPath)
	return filepath.Join(root, "apis", dir, "iac", "tf", "variables.tf")
}

// repoRoot walks up from this test file to the directory containing go.mod.
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
			t.Fatal("could not locate repo root (go.mod)")
		}
		dir = parent
	}
}
