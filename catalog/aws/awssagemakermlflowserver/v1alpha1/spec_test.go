package awssagemakermlflowserverv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsSagemakerMlflowServerSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSagemakerMlflowServerSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalServer is the smallest valid manifest: region, artifact store,
// and the access role.
func minimalServer() *AwsSagemakerMlflowServerSpec {
	return &AwsSagemakerMlflowServerSpec{
		Region:           "us-west-2",
		ArtifactStoreUri: "s3://my-mlflow/artifacts",
		RoleArn:          svr("arn:aws:iam::123456789012:role/mlflow-server"),
	}
}

var _ = ginkgo.Describe("AwsSagemakerMlflowServerSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalServer())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full surface configured", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalServer()
				spec.Size = "Medium"
				spec.MlflowVersion = "3.0"
				spec.AutomaticModelRegistration = true
				spec.WeeklyMaintenanceWindowStart = "TUE:03:30"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with a non-S3 artifact store", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalServer()
				spec.ArtifactStoreUri = "gs://my-mlflow/artifacts"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a bad size", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalServer()
				spec.Size = "XLarge"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a patch-level mlflow version", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalServer()
				spec.MlflowVersion = "3.0.0"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a bad maintenance window", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalServer()
				spec.WeeklyMaintenanceWindowStart = "TUE:25:00"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a missing role", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalServer()
				spec.RoleArn = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
