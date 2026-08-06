package kubernetesstorageclassv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestKubernetesStorageClassSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesStorageClassSpec Validation Suite")
}

func reclaimPolicy(p KubernetesStorageClassSpec_KubernetesStorageClassReclaimPolicy) *KubernetesStorageClassSpec_KubernetesStorageClassReclaimPolicy {
	return &p
}

func bindingMode(m KubernetesStorageClassSpec_KubernetesStorageClassVolumeBindingMode) *KubernetesStorageClassSpec_KubernetesStorageClassVolumeBindingMode {
	return &m
}

var _ = ginkgo.Describe("KubernetesStorageClassSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal spec (name + provisioner)", func() {
			spec := &KubernetesStorageClassSpec{
				Name:        "fast-ssd",
				Provisioner: "ebs.csi.aws.com",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a full EBS gp3 class", func() {
			spec := &KubernetesStorageClassSpec{
				Name:        "gp3-encrypted",
				Provisioner: "ebs.csi.aws.com",
				Parameters: map[string]string{
					"type":      "gp3",
					"iops":      "6000",
					"encrypted": "true",
				},
				ReclaimPolicy:        reclaimPolicy(KubernetesStorageClassSpec_retain),
				VolumeBindingMode:    bindingMode(KubernetesStorageClassSpec_wait_for_first_consumer),
				AllowVolumeExpansion: true,
				MountOptions:         []string{"noatime"},
				IsDefaultClass:       true,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts allowed topologies with zone requirements", func() {
			spec := &KubernetesStorageClassSpec{
				Name:              "zonal-ssd",
				Provisioner:       "pd.csi.storage.gke.io",
				VolumeBindingMode: bindingMode(KubernetesStorageClassSpec_wait_for_first_consumer),
				AllowedTopologies: []*KubernetesStorageClassTopologySelectorTerm{{
					MatchLabelExpressions: []*KubernetesStorageClassTopologySelectorLabelRequirement{{
						Key:    "topology.kubernetes.io/zone",
						Values: []string{"us-central1-a", "us-central1-b"},
					}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a dotted DNS subdomain name", func() {
			spec := &KubernetesStorageClassSpec{
				Name:        "tier1.fast.storage",
				Provisioner: "rancher.io/local-path",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a missing name", func() {
			spec := &KubernetesStorageClassSpec{Provisioner: "ebs.csi.aws.com"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase name", func() {
			spec := &KubernetesStorageClassSpec{Name: "FastSSD", Provisioner: "ebs.csi.aws.com"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a name with a trailing hyphen", func() {
			spec := &KubernetesStorageClassSpec{Name: "fast-", Provisioner: "ebs.csi.aws.com"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing provisioner", func() {
			spec := &KubernetesStorageClassSpec{Name: "fast-ssd"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an undefined reclaim policy enum value", func() {
			bad := KubernetesStorageClassSpec_KubernetesStorageClassReclaimPolicy(99)
			spec := &KubernetesStorageClassSpec{
				Name:          "fast-ssd",
				Provisioner:   "ebs.csi.aws.com",
				ReclaimPolicy: &bad,
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a topology term without requirements", func() {
			spec := &KubernetesStorageClassSpec{
				Name:              "zonal",
				Provisioner:       "ebs.csi.aws.com",
				AllowedTopologies: []*KubernetesStorageClassTopologySelectorTerm{{}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a topology requirement without values", func() {
			spec := &KubernetesStorageClassSpec{
				Name:        "zonal",
				Provisioner: "ebs.csi.aws.com",
				AllowedTopologies: []*KubernetesStorageClassTopologySelectorTerm{{
					MatchLabelExpressions: []*KubernetesStorageClassTopologySelectorLabelRequirement{{
						Key: "topology.kubernetes.io/zone",
					}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a topology requirement without a key", func() {
			spec := &KubernetesStorageClassSpec{
				Name:        "zonal",
				Provisioner: "ebs.csi.aws.com",
				AllowedTopologies: []*KubernetesStorageClassTopologySelectorTerm{{
					MatchLabelExpressions: []*KubernetesStorageClassTopologySelectorLabelRequirement{{
						Values: []string{"us-east-1a"},
					}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
