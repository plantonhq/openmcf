package gcpalloydbinstancev1alpha1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/shared"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "GcpAlloydbInstanceSpec Suite")
}

func litRef(v string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: v},
	}
}

func fromRef(name string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
			ValueFrom: &foreignkeyv1.ValueFromRef{Name: name},
		},
	}
}

func ptr(s string) *string {
	return &s
}

var _ = ginkgo.Describe("GcpAlloydbInstanceSpec", func() {
	var validator protovalidate.Validator

	ginkgo.BeforeEach(func() {
		var err error
		validator, err = protovalidate.New()
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	minimalReadPool := func() *GcpAlloydbInstance {
		return &GcpAlloydbInstance{
			ApiVersion: "gcp.planton.dev/v1alpha1",
			Kind:       "GcpAlloydbInstance",
			Metadata: &shared.CloudResourceMetadata{
				Name: "test-read-pool",
			},
			Spec: &GcpAlloydbInstanceSpec{
				Cluster:    litRef("projects/p/locations/us-central1/clusters/orders"),
				InstanceId: "orders-read-pool",
				ReadPoolConfig: &GcpAlloydbInstanceReadPoolConfig{
					NodeCount: 1,
				},
			},
		}
	}

	expectValid := func(r *GcpAlloydbInstance) {
		gomega.Expect(validator.Validate(r)).To(gomega.Succeed())
	}

	expectInvalid := func(r *GcpAlloydbInstance, substr string) {
		err := validator.Validate(r)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), substr)).To(
			gomega.BeTrue(), "expected error to contain %q, got: %s", substr, err)
	}

	ginkgo.It("accepts a minimal read pool instance", func() {
		expectValid(minimalReadPool())
	})

	ginkgo.It("accepts an explicit READ_POOL type", func() {
		r := minimalReadPool()
		r.Spec.InstanceType = ptr("READ_POOL")
		expectValid(r)
	})

	ginkgo.It("rejects a missing cluster", func() {
		r := minimalReadPool()
		r.Spec.Cluster = nil
		expectInvalid(r, "cluster")
	})

	ginkgo.It("accepts a cluster reference", func() {
		r := minimalReadPool()
		r.Spec.Cluster = fromRef("orders-cluster")
		expectValid(r)
	})

	ginkgo.It("rejects a missing instance_id", func() {
		r := minimalReadPool()
		r.Spec.InstanceId = ""
		expectInvalid(r, "instance_id")
	})

	ginkgo.It("rejects an unknown instance_type", func() {
		r := minimalReadPool()
		r.Spec.InstanceType = ptr("REPLICA")
		expectInvalid(r, "instance_type")
	})

	ginkgo.It("rejects a read pool without node_count", func() {
		r := minimalReadPool()
		r.Spec.ReadPoolConfig = nil
		expectInvalid(r, "node_count")
	})

	ginkgo.It("rejects read_pool_config on a PRIMARY instance", func() {
		r := minimalReadPool()
		r.Spec.InstanceType = ptr("PRIMARY")
		expectInvalid(r, "read_pool_config applies")
	})

	ginkgo.It("rejects both cpu_count and machine_type", func() {
		r := minimalReadPool()
		r.Spec.CpuCount = 2
		r.Spec.MachineType = "n2-highmem-2"
		expectInvalid(r, "only one of cpu_count or machine_type")
	})

	ginkgo.It("accepts cpu_count alone", func() {
		r := minimalReadPool()
		r.Spec.CpuCount = 2
		expectValid(r)
	})

	ginkgo.It("rejects authorized networks without public IP", func() {
		r := minimalReadPool()
		r.Spec.AuthorizedExternalNetworks = []*GcpAlloydbInstanceAuthorizedExternalNetwork{
			{CidrRange: "203.0.113.0/24"},
		}
		expectInvalid(r, "enable_public_ip")
	})

	ginkgo.It("accepts authorized networks with public IP", func() {
		r := minimalReadPool()
		r.Spec.EnablePublicIp = true
		r.Spec.AuthorizedExternalNetworks = []*GcpAlloydbInstanceAuthorizedExternalNetwork{
			{CidrRange: "203.0.113.0/24"},
		}
		expectValid(r)
	})

	ginkgo.It("rejects an invalid ssl_mode", func() {
		r := minimalReadPool()
		r.Spec.SslMode = "PLAINTEXT"
		expectInvalid(r, "ssl_mode")
	})

	ginkgo.It("rejects a wrong kind constant", func() {
		r := minimalReadPool()
		r.Kind = "GcpAlloydbCluster"
		expectInvalid(r, "kind")
	})
})
