package kuberneteshorizontalpodautoscalerv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestKubernetesHorizontalPodAutoscalerSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesHorizontalPodAutoscalerSpec Validation Suite")
}

func literal(v string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: v},
	}
}

func i32(v int32) *int32 { return &v }

func target(name string) *KubernetesHorizontalPodAutoscalerScaleTarget {
	return &KubernetesHorizontalPodAutoscalerScaleTarget{Name: literal(name)}
}

func cpuUtilizationMetric(percent int32) *KubernetesHorizontalPodAutoscalerMetric {
	return &KubernetesHorizontalPodAutoscalerMetric{
		Type: KubernetesHorizontalPodAutoscalerMetricType(1), // resource
		Resource: &KubernetesHorizontalPodAutoscalerResourceMetric{
			Name: "cpu",
			Target: &KubernetesHorizontalPodAutoscalerMetricTarget{
				Type:               KubernetesHorizontalPodAutoscalerMetricTargetType(1), // utilization
				AverageUtilization: i32(percent),
			},
		},
	}
}

var _ = ginkgo.Describe("KubernetesHorizontalPodAutoscalerSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal spec (default 80% CPU metric)", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "checkout-hpa",
				ScaleTarget: target("checkout"),
				MaxReplicas: 10,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace and target provided as references", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name: "checkout-hpa",
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
						ValueFrom: &foreignkeyv1.ValueFromRef{Name: "team-namespace"},
					},
				},
				ScaleTarget: &KubernetesHorizontalPodAutoscalerScaleTarget{
					Name: &foreignkeyv1.StringValueOrRef{
						LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
							ValueFrom: &foreignkeyv1.ValueFromRef{Name: "checkout-deployment"},
						},
					},
				},
				MinReplicas: i32(2),
				MaxReplicas: 20,
				Metrics:     []*KubernetesHorizontalPodAutoscalerMetric{cpuUtilizationMetric(60)},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a container-resource metric", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "app-hpa",
				ScaleTarget: target("app"),
				MaxReplicas: 5,
				Metrics: []*KubernetesHorizontalPodAutoscalerMetric{{
					Type: KubernetesHorizontalPodAutoscalerMetricType(2), // container_resource
					ContainerResource: &KubernetesHorizontalPodAutoscalerContainerResourceMetric{
						Name:      "cpu",
						Container: "app",
						Target: &KubernetesHorizontalPodAutoscalerMetricTarget{
							Type:               KubernetesHorizontalPodAutoscalerMetricTargetType(1),
							AverageUtilization: i32(70),
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a pods metric with average value", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "app-hpa",
				ScaleTarget: target("app"),
				MaxReplicas: 5,
				Metrics: []*KubernetesHorizontalPodAutoscalerMetric{{
					Type: KubernetesHorizontalPodAutoscalerMetricType(3), // pods
					Pods: &KubernetesHorizontalPodAutoscalerPodsMetric{
						Metric: &KubernetesHorizontalPodAutoscalerMetricIdentifier{Name: "requests_per_second"},
						Target: &KubernetesHorizontalPodAutoscalerMetricTarget{
							Type:         KubernetesHorizontalPodAutoscalerMetricTargetType(3), // average_value
							AverageValue: "100",
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an object metric with a raw value target", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "ingress-hpa",
				ScaleTarget: target("app"),
				MaxReplicas: 8,
				Metrics: []*KubernetesHorizontalPodAutoscalerMetric{{
					Type: KubernetesHorizontalPodAutoscalerMetricType(4), // object
					Object: &KubernetesHorizontalPodAutoscalerObjectMetric{
						DescribedObject: &KubernetesHorizontalPodAutoscalerObjectReference{
							ApiVersion: "networking.k8s.io/v1",
							Kind:       "Ingress",
							Name:       "main-ingress",
						},
						Metric: &KubernetesHorizontalPodAutoscalerMetricIdentifier{Name: "requests_per_second"},
						Target: &KubernetesHorizontalPodAutoscalerMetricTarget{
							Type:  KubernetesHorizontalPodAutoscalerMetricTargetType(2), // raw_value
							Value: "10k",
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an external metric with a scoped selector", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "worker-hpa",
				ScaleTarget: target("worker"),
				MaxReplicas: 30,
				Metrics: []*KubernetesHorizontalPodAutoscalerMetric{{
					Type: KubernetesHorizontalPodAutoscalerMetricType(5), // external
					External: &KubernetesHorizontalPodAutoscalerExternalMetric{
						Metric: &KubernetesHorizontalPodAutoscalerMetricIdentifier{
							Name:        "queue_messages_ready",
							MatchLabels: map[string]string{"queue": "orders"},
						},
						Target: &KubernetesHorizontalPodAutoscalerMetricTarget{
							Type:         KubernetesHorizontalPodAutoscalerMetricTargetType(3),
							AverageValue: "30",
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts behavior tuning with policies", func() {
			selectPolicy := KubernetesHorizontalPodAutoscalerScalingRules_min_change
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "tuned-hpa",
				ScaleTarget: target("app"),
				MaxReplicas: 20,
				Metrics:     []*KubernetesHorizontalPodAutoscalerMetric{cpuUtilizationMetric(60)},
				Behavior: &KubernetesHorizontalPodAutoscalerBehavior{
					ScaleUp: &KubernetesHorizontalPodAutoscalerScalingRules{
						StabilizationWindowSeconds: i32(0),
						Policies: []*KubernetesHorizontalPodAutoscalerScalingPolicy{{
							Type:          KubernetesHorizontalPodAutoscalerScalingPolicy_pods,
							Value:         4,
							PeriodSeconds: 15,
						}},
					},
					ScaleDown: &KubernetesHorizontalPodAutoscalerScalingRules{
						StabilizationWindowSeconds: i32(600),
						SelectPolicy:               &selectPolicy,
						Policies: []*KubernetesHorizontalPodAutoscalerScalingPolicy{{
							Type:          KubernetesHorizontalPodAutoscalerScalingPolicy_percent,
							Value:         10,
							PeriodSeconds: 60,
						}},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a disabled scale-down direction", func() {
			disabled := KubernetesHorizontalPodAutoscalerScalingRules_disabled
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "frozen-down",
				ScaleTarget: target("app"),
				MaxReplicas: 10,
				Behavior: &KubernetesHorizontalPodAutoscalerBehavior{
					ScaleDown: &KubernetesHorizontalPodAutoscalerScalingRules{SelectPolicy: &disabled},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a missing scale target", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{Name: "no-target", MaxReplicas: 5}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a scale target without a name", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "no-target-name",
				ScaleTarget: &KubernetesHorizontalPodAutoscalerScaleTarget{},
				MaxReplicas: 5,
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a DaemonSet scale target", func() {
			kind := "DaemonSet"
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name: "ds-hpa",
				ScaleTarget: &KubernetesHorizontalPodAutoscalerScaleTarget{
					Kind: &kind,
					Name: literal("node-agent"),
				},
				MaxReplicas: 5,
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects max_replicas below min_replicas", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "inverted",
				ScaleTarget: target("app"),
				MinReplicas: i32(5),
				MaxReplicas: 2,
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects zero max_replicas", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "no-ceiling",
				ScaleTarget: target("app"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a metric whose source does not match its type", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "mismatched",
				ScaleTarget: target("app"),
				MaxReplicas: 5,
				Metrics: []*KubernetesHorizontalPodAutoscalerMetric{{
					Type: KubernetesHorizontalPodAutoscalerMetricType(3), // pods
					Resource: &KubernetesHorizontalPodAutoscalerResourceMetric{
						Name: "cpu",
						Target: &KubernetesHorizontalPodAutoscalerMetricTarget{
							Type:               KubernetesHorizontalPodAutoscalerMetricTargetType(1),
							AverageUtilization: i32(60),
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a target whose value form does not match its type", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "mismatched-target",
				ScaleTarget: target("app"),
				MaxReplicas: 5,
				Metrics: []*KubernetesHorizontalPodAutoscalerMetric{{
					Type: KubernetesHorizontalPodAutoscalerMetricType(1),
					Resource: &KubernetesHorizontalPodAutoscalerResourceMetric{
						Name: "cpu",
						Target: &KubernetesHorizontalPodAutoscalerMetricTarget{
							Type:  KubernetesHorizontalPodAutoscalerMetricTargetType(1), // utilization
							Value: "60",
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a resource metric on an unsupported resource", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "bad-resource",
				ScaleTarget: target("app"),
				MaxReplicas: 5,
				Metrics: []*KubernetesHorizontalPodAutoscalerMetric{{
					Type: KubernetesHorizontalPodAutoscalerMetricType(1),
					Resource: &KubernetesHorizontalPodAutoscalerResourceMetric{
						Name: "disk",
						Target: &KubernetesHorizontalPodAutoscalerMetricTarget{
							Type:               KubernetesHorizontalPodAutoscalerMetricTargetType(1),
							AverageUtilization: i32(60),
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects utilization above 100 percent", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "over-100",
				ScaleTarget: target("app"),
				MaxReplicas: 5,
				Metrics:     []*KubernetesHorizontalPodAutoscalerMetric{cpuUtilizationMetric(150)},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a stabilization window above one hour", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "long-window",
				ScaleTarget: target("app"),
				MaxReplicas: 5,
				Behavior: &KubernetesHorizontalPodAutoscalerBehavior{
					ScaleDown: &KubernetesHorizontalPodAutoscalerScalingRules{
						StabilizationWindowSeconds: i32(7200),
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a disabled direction that lists policies", func() {
			disabled := KubernetesHorizontalPodAutoscalerScalingRules_disabled
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "contradiction",
				ScaleTarget: target("app"),
				MaxReplicas: 5,
				Behavior: &KubernetesHorizontalPodAutoscalerBehavior{
					ScaleDown: &KubernetesHorizontalPodAutoscalerScalingRules{
						SelectPolicy: &disabled,
						Policies: []*KubernetesHorizontalPodAutoscalerScalingPolicy{{
							Type:          KubernetesHorizontalPodAutoscalerScalingPolicy_pods,
							Value:         1,
							PeriodSeconds: 60,
						}},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a scaling policy with a zero value", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "zero-policy",
				ScaleTarget: target("app"),
				MaxReplicas: 5,
				Behavior: &KubernetesHorizontalPodAutoscalerBehavior{
					ScaleUp: &KubernetesHorizontalPodAutoscalerScalingRules{
						Policies: []*KubernetesHorizontalPodAutoscalerScalingPolicy{{
							Type:          KubernetesHorizontalPodAutoscalerScalingPolicy_percent,
							Value:         0,
							PeriodSeconds: 60,
						}},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a pods metric missing its identifier name", func() {
			spec := &KubernetesHorizontalPodAutoscalerSpec{
				Name:        "no-metric-name",
				ScaleTarget: target("app"),
				MaxReplicas: 5,
				Metrics: []*KubernetesHorizontalPodAutoscalerMetric{{
					Type: KubernetesHorizontalPodAutoscalerMetricType(3),
					Pods: &KubernetesHorizontalPodAutoscalerPodsMetric{
						Metric: &KubernetesHorizontalPodAutoscalerMetricIdentifier{},
						Target: &KubernetesHorizontalPodAutoscalerMetricTarget{
							Type:         KubernetesHorizontalPodAutoscalerMetricTargetType(3),
							AverageValue: "10",
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
