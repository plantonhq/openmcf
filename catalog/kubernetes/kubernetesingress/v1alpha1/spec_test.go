package kubernetesingressv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestKubernetesIngressSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesIngressSpec Validation Suite")
}

func svcRef(name string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: name},
	}
}

func pathType(t KubernetesIngressHttpPath_KubernetesIngressPathType) *KubernetesIngressHttpPath_KubernetesIngressPathType {
	return &t
}

// validBackend keeps the table-style specs readable.
func validBackend() *KubernetesIngressBackend {
	return &KubernetesIngressBackend{
		ServiceName: svcRef("app-svc"),
		PortNumber:  8080,
	}
}

var _ = ginkgo.Describe("KubernetesIngressSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal spec: one rule, one prefix path", func() {
			spec := &KubernetesIngressSpec{
				Name: "web",
				Rules: []*KubernetesIngressRule{{
					Host:  "app.example.com",
					Paths: []*KubernetesIngressHttpPath{{Path: "/", Backend: validBackend()}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a default backend with no rules", func() {
			spec := &KubernetesIngressSpec{
				Name:           "catch-all",
				DefaultBackend: validBackend(),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a host-less catch-all rule", func() {
			spec := &KubernetesIngressSpec{
				Name: "any-host",
				Rules: []*KubernetesIngressRule{{
					Paths: []*KubernetesIngressHttpPath{{Path: "/", Backend: validBackend()}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a wildcard host", func() {
			spec := &KubernetesIngressSpec{
				Name: "previews",
				Rules: []*KubernetesIngressRule{{
					Host:  "*.preview.example.com",
					Paths: []*KubernetesIngressHttpPath{{Path: "/", Backend: validBackend()}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a named backend port", func() {
			spec := &KubernetesIngressSpec{
				Name: "named-port",
				Rules: []*KubernetesIngressRule{{
					Host: "app.example.com",
					Paths: []*KubernetesIngressHttpPath{{
						Path: "/",
						Backend: &KubernetesIngressBackend{
							ServiceName: svcRef("app-svc"),
							PortName:    "http",
						},
					}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an implementation_specific path with no path string", func() {
			spec := &KubernetesIngressSpec{
				Name: "regex-router",
				Rules: []*KubernetesIngressRule{{
					Host: "app.example.com",
					Paths: []*KubernetesIngressHttpPath{{
						PathType: pathType(KubernetesIngressHttpPath_implementation_specific),
						Backend:  validBackend(),
					}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts TLS with a secret reference and hosts", func() {
			spec := &KubernetesIngressSpec{
				Name: "tls-ingress",
				Tls: []*KubernetesIngressTls{{
					Hosts:      []string{"app.example.com"},
					SecretName: svcRef("app-tls"),
				}},
				Rules: []*KubernetesIngressRule{{
					Host:  "app.example.com",
					Paths: []*KubernetesIngressHttpPath{{Path: "/", Backend: validBackend()}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an SNI-only TLS block without a secret", func() {
			spec := &KubernetesIngressSpec{
				Name: "sni-only",
				Tls: []*KubernetesIngressTls{{
					Hosts: []string{"*.example.com"},
				}},
				DefaultBackend: validBackend(),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a spec with neither rules nor default backend", func() {
			spec := &KubernetesIngressSpec{Name: "routes-nothing"}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an empty name", func() {
			spec := &KubernetesIngressSpec{DefaultBackend: validBackend()}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a rule with no paths", func() {
			spec := &KubernetesIngressSpec{
				Name:  "bad",
				Rules: []*KubernetesIngressRule{{Host: "app.example.com"}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a path not beginning with a slash", func() {
			spec := &KubernetesIngressSpec{
				Name: "bad",
				Rules: []*KubernetesIngressRule{{
					Host:  "app.example.com",
					Paths: []*KubernetesIngressHttpPath{{Path: "api", Backend: validBackend()}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a prefix path with no path string", func() {
			spec := &KubernetesIngressSpec{
				Name: "bad",
				Rules: []*KubernetesIngressRule{{
					Host: "app.example.com",
					Paths: []*KubernetesIngressHttpPath{{
						PathType: pathType(KubernetesIngressHttpPath_prefix),
						Backend:  validBackend(),
					}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a backend with BOTH port number and port name", func() {
			spec := &KubernetesIngressSpec{
				Name: "bad",
				DefaultBackend: &KubernetesIngressBackend{
					ServiceName: svcRef("app-svc"),
					PortNumber:  8080,
					PortName:    "http",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a backend with NEITHER port number nor port name", func() {
			spec := &KubernetesIngressSpec{
				Name: "bad",
				DefaultBackend: &KubernetesIngressBackend{
					ServiceName: svcRef("app-svc"),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a backend without a service name", func() {
			spec := &KubernetesIngressSpec{
				Name: "bad",
				DefaultBackend: &KubernetesIngressBackend{
					PortNumber: 8080,
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a backend port number out of range", func() {
			spec := &KubernetesIngressSpec{
				Name: "bad",
				DefaultBackend: &KubernetesIngressBackend{
					ServiceName: svcRef("app-svc"),
					PortNumber:  70000,
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase backend port name", func() {
			spec := &KubernetesIngressSpec{
				Name: "bad",
				DefaultBackend: &KubernetesIngressBackend{
					ServiceName: svcRef("app-svc"),
					PortName:    "HTTP",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase host", func() {
			spec := &KubernetesIngressSpec{
				Name: "bad",
				Rules: []*KubernetesIngressRule{{
					Host:  "App.Example.Com",
					Paths: []*KubernetesIngressHttpPath{{Path: "/", Backend: validBackend()}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a wildcard that is not the first label", func() {
			spec := &KubernetesIngressSpec{
				Name: "bad",
				Rules: []*KubernetesIngressRule{{
					Host:  "app.*.example.com",
					Paths: []*KubernetesIngressHttpPath{{Path: "/", Backend: validBackend()}},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid TLS host", func() {
			spec := &KubernetesIngressSpec{
				Name: "bad",
				Tls: []*KubernetesIngressTls{{
					Hosts: []string{"-bad.example.com"},
				}},
				DefaultBackend: validBackend(),
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid ingress class name", func() {
			spec := &KubernetesIngressSpec{
				Name:             "bad",
				IngressClassName: "Nginx Class",
				DefaultBackend:   validBackend(),
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})
	})
})
