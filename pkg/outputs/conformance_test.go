//go:build !codegen
// +build !codegen

package outputs

import (
	"path/filepath"
	"testing"

	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
)

// TestStackOutputsConformance is the standing guard against the systemic IaC
// output-drift class: an engine emits output names/shapes that do not flatten
// onto the kind's StackOutputs proto, silently leaving those proto fields empty.
// (The original bug: the Postgres tofu module emitted a flat
// "password_secret_name" output, which flattens to the key "password_secret_name"
// -- with no dot -- and therefore never populated the proto's nested
// password_secret{name,key} field, while the Pulumi module emitted the correct
// "password_secret.name". See the planton-postgres-iac-parity work.)
//
// Why this also enforces tofu<->pulumi parity: both engines feed the SAME generic
// transformer (TransformRaw -> Flatten -> populateMessage). So a single
// conformance bar per kind -- "this representative output set fully populates the
// proto with nothing left unmapped" -- when satisfied by each engine's emitted
// output set, guarantees the two engines produce the same typed StackOutputs.
//
// To extend coverage: add a case with the raw output shape an engine emits (scalars
// as strings; nested objects as map[string]interface{}, exactly how Terraform state
// and the Pulumi automation API surface them) and the proto fields it must populate.
func TestStackOutputsConformance(t *testing.T) {
	// A module dir with no transform override forces the generic reflection path,
	// which is the convention every in-repo module relies on (0 of 364 use an override).
	genericModuleDir := filepath.Join("testdata", "modules", "empty")

	cases := []struct {
		name string
		kind cloudresourcekind.CloudResourceKind
		// rawOutputs mirrors the post-Flatten-input shape both engines emit.
		rawOutputs map[string]interface{}
		// mustPopulate lists StackOutputs proto fields that MUST be set.
		mustPopulate []string
	}{
		{
			name: "KubernetesPostgres",
			kind: cloudresourcekind.CloudResourceKind_KubernetesPostgres,
			rawOutputs: map[string]interface{}{
				"namespace":            "gosilver-prod",
				"service":              "gosilver-prod-postgres-master",
				"port_forward_command": "kubectl port-forward -n gosilver-prod service/gosilver-prod-postgres-master 8080:8080",
				"kube_endpoint":        "gosilver-prod-postgres-master.gosilver-prod.svc.cluster.local",
				"external_hostname":    "gosilver-prod-postgres.planton.live",
				// Nested objects -- the shape that flattens to password_secret.name etc.
				"password_secret": map[string]interface{}{
					"name": "postgres.db-gosilver-prod-postgres.credentials.postgresql.acid.zalan.do",
					"key":  "password",
				},
				"username_secret": map[string]interface{}{
					"name": "postgres.db-gosilver-prod-postgres.credentials.postgresql.acid.zalan.do",
					"key":  "username",
				},
			},
			mustPopulate: []string{
				"namespace", "service", "port_forward_command", "kube_endpoint",
				"external_hostname", "password_secret", "username_secret",
			},
		},
		{
			// AwsSubnet: flat scalar outputs from both engines (subnet id/arn, AZ,
			// CIDR, route table id, region) must each land on the StackOutputs proto.
			name: "AwsSubnet",
			kind: cloudresourcekind.CloudResourceKind_AwsSubnet,
			rawOutputs: map[string]interface{}{
				"subnet_id":         "subnet-0abc123",
				"subnet_arn":        "arn:aws:ec2:us-west-2:123456789012:subnet/subnet-0abc123",
				"availability_zone": "us-west-2a",
				"cidr_block":        "10.0.1.0/24",
				"route_table_id":    "rtb-0abc123",
				"region":            "us-west-2",
			},
			mustPopulate: []string{
				"subnet_id", "subnet_arn", "availability_zone",
				"cidr_block", "route_table_id", "region",
			},
		},
		{
			// AwsInternetGateway: flat scalar outputs from both engines (gateway
			// id/arn, attached vpc id, region) must each land on the StackOutputs proto.
			name: "AwsInternetGateway",
			kind: cloudresourcekind.CloudResourceKind_AwsInternetGateway,
			rawOutputs: map[string]interface{}{
				"internet_gateway_id":  "igw-0abc123",
				"internet_gateway_arn": "arn:aws:ec2:us-west-2:123456789012:internet-gateway/igw-0abc123",
				"vpc_id":               "vpc-0abc123",
				"region":               "us-west-2",
			},
			mustPopulate: []string{
				"internet_gateway_id", "internet_gateway_arn", "vpc_id", "region",
			},
		},
		{
			// AwsEgressOnlyInternetGateway: flat scalar outputs from both engines
			// (gateway id, attached vpc id, region) must each land on the StackOutputs
			// proto. An egress-only gateway has no ARN, so none is emitted.
			name: "AwsEgressOnlyInternetGateway",
			kind: cloudresourcekind.CloudResourceKind_AwsEgressOnlyInternetGateway,
			rawOutputs: map[string]interface{}{
				"egress_only_internet_gateway_id": "eigw-0abc123",
				"vpc_id":                          "vpc-0abc123",
				"region":                          "us-west-2",
			},
			mustPopulate: []string{
				"egress_only_internet_gateway_id", "vpc_id", "region",
			},
		},
		{
			// AwsNatGateway: flat scalar outputs from both engines (gateway id,
			// public/private ip, ENI id, subnet id, region) must each land on the
			// StackOutputs proto. A NAT gateway has no ARN, so none is emitted.
			name: "AwsNatGateway",
			kind: cloudresourcekind.CloudResourceKind_AwsNatGateway,
			rawOutputs: map[string]interface{}{
				"nat_gateway_id":       "nat-0abc123",
				"public_ip":            "52.10.20.30",
				"private_ip":           "10.0.0.10",
				"network_interface_id": "eni-0abc123",
				"subnet_id":            "subnet-0abc123",
				"region":               "us-west-2",
			},
			mustPopulate: []string{
				"nat_gateway_id", "public_ip", "private_ip",
				"network_interface_id", "subnet_id", "region",
			},
		},
		{
			// AwsVpc: flat scalar outputs from both engines (vpc id/arn, primary and
			// IPv6 CIDR, owner, the route-table/default-resource ids, region) must
			// each land on the thin StackOutputs proto.
			name: "AwsVpc",
			kind: cloudresourcekind.CloudResourceKind_AwsVpc,
			rawOutputs: map[string]interface{}{
				"vpc_id":                    "vpc-0abc123",
				"vpc_arn":                   "arn:aws:ec2:us-west-2:123456789012:vpc/vpc-0abc123",
				"cidr_block":                "10.0.0.0/16",
				"ipv6_cidr_block":           "2600:1f18:abcd:1200::/56",
				"owner_id":                  "123456789012",
				"main_route_table_id":       "rtb-0abc123",
				"default_security_group_id": "sg-0abc123",
				"default_network_acl_id":    "acl-0abc123",
				"default_route_table_id":    "rtb-0abc123",
				"region":                    "us-west-2",
			},
			mustPopulate: []string{
				"vpc_id", "vpc_arn", "cidr_block", "ipv6_cidr_block", "owner_id",
				"main_route_table_id", "default_security_group_id",
				"default_network_acl_id", "default_route_table_id", "region",
			},
		},
		{
			// AwsIamPolicy: flat scalar outputs from both engines (policy arn/id/name)
			// must each land on the StackOutputs proto -- policy_arn is what role/user
			// attachments and permissions boundaries reference.
			name: "AwsIamPolicy",
			kind: cloudresourcekind.CloudResourceKind_AwsIamPolicy,
			rawOutputs: map[string]interface{}{
				"policy_arn":  "arn:aws:iam::123456789012:policy/s3-read-only",
				"policy_id":   "ANPAEXAMPLEID12345678",
				"policy_name": "s3-read-only",
			},
			mustPopulate: []string{"policy_arn", "policy_id", "policy_name"},
		},
		{
			// AwsIamInstanceProfile: flat scalar outputs from both engines (profile
			// arn/name/id and the carried role's name) must each land on the
			// StackOutputs proto -- instance_profile_arn is what EC2-shaped resources
			// reference.
			name: "AwsIamInstanceProfile",
			kind: cloudresourcekind.CloudResourceKind_AwsIamInstanceProfile,
			rawOutputs: map[string]interface{}{
				"instance_profile_arn":  "arn:aws:iam::123456789012:instance-profile/web-server",
				"instance_profile_name": "web-server",
				"instance_profile_id":   "AIPAEXAMPLEID12345678",
				"role_name":             "web-server-role",
			},
			mustPopulate: []string{
				"instance_profile_arn", "instance_profile_name",
				"instance_profile_id", "role_name",
			},
		},
		{
			// AwsIamRole: flat scalar outputs from both engines (role arn/name/id)
			// must each land on the StackOutputs proto. Guards the removal of the
			// role's former instance-profile outputs: EC2 delivery now composes
			// through AwsIamInstanceProfile, so the role emits only role-shaped
			// outputs.
			name: "AwsIamRole",
			kind: cloudresourcekind.CloudResourceKind_AwsIamRole,
			rawOutputs: map[string]interface{}{
				"role_arn":  "arn:aws:iam::123456789012:role/lambda-exec",
				"role_name": "lambda-exec",
				"role_id":   "AROAEXAMPLEID12345678",
			},
			mustPopulate: []string{"role_arn", "role_name", "role_id"},
		},
		{
			// AwsIamUser: flat scalar outputs from both engines (user arn/name/id,
			// access key id + base64 secret, console url) must each land on the
			// StackOutputs proto. The secret is base64-encoded by BOTH engines so the
			// emitted values are byte-identical.
			name: "AwsIamUser",
			kind: cloudresourcekind.CloudResourceKind_AwsIamUser,
			rawOutputs: map[string]interface{}{
				"user_arn":          "arn:aws:iam::123456789012:user/ci-deploy",
				"user_name":         "ci-deploy",
				"user_id":           "AIDAEXAMPLEID12345678",
				"access_key_id":     "AKIAEXAMPLEID1234567",
				"secret_access_key": "c2VjcmV0LWtleS1tYXRlcmlhbA==",
				"console_url":       "https://signin.aws.amazon.com/console",
			},
			mustPopulate: []string{
				"user_arn", "user_name", "user_id",
				"access_key_id", "secret_access_key", "console_url",
			},
		},
		{
			// AwsAlb: flat scalar outputs from both engines (arn/name/dns
			// name/hosted zone id/arn suffix) must each land on the StackOutputs
			// proto -- load_balancer_arn is what listeners attach through, the DNS
			// pair is what Route53 alias records consume, and arn_suffix is the
			// CloudWatch LoadBalancer dimension request-count autoscaling scopes on.
			name: "AwsAlb",
			kind: cloudresourcekind.CloudResourceKind_AwsAlb,
			rawOutputs: map[string]interface{}{
				"load_balancer_arn":            "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/demo/50dc6c495c0c9188",
				"load_balancer_name":           "demo",
				"load_balancer_dns_name":       "demo-1234567890.us-west-2.elb.amazonaws.com",
				"load_balancer_hosted_zone_id": "Z1H1FL5HABSF5",
				"arn_suffix":                   "app/demo/50dc6c495c0c9188",
			},
			mustPopulate: []string{
				"load_balancer_arn", "load_balancer_name",
				"load_balancer_dns_name", "load_balancer_hosted_zone_id",
				"arn_suffix",
			},
		},
		{
			// AwsNlb: the same four load-balancer scalars as AwsAlb. Guards the
			// load-balancer-only output shape: the NLB emits no listener or target
			// group outputs because those are first-class kinds with their own
			// outputs.
			name: "AwsNlb",
			kind: cloudresourcekind.CloudResourceKind_AwsNlb,
			rawOutputs: map[string]interface{}{
				"load_balancer_arn":            "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/net/demo/50dc6c495c0c9188",
				"load_balancer_name":           "demo",
				"load_balancer_dns_name":       "demo-1234567890.elb.us-west-2.amazonaws.com",
				"load_balancer_hosted_zone_id": "Z18D5FSROUN65G",
			},
			mustPopulate: []string{
				"load_balancer_arn", "load_balancer_name",
				"load_balancer_dns_name", "load_balancer_hosted_zone_id",
			},
		},
		{
			// AwsLbTargetGroup: flat scalar outputs from both engines (arn, the
			// possibly-truncated name, and the CloudWatch arn_suffix) must each land
			// on the StackOutputs proto -- target_group_arn is what listener forward
			// actions, ECS services, and ASG attachments reference.
			name: "AwsLbTargetGroup",
			kind: cloudresourcekind.CloudResourceKind_AwsLbTargetGroup,
			rawOutputs: map[string]interface{}{
				"target_group_arn":  "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/api/943f017f100becff",
				"target_group_name": "api",
				"arn_suffix":        "targetgroup/api/943f017f100becff",
			},
			mustPopulate: []string{"target_group_arn", "target_group_name", "arn_suffix"},
		},
		{
			// AwsLbListener: a single flat output -- listener_arn is what listener
			// rules attach through.
			name: "AwsLbListener",
			kind: cloudresourcekind.CloudResourceKind_AwsLbListener,
			rawOutputs: map[string]interface{}{
				"listener_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:listener/app/demo/50dc6c495c0c9188/f2f7dc8efc522ab2",
			},
			mustPopulate: []string{"listener_arn"},
		},
		{
			// AwsLbListenerRule: the rule ARN plus the AWS-assigned priority.
			// Priority is emitted as a STRING by both engines (Terraform's tostring
			// and the Pulumi module's strconv conversion) so the shapes stay
			// byte-identical -- this case guards that contract.
			name: "AwsLbListenerRule",
			kind: cloudresourcekind.CloudResourceKind_AwsLbListenerRule,
			rawOutputs: map[string]interface{}{
				"rule_arn": "arn:aws:elasticloadbalancing:us-west-2:123456789012:listener-rule/app/demo/50dc6c495c0c9188/f2f7dc8efc522ab2/9683b2d02a6cabee",
				"priority": "10",
			},
			mustPopulate: []string{"rule_arn", "priority"},
		},
		{
			// AwsEcsTaskDefinition: the revision-carrying ARN is the handle ECS
			// services reference (each new revision changes it and rolls the
			// service); revision is an int64 proto field fed from numeric engine
			// outputs, guarding the numeric-to-int64 flattening.
			name: "AwsEcsTaskDefinition",
			kind: cloudresourcekind.CloudResourceKind_AwsEcsTaskDefinition,
			rawOutputs: map[string]interface{}{
				"task_definition_arn":  "arn:aws:ecs:us-west-2:123456789012:task-definition/api:7",
				"arn_without_revision": "arn:aws:ecs:us-west-2:123456789012:task-definition/api",
				"family":               "api",
				"revision":             float64(7),
				"log_group_name":       "/ecs/api",
				"log_group_arn":        "arn:aws:logs:us-west-2:123456789012:log-group:/ecs/api",
			},
			mustPopulate: []string{
				"task_definition_arn", "arn_without_revision", "family",
				"revision", "log_group_name", "log_group_arn",
			},
		},
		{
			// AwsEcsService: flat scalar outputs -- the service ARN encodes both
			// the cluster and service names (the E2E verifier's key), and the
			// cluster/task-definition ARNs are republished resolved references.
			name: "AwsEcsService",
			kind: cloudresourcekind.CloudResourceKind_AwsEcsService,
			rawOutputs: map[string]interface{}{
				"service_arn":         "arn:aws:ecs:us-west-2:123456789012:service/prod/api",
				"service_name":        "api",
				"cluster_arn":         "arn:aws:ecs:us-west-2:123456789012:cluster/prod",
				"task_definition_arn": "arn:aws:ecs:us-west-2:123456789012:task-definition/api:7",
			},
			mustPopulate: []string{
				"service_arn", "service_name", "cluster_arn", "task_definition_arn",
			},
		},
		{
			// AwsLaunchTemplate: the template id/arn plus the two version numbers.
			// latest_version and default_version are int64 proto fields fed from
			// numeric engine outputs (Terraform's number, Pulumi's IntOutput) --
			// this case guards that numeric outputs flatten onto int64 fields.
			name: "AwsLaunchTemplate",
			kind: cloudresourcekind.CloudResourceKind_AwsLaunchTemplate,
			rawOutputs: map[string]interface{}{
				"launch_template_id":  "lt-0123456789abcdef0",
				"launch_template_arn": "arn:aws:ec2:us-west-2:123456789012:launch-template/lt-0123456789abcdef0",
				"latest_version":      float64(3),
				"default_version":     float64(3),
			},
			mustPopulate: []string{"launch_template_id", "launch_template_arn", "latest_version", "default_version"},
		},
		{
			// AwsAutoScalingGroup: flat scalar outputs -- the group name is the
			// CloudWatch dimension and ECS capacity-provider handle; the ARN scopes
			// IAM policies and EventBridge rules.
			name: "AwsAutoScalingGroup",
			kind: cloudresourcekind.CloudResourceKind_AwsAutoScalingGroup,
			rawOutputs: map[string]interface{}{
				"autoscaling_group_name": "web",
				"autoscaling_group_arn":  "arn:aws:autoscaling:us-west-2:123456789012:autoScalingGroup:uuid:autoScalingGroupName/web",
			},
			mustPopulate: []string{"autoscaling_group_name", "autoscaling_group_arn"},
		},
		{
			// AwsEksAddon: flat scalar outputs -- the ARN keys the E2E verifier
			// (it encodes cluster and add-on names); addon_version reports the
			// resolved AWS default when the spec pinned nothing.
			name: "AwsEksAddon",
			kind: cloudresourcekind.CloudResourceKind_AwsEksAddon,
			rawOutputs: map[string]interface{}{
				"addon_arn":     "arn:aws:eks:us-west-2:123456789012:addon/platform/vpc-cni/9ac7ab21-1a2b",
				"addon_name":    "vpc-cni",
				"addon_version": "v1.18.1-eksbuild.3",
			},
			mustPopulate: []string{"addon_arn", "addon_name", "addon_version"},
		},
		{
			// AwsEksFargateProfile: flat scalar outputs -- the ARN keys the E2E
			// verifier (it encodes cluster and profile names); status is ACTIVE
			// after a successful create.
			name: "AwsEksFargateProfile",
			kind: cloudresourcekind.CloudResourceKind_AwsEksFargateProfile,
			rawOutputs: map[string]interface{}{
				"fargate_profile_arn":  "arn:aws:eks:us-west-2:123456789012:fargateprofile/platform/serverless/9ac7ab21-1a2b",
				"fargate_profile_name": "serverless",
				"status":               "ACTIVE",
			},
			mustPopulate: []string{"fargate_profile_arn", "fargate_profile_name", "status"},
		},
		{
			// AwsEksAccessEntry: flat scalar outputs -- the entry ARN keys the E2E
			// verifier (it encodes the cluster and the principal identity), and the
			// resolved principal ARN is what downstream references consume.
			name: "AwsEksAccessEntry",
			kind: cloudresourcekind.CloudResourceKind_AwsEksAccessEntry,
			rawOutputs: map[string]interface{}{
				"access_entry_arn": "arn:aws:eks:us-west-2:123456789012:access-entry/platform/role/123456789012/TeamViewerRole/9ac7ab21-1a2b",
				"principal_arn":    "arn:aws:iam::123456789012:role/TeamViewerRole",
			},
			mustPopulate: []string{"access_entry_arn", "principal_arn"},
		},
		{
			// AwsVpcEndpoint: the endpoint id keys the E2E verifier; prefix_list_id
			// is the gateway-endpoint route/security-group handle; dns_name +
			// hosted_zone_id compose Route53 aliases to interface endpoints; and
			// network_interface_ids guards list outputs flattening onto a repeated
			// string field.
			name: "AwsVpcEndpoint",
			kind: cloudresourcekind.CloudResourceKind_AwsVpcEndpoint,
			rawOutputs: map[string]interface{}{
				"vpc_endpoint_id":       "vpce-0123456789abcdef0",
				"arn":                   "arn:aws:ec2:us-west-2:123456789012:vpc-endpoint/vpce-0123456789abcdef0",
				"state":                 "available",
				"prefix_list_id":        "pl-68a54001",
				"dns_name":              "vpce-0123456789abcdef0-abcd1234.sts.us-west-2.vpce.amazonaws.com",
				"hosted_zone_id":        "Z1K56Z6FNPJRR",
				"network_interface_ids": []interface{}{"eni-0123456789abcdef0", "eni-0f9e8d7c6b5a43210"},
			},
			mustPopulate: []string{
				"vpc_endpoint_id", "arn", "state", "prefix_list_id",
				"dns_name", "hosted_zone_id", "network_interface_ids",
			},
		},
		{
			// AwsRdsCluster: the identifier keys the E2E verifier; endpoint +
			// reader_endpoint are the connection handles downstream references
			// consume; master_user_secret_arn carries the AWS-managed credential
			// handle; and instance_endpoints guards list outputs flattening onto a
			// repeated string field (the folded per-name cluster instances).
			name: "AwsRdsCluster",
			kind: cloudresourcekind.CloudResourceKind_AwsRdsCluster,
			rawOutputs: map[string]interface{}{
				"cluster_identifier":              "orders-db",
				"arn":                             "arn:aws:rds:us-west-2:123456789012:cluster:orders-db",
				"cluster_resource_id":             "cluster-ABCDEFGHIJKL01234",
				"endpoint":                        "orders-db.cluster-abc123.us-west-2.rds.amazonaws.com",
				"reader_endpoint":                 "orders-db.cluster-ro-abc123.us-west-2.rds.amazonaws.com",
				"port":                            5432,
				"hosted_zone_id":                  "Z1PVIF0B656C1W",
				"engine_version_actual":           "16.4",
				"master_user_secret_arn":          "arn:aws:secretsmanager:us-west-2:123456789012:secret:rds!cluster-abc-def",
				"db_subnet_group_name":            "orders-db",
				"db_cluster_parameter_group_name": "default.aurora-postgresql16",
				"instance_endpoints":              []interface{}{"orders-db-writer.abc123.us-west-2.rds.amazonaws.com"},
			},
			mustPopulate: []string{
				"cluster_identifier", "arn", "cluster_resource_id", "endpoint",
				"reader_endpoint", "port", "hosted_zone_id", "engine_version_actual",
				"master_user_secret_arn", "db_subnet_group_name",
				"db_cluster_parameter_group_name", "instance_endpoints",
			},
		},
		{
			// AwsRdsInstance: the identifier keys the E2E verifier; endpoint is
			// address:port while address is the bare hostname (both are real AWS
			// attributes downstream references consume differently); resource_id is
			// the durable handle for IAM auth policies and point-in-time restores.
			name: "AwsRdsInstance",
			kind: cloudresourcekind.CloudResourceKind_AwsRdsInstance,
			rawOutputs: map[string]interface{}{
				"instance_identifier":    "billing-db",
				"arn":                    "arn:aws:rds:us-west-2:123456789012:db:billing-db",
				"resource_id":            "db-ABCDEFGHIJKL01234",
				"endpoint":               "billing-db.abc123.us-west-2.rds.amazonaws.com:5432",
				"address":                "billing-db.abc123.us-west-2.rds.amazonaws.com",
				"port":                   5432,
				"hosted_zone_id":         "Z1PVIF0B656C1W",
				"engine_version_actual":  "16.4",
				"master_user_secret_arn": "arn:aws:secretsmanager:us-west-2:123456789012:secret:rds!db-abc-def",
				"db_subnet_group_name":   "billing-db",
			},
			mustPopulate: []string{
				"instance_identifier", "arn", "resource_id", "endpoint", "address",
				"port", "hosted_zone_id", "engine_version_actual",
				"master_user_secret_arn", "db_subnet_group_name",
			},
		},
		{
			name: "AwsElasticacheUser",
			kind: cloudresourcekind.CloudResourceKind_AwsElasticacheUser,
			rawOutputs: map[string]interface{}{
				"user_id":   "app-cache-user",
				"arn":       "arn:aws:elasticache:us-west-2:123456789012:user:app-cache-user",
				"user_name": "app-cache-user",
			},
			mustPopulate: []string{"user_id", "arn", "user_name"},
		},
		{
			name: "AwsElasticacheUserGroup",
			kind: cloudresourcekind.CloudResourceKind_AwsElasticacheUserGroup,
			rawOutputs: map[string]interface{}{
				"user_group_id": "app-cache-group",
				"arn":           "arn:aws:elasticache:us-west-2:123456789012:usergroup:app-cache-group",
			},
			mustPopulate: []string{"user_group_id", "arn"},
		},
		{
			name: "AwsRedisElasticache",
			kind: cloudresourcekind.CloudResourceKind_AwsRedisElasticache,
			rawOutputs: map[string]interface{}{
				"replication_group_id":           "orders-cache",
				"primary_endpoint_address":       "orders-cache.abc123.usw2.cache.amazonaws.com",
				"reader_endpoint_address":        "orders-cache-ro.abc123.usw2.cache.amazonaws.com",
				"configuration_endpoint_address": "",
				"arn":                            "arn:aws:elasticache:us-west-2:123456789012:replicationgroup:orders-cache",
				"port":                           6379,
				"subnet_group_name":              "orders-cache",
				"parameter_group_name":           "orders-cache-custom",
				"engine_version_actual":          "7.1.0",
			},
			mustPopulate: []string{
				"replication_group_id", "primary_endpoint_address", "reader_endpoint_address",
				"arn", "port", "subnet_group_name", "parameter_group_name", "engine_version_actual",
			},
		},
		{
			name: "AwsMemcachedElasticache",
			kind: cloudresourcekind.CloudResourceKind_AwsMemcachedElasticache,
			rawOutputs: map[string]interface{}{
				"cluster_id":             "session-cache",
				"cluster_address":        "session-cache.abc123.cfg.usw2.cache.amazonaws.com",
				"configuration_endpoint": "session-cache.abc123.cfg.usw2.cache.amazonaws.com:11211",
				"arn":                    "arn:aws:elasticache:us-west-2:123456789012:cluster:session-cache",
				"port":                   11211,
				"subnet_group_name":      "session-cache",
				"parameter_group_name":   "session-cache-custom",
			},
			mustPopulate: []string{
				"cluster_id", "cluster_address", "configuration_endpoint",
				"arn", "port", "subnet_group_name", "parameter_group_name",
			},
		},
		{
			name: "AwsServerlessElasticache",
			kind: cloudresourcekind.CloudResourceKind_AwsServerlessElasticache,
			rawOutputs: map[string]interface{}{
				"arn":                     "arn:aws:elasticache:us-west-2:123456789012:serverlesscache:orders-srvless",
				"endpoint_address":        "orders-srvless-abc123.serverless.usw2.cache.amazonaws.com",
				"endpoint_port":           6379,
				"reader_endpoint_address": "orders-srvless-abc123-ro.serverless.usw2.cache.amazonaws.com",
				"reader_endpoint_port":    6380,
				"full_engine_version":     "7.1.0",
				"name":                    "orders-srvless",
			},
			mustPopulate: []string{
				"arn", "endpoint_address", "endpoint_port", "reader_endpoint_address",
				"reader_endpoint_port", "full_engine_version", "name",
			},
		},
		{
			// AwsDocumentDb: the identifier keys the E2E verifier; endpoint +
			// reader_endpoint are the connection handles downstream references
			// consume; master_user_secret_arn carries the AWS-managed credential
			// handle; and instance_endpoints guards list outputs flattening onto a
			// repeated string field (the folded per-name cluster instances).
			name: "AwsDocumentDb",
			kind: cloudresourcekind.CloudResourceKind_AwsDocumentDb,
			rawOutputs: map[string]interface{}{
				"cluster_identifier":              "catalog-docdb",
				"arn":                             "arn:aws:rds:us-west-2:123456789012:cluster:catalog-docdb",
				"cluster_resource_id":             "cluster-ABCDEFGHIJKL01234",
				"endpoint":                        "catalog-docdb.cluster-abc123.us-west-2.docdb.amazonaws.com",
				"reader_endpoint":                 "catalog-docdb.cluster-ro-abc123.us-west-2.docdb.amazonaws.com",
				"port":                            27017,
				"hosted_zone_id":                  "ZNKXH85TT8WVW",
				"engine_version_actual":           "5.0.0",
				"master_user_secret_arn":          "arn:aws:secretsmanager:us-west-2:123456789012:secret:rds!cluster-abc-def",
				"db_subnet_group_name":            "catalog-docdb",
				"db_cluster_parameter_group_name": "default.docdb5.0",
				"instance_endpoints":              []interface{}{"catalog-docdb-writer.abc123.us-west-2.docdb.amazonaws.com"},
			},
			mustPopulate: []string{
				"cluster_identifier", "arn", "cluster_resource_id", "endpoint",
				"reader_endpoint", "port", "hosted_zone_id", "engine_version_actual",
				"master_user_secret_arn", "db_subnet_group_name",
				"db_cluster_parameter_group_name", "instance_endpoints",
			},
		},
		{
			// AwsNeptuneCluster: the identifier keys the E2E verifier; the
			// cluster_resource_id is the durable handle IAM database-auth policies
			// scope to; and instance_endpoints guards list outputs flattening onto
			// a repeated string field. This case also guards the Terraform module's
			// first-ever outputs.tf (its absence was a live cross-engine parity bug).
			name: "AwsNeptuneCluster",
			kind: cloudresourcekind.CloudResourceKind_AwsNeptuneCluster,
			rawOutputs: map[string]interface{}{
				"cluster_identifier":                   "knowledge-graph",
				"arn":                                  "arn:aws:rds:us-west-2:123456789012:cluster:knowledge-graph",
				"cluster_resource_id":                  "cluster-ABCDEFGHIJKL01234",
				"endpoint":                             "knowledge-graph.cluster-abc123.us-west-2.neptune.amazonaws.com",
				"reader_endpoint":                      "knowledge-graph.cluster-ro-abc123.us-west-2.neptune.amazonaws.com",
				"port":                                 8182,
				"hosted_zone_id":                       "Z2T2AVZR3PGPQK",
				"engine_version_actual":                "1.4.5.1",
				"neptune_subnet_group_name":            "knowledge-graph",
				"neptune_cluster_parameter_group_name": "default.neptune1.4",
				"instance_endpoints":                   []interface{}{"knowledge-graph-writer.abc123.us-west-2.neptune.amazonaws.com"},
			},
			mustPopulate: []string{
				"cluster_identifier", "arn", "cluster_resource_id", "endpoint",
				"reader_endpoint", "port", "hosted_zone_id", "engine_version_actual",
				"neptune_subnet_group_name", "neptune_cluster_parameter_group_name",
				"instance_endpoints",
			},
		},
		{
			// AwsRedshiftCluster: the identifier keys the E2E verifier; endpoint +
			// dns_name are the connection handles downstream references consume;
			// cluster_namespace_arn is the data-sharing/Data-API handle; and
			// master_password_secret_arn carries the AWS-managed credential handle.
			name: "AwsRedshiftCluster",
			kind: cloudresourcekind.CloudResourceKind_AwsRedshiftCluster,
			rawOutputs: map[string]interface{}{
				"cluster_identifier":         "analytics-warehouse",
				"cluster_arn":                "arn:aws:redshift:us-west-2:123456789012:cluster:analytics-warehouse",
				"cluster_namespace_arn":      "arn:aws:redshift:us-west-2:123456789012:namespace:abc12345-6789-0abc-def1-234567890abc",
				"endpoint":                   "analytics-warehouse.abc123.us-west-2.redshift.amazonaws.com:5439",
				"dns_name":                   "analytics-warehouse.abc123.us-west-2.redshift.amazonaws.com",
				"database_name":              "analytics",
				"port":                       5439,
				"subnet_group_name":          "analytics-warehouse",
				"parameter_group_name":       "analytics-warehouse",
				"master_password_secret_arn": "arn:aws:secretsmanager:us-west-2:123456789012:secret:redshift!analytics-abc",
			},
			mustPopulate: []string{
				"cluster_identifier", "cluster_arn", "cluster_namespace_arn",
				"endpoint", "dns_name", "database_name", "port",
				"subnet_group_name", "parameter_group_name",
				"master_password_secret_arn",
			},
		},
		{
			// AwsRedshiftServerlessNamespace: namespace_name is the join key
			// workgroups attach with (downstream references resolve against
			// stack outputs, never metadata); admin_password_secret_arn
			// carries the AWS-managed credential handle.
			name: "AwsRedshiftServerlessNamespace",
			kind: cloudresourcekind.CloudResourceKind_AwsRedshiftServerlessNamespace,
			rawOutputs: map[string]interface{}{
				"namespace_name":            "analytics-data",
				"namespace_id":              "abc12345-6789-0abc-def1-234567890abc",
				"arn":                       "arn:aws:redshift-serverless:us-west-2:123456789012:namespace/abc12345-6789-0abc-def1-234567890abc",
				"db_name":                   "analytics",
				"admin_password_secret_arn": "arn:aws:secretsmanager:us-west-2:123456789012:secret:redshift!analytics-data-admin-abc",
			},
			mustPopulate: []string{
				"namespace_name", "namespace_id", "arn", "db_name",
				"admin_password_secret_arn",
			},
		},
		{
			// AwsRedshiftServerlessWorkgroup: workgroup_name keys the E2E
			// verifier and the credentials API; endpoint_address + port are
			// the connection handles downstream references consume.
			name: "AwsRedshiftServerlessWorkgroup",
			kind: cloudresourcekind.CloudResourceKind_AwsRedshiftServerlessWorkgroup,
			rawOutputs: map[string]interface{}{
				"workgroup_name":   "analytics-compute",
				"workgroup_id":     "def67890-1234-5abc-def6-789012345def",
				"arn":              "arn:aws:redshift-serverless:us-west-2:123456789012:workgroup/def67890-1234-5abc-def6-789012345def",
				"endpoint_address": "analytics-compute.123456789012.us-west-2.redshift-serverless.amazonaws.com",
				"port":             5439,
			},
			mustPopulate: []string{
				"workgroup_name", "workgroup_id", "arn", "endpoint_address",
				"port",
			},
		},
		{
			// AwsDynamodb: table_name/table_arn are the join keys IAM policies
			// and application config consume; stream_arn is what Lambda
			// event-source mappings attach to when streams are enabled.
			name: "AwsDynamodb",
			kind: cloudresourcekind.CloudResourceKind_AwsDynamodb,
			rawOutputs: map[string]interface{}{
				"table_name":   "orders",
				"table_arn":    "arn:aws:dynamodb:us-west-2:123456789012:table/orders",
				"table_id":     "orders",
				"stream_arn":   "arn:aws:dynamodb:us-west-2:123456789012:table/orders/stream/2026-07-04T00:00:00.000",
				"stream_label": "2026-07-04T00:00:00.000",
			},
			mustPopulate: []string{
				"table_name", "table_arn", "table_id", "stream_arn",
				"stream_label",
			},
		},
		{
			// AwsMskCluster: cluster_arn keys the E2E verifier and IAM
			// policies; the bootstrap_brokers_* family carries the
			// per-listener connection strings clients consume (each engine
			// emits every variant, empty when the listener is off); and
			// configuration_arn surfaces the module-managed configuration
			// folded from server_properties.
			name: "AwsMskCluster",
			kind: cloudresourcekind.CloudResourceKind_AwsMskCluster,
			rawOutputs: map[string]interface{}{
				"cluster_arn":                                   "arn:aws:kafka:us-west-2:123456789012:cluster/orders-streaming/abc12345-6789-0abc-def1-234567890abc-2",
				"cluster_name":                                  "orders-streaming",
				"cluster_uuid":                                  "abc12345-6789-0abc-def1-234567890abc-2",
				"current_version":                               "K3AEGXETSR30VB",
				"bootstrap_brokers":                             "b-1.orders.abc123.c2.kafka.us-west-2.amazonaws.com:9092",
				"bootstrap_brokers_tls":                         "b-1.orders.abc123.c2.kafka.us-west-2.amazonaws.com:9094",
				"bootstrap_brokers_sasl_iam":                    "b-1.orders.abc123.c2.kafka.us-west-2.amazonaws.com:9098",
				"bootstrap_brokers_sasl_scram":                  "b-1.orders.abc123.c2.kafka.us-west-2.amazonaws.com:9096",
				"bootstrap_brokers_public_tls":                  "b-1-public.orders.abc123.c2.kafka.us-west-2.amazonaws.com:9194",
				"bootstrap_brokers_public_sasl_iam":             "b-1-public.orders.abc123.c2.kafka.us-west-2.amazonaws.com:9198",
				"bootstrap_brokers_public_sasl_scram":           "b-1-public.orders.abc123.c2.kafka.us-west-2.amazonaws.com:9196",
				"bootstrap_brokers_vpc_connectivity_tls":        "b-1.orders.abc123.c2.kafka.us-west-2.amazonaws.com:14001",
				"bootstrap_brokers_vpc_connectivity_sasl_iam":   "b-1.orders.abc123.c2.kafka.us-west-2.amazonaws.com:14003",
				"bootstrap_brokers_vpc_connectivity_sasl_scram": "b-1.orders.abc123.c2.kafka.us-west-2.amazonaws.com:14002",
				"zookeeper_connect_string":                      "z-1.orders.abc123.c2.kafka.us-west-2.amazonaws.com:2181",
				"zookeeper_connect_string_tls":                  "z-1.orders.abc123.c2.kafka.us-west-2.amazonaws.com:2182",
				"configuration_arn":                             "arn:aws:kafka:us-west-2:123456789012:configuration/orders-streaming/def67890-1234-5abc-def6-789012345def-3",
			},
			mustPopulate: []string{
				"cluster_arn", "cluster_name", "cluster_uuid", "current_version",
				"bootstrap_brokers", "bootstrap_brokers_tls",
				"bootstrap_brokers_sasl_iam", "bootstrap_brokers_sasl_scram",
				"bootstrap_brokers_public_tls", "bootstrap_brokers_public_sasl_iam",
				"bootstrap_brokers_public_sasl_scram",
				"bootstrap_brokers_vpc_connectivity_tls",
				"bootstrap_brokers_vpc_connectivity_sasl_iam",
				"bootstrap_brokers_vpc_connectivity_sasl_scram",
				"zookeeper_connect_string", "zookeeper_connect_string_tls",
				"configuration_arn",
			},
		},
		{
			// AwsMskServerlessCluster: cluster_arn keys the E2E verifier and
			// the kafka-cluster:* IAM policies clients need;
			// bootstrap_brokers_sasl_iam is the only connection string
			// serverless MSK exposes (SASL/IAM is its sole auth scheme).
			name: "AwsMskServerlessCluster",
			kind: cloudresourcekind.CloudResourceKind_AwsMskServerlessCluster,
			rawOutputs: map[string]interface{}{
				"cluster_arn":                "arn:aws:kafka:us-west-2:123456789012:cluster/events-kafka/abc12345-6789-0abc-def1-234567890abc-s1",
				"cluster_name":               "events-kafka",
				"cluster_uuid":               "abc12345-6789-0abc-def1-234567890abc-s1",
				"bootstrap_brokers_sasl_iam": "boot-abc123.c1.kafka-serverless.us-west-2.amazonaws.com:9098",
			},
			mustPopulate: []string{
				"cluster_arn", "cluster_name", "cluster_uuid",
				"bootstrap_brokers_sasl_iam",
			},
		},
		{
			// AwsSecurityGroup: security_group_id is the join key every
			// attach-shaped kind references; security_group_arn is the form
			// IAM policy conditions expect; owner_id enables cross-account
			// rule references (<owner_id>/<group_id>).
			name: "AwsSecurityGroup",
			kind: cloudresourcekind.CloudResourceKind_AwsSecurityGroup,
			rawOutputs: map[string]interface{}{
				"security_group_id":  "sg-0123456789abcdef0",
				"security_group_arn": "arn:aws:ec2:us-west-2:123456789012:security-group/sg-0123456789abcdef0",
				"owner_id":           "123456789012",
			},
			mustPopulate: []string{
				"security_group_id", "security_group_arn", "owner_id",
			},
		},
		{
			// AwsLambda: function_name keys the E2E verifier; function_arn is
			// the join key for event-source mappings and IAM policies; invoke_arn
			// is what API Gateway integrations consume.
			name: "AwsLambda",
			kind: cloudresourcekind.CloudResourceKind_AwsLambda,
			rawOutputs: map[string]interface{}{
				"function_arn":   "arn:aws:lambda:us-west-2:123456789012:function:planton-oss-e2e-lambda-smoke",
				"function_name":  "planton-oss-e2e-lambda-smoke",
				"invoke_arn":     "arn:aws:apigateway:us-west-2:lambda:path/2015-03-31/functions/arn:aws:lambda:us-west-2:123456789012:function:planton-oss-e2e-lambda-smoke/invocations",
				"qualified_arn":  "",
				"version":        "",
				"function_url":   "",
				"alias_arns":     map[string]interface{}{},
				"log_group_name": "/aws/lambda/planton-oss-e2e-lambda-smoke",
			},
			mustPopulate: []string{
				"function_arn", "function_name", "invoke_arn", "log_group_name",
			},
		},
		{
			// AwsKmsKey: key_id keys the E2E verifier; key_arn is the join key
			// encryption-at-rest fields reference; alias_names carries the human-
			// friendly addresses SDK callers may use instead of the key ID.
			name: "AwsKmsKey",
			kind: cloudresourcekind.CloudResourceKind_AwsKmsKey,
			rawOutputs: map[string]interface{}{
				"key_id":      "12345678-1234-1234-1234-123456789012",
				"key_arn":     "arn:aws:kms:us-west-2:123456789012:key/12345678-1234-1234-1234-123456789012",
				"alias_names": []interface{}{"alias/planton-oss-e2e-kms-smoke"},
			},
			mustPopulate: []string{
				"key_id", "key_arn", "alias_names",
			},
		},
		{
			// AwsSqsQueue: queue_url is the SQS API handle; queue_arn is the
			// IAM/cross-service join key (DLQ targets, SNS subscriptions,
			// Lambda event source mappings); queue_name keys the E2E verifier.
			name: "AwsSqsQueue",
			kind: cloudresourcekind.CloudResourceKind_AwsSqsQueue,
			rawOutputs: map[string]interface{}{
				"queue_url":  "https://sqs.us-west-2.amazonaws.com/123456789012/planton-oss-e2e-sqs-smoke",
				"queue_arn":  "arn:aws:sqs:us-west-2:123456789012:planton-oss-e2e-sqs-smoke",
				"queue_name": "planton-oss-e2e-sqs-smoke",
			},
			mustPopulate: []string{
				"queue_url", "queue_arn", "queue_name",
			},
		},
		{
			// AwsSnsTopic: topic_arn is the subscription/EventBridge join key;
			// topic_name keys the E2E verifier; owner and beginning_archive_time
			// surface FIFO archive metadata when enabled.
			name: "AwsSnsTopic",
			kind: cloudresourcekind.CloudResourceKind_AwsSnsTopic,
			rawOutputs: map[string]interface{}{
				"topic_arn":              "arn:aws:sns:us-west-2:123456789012:planton-oss-e2e-sns-smoke",
				"topic_name":             "planton-oss-e2e-sns-smoke",
				"owner":                  "123456789012",
				"beginning_archive_time": "2026-07-04T12:00:00Z",
			},
			mustPopulate: []string{
				"topic_arn", "topic_name", "owner", "beginning_archive_time",
			},
		},
		{
			// AwsSnsSubscription: subscription_arn is the AWS identity and
			// unsubscribe handle; owner_id supports cross-account wiring;
			// pending_confirmation and confirmation_was_authenticated surface
			// the HTTP/email handshake lifecycle.
			name: "AwsSnsSubscription",
			kind: cloudresourcekind.CloudResourceKind_AwsSnsSubscription,
			rawOutputs: map[string]interface{}{
				"subscription_arn":               "arn:aws:sns:us-west-2:123456789012:planton-oss-e2e-sns-smoke:01234567-89ab-cdef-0123-456789abcdef",
				"owner_id":                       "123456789012",
				"pending_confirmation":           true,
				"confirmation_was_authenticated": true,
			},
			mustPopulate: []string{
				"subscription_arn", "owner_id",
				"pending_confirmation", "confirmation_was_authenticated",
			},
		},
		{
			// AwsEventBridgeBus: bus_name keys the E2E verifier and rule
			// event_bus_name references; bus_arn is the IAM/cross-account
			// join key.
			name: "AwsEventBridgeBus",
			kind: cloudresourcekind.CloudResourceKind_AwsEventBridgeBus,
			rawOutputs: map[string]interface{}{
				"bus_name": "planton-oss-e2e-eventbridge-bus-smoke",
				"bus_arn":  "arn:aws:events:us-west-2:123456789012:event-bus/planton-oss-e2e-eventbridge-bus-smoke",
			},
			mustPopulate: []string{
				"bus_name", "bus_arn",
			},
		},
		{
			// AwsEventBridgeRule: rule_arn is the IAM/monitoring join key;
			// rule_name keys the E2E verifier and EventBridge API calls.
			name: "AwsEventBridgeRule",
			kind: cloudresourcekind.CloudResourceKind_AwsEventBridgeRule,
			rawOutputs: map[string]interface{}{
				"rule_arn":  "arn:aws:events:us-west-2:123456789012:rule/planton-oss-e2e-eventbridge-bus-smoke/planton-oss-e2e-rule-smoke",
				"rule_name": "planton-oss-e2e-rule-smoke",
			},
			mustPopulate: []string{
				"rule_arn", "rule_name",
			},
		},
		{
			// AwsLambdaEventSourceMapping: uuid keys the E2E verifier;
			// mapping_arn and function_arn are the join keys downstream
			// automation consumes; state surfaces the last observed lifecycle.
			name: "AwsLambdaEventSourceMapping",
			kind: cloudresourcekind.CloudResourceKind_AwsLambdaEventSourceMapping,
			rawOutputs: map[string]interface{}{
				"uuid":         "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
				"mapping_arn":  "arn:aws:lambda:us-west-2:123456789012:event-source-mapping:a1b2c3d4-e5f6-7890-abcd-ef1234567890",
				"function_arn": "arn:aws:lambda:us-west-2:123456789012:function:planton-oss-e2e-lambda-smoke",
				"state":        "Enabled",
			},
			mustPopulate: []string{
				"uuid", "mapping_arn", "function_arn", "state",
			},
		},
		{
			// AwsMwaaEnvironment: environment_name keys the E2E verifier;
			// webserver_url is the operator's handle on the Airflow UI; the
			// two *_vpc_endpoint_service outputs are what CUSTOMER endpoint
			// management composes AwsVpcEndpoint nodes against.
			name: "AwsMwaaEnvironment",
			kind: cloudresourcekind.CloudResourceKind_AwsMwaaEnvironment,
			rawOutputs: map[string]interface{}{
				"environment_arn":                "arn:aws:airflow:us-west-2:123456789012:environment/prod-airflow",
				"environment_name":               "prod-airflow",
				"webserver_url":                  "abc123de-f456-7890-abcd-ef1234567890.c2.us-west-2.airflow.amazonaws.com",
				"airflow_version":                "2.10.1",
				"service_role_arn":               "arn:aws:iam::123456789012:role/aws-service-role/airflow.amazonaws.com/AWSServiceRoleForAmazonMWAA",
				"environment_class":              "mw1.medium",
				"status":                         "AVAILABLE",
				"created_at":                     "2026-07-04T12:00:00Z",
				"database_vpc_endpoint_service":  "com.amazonaws.vpce.us-west-2.vpce-svc-0123456789abcdef0",
				"webserver_vpc_endpoint_service": "com.amazonaws.vpce.us-west-2.vpce-svc-0fedcba9876543210",
			},
			mustPopulate: []string{
				"environment_arn", "environment_name", "webserver_url",
				"airflow_version", "service_role_arn", "environment_class",
				"status", "created_at", "database_vpc_endpoint_service",
				"webserver_vpc_endpoint_service",
			},
		},
		{
			// AwsOpenSearchDomain: domain_name keys the E2E verifier;
			// endpoint + dashboard_endpoint are the connection handles
			// downstream references consume; the *_v2 trio carries the
			// dual-stack endpoint surface added this session.
			name: "AwsOpenSearchDomain",
			kind: cloudresourcekind.CloudResourceKind_AwsOpenSearchDomain,
			rawOutputs: map[string]interface{}{
				"domain_id":                         "123456789012/search-logs",
				"domain_name":                       "search-logs",
				"domain_arn":                        "arn:aws:es:us-west-2:123456789012:domain/search-logs",
				"endpoint":                          "search-search-logs-abc123.us-west-2.es.amazonaws.com",
				"dashboard_endpoint":                "search-search-logs-abc123.us-west-2.es.amazonaws.com/_dashboards",
				"endpoint_v2":                       "search-logs-abc123.us-west-2.aos.amazonaws.com",
				"dashboard_endpoint_v2":             "search-logs-abc123.us-west-2.aos.amazonaws.com/_dashboards",
				"domain_endpoint_v2_hosted_zone_id": "Z1H1FL5HABSF5",
			},
			mustPopulate: []string{
				"domain_id", "domain_name", "domain_arn", "endpoint",
				"dashboard_endpoint", "endpoint_v2", "dashboard_endpoint_v2",
				"domain_endpoint_v2_hosted_zone_id",
			},
		},
		{
			// AwsEc2Instance: instance_id is the join key target groups
			// register; the address quartet carries the connection surface
			// (public values empty for private-only instances -- both
			// engines emit them regardless).
			name: "AwsEc2Instance",
			kind: cloudresourcekind.CloudResourceKind_AwsEc2Instance,
			rawOutputs: map[string]interface{}{
				"instance_id":                  "i-0123456789abcdef0",
				"arn":                          "arn:aws:ec2:us-west-2:123456789012:instance/i-0123456789abcdef0",
				"instance_state":               "running",
				"availability_zone":            "us-west-2a",
				"private_ip":                   "10.0.1.15",
				"private_dns":                  "ip-10-0-1-15.us-west-2.compute.internal",
				"public_ip":                    "",
				"public_dns":                   "",
				"primary_network_interface_id": "eni-0123456789abcdef0",
			},
			mustPopulate: []string{
				"instance_id", "arn", "instance_state", "availability_zone",
				"private_ip", "private_dns", "primary_network_interface_id",
			},
		},
		{
			// AwsEcsCluster: cluster_arn is the join key AwsEcsService
			// references; capacity_provider_names is the strategy
			// vocabulary (built-ins plus folded EC2 providers) and
			// capacity_provider_arns the folded providers' identities --
			// both list outputs, guarding list flattening.
			name: "AwsEcsCluster",
			kind: cloudresourcekind.CloudResourceKind_AwsEcsCluster,
			rawOutputs: map[string]interface{}{
				"cluster_name":            "prod-apps",
				"cluster_arn":             "arn:aws:ecs:us-west-2:123456789012:cluster/prod-apps",
				"capacity_provider_names": []interface{}{"FARGATE", "FARGATE_SPOT", "general-purpose"},
				"capacity_provider_arns":  []interface{}{"arn:aws:ecs:us-west-2:123456789012:capacity-provider/general-purpose"},
			},
			mustPopulate: []string{
				"cluster_name", "cluster_arn", "capacity_provider_names",
				"capacity_provider_arns",
			},
		},
		{
			// AwsS3Bucket: bucket_id (name) and bucket_arn are the join keys the
			// catalog's 12 consumer fields reference; the regional domain doubles
			// as the CloudFront origin domain; the website pair is exported as
			// empty strings when hosting is off so the output contract stays
			// shape-stable across both engines.
			name: "AwsS3Bucket",
			kind: cloudresourcekind.CloudResourceKind_AwsS3Bucket,
			rawOutputs: map[string]interface{}{
				"bucket_id":                   "planton-oss-e2e-awss3bucket-smoke",
				"bucket_arn":                  "arn:aws:s3:::planton-oss-e2e-awss3bucket-smoke",
				"region":                      "us-west-2",
				"bucket_regional_domain_name": "planton-oss-e2e-awss3bucket-smoke.s3.us-west-2.amazonaws.com",
				"bucket_domain_name":          "planton-oss-e2e-awss3bucket-smoke.s3.amazonaws.com",
				"hosted_zone_id":              "Z3BJ6K6RIION7M",
				"website_endpoint":            "planton-oss-e2e-awss3bucket-smoke.s3-website-us-west-2.amazonaws.com",
				"website_domain":              "s3-website-us-west-2.amazonaws.com",
			},
			mustPopulate: []string{
				"bucket_id", "bucket_arn", "region",
				"bucket_regional_domain_name", "bucket_domain_name",
				"hosted_zone_id", "website_endpoint", "website_domain",
			},
		},
		{
			// AwsKinesisStream: stream_arn is the join key consumers, Lambda
			// event source mappings, DynamoDB streaming destinations, and
			// Firehose sources reference; stream_name keys the E2E verifier.
			name: "AwsKinesisStream",
			kind: cloudresourcekind.CloudResourceKind_AwsKinesisStream,
			rawOutputs: map[string]interface{}{
				"stream_arn":  "arn:aws:kinesis:us-west-2:123456789012:stream/planton-oss-e2e-kinesis-smoke",
				"stream_name": "planton-oss-e2e-kinesis-smoke",
			},
			mustPopulate: []string{
				"stream_arn", "stream_name",
			},
		},
		{
			// AwsKinesisStreamConsumer: consumer_arn is the enhanced-fan-out
			// identity SubscribeToShard callers use; consumer_name keys the
			// E2E verifier; stream_arn echoes the parent join key.
			name: "AwsKinesisStreamConsumer",
			kind: cloudresourcekind.CloudResourceKind_AwsKinesisStreamConsumer,
			rawOutputs: map[string]interface{}{
				"consumer_arn":       "arn:aws:kinesis:us-west-2:123456789012:stream/planton-oss-e2e-kinesis-smoke/consumer/planton-oss-e2e-consumer-smoke:1751700000",
				"consumer_name":      "planton-oss-e2e-consumer-smoke",
				"stream_arn":         "arn:aws:kinesis:us-west-2:123456789012:stream/planton-oss-e2e-kinesis-smoke",
				"creation_timestamp": "2026-07-05T12:00:00Z",
			},
			mustPopulate: []string{
				"consumer_arn", "consumer_name", "stream_arn", "creation_timestamp",
			},
		},
		{
			// AwsKinesisFirehose: delivery_stream_arn is the IAM/EventBridge
			// join key; delivery_stream_name keys the E2E verifier and the
			// MSK broker-log delivery reference.
			name: "AwsKinesisFirehose",
			kind: cloudresourcekind.CloudResourceKind_AwsKinesisFirehose,
			rawOutputs: map[string]interface{}{
				"delivery_stream_arn":  "arn:aws:firehose:us-west-2:123456789012:deliverystream/planton-oss-e2e-firehose-smoke",
				"delivery_stream_name": "planton-oss-e2e-firehose-smoke",
			},
			mustPopulate: []string{
				"delivery_stream_arn", "delivery_stream_name",
			},
		},
		{
			// AwsEcrRepo: repository_url is what docker push/pull targets;
			// repository_arn scopes IAM policies; repository_name keys the
			// E2E verifier; registry_id is the owning account.
			name: "AwsEcrRepo",
			kind: cloudresourcekind.CloudResourceKind_AwsEcrRepo,
			rawOutputs: map[string]interface{}{
				"repository_name": "planton-oss-e2e/full-surface",
				"repository_url":  "123456789012.dkr.ecr.us-west-2.amazonaws.com/planton-oss-e2e/full-surface",
				"repository_arn":  "arn:aws:ecr:us-west-2:123456789012:repository/planton-oss-e2e/full-surface",
				"registry_id":     "123456789012",
			},
			mustPopulate: []string{
				"repository_name", "repository_url", "repository_arn", "registry_id",
			},
		},
		{
			// AwsRoute53Zone: zone_id is the join key every DNS-composing
			// resource references (records, ACM validation, ALB/NLB alias
			// registration) and keys the E2E verifier; nameservers carry the
			// registrar delegation values.
			name: "AwsRoute53Zone",
			kind: cloudresourcekind.CloudResourceKind_AwsRoute53Zone,
			rawOutputs: map[string]interface{}{
				"zone_id":             "Z1D633PJN98FT9",
				"zone_name":           "example.com",
				"nameservers":         []interface{}{"ns-1.awsdns-01.org", "ns-2.awsdns-02.com"},
				"primary_name_server": "ns-1.awsdns-01.org",
				"zone_arn":            "arn:aws:route53:::hostedzone/Z1D633PJN98FT9",
			},
			mustPopulate: []string{
				"zone_id", "zone_name", "nameservers", "primary_name_server", "zone_arn",
			},
		},
		{
			// AwsRoute53DnsRecord: fqdn + record_type + zone_id together key
			// the E2E verifier (a record has no standalone describe API);
			// is_alias and set_identifier echo the record's shape.
			name: "AwsRoute53DnsRecord",
			kind: cloudresourcekind.CloudResourceKind_AwsRoute53DnsRecord,
			rawOutputs: map[string]interface{}{
				"fqdn":           "canary.example.com",
				"record_type":    "A",
				"zone_id":        "Z1D633PJN98FT9",
				"is_alias":       false,
				"set_identifier": "canary",
			},
			mustPopulate: []string{
				"fqdn", "record_type", "zone_id", "set_identifier",
			},
		},
		{
			// AwsRoute53HealthCheck: health_check_id is what DNS records
			// reference (health_check_id) and calculated parents aggregate;
			// it also keys the E2E verifier.
			name: "AwsRoute53HealthCheck",
			kind: cloudresourcekind.CloudResourceKind_AwsRoute53HealthCheck,
			rawOutputs: map[string]interface{}{
				"health_check_id":  "abcdef11-2222-3333-4444-555555fedcba",
				"health_check_arn": "arn:aws:route53:::healthcheck/abcdef11-2222-3333-4444-555555fedcba",
			},
			mustPopulate: []string{
				"health_check_id", "health_check_arn",
			},
		},
		{
			// AwsCloudwatchLogGroup: log_group_arn is the FK target for Step
			// Functions logging, Route 53 query logging, API Gateway access
			// logs, and OpenSearch log publishing; log_group_name is the join
			// key for name-addressed consumers (ECS awslogs, ElastiCache) and
			// the E2E verifier.
			name: "AwsCloudwatchLogGroup",
			kind: cloudresourcekind.CloudResourceKind_AwsCloudwatchLogGroup,
			rawOutputs: map[string]interface{}{
				"log_group_arn":  "arn:aws:logs:us-west-2:123456789012:log-group:app-logs",
				"log_group_name": "app-logs",
			},
			mustPopulate: []string{
				"log_group_arn", "log_group_name",
			},
		},
		{
			// AwsCloudwatchAlarm: alarm_arn is referenced by ECS service
			// rollback alarms and ASG instance-refresh alarms; alarm_name is
			// the join key composite alarm rules and actions suppressors use,
			// and keys the E2E verifier.
			name: "AwsCloudwatchAlarm",
			kind: cloudresourcekind.CloudResourceKind_AwsCloudwatchAlarm,
			rawOutputs: map[string]interface{}{
				"alarm_arn":  "arn:aws:cloudwatch:us-west-2:123456789012:alarm:cpu-high",
				"alarm_name": "cpu-high",
			},
			mustPopulate: []string{
				"alarm_arn", "alarm_name",
			},
		},
		{
			// AwsCloudwatchCompositeAlarm: alarm_name is how parent composite
			// alarms reference this one inside their own rule expressions;
			// it also keys the E2E verifier.
			name: "AwsCloudwatchCompositeAlarm",
			kind: cloudresourcekind.CloudResourceKind_AwsCloudwatchCompositeAlarm,
			rawOutputs: map[string]interface{}{
				"alarm_arn":  "arn:aws:cloudwatch:us-west-2:123456789012:alarm:shared-cause",
				"alarm_name": "shared-cause",
			},
			mustPopulate: []string{
				"alarm_arn", "alarm_name",
			},
		},
		{
			// AwsStepFunction: state_machine_arn is the FK target for
			// EventBridge targets and API Gateway service integrations;
			// state_machine_version_arn pins consumers to a published
			// snapshot when spec.publish is set. The name keys the E2E
			// verifier.
			name: "AwsStepFunction",
			kind: cloudresourcekind.CloudResourceKind_AwsStepFunction,
			rawOutputs: map[string]interface{}{
				"state_machine_arn":         "arn:aws:states:us-west-2:123456789012:stateMachine:orders",
				"state_machine_name":        "orders",
				"state_machine_version_arn": "arn:aws:states:us-west-2:123456789012:stateMachine:orders:1",
				"revision_id":               "aaaa1111-bbbb-2222-cccc-333344445555",
				"status":                    "ACTIVE",
				"creation_date":             "2026-07-07T00:00:00Z",
			},
			mustPopulate: []string{
				"state_machine_arn", "state_machine_name",
				"state_machine_version_arn", "revision_id", "status", "creation_date",
			},
		},
		{
			// AwsHttpApiGateway: api_id is the join key domain mappings
			// reference; execution_arn feeds Lambda resource policies;
			// stage_name composes into domain mappings.
			name: "AwsHttpApiGateway",
			kind: cloudresourcekind.CloudResourceKind_AwsHttpApiGateway,
			rawOutputs: map[string]interface{}{
				"api_id":           "a1b2c3d4",
				"api_endpoint":     "https://a1b2c3d4.execute-api.us-west-2.amazonaws.com",
				"api_arn":          "arn:aws:apigateway:us-west-2::/apis/a1b2c3d4",
				"execution_arn":    "arn:aws:execute-api:us-west-2:123456789012:a1b2c3d4",
				"stage_invoke_url": "https://a1b2c3d4.execute-api.us-west-2.amazonaws.com",
				"stage_name":       "$default",
			},
			mustPopulate: []string{
				"api_id", "api_endpoint", "api_arn",
				"execution_arn", "stage_invoke_url", "stage_name",
			},
		},
		{
			// AwsHttpApiVpcLink: vpc_link_id is what private integrations set
			// as connection_id; it also keys the E2E verifier.
			name: "AwsHttpApiVpcLink",
			kind: cloudresourcekind.CloudResourceKind_AwsHttpApiVpcLink,
			rawOutputs: map[string]interface{}{
				"vpc_link_id":  "abc123",
				"vpc_link_arn": "arn:aws:apigateway:us-west-2::/vpclinks/abc123",
			},
			mustPopulate: []string{"vpc_link_id", "vpc_link_arn"},
		},
		{
			// AwsHttpApiDomain: target_domain_name + hosted_zone_id are the
			// DNS composition surface (a Route 53 alias record targets them);
			// domain_name is the domain's join key and keys the E2E verifier.
			name: "AwsHttpApiDomain",
			kind: cloudresourcekind.CloudResourceKind_AwsHttpApiDomain,
			rawOutputs: map[string]interface{}{
				"domain_name":        "api.example.com",
				"domain_name_arn":    "arn:aws:apigateway:us-west-2::/domainnames/api.example.com",
				"target_domain_name": "d-abc123.execute-api.us-west-2.amazonaws.com",
				"hosted_zone_id":     "Z2OJLYMUO9EFXC",
			},
			mustPopulate: []string{
				"domain_name", "domain_name_arn", "target_domain_name", "hosted_zone_id",
			},
		},
		{
			// AwsCognitoUserPool: issuer is the JWT-authorizer join key (the
			// scheme-carrying spelling of user_pool_endpoint); user_pool_domain
			// is the RAW domain string ALB authenticate-cognito actions take;
			// the CloudFront trio composes a custom domain's DNS alias record.
			name: "AwsCognitoUserPool",
			kind: cloudresourcekind.CloudResourceKind_AwsCognitoUserPool,
			rawOutputs: map[string]interface{}{
				"user_pool_id":                "us-west-2_Ab1Cd2EfG",
				"user_pool_arn":               "arn:aws:cognito-idp:us-west-2:123456789012:userpool/us-west-2_Ab1Cd2EfG",
				"user_pool_endpoint":          "cognito-idp.us-west-2.amazonaws.com/us-west-2_Ab1Cd2EfG",
				"issuer":                      "https://cognito-idp.us-west-2.amazonaws.com/us-west-2_Ab1Cd2EfG",
				"user_pool_domain":            "myapp-auth",
				"hosted_ui_url":               "https://myapp-auth.auth.us-west-2.amazoncognito.com",
				"cloudfront_distribution":     "d111abcdef8.cloudfront.net",
				"cloudfront_distribution_arn": "arn:aws:cloudfront::123456789012:distribution/E1ABCDEF",
				"cloudfront_hosted_zone_id":   "Z2FDTNDATAQYW2",
			},
			mustPopulate: []string{
				"user_pool_id", "user_pool_arn", "user_pool_endpoint", "issuer",
				"user_pool_domain", "hosted_ui_url",
				"cloudfront_distribution", "cloudfront_distribution_arn", "cloudfront_hosted_zone_id",
			},
		},
		{
			// AwsCognitoIdentityProvider: provider_name is the sole
			// integration identifier (IdPs have no ARN) -- app clients list it
			// in supported_identity_providers.
			name: "AwsCognitoIdentityProvider",
			kind: cloudresourcekind.CloudResourceKind_AwsCognitoIdentityProvider,
			rawOutputs: map[string]interface{}{
				"provider_name": "Google",
				"provider_type": "Google",
				"user_pool_id":  "us-west-2_Ab1Cd2EfG",
			},
			mustPopulate: []string{"provider_name", "provider_type", "user_pool_id"},
		},
		{
			// AwsCognitoUserPoolClient: client_id is the join key JWT
			// authorizers list as an audience and ALB authenticate-cognito
			// actions take as user_pool_client_id; the secret only exists for
			// confidential clients.
			name: "AwsCognitoUserPoolClient",
			kind: cloudresourcekind.CloudResourceKind_AwsCognitoUserPoolClient,
			rawOutputs: map[string]interface{}{
				"client_id":     "1a2b3c4d5e6f7g8h9i0j",
				"client_secret": "shhh-not-a-real-secret",
				"user_pool_id":  "us-west-2_Ab1Cd2EfG",
			},
			mustPopulate: []string{"client_id", "client_secret", "user_pool_id"},
		},
		{
			// AwsCognitoResourceServer: scope_identifiers are the exact
			// strings app clients list in allowed_oauth_scopes; the identifier
			// keys the E2E verifier within its pool.
			name: "AwsCognitoResourceServer",
			kind: cloudresourcekind.CloudResourceKind_AwsCognitoResourceServer,
			rawOutputs: map[string]interface{}{
				"resource_server_identifier": "https://api.example.com",
				"scope_identifiers":          []interface{}{"https://api.example.com/read", "https://api.example.com/orders:write"},
				"user_pool_id":               "us-west-2_Ab1Cd2EfG",
			},
			mustPopulate: []string{"resource_server_identifier", "scope_identifiers", "user_pool_id"},
		},
		{
			// Guards the externaldns tofu module's output rename to solver_sa: the
			// module previously emitted "service_account_name", which does not flatten
			// onto the KubernetesExternalDnsStackOutputs.solver_sa proto field (the
			// Pulumi module already exported "solver_sa"). Both engines now emit the
			// same three outputs.
			name: "KubernetesExternalDns",
			kind: cloudresourcekind.CloudResourceKind_KubernetesExternalDns,
			rawOutputs: map[string]interface{}{
				"namespace":    "external-dns",
				"release_name": "gosilver-in-external-dns",
				"solver_sa":    "gosilver-in-external-dns",
			},
			mustPopulate: []string{"namespace", "release_name", "solver_sa"},
		},
		{
			// CloudflareR2Bucket: both engines emit the same outputs -- bucket name,
			// path-style S3 URL, the list of custom-domain URLs, and the managed
			// r2.dev public URL -- each of which must land on the StackOutputs proto.
			name: "CloudflareR2Bucket",
			kind: cloudresourcekind.CloudResourceKind_CloudflareR2Bucket,
			rawOutputs: map[string]interface{}{
				"bucket_name":        "media-assets",
				"bucket_url":         "https://00000000000000000000000000000000.r2.cloudflarestorage.com/media-assets",
				"custom_domain_urls": []interface{}{"https://media.example.com", "https://cdn.example.com"},
				"public_url":         "https://pub-0123456789abcdef.r2.dev",
			},
			mustPopulate: []string{"bucket_name", "bucket_url", "custom_domain_urls", "public_url"},
		},
		{
			// CloudflareD1Database: both engines emit the database id and name as
			// flat scalars (a Worker reaches D1 through its binding; no DSN exists).
			name: "CloudflareD1Database",
			kind: cloudresourcekind.CloudResourceKind_CloudflareD1Database,
			rawOutputs: map[string]interface{}{
				"database_id":   "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
				"database_name": "app-prod-db",
			},
			mustPopulate: []string{"database_id", "database_name"},
		},
		{
			// CloudflareKvNamespace: both engines emit the namespace id and the
			// url-encoding support flag.
			name: "CloudflareKvNamespace",
			kind: cloudresourcekind.CloudResourceKind_CloudflareKvNamespace,
			rawOutputs: map[string]interface{}{
				"namespace_id":          "0f1e2d3c4b5a69788796a5b4c3d2e1f0",
				"supports_url_encoding": true,
			},
			mustPopulate: []string{"namespace_id", "supports_url_encoding"},
		},
		{
			// CloudflareWorkersKvPair: both engines emit the entry key and the
			// namespace it was written to.
			name: "CloudflareWorkersKvPair",
			kind: cloudresourcekind.CloudResourceKind_CloudflareWorkersKvPair,
			rawOutputs: map[string]interface{}{
				"key_name":     "feature.new-dashboard",
				"namespace_id": "0f1e2d3c4b5a69788796a5b4c3d2e1f0",
			},
			mustPopulate: []string{"key_name", "namespace_id"},
		},
		{
			// CloudflareHyperdriveConfig: both engines emit the config id and name.
			name: "CloudflareHyperdriveConfig",
			kind: cloudresourcekind.CloudResourceKind_CloudflareHyperdriveConfig,
			rawOutputs: map[string]interface{}{
				"hyperdrive_id": "a1b2c3d4e5f60718293a4b5c6d7e8f90",
				"name":          "app-prod-pg",
			},
			mustPopulate: []string{"hyperdrive_id", "name"},
		},
		{
			// CloudflareQueue: both engines emit the queue id and name (referenced by
			// consumers, worker producer bindings, and R2 event notifications).
			name: "CloudflareQueue",
			kind: cloudresourcekind.CloudResourceKind_CloudflareQueue,
			rawOutputs: map[string]interface{}{
				"queue_id":    "a1b2c3d4e5f60718293a4b5c6d7e8f90",
				"queue_name":  "orders-queue",
				"created_on":  "2026-06-25T00:00:00Z",
				"modified_on": "2026-06-25T00:00:00Z",
			},
			mustPopulate: []string{"queue_id", "queue_name"},
		},
		{
			// CloudflarePagesProject: both engines emit the project name, its
			// pages.dev subdomain, attached custom domains, and creation time.
			name: "CloudflarePagesProject",
			kind: cloudresourcekind.CloudResourceKind_CloudflarePagesProject,
			rawOutputs: map[string]interface{}{
				"project_name": "marketing-site",
				"subdomain":    "marketing-site.pages.dev",
				"domains":      []interface{}{"www.example.com"},
				"created_on":   "2026-06-25T00:00:00Z",
			},
			mustPopulate: []string{"project_name", "subdomain"},
		},
		{
			// CloudflareDnsRecord: both engines emit the record id, name, type and
			// proxied flag as flat scalars onto the StackOutputs proto.
			name: "CloudflareDnsRecord",
			kind: cloudresourcekind.CloudResourceKind_CloudflareDnsRecord,
			rawOutputs: map[string]interface{}{
				"record_id":   "372e67954025e0ba6aaa6d586b9e0b59",
				"record_name": "www",
				"record_type": "A",
				"proxied":     true,
			},
			mustPopulate: []string{"record_id", "record_name", "record_type", "proxied"},
		},
		{
			// CloudflareDnsZone: both engines emit the zone id (scalar) and the
			// assigned nameservers (repeated string) onto the StackOutputs proto.
			name: "CloudflareDnsZone",
			kind: cloudresourcekind.CloudResourceKind_CloudflareDnsZone,
			rawOutputs: map[string]interface{}{
				"zone_id":                 "023e105f4ecef8ad9ca31a8372d0c353",
				"nameservers":             []interface{}{"ns1.cloudflare.com", "ns2.cloudflare.com"},
				"status":                  "active",
				"dnssec_status":           "active",
				"dnssec_ds":               "example.com. 3600 IN DS 2371 13 2 ABCDEF",
				"dnssec_digest":           "abcdef0123456789",
				"dnssec_digest_type":      "2",
				"dnssec_digest_algorithm": "SHA256",
				"dnssec_algorithm":        "13",
				"dnssec_key_tag":          "2371",
				"dnssec_public_key":       "mdsswUyr3DPW132mOi8V9xESWE8jTo0d",
				"dnssec_flags":            "257",
			},
			mustPopulate: []string{"zone_id", "nameservers", "status", "dnssec_ds", "dnssec_key_tag"},
		},
		{
			// CloudflareRuleset: both engines emit ruleset id, version, and the
			// zone_id/phase pass-throughs as flat scalars onto the proto.
			name: "CloudflareRuleset",
			kind: cloudresourcekind.CloudResourceKind_CloudflareRuleset,
			rawOutputs: map[string]interface{}{
				"ruleset_id": "2f2feab2026849078ba485f918791bdc",
				"version":    "3",
				"zone_id":    "023e105f4ecef8ad9ca31a8372d0c353",
				"phase":      "http_request_origin",
			},
			mustPopulate: []string{"ruleset_id", "version", "zone_id", "phase"},
		},
		{
			// CloudflareLoadBalancer: both engines emit the load balancer id,
			// hostname, and cname target as flat scalars onto the proto.
			name: "CloudflareLoadBalancer",
			kind: cloudresourcekind.CloudResourceKind_CloudflareLoadBalancer,
			rawOutputs: map[string]interface{}{
				"load_balancer_id":              "699d98642c564d2e855e9661899b7252",
				"load_balancer_dns_record_name": "lb.example.com",
				"load_balancer_cname_target":    "699d98642c564d2e855e9661899b7252",
			},
			mustPopulate: []string{"load_balancer_id", "load_balancer_dns_record_name", "load_balancer_cname_target"},
		},
		{
			// CloudflareLoadBalancerPool: both engines emit the pool id and name
			// (account-scoped pool referenced by load balancers).
			name: "CloudflareLoadBalancerPool",
			kind: cloudresourcekind.CloudResourceKind_CloudflareLoadBalancerPool,
			rawOutputs: map[string]interface{}{
				"pool_id":   "17b5962d775c646f3f9725cbc7a53df4",
				"pool_name": "web-pool",
			},
			mustPopulate: []string{"pool_id", "pool_name"},
		},
		{
			// CloudflareLoadBalancerMonitor: both engines emit the monitor id and
			// its protocol (account-scoped health check referenced by pools).
			name: "CloudflareLoadBalancerMonitor",
			kind: cloudresourcekind.CloudResourceKind_CloudflareLoadBalancerMonitor,
			rawOutputs: map[string]interface{}{
				"monitor_id":   "f1aba936b94213e5b8dca0c0dbf1f9cc",
				"monitor_type": "https",
			},
			mustPopulate: []string{"monitor_id", "monitor_type"},
		},
		{
			// CloudflareWorker: both engines emit the script id and name (scalars)
			// and the custom-domain hostnames / route patterns (repeated strings).
			name: "CloudflareWorker",
			kind: cloudresourcekind.CloudResourceKind_CloudflareWorker,
			rawOutputs: map[string]interface{}{
				"script_id":               "my-worker",
				"script_name":             "my-worker",
				"custom_domain_hostnames": []interface{}{"api.example.com"},
				"route_patterns":          []interface{}{"api.example.com/*"},
			},
			mustPopulate: []string{"script_id", "script_name", "custom_domain_hostnames", "route_patterns"},
		},
		{
			// CloudflareZeroTrustAccessApplication: both engines emit the
			// application id, audience tag, protected domain, and SaaS material.
			name: "CloudflareZeroTrustAccessApplication",
			kind: cloudresourcekind.CloudResourceKind_CloudflareZeroTrustAccessApplication,
			rawOutputs: map[string]interface{}{
				"application_id":     "f174e90a-fafe-4643-bbbc-4a0ed4fc8415",
				"aud":                "8a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b",
				"domain":             "dashboard.example.com",
				"saas_client_id":     "client-abc",
				"saas_client_secret": "secret-xyz",
				"saas_public_key":    "MIIBIjANBgkqh...",
				"saas_sso_endpoint":  "https://example.cloudflareaccess.com/cdn-cgi/access/sso/saml/abc",
				"saas_idp_entity_id": "https://example.cloudflareaccess.com",
			},
			mustPopulate: []string{
				"application_id", "aud", "domain", "saas_client_id", "saas_client_secret",
				"saas_public_key", "saas_sso_endpoint", "saas_idp_entity_id",
			},
		},
		{
			// CloudflareZeroTrustAccessPolicy: both engines emit the policy id.
			name: "CloudflareZeroTrustAccessPolicy",
			kind: cloudresourcekind.CloudResourceKind_CloudflareZeroTrustAccessPolicy,
			rawOutputs: map[string]interface{}{
				"policy_id": "699d98642c564d2e855e9661899b7252",
			},
			mustPopulate: []string{"policy_id"},
		},
		{
			// CloudflareZeroTrustAccessGroup: both engines emit the group id.
			name: "CloudflareZeroTrustAccessGroup",
			kind: cloudresourcekind.CloudResourceKind_CloudflareZeroTrustAccessGroup,
			rawOutputs: map[string]interface{}{
				"group_id": "aa9d98642c564d2e855e9661899b7252",
			},
			mustPopulate: []string{"group_id"},
		},
		{
			// CloudflareZeroTrustTunnel: both engines emit flat scalar outputs --
			// tunnel id, CNAME target, the (sensitive) connector token, status, the
			// account tag, and the creation timestamp.
			name: "CloudflareZeroTrustTunnel",
			kind: cloudresourcekind.CloudResourceKind_CloudflareZeroTrustTunnel,
			rawOutputs: map[string]interface{}{
				"tunnel_id":     "f70ff985-a4ef-4643-bbbc-4a0ed4fc8415",
				"tunnel_cname":  "f70ff985-a4ef-4643-bbbc-4a0ed4fc8415.cfargotunnel.com",
				"tunnel_token":  "eyJhIjoiMDc0NzU1YTc4ZDhlIn0=",
				"tunnel_status": "healthy",
				"account_tag":   "074755a78d8e8f77c119a90a125e8a06",
				"created_on":    "2026-06-25T12:00:00Z",
			},
			mustPopulate: []string{
				"tunnel_id", "tunnel_cname", "tunnel_token",
				"tunnel_status", "account_tag", "created_on",
			},
		},
		{
			// CloudflareZeroTrustTunnelVirtualNetwork: both engines emit the virtual
			// network id and name.
			name: "CloudflareZeroTrustTunnelVirtualNetwork",
			kind: cloudresourcekind.CloudResourceKind_CloudflareZeroTrustTunnelVirtualNetwork,
			rawOutputs: map[string]interface{}{
				"virtual_network_id":   "aaaa1111-bbbb-2222-cccc-333344445555",
				"virtual_network_name": "prod-vnet",
			},
			mustPopulate: []string{"virtual_network_id", "virtual_network_name"},
		},
		{
			// CloudflareZeroTrustTunnelRoute: both engines emit the route id and the
			// advertised CIDR.
			name: "CloudflareZeroTrustTunnelRoute",
			kind: cloudresourcekind.CloudResourceKind_CloudflareZeroTrustTunnelRoute,
			rawOutputs: map[string]interface{}{
				"route_id": "b8f2e1c0-1111-2222-3333-444455556666",
				"network":  "10.0.0.0/24",
			},
			mustPopulate: []string{"route_id", "network"},
		},
		{
			// CloudflareList: both engines emit the list id, name, and kind.
			name: "CloudflareList",
			kind: cloudresourcekind.CloudResourceKind_CloudflareList,
			rawOutputs: map[string]interface{}{
				"list_id": "2c0fc9fa937b11eaa1b71c4d701ab86e",
				"name":    "office_allowlist",
				"kind":    "ip",
			},
			mustPopulate: []string{"list_id", "name", "kind"},
		},
		{
			// CloudflareListItem: both engines emit the item id and parent list id.
			name: "CloudflareListItem",
			kind: cloudresourcekind.CloudResourceKind_CloudflareListItem,
			rawOutputs: map[string]interface{}{
				"item_id": "70c4e0c9b0e34f1a9b6f2d3c4a5b6c7d",
				"list_id": "2c0fc9fa937b11eaa1b71c4d701ab86e",
			},
			mustPopulate: []string{"item_id", "list_id"},
		},
		{
			// CloudflareTurnstileWidget: both engines emit the site key, the
			// (sensitive) secret, and timestamps.
			name: "CloudflareTurnstileWidget",
			kind: cloudresourcekind.CloudResourceKind_CloudflareTurnstileWidget,
			rawOutputs: map[string]interface{}{
				"sitekey":     "0x4AAAAAAA_examplesitekey",
				"secret":      "0x4AAAAAAA_examplesecretkey",
				"created_on":  "2026-06-25T00:00:00Z",
				"modified_on": "2026-06-25T00:00:00Z",
			},
			mustPopulate: []string{"sitekey", "secret"},
		},
		{
			// CloudflareEmailRoutingZone: both engines emit the zone id, enabled
			// flag, status, and name.
			name: "CloudflareEmailRoutingZone",
			kind: cloudresourcekind.CloudResourceKind_CloudflareEmailRoutingZone,
			rawOutputs: map[string]interface{}{
				"zone_id": "023e105f4ecef8ad9ca31a8372d0c353",
				"enabled": "true",
				"status":  "ready",
				"name":    "example.com",
			},
			mustPopulate: []string{"zone_id", "status", "name"},
		},
		{
			// CloudflareEmailRoutingRule: both engines emit the rule id and zone id.
			name: "CloudflareEmailRoutingRule",
			kind: cloudresourcekind.CloudResourceKind_CloudflareEmailRoutingRule,
			rawOutputs: map[string]interface{}{
				"rule_id": "a1b2c3d4e5f60718293a4b5c6d7e8f90",
				"zone_id": "023e105f4ecef8ad9ca31a8372d0c353",
			},
			mustPopulate: []string{"rule_id", "zone_id"},
		},
		{
			// CloudflareEmailRoutingAddress: both engines emit the address id,
			// email, and timestamps.
			name: "CloudflareEmailRoutingAddress",
			kind: cloudresourcekind.CloudResourceKind_CloudflareEmailRoutingAddress,
			rawOutputs: map[string]interface{}{
				"address_id": "b8f2e1c0a1b2c3d4e5f60718293a4b5c",
				"email":      "ops@example.com",
				"created":    "2026-06-25T00:00:00Z",
			},
			mustPopulate: []string{"address_id", "email"},
		},
		{
			// CloudflareOriginCaCertificate: both engines emit the certificate id,
			// the certificate PEM, the (sensitive) generated private key, and expiry.
			name: "CloudflareOriginCaCertificate",
			kind: cloudresourcekind.CloudResourceKind_CloudflareOriginCaCertificate,
			rawOutputs: map[string]interface{}{
				"certificate_id": "b8f2e1c0a1b2c3d4e5f60718293a4b5c",
				"certificate":    "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----",
				"private_key":    "-----BEGIN PRIVATE KEY-----\nMIIE\n-----END PRIVATE KEY-----",
				"expires_on":     "2041-06-25T00:00:00Z",
			},
			mustPopulate: []string{"certificate_id", "certificate"},
		},
		{
			// CloudflareCertificatePack: both engines emit the pack id, status, and
			// primary certificate id.
			name: "CloudflareCertificatePack",
			kind: cloudresourcekind.CloudResourceKind_CloudflareCertificatePack,
			rawOutputs: map[string]interface{}{
				"certificate_pack_id": "3822ff90e3534420ac41fc7e4a1f4b07",
				"status":              "active",
				"primary_certificate": "caa875a3-b2f0-4f7e-9a1e-0d2b4c6e8f10",
			},
			mustPopulate: []string{"certificate_pack_id", "status"},
		},
		{
			// CloudflareCustomHostname: both engines emit the hostname id, status,
			// the ownership-verification records, and the creation timestamp.
			name: "CloudflareCustomHostname",
			kind: cloudresourcekind.CloudResourceKind_CloudflareCustomHostname,
			rawOutputs: map[string]interface{}{
				"custom_hostname_id":               "0d89c70f8d4f4b1aa1b5d2e3f4a5b6c7",
				"status":                           "pending",
				"ownership_verification_name":      "_cf-custom-hostname.support.acme.com",
				"ownership_verification_type":      "txt",
				"ownership_verification_value":     "1f2e3d4c5b6a7988",
				"ownership_verification_http_url":  "http://support.acme.com/.well-known/cf-custom-hostname-challenge/0d89",
				"ownership_verification_http_body": "1f2e3d4c5b6a7988",
				"verification_errors":              []interface{}{},
				"created_at":                       "2026-06-25T00:00:00Z",
			},
			mustPopulate: []string{"custom_hostname_id", "status"},
		},
		{
			// CloudflareCustomHostnameFallbackOrigin: both engines emit status and
			// timestamps for the zone's fallback origin.
			name: "CloudflareCustomHostnameFallbackOrigin",
			kind: cloudresourcekind.CloudResourceKind_CloudflareCustomHostnameFallbackOrigin,
			rawOutputs: map[string]interface{}{
				"status":     "active",
				"created_at": "2026-06-25T00:00:00Z",
				"updated_at": "2026-06-25T00:00:00Z",
				"errors":     []interface{}{},
			},
			mustPopulate: []string{"status"},
		},
		{
			// AwsCertManagerCert: cert_arn is the join key every TLS consumer
			// references (listeners, CloudFront, Cognito, OpenSearch, Client
			// VPN); domain_validation_records guards the repeated-message
			// shape external-DNS users consume to create their validation
			// CNAMEs; status is what the E2E verifier keys on (a no-zone cert
			// rests in PENDING_VALIDATION).
			name: "AwsCertManagerCert",
			kind: cloudresourcekind.CloudResourceKind_AwsCertManagerCert,
			rawOutputs: map[string]interface{}{
				"cert_arn": "arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012",
				"status":   "PENDING_VALIDATION",
				"domain_validation_records": []interface{}{
					map[string]interface{}{
						"domain_name":  "example.com",
						"record_name":  "_3839f23e624907e70b9e.example.com.",
						"record_type":  "CNAME",
						"record_value": "_632077f7a35f9d.mhbtsbpdnt.acm-validations.aws.",
					},
				},
				"not_before":       "",
				"not_after":        "",
				"certificate_type": "AMAZON_ISSUED",
			},
			mustPopulate: []string{
				"cert_arn", "status", "certificate_type", "domain_validation_records",
			},
		},
		{
			// AwsCloudFront: distribution_id keys the E2E verifier and
			// invalidation requests; domain_name + hosted_zone_id are what
			// Route53 alias records compose against; distribution_arn is the
			// WAF-association join key.
			name: "AwsCloudFront",
			kind: cloudresourcekind.CloudResourceKind_AwsCloudFront,
			rawOutputs: map[string]interface{}{
				"distribution_id":  "E2ABCDEF123456",
				"distribution_arn": "arn:aws:cloudfront::123456789012:distribution/E2ABCDEF123456",
				"domain_name":      "d123abc456def.cloudfront.net",
				"hosted_zone_id":   "Z2FDTNDATAQYW2",
				"status":           "Deployed",
			},
			mustPopulate: []string{
				"distribution_id", "distribution_arn", "domain_name",
				"hosted_zone_id", "status",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidateOverride(tc.kind, genericModuleDir, tc.rawOutputs)
			if err != nil {
				t.Fatalf("ValidateOverride failed: %v", err)
			}
			if len(result.SchemaErrors) != 0 {
				t.Fatalf("unexpected schema errors: %v", result.SchemaErrors)
			}
			if result.DryRun == nil {
				t.Fatal("expected a dry-run result")
			}

			// Core invariant: every emitted output lands on a proto field. A
			// regression to a flat/mismatched output name surfaces here.
			if len(result.DryRun.UnmappedOutputs) != 0 {
				t.Errorf("%s: outputs did not map onto the StackOutputs proto: %v",
					tc.kind.String(), result.DryRun.UnmappedOutputs)
			}

			populated := make(map[string]bool, len(result.DryRun.PopulatedFields))
			for _, f := range result.DryRun.PopulatedFields {
				populated[f.ProtoField] = true
			}
			for _, field := range tc.mustPopulate {
				if !populated[field] {
					t.Errorf("%s: expected proto field %q to be populated, but it was not",
						tc.kind.String(), field)
				}
			}
		})
	}
}

// TestStackOutputsConformance_DetectsFlatSecretDrift proves the guard actually
// catches the historical drift: the pre-fix Postgres tofu module emitted flat
// "password_secret_name"/"password_secret_key" outputs, which do NOT flatten onto
// the proto's password_secret{name,key} field. The guard must flag both the
// unmapped output and the unpopulated proto field.
func TestStackOutputsConformance_DetectsFlatSecretDrift(t *testing.T) {
	genericModuleDir := filepath.Join("testdata", "modules", "empty")
	kind := cloudresourcekind.CloudResourceKind_KubernetesPostgres

	flatDriftOutputs := map[string]interface{}{
		"namespace":            "gosilver-prod",
		"password_secret_name": "postgres.db-gosilver-prod-postgres.credentials.postgresql.acid.zalan.do",
		"password_secret_key":  "password",
	}

	result, err := ValidateOverride(kind, genericModuleDir, flatDriftOutputs)
	if err != nil {
		t.Fatalf("ValidateOverride failed: %v", err)
	}
	if result.DryRun == nil {
		t.Fatal("expected a dry-run result")
	}

	if len(result.DryRun.UnmappedOutputs) == 0 {
		t.Error("expected the flat password_secret_name/_key outputs to be reported as unmapped, but none were")
	}
	for _, f := range result.DryRun.PopulatedFields {
		if f.ProtoField == "password_secret" {
			t.Error("flat outputs must NOT populate the nested password_secret proto field")
		}
	}
}
