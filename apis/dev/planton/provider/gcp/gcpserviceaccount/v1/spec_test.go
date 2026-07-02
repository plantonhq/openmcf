package gcpserviceaccountv1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/apis/dev/planton/shared"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "GcpServiceAccountSpec Suite")
}

func boolPtr(b bool) *bool {
	return &b
}

var _ = ginkgo.Describe("GcpServiceAccountSpec", func() {
	var validator protovalidate.Validator

	ginkgo.BeforeEach(func() {
		var err error
		validator, err = protovalidate.New()
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	// Helper to build a minimal valid GcpServiceAccount.
	minimal := func() *GcpServiceAccount {
		return &GcpServiceAccount{
			ApiVersion: "gcp.planton.dev/v1",
			Kind:       "GcpServiceAccount",
			Metadata: &shared.CloudResourceMetadata{
				Name: "test-service-account",
			},
			Spec: &GcpServiceAccountSpec{
				ServiceAccountId: "test-sa-123",
			},
		}
	}

	// ──────────────── Positive Cases ────────────────

	ginkgo.It("should accept a minimal spec with only service_account_id", func() {
		err := validator.Validate(minimal())
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("should accept a fully populated spec", func() {
		msg := minimal()
		msg.Spec.ProjectId = &foreignkeyv1.StringValueOrRef{
			LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "my-gcp-project"},
		}
		msg.Spec.DisplayName = "CI/CD Deployer"
		msg.Spec.Description = "Deploys application releases from the pipeline"
		msg.Spec.Disabled = boolPtr(false)
		msg.Spec.CreateKey = boolPtr(true)
		msg.Spec.ProjectIamRoles = []string{"roles/logging.logWriter", "roles/storage.objectViewer"}
		msg.Spec.OrgId = "123456789012"
		msg.Spec.OrgIamRoles = []string{"roles/resourcemanager.organizationViewer"}
		err := validator.Validate(msg)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("should accept a 6-character service_account_id (lower bound)", func() {
		msg := minimal()
		msg.Spec.ServiceAccountId = "abc123"
		err := validator.Validate(msg)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("should accept a 30-character service_account_id (upper bound)", func() {
		msg := minimal()
		msg.Spec.ServiceAccountId = "a" + strings.Repeat("b", 29)
		err := validator.Validate(msg)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("should accept a description at exactly 256 characters", func() {
		msg := minimal()
		msg.Spec.Description = strings.Repeat("d", 256)
		err := validator.Validate(msg)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("should accept org_iam_roles when org_id is set", func() {
		msg := minimal()
		msg.Spec.OrgId = "123456789012"
		msg.Spec.OrgIamRoles = []string{"roles/resourcemanager.organizationViewer"}
		err := validator.Validate(msg)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	// ──────────────── Negative Cases ────────────────

	ginkgo.It("should reject a missing service_account_id", func() {
		msg := minimal()
		msg.Spec.ServiceAccountId = ""
		err := validator.Validate(msg)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject a service_account_id shorter than 6 characters", func() {
		msg := minimal()
		msg.Spec.ServiceAccountId = "abc12"
		err := validator.Validate(msg)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject a service_account_id longer than 30 characters", func() {
		msg := minimal()
		msg.Spec.ServiceAccountId = "a" + strings.Repeat("b", 30)
		err := validator.Validate(msg)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject a service_account_id starting with a digit", func() {
		msg := minimal()
		msg.Spec.ServiceAccountId = "1invalid-sa"
		err := validator.Validate(msg)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject a service_account_id with uppercase letters", func() {
		msg := minimal()
		msg.Spec.ServiceAccountId = "Invalid-SA-Id"
		err := validator.Validate(msg)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject a service_account_id ending with a hyphen", func() {
		msg := minimal()
		msg.Spec.ServiceAccountId = "test-sa-123-"
		err := validator.Validate(msg)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject a description longer than 256 characters", func() {
		msg := minimal()
		msg.Spec.Description = strings.Repeat("d", 257)
		err := validator.Validate(msg)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("should reject org_iam_roles without org_id", func() {
		msg := minimal()
		msg.Spec.OrgIamRoles = []string{"roles/resourcemanager.organizationViewer"}
		err := validator.Validate(msg)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})
})
