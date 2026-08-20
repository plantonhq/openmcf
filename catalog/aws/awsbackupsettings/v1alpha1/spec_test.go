package awsbackupsettingsv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsBackupSettingsSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBackupSettingsSpec Validation Suite")
}

// minimalRegionArm is the smallest valid instance: the region arm
// opting one type in.
func minimalRegionArm() *AwsBackupSettingsSpec {
	return &AwsBackupSettingsSpec{
		Region: "us-west-2",
		RegionSettings: &AwsBackupSettingsRegion{
			ResourceTypeOptInPreference: map[string]bool{"EBS": true},
		},
	}
}

var _ = ginkgo.Describe("AwsBackupSettingsSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the region arm alone", func() {
			gomega.Expect(protovalidate.Validate(minimalRegionArm())).To(gomega.BeNil())
		})

		ginkgo.It("accepts the global arm alone", func() {
			spec := &AwsBackupSettingsSpec{
				Region: "us-west-2",
				Global: &AwsBackupSettingsGlobal{
					Settings: map[string]string{"isCrossAccountBackupEnabled": "true"},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts both arms together with management preferences", func() {
			spec := minimalRegionArm()
			spec.Global = &AwsBackupSettingsGlobal{
				Settings: map[string]string{"isCrossAccountBackupEnabled": "false"},
			}
			spec.RegionSettings.ResourceTypeManagementPreference = map[string]bool{"DynamoDB": true}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalRegionArm()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an instance with neither arm", func() {
			spec := &AwsBackupSettingsSpec{Region: "us-west-2"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty global settings map", func() {
			spec := &AwsBackupSettingsSpec{
				Region: "us-west-2",
				Global: &AwsBackupSettingsGlobal{},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty opt-in preference map", func() {
			spec := &AwsBackupSettingsSpec{
				Region:         "us-west-2",
				RegionSettings: &AwsBackupSettingsRegion{},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
