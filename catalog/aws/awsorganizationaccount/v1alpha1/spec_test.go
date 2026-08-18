package awsorganizationaccountv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsOrganizationAccountSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsOrganizationAccountSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalAccount is the smallest valid instance: a member account
// under the organization root with AWS defaults everywhere.
func minimalAccount() *AwsOrganizationAccountSpec {
	return &AwsOrganizationAccountSpec{
		Region:      "us-west-2",
		AccountName: "Workloads Production",
		Email:       "aws-prod@example.com",
	}
}

func billingContact() *AwsOrganizationAccountAlternateContact {
	return &AwsOrganizationAccountAlternateContact{
		Name:         "Jane Doe",
		Title:        "Finance Lead",
		EmailAddress: "billing@example.com",
		PhoneNumber:  "+1 555 0100",
	}
}

var _ = ginkgo.Describe("AwsOrganizationAccountSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal account", func() {
			gomega.Expect(protovalidate.Validate(minimalAccount())).To(gomega.BeNil())
		})

		ginkgo.It("accepts the full surface: placement, role, billing access, contacts, regions", func() {
			spec := minimalAccount()
			spec.ParentId = svr("ou-abcd-12345678")
			spec.RoleName = "OrgBootstrapRole"
			spec.IamUserAccessToBilling = "DENY"
			spec.CloseOnDeletion = true
			spec.AlternateContacts = &AwsOrganizationAccountAlternateContacts{
				Billing:  billingContact(),
				Security: &AwsOrganizationAccountAlternateContact{Name: "Sec Team", Title: "CISO", EmailAddress: "sec@example.com", PhoneNumber: "(555) 0101"},
			}
			spec.PrimaryContact = &AwsOrganizationAccountPrimaryContact{
				FullName:      "Jane Doe",
				AddressLine_1: "1 Main St",
				City:          "Seattle",
				StateOrRegion: "WA",
				PostalCode:    "98101",
				CountryCode:   "US",
				PhoneNumber:   "+1 555 0100",
			}
			spec.Regions = []*AwsOrganizationAccountRegion{
				{RegionName: "ap-southeast-3", Enabled: true},
				{RegionName: "eu-central-2", Enabled: false},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a literal root parent", func() {
			spec := minimalAccount()
			spec.ParentId = svr("r-abcd")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalAccount()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing account name", func() {
			spec := minimalAccount()
			spec.AccountName = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an account name above 50 characters", func() {
			spec := minimalAccount()
			name := ""
			for i := 0; i < 51; i++ {
				name += "a"
			}
			spec.AccountName = name
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed email", func() {
			spec := minimalAccount()
			spec.Email = "not-an-email"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a literal parent that is neither a root nor an OU", func() {
			spec := minimalAccount()
			spec.ParentId = svr("111111111111")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a role name with illegal characters", func() {
			spec := minimalAccount()
			spec.RoleName = "role name with spaces"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown billing-access value", func() {
			spec := minimalAccount()
			spec.IamUserAccessToBilling = "READONLY"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an alternate contact with a malformed phone number", func() {
			spec := minimalAccount()
			contact := billingContact()
			contact.PhoneNumber = "call me"
			spec.AlternateContacts = &AwsOrganizationAccountAlternateContacts{Billing: contact}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a primary contact without the leading + on the phone number", func() {
			spec := minimalAccount()
			spec.PrimaryContact = &AwsOrganizationAccountPrimaryContact{
				FullName:      "Jane Doe",
				AddressLine_1: "1 Main St",
				City:          "Seattle",
				PostalCode:    "98101",
				CountryCode:   "US",
				PhoneNumber:   "555 0100",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a primary contact with a lowercase country code", func() {
			spec := minimalAccount()
			spec.PrimaryContact = &AwsOrganizationAccountPrimaryContact{
				FullName:      "Jane Doe",
				AddressLine_1: "1 Main St",
				City:          "Seattle",
				PostalCode:    "98101",
				CountryCode:   "us",
				PhoneNumber:   "+1 555 0100",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate region entries", func() {
			spec := minimalAccount()
			spec.Regions = []*AwsOrganizationAccountRegion{
				{RegionName: "ap-southeast-3", Enabled: true},
				{RegionName: "ap-southeast-3", Enabled: false},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
