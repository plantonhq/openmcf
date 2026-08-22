package digitaloceandatabasereplicav1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	digitalocean "github.com/plantonhq/planton/catalog/digitalocean"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestDigitalOceanDatabaseReplicaSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanDatabaseReplicaSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanDatabaseReplicaSpec validations", func() {

	newRef := func(value string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: value},
		}
	}

	makeValidSpec := func() *DigitalOceanDatabaseReplicaSpec {
		return &DigitalOceanDatabaseReplicaSpec{
			Cluster:     newRef("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
			ReplicaName: "read-replica",
			Region:      digitalocean.DigitalOceanRegion_nyc3,
			Size:        "db-s-1vcpu-1gb",
		}
	}

	ginkgo.Context("Required fields", func() {
		ginkgo.It("accepts a minimal valid spec", func() {
			err := protovalidate.Validate(makeValidSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing cluster", func() {
			spec := makeValidSpec()
			spec.Cluster = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing replica_name", func() {
			spec := makeValidSpec()
			spec.ReplicaName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing region", func() {
			spec := makeValidSpec()
			spec.Region = digitalocean.DigitalOceanRegion_digital_ocean_region_unspecified
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing size", func() {
			spec := makeValidSpec()
			spec.Size = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Optional fields", func() {
		ginkgo.It("accepts a spec with vpc, storage, and tags", func() {
			spec := makeValidSpec()
			spec.Vpc = newRef("bbbbbbbb-cccc-dddd-eeee-ffffffffffff")
			spec.StorageSizeMib = 30720
			spec.Tags = []string{"env:prod", "replica"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects a tag with forbidden characters", func() {
			spec := makeValidSpec()
			spec.Tags = []string{"has spaces"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
