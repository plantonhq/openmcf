package kubernetesdaemonsetv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	kubernetes "github.com/plantonhq/planton/catalog/kubernetes"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func validContainer() *kubernetes.WorkloadContainer {
	return &kubernetes.WorkloadContainer{
		Image: &kubernetes.WorkloadContainerImage{
			Repo: "ghcr.io/acme/node-agent",
			Tag:  "v2.1.0",
		},
	}
}

func validSpec() *KubernetesDaemonSetSpec {
	return &KubernetesDaemonSetSpec{
		Namespace: literal("kube-system"),
		Container: &KubernetesDaemonSetContainer{
			App: validContainer(),
		},
	}
}

func TestKubernetesDaemonSetSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesDaemonSetSpec Validation Suite")
}

var _ = ginkgo.Describe("KubernetesDaemonSetSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal spec with namespace and app container", func() {
			err := protovalidate.Validate(validSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a full-featured node agent with host access, sidecars, and scheduling", func() {
			app := validContainer()
			app.Name = "agent"
			app.ImagePullPolicy = "IfNotPresent"
			app.Ports = []*kubernetes.WorkloadContainerPort{
				{Name: "metrics", ContainerPort: 9100, HostPort: 9100},
			}
			app.SecurityContext = &kubernetes.WorkloadContainerSecurityContext{
				Privileged: true,
			}
			spec := &KubernetesDaemonSetSpec{
				Namespace:       literal("monitoring"),
				CreateNamespace: true,
				Container: &KubernetesDaemonSetContainer{
					App: app,
					Sidecars: []*kubernetes.WorkloadContainer{
						{
							Name: "config-reloader",
							Image: &kubernetes.WorkloadContainerImage{
								Repo: "ghcr.io/acme/reloader",
								Tag:  "v0.5.0",
							},
						},
					},
				},
				Pod: &kubernetes.WorkloadPod{
					ServiceAccount: literal("node-agent-identity"),
					HostNetwork:    true,
					HostPid:        true,
					DnsPolicy:      "ClusterFirstWithHostNet",
					Scheduling: &kubernetes.WorkloadScheduling{
						Tolerations: []*kubernetes.WorkloadToleration{
							{Key: "node-role.kubernetes.io/control-plane", Operator: "Exists", Effect: "NoSchedule"},
						},
						NodeAffinity: &kubernetes.WorkloadNodeAffinity{
							Required: []*kubernetes.WorkloadNodeSelectorTerm{
								{
									MatchExpressions: []*kubernetes.WorkloadNodeSelectorRequirement{
										{Key: "kubernetes.io/os", Operator: "In", Values: []string{"linux"}},
									},
								},
							},
						},
					},
					PriorityClassName: "system-node-critical",
				},
				UpdateStrategy: &KubernetesDaemonSetUpdateStrategy{
					Type:           "RollingUpdate",
					MaxUnavailable: "10%",
				},
				MinReadySeconds:      proto.Int32(15),
				RevisionHistoryLimit: proto.Int32(5),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts the OnDelete update strategy", func() {
			spec := validSpec()
			spec.UpdateStrategy = &KubernetesDaemonSetUpdateStrategy{Type: "OnDelete"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a surging rolling update with max_unavailable 0", func() {
			spec := validSpec()
			spec.UpdateStrategy = &KubernetesDaemonSetUpdateStrategy{
				Type:           "RollingUpdate",
				MaxUnavailable: "0",
				MaxSurge:       "1",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a node affinity requirement using Exists with no values", func() {
			spec := validSpec()
			spec.Pod = &kubernetes.WorkloadPod{
				Scheduling: &kubernetes.WorkloadScheduling{
					NodeAffinity: &kubernetes.WorkloadNodeAffinity{
						Required: []*kubernetes.WorkloadNodeSelectorTerm{
							{
								MatchExpressions: []*kubernetes.WorkloadNodeSelectorRequirement{
									{Key: "node.kubernetes.io/gpu", Operator: "Exists"},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a Gt node affinity requirement with exactly one value", func() {
			spec := validSpec()
			spec.Pod = &kubernetes.WorkloadPod{
				Scheduling: &kubernetes.WorkloadScheduling{
					NodeAffinity: &kubernetes.WorkloadNodeAffinity{
						Preferred: []*kubernetes.WorkloadPreferredNodeSelectorTerm{
							{
								Weight: 50,
								Term: &kubernetes.WorkloadNodeSelectorTerm{
									MatchExpressions: []*kubernetes.WorkloadNodeSelectorRequirement{
										{Key: "acme.io/cpu-count", Operator: "Gt", Values: []string{"8"}},
									},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a host alias with a valid IP and hostnames", func() {
			spec := validSpec()
			spec.Pod = &kubernetes.WorkloadPod{
				HostAliases: []*kubernetes.WorkloadHostAlias{
					{Ip: "10.0.0.5", Hostnames: []string{"registry.internal"}},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a spec without a namespace", func() {
			spec := validSpec()
			spec.Namespace = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a spec without a container", func() {
			spec := validSpec()
			spec.Container = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a container group without an app container", func() {
			spec := validSpec()
			spec.Container = &KubernetesDaemonSetContainer{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an app image without a repo", func() {
			spec := validSpec()
			spec.Container.App.Image = &kubernetes.WorkloadContainerImage{Tag: "v2.1.0"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an unnamed sidecar", func() {
			spec := validSpec()
			spec.Container.Sidecars = []*kubernetes.WorkloadContainer{
				{
					Image: &kubernetes.WorkloadContainerImage{Repo: "busybox", Tag: "1.36"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown update strategy type", func() {
			spec := validSpec()
			spec.UpdateStrategy = &KubernetesDaemonSetUpdateStrategy{Type: "Recreate"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an update strategy with both max_unavailable and max_surge at zero", func() {
			spec := validSpec()
			spec.UpdateStrategy = &KubernetesDaemonSetUpdateStrategy{
				Type:           "RollingUpdate",
				MaxUnavailable: "0%",
				MaxSurge:       "0",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a negative min_ready_seconds", func() {
			spec := validSpec()
			spec.MinReadySeconds = proto.Int32(-1)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a negative revision_history_limit", func() {
			spec := validSpec()
			spec.RevisionHistoryLimit = proto.Int32(-1)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown image pull policy", func() {
			spec := validSpec()
			spec.Container.App.ImagePullPolicy = "Sometimes"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a host port outside 1-65535", func() {
			spec := validSpec()
			spec.Container.App.Ports = []*kubernetes.WorkloadContainerPort{
				{Name: "metrics", ContainerPort: 9100, HostPort: 70000},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a port without a name", func() {
			spec := validSpec()
			spec.Container.App.Ports = []*kubernetes.WorkloadContainerPort{
				{ContainerPort: 9100},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an In node affinity requirement with no values", func() {
			spec := validSpec()
			spec.Pod = &kubernetes.WorkloadPod{
				Scheduling: &kubernetes.WorkloadScheduling{
					NodeAffinity: &kubernetes.WorkloadNodeAffinity{
						Required: []*kubernetes.WorkloadNodeSelectorTerm{
							{
								MatchExpressions: []*kubernetes.WorkloadNodeSelectorRequirement{
									{Key: "kubernetes.io/os", Operator: "In"},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an Exists node affinity requirement carrying values", func() {
			spec := validSpec()
			spec.Pod = &kubernetes.WorkloadPod{
				Scheduling: &kubernetes.WorkloadScheduling{
					NodeAffinity: &kubernetes.WorkloadNodeAffinity{
						Required: []*kubernetes.WorkloadNodeSelectorTerm{
							{
								MatchExpressions: []*kubernetes.WorkloadNodeSelectorRequirement{
									{Key: "node.kubernetes.io/gpu", Operator: "Exists", Values: []string{"true"}},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a Gt node affinity requirement with two values", func() {
			spec := validSpec()
			spec.Pod = &kubernetes.WorkloadPod{
				Scheduling: &kubernetes.WorkloadScheduling{
					NodeAffinity: &kubernetes.WorkloadNodeAffinity{
						Required: []*kubernetes.WorkloadNodeSelectorTerm{
							{
								MatchExpressions: []*kubernetes.WorkloadNodeSelectorRequirement{
									{Key: "acme.io/cpu-count", Operator: "Gt", Values: []string{"4", "8"}},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown node affinity operator", func() {
			spec := validSpec()
			spec.Pod = &kubernetes.WorkloadPod{
				Scheduling: &kubernetes.WorkloadScheduling{
					NodeAffinity: &kubernetes.WorkloadNodeAffinity{
						Required: []*kubernetes.WorkloadNodeSelectorTerm{
							{
								MatchExpressions: []*kubernetes.WorkloadNodeSelectorRequirement{
									{Key: "kubernetes.io/os", Operator: "Equals", Values: []string{"linux"}},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a preferred node affinity term with weight above 100", func() {
			spec := validSpec()
			spec.Pod = &kubernetes.WorkloadPod{
				Scheduling: &kubernetes.WorkloadScheduling{
					NodeAffinity: &kubernetes.WorkloadNodeAffinity{
						Preferred: []*kubernetes.WorkloadPreferredNodeSelectorTerm{
							{
								Weight: 101,
								Term: &kubernetes.WorkloadNodeSelectorTerm{
									MatchExpressions: []*kubernetes.WorkloadNodeSelectorRequirement{
										{Key: "kubernetes.io/os", Operator: "Exists"},
									},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a host alias with an invalid IP", func() {
			spec := validSpec()
			spec.Pod = &kubernetes.WorkloadPod{
				HostAliases: []*kubernetes.WorkloadHostAlias{
					{Ip: "not-an-ip", Hostnames: []string{"registry.internal"}},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a host alias without hostnames", func() {
			spec := validSpec()
			spec.Pod = &kubernetes.WorkloadPod{
				HostAliases: []*kubernetes.WorkloadHostAlias{
					{Ip: "10.0.0.5"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown DNS policy", func() {
			spec := validSpec()
			spec.Pod = &kubernetes.WorkloadPod{DnsPolicy: "NodeFirst"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})
	})
})
