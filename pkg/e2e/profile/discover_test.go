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
		// Kubernetes kinds resolve through the registry too; Signoz's enum
		// name matches its TestKubernetesSignoz_* entrypoints (the retired
		// hand-maintained table had it wrong as "KubernetesSigNoz").
		{component: "kubernetessignoz", want: "KubernetesSignoz"},
		{component: "kubernetesgatewayapicrds", want: "KubernetesGatewayApiCrds"},
		// ArgoCD's test entrypoints deviate from the enum name
		// (KubernetesArgocd); the explicit override keeps its green CI lane
		// matching TestKubernetesArgoCD_*.
		{component: "kubernetesargocd", want: "KubernetesArgoCD"},
	}
	for _, tc := range tests {
		t.Run(tc.component, func(t *testing.T) {
			if got := toPascalCase(tc.component); got != tc.want {
				t.Errorf("toPascalCase(%q) = %q, want %q", tc.component, got, tc.want)
			}
		})
	}
}
