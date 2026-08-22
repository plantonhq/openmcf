package digitaloceandropletautoscalepoolv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	digitaloceanprovider "github.com/plantonhq/planton/catalog/digitalocean"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestDigitalOceanDropletAutoscalePoolSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanDropletAutoscalePoolSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanDropletAutoscalePoolSpec validations", func() {

	newRef := func(value string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: value},
		}
	}

	makeValidTemplate := func() *DigitalOceanDropletAutoscalePoolTemplate {
		return &DigitalOceanDropletAutoscalePoolTemplate{
			Size:    "s-1vcpu-1gb",
			Region:  digitaloceanprovider.DigitalOceanRegion_nyc3,
			Image:   "ubuntu-24-04-x64",
			SshKeys: []*fk.StringValueOrRef{newRef("12345678")},
		}
	}

	makeValidStaticSpec := func() *DigitalOceanDropletAutoscalePoolSpec {
		return &DigitalOceanDropletAutoscalePoolSpec{
			PoolName: "web-pool",
			Scaling: &DigitalOceanDropletAutoscalePoolSpec_Static{
				Static: &DigitalOceanDropletAutoscalePoolStaticScale{
					TargetInstances: 2,
				},
			},
			DropletTemplate: makeValidTemplate(),
		}
	}

	makeValidDynamicSpec := func() *DigitalOceanDropletAutoscalePoolSpec {
		return &DigitalOceanDropletAutoscalePoolSpec{
			PoolName: "web-pool",
			Scaling: &DigitalOceanDropletAutoscalePoolSpec_Dynamic{
				Dynamic: &DigitalOceanDropletAutoscalePoolDynamicScale{
					MinInstances:         1,
					MaxInstances:         5,
					TargetCpuUtilization: proto.Float64(0.7),
				},
			},
			DropletTemplate: makeValidTemplate(),
		}
	}

	ginkgo.Context("Required fields and the scaling oneof", func() {
		ginkgo.It("accepts a valid static pool", func() {
			err := protovalidate.Validate(makeValidStaticSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a valid dynamic pool", func() {
			err := protovalidate.Validate(makeValidDynamicSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing pool_name", func() {
			spec := makeValidStaticSpec()
			spec.PoolName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with no scaling mode", func() {
			spec := makeValidStaticSpec()
			spec.Scaling = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing droplet_template", func() {
			spec := makeValidStaticSpec()
			spec.DropletTemplate = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Static scaling", func() {
		ginkgo.It("rejects a zero target", func() {
			spec := makeValidStaticSpec()
			spec.GetStatic().TargetInstances = 0
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Dynamic scaling", func() {
		ginkgo.It("accepts min equal to max", func() {
			spec := makeValidDynamicSpec()
			spec.GetDynamic().MinInstances = 3
			spec.GetDynamic().MaxInstances = 3
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a memory-only target", func() {
			spec := makeValidDynamicSpec()
			spec.GetDynamic().TargetCpuUtilization = nil
			spec.GetDynamic().TargetMemoryUtilization = proto.Float64(0.6)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts both targets with a cooldown", func() {
			spec := makeValidDynamicSpec()
			spec.GetDynamic().TargetMemoryUtilization = proto.Float64(0.6)
			spec.GetDynamic().CooldownMinutes = proto.Uint32(10)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects a zero min", func() {
			spec := makeValidDynamicSpec()
			spec.GetDynamic().MinInstances = 0
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects min greater than max", func() {
			spec := makeValidDynamicSpec()
			spec.GetDynamic().MinInstances = 6
			spec.GetDynamic().MaxInstances = 5
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a dynamic pool with no utilization target", func() {
			spec := makeValidDynamicSpec()
			spec.GetDynamic().TargetCpuUtilization = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an explicit zero CPU target (unsendable upstream)", func() {
			spec := makeValidDynamicSpec()
			spec.GetDynamic().TargetCpuUtilization = proto.Float64(0)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a CPU target above 1", func() {
			spec := makeValidDynamicSpec()
			spec.GetDynamic().TargetCpuUtilization = proto.Float64(1.5)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a memory target above 1", func() {
			spec := makeValidDynamicSpec()
			spec.GetDynamic().TargetMemoryUtilization = proto.Float64(1.01)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a zero cooldown", func() {
			spec := makeValidDynamicSpec()
			spec.GetDynamic().CooldownMinutes = proto.Uint32(0)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Droplet template", func() {
		ginkgo.It("rejects a template with missing size", func() {
			spec := makeValidStaticSpec()
			spec.DropletTemplate.Size = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a template with missing region", func() {
			spec := makeValidStaticSpec()
			spec.DropletTemplate.Region = digitaloceanprovider.DigitalOceanRegion_digital_ocean_region_unspecified
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a template with missing image", func() {
			spec := makeValidStaticSpec()
			spec.DropletTemplate.Image = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a template with no ssh keys", func() {
			spec := makeValidStaticSpec()
			spec.DropletTemplate.SshKeys = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts ssh keys by reference", func() {
			spec := makeValidStaticSpec()
			spec.DropletTemplate.SshKeys = []*fk.StringValueOrRef{
				{
					LiteralOrRef: &fk.StringValueOrRef_ValueFrom{
						ValueFrom: &fk.ValueFromRef{Name: "ops-key"},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts the full optional surface", func() {
			spec := makeValidStaticSpec()
			spec.DropletTemplate.Vpc = newRef("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
			spec.DropletTemplate.ProjectId = newRef("ffffffff-1111-2222-3333-444444444444")
			spec.DropletTemplate.Tags = []string{"web", "autoscaled"}
			spec.DropletTemplate.WithDropletAgent = true
			spec.DropletTemplate.Ipv6 = true
			spec.DropletTemplate.UserData = "#cloud-config\npackages: [nginx]\n"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects a tag with illegal characters", func() {
			spec := makeValidStaticSpec()
			spec.DropletTemplate.Tags = []string{"bad tag with spaces"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
