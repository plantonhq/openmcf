package awsbackupframeworkv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsBackupFrameworkSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBackupFrameworkSpec Validation Suite")
}

// minimalFramework is the smallest valid instance: one parameterless
// control.
func minimalFramework() *AwsBackupFrameworkSpec {
	return &AwsBackupFrameworkSpec{
		Region:        "us-west-2",
		FrameworkName: "backup_posture",
		Controls: []*AwsBackupFrameworkControl{{
			Name: "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN",
		}},
	}
}

var _ = ginkgo.Describe("AwsBackupFrameworkSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal framework", func() {
			gomega.Expect(protovalidate.Validate(minimalFramework())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a parameterized control with a scope", func() {
			spec := minimalFramework()
			spec.Controls = append(spec.Controls, &AwsBackupFrameworkControl{
				Name: "BACKUP_RECOVERY_POINT_MINIMUM_RETENTION_CHECK",
				InputParameters: []*AwsBackupFrameworkControlInputParameter{{
					Name:  "requiredRetentionDays",
					Value: "35",
				}},
				Scope: &AwsBackupFrameworkControlScope{
					ComplianceResourceTypes: []string{"EBS"},
					Tags:                    map[string]string{"environment": "prod"},
				},
			})
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalFramework()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a hyphenated framework name", func() {
			spec := minimalFramework()
			spec.FrameworkName = "backup-posture"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a framework name starting with a digit", func() {
			spec := minimalFramework()
			spec.FrameworkName = "1backup"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a framework without controls", func() {
			spec := minimalFramework()
			spec.Controls = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate control names", func() {
			spec := minimalFramework()
			spec.Controls = append(spec.Controls, &AwsBackupFrameworkControl{
				Name: "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN",
			})
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a scope tag map with more than one pair", func() {
			spec := minimalFramework()
			spec.Controls[0].Scope = &AwsBackupFrameworkControlScope{
				Tags: map[string]string{"a": "1", "b": "2"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects more than 100 compliance resource IDs", func() {
			ids := make([]string, 101)
			for i := range ids {
				ids[i] = "vol-0abc"
			}
			spec := minimalFramework()
			spec.Controls[0].Scope = &AwsBackupFrameworkControlScope{ComplianceResourceIds: ids}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
