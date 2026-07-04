package profile

import "testing"

func TestToPascalCase_RegistryLookup(t *testing.T) {
	tests := []struct {
		component string
		want      string
	}{
		{component: "awslambda", want: "AwsLambda"},
		{component: "awskmskey", want: "AwsKmsKey"},
		{component: "awslambdaeventsourcemapping", want: "AwsLambdaEventSourceMapping"},
		{component: "awssecuritygroup", want: "AwsSecurityGroup"},
	}
	for _, tc := range tests {
		t.Run(tc.component, func(t *testing.T) {
			if got := toPascalCase(tc.component); got != tc.want {
				t.Errorf("toPascalCase(%q) = %q, want %q", tc.component, got, tc.want)
			}
		})
	}
}
