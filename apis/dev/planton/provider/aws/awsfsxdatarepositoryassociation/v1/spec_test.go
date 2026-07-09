package awsfsxdatarepositoryassociationv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsFsxDataRepositoryAssociationSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsFsxDataRepositoryAssociationSpec Validation Suite")
}

func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}

var _ = ginkgo.Describe("AwsFsxDataRepositoryAssociationSpec validations", func() {
	var spec *AwsFsxDataRepositoryAssociationSpec

	ginkgo.BeforeEach(func() {
		spec = &AwsFsxDataRepositoryAssociationSpec{
			Region:             "us-west-2",
			FileSystemId:       strRef("fs-0123456789abcdef0"),
			FileSystemPath:     "/datasets",
			DataRepositoryPath: "s3://training-data/2026/",
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal valid spec", func() {
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts full bidirectional sync policies", func() {
		spec.AutoImportEvents = []string{"NEW", "CHANGED", "DELETED"}
		spec.AutoExportEvents = []string{"NEW", "CHANGED"}
		spec.ImportedFileChunkSize = int32Ptr(2048)
		spec.BatchImportMetaDataOnCreate = true
		spec.DeleteDataInFilesystem = false
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts the root file-system path", func() {
		spec.FileSystemPath = "/"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Field-level validation failures
	// -------------------------------------------------------------------------

	ginkgo.Context("field-level validations", func() {
		ginkgo.It("fails when file_system_id is missing", func() {
			spec.FileSystemId = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when file_system_path does not begin with '/'", func() {
			spec.FileSystemPath = "datasets"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when data_repository_path is not an S3 URI", func() {
			spec.DataRepositoryPath = "https://training-data"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when imported_file_chunk_size is out of range", func() {
			spec.ImportedFileChunkSize = int32Ptr(512001)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: event sets
	// -------------------------------------------------------------------------

	ginkgo.Context("sync event sets", func() {
		ginkgo.It("fails on an unknown import event", func() {
			spec.AutoImportEvents = []string{"CREATED"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on an unknown export event", func() {
			spec.AutoExportEvents = []string{"ALL"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on duplicate events", func() {
			spec.AutoImportEvents = []string{"NEW", "NEW"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
