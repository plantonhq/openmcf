package kubernetesconfigmapv1alpha1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestKubernetesConfigMapSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesConfigMapSpec Validation Suite")
}

var _ = ginkgo.Describe("KubernetesConfigMapSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal valid spec with data entries", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "app-config",
				Data: map[string]string{
					"application.properties": "server.port=8080",
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts an empty ConfigMap with no data at all", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "name-reservation",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a name with dots (DNS subdomain)", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "my.dotted.configmap",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a spec without a namespace", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "default-ns-config",
				Data: map[string]string{"key": "value"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace provided as a literal value", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "prod-config",
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "prod"},
				},
				Data: map[string]string{"key": "value"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace provided as a resource reference", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "ref-ns-config",
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
						ValueFrom: &foreignkeyv1.ValueFromRef{Name: "team-namespace"},
					},
				},
				Data: map[string]string{"key": "value"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts data keys with dots, hyphens, and underscores", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "key-charset-config",
				Data: map[string]string{
					"app.properties": "a=1",
					"log-level":      "debug",
					"FEATURE_FLAG":   "on",
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts binary_data with valid base64 values", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "binary-config",
				BinaryData: map[string]string{
					"logo.png":  "iVBORw0KGgo=",
					"cert.der":  "AQIDBA==",
					"empty.bin": "",
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts data and binary_data with disjoint keys", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "mixed-config",
				Data: map[string]string{
					"config.yaml": "replicas: 3",
				},
				BinaryData: map[string]string{
					"keystore.jks": "AQIDBA==",
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a ConfigMap with immutable set to true", func() {
			spec := &KubernetesConfigMapSpec{
				Name:      "app-config-v42",
				Immutable: true,
				Data:      map[string]string{"key": "frozen-value"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a spec with labels and annotations", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "labeled-config",
				Labels: map[string]string{
					"team": "platform",
				},
				Annotations: map[string]string{
					"description": "platform configuration",
				},
				Data: map[string]string{"key": "value"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects an empty name", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a name with uppercase letters", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "MyConfigMap",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a name with underscores", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "app_config",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a name ending with a hyphen", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "app-config-",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a name longer than 253 characters", func() {
			spec := &KubernetesConfigMapSpec{
				Name: strings.Repeat("a", 254),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an empty namespace message", func() {
			spec := &KubernetesConfigMapSpec{
				Name:      "app-config",
				Namespace: &foreignkeyv1.StringValueOrRef{},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a data key with invalid characters", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "bad-key-config",
				Data: map[string]string{
					"invalid/key": "value",
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a binary_data key with invalid characters", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "bad-binary-key-config",
				BinaryData: map[string]string{
					"invalid key": "AQIDBA==",
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a binary_data value that is not valid base64", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "bad-base64-config",
				BinaryData: map[string]string{
					"blob.bin": "not-valid-base64!!!",
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects the same key appearing in both data and binary_data", func() {
			spec := &KubernetesConfigMapSpec{
				Name: "overlap-config",
				Data: map[string]string{
					"shared-key": "text-value",
				},
				BinaryData: map[string]string{
					"shared-key": "AQIDBA==",
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})
	})
})
