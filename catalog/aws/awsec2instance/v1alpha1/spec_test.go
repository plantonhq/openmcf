package awsec2instancev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsEc2InstanceSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsEc2InstanceSpec Validation Suite")
}

func literal(v string) *fk.StringValueOrRef {
	return &fk.StringValueOrRef{LiteralOrRef: &fk.StringValueOrRef_Value{Value: v}}
}

var _ = ginkgo.Describe("AwsEc2InstanceSpec validations", func() {
	var spec *AwsEc2InstanceSpec

	ginkgo.BeforeEach(func() {
		spec = &AwsEc2InstanceSpec{
			Region:       "us-west-2",
			Ami:          "ami-0123456789abcdef0",
			InstanceType: "t4g.nano",
			SubnetId:     literal("subnet-aaa111"),
			SecurityGroupIds: []*fk.StringValueOrRef{
				literal("sg-000111222"),
			},
			InstanceProfile: literal("web-server-profile"),
		}
	})

	ginkgo.Context("launch source", func() {
		ginkgo.It("accepts a minimal ami + instance_type spec", func() {
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a launch-template-only spec", func() {
			spec.Ami = ""
			spec.InstanceType = ""
			spec.LaunchTemplate = &AwsEc2InstanceLaunchTemplate{
				Id: literal("lt-0123456789abcdef0"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails when neither ami nor launch_template is set", func() {
			spec.Ami = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when neither instance_type nor launch_template is set", func() {
			spec.InstanceType = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when ami does not start with ami-", func() {
			spec.Ami = "image-123"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when the launch template sets both id and name", func() {
			spec.LaunchTemplate = &AwsEc2InstanceLaunchTemplate{
				Id:   literal("lt-0123456789abcdef0"),
				Name: "golden-template",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts a launch template identified by name", func() {
			spec.LaunchTemplate = &AwsEc2InstanceLaunchTemplate{
				Name:    "golden-template",
				Version: "$Latest",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("region", func() {
		ginkgo.It("fails when region is empty", func() {
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("networking", func() {
		ginkgo.It("accepts a spec without subnet or security groups (default VPC)", func() {
			spec.SubnetId = nil
			spec.SecurityGroupIds = nil
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails when a primary ENI is combined with a subnet", func() {
			spec.PrimaryNetworkInterfaceId = "eni-0123456789abcdef0"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts a primary ENI with inline networking cleared", func() {
			spec.PrimaryNetworkInterfaceId = "eni-0123456789abcdef0"
			spec.SubnetId = nil
			spec.SecurityGroupIds = nil
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails when ipv6_address_count and ipv6_addresses are both set", func() {
			spec.Ipv6AddressCount = 1
			spec.Ipv6Addresses = []string{"2600:1f14::1"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a secondary interface targets card 0", func() {
			spec.SecondaryNetworkInterfaces = []*AwsEc2InstanceSecondaryNetworkInterface{{
				NetworkCardIndex: 0,
				SubnetId:         literal("subnet-bbb222"),
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a secondary interface has no subnet", func() {
			spec.SecondaryNetworkInterfaces = []*AwsEc2InstanceSecondaryNetworkInterface{{
				NetworkCardIndex: 1,
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts a valid secondary interface", func() {
			spec.SecondaryNetworkInterfaces = []*AwsEc2InstanceSecondaryNetworkInterface{{
				NetworkCardIndex:      1,
				SubnetId:              literal("subnet-bbb222"),
				PrivateIpAddressCount: 2,
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails on an invalid private DNS hostname type", func() {
			spec.PrivateDnsNameOptions = &AwsEc2InstancePrivateDnsNameOptions{HostnameType: "fqdn"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("storage", func() {
		ginkgo.It("accepts a gp3 root volume with throughput and iops", func() {
			spec.RootBlockDevice = &AwsEc2InstanceRootBlockDevice{
				VolumeSizeGb:    30,
				VolumeType:      "gp3",
				Iops:            4000,
				ThroughputMibps: 250,
				Encrypted:       true,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails on an invalid root volume type", func() {
			spec.RootBlockDevice = &AwsEc2InstanceRootBlockDevice{VolumeType: "gp4"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when throughput is set on a non-gp3 root volume", func() {
			spec.RootBlockDevice = &AwsEc2InstanceRootBlockDevice{
				VolumeType:      "gp2",
				ThroughputMibps: 250,
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when iops is set on an st1 root volume", func() {
			spec.RootBlockDevice = &AwsEc2InstanceRootBlockDevice{
				VolumeType: "st1",
				Iops:       3000,
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when an EBS device has neither size nor snapshot", func() {
			spec.EbsBlockDevices = []*AwsEc2InstanceEbsBlockDevice{{
				DeviceName: "/dev/sdf",
				VolumeType: "gp3",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts an EBS device sized by snapshot", func() {
			spec.EbsBlockDevices = []*AwsEc2InstanceEbsBlockDevice{{
				DeviceName: "/dev/sdf",
				SnapshotId: "snap-0123456789abcdef0",
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails when an ephemeral mapping sets virtual_name and no_device", func() {
			spec.EphemeralBlockDevices = []*AwsEc2InstanceEphemeralBlockDevice{{
				DeviceName:  "/dev/sdc",
				VirtualName: "ephemeral0",
				NoDevice:    true,
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("instance posture", func() {
		ginkgo.It("accepts hardened metadata options", func() {
			spec.MetadataOptions = &AwsEc2InstanceMetadataOptions{
				HttpEndpoint:            "enabled",
				HttpTokens:              "required",
				HttpPutResponseHopLimit: 2,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails on an invalid http_tokens value", func() {
			spec.MetadataOptions = &AwsEc2InstanceMetadataOptions{HttpTokens: "mandatory"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on an out-of-range hop limit", func() {
			spec.MetadataOptions = &AwsEc2InstanceMetadataOptions{HttpPutResponseHopLimit: 65}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on an invalid cpu_credits value", func() {
			spec.CpuCredits = "bursty"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on threads_per_core outside 1..2", func() {
			spec.CpuOptions = &AwsEc2InstanceCpuOptions{ThreadsPerCore: 4}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when enclaves and hibernation are both enabled", func() {
			spec.EnclaveEnabled = true
			spec.HibernationEnabled = true
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on an invalid auto_recovery value", func() {
			spec.AutoRecovery = "enabled"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on an invalid shutdown behavior", func() {
			spec.InstanceInitiatedShutdownBehavior = "hibernate"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts protection flags with explicit values", func() {
			spec.DisableApiStop = proto.Bool(true)
			spec.DisableApiTermination = proto.Bool(true)
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("purchase options", func() {
		ginkgo.It("accepts a one-time spot request", func() {
			spec.SpotOptions = &AwsEc2InstanceSpotOptions{
				SpotInstanceType: "one-time",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails when stop interruption behavior lacks a persistent request", func() {
			spec.SpotOptions = &AwsEc2InstanceSpotOptions{
				SpotInstanceType:             "one-time",
				InstanceInterruptionBehavior: "stop",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when valid_until is set on a one-time request", func() {
			spec.SpotOptions = &AwsEc2InstanceSpotOptions{
				SpotInstanceType: "one-time",
				ValidUntil:       "2027-01-01T00:00:00Z",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a capacity reservation sets preference and a target", func() {
			spec.CapacityReservation = &AwsEc2InstanceCapacityReservation{
				Preference:            "open",
				CapacityReservationId: "cr-0123456789abcdef0",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a capacity reservation targets an id and a group", func() {
			spec.CapacityReservation = &AwsEc2InstanceCapacityReservation{
				CapacityReservationId:               "cr-0123456789abcdef0",
				CapacityReservationResourceGroupArn: "arn:aws:resource-groups:us-west-2:123456789012:group/crg",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("placement", func() {
		ginkgo.It("fails on an invalid tenancy", func() {
			spec.Placement = &AwsEc2InstancePlacement{Tenancy: "shared"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when group_name and group_id are both set", func() {
			spec.Placement = &AwsEc2InstancePlacement{
				GroupName: "cluster-pg",
				GroupId:   "pg-0123456789abcdef0",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a partition number is set without a placement group", func() {
			spec.Placement = &AwsEc2InstancePlacement{PartitionNumber: 2}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts a partition within a named group", func() {
			spec.Placement = &AwsEc2InstancePlacement{
				GroupName:       "partition-pg",
				PartitionNumber: 2,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("user data", func() {
		ginkgo.It("fails when user_data and user_data_base64 are both set", func() {
			spec.UserData = "#!/bin/bash\necho hello"
			spec.UserDataBase64 = "IyEvYmluL2Jhc2gKZWNobyBoZWxsbw=="
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts user data with literal template introducers", func() {
			spec.UserData = "#!/bin/bash\necho ${HOSTNAME} > /etc/motd"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})
})
