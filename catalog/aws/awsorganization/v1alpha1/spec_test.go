package awsorganizationv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsOrganizationSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsOrganizationSpec Validation Suite")
}

// minimalOrganization is the smallest valid instance: an all-features
// organization with no advanced arms.
func minimalOrganization() *AwsOrganizationSpec {
	return &AwsOrganizationSpec{Region: "us-west-2"}
}

func delegationPolicy() *structpb.Struct {
	policy, err := structpb.NewStruct(map[string]any{
		"Version": "2012-10-17",
		"Statement": []any{map[string]any{
			"Effect":    "Allow",
			"Principal": map[string]any{"AWS": "111111111111"},
			"Action":    "organizations:DescribeOrganization",
			"Resource":  "*",
		}},
	})
	if err != nil {
		panic(err)
	}
	return policy
}

var _ = ginkgo.Describe("AwsOrganizationSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal organization", func() {
			gomega.Expect(protovalidate.Validate(minimalOrganization())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an explicit ALL feature set with every advanced arm", func() {
			spec := minimalOrganization()
			spec.FeatureSet = "ALL"
			spec.AwsServiceAccessPrincipals = []string{"cloudtrail.amazonaws.com", "account.amazonaws.com"}
			spec.EnabledPolicyTypes = []string{"SERVICE_CONTROL_POLICY", "TAG_POLICY"}
			spec.DelegatedAdministrators = []*AwsOrganizationDelegatedAdministrator{{
				AccountId:        "111111111111",
				ServicePrincipal: "config.amazonaws.com",
			}}
			spec.ResourcePolicy = delegationPolicy()
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a consolidated-billing organization with no advanced arms", func() {
			spec := minimalOrganization()
			spec.FeatureSet = "CONSOLIDATED_BILLING"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a china-partition service principal", func() {
			spec := minimalOrganization()
			spec.AwsServiceAccessPrincipals = []string{"cloudtrail.amazonaws.com.cn"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts one account delegated for two different services", func() {
			spec := minimalOrganization()
			spec.DelegatedAdministrators = []*AwsOrganizationDelegatedAdministrator{
				{AccountId: "111111111111", ServicePrincipal: "config.amazonaws.com"},
				{AccountId: "111111111111", ServicePrincipal: "guardduty.amazonaws.com"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts root-access management alongside IAM trusted access", func() {
			spec := minimalOrganization()
			spec.AwsServiceAccessPrincipals = []string{"iam.amazonaws.com"}
			spec.RootAccessManagement = &AwsOrganizationRootAccessManagement{
				EnabledFeatures: []string{"RootCredentialsManagement", "RootSessions"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalOrganization()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown feature set", func() {
			spec := minimalOrganization()
			spec.FeatureSet = "BILLING_ONLY"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects service access principals under consolidated billing", func() {
			spec := minimalOrganization()
			spec.FeatureSet = "CONSOLIDATED_BILLING"
			spec.AwsServiceAccessPrincipals = []string{"cloudtrail.amazonaws.com"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects enabled policy types under consolidated billing", func() {
			spec := minimalOrganization()
			spec.FeatureSet = "CONSOLIDATED_BILLING"
			spec.EnabledPolicyTypes = []string{"SERVICE_CONTROL_POLICY"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects delegated administrators under consolidated billing", func() {
			spec := minimalOrganization()
			spec.FeatureSet = "CONSOLIDATED_BILLING"
			spec.DelegatedAdministrators = []*AwsOrganizationDelegatedAdministrator{{
				AccountId:        "111111111111",
				ServicePrincipal: "config.amazonaws.com",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a resource policy under consolidated billing", func() {
			spec := minimalOrganization()
			spec.FeatureSet = "CONSOLIDATED_BILLING"
			spec.ResourcePolicy = delegationPolicy()
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed service principal", func() {
			spec := minimalOrganization()
			spec.AwsServiceAccessPrincipals = []string{"not-a-principal"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate service principals", func() {
			spec := minimalOrganization()
			spec.AwsServiceAccessPrincipals = []string{"cloudtrail.amazonaws.com", "cloudtrail.amazonaws.com"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown policy type", func() {
			spec := minimalOrganization()
			spec.EnabledPolicyTypes = []string{"FIREWALL_POLICY"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a delegated administrator with a malformed account ID", func() {
			spec := minimalOrganization()
			spec.DelegatedAdministrators = []*AwsOrganizationDelegatedAdministrator{{
				AccountId:        "1234",
				ServicePrincipal: "config.amazonaws.com",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate (account, service) delegation pairs", func() {
			spec := minimalOrganization()
			spec.DelegatedAdministrators = []*AwsOrganizationDelegatedAdministrator{
				{AccountId: "111111111111", ServicePrincipal: "config.amazonaws.com"},
				{AccountId: "111111111111", ServicePrincipal: "config.amazonaws.com"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects root-access management without IAM trusted access", func() {
			spec := minimalOrganization()
			spec.RootAccessManagement = &AwsOrganizationRootAccessManagement{
				EnabledFeatures: []string{"RootCredentialsManagement"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects root-access management with no features", func() {
			spec := minimalOrganization()
			spec.AwsServiceAccessPrincipals = []string{"iam.amazonaws.com"}
			spec.RootAccessManagement = &AwsOrganizationRootAccessManagement{}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown root-access feature", func() {
			spec := minimalOrganization()
			spec.AwsServiceAccessPrincipals = []string{"iam.amazonaws.com"}
			spec.RootAccessManagement = &AwsOrganizationRootAccessManagement{
				EnabledFeatures: []string{"RootPasswordRotation"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
