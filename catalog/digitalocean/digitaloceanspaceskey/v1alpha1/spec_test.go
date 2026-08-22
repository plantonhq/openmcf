package digitaloceanspaceskeyv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestDigitalOceanSpacesKeySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanSpacesKeySpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanSpacesKeySpec validations", func() {

	newBucketRef := func(bucketName string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: bucketName},
		}
	}

	makeValidSpec := func() *DigitalOceanSpacesKeySpec {
		return &DigitalOceanSpacesKeySpec{
			KeyName: "ci-uploads",
			Grants: []*DigitalOceanSpacesKeyGrant{
				{
					Bucket:     newBucketRef("app-assets"),
					Permission: "readwrite",
				},
			},
		}
	}

	ginkgo.Context("Required fields", func() {
		ginkgo.It("accepts a valid spec with a bucket grant", func() {
			err := protovalidate.Validate(makeValidSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a spec with no grants (a key with no access)", func() {
			spec := makeValidSpec()
			spec.Grants = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing key_name", func() {
			spec := makeValidSpec()
			spec.KeyName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Grant permission wall", func() {
		ginkgo.It("accepts a read grant on a bucket", func() {
			spec := makeValidSpec()
			spec.Grants[0].Permission = "read"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a bucket grant by reference", func() {
			spec := makeValidSpec()
			spec.Grants[0].Bucket = &fk.StringValueOrRef{
				LiteralOrRef: &fk.StringValueOrRef_ValueFrom{
					ValueFrom: &fk.ValueFromRef{Name: "app-assets-bucket"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects an empty permission (the provider's silent fall-through)", func() {
			spec := makeValidSpec()
			spec.Grants[0].Permission = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a permission outside the wall", func() {
			spec := makeValidSpec()
			spec.Grants[0].Permission = "write"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("The full-access grant grammar", func() {
		ginkgo.It("accepts fullaccess with no bucket", func() {
			spec := makeValidSpec()
			spec.Grants[0] = &DigitalOceanSpacesKeyGrant{Permission: "fullaccess"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects fullaccess naming a bucket", func() {
			spec := makeValidSpec()
			spec.Grants[0].Permission = "fullaccess"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a read grant without a bucket", func() {
			spec := makeValidSpec()
			spec.Grants[0] = &DigitalOceanSpacesKeyGrant{Permission: "read"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a readwrite grant without a bucket", func() {
			spec := makeValidSpec()
			spec.Grants[0] = &DigitalOceanSpacesKeyGrant{Permission: "readwrite"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts mixed per-bucket grants across buckets", func() {
			spec := makeValidSpec()
			spec.Grants = append(spec.Grants, &DigitalOceanSpacesKeyGrant{
				Bucket:     newBucketRef("app-logs"),
				Permission: "read",
			})
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})
})
