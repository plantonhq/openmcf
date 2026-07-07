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
			// GcpServiceAccount: flat scalar outputs from both engines (email, the
			// ready-made IAM member string, stable unique id, fully-qualified name,
			// and the optional key) must each land on the StackOutputs proto.
			name: "GcpServiceAccount",
			kind: cloudresourcekind.CloudResourceKind_GcpServiceAccount,
			rawOutputs: map[string]interface{}{
				"email":      "my-sa@my-project.iam.gserviceaccount.com",
				"member":     "serviceAccount:my-sa@my-project.iam.gserviceaccount.com",
				"unique_id":  "112233445566778899000",
				"name":       "projects/my-project/serviceAccounts/my-sa@my-project.iam.gserviceaccount.com",
				"key_base64": "eyJ0eXBlIjoic2VydmljZV9hY2NvdW50In0=",
			},
			mustPopulate: []string{
				"email", "member", "unique_id", "name", "key_base64",
			},
		},
		{
			// GcpIamCustomRole: flat scalar outputs from both engines (the grantable
			// fully-qualified role name, the bare role id, and the soft-delete flag)
			// must each land on the StackOutputs proto.
			name: "GcpIamCustomRole",
			kind: cloudresourcekind.CloudResourceKind_GcpIamCustomRole,
			rawOutputs: map[string]interface{}{
				"name":    "projects/my-project/roles/logBucketWriter",
				"role_id": "logBucketWriter",
				"deleted": "false",
			},
			mustPopulate: []string{
				"name", "role_id",
			},
		},
		{
			// GcpProjectIamMember: the grant tuple echoed by both engines (project,
			// role, member, policy etag) must each land on the StackOutputs proto.
			name: "GcpProjectIamMember",
			kind: cloudresourcekind.CloudResourceKind_GcpProjectIamMember,
			rawOutputs: map[string]interface{}{
				"project_id": "my-project",
				"role":       "projects/my-project/roles/logBucketWriter",
				"member":     "serviceAccount:my-sa@my-project.iam.gserviceaccount.com",
				"etag":       "BwYn2FQlJeM=",
			},
			mustPopulate: []string{
				"project_id", "role", "member", "etag",
			},
		},
		{
			// GcpWorkloadIdentityPool: flat scalar outputs from both engines (the
			// full pool resource name principals embed, the bare pool id providers
			// reference, and the lifecycle state) must each land on the
			// StackOutputs proto.
			name: "GcpWorkloadIdentityPool",
			kind: cloudresourcekind.CloudResourceKind_GcpWorkloadIdentityPool,
			rawOutputs: map[string]interface{}{
				"name":                      "projects/123456789/locations/global/workloadIdentityPools/github-actions",
				"workload_identity_pool_id": "github-actions",
				"state":                     "ACTIVE",
			},
			mustPopulate: []string{
				"name", "workload_identity_pool_id", "state",
			},
		},
		{
			// GcpWorkloadIdentityPoolProvider: flat scalar outputs from both
			// engines (the full provider resource name — the token-exchange
			// audience — the bare provider id, and the lifecycle state) must each
			// land on the StackOutputs proto.
			name: "GcpWorkloadIdentityPoolProvider",
			kind: cloudresourcekind.CloudResourceKind_GcpWorkloadIdentityPoolProvider,
			rawOutputs: map[string]interface{}{
				"name":                               "projects/123456789/locations/global/workloadIdentityPools/github-actions/providers/github-oidc",
				"workload_identity_pool_provider_id": "github-oidc",
				"state":                              "ACTIVE",
			},
			mustPopulate: []string{
				"name", "workload_identity_pool_provider_id", "state",
			},
		},
		{
			// GcpHealthCheck: flat scalar outputs from both engines (the self-link
			// backend services reference, the cloud-side name, the computed probe
			// type, and the scope-marking region — empty for global) must each
			// land on the StackOutputs proto.
			name: "GcpHealthCheck",
			kind: cloudresourcekind.CloudResourceKind_GcpHealthCheck,
			rawOutputs: map[string]interface{}{
				"self_link":         "https://www.googleapis.com/compute/v1/projects/my-project/global/healthChecks/web-probe",
				"health_check_name": "web-probe",
				"type":              "HTTP",
				"region":            "",
			},
			mustPopulate: []string{
				"self_link", "health_check_name", "type",
			},
		},
		{
			// GcpBackendBucket: flat scalar outputs from both engines (the
			// self-link URL maps reference, the cloud-side name, and the origin
			// bucket) must each land on the StackOutputs proto.
			name: "GcpBackendBucket",
			kind: cloudresourcekind.CloudResourceKind_GcpBackendBucket,
			rawOutputs: map[string]interface{}{
				"self_link":           "https://www.googleapis.com/compute/v1/projects/my-project/global/backendBuckets/static-assets",
				"backend_bucket_name": "static-assets",
				"bucket_name":         "my-assets-bucket",
			},
			mustPopulate: []string{
				"self_link", "backend_bucket_name", "bucket_name",
			},
		},
		{
			// GcpBackendService: flat scalar outputs from both engines (the
			// self-link URL maps reference, the cloud-side name, the numeric id,
			// and the concurrency fingerprint) must each land on the
			// StackOutputs proto.
			name: "GcpBackendService",
			kind: cloudresourcekind.CloudResourceKind_GcpBackendService,
			rawOutputs: map[string]interface{}{
				"self_link":            "https://www.googleapis.com/compute/v1/projects/my-project/global/backendServices/web-backend",
				"backend_service_name": "web-backend",
				"generated_id":         "1234567890123456789",
				"fingerprint":          "BwYn2FQlJeM=",
			},
			mustPopulate: []string{
				"self_link", "backend_service_name", "generated_id", "fingerprint",
			},
		},
		{
			// GcpRegionNetworkEndpointGroup: flat scalar outputs from both engines
			// (self-link, name, endpoint type, region) must land on StackOutputs.
			name: "GcpRegionNetworkEndpointGroup",
			kind: cloudresourcekind.CloudResourceKind_GcpRegionNetworkEndpointGroup,
			rawOutputs: map[string]interface{}{
				"self_link":                   "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1/networkEndpointGroups/my-neg",
				"network_endpoint_group_name": "my-neg",
				"network_endpoint_type":       "SERVERLESS",
				"region":                      "us-central1",
			},
			mustPopulate: []string{
				"self_link", "network_endpoint_group_name", "network_endpoint_type", "region",
			},
		},
		{
			// GcpUrlMap: flat scalar outputs from both engines (self-link, name,
			// numeric id, fingerprint) must land on StackOutputs.
			name: "GcpUrlMap",
			kind: cloudresourcekind.CloudResourceKind_GcpUrlMap,
			rawOutputs: map[string]interface{}{
				"self_link":    "https://www.googleapis.com/compute/v1/projects/my-project/global/urlMaps/my-map",
				"url_map_name": "my-map",
				"map_id":       "1234567890123456789",
				"fingerprint":  "BwYn2FQlJeM=",
			},
			mustPopulate: []string{
				"self_link", "url_map_name", "map_id", "fingerprint",
			},
		},
		{
			// GcpManagedSslCertificate: flat scalar outputs from both engines
			// (self-link, name, id, expire_time) must land on StackOutputs.
			name: "GcpManagedSslCertificate",
			kind: cloudresourcekind.CloudResourceKind_GcpManagedSslCertificate,
			rawOutputs: map[string]interface{}{
				"self_link":        "https://www.googleapis.com/compute/v1/projects/my-project/global/sslCertificates/my-cert",
				"certificate_name": "my-cert",
				"certificate_id":   "1234567890123456789",
				"expire_time":      "2027-01-01T00:00:00Z",
			},
			mustPopulate: []string{
				"self_link", "certificate_name", "certificate_id", "expire_time",
			},
		},
		{
			// GcpTargetHttpProxy: flat scalar outputs from both engines
			// (self-link, name, numeric id, fingerprint) must land on StackOutputs.
			name: "GcpTargetHttpProxy",
			kind: cloudresourcekind.CloudResourceKind_GcpTargetHttpProxy,
			rawOutputs: map[string]interface{}{
				"self_link":   "https://www.googleapis.com/compute/v1/projects/my-project/global/targetHttpProxies/my-proxy",
				"proxy_name":  "my-proxy",
				"proxy_id":    "1234567890123456789",
				"fingerprint": "BwYn2FQlJeM=",
			},
			mustPopulate: []string{
				"self_link", "proxy_name", "proxy_id", "fingerprint",
			},
		},
		{
			// GcpTargetHttpsProxy: flat scalar outputs from both engines
			// (self-link, name, numeric id, fingerprint) must land on StackOutputs.
			name: "GcpTargetHttpsProxy",
			kind: cloudresourcekind.CloudResourceKind_GcpTargetHttpsProxy,
			rawOutputs: map[string]interface{}{
				"self_link":   "https://www.googleapis.com/compute/v1/projects/my-project/global/targetHttpsProxies/my-proxy",
				"proxy_name":  "my-proxy",
				"proxy_id":    "1234567890123456789",
				"fingerprint": "BwYn2FQlJeM=",
			},
			mustPopulate: []string{
				"self_link", "proxy_name", "proxy_id", "fingerprint",
			},
		},
		{
			// GcpGlobalForwardingRule: flat scalar outputs from both engines (the
			// VIP, self-link, name, numeric id, and the PSC connection fields)
			// must land on StackOutputs.
			name: "GcpGlobalForwardingRule",
			kind: cloudresourcekind.CloudResourceKind_GcpGlobalForwardingRule,
			rawOutputs: map[string]interface{}{
				"ip_address":            "34.120.1.2",
				"self_link":             "https://www.googleapis.com/compute/v1/projects/my-project/global/forwardingRules/my-frontend",
				"forwarding_rule_name":  "my-frontend",
				"forwarding_rule_id":    "1234567890123456789",
				"psc_connection_id":     "1111222233334444",
				"psc_connection_status": "ACCEPTED",
			},
			mustPopulate: []string{
				"ip_address", "self_link", "forwarding_rule_name",
				"forwarding_rule_id", "psc_connection_id", "psc_connection_status",
			},
		},
		{
			// GcpSslPolicy: flat scalar outputs plus the repeated enabled_features
			// cipher list both engines emit (per-index keys) must land on
			// StackOutputs.
			name: "GcpSslPolicy",
			kind: cloudresourcekind.CloudResourceKind_GcpSslPolicy,
			rawOutputs: map[string]interface{}{
				"self_link":          "https://www.googleapis.com/compute/v1/projects/my-project/global/sslPolicies/my-policy",
				"ssl_policy_name":    "my-policy",
				"enabled_features.0": "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
				"enabled_features.1": "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
				"region":             "",
			},
			mustPopulate: []string{
				"self_link", "ssl_policy_name", "enabled_features",
			},
		},
		{
			// GcpSslCertificate: flat scalar outputs from both engines
			// (self-link, name, id, expiry, scope region) must land on
			// StackOutputs. The private key is write-only and never an output.
			name: "GcpSslCertificate",
			kind: cloudresourcekind.CloudResourceKind_GcpSslCertificate,
			rawOutputs: map[string]interface{}{
				"self_link":        "https://www.googleapis.com/compute/v1/projects/my-project/global/sslCertificates/my-cert",
				"certificate_name": "my-cert",
				"certificate_id":   "1234567890123456789",
				"expire_time":      "2036-06-30T12:36:27Z",
				"region":           "",
			},
			mustPopulate: []string{
				"self_link", "certificate_name", "certificate_id", "expire_time",
			},
		},
		{
			// GcpGlobalAddress: name output added for service networking composition.
			name: "GcpGlobalAddress",
			kind: cloudresourcekind.CloudResourceKind_GcpGlobalAddress,
			rawOutputs: map[string]interface{}{
				"address":            "10.100.0.0",
				"self_link":          "https://www.googleapis.com/compute/v1/projects/my-project/global/addresses/vpc-peering-range",
				"creation_timestamp": "2026-01-01T00:00:00Z",
				"name":               "vpc-peering-range",
			},
			mustPopulate: []string{
				"address", "self_link", "creation_timestamp", "name",
			},
		},
		{
			// GcpServiceNetworkingConnection: peering + network outputs from both engines.
			name: "GcpServiceNetworkingConnection",
			kind: cloudresourcekind.CloudResourceKind_GcpServiceNetworkingConnection,
			rawOutputs: map[string]interface{}{
				"peering": "servicenetworking-googleapis-com",
				"network": "projects/my-project/global/networks/app-vpc",
			},
			mustPopulate: []string{"peering", "network"},
		},
		{
			// GcpAddress: regional reservation outputs including plain spec region.
			name: "GcpAddress",
			kind: cloudresourcekind.CloudResourceKind_GcpAddress,
			rawOutputs: map[string]interface{}{
				"address":   "203.0.113.10",
				"self_link": "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1/addresses/nat-ip",
				"name":      "nat-ip",
				"region":    "us-central1",
			},
			mustPopulate: []string{"address", "self_link", "name", "region"},
		},
		{
			// GcpCloudSql: instance outputs from both engines, including the
			// service-account identity and the PSC-only fields (empty on
			// non-PSC instances).
			name: "GcpCloudSql",
			kind: cloudresourcekind.CloudResourceKind_GcpCloudSql,
			rawOutputs: map[string]interface{}{
				"instance_name":               "orders-db",
				"connection_name":             "my-project:us-central1:orders-db",
				"private_ip":                  "10.20.0.5",
				"public_ip":                   "",
				"self_link":                   "https://sqladmin.googleapis.com/sql/v1beta4/projects/my-project/instances/orders-db",
				"service_account_email":       "p1234-abcdef@gcp-sa-cloud-sql.iam.gserviceaccount.com",
				"dns_name":                    "",
				"psc_service_attachment_link": "",
			},
			mustPopulate: []string{
				"instance_name", "connection_name", "self_link", "service_account_email",
			},
		},
		{
			// GcpCloudSqlDatabase: database name + self link.
			name: "GcpCloudSqlDatabase",
			kind: cloudresourcekind.CloudResourceKind_GcpCloudSqlDatabase,
			rawOutputs: map[string]interface{}{
				"database_name": "orders",
				"self_link":     "https://sqladmin.googleapis.com/sql/v1beta4/projects/my-project/instances/orders-db/databases/orders",
			},
			mustPopulate: []string{"database_name", "self_link"},
		},
		{
			// GcpCloudSqlUser: the stored user name (IAM users on MySQL come
			// back truncated before the @).
			name: "GcpCloudSqlUser",
			kind: cloudresourcekind.CloudResourceKind_GcpCloudSqlUser,
			rawOutputs: map[string]interface{}{
				"user_name":     "orders-app",
				"instance_name": "orders-db",
			},
			mustPopulate: []string{"user_name", "instance_name"},
		},
		{
			// GcpRedisInstance: endpoint scalars, the secret AUTH string, the
			// repeated CA-cert PEMs (populated when TLS is on), the import/export
			// IAM identity, and the effective reserved range.
			name: "GcpRedisInstance",
			kind: cloudresourcekind.CloudResourceKind_GcpRedisInstance,
			rawOutputs: map[string]interface{}{
				"host":                        "10.118.0.4",
				"port":                        "6379",
				"read_endpoint":               "10.118.0.5",
				"read_endpoint_port":          "6379",
				"current_location_id":         "us-central1-a",
				"auth_string":                 "d1f0e2c3-4b5a-6789-abcd-ef0123456789",
				"server_ca_certs":             []interface{}{"-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----"},
				"persistence_iam_identity":    "serviceAccount:service-1234@gcp-sa-redis.iam.gserviceaccount.com",
				"effective_reserved_ip_range": "10.118.0.0/29",
				"instance_name":               "session-store-prod",
				"region":                      "us-central1",
			},
			mustPopulate: []string{
				"host", "port", "read_endpoint", "read_endpoint_port",
				"current_location_id", "auth_string", "server_ca_certs",
				"persistence_iam_identity", "effective_reserved_ip_range",
				"instance_name", "region",
			},
		},
		{
			// GcpGkeCluster: control-plane endpoint + CA trust anchor, the
			// Workload Identity pool, the fully qualified cluster ID, and the
			// name/location handles node pools compose against.
			name: "GcpGkeCluster",
			kind: cloudresourcekind.CloudResourceKind_GcpGkeCluster,
			rawOutputs: map[string]interface{}{
				"endpoint":               "34.72.10.11",
				"cluster_ca_certificate": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t",
				"workload_identity_pool": "my-project.svc.id.goog",
				"cluster_id":             "projects/my-project/locations/us-central1/clusters/prod-primary",
				"name":                   "prod-primary",
				"location":               "us-central1",
				"self_link":              "https://container.googleapis.com/v1/projects/my-project/locations/us-central1/clusters/prod-primary",
				"master_version":         "1.31.4-gke.1256000",
			},
			mustPopulate: []string{
				"endpoint", "cluster_ca_certificate", "workload_identity_pool",
				"cluster_id", "name", "location", "self_link", "master_version",
			},
		},
		{
			// GcpGkeNodePool: the pool name/location handles, the backing
			// instance groups, the effective sizing bounds, and the fully
			// qualified pool ID downstream services (e.g. Dataproc on GKE)
			// reference.
			name: "GcpGkeNodePool",
			kind: cloudresourcekind.CloudResourceKind_GcpGkeNodePool,
			rawOutputs: map[string]interface{}{
				"node_pool_name": "general-pool",
				"instance_group_urls": []interface{}{
					"https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a/instanceGroupManagers/gke-prod-primary-general-pool-grp",
				},
				"min_nodes":          "1",
				"max_nodes":          "5",
				"current_node_count": "2",
				"node_pool_id":       "projects/my-project/locations/us-central1/clusters/prod-primary/nodePools/general-pool",
				"location":           "us-central1",
				"version":            "1.31.4-gke.1256000",
			},
			mustPopulate: []string{
				"node_pool_name", "instance_group_urls", "min_nodes",
				"max_nodes", "current_node_count", "node_pool_id",
				"location", "version",
			},
		},
		{
			// GcpCloudRun: the serving URL, the service-name handle serverless
			// NEGs reference, the latest ready revision, and the identifiers
			// API callers use.
			name: "GcpCloudRun",
			kind: cloudresourcekind.CloudResourceKind_GcpCloudRun,
			rawOutputs: map[string]interface{}{
				"url":          "https://my-api-abc123-uc.a.run.app",
				"service_name": "my-api",
				"revision":     "my-api-00042-abc",
				"location":     "us-central1",
				"uid":          "12345678-1234-1234-1234-123456789012",
				"urls": []interface{}{
					"https://my-api-abc123-uc.a.run.app",
				},
			},
			mustPopulate: []string{
				"url", "service_name", "revision", "location", "uid", "urls",
			},
		},
		{
			// GcpRouterNat: NAT name, router self link, and the manual NAT IP
			// self links (empty for auto-allocation).
			name: "GcpRouterNat",
			kind: cloudresourcekind.CloudResourceKind_GcpRouterNat,
			rawOutputs: map[string]interface{}{
				"name":             "prod-nat",
				"router_self_link": "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1/routers/prod-router",
				"nat_ips":          []interface{}{"https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1/addresses/egress-ip-a"},
			},
			mustPopulate: []string{"name", "router_self_link", "nat_ips"},
		},
		{
			// GcpVpcNetwork: deep-rebuilt outputs — PSA fields removed, gateway + ULA added.
			name: "GcpVpcNetwork",
			kind: cloudresourcekind.CloudResourceKind_GcpVpcNetwork,
			rawOutputs: map[string]interface{}{
				"network_self_link":   "https://www.googleapis.com/compute/v1/projects/my-project/global/networks/app-vpc",
				"network_name":        "app-vpc",
				"network_id":          "projects/my-project/global/networks/app-vpc",
				"gateway_ipv4":        "10.128.0.1",
				"internal_ipv6_range": "fd20:1234:5678::/48",
			},
			mustPopulate: []string{
				"network_self_link", "network_name", "network_id",
			},
		},
		{
			// GcpSubnetwork: scalar outputs plus the per-index secondary-range
			// exports both engines emit must land on the StackOutputs proto,
			// including the repeated secondary_ranges message.
			name: "GcpSubnetwork",
			kind: cloudresourcekind.CloudResourceKind_GcpSubnetwork,
			rawOutputs: map[string]interface{}{
				"subnetwork_self_link":             "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1/subnetworks/app-subnet",
				"subnetwork_name":                  "app-subnet",
				"region":                           "us-central1",
				"ip_cidr_range":                    "10.10.0.0/20",
				"gateway_address":                  "10.10.0.1",
				"subnetwork_id":                    "1234567890123456789",
				"internal_ipv6_prefix":             "",
				"external_ipv6_prefix":             "",
				"secondary_ranges.0.range_name":    "pods",
				"secondary_ranges.0.ip_cidr_range": "10.16.0.0/14",
				"secondary_ranges.1.range_name":    "services",
				"secondary_ranges.1.ip_cidr_range": "10.20.0.0/20",
			},
			mustPopulate: []string{
				"subnetwork_self_link", "subnetwork_name", "region",
				"ip_cidr_range", "gateway_address", "subnetwork_id",
			},
		},
		{
			// GcpAlloydbCluster: the fully qualified cluster path (the FK target
			// instance and user kinds parent by), the short name, and the bundled
			// primary instance's connection endpoint.
			name: "GcpAlloydbCluster",
			kind: cloudresourcekind.CloudResourceKind_GcpAlloydbCluster,
			rawOutputs: map[string]interface{}{
				"cluster_id":            "projects/my-project/locations/us-central1/clusters/orders-alloydb",
				"cluster_name":          "orders-alloydb",
				"primary_instance_ip":   "10.30.0.5",
				"primary_instance_name": "projects/my-project/locations/us-central1/clusters/orders-alloydb/instances/primary",
				"database_version":      "POSTGRES_16",
				"state":                 "READY",
			},
			mustPopulate: []string{
				"cluster_id", "cluster_name", "primary_instance_ip",
				"primary_instance_name", "database_version", "state",
			},
		},
		{
			// GcpAlloydbInstance: the fully qualified instance path, its private
			// connection endpoint, and lifecycle state.
			name: "GcpAlloydbInstance",
			kind: cloudresourcekind.CloudResourceKind_GcpAlloydbInstance,
			rawOutputs: map[string]interface{}{
				"instance_name": "projects/my-project/locations/us-central1/clusters/orders-alloydb/instances/read-pool",
				"ip_address":    "10.30.0.7",
				"state":         "READY",
			},
			mustPopulate: []string{"instance_name", "ip_address", "state"},
		},
		{
			// GcpAlloydbUser: the fully qualified user path plus the id and
			// cluster handles.
			name: "GcpAlloydbUser",
			kind: cloudresourcekind.CloudResourceKind_GcpAlloydbUser,
			rawOutputs: map[string]interface{}{
				"name":       "projects/my-project/locations/us-central1/clusters/orders-alloydb/users/orders-app",
				"user_id":    "orders-app",
				"cluster_id": "projects/my-project/locations/us-central1/clusters/orders-alloydb",
			},
			mustPopulate: []string{"name", "user_id", "cluster_id"},
		},
		{
			// GcpDnsZone: the numeric zone id, the zone-name handle GcpDnsRecord
			// composes against, and the delegated nameserver set.
			name: "GcpDnsZone",
			kind: cloudresourcekind.CloudResourceKind_GcpDnsZone,
			rawOutputs: map[string]interface{}{
				"zone_id":   "1234567890123456789",
				"zone_name": "example-com",
				"nameservers": []interface{}{
					"ns-cloud-a1.googledomains.com.",
					"ns-cloud-a2.googledomains.com.",
				},
			},
			mustPopulate: []string{"zone_id", "zone_name", "nameservers"},
		},
		{
			// GcpCloudRunJob: the job-name handle gcloud/Scheduler trigger, the
			// region, the server-assigned uid, and the latest execution (empty
			// until first run).
			name: "GcpCloudRunJob",
			kind: cloudresourcekind.CloudResourceKind_GcpCloudRunJob,
			rawOutputs: map[string]interface{}{
				"job_name":                 "nightly-etl",
				"location":                 "us-central1",
				"uid":                      "12345678-1234-1234-1234-123456789012",
				"latest_created_execution": "",
			},
			mustPopulate: []string{"job_name", "location", "uid"},
		},
		{
			// GcpSpannerInstance: the fully qualified instance path (the IAM/API
			// handle), the short name downstream databases and backup schedules
			// reference, the lifecycle state, and the geographic config.
			name: "GcpSpannerInstance",
			kind: cloudresourcekind.CloudResourceKind_GcpSpannerInstance,
			rawOutputs: map[string]interface{}{
				"instance_id":   "projects/prod-project/instances/orders-spanner",
				"instance_name": "orders-spanner",
				"state":         "READY",
				"config":        "regional-us-central1",
			},
			mustPopulate: []string{"instance_id", "instance_name", "state", "config"},
		},
		{
			// GcpSpannerDatabase: the fully qualified database path (the IAM/API
			// handle), the short name backup schedules reference, and the
			// lifecycle state.
			name: "GcpSpannerDatabase",
			kind: cloudresourcekind.CloudResourceKind_GcpSpannerDatabase,
			rawOutputs: map[string]interface{}{
				"database_id":   "projects/prod-project/instances/orders-spanner/databases/orders",
				"database_name": "orders",
				"state":         "READY",
			},
			mustPopulate: []string{"database_id", "database_name", "state"},
		},
		{
			// GcpSpannerBackupSchedule: the fully qualified schedule path (the
			// API handle) and the short name within the database.
			name: "GcpSpannerBackupSchedule",
			kind: cloudresourcekind.CloudResourceKind_GcpSpannerBackupSchedule,
			rawOutputs: map[string]interface{}{
				"schedule_id":   "projects/prod-project/instances/orders-spanner/databases/orders/backupSchedules/daily-backups",
				"schedule_name": "daily-backups",
			},
			mustPopulate: []string{"schedule_id", "schedule_name"},
		},
		{
			// GcpBigQueryDataset: the short dataset id SQL queries reference,
			// the self link, resolved project, creation time, location, and etag.
			name: "GcpBigQueryDataset",
			kind: cloudresourcekind.CloudResourceKind_GcpBigQueryDataset,
			rawOutputs: map[string]interface{}{
				"dataset_id":    "analytics_prod",
				"self_link":     "https://bigquery.googleapis.com/bigquery/v2/projects/prod-project/datasets/analytics_prod",
				"project":       "prod-project",
				"creation_time": int64(1700000000000),
				"location":      "US",
				"etag":          "abc123",
			},
			mustPopulate: []string{"dataset_id", "self_link", "project", "creation_time", "location", "etag"},
		},
		{
			// GcpBigQueryTable: the short table id, self link, resolved project,
			// parent dataset, table type, location, creation time, and the
			// pre-assembled dotted handle Pub/Sub BigQuery delivery consumes.
			name: "GcpBigQueryTable",
			kind: cloudresourcekind.CloudResourceKind_GcpBigQueryTable,
			rawOutputs: map[string]interface{}{
				"table_id":       "events_raw",
				"self_link":      "https://bigquery.googleapis.com/bigquery/v2/projects/prod-project/datasets/analytics_prod/tables/events_raw",
				"project":        "prod-project",
				"dataset_id":     "analytics_prod",
				"type":           "TABLE",
				"location":       "US",
				"creation_time":  int64(1700000000000),
				"qualified_name": "prod-project.analytics_prod.events_raw",
			},
			mustPopulate: []string{"table_id", "self_link", "project", "dataset_id", "type", "location", "creation_time", "qualified_name"},
		},
		{
			// GcpPubSubSchema: the fully qualified schema path a topic's
			// schema_settings.schema reference consumes, and the short name.
			name: "GcpPubSubSchema",
			kind: cloudresourcekind.CloudResourceKind_GcpPubSubSchema,
			rawOutputs: map[string]interface{}{
				"schema_id":   "projects/prod-project/schemas/order-events",
				"schema_name": "order-events",
			},
			mustPopulate: []string{"schema_id", "schema_name"},
		},
		{
			// GcpPubSubTopic: the fully qualified topic path subscriptions
			// and event triggers consume, and the short name.
			name: "GcpPubSubTopic",
			kind: cloudresourcekind.CloudResourceKind_GcpPubSubTopic,
			rawOutputs: map[string]interface{}{
				"topic_id":   "projects/prod-project/topics/order-events",
				"topic_name": "order-events",
			},
			mustPopulate: []string{"topic_id", "topic_name"},
		},
		{
			// GcpPubSubSubscription: the fully qualified subscription path
			// consumers and monitoring reference, and the short name.
			name: "GcpPubSubSubscription",
			kind: cloudresourcekind.CloudResourceKind_GcpPubSubSubscription,
			rawOutputs: map[string]interface{}{
				"subscription_id":   "projects/prod-project/subscriptions/order-events-worker",
				"subscription_name": "order-events-worker",
			},
			mustPopulate: []string{"subscription_id", "subscription_name"},
		},
		{
			// GcpServerlessVpcConnector: the short connector name, the fully
			// qualified path serverless workloads attach to, the reconciled
			// state, and the plain region name.
			name: "GcpServerlessVpcConnector",
			kind: cloudresourcekind.CloudResourceKind_GcpServerlessVpcConnector,
			rawOutputs: map[string]interface{}{
				"name":      "svc-egress",
				"self_link": "projects/prod-project/locations/us-central1/connectors/svc-egress",
				"state":     "READY",
				"region":    "us-central1",
			},
			mustPopulate: []string{"name", "self_link", "state", "region"},
		},
		{
			// GcpCloudFunction: the fully qualified function path, the bare
			// name serverless NEGs reference, both serving URLs, the
			// underlying Cloud Run service, runtime identity, state,
			// environment, and update time (eventarc trigger empty for HTTP
			// functions).
			name: "GcpCloudFunction",
			kind: cloudresourcekind.CloudResourceKind_GcpCloudFunction,
			rawOutputs: map[string]interface{}{
				"function_id":           "projects/prod-project/locations/us-central1/functions/hello-api",
				"function_url":          "https://us-central1-prod-project.cloudfunctions.net/hello-api",
				"service_account_email": "fn-runtime@prod-project.iam.gserviceaccount.com",
				"state":                 "ACTIVE",
				"cloud_run_service_id":  "projects/prod-project/locations/us-central1/services/hello-api",
				"eventarc_trigger_id":   "",
				"name":                  "hello-api",
				"uri":                   "https://hello-api-abc123-uc.a.run.app",
				"environment":           "GEN_2",
				"update_time":           "2026-07-05T12:00:00Z",
			},
			mustPopulate: []string{
				"function_id", "function_url", "service_account_email", "state",
				"cloud_run_service_id", "name", "uri", "environment", "update_time",
			},
		},
		{
			// GcpServiceConnectionPolicy: the fully qualified policy path, the
			// short name, the connectivity mechanism the automation reports,
			// and the change-detection etag.
			name: "GcpServiceConnectionPolicy",
			kind: cloudresourcekind.CloudResourceKind_GcpServiceConnectionPolicy,
			rawOutputs: map[string]interface{}{
				"policy_id":      "projects/prod-project/locations/us-central1/serviceConnectionPolicies/memorystore-policy",
				"name":           "memorystore-policy",
				"infrastructure": "PSC",
				"etag":           "abc123etag",
			},
			mustPopulate: []string{"policy_id", "name", "infrastructure", "etag"},
		},
		{
			// GcpMemorystoreInstance: the PSC discovery endpoint (address +
			// numeric port), the server uid, the node memory a node_type
			// implies (float), the full resource path DR secondaries
			// reference, and the backup collection automated backups land in.
			name: "GcpMemorystoreInstance",
			kind: cloudresourcekind.CloudResourceKind_GcpMemorystoreInstance,
			rawOutputs: map[string]interface{}{
				"discovery_address": "10.9.0.5",
				"discovery_port":    6379,
				"instance_uid":      "a1b2c3d4-uid",
				"node_size_gb":      1.4,
				"name":              "projects/prod-project/locations/us-central1/instances/prod-cache",
				"backup_collection": "projects/prod-project/locations/us-central1/backupCollections/col-1",
			},
			mustPopulate: []string{
				"discovery_address", "discovery_port", "instance_uid",
				"node_size_gb", "name", "backup_collection",
			},
		},
		{
			// GcpBigtableInstance: the fully qualified instance path and the
			// short name client libraries connect with.
			name: "GcpBigtableInstance",
			kind: cloudresourcekind.CloudResourceKind_GcpBigtableInstance,
			rawOutputs: map[string]interface{}{
				"instance_id":   "projects/prod-project/instances/prod-bigtable",
				"instance_name": "prod-bigtable",
			},
			mustPopulate: []string{"instance_id", "instance_name"},
		},
		{
			// GcpBigtableTable: the fully qualified table path, the short
			// name clients open, and the parent instance.
			name: "GcpBigtableTable",
			kind: cloudresourcekind.CloudResourceKind_GcpBigtableTable,
			rawOutputs: map[string]interface{}{
				"table_id":      "projects/prod-project/instances/prod-bigtable/tables/events",
				"table_name":    "events",
				"instance_name": "prod-bigtable",
			},
			mustPopulate: []string{"table_id", "table_name", "instance_name"},
		},
		{
			// GcpFirestoreDatabase: the fully qualified database path, the
			// name clients connect with, the server-generated uid, and the
			// PITR/version-retention posture timestamps.
			name: "GcpFirestoreDatabase",
			kind: cloudresourcekind.CloudResourceKind_GcpFirestoreDatabase,
			rawOutputs: map[string]interface{}{
				"database_id":              "projects/prod-project/databases/orders-db",
				"database_name":            "orders-db",
				"uid":                      "8d68546e-3c88-4244-8722-0a4b0a4b0a4b",
				"create_time":              "2026-07-05T10:00:00Z",
				"earliest_version_time":    "2026-07-05T10:00:00Z",
				"version_retention_period": "3600s",
				"key_prefix":               "",
				"update_time":              "2026-07-05T10:00:00Z",
			},
			mustPopulate: []string{
				"database_id", "database_name", "uid",
				"create_time", "earliest_version_time",
				"version_retention_period", "update_time",
			},
		},
		{
			// GcpFirestoreBackupSchedule: the server-assigned schedule id and
			// the parent database name the verifier reassembles the resource
			// path from.
			name: "GcpFirestoreBackupSchedule",
			kind: cloudresourcekind.CloudResourceKind_GcpFirestoreBackupSchedule,
			rawOutputs: map[string]interface{}{
				"schedule_id": "8d68546e-3c88-4244-8722-0a4b0a4b0a4b",
				"database":    "orders-db",
			},
			mustPopulate: []string{"schedule_id", "database"},
		},
		{
			// GcpFirestoreIndex: the server-defined index resource path and
			// the collection group it serves.
			name: "GcpFirestoreIndex",
			kind: cloudresourcekind.CloudResourceKind_GcpFirestoreIndex,
			rawOutputs: map[string]interface{}{
				"index_id":   "projects/prod-project/databases/orders-db/collectionGroups/orders/indexes/CICAgJjF6JEK",
				"collection": "orders",
			},
			mustPopulate: []string{"index_id", "collection"},
		},
		{
			// GcpDataprocCluster: the fully qualified cluster path (the
			// composition handle downstream spark-history-server references
			// consume), the short name, and the staging bucket in use.
			name: "GcpDataprocCluster",
			kind: cloudresourcekind.CloudResourceKind_GcpDataprocCluster,
			rawOutputs: map[string]interface{}{
				"cluster_id":     "projects/prod-project/regions/us-central1/clusters/etl-cluster",
				"cluster_name":   "etl-cluster",
				"staging_bucket": "dataproc-staging-us-central1-123456789012-abcdef",
			},
			mustPopulate: []string{"cluster_id", "cluster_name", "staging_bucket"},
		},
		{
			// GcpDataprocAutoscalingPolicy: the fully qualified policy path
			// (what a cluster's autoscaling_policy_uri reference resolves to)
			// plus the plain id and region.
			name: "GcpDataprocAutoscalingPolicy",
			kind: cloudresourcekind.CloudResourceKind_GcpDataprocAutoscalingPolicy,
			rawOutputs: map[string]interface{}{
				"name":      "projects/prod-project/locations/us-central1/autoscalingPolicies/batch-scaling",
				"policy_id": "batch-scaling",
				"location":  "us-central1",
			},
			mustPopulate: []string{"name", "policy_id", "location"},
		},
		{
			// GcpCloudComposerEnvironment: the fully qualified environment
			// path, the short name, and the assembled-stack handles (Airflow
			// UI, DAG bucket prefix, underlying GKE cluster).
			name: "GcpCloudComposerEnvironment",
			kind: cloudresourcekind.CloudResourceKind_GcpCloudComposerEnvironment,
			rawOutputs: map[string]interface{}{
				"environment_id":   "projects/prod-project/locations/us-central1/environments/data-pipelines",
				"environment_name": "data-pipelines",
				"airflow_uri":      "https://12345678-dot-us-central1.composer.googleusercontent.com",
				"dag_gcs_prefix":   "gs://us-central1-data-pipelines-abcdef-bucket/dags",
				"gke_cluster":      "projects/prod-project/locations/us-central1/clusters/us-central1-data-pipelines-abcdef-gke",
			},
			mustPopulate: []string{
				"environment_id", "environment_name", "airflow_uri",
				"dag_gcs_prefix", "gke_cluster",
			},
		},
		{
			// GcpCloudComposerUserWorkloadsSecret: the fully qualified secret
			// path and the Kubernetes Secret name DAGs reference. The secret
			// data is deliberately never an output.
			name: "GcpCloudComposerUserWorkloadsSecret",
			kind: cloudresourcekind.CloudResourceKind_GcpCloudComposerUserWorkloadsSecret,
			rawOutputs: map[string]interface{}{
				"name":        "projects/prod-project/locations/us-central1/environments/data-pipelines/userWorkloadsSecrets/airflow-connections",
				"secret_name": "airflow-connections",
			},
			mustPopulate: []string{"name", "secret_name"},
		},
		{
			// GcpCloudComposerUserWorkloadsConfigMap: the fully qualified
			// config-map path and the Kubernetes ConfigMap name DAGs reference.
			name: "GcpCloudComposerUserWorkloadsConfigMap",
			kind: cloudresourcekind.CloudResourceKind_GcpCloudComposerUserWorkloadsConfigMap,
			rawOutputs: map[string]interface{}{
				"name":            "projects/prod-project/locations/us-central1/environments/data-pipelines/userWorkloadsConfigMaps/dag-configuration",
				"config_map_name": "dag-configuration",
			},
			mustPopulate: []string{"name", "config_map_name"},
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
