package workloadpod

import (
	"strings"
	"testing"

	kubernetesv1 "github.com/plantonhq/planton/catalog/kubernetes"
)

func TestBuildImagePullSecretData_NilWhenNothingIsDeclared(t *testing.T) {
	for name, pod := range map[string]*kubernetesv1.WorkloadPod{
		"nil pod":       nil,
		"no registries": {},
	} {
		data, err := BuildImagePullSecretData(pod)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", name, err)
		}
		if data != nil {
			t.Fatalf("%s: a pod that declares no registry gets no Secret data, got %v", name, data)
		}
	}
}

func TestBuildImagePullSecretData_OneDockerConfigForEveryDeclaredRegistry(t *testing.T) {
	data, err := BuildImagePullSecretData(&kubernetesv1.WorkloadPod{
		ImageRegistries: []*kubernetesv1.WorkloadImageRegistry{
			{Server: "ghcr.io", Username: "acme-pull-bot", Password: "resolved-token"},
			{Server: "quay.io", Username: "acme+robot", Password: "robot-token", Email: "ops@acme.io"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	doc, ok := data[ImagePullSecretDataKey]
	if !ok || len(data) != 1 {
		t.Fatalf("want exactly the %s key, got %v", ImagePullSecretDataKey, data)
	}
	for _, want := range []string{`"ghcr.io"`, `"quay.io"`, `"acme-pull-bot"`, `"resolved-token"`, `"ops@acme.io"`} {
		if !strings.Contains(doc, want) {
			t.Fatalf("document must carry %s, got %s", want, doc)
		}
	}
}

func TestBuildImagePullSecretData_RefusesADuplicateServer(t *testing.T) {
	_, err := BuildImagePullSecretData(&kubernetesv1.WorkloadPod{
		ImageRegistries: []*kubernetesv1.WorkloadImageRegistry{
			{Server: "ghcr.io", Username: "a", Password: "p"},
			{Server: "ghcr.io", Username: "b", Password: "p"},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "one login per registry") {
		t.Fatalf("a duplicate server must be refused with the sentence, got: %v", err)
	}
}
