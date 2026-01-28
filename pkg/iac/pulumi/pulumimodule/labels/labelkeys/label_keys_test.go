package labelkeys

import (
	"testing"
)

func TestLabelConversionPrometheusFormat(t *testing.T) {
	testCases := []struct {
		testName                string
		inputLabel              string
		expectedPrometheusLabel string
	}{
		{
			testName:                "openmcf org label should be converted to prometheus label",
			inputLabel:              "openmcf.org/org",
			expectedPrometheusLabel: "openmcf_org_org",
		},
		{
			testName:                "openmcf service label should be converted to prometheus label",
			inputLabel:              "openmcf.org/service",
			expectedPrometheusLabel: "openmcf_org_service",
		},
		{
			testName:                "openmcf service-env label should be converted to prometheus label",
			inputLabel:              "openmcf.org/env",
			expectedPrometheusLabel: "openmcf_org_env",
		},
		{
			testName:                "openmcf kind label should be converted to prometheus label",
			inputLabel:              "openmcf.org/kind",
			expectedPrometheusLabel: "openmcf_org_kind",
		},
		{
			testName:                "openmcf id label should be converted to prometheus label",
			inputLabel:              "openmcf.org/id",
			expectedPrometheusLabel: "openmcf_org_id",
		},
	}
	t.Run("test openmcf label conversion to prometheus format labels", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.testName, func(t *testing.T) {
				r := WithPrometheusFormat(tc.inputLabel)
				if r != tc.expectedPrometheusLabel {
					t.Errorf("expected: %s, got: %s", tc.expectedPrometheusLabel, r)
				}
			})
		}
	})
}
