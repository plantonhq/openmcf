package digitaloceanmonitoralertv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestDigitalOceanMonitorAlertSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanMonitorAlertSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanMonitorAlertSpec validations", func() {

	newRef := func(value string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: value},
		}
	}

	makeValidSpec := func() *DigitalOceanMonitorAlertSpec {
		return &DigitalOceanMonitorAlertSpec{
			Description: "CPU is running hot",
			MetricType:  "v1/insights/droplet/cpu",
			Compare:     "GreaterThan",
			Value:       90,
			Window:      "5m",
			DropletIds:  []*fk.StringValueOrRef{newRef("123456789")},
			Alerts: &DigitalOceanMonitorAlertNotifications{
				Emails: []string{"ops@example.com"},
			},
		}
	}

	ginkgo.Context("Required fields", func() {
		ginkgo.It("accepts a minimal valid spec", func() {
			err := protovalidate.Validate(makeValidSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing description", func() {
			spec := makeValidSpec()
			spec.Description = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing metric_type", func() {
			spec := makeValidSpec()
			spec.MetricType = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing alerts", func() {
			spec := makeValidSpec()
			spec.Alerts = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing window", func() {
			spec := makeValidSpec()
			spec.Window = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Metric type", func() {
		ginkgo.It("accepts a value from each metric family", func() {
			cases := []struct {
				metric string
				adjust func(*DigitalOceanMonitorAlertSpec)
			}{
				{"v1/insights/droplet/memory_utilization_percent", func(s *DigitalOceanMonitorAlertSpec) {}},
				{"v1/insights/lbaas/droplet_health", func(s *DigitalOceanMonitorAlertSpec) {
					s.DropletIds = nil
					s.LoadBalancerIds = []*fk.StringValueOrRef{newRef("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")}
				}},
				{"v1/dbaas/alerts/cpu_alerts", func(s *DigitalOceanMonitorAlertSpec) {
					s.DropletIds = nil
					s.DatabaseClusterIds = []*fk.StringValueOrRef{newRef("bbbbbbbb-cccc-dddd-eeee-ffffffffffff")}
				}},
			}
			for _, c := range cases {
				spec := makeValidSpec()
				spec.MetricType = c.metric
				c.adjust(spec)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil(), "metric %q should validate", c.metric)
			}
		})

		ginkgo.It("rejects an unknown metric", func() {
			spec := makeValidSpec()
			spec.MetricType = "v1/insights/droplet/cpu_utilization_percent"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Compare and window", func() {
		ginkgo.It("rejects the snake_case comparison spelling (that is the uptime API's)", func() {
			spec := makeValidSpec()
			spec.Compare = "greater_than"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown window", func() {
			spec := makeValidSpec()
			spec.Window = "15m"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a negative threshold value", func() {
			spec := makeValidSpec()
			spec.Value = -1
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Metric-family target pairing", func() {
		ginkgo.It("rejects a droplet metric targeting load balancers", func() {
			spec := makeValidSpec()
			spec.LoadBalancerIds = []*fk.StringValueOrRef{newRef("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a load-balancer metric targeting droplets", func() {
			spec := makeValidSpec()
			spec.MetricType = "v1/insights/lbaas/droplet_health"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a load-balancer metric targeting tags", func() {
			spec := makeValidSpec()
			spec.MetricType = "v1/insights/lbaas/droplet_health"
			spec.DropletIds = nil
			spec.Tags = []string{"web"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a database metric targeting droplets", func() {
			spec := makeValidSpec()
			spec.MetricType = "v1/dbaas/alerts/cpu_alerts"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts a droplet metric targeting tags only", func() {
			spec := makeValidSpec()
			spec.DropletIds = nil
			spec.Tags = []string{"web", "prod:frontend"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a droplet metric with no targets at all (tag or entity set at deploy time is the API's call)", func() {
			spec := makeValidSpec()
			spec.DropletIds = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid tag", func() {
			spec := makeValidSpec()
			spec.Tags = []string{"has space"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Notification channels", func() {
		ginkgo.It("rejects alerts with no channel at all", func() {
			spec := makeValidSpec()
			spec.Alerts = &DigitalOceanMonitorAlertNotifications{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts slack as the only channel", func() {
			spec := makeValidSpec()
			spec.Alerts = &DigitalOceanMonitorAlertNotifications{
				Slack: []*DigitalOceanMonitorAlertSlack{{
					Channel: "#alerts",
					Url:     "https://hooks.slack.com/services/T0/B0/XXXX",
				}},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects a slack target without a url", func() {
			spec := makeValidSpec()
			spec.Alerts = &DigitalOceanMonitorAlertNotifications{
				Slack: []*DigitalOceanMonitorAlertSlack{{Channel: "#alerts"}},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a slack target without a channel", func() {
			spec := makeValidSpec()
			spec.Alerts = &DigitalOceanMonitorAlertNotifications{
				Slack: []*DigitalOceanMonitorAlertSlack{{Url: "https://hooks.slack.com/services/T0/B0/XXXX"}},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
