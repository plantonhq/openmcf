package awssagemakermlflowappv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsSagemakerMlflowAppSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSagemakerMlflowAppSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalApp is the smallest valid manifest: region, artifact store,
// and the access role.
func minimalApp() *AwsSagemakerMlflowAppSpec {
	return &AwsSagemakerMlflowAppSpec{
		Region:           "us-west-2",
		ArtifactStoreUri: "s3://my-mlflow/artifacts",
		RoleArn:          svr("arn:aws:iam::123456789012:role/mlflow-app"),
	}
}

var _ = ginkgo.Describe("AwsSagemakerMlflowAppSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalApp())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full surface configured", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalApp()
				spec.AccountDefaultStatus = "ENABLED"
				spec.DefaultDomainIds = []*foreignkeyv1.StringValueOrRef{svr("d-abc123")}
				spec.ModelRegistrationMode = "AutoModelRegistrationEnabled"
				spec.WeeklyMaintenanceWindowStart = "SUN:03:00"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with an empty artifact store", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalApp()
				spec.ArtifactStoreUri = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a bad account default status", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalApp()
				spec.AccountDefaultStatus = "ON"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a bad registration mode", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalApp()
				spec.ModelRegistrationMode = "Enabled"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a bad maintenance window day", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalApp()
				spec.WeeklyMaintenanceWindowStart = "XYZ:03:00"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a missing role", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalApp()
				spec.RoleArn = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
