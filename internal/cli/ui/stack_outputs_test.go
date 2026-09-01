package ui

import (
	"strings"
	"testing"

	"github.com/plantonhq/planton/pkg/outputs"
	"github.com/stretchr/testify/assert"
)

// TestFormatStackOutputLines_SensitiveValuesNeverRender pins the masking law:
// a sensitive output's VALUE must never appear in rendered output — apply
// runs in CI as much as on laptops, and CI logs are persistent and shared.
func TestFormatStackOutputLines_SensitiveValuesNeverRender(t *testing.T) {
	result := &outputs.CaptureResult{
		Flat: map[string]string{
			"host":         "10.0.0.5",
			"password":     "sup3r-s3cr3t",
			"conn.uri":     "redis://:sup3r-s3cr3t@10.0.0.5",
			"conn.port":    "6379",
			"cluster_name": "cache-prod",
		},
		Sensitive: map[string]bool{
			"password": true,
			"conn":     true,
		},
	}

	lines := FormatStackOutputLines(result)
	rendered := strings.Join(lines, "\n")

	assert.NotContains(t, rendered, "sup3r-s3cr3t")
	assert.Contains(t, rendered, "10.0.0.5")
	assert.Contains(t, rendered, "cache-prod")
	assert.Contains(t, rendered, "(sensitive)")

	// Dotted keys inherit their root output's sensitivity — conn.uri and
	// conn.port both mask because conn is secret.
	for _, line := range lines {
		if strings.Contains(line, "conn.") {
			assert.Contains(t, line, "(sensitive)")
		}
	}
}

func TestFormatStackOutputLines_EmptyAndNilRenderNothing(t *testing.T) {
	assert.Empty(t, FormatStackOutputLines(nil))
	assert.Empty(t, FormatStackOutputLines(&outputs.CaptureResult{}))
}

func TestCaptureResultIsSensitive_RootAndDottedKeys(t *testing.T) {
	result := &outputs.CaptureResult{Sensitive: map[string]bool{"password": true}}

	assert.True(t, result.IsSensitive("password"))
	assert.True(t, result.IsSensitive("password.0"))
	assert.False(t, result.IsSensitive("host"))
	assert.False(t, result.IsSensitive("passwordless"))
}
