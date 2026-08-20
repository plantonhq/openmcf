package digitaloceandatabasekafkaschemav1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestDigitalOceanDatabaseKafkaSchemaSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanDatabaseKafkaSchemaSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanDatabaseKafkaSchemaSpec validations", func() {

	newClusterRef := func(clusterId string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: clusterId},
		}
	}

	makeValidSpec := func() *DigitalOceanDatabaseKafkaSchemaSpec {
		return &DigitalOceanDatabaseKafkaSchemaSpec{
			Cluster:     newClusterRef("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
			SubjectName: "orders-value",
			SchemaType:  "avro",
			Schema:      `{"type":"record","name":"Order","fields":[{"name":"id","type":"string"}]}`,
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
					ValueFrom: &fk.ValueFromRef{Name: "my-kafka-cluster"},
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

		ginkgo.It("rejects spec with missing subject_name", func() {
			spec := makeValidSpec()
			spec.SubjectName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing schema_type", func() {
			spec := makeValidSpec()
			spec.SchemaType = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing schema", func() {
			spec := makeValidSpec()
			spec.Schema = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Schema type wall", func() {
		ginkgo.It("accepts every schema type the provider allows", func() {
			for _, schemaType := range []string{"avro", "json", "protobuf"} {
				spec := makeValidSpec()
				spec.SchemaType = schemaType
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			}
		})

		ginkgo.It("rejects an unknown schema_type", func() {
			spec := makeValidSpec()
			spec.SchemaType = "thrift"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a case-mismatched schema_type (values are case-sensitive)", func() {
			spec := makeValidSpec()
			spec.SchemaType = "AVRO"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
