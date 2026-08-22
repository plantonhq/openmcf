package digitaloceandatabasedbv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestDigitalOceanDatabaseDbSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanDatabaseDbSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanDatabaseDbSpec validations", func() {

	newClusterRef := func(clusterId string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: clusterId},
		}
	}

	makeValidSpec := func() *DigitalOceanDatabaseDbSpec {
		return &DigitalOceanDatabaseDbSpec{
			Cluster:      newClusterRef("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
			DatabaseName: "orders",
		}
	}

	ginkgo.Context("Required fields", func() {
		ginkgo.It("accepts a minimal valid spec", func() {
			err := protovalidate.Validate(makeValidSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a cluster reference by name", func() {
			spec := makeValidSpec()
			spec.Cluster = &fk.StringValueOrRef{
				LiteralOrRef: &fk.StringValueOrRef_ValueFrom{
					ValueFrom: &fk.ValueFromRef{Name: "my-cluster"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing cluster", func() {
			spec := makeValidSpec()
			spec.Cluster = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing database_name", func() {
			spec := makeValidSpec()
			spec.DatabaseName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
