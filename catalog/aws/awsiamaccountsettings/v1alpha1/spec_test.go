package awsiamaccountsettingsv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsIamAccountSettingsSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsIamAccountSettingsSpec Validation Suite")
}

func intPtr(v int32) *int32 { return &v }

func boolPtr(v bool) *bool { return &v }

// minimalSettings is the smallest valid instance: the STS token
// version arm alone.
func minimalSettings() *AwsIamAccountSettingsSpec {
	return &AwsIamAccountSettingsSpec{
		Region: "us-east-1",
		Sts: &AwsIamAccountSettingsSts{
			GlobalEndpointTokenVersion: "v2Token",
		},
	}
}

var _ = ginkgo.Describe("AwsIamAccountSettingsSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal STS-only instance", func() {
			gomega.Expect(protovalidate.Validate(minimalSettings())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an account alias", func() {
			spec := minimalSettings()
			spec.AccountAlias = "acme-platform"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a full password policy", func() {
			spec := &AwsIamAccountSettingsSpec{
				Region: "us-east-1",
				PasswordPolicy: &AwsIamAccountSettingsPasswordPolicy{
					MinimumPasswordLength:      intPtr(14),
					RequireLowercaseCharacters: true,
					RequireNumbers:             true,
					RequireSymbols:             true,
					RequireUppercaseCharacters: true,
					AllowUsersToChangePassword: boolPtr(true),
					MaxPasswordAge:             intPtr(90),
					PasswordReusePrevention:    intPtr(24),
					HardExpiry:                 false,
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalSettings()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an instance managing no arm", func() {
			spec := &AwsIamAccountSettingsSpec{Region: "us-east-1"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an alias with uppercase characters", func() {
			spec := minimalSettings()
			spec.AccountAlias = "Acme-Platform"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an alias with consecutive hyphens", func() {
			spec := minimalSettings()
			spec.AccountAlias = "acme--platform"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an alias ending with a hyphen", func() {
			spec := minimalSettings()
			spec.AccountAlias = "acme-platform-"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an alias under 3 characters", func() {
			spec := minimalSettings()
			spec.AccountAlias = "ac"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a minimum password length under AWS's floor", func() {
			spec := minimalSettings()
			spec.PasswordPolicy = &AwsIamAccountSettingsPasswordPolicy{
				MinimumPasswordLength: intPtr(4),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an out-of-range password age", func() {
			spec := minimalSettings()
			spec.PasswordPolicy = &AwsIamAccountSettingsPasswordPolicy{
				MaxPasswordAge: intPtr(2000),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an out-of-range reuse prevention", func() {
			spec := minimalSettings()
			spec.PasswordPolicy = &AwsIamAccountSettingsPasswordPolicy{
				PasswordReusePrevention: intPtr(30),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown STS token version", func() {
			spec := minimalSettings()
			spec.Sts.GlobalEndpointTokenVersion = "v3Token"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
