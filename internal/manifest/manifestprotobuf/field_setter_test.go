package manifestprotobuf

import (
	kubernetesvalkeyv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesvalkey/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestSetProtoField(t *testing.T) {
	tests := []struct {
		name      string
		message   proto.Message
		fieldPath string
		value     interface{}
		expected  proto.Message
		expectErr bool
	}{
		{
			name: "Set existing string field in snake case",
			message: &kubernetesvalkeyv1.KubernetesValkey{Spec: &kubernetesvalkeyv1.KubernetesValkeySpec{
				Config: &kubernetesvalkeyv1.KubernetesValkeyConfig{
					MaxMemory: "256mb",
				},
			}},
			fieldPath: "spec.config.max_memory",
			value:     "512mb",
			expected: &kubernetesvalkeyv1.KubernetesValkey{Spec: &kubernetesvalkeyv1.KubernetesValkeySpec{
				Config: &kubernetesvalkeyv1.KubernetesValkeyConfig{
					MaxMemory: "512mb",
				},
			}},
			expectErr: false,
		},
		{
			name: "Set existing string field in camel case",
			message: &kubernetesvalkeyv1.KubernetesValkey{Spec: &kubernetesvalkeyv1.KubernetesValkeySpec{
				Config: &kubernetesvalkeyv1.KubernetesValkeyConfig{
					MaxMemory: "256mb",
				},
			}},
			fieldPath: "spec.config.maxMemory",
			value:     "512mb",
			expected: &kubernetesvalkeyv1.KubernetesValkey{Spec: &kubernetesvalkeyv1.KubernetesValkeySpec{
				Config: &kubernetesvalkeyv1.KubernetesValkeyConfig{
					MaxMemory: "512mb",
				},
			}},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SetProtoField(tt.message, tt.fieldPath, tt.value)
			if tt.expectErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error: %v", err)
				assert.True(t, proto.Equal(tt.expected, result), "Expected %v but got %v", tt.expected, result)
			}
		})
	}
}
