package gcpserverlessvpcconnectorv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestGcpServerlessVpcConnectorSpec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GcpServerlessVpcConnectorSpec Validation Suite")
}

var _ = Describe("GcpServerlessVpcConnectorSpec validations", func() {

	strVal := func(v string) *foreignkeyv1.StringValueOrRef {
		return &foreignkeyv1.StringValueOrRef{
			LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: v},
		}
	}

	strRef := func(name string) *foreignkeyv1.StringValueOrRef {
		return &foreignkeyv1.StringValueOrRef{
			LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
				ValueFrom: &foreignkeyv1.ValueFromRef{Name: name},
			},
		}
	}

	// Network placement is the baseline valid shape.
	makeValidSpec := func() *GcpServerlessVpcConnectorSpec {
		return &GcpServerlessVpcConnectorSpec{
			ProjectId:   strVal("my-gcp-project"),
			Region:      "us-central1",
			Network:     strVal("my-vpc"),
			IpCidrRange: "10.8.0.0/28",
		}
	}

	makeValidSubnetSpec := func() *GcpServerlessVpcConnectorSpec {
		return &GcpServerlessVpcConnectorSpec{
			ProjectId: strVal("my-gcp-project"),
			Region:    "us-central1",
			Subnet: &GcpServerlessVpcConnectorSubnet{
				Name: strVal("connector-subnet"),
			},
		}
	}

	Context("Required fields", func() {
		It("accepts a minimal network-placement spec", func() {
			Expect(protovalidate.Validate(makeValidSpec())).To(BeNil())
		})

		It("accepts a minimal subnet-placement spec", func() {
			Expect(protovalidate.Validate(makeValidSubnetSpec())).To(BeNil())
		})

		It("accepts a spec without project_id (ambient project)", func() {
			spec := makeValidSpec()
			spec.ProjectId = nil
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects spec with missing region", func() {
			spec := makeValidSpec()
			spec.Region = ""
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Region validation", func() {
		It("accepts a standard region", func() {
			spec := makeValidSpec()
			spec.Region = "europe-west1"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a multi-digit region", func() {
			spec := makeValidSpec()
			spec.Region = "europe-west12"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects a zone as region", func() {
			spec := makeValidSpec()
			spec.Region = "us-central1-a"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Connector name validation", func() {
		It("accepts a valid connector name", func() {
			spec := makeValidSpec()
			spec.ConnectorName = "svc-egress"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a 25-character name (the GCP maximum)", func() {
			spec := makeValidSpec()
			spec.ConnectorName = "a234567890123456789012345"
			Expect(len(spec.ConnectorName)).To(Equal(25))
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects a 26-character name", func() {
			spec := makeValidSpec()
			spec.ConnectorName = "a2345678901234567890123456"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects an uppercase name", func() {
			spec := makeValidSpec()
			spec.ConnectorName = "MyConnector"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a name ending with a hyphen", func() {
			spec := makeValidSpec()
			spec.ConnectorName = "connector-"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Placement exclusivity (network XOR subnet)", func() {
		It("rejects a spec with neither network nor subnet", func() {
			spec := makeValidSpec()
			spec.Network = nil
			spec.IpCidrRange = ""
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a spec with both network and subnet", func() {
			spec := makeValidSpec()
			spec.Subnet = &GcpServerlessVpcConnectorSubnet{Name: strVal("connector-subnet")}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts a network placement expressed as a reference", func() {
			spec := makeValidSpec()
			spec.Network = strRef("my-vpc-resource")
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects a network reference with both network and subnet set", func() {
			spec := makeValidSpec()
			spec.Network = strRef("my-vpc-resource")
			spec.Subnet = &GcpServerlessVpcConnectorSubnet{Name: strVal("connector-subnet")}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Network placement CIDR coherence", func() {
		It("rejects network placement without ip_cidr_range", func() {
			spec := makeValidSpec()
			spec.IpCidrRange = ""
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects ip_cidr_range on subnet placement", func() {
			spec := makeValidSubnetSpec()
			spec.IpCidrRange = "10.8.0.0/28"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a non-/28 range", func() {
			spec := makeValidSpec()
			spec.IpCidrRange = "10.8.0.0/24"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a malformed CIDR", func() {
			spec := makeValidSpec()
			spec.IpCidrRange = "10.8.0/28"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Subnet placement", func() {
		It("rejects a subnet block without a name", func() {
			spec := makeValidSubnetSpec()
			spec.Subnet.Name = nil
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts a subnet name expressed as a reference", func() {
			spec := makeValidSubnetSpec()
			spec.Subnet.Name = strRef("my-connector-subnet")
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a Shared-VPC host project on the subnet", func() {
			spec := makeValidSubnetSpec()
			spec.Subnet.ProjectId = "shared-vpc-host-project"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})
	})

	Context("Machine type validation", func() {
		It("accepts each supported machine type", func() {
			for _, mt := range []string{"f1-micro", "e2-micro", "e2-standard-4"} {
				spec := makeValidSpec()
				spec.MachineType = mt
				Expect(protovalidate.Validate(spec)).To(BeNil())
			}
		})

		It("rejects an unsupported machine type", func() {
			spec := makeValidSpec()
			spec.MachineType = "n1-standard-1"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Instance scaling bounds", func() {
		It("accepts valid min/max instances", func() {
			spec := makeValidSpec()
			spec.MinInstances = proto.Int32(2)
			spec.MaxInstances = proto.Int32(10)
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects min_instances below 2", func() {
			spec := makeValidSpec()
			spec.MinInstances = proto.Int32(1)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects min_instances above 9", func() {
			spec := makeValidSpec()
			spec.MinInstances = proto.Int32(10)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects max_instances below 3", func() {
			spec := makeValidSpec()
			spec.MaxInstances = proto.Int32(2)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects max_instances above 10", func() {
			spec := makeValidSpec()
			spec.MaxInstances = proto.Int32(11)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects min_instances equal to max_instances", func() {
			spec := makeValidSpec()
			spec.MinInstances = proto.Int32(5)
			spec.MaxInstances = proto.Int32(5)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects min_instances greater than max_instances", func() {
			spec := makeValidSpec()
			spec.MinInstances = proto.Int32(8)
			spec.MaxInstances = proto.Int32(3)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts min alone and max alone", func() {
			spec := makeValidSpec()
			spec.MinInstances = proto.Int32(3)
			Expect(protovalidate.Validate(spec)).To(BeNil())

			spec = makeValidSpec()
			spec.MaxInstances = proto.Int32(4)
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})
	})
})
