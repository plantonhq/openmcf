package gcpgkeclusterv1alpha1

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"

	"buf.build/go/protovalidate"
	"github.com/plantonhq/planton/apis/dev/planton/shared"
)

func TestGcpGkeClusterSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "GcpGkeClusterSpec Custom Validation Tests")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func ref(name string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
			ValueFrom: &foreignkeyv1.ValueFromRef{Name: name},
		},
	}
}

// minimalSpec returns a valid baseline spec every case mutates from.
func minimalSpec() *GcpGkeClusterSpec {
	return &GcpGkeClusterSpec{
		ProjectId:  literal("test-project-123"),
		Location:   "us-central1",
		Network:    literal("projects/test-project-123/global/networks/test-vpc"),
		Subnetwork: literal("projects/test-project-123/regions/us-central1/subnetworks/test-subnet"),
	}
}

func newCluster(spec *GcpGkeClusterSpec) *GcpGkeCluster {
	return &GcpGkeCluster{
		ApiVersion: "gcp.planton.dev/v1alpha1",
		Kind:       "GcpGkeCluster",
		Metadata: &shared.CloudResourceMetadata{
			Name: "test-gke-cluster",
		},
		Spec: spec,
	}
}

var _ = ginkgo.Describe("GcpGkeClusterSpec Custom Validation Tests", func() {

	ginkgo.Describe("valid configurations", func() {

		ginkgo.It("accepts a minimal regional cluster", func() {
			err := protovalidate.Validate(newCluster(minimalSpec()))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a zonal location", func() {
			spec := minimalSpec()
			spec.Location = "us-central1-a"
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a multi-digit region and zone", func() {
			spec := minimalSpec()
			spec.Location = "europe-west12"
			gomega.Expect(protovalidate.Validate(newCluster(spec))).To(gomega.BeNil())
			spec.Location = "me-central2-a"
			gomega.Expect(protovalidate.Validate(newCluster(spec))).To(gomega.BeNil())
		})

		ginkgo.It("accepts network and subnetwork by reference", func() {
			spec := minimalSpec()
			spec.Network = ref("my-vpc")
			spec.Subnetwork = ref("my-subnet")
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts an explicit cluster_name", func() {
			spec := minimalSpec()
			spec.ClusterName = "prod-primary"
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts named secondary ranges", func() {
			spec := minimalSpec()
			spec.IpAllocation = &GcpGkeClusterIpAllocation{
				ClusterSecondaryRangeName:  literal("pods-range"),
				ServicesSecondaryRangeName: literal("services-range"),
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts GKE-managed CIDR blocks", func() {
			spec := minimalSpec()
			spec.IpAllocation = &GcpGkeClusterIpAllocation{
				ClusterIpv4CidrBlock:  "10.4.0.0/14",
				ServicesIpv4CidrBlock: "10.8.0.0/20",
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a private cluster with a /28 master range", func() {
			spec := minimalSpec()
			spec.PrivateCluster = &GcpGkeClusterPrivateCluster{
				EnablePrivateNodes:  true,
				MasterIpv4CidrBlock: "172.16.0.16/28",
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a private-only endpoint with private nodes", func() {
			spec := minimalSpec()
			spec.PrivateCluster = &GcpGkeClusterPrivateCluster{
				EnablePrivateNodes:    true,
				EnablePrivateEndpoint: true,
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts an Autopilot cluster", func() {
			spec := minimalSpec()
			spec.EnableAutopilot = true
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts Autopilot with allow_net_admin", func() {
			spec := minimalSpec()
			spec.EnableAutopilot = true
			spec.AllowNetAdmin = true
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a daily maintenance window", func() {
			spec := minimalSpec()
			spec.MaintenancePolicy = &GcpGkeClusterMaintenancePolicy{
				DailyWindow: &GcpGkeClusterDailyMaintenanceWindow{StartTime: "03:00"},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a recurring window with exclusions", func() {
			spec := minimalSpec()
			spec.MaintenancePolicy = &GcpGkeClusterMaintenancePolicy{
				RecurringWindow: &GcpGkeClusterRecurringMaintenanceWindow{
					StartTime:  "2026-01-01T02:00:00Z",
					EndTime:    "2026-01-01T06:00:00Z",
					Recurrence: "FREQ=WEEKLY;BYDAY=SA,SU",
				},
				Exclusions: []*GcpGkeClusterMaintenanceExclusion{{
					ExclusionName: "year-end-freeze",
					StartTime:     "2026-12-15T00:00:00Z",
					EndTime:       "2027-01-05T00:00:00Z",
					Scope:         "NO_MINOR_UPGRADES",
				}},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts NAP with resource limits", func() {
			spec := minimalSpec()
			spec.ClusterAutoscaling = &GcpGkeClusterAutoscaling{
				Enabled: true,
				ResourceLimits: []*GcpGkeClusterAutoscalingResourceLimit{
					{ResourceType: "cpu", Minimum: 4, Maximum: 64},
					{ResourceType: "memory", Minimum: 16, Maximum: 256},
				},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts NAP defaults with a service-account reference", func() {
			spec := minimalSpec()
			spec.ClusterAutoscaling = &GcpGkeClusterAutoscaling{
				Enabled: true,
				ResourceLimits: []*GcpGkeClusterAutoscalingResourceLimit{
					{ResourceType: "cpu", Maximum: 32},
				},
				AutoProvisioningDefaults: &GcpGkeClusterAutoProvisioningDefaults{
					ServiceAccount: ref("gke-nodes-sa"),
					DiskType:       "pd-balanced",
				},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts CMEK database encryption with a key reference", func() {
			spec := minimalSpec()
			spec.DatabaseEncryption = &GcpGkeClusterDatabaseEncryption{
				State:   "ENCRYPTED",
				KeyName: ref("etcd-secrets-key"),
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts DECRYPTED state without a key", func() {
			spec := minimalSpec()
			spec.DatabaseEncryption = &GcpGkeClusterDatabaseEncryption{State: "DECRYPTED"}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts master authorized networks", func() {
			spec := minimalSpec()
			spec.MasterAuthorizedNetworks = &GcpGkeClusterMasterAuthorizedNetworks{
				CidrBlocks: []*GcpGkeClusterMasterAuthorizedNetworkCidr{
					{CidrBlock: "203.0.113.0/24", DisplayName: "office"},
				},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts logging and monitoring component lists", func() {
			spec := minimalSpec()
			spec.Logging = &GcpGkeClusterLogging{
				Components: []string{"SYSTEM_COMPONENTS", "WORKLOADS", "APISERVER"},
			}
			spec.Monitoring = &GcpGkeClusterMonitoring{
				Components:          []string{"SYSTEM_COMPONENTS", "APISERVER"},
				AutoMonitoringScope: "ALL",
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts notification pub/sub with a topic reference", func() {
			spec := minimalSpec()
			spec.NotificationPubsub = &GcpGkeClusterNotificationPubSub{
				Enabled:    true,
				Topic:      ref("gke-events"),
				EventTypes: []string{"UPGRADE_EVENT", "SECURITY_BULLETIN_EVENT"},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a resource usage export destination", func() {
			spec := minimalSpec()
			spec.ResourceUsageExport = &GcpGkeClusterResourceUsageExport{
				BigqueryDatasetId: literal("gke_usage"),
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts the full networking surface", func() {
			spec := minimalSpec()
			spec.DatapathProvider = "ADVANCED_DATAPATH"
			spec.EnableFqdnNetworkPolicy = true
			spec.EnableCiliumClusterwideNetworkPolicy = true
			spec.GatewayApiChannel = "CHANNEL_STANDARD"
			spec.PrivateIpv6GoogleAccess = "PRIVATE_IPV6_GOOGLE_ACCESS_TO_GOOGLE"
			spec.InTransitEncryption = "IN_TRANSIT_ENCRYPTION_INTER_NODE_TRANSPARENT"
			spec.TotalEgressBandwidthTier = "TIER_1"
			spec.DnsConfig = &GcpGkeClusterDnsConfig{
				ClusterDns:      "CLOUD_DNS",
				ClusterDnsScope: "VPC_SCOPE",
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts node_locations zones", func() {
			spec := minimalSpec()
			spec.NodeLocations = []string{"us-central1-a", "us-central1-b", "europe-west12-a"}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts confidential nodes with an instance type", func() {
			spec := minimalSpec()
			spec.ConfidentialNodes = &GcpGkeClusterConfidentialNodes{
				Enabled:                  true,
				ConfidentialInstanceType: "SEV_SNP",
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("required fields and formats", func() {

		ginkgo.It("rejects a missing location", func() {
			spec := minimalSpec()
			spec.Location = ""
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid location format", func() {
			spec := minimalSpec()
			spec.Location = "US-CENTRAL1"
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing network", func() {
			spec := minimalSpec()
			spec.Network = nil
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing subnetwork", func() {
			spec := minimalSpec()
			spec.Subnetwork = nil
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid cluster_name", func() {
			spec := minimalSpec()
			spec.ClusterName = "Prod_Cluster"
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an over-long cluster_name", func() {
			spec := minimalSpec()
			spec.ClusterName = "a-very-long-cluster-name-that-exceeds-forty-characters"
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid node_locations zone", func() {
			spec := minimalSpec()
			spec.NodeLocations = []string{"us-central1"}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects default_max_pods_per_node out of range", func() {
			spec := minimalSpec()
			four := int32(4)
			spec.DefaultMaxPodsPerNode = &four
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid datapath_provider", func() {
			spec := minimalSpec()
			spec.DatapathProvider = "EBPF"
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid gateway_api_channel", func() {
			spec := minimalSpec()
			spec.GatewayApiChannel = "CHANNEL_EXPERIMENTAL"
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid logging component", func() {
			spec := minimalSpec()
			spec.Logging = &GcpGkeClusterLogging{Components: []string{"EVERYTHING"}}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid notification event type", func() {
			spec := minimalSpec()
			spec.NotificationPubsub = &GcpGkeClusterNotificationPubSub{
				Enabled:    true,
				Topic:      literal("projects/p/topics/t"),
				EventTypes: []string{"EVERY_EVENT"},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid binary authorization mode", func() {
			spec := minimalSpec()
			spec.BinaryAuthorizationEvaluationMode = "ENFORCE_EVERYTHING"
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid hpa_profile", func() {
			spec := minimalSpec()
			spec.HpaProfile = "FAST"
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("ip_allocation coherence", func() {

		ginkgo.It("rejects a cluster range name combined with a cluster CIDR block", func() {
			spec := minimalSpec()
			spec.IpAllocation = &GcpGkeClusterIpAllocation{
				ClusterSecondaryRangeName: literal("pods-range"),
				ClusterIpv4CidrBlock:      "10.4.0.0/14",
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a services range name combined with a services CIDR block", func() {
			spec := minimalSpec()
			spec.IpAllocation = &GcpGkeClusterIpAllocation{
				ServicesSecondaryRangeName: literal("services-range"),
				ServicesIpv4CidrBlock:      "10.8.0.0/20",
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid stack_type", func() {
			spec := minimalSpec()
			ipv6 := "IPV6"
			spec.IpAllocation = &GcpGkeClusterIpAllocation{StackType: &ipv6}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("private cluster coherence", func() {

		ginkgo.It("rejects a private endpoint without private nodes", func() {
			spec := minimalSpec()
			spec.PrivateCluster = &GcpGkeClusterPrivateCluster{
				EnablePrivateEndpoint: true,
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a master CIDR combined with an endpoint subnetwork", func() {
			spec := minimalSpec()
			spec.PrivateCluster = &GcpGkeClusterPrivateCluster{
				EnablePrivateNodes:        true,
				MasterIpv4CidrBlock:       "172.16.0.16/28",
				PrivateEndpointSubnetwork: literal("projects/p/regions/us-central1/subnetworks/cp"),
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a non-/28 master range", func() {
			spec := minimalSpec()
			spec.PrivateCluster = &GcpGkeClusterPrivateCluster{
				EnablePrivateNodes:  true,
				MasterIpv4CidrBlock: "172.16.0.0/24",
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("maintenance policy coherence", func() {

		ginkgo.It("rejects a policy without any window", func() {
			spec := minimalSpec()
			spec.MaintenancePolicy = &GcpGkeClusterMaintenancePolicy{}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects both daily and recurring windows", func() {
			spec := minimalSpec()
			spec.MaintenancePolicy = &GcpGkeClusterMaintenancePolicy{
				DailyWindow: &GcpGkeClusterDailyMaintenanceWindow{StartTime: "03:00"},
				RecurringWindow: &GcpGkeClusterRecurringMaintenanceWindow{
					StartTime:  "2026-01-01T02:00:00Z",
					EndTime:    "2026-01-01T06:00:00Z",
					Recurrence: "FREQ=DAILY",
				},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid daily start_time", func() {
			spec := minimalSpec()
			spec.MaintenancePolicy = &GcpGkeClusterMaintenancePolicy{
				DailyWindow: &GcpGkeClusterDailyMaintenanceWindow{StartTime: "25:00"},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid exclusion scope", func() {
			spec := minimalSpec()
			spec.MaintenancePolicy = &GcpGkeClusterMaintenancePolicy{
				DailyWindow: &GcpGkeClusterDailyMaintenanceWindow{StartTime: "03:00"},
				Exclusions: []*GcpGkeClusterMaintenanceExclusion{{
					ExclusionName: "freeze",
					StartTime:     "2026-12-15T00:00:00Z",
					EndTime:       "2027-01-05T00:00:00Z",
					Scope:         "NO_TOUCHING",
				}},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("cluster autoscaling coherence", func() {

		ginkgo.It("rejects enabled NAP without resource limits", func() {
			spec := minimalSpec()
			spec.ClusterAutoscaling = &GcpGkeClusterAutoscaling{Enabled: true}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects resource limits on a disabled NAP", func() {
			spec := minimalSpec()
			spec.ClusterAutoscaling = &GcpGkeClusterAutoscaling{
				ResourceLimits: []*GcpGkeClusterAutoscalingResourceLimit{
					{ResourceType: "cpu", Maximum: 32},
				},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a limit without a maximum", func() {
			spec := minimalSpec()
			spec.ClusterAutoscaling = &GcpGkeClusterAutoscaling{
				Enabled: true,
				ResourceLimits: []*GcpGkeClusterAutoscalingResourceLimit{
					{ResourceType: "cpu"},
				},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid NAP disk type", func() {
			spec := minimalSpec()
			spec.ClusterAutoscaling = &GcpGkeClusterAutoscaling{
				Enabled: true,
				ResourceLimits: []*GcpGkeClusterAutoscalingResourceLimit{
					{ResourceType: "cpu", Maximum: 32},
				},
				AutoProvisioningDefaults: &GcpGkeClusterAutoProvisioningDefaults{
					DiskType: "local-ssd",
				},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("autopilot conflict rules", func() {

		ginkgo.It("rejects Autopilot with cluster_autoscaling", func() {
			spec := minimalSpec()
			spec.EnableAutopilot = true
			spec.ClusterAutoscaling = &GcpGkeClusterAutoscaling{
				Enabled: true,
				ResourceLimits: []*GcpGkeClusterAutoscalingResourceLimit{
					{ResourceType: "cpu", Maximum: 32},
				},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects Autopilot with default_max_pods_per_node", func() {
			spec := minimalSpec()
			spec.EnableAutopilot = true
			maxPods := int32(64)
			spec.DefaultMaxPodsPerNode = &maxPods
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects Autopilot with intranode visibility", func() {
			spec := minimalSpec()
			spec.EnableAutopilot = true
			spec.EnableIntranodeVisibility = true
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects Autopilot with Calico network policy", func() {
			spec := minimalSpec()
			spec.EnableAutopilot = true
			spec.EnableNetworkPolicy = true
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects Autopilot with an explicit shielded-nodes flag", func() {
			spec := minimalSpec()
			spec.EnableAutopilot = true
			shielded := true
			spec.EnableShieldedNodes = &shielded
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects Autopilot with the dns-cache addon", func() {
			spec := minimalSpec()
			spec.EnableAutopilot = true
			spec.Addons = &GcpGkeClusterAddons{DnsCacheEnabled: true}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects allow_net_admin on a Standard cluster", func() {
			spec := minimalSpec()
			spec.AllowNetAdmin = true
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("security block coherence", func() {

		ginkgo.It("rejects ENCRYPTED state without a key", func() {
			spec := minimalSpec()
			spec.DatabaseEncryption = &GcpGkeClusterDatabaseEncryption{State: "ENCRYPTED"}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid encryption state", func() {
			spec := minimalSpec()
			spec.DatabaseEncryption = &GcpGkeClusterDatabaseEncryption{State: "MAYBE"}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid security posture mode", func() {
			spec := minimalSpec()
			spec.SecurityPosture = &GcpGkeClusterSecurityPosture{Mode: "PARANOID"}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid confidential instance type", func() {
			spec := minimalSpec()
			spec.ConfidentialNodes = &GcpGkeClusterConfidentialNodes{
				Enabled:                  true,
				ConfidentialInstanceType: "SGX",
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid anonymous authentication mode", func() {
			spec := minimalSpec()
			spec.AnonymousAuthenticationMode = "DISABLED"
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("dns and notification coherence", func() {

		ginkgo.It("rejects an additive VPC scope domain without Cloud DNS", func() {
			spec := minimalSpec()
			spec.DnsConfig = &GcpGkeClusterDnsConfig{
				AdditiveVpcScopeDnsDomain: "svc.example.internal",
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a DNS scope without Cloud DNS", func() {
			spec := minimalSpec()
			spec.DnsConfig = &GcpGkeClusterDnsConfig{ClusterDnsScope: "VPC_SCOPE"}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects enabled notifications without a topic", func() {
			spec := minimalSpec()
			spec.NotificationPubsub = &GcpGkeClusterNotificationPubSub{Enabled: true}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a usage export without a dataset", func() {
			spec := minimalSpec()
			spec.ResourceUsageExport = &GcpGkeClusterResourceUsageExport{}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid master authorized CIDR", func() {
			spec := minimalSpec()
			spec.MasterAuthorizedNetworks = &GcpGkeClusterMasterAuthorizedNetworks{
				CidrBlocks: []*GcpGkeClusterMasterAuthorizedNetworkCidr{
					{CidrBlock: "office-network"},
				},
			}
			err := protovalidate.Validate(newCluster(spec))
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
