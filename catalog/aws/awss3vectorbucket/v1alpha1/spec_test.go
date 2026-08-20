package awss3vectorbucketv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsS3VectorBucketSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsS3VectorBucketSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func embeddingsIndex() *AwsS3VectorsIndex {
	return &AwsS3VectorsIndex{
		Name:           "kb-embeddings",
		Dimension:      1024,
		DistanceMetric: "cosine",
	}
}

func minimalVectorBucket() *AwsS3VectorBucketSpec {
	return &AwsS3VectorBucketSpec{Region: "us-east-1"}
}

var _ = ginkgo.Describe("AwsS3VectorBucketSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal empty bucket", func() {
			gomega.Expect(protovalidate.Validate(minimalVectorBucket())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an index sized for Titan embeddings", func() {
			spec := minimalVectorBucket()
			spec.Indexes = []*AwsS3VectorsIndex{embeddingsIndex()}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts KMS encryption with a key", func() {
			spec := minimalVectorBucket()
			spec.Encryption = &AwsS3VectorsEncryption{
				SseType:   "aws:kms",
				KmsKeyArn: literal("arn:aws:kms:us-east-1:111122223333:key/abc"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts non-filterable metadata keys", func() {
			index := embeddingsIndex()
			index.NonFilterableMetadataKeys = []string{"source_text", "source_uri"}
			spec := minimalVectorBucket()
			spec.Indexes = []*AwsS3VectorsIndex{index}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects duplicate index names", func() {
			spec := minimalVectorBucket()
			spec.Indexes = []*AwsS3VectorsIndex{embeddingsIndex(), embeddingsIndex()}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a KMS key under AES256", func() {
			spec := minimalVectorBucket()
			spec.Encryption = &AwsS3VectorsEncryption{
				SseType:   "AES256",
				KmsKeyArn: literal("arn:aws:kms:us-east-1:111122223333:key/abc"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a dimension above 4096", func() {
			index := embeddingsIndex()
			index.Dimension = 5000
			spec := minimalVectorBucket()
			spec.Indexes = []*AwsS3VectorsIndex{index}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a zero dimension", func() {
			index := embeddingsIndex()
			index.Dimension = 0
			spec := minimalVectorBucket()
			spec.Indexes = []*AwsS3VectorsIndex{index}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown distance metric", func() {
			index := embeddingsIndex()
			index.DistanceMetric = "manhattan"
			spec := minimalVectorBucket()
			spec.Indexes = []*AwsS3VectorsIndex{index}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase index name", func() {
			index := embeddingsIndex()
			index.Name = "KB-Embeddings"
			spec := minimalVectorBucket()
			spec.Indexes = []*AwsS3VectorsIndex{index}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects eleven non-filterable keys", func() {
			index := embeddingsIndex()
			for _, key := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"} {
				index.NonFilterableMetadataKeys = append(index.NonFilterableMetadataKeys, key)
			}
			spec := minimalVectorBucket()
			spec.Indexes = []*AwsS3VectorsIndex{index}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
