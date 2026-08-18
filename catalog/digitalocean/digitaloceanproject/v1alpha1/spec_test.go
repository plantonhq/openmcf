package digitaloceanprojectv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestDigitalOceanProjectSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanProjectSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanProjectSpec validations", func() {

	newRef := func(value string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: value},
		}
	}

	makeValidSpec := func() *DigitalOceanProjectSpec {
		return &DigitalOceanProjectSpec{
			ProjectName: "web-production",
		}
	}

	ginkgo.Context("Required fields", func() {
		ginkgo.It("accepts a minimal valid spec (name only)", func() {
			err := protovalidate.Validate(makeValidSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing project_name", func() {
			spec := makeValidSpec()
			spec.ProjectName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects project_name longer than 175 characters", func() {
			spec := makeValidSpec()
			long := make([]byte, 176)
			for i := range long {
				long[i] = 'a'
			}
			spec.ProjectName = string(long)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Purpose", func() {
		ginkgo.It("accepts a standard purpose", func() {
			spec := makeValidSpec()
			spec.Purpose = "Web Application"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts free-text purpose (canonicalized by the API)", func() {
			spec := makeValidSpec()
			spec.Purpose = "My Basic Web App"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects a purpose starting with Other: (the permanent-diff trap)", func() {
			spec := makeValidSpec()
			spec.Purpose = "Other: something"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects purpose longer than 255 characters", func() {
			spec := makeValidSpec()
			long := make([]byte, 256)
			for i := range long {
				long[i] = 'p'
			}
			spec.Purpose = string(long)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Environment", func() {
		ginkgo.It("accepts each of the three lowercase environments", func() {
			for _, env := range []string{"development", "staging", "production"} {
				spec := makeValidSpec()
				spec.Environment = env
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil(), "environment %q should validate", env)
			}
		})

		ginkgo.It("accepts an unset environment", func() {
			spec := makeValidSpec()
			spec.Environment = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects the capitalized spelling (lowercase is canonical here)", func() {
			spec := makeValidSpec()
			spec.Environment = "Production"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown environment", func() {
			spec := makeValidSpec()
			spec.Environment = "qa"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Resources", func() {
		ginkgo.It("accepts literal URNs across families", func() {
			spec := makeValidSpec()
			spec.Resources = []*fk.StringValueOrRef{
				newRef("do:droplet:123456"),
				newRef("do:dbaas:6ec9c684-aaaa-bbbb-cccc-9ffb47dd4f8f"),
				newRef("do:space:my-bucket"),
				newRef("do:domain:example.com"),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	ginkgo.Context("Full spec", func() {
		ginkgo.It("accepts a fully populated spec", func() {
			spec := &DigitalOceanProjectSpec{
				ProjectName: "web-production",
				Description: "Production web workloads",
				Purpose:     "Web Application",
				Environment: "production",
				IsDefault:   false,
				Resources:   []*fk.StringValueOrRef{newRef("do:droplet:123456")},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})
})
