package awss3directorybucketv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsS3DirectoryBucketSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsS3DirectoryBucketSpec Validation Suite")
}

func minimalBucket() *AwsS3DirectoryBucketSpec {
	return &AwsS3DirectoryBucketSpec{
		Region: "us-east-1",
		ZoneId: "use1-az4",
	}
}

var _ = ginkgo.Describe("AwsS3DirectoryBucketSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal availability-zone bucket", func() {
			gomega.Expect(protovalidate.Validate(minimalBucket())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an explicit matched redundancy", func() {
			spec := minimalBucket()
			spec.DataRedundancy = "SingleAvailabilityZone"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a Local Zone bucket with matched redundancy", func() {
			spec := minimalBucket()
			spec.ZoneId = "usw2-lax1-az1"
			spec.ZoneType = "LocalZone"
			spec.DataRedundancy = "SingleLocalZone"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts force_destroy", func() {
			spec := minimalBucket()
			spec.ForceDestroy = true
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing zone id", func() {
			spec := minimalBucket()
			spec.ZoneId = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase zone id", func() {
			spec := minimalBucket()
			spec.ZoneId = "USE1-AZ4"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown zone type", func() {
			spec := minimalBucket()
			spec.ZoneType = "WavelengthZone"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects Local Zone redundancy on an availability zone", func() {
			spec := minimalBucket()
			spec.DataRedundancy = "SingleLocalZone"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects availability-zone redundancy on a Local Zone", func() {
			spec := minimalBucket()
			spec.ZoneType = "LocalZone"
			spec.DataRedundancy = "SingleAvailabilityZone"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
