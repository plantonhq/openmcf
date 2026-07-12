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
			// name/hosted zone id) must each land on the StackOutputs proto --
			// load_balancer_arn is what listeners attach through, and the DNS pair
			// is what Route53 alias records consume.
			name: "AwsAlb",
			kind: cloudresourcekind.CloudResourceKind_AwsAlb,
			rawOutputs: map[string]interface{}{
				"load_balancer_arn":            "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/demo/50dc6c495c0c9188",
				"load_balancer_name":           "demo",
				"load_balancer_dns_name":       "demo-1234567890.us-west-2.elb.amazonaws.com",
				"load_balancer_hosted_zone_id": "Z1H1FL5HABSF5",
			},
			mustPopulate: []string{
				"load_balancer_arn", "load_balancer_name",
				"load_balancer_dns_name", "load_balancer_hosted_zone_id",
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
			// AzureResourceGroup: flat scalar outputs from both engines (ARM id,
			// name, region) must each land on the StackOutputs proto --
			// resource_group_name is the FK target every other Azure kind
			// references, and resource_group_id is the default scope for role
			// assignments.
			name: "AzureResourceGroup",
			kind: cloudresourcekind.CloudResourceKind_AzureResourceGroup,
			rawOutputs: map[string]interface{}{
				"resource_group_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg",
				"resource_group_name": "platform-rg",
				"region":              "eastus",
			},
			mustPopulate: []string{"resource_group_id", "resource_group_name", "region"},
		},
		{
			// AzureRoleAssignment: flat scalar outputs from both engines (the
			// fully-scoped assignment id, GUID name, scope, resolved role
			// definition id, principal id/type) must each land on the StackOutputs
			// proto -- role_assignment_id is what the authorization API and the
			// E2E verifier key on.
			name: "AzureRoleAssignment",
			kind: cloudresourcekind.CloudResourceKind_AzureRoleAssignment,
			rawOutputs: map[string]interface{}{
				"role_assignment_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg/providers/Microsoft.Authorization/roleAssignments/a67e1183-4b2d-4b6e-93f1-2b2b8d2e1c11",
				"name":               "a67e1183-4b2d-4b6e-93f1-2b2b8d2e1c11",
				"scope":              "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg",
				"role_definition_id": "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
				"principal_id":       "11111111-1111-1111-1111-111111111111",
				"principal_type":     "ServicePrincipal",
			},
			mustPopulate: []string{
				"role_assignment_id", "name", "scope",
				"role_definition_id", "principal_id", "principal_type",
			},
		},
		{
			// AzureRoleDefinition: scalar outputs plus a repeated string
			// (assignable_scopes) from both engines must land on the
			// StackOutputs proto -- role_definition_id carries the fully-scoped
			// ARM id (what an AzureRoleAssignment binds and what the E2E
			// verifier keys on), not the bare GUID (that is
			// role_definition_guid).
			name: "AzureRoleDefinition",
			kind: cloudresourcekind.CloudResourceKind_AzureRoleDefinition,
			rawOutputs: map[string]interface{}{
				"role_definition_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c",
				"role_definition_guid": "b24988ac-6180-42a0-ab88-20f7382dd24c",
				"role_name":            "acme-vm-operator",
				"scope":                "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg",
				"assignable_scopes":    []interface{}{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg"},
			},
			mustPopulate: []string{
				"role_definition_id", "role_definition_guid", "role_name",
				"scope", "assignable_scopes",
			},
		},
		{
			// AzureUserAssignedIdentity: the identity's three identifiers plus
			// its ARM id must land on the StackOutputs proto -- principal_id is
			// what role assignments grant to, client_id is what workloads
			// present to authenticate, identity_id is what consuming resources
			// and federated credentials attach to.
			name: "AzureUserAssignedIdentity",
			kind: cloudresourcekind.CloudResourceKind_AzureUserAssignedIdentity,
			rawOutputs: map[string]interface{}{
				"identity_id":  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/payments-api",
				"principal_id": "11111111-1111-1111-1111-111111111111",
				"client_id":    "22222222-2222-2222-2222-222222222222",
				"tenant_id":    "33333333-3333-3333-3333-333333333333",
			},
			mustPopulate: []string{
				"identity_id", "principal_id", "client_id", "tenant_id",
			},
		},
		{
			// AzureFederatedIdentityCredential: both engines export the
			// credential's ARM id plus the trust coordinates (issuer /
			// subject / audience) as deployed. audience is a single string on
			// the proto even though ARM's wire shape is a one-element list --
			// the Terraform module exports the sole element, matching the
			// Pulumi provider's flattened attribute.
			name: "AzureFederatedIdentityCredential",
			kind: cloudresourcekind.CloudResourceKind_AzureFederatedIdentityCredential,
			rawOutputs: map[string]interface{}{
				"federated_identity_credential_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/ci-deployer/federatedIdentityCredentials/github-main-branch",
				"name":                             "github-main-branch",
				"user_assigned_identity_id":        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/ci-deployer",
				"issuer":                           "https://token.actions.githubusercontent.com",
				"subject":                          "repo:acme/platform:ref:refs/heads/main",
				"audience":                         "api://AzureADTokenExchange",
			},
			mustPopulate: []string{
				"federated_identity_credential_id", "name",
				"user_assigned_identity_id", "issuer", "subject", "audience",
			},
		},
		{
			// AzureVirtualNetwork: scalar identifiers plus a repeated string
			// (address_spaces) from both engines must land on the StackOutputs
			// proto -- virtual_network_id is the join key subnets, peerings,
			// and DNS links attach through, and address_spaces reflects the
			// ACTUAL ranges (IPAM-provisioned when pools delegate allocation).
			name: "AzureVirtualNetwork",
			kind: cloudresourcekind.CloudResourceKind_AzureVirtualNetwork,
			rawOutputs: map[string]interface{}{
				"virtual_network_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet",
				"virtual_network_name": "prod-vnet",
				"guid":                 "44444444-4444-4444-4444-444444444444",
				"address_spaces":       []interface{}{"10.0.0.0/16", "10.1.0.0/16"},
			},
			mustPopulate: []string{
				"virtual_network_id", "virtual_network_name", "guid",
				"address_spaces",
			},
		},
		{
			// AzureRouteTable: both engines export the table's ARM id (the
			// join key subnets use to attach the table's routing policy) and
			// its name.
			name: "AzureRouteTable",
			kind: cloudresourcekind.CloudResourceKind_AzureRouteTable,
			rawOutputs: map[string]interface{}{
				"route_table_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/routeTables/egress-via-firewall",
				"route_table_name": "egress-via-firewall",
			},
			mustPopulate: []string{
				"route_table_id", "route_table_name",
			},
		},
		{
			// AzurePrivateDnsZone: the zone's ARM id (the join key links,
			// private endpoints, and databases attach through), its DNS name,
			// and its resource group (echoed for tooling that joins on
			// name+RG rather than parsing ARM ids).
			name: "AzurePrivateDnsZone",
			kind: cloudresourcekind.CloudResourceKind_AzurePrivateDnsZone,
			rawOutputs: map[string]interface{}{
				"zone_id":             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/privateDnsZones/privatelink.postgres.database.azure.com",
				"zone_name":           "privatelink.postgres.database.azure.com",
				"resource_group_name": "network-rg",
			},
			mustPopulate: []string{
				"zone_id", "zone_name", "resource_group_name",
			},
		},
		{
			// AzurePrivateDnsZoneVirtualNetworkLink: both engines export the
			// link's ARM id (a child of the zone:
			// {zone-id}/virtualNetworkLinks/{name}) and its name.
			name: "AzurePrivateDnsZoneVirtualNetworkLink",
			kind: cloudresourcekind.CloudResourceKind_AzurePrivateDnsZoneVirtualNetworkLink,
			rawOutputs: map[string]interface{}{
				"link_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/privateDnsZones/privatelink.postgres.database.azure.com/virtualNetworkLinks/hub-vnet",
				"link_name": "hub-vnet",
			},
			mustPopulate: []string{
				"link_id", "link_name",
			},
		},
		{
			// AzureSubnet: the catalog's most-referenced join key (subnet_id)
			// plus a repeated string (address_prefixes reflects the ACTUAL
			// ranges, IPAM-provisioned when a pool delegates allocation) and
			// the parent coordinates derived from the referenced network's
			// ARM id.
			name: "AzureSubnet",
			kind: cloudresourcekind.CloudResourceKind_AzureSubnet,
			rawOutputs: map[string]interface{}{
				"subnet_id":            "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app",
				"subnet_name":          "app",
				"address_prefixes":     []interface{}{"10.0.1.0/24"},
				"virtual_network_name": "prod-vnet",
				"resource_group_name":  "network-rg",
			},
			mustPopulate: []string{
				"subnet_id", "subnet_name", "address_prefixes",
				"virtual_network_name", "resource_group_name",
			},
		},
		{
			// AzureNetworkSecurityGroup: both engines export the group's ARM
			// id (the join key subnets attach through) and its name.
			name: "AzureNetworkSecurityGroup",
			kind: cloudresourcekind.CloudResourceKind_AzureNetworkSecurityGroup,
			rawOutputs: map[string]interface{}{
				"network_security_group_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/networkSecurityGroups/web-tier",
				"network_security_group_name": "web-tier",
			},
			mustPopulate: []string{
				"network_security_group_id", "network_security_group_name",
			},
		},
		{
			// AzurePublicIp: the address's ARM id (the join key gateways and
			// load balancers attach), the allocated address itself, and the
			// Azure-managed FQDN when a DNS label is set.
			name: "AzurePublicIp",
			kind: cloudresourcekind.CloudResourceKind_AzurePublicIp,
			rawOutputs: map[string]interface{}{
				"public_ip_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/publicIPAddresses/prod-frontend",
				"ip_address":     "20.42.1.1",
				"fqdn":           "prod-gateway.eastus.cloudapp.azure.com",
				"public_ip_name": "prod-frontend",
			},
			mustPopulate: []string{
				"public_ip_id", "ip_address", "fqdn", "public_ip_name",
			},
		},
		{
			// AzurePublicIpPrefix: the prefix's ARM id (referenced by public
			// IPs and NAT gateway associations) and the ACTUAL reserved CIDR
			// -- known only after creation, the value partners allowlist.
			name: "AzurePublicIpPrefix",
			kind: cloudresourcekind.CloudResourceKind_AzurePublicIpPrefix,
			rawOutputs: map[string]interface{}{
				"public_ip_prefix_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/publicIPPrefixes/prod-egress",
				"ip_prefix":             "20.42.0.16/28",
				"public_ip_prefix_name": "prod-egress",
			},
			mustPopulate: []string{
				"public_ip_prefix_id", "ip_prefix", "public_ip_prefix_name",
			},
		},
		{
			// AzureApplicationSecurityGroup: the group's ARM id is the
			// composition seam -- NIC ip configurations, scale-set network
			// profiles, and NSG rules reference it to declare membership or
			// target the group.
			name: "AzureApplicationSecurityGroup",
			kind: cloudresourcekind.CloudResourceKind_AzureApplicationSecurityGroup,
			rawOutputs: map[string]interface{}{
				"application_security_group_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/applicationSecurityGroups/web-tier",
				"application_security_group_name": "web-tier",
			},
			mustPopulate: []string{
				"application_security_group_id", "application_security_group_name",
			},
		},
		{
			// AzurePrivateEndpoint: the endpoint's ARM id, its private IP
			// (the address the service FQDN resolves to inside the VNet),
			// and the auto-created NIC's id.
			name: "AzurePrivateEndpoint",
			kind: cloudresourcekind.CloudResourceKind_AzurePrivateEndpoint,
			rawOutputs: map[string]interface{}{
				"private_endpoint_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/privateEndpoints/pg-pe",
				"private_endpoint_name": "pg-pe",
				"private_ip_address":    "10.0.1.10",
				"network_interface_id":  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/networkInterfaces/pg-pe-nic",
			},
			mustPopulate: []string{
				"private_endpoint_id", "private_endpoint_name", "private_ip_address", "network_interface_id",
			},
		},
		{
			// AzureDiskEncryptionSet: the set's ARM id (referenced by disks,
			// VMs, and scale sets) plus the system-assigned identity's
			// principal/tenant (the grant target for Key Vault crypto access).
			name: "AzureDiskEncryptionSet",
			kind: cloudresourcekind.CloudResourceKind_AzureDiskEncryptionSet,
			rawOutputs: map[string]interface{}{
				"disk_encryption_set_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg/providers/Microsoft.Compute/diskEncryptionSets/prod-des",
				"disk_encryption_set_name": "prod-des",
				"identity_principal_id":    "11111111-2222-3333-4444-555555555555",
				"identity_tenant_id":       "99999999-8888-7777-6666-555555555555",
			},
			mustPopulate: []string{
				"disk_encryption_set_id", "disk_encryption_set_name",
			},
		},
		{
			// AzureMssqlFailoverGroup: the group's ARM id plus the
			// DNS-composed listener endpoints -- the read-write listener is
			// the failover-following connection target downstream apps use.
			name: "AzureMssqlFailoverGroup",
			kind: cloudresourcekind.CloudResourceKind_AzureMssqlFailoverGroup,
			rawOutputs: map[string]interface{}{
				"failover_group_id":            "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Sql/servers/primary/failoverGroups/prod-sql-fog",
				"failover_group_name":          "prod-sql-fog",
				"read_write_listener_endpoint": "prod-sql-fog.database.windows.net",
				"read_only_listener_endpoint":  "prod-sql-fog.secondary.database.windows.net",
			},
			mustPopulate: []string{
				"failover_group_id", "failover_group_name", "read_write_listener_endpoint", "read_only_listener_endpoint",
			},
		},
		{
			// AzureMonitorActivityLogAlert: the alert's ARM id and name.
			name: "AzureMonitorActivityLogAlert",
			kind: cloudresourcekind.CloudResourceKind_AzureMonitorActivityLogAlert,
			rawOutputs: map[string]interface{}{
				"activity_log_alert_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Insights/activityLogAlerts/vm-delete-alert",
				"activity_log_alert_name": "vm-delete-alert",
			},
			mustPopulate: []string{
				"activity_log_alert_id", "activity_log_alert_name",
			},
		},
		{
			// AzureApplicationInsightsStandardWebTest: the test's ARM id
			// (referenced by a metric alert's web-test criteria), its name,
			// and the synthetic monitor id.
			name: "AzureApplicationInsightsStandardWebTest",
			kind: cloudresourcekind.CloudResourceKind_AzureApplicationInsightsStandardWebTest,
			rawOutputs: map[string]interface{}{
				"web_test_id":          "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Insights/webTests/homepage-health",
				"web_test_name":        "homepage-health",
				"synthetic_monitor_id": "homepage-health",
			},
			mustPopulate: []string{
				"web_test_id", "web_test_name", "synthetic_monitor_id",
			},
		},
		{
			// AzureLoadBalancer: the name-keyed maps are the composition
			// seams -- backend_pool_ids is what NIC ip_configurations and
			// scale-set network profiles join, nat_rule_ids is what a NIC's
			// NAT-rule association completes, probe_ids is what a scale
			// set's rolling-upgrade health probe references.
			name: "AzureLoadBalancer",
			kind: cloudresourcekind.CloudResourceKind_AzureLoadBalancer,
			rawOutputs: map[string]interface{}{
				"load_balancer_id":     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Network/loadBalancers/app-lb",
				"load_balancer_name":   "app-lb",
				"private_ip_address":   "10.0.1.6",
				"private_ip_addresses": []interface{}{"10.0.1.6"},
				"frontend_ip_configuration_ids": map[string]interface{}{
					"internal": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Network/loadBalancers/app-lb/frontendIPConfigurations/internal",
				},
				"backend_pool_ids": map[string]interface{}{
					"web": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Network/loadBalancers/app-lb/backendAddressPools/web",
				},
				"probe_ids": map[string]interface{}{
					"http-health": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Network/loadBalancers/app-lb/probes/http-health",
				},
				"nat_rule_ids": map[string]interface{}{
					"ssh-admin": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Network/loadBalancers/app-lb/inboundNatRules/ssh-admin",
				},
			},
			mustPopulate: []string{
				"load_balancer_id", "load_balancer_name",
				"private_ip_address", "private_ip_addresses",
				"frontend_ip_configuration_ids", "backend_pool_ids",
				"probe_ids", "nat_rule_ids",
			},
		},
		{
			// AzureNatGateway: the gateway's ARM id (the join key subnets
			// attach through), its name, and the ARM-assigned GUID.
			name: "AzureNatGateway",
			kind: cloudresourcekind.CloudResourceKind_AzureNatGateway,
			rawOutputs: map[string]interface{}{
				"nat_gateway_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/natGateways/prod-egress",
				"nat_gateway_name": "prod-egress",
				"resource_guid":    "55555555-5555-5555-5555-555555555555",
			},
			mustPopulate: []string{
				"nat_gateway_id", "nat_gateway_name", "resource_guid",
			},
		},
		{
			// AzureVirtualNetworkPeering: one direction's ARM id (a child of
			// the local network: {vnet-id}/virtualNetworkPeerings/{name}),
			// its name, and the local coordinates derived from the referenced
			// network's ARM id.
			name: "AzureVirtualNetworkPeering",
			kind: cloudresourcekind.CloudResourceKind_AzureVirtualNetworkPeering,
			rawOutputs: map[string]interface{}{
				"peering_id":           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hub-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/virtualNetworkPeerings/hub-to-spoke1",
				"peering_name":         "hub-to-spoke1",
				"virtual_network_name": "hub-vnet",
				"resource_group_name":  "hub-rg",
			},
			mustPopulate: []string{
				"peering_id", "peering_name", "virtual_network_name",
				"resource_group_name",
			},
		},
		{
			// AzureAksCluster: cluster_id is the parent seam every standalone
			// AzureAksNodePool consumes; oidc_issuer_url is the trust anchor
			// AzureFederatedIdentityCredential binds to.
			name: "AzureAksCluster",
			kind: cloudresourcekind.CloudResourceKind_AzureAksCluster,
			rawOutputs: map[string]interface{}{
				"cluster_id":                    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/aks-rg/providers/Microsoft.ContainerService/managedClusters/prod-aks",
				"cluster_name":                  "prod-aks",
				"fqdn":                          "prod-aks-abc123.hcp.eastus.azmk8s.io",
				"private_fqdn":                  "",
				"portal_fqdn":                   "prod-aks-abc123.privatelink.eastus.azmk8s.io",
				"oidc_issuer_url":               "https://eastus.oic.prod-aks.abc123.azmk8s.io/00000000-0000-0000-0000-000000000000/",
				"node_resource_group":           "MC_aks-rg_prod-aks_eastus",
				"node_resource_group_id":        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/MC_aks-rg_prod-aks_eastus",
				"cluster_kubeconfig":            "YXBpVmVyc2lvbjogdjE=",
				"cluster_identity_principal_id": "11111111-1111-1111-1111-111111111111",
				"kubelet_identity_object_id":    "22222222-2222-2222-2222-222222222222",
				"kubelet_identity_client_id":    "33333333-3333-3333-3333-333333333333",
				"current_kubernetes_version":    "1.35.2",
			},
			mustPopulate: []string{
				"cluster_id", "cluster_name", "fqdn", "oidc_issuer_url",
				"node_resource_group", "node_resource_group_id",
				"cluster_identity_principal_id", "current_kubernetes_version",
			},
		},
		{
			// AzureAksNodePool: the pool's ARM id (node_pool_id) is the
			// verification key; node_image_version reflects the OS patch
			// level actually rolled out.
			name: "AzureAksNodePool",
			kind: cloudresourcekind.CloudResourceKind_AzureAksNodePool,
			rawOutputs: map[string]interface{}{
				"node_pool_id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/aks-rg/providers/Microsoft.ContainerService/managedClusters/prod-aks/agentPools/general",
				"node_pool_name":     "general",
				"node_image_version": "AKSUbuntu-2204gen2containerd-202502.03.0",
			},
			mustPopulate: []string{
				"node_pool_id", "node_pool_name", "node_image_version",
			},
		},
		{
			// AzureContainerRegistry: the registry's ARM id is the seam AKS
			// clusters and AcrPull/AcrPush role assignments scope to;
			// login_server is what images are tagged with; the admin
			// credentials and data-endpoint hostnames are feature-gated
			// outputs (empty/absent when their features are off).
			name: "AzureContainerRegistry",
			kind: cloudresourcekind.CloudResourceKind_AzureContainerRegistry,
			rawOutputs: map[string]interface{}{
				"container_registry_id":                 "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg/providers/Microsoft.ContainerRegistry/registries/prodimages",
				"container_registry_name":               "prodimages",
				"login_server":                          "prodimages.azurecr.io",
				"admin_username":                        "prodimages",
				"admin_password":                        "s3cr3t-rotatable",
				"system_assigned_identity_principal_id": "44444444-4444-4444-4444-444444444444",
				"data_endpoint_host_names":              []interface{}{"prodimages.eastus.data.azurecr.io"},
			},
			mustPopulate: []string{
				"container_registry_id", "container_registry_name",
				"login_server", "admin_username", "admin_password",
				"system_assigned_identity_principal_id",
				"data_endpoint_host_names",
			},
		},
		{
			// AzureNetworkInterface: the NIC's ARM id is the seam
			// AzureVirtualMachine.network_interface_ids consumes; the
			// private address is what backends and DNS records key on; the
			// MAC populates only once attached to a running VM.
			name: "AzureNetworkInterface",
			kind: cloudresourcekind.CloudResourceKind_AzureNetworkInterface,
			rawOutputs: map[string]interface{}{
				"network_interface_id":        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Network/networkInterfaces/app-nic",
				"network_interface_name":      "app-nic",
				"private_ip_address":          "10.0.1.4",
				"private_ip_addresses":        []interface{}{"10.0.1.4"},
				"mac_address":                 "00-0D-3A-1B-2C-3D",
				"internal_domain_name_suffix": "abc123.bx.internal.cloudapp.net",
			},
			mustPopulate: []string{
				"network_interface_id", "network_interface_name",
				"private_ip_address", "private_ip_addresses",
				"mac_address", "internal_domain_name_suffix",
			},
		},
		{
			// AzureManagedDisk: the disk's ARM id is the seam
			// AzureVirtualMachine.data_disk_attachments consumes; the
			// actual size matters for COPY/FROM_IMAGE disks that inherited
			// the source's size.
			name: "AzureManagedDisk",
			kind: cloudresourcekind.CloudResourceKind_AzureManagedDisk,
			rawOutputs: map[string]interface{}{
				"disk_id":      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Compute/disks/orders-db-data",
				"disk_name":    "orders-db-data",
				"disk_size_gb": 256,
			},
			mustPopulate: []string{
				"disk_id", "disk_name", "disk_size_gb",
			},
		},
		{
			// AzureVirtualMachineScaleSet: the scale set's ARM id is the
			// seam a standalone VM's Flexible-attach consumes and what
			// autoscale/monitoring scope to; the system-assigned principal
			// is the AzureRoleAssignment seam (UNIFORM sets).
			name: "AzureVirtualMachineScaleSet",
			kind: cloudresourcekind.CloudResourceKind_AzureVirtualMachineScaleSet,
			rawOutputs: map[string]interface{}{
				"scale_set_id":                          "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Compute/virtualMachineScaleSets/web-fleet",
				"scale_set_name":                        "web-fleet",
				"unique_id":                             "88888888-8888-8888-8888-888888888888",
				"system_assigned_identity_principal_id": "99999999-9999-9999-9999-999999999999",
			},
			mustPopulate: []string{
				"scale_set_id", "scale_set_name", "unique_id",
				"system_assigned_identity_principal_id",
			},
		},
		{
			// AzureVirtualMachine: vm_id is what grants and policies scope
			// to; the identity principal is the AzureRoleAssignment seam;
			// the IP conveniences aggregate from the referenced NICs.
			name: "AzureVirtualMachine",
			kind: cloudresourcekind.CloudResourceKind_AzureVirtualMachine,
			rawOutputs: map[string]interface{}{
				"vm_id":                                 "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Compute/virtualMachines/app-vm",
				"vm_name":                               "app-vm",
				"virtual_machine_guid":                  "66666666-6666-6666-6666-666666666666",
				"private_ip_address":                    "10.0.1.4",
				"public_ip_address":                     "",
				"computer_name":                         "app-vm",
				"system_assigned_identity_principal_id": "77777777-7777-7777-7777-777777777777",
			},
			mustPopulate: []string{
				"vm_id", "vm_name", "virtual_machine_guid",
				"private_ip_address", "computer_name",
				"system_assigned_identity_principal_id",
			},
		},
		{
			// AzureKeyVault: the vault's ARM id is the seam keys,
			// certificates, VM/VMSS secret blocks, and vault-scoped role
			// assignments reference; vault_uri is the data-plane endpoint
			// applications call.
			name: "AzureKeyVault",
			kind: cloudresourcekind.CloudResourceKind_AzureKeyVault,
			rawOutputs: map[string]interface{}{
				"key_vault_id":        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/security-rg/providers/Microsoft.KeyVault/vaults/platform-kv",
				"key_vault_name":      "platform-kv",
				"vault_uri":           "https://platform-kv.vault.azure.net/",
				"tenant_id":           "255aadad-0000-0000-0000-000000000000",
				"resource_group_name": "security-rg",
			},
			mustPopulate: []string{
				"key_vault_id", "key_vault_name", "vault_uri",
				"tenant_id", "resource_group_name",
			},
		},
		{
			// AzureKeyVaultKey: versionless_id is the CMK seam consumers
			// (ACR encryption) reference so rotation propagates; key_id
			// pins a version (the AKS KMS grain); the ARM proxy ids serve
			// control-plane integrations.
			name: "AzureKeyVaultKey",
			kind: cloudresourcekind.CloudResourceKind_AzureKeyVaultKey,
			rawOutputs: map[string]interface{}{
				"key_id":                  "https://platform-kv.vault.azure.net/keys/storage-cmk/abc123def456",
				"versionless_id":          "https://platform-kv.vault.azure.net/keys/storage-cmk",
				"key_name":                "storage-cmk",
				"version":                 "abc123def456",
				"resource_id":             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/security-rg/providers/Microsoft.KeyVault/vaults/platform-kv/keys/storage-cmk/versions/abc123def456",
				"resource_versionless_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/security-rg/providers/Microsoft.KeyVault/vaults/platform-kv/keys/storage-cmk",
				"public_key_pem":          "-----BEGIN PUBLIC KEY-----\nMIIBIjAN...\n-----END PUBLIC KEY-----",
				"public_key_openssh":      "ssh-rsa AAAAB3NzaC1yc2E...",
			},
			mustPopulate: []string{
				"key_id", "versionless_id", "key_name", "version",
				"resource_id", "resource_versionless_id",
				"public_key_pem", "public_key_openssh",
			},
		},
		{
			// AzureApplicationGateway: the name-keyed maps are the
			// composition seams -- backend_address_pool_ids is what NIC
			// ip_configurations and scale-set network profiles join;
			// frontend_ip_configuration_ids chains frontends; a private
			// frontend's address is what internal DNS records point at.
			name: "AzureApplicationGateway",
			kind: cloudresourcekind.CloudResourceKind_AzureApplicationGateway,
			rawOutputs: map[string]interface{}{
				"application_gateway_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Network/applicationGateways/web-gateway",
				"application_gateway_name": "web-gateway",
				"backend_address_pool_ids": map[string]interface{}{
					"web": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Network/applicationGateways/web-gateway/backendAddressPools/web",
				},
				"frontend_ip_configuration_ids": map[string]interface{}{
					"public": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Network/applicationGateways/web-gateway/frontendIPConfigurations/public",
				},
				"private_ip_address":   "10.0.2.10",
				"private_ip_addresses": []interface{}{"10.0.2.10"},
			},
			mustPopulate: []string{
				"application_gateway_id", "application_gateway_name",
				"backend_address_pool_ids", "frontend_ip_configuration_ids",
				"private_ip_address", "private_ip_addresses",
			},
		},
		{
			// AzureWebApplicationFirewallPolicy: policy_id is the seam
			// Application Gateways attach the policy through -- gateway-wide,
			// per listener, and per URL path rule.
			name: "AzureWebApplicationFirewallPolicy",
			kind: cloudresourcekind.CloudResourceKind_AzureWebApplicationFirewallPolicy,
			rawOutputs: map[string]interface{}{
				"policy_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/applicationGatewayWebApplicationFirewallPolicies/org-waf-baseline",
				"policy_name": "org-waf-baseline",
			},
			mustPopulate: []string{
				"policy_id", "policy_name",
			},
		},
		{
			// AzurePostgresqlFlexibleServer: fqdn + administrator_login are
			// what applications build connection strings from; server_id is
			// the seam private endpoints and replica/restore servers
			// (source_server_id) reference; database_ids is the name-keyed
			// map seam for per-database references.
			name: "AzurePostgresqlFlexibleServer",
			kind: cloudresourcekind.CloudResourceKind_AzurePostgresqlFlexibleServer,
			rawOutputs: map[string]interface{}{
				"server_id":           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/data-rg/providers/Microsoft.DBforPostgreSQL/flexibleServers/orders-pg",
				"server_name":         "orders-pg",
				"fqdn":                "orders-pg.postgres.database.azure.com",
				"administrator_login": "pgadmin",
				"database_ids": map[string]interface{}{
					"orders": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/data-rg/providers/Microsoft.DBforPostgreSQL/flexibleServers/orders-pg/databases/orders",
				},
				"identity_principal_id": "44444444-4444-4444-4444-444444444444",
			},
			mustPopulate: []string{
				"server_id", "server_name", "fqdn", "administrator_login",
				"database_ids", "identity_principal_id",
			},
		},
		{
			// AzureMysqlFlexibleServer: fqdn + administrator_login are what
			// applications build connection strings from; server_id is the
			// seam private endpoints and replica/restore servers
			// (source_server_id) reference; database_ids is the name-keyed
			// map seam; replica_capacity sizes replica topologies.
			name: "AzureMysqlFlexibleServer",
			kind: cloudresourcekind.CloudResourceKind_AzureMysqlFlexibleServer,
			rawOutputs: map[string]interface{}{
				"server_id":           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/data-rg/providers/Microsoft.DBforMySQL/flexibleServers/orders-mysql",
				"server_name":         "orders-mysql",
				"fqdn":                "orders-mysql.mysql.database.azure.com",
				"administrator_login": "mysqladmin",
				"database_ids": map[string]interface{}{
					"orders": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/data-rg/providers/Microsoft.DBforMySQL/flexibleServers/orders-mysql/databases/orders",
				},
				"replica_capacity": 10,
			},
			mustPopulate: []string{
				"server_id", "server_name", "fqdn", "administrator_login",
				"database_ids", "replica_capacity",
			},
		},
		{
			// AzureMssqlServer: server_id is the parent seam
			// AzureMssqlDatabase and AzureMssqlElasticPool reference (and
			// AzurePrivateEndpoint's connection target); fqdn +
			// administrator_login build connection strings.
			name: "AzureMssqlServer",
			kind: cloudresourcekind.CloudResourceKind_AzureMssqlServer,
			rawOutputs: map[string]interface{}{
				"server_id":             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/data-rg/providers/Microsoft.Sql/servers/orders-sql",
				"server_name":           "orders-sql",
				"fqdn":                  "orders-sql.database.windows.net",
				"administrator_login":   "sqladmin",
				"identity_principal_id": "44444444-4444-4444-4444-444444444444",
			},
			mustPopulate: []string{
				"server_id", "server_name", "fqdn", "administrator_login",
				"identity_principal_id",
			},
		},
		{
			// AzureMssqlDatabase: database_id is the seam
			// copy/secondary/restore databases reference
			// (creation_source_database_id).
			name: "AzureMssqlDatabase",
			kind: cloudresourcekind.CloudResourceKind_AzureMssqlDatabase,
			rawOutputs: map[string]interface{}{
				"database_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/data-rg/providers/Microsoft.Sql/servers/orders-sql/databases/orders",
				"database_name": "orders",
			},
			mustPopulate: []string{
				"database_id", "database_name",
			},
		},
		{
			// AzureMssqlElasticPool: elastic_pool_id is the seam pooled
			// databases attach through (AzureMssqlDatabase.elastic_pool_id).
			name: "AzureMssqlElasticPool",
			kind: cloudresourcekind.CloudResourceKind_AzureMssqlElasticPool,
			rawOutputs: map[string]interface{}{
				"elastic_pool_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/data-rg/providers/Microsoft.Sql/servers/orders-sql/elasticPools/tenant-pool",
				"elastic_pool_name": "tenant-pool",
			},
			mustPopulate: []string{
				"elastic_pool_id", "elastic_pool_name",
			},
		},
		{
			// AzureStorageAccount: storage_account_id is the parent seam
			// AzureStorageContainer references (and data-plane role
			// assignments scope to); the name + primary_access_key pair is
			// what Function App / Linux Web App storage bindings consume;
			// the endpoints are what applications and CDN origins connect
			// to.
			name: "AzureStorageAccount",
			kind: cloudresourcekind.CloudResourceKind_AzureStorageAccount,
			rawOutputs: map[string]interface{}{
				"storage_account_id":               "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Storage/storageAccounts/plantonappstorage",
				"storage_account_name":             "plantonappstorage",
				"resource_group_name":              "app-rg",
				"primary_blob_endpoint":            "https://plantonappstorage.blob.core.windows.net/",
				"primary_blob_host":                "plantonappstorage.blob.core.windows.net",
				"primary_queue_endpoint":           "https://plantonappstorage.queue.core.windows.net/",
				"primary_table_endpoint":           "https://plantonappstorage.table.core.windows.net/",
				"primary_file_endpoint":            "https://plantonappstorage.file.core.windows.net/",
				"primary_dfs_endpoint":             "https://plantonappstorage.dfs.core.windows.net/",
				"primary_web_endpoint":             "https://plantonappstorage.z13.web.core.windows.net/",
				"primary_web_host":                 "plantonappstorage.z13.web.core.windows.net",
				"secondary_blob_endpoint":          "https://plantonappstorage-secondary.blob.core.windows.net/",
				"primary_access_key":               "base64keymaterial==",
				"secondary_access_key":             "base64keymaterial2==",
				"primary_connection_string":        "DefaultEndpointsProtocol=https;AccountName=plantonappstorage;AccountKey=base64keymaterial==;EndpointSuffix=core.windows.net",
				"secondary_connection_string":      "DefaultEndpointsProtocol=https;AccountName=plantonappstorage;AccountKey=base64keymaterial2==;EndpointSuffix=core.windows.net",
				"primary_blob_connection_string":   "DefaultEndpointsProtocol=https;BlobEndpoint=https://plantonappstorage.blob.core.windows.net/;AccountName=plantonappstorage;AccountKey=base64keymaterial==",
				"secondary_blob_connection_string": "DefaultEndpointsProtocol=https;BlobEndpoint=https://plantonappstorage-secondary.blob.core.windows.net/;AccountName=plantonappstorage;AccountKey=base64keymaterial2==",
				"identity_principal_id":            "44444444-4444-4444-4444-444444444444",
			},
			mustPopulate: []string{
				"storage_account_id", "storage_account_name", "resource_group_name",
				"primary_blob_endpoint", "primary_blob_host", "primary_queue_endpoint",
				"primary_table_endpoint", "primary_file_endpoint", "primary_dfs_endpoint",
				"primary_web_endpoint", "primary_web_host", "secondary_blob_endpoint",
				"primary_access_key", "secondary_access_key", "primary_connection_string",
				"secondary_connection_string", "primary_blob_connection_string",
				"secondary_blob_connection_string", "identity_principal_id",
			},
		},
		{
			// AzureStorageContainer: container_id is the scope data-plane
			// role assignments target for container-level access; the
			// account/container name pair is what SDK clients and function
			// bindings consume.
			name: "AzureStorageContainer",
			kind: cloudresourcekind.CloudResourceKind_AzureStorageContainer,
			rawOutputs: map[string]interface{}{
				"container_id":         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Storage/storageAccounts/plantonappstorage/blobServices/default/containers/uploads",
				"container_name":       "uploads",
				"storage_account_name": "plantonappstorage",
			},
			mustPopulate: []string{
				"container_id", "container_name", "storage_account_name",
			},
		},
		{
			// AzureStorageShare: share_id is the management identity;
			// rbac_scope_id is the DIFFERENT segment Azure Files data-plane
			// role assignments scope to; the account/share name pair is
			// what mount commands and CSI volume definitions consume.
			name: "AzureStorageShare",
			kind: cloudresourcekind.CloudResourceKind_AzureStorageShare,
			rawOutputs: map[string]interface{}{
				"share_id":             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Storage/storageAccounts/plantonappstorage/fileServices/default/shares/team-files",
				"rbac_scope_id":        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Storage/storageAccounts/plantonappstorage/fileServices/default/fileshares/team-files",
				"share_name":           "team-files",
				"storage_account_name": "plantonappstorage",
			},
			mustPopulate: []string{
				"share_id", "rbac_scope_id", "share_name", "storage_account_name",
			},
		},
		{
			// AzureStorageQueue: queue_id is the scope data-plane role
			// assignments target for queue-level access; the account/queue
			// name pair is what SDK clients and Functions queue triggers
			// consume.
			name: "AzureStorageQueue",
			kind: cloudresourcekind.CloudResourceKind_AzureStorageQueue,
			rawOutputs: map[string]interface{}{
				"queue_id":             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Storage/storageAccounts/plantonappstorage/queueServices/default/queues/work-items",
				"queue_name":           "work-items",
				"storage_account_name": "plantonappstorage",
			},
			mustPopulate: []string{
				"queue_id", "queue_name", "storage_account_name",
			},
		},
		{
			// AzureStorageTable: table_id carries the resource-manager id
			// from BOTH engines (the addressing parity exception never
			// touches outputs); the account/table name pair is what SDK
			// clients and Functions table bindings consume.
			name: "AzureStorageTable",
			kind: cloudresourcekind.CloudResourceKind_AzureStorageTable,
			rawOutputs: map[string]interface{}{
				"table_id":             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Storage/storageAccounts/plantonappstorage/tableServices/default/tables/AppEntities",
				"table_name":           "AppEntities",
				"storage_account_name": "plantonappstorage",
			},
			mustPopulate: []string{
				"table_id", "table_name", "storage_account_name",
			},
		},
		{
			// AzureStorageEncryptionScope: encryption_scope_name is the
			// seam containers (default_encryption_scope) and ADLS
			// filesystems reference for sub-account key isolation.
			name: "AzureStorageEncryptionScope",
			kind: cloudresourcekind.CloudResourceKind_AzureStorageEncryptionScope,
			rawOutputs: map[string]interface{}{
				"encryption_scope_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Storage/storageAccounts/plantonappstorage/encryptionScopes/tenant42scope",
				"encryption_scope_name": "tenant42scope",
				"storage_account_name":  "plantonappstorage",
			},
			mustPopulate: []string{
				"encryption_scope_id", "encryption_scope_name", "storage_account_name",
			},
		},
		{
			// AzureStorageDataLakeGen2Filesystem: filesystem_id is the ARM
			// container-proxy ID data-plane role assignments scope to
			// (ADLS filesystems surface in ARM as blob containers).
			name: "AzureStorageDataLakeGen2Filesystem",
			kind: cloudresourcekind.CloudResourceKind_AzureStorageDataLakeGen2Filesystem,
			rawOutputs: map[string]interface{}{
				"filesystem_id":        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/lake-rg/providers/Microsoft.Storage/storageAccounts/plantonlake/blobServices/default/containers/raw-zone",
				"filesystem_name":      "raw-zone",
				"storage_account_name": "plantonlake",
			},
			mustPopulate: []string{
				"filesystem_id", "filesystem_name", "storage_account_name",
			},
		},
		{
			// AzureStorageLocalUser: sftp_username is the composed login
			// clients connect with; sid and password are the
			// secret-bearing credential outputs.
			name: "AzureStorageLocalUser",
			kind: cloudresourcekind.CloudResourceKind_AzureStorageLocalUser,
			rawOutputs: map[string]interface{}{
				"local_user_id":        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/exchange-rg/providers/Microsoft.Storage/storageAccounts/plantonsftp/localUsers/partner01",
				"user_name":            "partner01",
				"sftp_username":        "plantonsftp.partner01",
				"sid":                  "S-1-2-0-3895023191-1105595861-2277418014-1116",
				"password":             "generated-once-by-azure",
				"storage_account_name": "plantonsftp",
			},
			mustPopulate: []string{
				"local_user_id", "user_name", "sftp_username", "sid", "password", "storage_account_name",
			},
		},
		{
			// AzureStorageObjectReplication: one logical policy
			// materialized on BOTH accounts under one GUID -- two ARM IDs
			// plus the shared policy_id monitoring keys on.
			name: "AzureStorageObjectReplication",
			kind: cloudresourcekind.CloudResourceKind_AzureStorageObjectReplication,
			rawOutputs: map[string]interface{}{
				"source_object_replication_id":      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/dr-rg/providers/Microsoft.Storage/storageAccounts/plantonorsrc/objectReplicationPolicies/6a2f5b7e-1c3d-4e5f-8a9b-0c1d2e3f4a5b",
				"destination_object_replication_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/dr-rg/providers/Microsoft.Storage/storageAccounts/plantonordst/objectReplicationPolicies/6a2f5b7e-1c3d-4e5f-8a9b-0c1d2e3f4a5b",
				"policy_id":                         "6a2f5b7e-1c3d-4e5f-8a9b-0c1d2e3f4a5b",
			},
			mustPopulate: []string{
				"source_object_replication_id", "destination_object_replication_id", "policy_id",
			},
		},
		{
			// AzureKeyVaultCertificate: the secret face
			// (versionless_secret_id) is the seam TLS terminators
			// (Application Gateway) consume so renewals propagate; the
			// thumbprint serves fingerprint-pinning integrations.
			name: "AzureKeyVaultCertificate",
			kind: cloudresourcekind.CloudResourceKind_AzureKeyVaultCertificate,
			rawOutputs: map[string]interface{}{
				"certificate_id":                  "https://platform-kv.vault.azure.net/certificates/internal-tls/fed321cba654",
				"versionless_id":                  "https://platform-kv.vault.azure.net/certificates/internal-tls",
				"secret_id":                       "https://platform-kv.vault.azure.net/secrets/internal-tls/fed321cba654",
				"versionless_secret_id":           "https://platform-kv.vault.azure.net/secrets/internal-tls",
				"certificate_name":                "internal-tls",
				"version":                         "fed321cba654",
				"thumbprint":                      "9F3C4E2A1B0D8765F4E3D2C1B0A99887E6D5C4B3",
				"resource_manager_id":             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/security-rg/providers/Microsoft.KeyVault/vaults/platform-kv/certificates/internal-tls/versions/fed321cba654",
				"resource_manager_versionless_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/security-rg/providers/Microsoft.KeyVault/vaults/platform-kv/certificates/internal-tls",
			},
			mustPopulate: []string{
				"certificate_id", "versionless_id", "secret_id",
				"versionless_secret_id", "certificate_name", "version",
				"thumbprint", "resource_manager_id",
				"resource_manager_versionless_id",
			},
		},
		{
			// AzureRedisCache: redis_cache_id is what the linked-server,
			// access-policy, and private-endpoint kinds reference; region
			// is the linked-server location seam; both key faces stay live
			// so clients rotate with zero downtime.
			name: "AzureRedisCache",
			kind: cloudresourcekind.CloudResourceKind_AzureRedisCache,
			rawOutputs: map[string]interface{}{
				"redis_cache_id":              "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cache/redis/app-cache",
				"redis_cache_name":            "app-cache",
				"region":                      "eastus",
				"resource_group_name":         "app-rg",
				"hostname":                    "app-cache.redis.cache.windows.net",
				"port":                        6379,
				"ssl_port":                    6380,
				"primary_access_key":          "primary-key-value",
				"secondary_access_key":        "secondary-key-value",
				"primary_connection_string":   "app-cache.redis.cache.windows.net:6380,password=primary-key-value,ssl=True,abortConnect=False",
				"secondary_connection_string": "app-cache.redis.cache.windows.net:6380,password=secondary-key-value,ssl=True,abortConnect=False",
				"identity_principal_id":       "11111111-2222-3333-4444-555555555555",
			},
			mustPopulate: []string{
				"redis_cache_id", "redis_cache_name", "region",
				"resource_group_name", "hostname", "port", "ssl_port",
				"primary_access_key", "secondary_access_key",
				"primary_connection_string", "secondary_connection_string",
				"identity_principal_id",
			},
		},
		{
			// AzureRedisLinkedServer: the geo hostname follows the CURRENT
			// primary across failovers -- the stable endpoint applications
			// point at instead of either cache's own hostname.
			name: "AzureRedisLinkedServer",
			kind: cloudresourcekind.CloudResourceKind_AzureRedisLinkedServer,
			rawOutputs: map[string]interface{}{
				"linked_server_id":                 "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-east/providers/Microsoft.Cache/redis/app-cache-east/linkedServers/app-cache-west",
				"linked_server_name":               "app-cache-west",
				"geo_replicated_primary_host_name": "app-cache-east.geo.redis.cache.windows.net",
			},
			mustPopulate: []string{
				"linked_server_id", "linked_server_name",
				"geo_replicated_primary_host_name",
			},
		},
		{
			// AzureRedisCacheAccessPolicy: access_policy_name is the seam
			// assignments reference to grant the policy to an identity.
			name: "AzureRedisCacheAccessPolicy",
			kind: cloudresourcekind.CloudResourceKind_AzureRedisCacheAccessPolicy,
			rawOutputs: map[string]interface{}{
				"access_policy_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cache/redis/app-cache/accessPolicies/app-read-only",
				"access_policy_name": "app-read-only",
			},
			mustPopulate: []string{
				"access_policy_id", "access_policy_name",
			},
		},
		{
			// AzureManagedRedis: managed_redis_id is what the
			// geo-replication and access-policy-assignment kinds
			// reference; database_id is the grant/link scope; the keys
			// populate only while access-keys authentication is enabled
			// (keyless is the default).
			name: "AzureManagedRedis",
			kind: cloudresourcekind.CloudResourceKind_AzureManagedRedis,
			rawOutputs: map[string]interface{}{
				"managed_redis_id":      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cache/redisEnterprise/app-cache",
				"managed_redis_name":    "app-cache",
				"region":                "eastus",
				"resource_group_name":   "app-rg",
				"hostname":              "app-cache.eastus.redis.azure.net",
				"database_id":           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cache/redisEnterprise/app-cache/databases/default",
				"port":                  10000,
				"primary_access_key":    "primary-key-value",
				"secondary_access_key":  "secondary-key-value",
				"identity_principal_id": "11111111-2222-3333-4444-555555555555",
			},
			mustPopulate: []string{
				"managed_redis_id", "managed_redis_name", "region",
				"resource_group_name", "hostname", "database_id", "port",
				"primary_access_key", "secondary_access_key",
				"identity_principal_id",
			},
		},
		{
			// AzureServicePlan: service_plan_id is what the web/function
			// app kinds reference; kind and reserved are Azure-computed
			// attributes read back after creation.
			name: "AzureServicePlan",
			kind: cloudresourcekind.CloudResourceKind_AzureServicePlan,
			rawOutputs: map[string]interface{}{
				"service_plan_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Web/serverFarms/app-plan",
				"service_plan_name": "app-plan",
				"os_type":           "Linux",
				"sku_name":          "P1v3",
				"kind":              "linux",
				"reserved":          true,
			},
			mustPopulate: []string{
				"service_plan_id", "service_plan_name", "os_type",
				"sku_name", "kind", "reserved",
			},
		},
		{
			// AzureLinuxWebApp: default_hostname is the app's endpoint;
			// the outbound IP sets arrive as real lists from both
			// engines; the site credential populates while basic-auth
			// publishing is enabled.
			name: "AzureLinuxWebApp",
			kind: cloudresourcekind.CloudResourceKind_AzureLinuxWebApp,
			rawOutputs: map[string]interface{}{
				"web_app_id":                     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Web/sites/app-web",
				"default_hostname":               "app-web.azurewebsites.net",
				"outbound_ip_addresses":          []interface{}{"20.1.2.3", "20.1.2.4"},
				"possible_outbound_ip_addresses": []interface{}{"20.1.2.3", "20.1.2.4", "20.1.2.5"},
				"identity_principal_id":          "11111111-2222-3333-4444-555555555555",
				"identity_tenant_id":             "99999999-8888-7777-6666-555555555555",
				"custom_domain_verification_id":  "ABCD1234",
				"kind":                           "app,linux",
				"hosting_environment_id":         "",
				"site_credential_name":           "$app-web",
				"site_credential_password":       "publish-password",
			},
			mustPopulate: []string{
				"web_app_id", "default_hostname", "outbound_ip_addresses",
				"possible_outbound_ip_addresses", "identity_principal_id",
				"identity_tenant_id", "custom_domain_verification_id",
				"kind", "site_credential_name", "site_credential_password",
			},
		},
		{
			// AzureFunctionApp: default_hostname serves HTTP triggers;
			// the outbound IP sets arrive as real lists from both
			// engines; the site credential populates while basic-auth
			// publishing is enabled.
			name: "AzureFunctionApp",
			kind: cloudresourcekind.CloudResourceKind_AzureFunctionApp,
			rawOutputs: map[string]interface{}{
				"function_app_id":                "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Web/sites/app-fn",
				"default_hostname":               "app-fn.azurewebsites.net",
				"outbound_ip_addresses":          []interface{}{"20.1.2.3", "20.1.2.4"},
				"possible_outbound_ip_addresses": []interface{}{"20.1.2.3", "20.1.2.4", "20.1.2.5"},
				"identity_principal_id":          "11111111-2222-3333-4444-555555555555",
				"identity_tenant_id":             "99999999-8888-7777-6666-555555555555",
				"custom_domain_verification_id":  "ABCD1234",
				"kind":                           "functionapp,linux",
				"hosting_environment_id":         "",
				"site_credential_name":           "$app-fn",
				"site_credential_password":       "publish-password",
			},
			mustPopulate: []string{
				"function_app_id", "default_hostname",
				"outbound_ip_addresses", "possible_outbound_ip_addresses",
				"identity_principal_id", "identity_tenant_id",
				"custom_domain_verification_id", "kind",
				"site_credential_name", "site_credential_password",
			},
		},
		{
			// AzureManagedRedisGeoReplication: the group has no ARM
			// object of its own -- its resource ID is the managing
			// cluster's ARM ID.
			name: "AzureManagedRedisGeoReplication",
			kind: cloudresourcekind.CloudResourceKind_AzureManagedRedisGeoReplication,
			rawOutputs: map[string]interface{}{
				"geo_replication_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-east/providers/Microsoft.Cache/redisEnterprise/app-cache-east",
			},
			mustPopulate: []string{
				"geo_replication_id",
			},
		},
		{
			// AzureManagedRedisAccessPolicyAssignment: Azure names the
			// assignment after the granted object ID, so the name equals
			// the principal's GUID.
			name: "AzureManagedRedisAccessPolicyAssignment",
			kind: cloudresourcekind.CloudResourceKind_AzureManagedRedisAccessPolicyAssignment,
			rawOutputs: map[string]interface{}{
				"access_policy_assignment_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cache/redisEnterprise/app-cache/databases/default/accessPolicyAssignments/11111111-2222-3333-4444-555555555555",
				"access_policy_assignment_name": "11111111-2222-3333-4444-555555555555",
			},
			mustPopulate: []string{
				"access_policy_assignment_id", "access_policy_assignment_name",
			},
		},
		{
			// AzureRedisCacheAccessPolicyAssignment: the grant half of the
			// keyless cache story -- id and name identify the grant for
			// audits and teardown.
			name: "AzureRedisCacheAccessPolicyAssignment",
			kind: cloudresourcekind.CloudResourceKind_AzureRedisCacheAccessPolicyAssignment,
			rawOutputs: map[string]interface{}{
				"access_policy_assignment_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cache/redis/app-cache/accessPolicyAssignments/app-identity-data-reader",
				"access_policy_assignment_name": "app-identity-data-reader",
			},
			mustPopulate: []string{
				"access_policy_assignment_id", "access_policy_assignment_name",
			},
		},
		{
			// AzureCosmosdbAccount: cosmosdb_account_id is what the SQL/Mongo
			// database kinds and private endpoints reference; the keys and
			// ready-made connection strings are the credential surface; the
			// per-region endpoint lists are repeated string outputs.
			name: "AzureCosmosdbAccount",
			kind: cloudresourcekind.CloudResourceKind_AzureCosmosdbAccount,
			rawOutputs: map[string]interface{}{
				"cosmosdb_account_id":                          "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.DocumentDB/databaseAccounts/app-cosmos",
				"cosmosdb_account_name":                        "app-cosmos",
				"endpoint":                                     "https://app-cosmos.documents.azure.com:443/",
				"read_endpoints":                               []interface{}{"https://app-cosmos-eastus.documents.azure.com:443/", "https://app-cosmos-westus2.documents.azure.com:443/"},
				"write_endpoints":                              []interface{}{"https://app-cosmos-eastus.documents.azure.com:443/"},
				"primary_key":                                  "primary-key-value",
				"secondary_key":                                "secondary-key-value",
				"primary_readonly_key":                         "primary-readonly-key-value",
				"secondary_readonly_key":                       "secondary-readonly-key-value",
				"primary_sql_connection_string":                "AccountEndpoint=https://app-cosmos.documents.azure.com:443/;AccountKey=primary-key-value;",
				"secondary_sql_connection_string":              "AccountEndpoint=https://app-cosmos.documents.azure.com:443/;AccountKey=secondary-key-value;",
				"primary_readonly_sql_connection_string":       "AccountEndpoint=https://app-cosmos.documents.azure.com:443/;AccountKey=primary-readonly-key-value;",
				"secondary_readonly_sql_connection_string":     "AccountEndpoint=https://app-cosmos.documents.azure.com:443/;AccountKey=secondary-readonly-key-value;",
				"primary_mongodb_connection_string":            "mongodb://app-cosmos:primary-key-value@app-cosmos.mongo.cosmos.azure.com:10255/?ssl=true",
				"secondary_mongodb_connection_string":          "mongodb://app-cosmos:secondary-key-value@app-cosmos.mongo.cosmos.azure.com:10255/?ssl=true",
				"primary_readonly_mongodb_connection_string":   "mongodb://app-cosmos:primary-readonly-key-value@app-cosmos.mongo.cosmos.azure.com:10255/?ssl=true",
				"secondary_readonly_mongodb_connection_string": "mongodb://app-cosmos:secondary-readonly-key-value@app-cosmos.mongo.cosmos.azure.com:10255/?ssl=true",
				"identity_principal_id":                        "11111111-2222-3333-4444-555555555555",
			},
			mustPopulate: []string{
				"cosmosdb_account_id", "cosmosdb_account_name", "endpoint",
				"read_endpoints", "write_endpoints",
				"primary_key", "secondary_key",
				"primary_readonly_key", "secondary_readonly_key",
				"primary_sql_connection_string", "secondary_sql_connection_string",
				"primary_readonly_sql_connection_string", "secondary_readonly_sql_connection_string",
				"primary_mongodb_connection_string", "secondary_mongodb_connection_string",
				"primary_readonly_mongodb_connection_string", "secondary_readonly_mongodb_connection_string",
				"identity_principal_id",
			},
		},
		{
			// AzureCosmosdbSqlDatabase: sql_database_id is the seam
			// containers (AzureCosmosdbSqlContainer.sql_database_id)
			// reference; the account/database name pair is what SDK calls
			// consume inside the account's connection.
			name: "AzureCosmosdbSqlDatabase",
			kind: cloudresourcekind.CloudResourceKind_AzureCosmosdbSqlDatabase,
			rawOutputs: map[string]interface{}{
				"sql_database_id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.DocumentDB/databaseAccounts/app-cosmos/sqlDatabases/app-data",
				"sql_database_name":     "app-data",
				"cosmosdb_account_name": "app-cosmos",
			},
			mustPopulate: []string{
				"sql_database_id", "sql_database_name", "cosmosdb_account_name",
			},
		},
		{
			// AzureCosmosdbSqlContainer: sql_container_id is the
			// management identity and the container-level data-plane RBAC
			// scope; the name triple addresses the container inside the
			// account's connection.
			name: "AzureCosmosdbSqlContainer",
			kind: cloudresourcekind.CloudResourceKind_AzureCosmosdbSqlContainer,
			rawOutputs: map[string]interface{}{
				"sql_container_id":      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.DocumentDB/databaseAccounts/app-cosmos/sqlDatabases/app-data/containers/orders",
				"sql_container_name":    "orders",
				"sql_database_name":     "app-data",
				"cosmosdb_account_name": "app-cosmos",
			},
			mustPopulate: []string{
				"sql_container_id", "sql_container_name",
				"sql_database_name", "cosmosdb_account_name",
			},
		},
		{
			// AzureCosmosdbMongoDatabase: mongo_database_id is the seam
			// collections (AzureCosmosdbMongoCollection.mongo_database_id)
			// reference.
			name: "AzureCosmosdbMongoDatabase",
			kind: cloudresourcekind.CloudResourceKind_AzureCosmosdbMongoDatabase,
			rawOutputs: map[string]interface{}{
				"mongo_database_id":     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.DocumentDB/databaseAccounts/app-cosmos-mongo/mongodbDatabases/app-data",
				"mongo_database_name":   "app-data",
				"cosmosdb_account_name": "app-cosmos-mongo",
			},
			mustPopulate: []string{
				"mongo_database_id", "mongo_database_name", "cosmosdb_account_name",
			},
		},
		{
			// AzureCosmosdbMongoCollection: the name triple addresses the
			// collection inside the account's Mongo connection string.
			name: "AzureCosmosdbMongoCollection",
			kind: cloudresourcekind.CloudResourceKind_AzureCosmosdbMongoCollection,
			rawOutputs: map[string]interface{}{
				"mongo_collection_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.DocumentDB/databaseAccounts/app-cosmos-mongo/mongodbDatabases/app-data/collections/events",
				"mongo_collection_name": "events",
				"mongo_database_name":   "app-data",
				"cosmosdb_account_name": "app-cosmos-mongo",
			},
			mustPopulate: []string{
				"mongo_collection_id", "mongo_collection_name",
				"mongo_database_name", "cosmosdb_account_name",
			},
		},
		{
			// AzureCosmosdbSqlRoleDefinition: role_definition_id is the
			// fully-scoped ARM id an AzureCosmosdbSqlRoleAssignment's
			// role_definition_id field consumes with zero translation.
			name: "AzureCosmosdbSqlRoleDefinition",
			kind: cloudresourcekind.CloudResourceKind_AzureCosmosdbSqlRoleDefinition,
			rawOutputs: map[string]interface{}{
				"role_definition_id":    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.DocumentDB/databaseAccounts/app-cosmos/sqlRoleDefinitions/9b7f3f6a-2f0e-4b9a-8f0d-2f6a8f0d2f6a",
				"role_definition_guid":  "9b7f3f6a-2f0e-4b9a-8f0d-2f6a8f0d2f6a",
				"role_name":             "app-reader",
				"cosmosdb_account_name": "app-cosmos",
			},
			mustPopulate: []string{
				"role_definition_id", "role_definition_guid",
				"role_name", "cosmosdb_account_name",
			},
		},
		{
			// AzureCosmosdbSqlRoleAssignment: the grant record's ARM
			// identity, exported for audit trails and cross-references.
			name: "AzureCosmosdbSqlRoleAssignment",
			kind: cloudresourcekind.CloudResourceKind_AzureCosmosdbSqlRoleAssignment,
			rawOutputs: map[string]interface{}{
				"role_assignment_id":    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.DocumentDB/databaseAccounts/app-cosmos/sqlRoleAssignments/7c1de3f8-5a4b-4c2d-9e8f-1a2b3c4d5e6f",
				"role_assignment_guid":  "7c1de3f8-5a4b-4c2d-9e8f-1a2b3c4d5e6f",
				"cosmosdb_account_name": "app-cosmos",
			},
			mustPopulate: []string{
				"role_assignment_id", "role_assignment_guid",
				"cosmosdb_account_name",
			},
		},
		{
			// AzureFrontDoorProfile: profile_id is the parent seam every
			// Front Door delivery kind (endpoint, origin group) references;
			// identity_principal_id is the Key Vault grant target for
			// bring-your-own TLS certificates.
			name: "AzureFrontDoorProfile",
			kind: cloudresourcekind.CloudResourceKind_AzureFrontDoorProfile,
			rawOutputs: map[string]interface{}{
				"profile_id":            "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cdn/profiles/app-fd",
				"profile_name":          "app-fd",
				"resource_guid":         "11111111-2222-3333-4444-555555555555",
				"identity_principal_id": "66666666-7777-8888-9999-000000000000",
			},
			mustPopulate: []string{
				"profile_id", "profile_name", "resource_guid", "identity_principal_id",
			},
		},
		{
			// AzureFrontDoorEndpoint: endpoint_id is the route's parent
			// seam; host_name is the generated *.azurefd.net hostname DNS
			// records CNAME onto.
			name: "AzureFrontDoorEndpoint",
			kind: cloudresourcekind.CloudResourceKind_AzureFrontDoorEndpoint,
			rawOutputs: map[string]interface{}{
				"endpoint_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cdn/profiles/app-fd/afdEndpoints/web",
				"endpoint_name": "web",
				"host_name":     "web-abc123.z01.azurefd.net",
			},
			mustPopulate: []string{
				"endpoint_id", "endpoint_name", "host_name",
			},
		},
		{
			// AzureFrontDoorOriginGroup: origin_group_id is what origins
			// reference as parent and routes reference as destination.
			name: "AzureFrontDoorOriginGroup",
			kind: cloudresourcekind.CloudResourceKind_AzureFrontDoorOriginGroup,
			rawOutputs: map[string]interface{}{
				"origin_group_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cdn/profiles/app-fd/originGroups/api-backends",
				"origin_group_name": "api-backends",
			},
			mustPopulate: []string{
				"origin_group_id", "origin_group_name",
			},
		},
		{
			// AzureFrontDoorOrigin: origin_id is what routes list in
			// origin_ids to sequence deployment after the backends exist.
			name: "AzureFrontDoorOrigin",
			kind: cloudresourcekind.CloudResourceKind_AzureFrontDoorOrigin,
			rawOutputs: map[string]interface{}{
				"origin_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cdn/profiles/app-fd/originGroups/api-backends/origins/primary",
				"origin_name": "primary",
			},
			mustPopulate: []string{
				"origin_id", "origin_name",
			},
		},
		{
			// AzureFrontDoorRoute: the traffic-serving edge of the Front
			// Door graph; no hostname output on purpose (it lives on the
			// endpoint).
			name: "AzureFrontDoorRoute",
			kind: cloudresourcekind.CloudResourceKind_AzureFrontDoorRoute,
			rawOutputs: map[string]interface{}{
				"route_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cdn/profiles/app-fd/afdEndpoints/web/routes/default",
				"route_name": "default",
			},
			mustPopulate: []string{
				"route_id", "route_name",
			},
		},
		{
			// AzureFrontDoorRuleSet: rule_set_id is what routes reference
			// in rule_set_ids to attach the delivery policy; the folded
			// rules export no ids on purpose (nothing references a rule).
			name: "AzureFrontDoorRuleSet",
			kind: cloudresourcekind.CloudResourceKind_AzureFrontDoorRuleSet,
			rawOutputs: map[string]interface{}{
				"rule_set_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cdn/profiles/app-fd/ruleSets/deliverypolicy",
				"rule_set_name": "deliverypolicy",
			},
			mustPopulate: []string{
				"rule_set_id", "rule_set_name",
			},
		},
		{
			// AzureFrontDoorCustomDomain: custom_domain_id is the route
			// attach seam; validation_token is the DNS TXT challenge the
			// operator publishes at _dnsauth.<host_name>.
			name: "AzureFrontDoorCustomDomain",
			kind: cloudresourcekind.CloudResourceKind_AzureFrontDoorCustomDomain,
			rawOutputs: map[string]interface{}{
				"custom_domain_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cdn/profiles/app-fd/customDomains/www-example-com",
				"host_name":        "www.example.com",
				"validation_token": "_zy4mhfvswrqzeqmnyv6gjr26xk1mbrv",
				"expiration_date":  "2026-07-16T00:00:00.000Z",
			},
			mustPopulate: []string{
				"custom_domain_id", "host_name", "validation_token", "expiration_date",
			},
		},
		{
			// AzureFrontDoorSecret: secret_id is the custom domain's
			// tls.secret_id seam; the SANs are read back from the wrapped
			// Key Vault certificate.
			name: "AzureFrontDoorSecret",
			kind: cloudresourcekind.CloudResourceKind_AzureFrontDoorSecret,
			rawOutputs: map[string]interface{}{
				"secret_id":                 "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cdn/profiles/app-fd/secrets/wildcard-example-com",
				"secret_name":               "wildcard-example-com",
				"subject_alternative_names": []interface{}{"*.example.com", "example.com"},
			},
			mustPopulate: []string{
				"secret_id", "secret_name", "subject_alternative_names",
			},
		},
		{
			// AzureFrontDoorFirewallPolicy: firewall_policy_id is what the
			// security policy references in firewall_policy_id to attach
			// the WAF to a profile's domains.
			name: "AzureFrontDoorFirewallPolicy",
			kind: cloudresourcekind.CloudResourceKind_AzureFrontDoorFirewallPolicy,
			rawOutputs: map[string]interface{}{
				"firewall_policy_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/edgewaf",
				"firewall_policy_name": "edgewaf",
			},
			mustPopulate: []string{
				"firewall_policy_id", "firewall_policy_name",
			},
		},
		{
			// AzureFrontDoorSecurityPolicy: the association itself --
			// nothing composes on it; the id serves operational
			// addressing.
			name: "AzureFrontDoorSecurityPolicy",
			kind: cloudresourcekind.CloudResourceKind_AzureFrontDoorSecurityPolicy,
			rawOutputs: map[string]interface{}{
				"security_policy_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Cdn/profiles/app-fd/securityPolicies/edge-waf-attach",
				"security_policy_name": "edge-waf-attach",
			},
			mustPopulate: []string{
				"security_policy_id", "security_policy_name",
			},
		},
		{
			// AzureContainerAppEnvironment: environment_id is what every
			// kind living inside the environment references; the
			// platform-reserved values only populate for VNet-injected
			// environments but the output shape stays constant.
			name: "AzureContainerAppEnvironment",
			kind: cloudresourcekind.CloudResourceKind_AzureContainerAppEnvironment,
			rawOutputs: map[string]interface{}{
				"environment_id":                   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.App/managedEnvironments/app-env",
				"environment_name":                 "app-env",
				"default_domain":                   "app-env.eastus.azurecontainerapps.io",
				"static_ip_address":                "20.1.2.3",
				"platform_reserved_cidr":           "10.0.0.0/24",
				"platform_reserved_dns_ip_address": "10.0.0.2",
				"docker_bridge_cidr":               "172.17.0.1/16",
				"custom_domain_verification_id":    "ABCD1234",
				"identity_principal_id":            "11111111-2222-3333-4444-555555555555",
			},
			mustPopulate: []string{
				"environment_id", "environment_name", "default_domain",
				"static_ip_address", "platform_reserved_cidr",
				"platform_reserved_dns_ip_address", "docker_bridge_cidr",
				"custom_domain_verification_id", "identity_principal_id",
			},
		},
		{
			// AzureContainerApp: ingress_fqdn is the user-facing endpoint;
			// the outbound IPs arrive as a real list from both engines;
			// custom_domain_verification_id is provider-Sensitive (the TF
			// output carries sensitive = true).
			name: "AzureContainerApp",
			kind: cloudresourcekind.CloudResourceKind_AzureContainerApp,
			rawOutputs: map[string]interface{}{
				"container_app_id":              "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.App/containerApps/app-web",
				"container_app_name":            "app-web",
				"latest_revision_name":          "app-web--abc123",
				"latest_revision_fqdn":          "app-web--abc123.app-env.eastus.azurecontainerapps.io",
				"outbound_ip_addresses":         []interface{}{"20.1.2.3", "20.1.2.4"},
				"ingress_fqdn":                  "app-web.app-env.eastus.azurecontainerapps.io",
				"custom_domain_verification_id": "ABCD1234",
				"identity_principal_id":         "11111111-2222-3333-4444-555555555555",
			},
			mustPopulate: []string{
				"container_app_id", "container_app_name",
				"latest_revision_name", "latest_revision_fqdn",
				"outbound_ip_addresses", "ingress_fqdn",
				"custom_domain_verification_id", "identity_principal_id",
			},
		},
		{
			// AzureContainerAppJob: job_id is the handle for starting
			// manual executions; event_stream_endpoint feeds execution
			// monitoring.
			name: "AzureContainerAppJob",
			kind: cloudresourcekind.CloudResourceKind_AzureContainerAppJob,
			rawOutputs: map[string]interface{}{
				"job_id":                "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.App/jobs/nightly-report",
				"job_name":              "nightly-report",
				"event_stream_endpoint": "https://eastus.azurecontainerapps.dev/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/containerAppJobs/nightly-report/eventstream",
				"outbound_ip_addresses": []interface{}{"20.1.2.3", "20.1.2.4"},
				"identity_principal_id": "11111111-2222-3333-4444-555555555555",
			},
			mustPopulate: []string{
				"job_id", "job_name", "event_stream_endpoint",
				"outbound_ip_addresses", "identity_principal_id",
			},
		},
		{
			// AzureContainerAppEnvironmentStorage: storage_name is the
			// seam app and job volumes reference in storage_name.
			name: "AzureContainerAppEnvironmentStorage",
			kind: cloudresourcekind.CloudResourceKind_AzureContainerAppEnvironmentStorage,
			rawOutputs: map[string]interface{}{
				"storage_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.App/managedEnvironments/app-env/storages/app-data",
				"storage_name": "app-data",
			},
			mustPopulate: []string{
				"storage_id", "storage_name",
			},
		},
		{
			// AzureContainerAppEnvironmentDaprComponent: component_name is
			// what application code passes to the Dapr API.
			name: "AzureContainerAppEnvironmentDaprComponent",
			kind: cloudresourcekind.CloudResourceKind_AzureContainerAppEnvironmentDaprComponent,
			rawOutputs: map[string]interface{}{
				"dapr_component_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.App/managedEnvironments/app-env/daprComponents/statestore",
				"component_name":    "statestore",
			},
			mustPopulate: []string{
				"dapr_component_id", "component_name",
			},
		},
		{
			// AzureContainerAppEnvironmentCertificate: certificate_id is the
			// binding seam AzureContainerAppCustomDomain consumes; the
			// certificate facts feed expiry monitoring.
			name: "AzureContainerAppEnvironmentCertificate",
			kind: cloudresourcekind.CloudResourceKind_AzureContainerAppEnvironmentCertificate,
			rawOutputs: map[string]interface{}{
				"certificate_id":  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.App/managedEnvironments/app-env/certificates/app.example.com",
				"subject_name":    "CN=app.example.com",
				"issuer":          "CN=R11, O=Let's Encrypt, C=US",
				"issue_date":      "2026-07-01T00:00:00+00:00",
				"expiration_date": "2026-09-29T00:00:00+00:00",
				"thumbprint":      "A1B2C3D4E5F60718293A4B5C6D7E8F9012345678",
			},
			mustPopulate: []string{
				"certificate_id", "subject_name", "issuer",
				"issue_date", "expiration_date", "thumbprint",
			},
		},
		{
			// AzureContainerAppEnvironmentManagedCertificate:
			// certificate_id identifies the Azure-issued certificate;
			// validation_token is informational once issuance completes.
			name: "AzureContainerAppEnvironmentManagedCertificate",
			kind: cloudresourcekind.CloudResourceKind_AzureContainerAppEnvironmentManagedCertificate,
			rawOutputs: map[string]interface{}{
				"certificate_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.App/managedEnvironments/app-env/managedCertificates/app-example-com",
				"validation_token": "0123456789abcdef0123456789abcdef",
			},
			mustPopulate: []string{
				"certificate_id", "validation_token",
			},
		},
		{
			// AzureContainerAppCustomDomain: custom_domain_id is the
			// providers' synthetic binding identifier (the binding lives
			// inside the app's ingress configuration, not as its own ARM
			// resource). managed_certificate_id is legitimately empty for
			// bring-your-own bindings and until Azure attaches the managed
			// certificate, so it is not asserted.
			name: "AzureContainerAppCustomDomain",
			kind: cloudresourcekind.CloudResourceKind_AzureContainerAppCustomDomain,
			rawOutputs: map[string]interface{}{
				"custom_domain_id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.App/containerApps/web-app/customDomainName/app.example.com",
				"managed_certificate_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.App/managedEnvironments/app-env/managedCertificates/app-example-com",
			},
			mustPopulate: []string{
				"custom_domain_id", "managed_certificate_id",
			},
		},
		{
			// AzureDnsZone: zone_name (with resource_group_name) is the join
			// key AzureDnsRecord addresses record sets through; zone_id is
			// the ARM seam for kinds watching the zone (Front Door custom
			// domains, AKS web-app routing); name_servers is the registrar
			// delegation handoff.
			name: "AzureDnsZone",
			kind: cloudresourcekind.CloudResourceKind_AzureDnsZone,
			rawOutputs: map[string]interface{}{
				"zone_id":                   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/dns-rg/providers/Microsoft.Network/dnsZones/example.com",
				"zone_name":                 "example.com",
				"resource_group_name":       "dns-rg",
				"name_servers":              []interface{}{"ns1-05.azure-dns.com.", "ns2-05.azure-dns.net.", "ns3-05.azure-dns.org.", "ns4-05.azure-dns.info."},
				"max_number_of_record_sets": 10000,
			},
			mustPopulate: []string{
				"zone_id", "zone_name", "resource_group_name",
				"name_servers", "max_number_of_record_sets",
			},
		},
		{
			// AzureDnsRecord: record_id embeds the record type as its own
			// ARM path segment; fqdn is DNS's own trailing-dot spelling.
			name: "AzureDnsRecord",
			kind: cloudresourcekind.CloudResourceKind_AzureDnsRecord,
			rawOutputs: map[string]interface{}{
				"record_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/dns-rg/providers/Microsoft.Network/dnsZones/example.com/MX/@",
				"fqdn":      "example.com.",
			},
			mustPopulate: []string{
				"record_id", "fqdn",
			},
		},
		{
			// AzureLogAnalyticsWorkspace: workspace_id (the ARM id) is the FK
			// seam App Insights / AKS / Container Apps / diagnostic settings
			// reference; workspace_customer_id is the agent-facing GUID the
			// provider confusingly calls workspace_id.
			name: "AzureLogAnalyticsWorkspace",
			kind: cloudresourcekind.CloudResourceKind_AzureLogAnalyticsWorkspace,
			rawOutputs: map[string]interface{}{
				"workspace_id":          "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/obs-rg/providers/Microsoft.OperationalInsights/workspaces/platform-law",
				"workspace_name":        "platform-law",
				"workspace_customer_id": "11111111-2222-3333-4444-555555555555",
				"resource_group_name":   "obs-rg",
				"primary_shared_key":    "cHJpbWFyeS1rZXk=",
				"secondary_shared_key":  "c2Vjb25kYXJ5LWtleQ==",
				"identity_principal_id": "99999999-8888-7777-6666-555555555555",
			},
			mustPopulate: []string{
				"workspace_id", "workspace_name", "workspace_customer_id",
				"resource_group_name", "primary_shared_key",
				"secondary_shared_key", "identity_principal_id",
			},
		},
		{
			// AzureApplicationInsights: connection_string is the seam the
			// app-hosting kinds reference.
			name: "AzureApplicationInsights",
			kind: cloudresourcekind.CloudResourceKind_AzureApplicationInsights,
			rawOutputs: map[string]interface{}{
				"application_insights_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/obs-rg/providers/Microsoft.Insights/components/platform-appinsights",
				"application_insights_name": "platform-appinsights",
				"instrumentation_key":       "22222222-3333-4444-5555-666666666666",
				"connection_string":         "InstrumentationKey=22222222-3333-4444-5555-666666666666;IngestionEndpoint=https://eastus-8.in.applicationinsights.azure.com/",
				"app_id":                    "77777777-8888-9999-aaaa-bbbbbbbbbbbb",
			},
			mustPopulate: []string{
				"application_insights_id", "application_insights_name",
				"instrumentation_key", "connection_string", "app_id",
			},
		},
		{
			// AzureMonitorDiagnosticSetting: the id is the CONSTRUCTED ARM
			// extension-resource id (the provider's own state id is a
			// "{target}|{name}" composite no API consumes).
			name: "AzureMonitorDiagnosticSetting",
			kind: cloudresourcekind.CloudResourceKind_AzureMonitorDiagnosticSetting,
			rawOutputs: map[string]interface{}{
				"diagnostic_setting_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/obs-rg/providers/Microsoft.KeyVault/vaults/app-vault/providers/Microsoft.Insights/diagnosticSettings/route-to-law",
				"diagnostic_setting_name": "route-to-law",
				"target_resource_id":      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/obs-rg/providers/Microsoft.KeyVault/vaults/app-vault",
			},
			mustPopulate: []string{
				"diagnostic_setting_id", "diagnostic_setting_name", "target_resource_id",
			},
		},
		{
			// AzureMonitorActionGroup: action_group_id is the seam alert
			// rules reference.
			name: "AzureMonitorActionGroup",
			kind: cloudresourcekind.CloudResourceKind_AzureMonitorActionGroup,
			rawOutputs: map[string]interface{}{
				"action_group_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/obs-rg/providers/Microsoft.Insights/actionGroups/platform-oncall",
				"action_group_name": "platform-oncall",
			},
			mustPopulate: []string{
				"action_group_id", "action_group_name",
			},
		},
		{
			name: "AzureMonitorMetricAlert",
			kind: cloudresourcekind.CloudResourceKind_AzureMonitorMetricAlert,
			rawOutputs: map[string]interface{}{
				"metric_alert_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/obs-rg/providers/Microsoft.Insights/metricAlerts/storage-availability",
				"metric_alert_name": "storage-availability",
			},
			mustPopulate: []string{
				"metric_alert_id", "metric_alert_name",
			},
		},
		{
			name: "AzureMonitorScheduledQueryAlert",
			kind: cloudresourcekind.CloudResourceKind_AzureMonitorScheduledQueryAlert,
			rawOutputs: map[string]interface{}{
				"scheduled_query_alert_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/obs-rg/providers/Microsoft.Insights/scheduledQueryRules/error-spike",
				"scheduled_query_alert_name": "error-spike",
				"identity_principal_id":      "33333333-4444-5555-6666-777777777777",
			},
			mustPopulate: []string{
				"scheduled_query_alert_id", "scheduled_query_alert_name",
				"identity_principal_id",
			},
		},
		{
			// AzureServiceBusNamespace: namespace_id is the parent seam every
			// Service Bus child kind references (queue, topic, authorization
			// rule, geo-DR pairing); the root SAS rule's four credential
			// faces are the quick-start connection surface.
			name: "AzureServiceBusNamespace",
			kind: cloudresourcekind.CloudResourceKind_AzureServiceBusNamespace,
			rawOutputs: map[string]interface{}{
				"namespace_id":                        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/msg-rg/providers/Microsoft.ServiceBus/namespaces/orders-bus",
				"namespace_name":                      "orders-bus",
				"endpoint":                            "https://orders-bus.servicebus.windows.net:443/",
				"identity_principal_id":               "55555555-6666-7777-8888-999999999999",
				"default_primary_connection_string":   "Endpoint=sb://orders-bus.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=base64key==",
				"default_secondary_connection_string": "Endpoint=sb://orders-bus.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=base64key2==",
				"default_primary_key":                 "base64key==",
				"default_secondary_key":               "base64key2==",
			},
			mustPopulate: []string{
				"namespace_id", "namespace_name", "endpoint",
				"identity_principal_id", "default_primary_connection_string",
				"default_secondary_connection_string", "default_primary_key",
				"default_secondary_key",
			},
		},
		{
			// AzureServiceBusQueue: queue_id is the data-plane RBAC scope and
			// the parent seam queue-scoped SAS rules reference; the
			// namespace/queue name pair is what SDK clients and function
			// bindings consume.
			name: "AzureServiceBusQueue",
			kind: cloudresourcekind.CloudResourceKind_AzureServiceBusQueue,
			rawOutputs: map[string]interface{}{
				"queue_id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/msg-rg/providers/Microsoft.ServiceBus/namespaces/orders-bus/queues/orders",
				"queue_name":     "orders",
				"namespace_name": "orders-bus",
			},
			mustPopulate: []string{
				"queue_id", "queue_name", "namespace_name",
			},
		},
		{
			// AzureServiceBusTopic: topic_id is the parent seam subscriptions
			// and topic-scoped SAS rules reference.
			name: "AzureServiceBusTopic",
			kind: cloudresourcekind.CloudResourceKind_AzureServiceBusTopic,
			rawOutputs: map[string]interface{}{
				"topic_id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/msg-rg/providers/Microsoft.ServiceBus/namespaces/orders-bus/topics/events",
				"topic_name":     "events",
				"namespace_name": "orders-bus",
			},
			mustPopulate: []string{
				"topic_id", "topic_name", "namespace_name",
			},
		},
		{
			// AzureServiceBusSubscription: consumers receive by the
			// namespace/topic/subscription triple.
			name: "AzureServiceBusSubscription",
			kind: cloudresourcekind.CloudResourceKind_AzureServiceBusSubscription,
			rawOutputs: map[string]interface{}{
				"subscription_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/msg-rg/providers/Microsoft.ServiceBus/namespaces/orders-bus/topics/events/subscriptions/audit",
				"subscription_name": "audit",
				"topic_name":        "events",
				"namespace_name":    "orders-bus",
			},
			mustPopulate: []string{
				"subscription_id", "subscription_name", "topic_name",
				"namespace_name",
			},
		},
		{
			// AzureServiceBusAuthorizationRule: authorization_rule_id is the
			// seam the geo-DR pairing's alias_authorization_rule_id consumes;
			// the six key/connection-string faces are the least-privilege
			// credential surface applications hold.
			name: "AzureServiceBusAuthorizationRule",
			kind: cloudresourcekind.CloudResourceKind_AzureServiceBusAuthorizationRule,
			rawOutputs: map[string]interface{}{
				"authorization_rule_id":             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/msg-rg/providers/Microsoft.ServiceBus/namespaces/orders-bus/queues/orders/authorizationRules/orders-sender",
				"rule_name":                         "orders-sender",
				"primary_key":                       "base64key==",
				"secondary_key":                     "base64key2==",
				"primary_connection_string":         "Endpoint=sb://orders-bus.servicebus.windows.net/;SharedAccessKeyName=orders-sender;SharedAccessKey=base64key==;EntityPath=orders",
				"secondary_connection_string":       "Endpoint=sb://orders-bus.servicebus.windows.net/;SharedAccessKeyName=orders-sender;SharedAccessKey=base64key2==;EntityPath=orders",
				"primary_connection_string_alias":   "",
				"secondary_connection_string_alias": "",
			},
			mustPopulate: []string{
				"authorization_rule_id", "rule_name", "primary_key",
				"secondary_key", "primary_connection_string",
				"secondary_connection_string",
			},
		},
		{
			// AzureServiceBusDisasterRecoveryConfig: the alias connection
			// strings are what DR-aware clients hold -- they survive a
			// failover without reconfiguration.
			name: "AzureServiceBusDisasterRecoveryConfig",
			kind: cloudresourcekind.CloudResourceKind_AzureServiceBusDisasterRecoveryConfig,
			rawOutputs: map[string]interface{}{
				"disaster_recovery_config_id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/msg-rg/providers/Microsoft.ServiceBus/namespaces/orders-bus-eastus/disasterRecoveryConfigs/orders-bus-alias",
				"alias_name":                        "orders-bus-alias",
				"primary_connection_string_alias":   "Endpoint=sb://orders-bus-alias.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=base64key==",
				"secondary_connection_string_alias": "Endpoint=sb://orders-bus-alias.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=base64key2==",
				"default_primary_key":               "base64key==",
				"default_secondary_key":             "base64key2==",
			},
			mustPopulate: []string{
				"disaster_recovery_config_id", "alias_name",
				"primary_connection_string_alias",
				"secondary_connection_string_alias", "default_primary_key",
				"default_secondary_key",
			},
		},
		{
			// AzureEventHubNamespace: namespace_id is the parent seam every
			// Event Hubs child kind references (hub, authorization rule,
			// schema group, geo-DR pairing, CMK); the root SAS rule's
			// credential faces (incl. the geo-DR alias pair) are the
			// quick-start connection surface.
			name: "AzureEventHubNamespace",
			kind: cloudresourcekind.CloudResourceKind_AzureEventHubNamespace,
			rawOutputs: map[string]interface{}{
				"namespace_id":                              "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/stream-rg/providers/Microsoft.EventHub/namespaces/telemetry-hubs",
				"namespace_name":                            "telemetry-hubs",
				"identity_principal_id":                     "55555555-6666-7777-8888-999999999999",
				"default_primary_connection_string":         "Endpoint=sb://telemetry-hubs.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=base64key==",
				"default_secondary_connection_string":       "Endpoint=sb://telemetry-hubs.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=base64key2==",
				"default_primary_key":                       "base64key==",
				"default_secondary_key":                     "base64key2==",
				"default_primary_connection_string_alias":   "",
				"default_secondary_connection_string_alias": "",
			},
			mustPopulate: []string{
				"namespace_id", "namespace_name", "identity_principal_id",
				"default_primary_connection_string",
				"default_secondary_connection_string", "default_primary_key",
				"default_secondary_key",
			},
		},
		{
			// AzureEventHub: event_hub_id is the parent seam consumer groups
			// and hub-scoped SAS rules reference; partition_ids is the
			// repeated output partition-aware consumers enumerate.
			name: "AzureEventHub",
			kind: cloudresourcekind.CloudResourceKind_AzureEventHub,
			rawOutputs: map[string]interface{}{
				"event_hub_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/stream-rg/providers/Microsoft.EventHub/namespaces/telemetry-hubs/eventhubs/telemetry",
				"event_hub_name": "telemetry",
				"partition_ids":  []interface{}{"0", "1", "2", "3"},
			},
			mustPopulate: []string{
				"event_hub_id", "event_hub_name", "partition_ids",
			},
		},
		{
			// AzureEventHubConsumerGroup: the group name is what consumer
			// applications pass to their SDK client alongside the hub name.
			name: "AzureEventHubConsumerGroup",
			kind: cloudresourcekind.CloudResourceKind_AzureEventHubConsumerGroup,
			rawOutputs: map[string]interface{}{
				"consumer_group_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/stream-rg/providers/Microsoft.EventHub/namespaces/telemetry-hubs/eventhubs/telemetry/consumergroups/analytics",
				"consumer_group_name": "analytics",
			},
			mustPopulate: []string{
				"consumer_group_id", "consumer_group_name",
			},
		},
		{
			// AzureEventHubAuthorizationRule: identical credential faces
			// regardless of scope; the alias pair is only populated when a
			// geo-DR pairing exists.
			name: "AzureEventHubAuthorizationRule",
			kind: cloudresourcekind.CloudResourceKind_AzureEventHubAuthorizationRule,
			rawOutputs: map[string]interface{}{
				"authorization_rule_id":             "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/stream-rg/providers/Microsoft.EventHub/namespaces/telemetry-hubs/eventhubs/telemetry/authorizationRules/producer-send",
				"rule_name":                         "producer-send",
				"primary_key":                       "base64key==",
				"secondary_key":                     "base64key2==",
				"primary_connection_string":         "Endpoint=sb://telemetry-hubs.servicebus.windows.net/;SharedAccessKeyName=producer-send;SharedAccessKey=base64key==;EntityPath=telemetry",
				"secondary_connection_string":       "Endpoint=sb://telemetry-hubs.servicebus.windows.net/;SharedAccessKeyName=producer-send;SharedAccessKey=base64key2==;EntityPath=telemetry",
				"primary_connection_string_alias":   "",
				"secondary_connection_string_alias": "",
			},
			mustPopulate: []string{
				"authorization_rule_id", "rule_name", "primary_key",
				"secondary_key", "primary_connection_string",
				"secondary_connection_string",
			},
		},
		{
			// AzureEventHubDisasterRecoveryConfig: alias credentials
			// deliberately live on the namespace/authorization-rule kinds
			// (Azure's own surface) -- this kind exports the pairing
			// identity only.
			name: "AzureEventHubDisasterRecoveryConfig",
			kind: cloudresourcekind.CloudResourceKind_AzureEventHubDisasterRecoveryConfig,
			rawOutputs: map[string]interface{}{
				"disaster_recovery_config_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/stream-rg/providers/Microsoft.EventHub/namespaces/telemetry-hubs/disasterRecoveryConfigs/telemetry-alias",
				"alias_name":                  "telemetry-alias",
			},
			mustPopulate: []string{
				"disaster_recovery_config_id", "alias_name",
			},
		},
		{
			// AzureEventHubSchemaGroup: the group name is what
			// schema-registry serializers address at runtime.
			name: "AzureEventHubSchemaGroup",
			kind: cloudresourcekind.CloudResourceKind_AzureEventHubSchemaGroup,
			rawOutputs: map[string]interface{}{
				"schema_group_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/stream-rg/providers/Microsoft.EventHub/namespaces/telemetry-hubs/schemagroups/telemetry-schemas",
				"schema_group_name": "telemetry-schemas",
			},
			mustPopulate: []string{
				"schema_group_id", "schema_group_name",
			},
		},
		{
			// AzureEventHubCluster: cluster_id is the seam a namespace's
			// dedicated_cluster_id references for single-tenant placement.
			name: "AzureEventHubCluster",
			kind: cloudresourcekind.CloudResourceKind_AzureEventHubCluster,
			rawOutputs: map[string]interface{}{
				"cluster_id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/stream-rg/providers/Microsoft.EventHub/clusters/streaming-dedicated",
				"cluster_name": "streaming-dedicated",
			},
			mustPopulate: []string{
				"cluster_id", "cluster_name",
			},
		},
		{
			// AzureEventHubNamespaceCustomerManagedKey: the configuration is
			// a property of the namespace (no ARM object of its own), so its
			// identity output is the namespace's ARM id.
			name: "AzureEventHubNamespaceCustomerManagedKey",
			kind: cloudresourcekind.CloudResourceKind_AzureEventHubNamespaceCustomerManagedKey,
			rawOutputs: map[string]interface{}{
				"customer_managed_key_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/stream-rg/providers/Microsoft.EventHub/namespaces/telemetry-hubs",
			},
			mustPopulate: []string{
				"customer_managed_key_id",
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
