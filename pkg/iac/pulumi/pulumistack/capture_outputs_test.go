package pulumistack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectSecretOutputs_DifferingValueIsSecret(t *testing.T) {
	masked := map[string]interface{}{
		"host":     "10.0.0.5",
		"password": "[secret]",
	}
	shown := map[string]interface{}{
		"host":     "10.0.0.5",
		"password": "s3cr3t",
	}

	sensitive := detectSecretOutputs(masked, shown)

	assert.False(t, sensitive["host"])
	assert.True(t, sensitive["password"])
}

func TestDetectSecretOutputs_MissingFromMaskedIsSecret(t *testing.T) {
	masked := map[string]interface{}{"host": "10.0.0.5"}
	shown := map[string]interface{}{"host": "10.0.0.5", "token": "abc"}

	sensitive := detectSecretOutputs(masked, shown)

	assert.True(t, sensitive["token"])
}

func TestDetectSecretOutputs_CompositeValuesCompareDeeply(t *testing.T) {
	masked := map[string]interface{}{
		"endpoints": []interface{}{"a", "b"},
		"conn":      map[string]interface{}{"uri": "[secret]"},
	}
	shown := map[string]interface{}{
		"endpoints": []interface{}{"a", "b"},
		"conn":      map[string]interface{}{"uri": "redis://:pw@host"},
	}

	sensitive := detectSecretOutputs(masked, shown)

	assert.False(t, sensitive["endpoints"])
	assert.True(t, sensitive["conn"])
}
