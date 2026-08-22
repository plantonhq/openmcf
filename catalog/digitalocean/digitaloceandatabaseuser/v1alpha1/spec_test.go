package digitaloceandatabaseuserv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestDigitalOceanDatabaseUserSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanDatabaseUserSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanDatabaseUserSpec validations", func() {

	newClusterRef := func(clusterId string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: clusterId},
		}
	}

	makeValidSpec := func() *DigitalOceanDatabaseUserSpec {
		return &DigitalOceanDatabaseUserSpec{
			Cluster:  newClusterRef("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
			UserName: "app-user",
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

		ginkgo.It("rejects spec with missing user_name", func() {
			spec := makeValidSpec()
			spec.UserName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("mysql_auth_plugin", func() {
		ginkgo.It("accepts caching_sha2_password", func() {
			spec := makeValidSpec()
			spec.MysqlAuthPlugin = "caching_sha2_password"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts mysql_native_password", func() {
			spec := makeValidSpec()
			spec.MysqlAuthPlugin = "mysql_native_password"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts unset (DigitalOcean default applies)", func() {
			spec := makeValidSpec()
			spec.MysqlAuthPlugin = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown auth plugin", func() {
			spec := makeValidSpec()
			spec.MysqlAuthPlugin = "sha256_password"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Kafka ACLs", func() {
		ginkgo.It("accepts a valid topic ACL", func() {
			spec := makeValidSpec()
			spec.Settings = &DigitalOceanDatabaseUserSettings{
				KafkaAcls: []*DigitalOceanDatabaseUserKafkaAcl{
					{Topic: "events-*", Permission: "produceconsume"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects an ACL with an empty topic", func() {
			spec := makeValidSpec()
			spec.Settings = &DigitalOceanDatabaseUserSettings{
				KafkaAcls: []*DigitalOceanDatabaseUserKafkaAcl{
					{Topic: "", Permission: "consume"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an ACL with an unknown permission", func() {
			spec := makeValidSpec()
			spec.Settings = &DigitalOceanDatabaseUserSettings{
				KafkaAcls: []*DigitalOceanDatabaseUserKafkaAcl{
					{Topic: "events", Permission: "read_write"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("OpenSearch ACLs", func() {
		ginkgo.It("accepts a valid index ACL", func() {
			spec := makeValidSpec()
			spec.Settings = &DigitalOceanDatabaseUserSettings{
				OpensearchAcls: []*DigitalOceanDatabaseUserOpenSearchAcl{
					{Index: "logs-*", Permission: "readwrite"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects an ACL with an empty index", func() {
			spec := makeValidSpec()
			spec.Settings = &DigitalOceanDatabaseUserSettings{
				OpensearchAcls: []*DigitalOceanDatabaseUserOpenSearchAcl{
					{Index: "", Permission: "read"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an ACL with an unknown permission", func() {
			spec := makeValidSpec()
			spec.Settings = &DigitalOceanDatabaseUserSettings{
				OpensearchAcls: []*DigitalOceanDatabaseUserOpenSearchAcl{
					{Index: "logs", Permission: "all"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Settings message", func() {
		ginkgo.It("accepts an empty settings message", func() {
			spec := makeValidSpec()
			spec.Settings = &DigitalOceanDatabaseUserSettings{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})
})
