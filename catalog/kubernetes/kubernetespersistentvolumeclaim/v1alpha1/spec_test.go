package kubernetespersistentvolumeclaimv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestKubernetesPersistentVolumeClaimSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesPersistentVolumeClaimSpec Validation Suite")
}

func literal(v string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: v},
	}
}

var _ = ginkgo.Describe("KubernetesPersistentVolumeClaimSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal spec (name + storage request)", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:           "data",
				StorageRequest: "10Gi",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace provided as a resource reference", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:           "data",
				StorageRequest: "10Gi",
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
						ValueFrom: &foreignkeyv1.ValueFromRef{Name: "team-namespace"},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a full dynamic-provisioning claim", func() {
			mode := KubernetesPersistentVolumeClaimSpec_block
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:             "shared-cache",
				Namespace:        literal("apps"),
				AccessModes:      []string{"ReadWriteMany"},
				StorageRequest:   "100Gi",
				StorageLimit:     "200Gi",
				StorageClassName: literal("fast-ssd"),
				VolumeMode:       &mode,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a static-binding claim (pre-provisioned volume)", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:                       "restored-data",
				StorageRequest:             "50Gi",
				DisableDynamicProvisioning: true,
				VolumeName:                 "pv-restored-0001",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a selector-narrowed static claim", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:                       "archive",
				StorageRequest:             "1Ti",
				DisableDynamicProvisioning: true,
				Selector: &KubernetesPersistentVolumeClaimLabelSelector{
					MatchLabels: map[string]string{"tier": "archive"},
					MatchExpressions: []*KubernetesPersistentVolumeClaimLabelSelectorRequirement{{
						Key:      "zone",
						Operator: "In",
						Values:   []string{"us-east-1a"},
					}},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a clone data source", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:           "clone-of-data",
				StorageRequest: "10Gi",
				DataSource: &KubernetesPersistentVolumeClaimDataSource{
					Kind: KubernetesPersistentVolumeClaimDataSource_persistent_volume_claim,
					Name: "data",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a snapshot restore data source", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:           "restored",
				StorageRequest: "10Gi",
				DataSource: &KubernetesPersistentVolumeClaimDataSource{
					Kind: KubernetesPersistentVolumeClaimDataSource_volume_snapshot,
					Name: "nightly-backup",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a missing name", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{StorageRequest: "10Gi"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing storage request", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{Name: "data"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed storage request", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{Name: "data", StorageRequest: "ten-gigs"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed storage limit", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{Name: "data", StorageRequest: "10Gi", StorageLimit: "big"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid access mode", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:           "data",
				StorageRequest: "10Gi",
				AccessModes:    []string{"ReadWriteEverywhere"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects combining disable_dynamic_provisioning with a storage class", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:                       "data",
				StorageRequest:             "10Gi",
				DisableDynamicProvisioning: true,
				StorageClassName:           literal("fast-ssd"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a data source without a name", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:           "clone",
				StorageRequest: "10Gi",
				DataSource: &KubernetesPersistentVolumeClaimDataSource{
					Kind: KubernetesPersistentVolumeClaimDataSource_persistent_volume_claim,
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a data source without a kind", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:           "clone",
				StorageRequest: "10Gi",
				DataSource:     &KubernetesPersistentVolumeClaimDataSource{Name: "data"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects selector In requirements without values", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:           "data",
				StorageRequest: "10Gi",
				Selector: &KubernetesPersistentVolumeClaimLabelSelector{
					MatchExpressions: []*KubernetesPersistentVolumeClaimLabelSelectorRequirement{{
						Key:      "tier",
						Operator: "In",
					}},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects selector Exists requirements carrying values", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{
				Name:           "data",
				StorageRequest: "10Gi",
				Selector: &KubernetesPersistentVolumeClaimLabelSelector{
					MatchExpressions: []*KubernetesPersistentVolumeClaimLabelSelectorRequirement{{
						Key:      "tier",
						Operator: "Exists",
						Values:   []string{"archive"},
					}},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase name", func() {
			spec := &KubernetesPersistentVolumeClaimSpec{Name: "Data", StorageRequest: "10Gi"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
