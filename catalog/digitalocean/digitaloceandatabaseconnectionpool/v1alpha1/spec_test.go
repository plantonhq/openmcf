package digitaloceandatabaseconnectionpoolv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestDigitalOceanDatabaseConnectionPoolSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanDatabaseConnectionPoolSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanDatabaseConnectionPoolSpec validations", func() {

	newClusterRef := func(clusterId string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: clusterId},
		}
	}

	makeValidSpec := func() *DigitalOceanDatabaseConnectionPoolSpec {
		return &DigitalOceanDatabaseConnectionPoolSpec{
			Cluster:  newClusterRef("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
			PoolName: "app-pool",
			Mode:     "transaction",
			Size:     10,
			DbName:   "defaultdb",
			User:     "app-user",
		}
	}

	ginkgo.Context("Required fields", func() {
		ginkgo.It("accepts a valid spec with a user", func() {
			err := protovalidate.Validate(makeValidSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts an inbound-user pool (user omitted)", func() {
			spec := makeValidSpec()
			spec.User = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing cluster", func() {
			spec := makeValidSpec()
			spec.Cluster = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing pool_name", func() {
			spec := makeValidSpec()
			spec.PoolName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing mode", func() {
			spec := makeValidSpec()
			spec.Mode = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing db_name", func() {
			spec := makeValidSpec()
			spec.DbName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("pool_name bounds", func() {
		ginkgo.It("rejects a name shorter than 3 characters", func() {
			spec := makeValidSpec()
			spec.PoolName = "ab"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a name longer than 63 characters", func() {
			spec := makeValidSpec()
			spec.PoolName = "a123456789b123456789c123456789d123456789e123456789f123456789g123"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("mode", func() {
		ginkgo.It("accepts session mode", func() {
			spec := makeValidSpec()
			spec.Mode = "session"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts statement mode", func() {
			spec := makeValidSpec()
			spec.Mode = "statement"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown mode", func() {
			spec := makeValidSpec()
			spec.Mode = "pooled"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("size", func() {
		ginkgo.It("rejects a zero size", func() {
			spec := makeValidSpec()
			spec.Size = 0
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a negative size", func() {
			spec := makeValidSpec()
			spec.Size = -5
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
